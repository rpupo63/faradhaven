package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type CharacterRepository interface {
	FindAll() ([]*models.Character, error)
	FindByID(id uuid.UUID) (*models.Character, error)
	FindByIDWithSkills(id uuid.UUID) (*models.Character, error)
	FindByIDWithRelations(id uuid.UUID) (*models.Character, error)
	FindByIDWithInventory(id uuid.UUID) (*models.Character, error)
	FindByUserID(userID uuid.UUID) ([]*models.Character, error)
	Add(character *models.Character) error
	Update(character *models.Character) error
	Delete(id uuid.UUID) error
	ReplaceSkillProficiencies(characterID uuid.UUID, skillIDs []string) error
	UpdateComponentCount(characterID uuid.UUID, componentID uuid.UUID, delta int) error
	ClearComponentsForCharacter(characterID uuid.UUID) error
	GetDB() *gorm.DB
}

type CharacterRepo struct {
	db *gorm.DB
}

func NewCharacterRepo(db *gorm.DB) *CharacterRepo {
	return &CharacterRepo{db}
}

func (r *CharacterRepo) GetDB() *gorm.DB {
	return r.db
}

// FindAll returns all characters
func (r *CharacterRepo) FindAll() ([]*models.Character, error) {
	var characters []*models.Character
	err := r.db.Preload("Race").Preload("Class").Preload("Archetype").Preload("Components.Component").Find(&characters).Error
	return characters, err
}

// FindByID returns a character by ID
func (r *CharacterRepo) FindByID(id uuid.UUID) (*models.Character, error) {
	var character models.Character
	err := r.db.Preload("Race").Preload("Class").Preload("Archetype").Preload("Components.Component").First(&character, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &character, nil
}

// FindByIDWithSkills returns a character by ID with SkillProficiencies preloaded,
// SkillProficiencyIDs populated, and Race.Traits.Options preloaded (for character sheet racial traits)
func (r *CharacterRepo) FindByIDWithSkills(id uuid.UUID) (*models.Character, error) {
	var character models.Character
	if err := r.db.Preload("Race").Preload("Class").Preload("Archetype").Preload("Race.Traits.Options").Preload("Components.Component").First(&character, "id = ?", id).Error; err != nil {
		return nil, err
	}
	var skills []models.CharacterSkill
	if err := r.db.Where("character_id = ?", id).Find(&skills).Error; err != nil {
		return nil, err
	}
	for _, s := range skills {
		if s.Proficient {
			character.SkillProficiencyIDs = append(character.SkillProficiencyIDs, s.SkillID)
		}
	}
	return &character, nil
}

// FindByIDWithRelations returns a character by ID with all necessary relations for Lorewright features
func (r *CharacterRepo) FindByIDWithRelations(id uuid.UUID) (*models.Character, error) {
	var character models.Character
	if err := r.db.
		Preload("Race").
		Preload("Class").
		Preload("Class.Levels"). // Crucial for Lorewright slot checks
		Preload("Archetype").
		Preload("Party"). // NEW: Preload Party
		Preload("Race.Traits.Options").
		Preload("Components.Component").
		First(&character, "id = ?", id).Error; err != nil {
		return nil, err
	}
	// Manually populate SkillProficiencyIDs as it's a computed field
	var skills []models.CharacterSkill
	if err := r.db.Where("character_id = ?", id).Find(&skills).Error; err != nil {
		return nil, err
	}
	for _, s := range skills {
		if s.Proficient {
			character.SkillProficiencyIDs = append(character.SkillProficiencyIDs, s.SkillID)
		}
	}
	return &character, nil
}

// FindByUserID returns all characters for a user
func (r *CharacterRepo) FindByUserID(userID uuid.UUID) ([]*models.Character, error) {
	var characters []*models.Character
	err := r.db.Preload("Race").Preload("Class").Preload("Archetype").Preload("Components.Component").Where("user_id = ?", userID).Find(&characters).Error
	return characters, err
}

// FindByIDWithInventory returns a character by ID with Items and Weapons preloaded
func (r *CharacterRepo) FindByIDWithInventory(id uuid.UUID) (*models.Character, error) {
	var character models.Character
	err := r.db.Preload("Race").Preload("Class").Preload("Archetype").Preload("Items").Preload("Weapons").First(&character, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &character, nil
}

// Add inserts a new character
func (r *CharacterRepo) Add(character *models.Character) error {
	now := time.Now()
	character.CreatedAt = now
	character.UpdatedAt = now
	return r.db.Create(character).Error
}

// Update updates an existing character
func (r *CharacterRepo) Update(character *models.Character) error {
	character.UpdatedAt = time.Now()
	return r.db.Save(character).Error
}

// Delete removes a character by ID (cascades to CharacterSkills)
func (r *CharacterRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Character{}, "id = ?", id).Error
}

// ReplaceSkillProficiencies replaces all skill proficiencies for a character
func (r *CharacterRepo) ReplaceSkillProficiencies(characterID uuid.UUID, skillIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing proficiencies
		if err := tx.Where("character_id = ?", characterID).Delete(&models.CharacterSkill{}).Error; err != nil {
			return err
		}

		// Insert new proficiencies
		for _, skillID := range skillIDs {
			cs := models.CharacterSkill{
				CharacterID: characterID,
				SkillID:     skillID,
				Proficient:  true,
			}
			if err := tx.Create(&cs).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateComponentCount updates the count of a specific component for a character
func (r *CharacterRepo) UpdateComponentCount(characterID uuid.UUID, componentID uuid.UUID, delta int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var charComp models.CharacterComponent
		err := tx.Where("character_id = ? AND component_id = ?", characterID, componentID).First(&charComp).Error

		if err == gorm.ErrRecordNotFound {
			if delta <= 0 {
				return nil // Nothing to deduct
			}
			// Create new record if adding
			charComp = models.CharacterComponent{
				CharacterID: characterID,
				ComponentID: componentID,
				Count:       delta,
			}
			return tx.Create(&charComp).Error
		} else if err != nil {
			return err
		}

		charComp.Count += delta
		if charComp.Count < 0 {
			charComp.Count = 0
		}

		return tx.Save(&charComp).Error
	})
}

// ClearComponentsForCharacter removes all component inventory for a character
func (r *CharacterRepo) ClearComponentsForCharacter(characterID uuid.UUID) error {
	return r.db.Where("character_id = ?", characterID).Delete(&models.CharacterComponent{}).Error
}
