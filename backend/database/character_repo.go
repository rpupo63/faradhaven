package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type CharacterRepository interface {
	FindAll() ([]*models.Character, error)
	FindAllPaginated(page, limit int) ([]*models.Character, int64, error)
	FindByID(id uuid.UUID) (*models.Character, error)
	FindByIDWithSkills(id uuid.UUID) (*models.Character, error)
	FindByIDWithRelations(id uuid.UUID) (*models.Character, error)
	FindByIDForSheet(id uuid.UUID) (*models.Character, error)
	FindByIDWithInventory(id uuid.UUID) (*models.Character, error)
	FindByUserID(userID uuid.UUID) ([]*models.Character, error)
	FindByUserIDPaginated(userID uuid.UUID, page, limit int) ([]*models.Character, int64, error)
	Add(character *models.Character) error
	Update(character *models.Character) error
	Delete(id uuid.UUID) error
	CharacterBelongsToUser(characterID, userID uuid.UUID) (bool, error)
	ReplaceSkillProficiencies(characterID uuid.UUID, skillIDs []string) error
	UpdateComponentCount(characterID uuid.UUID, componentID uuid.UUID, delta int) error
	ClearComponentsForCharacter(characterID uuid.UUID) error
	AppendWeapon(characterID uuid.UUID, weaponID uuid.UUID) error
	AppendItem(characterID uuid.UUID, itemID uuid.UUID) error
	UpdateMoney(id uuid.UUID, money int64) error
	UpdateHP(id uuid.UUID, hp int) error
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
	err := r.db.Preload("Race").Preload("Class").Preload("Archetype").
		Order("name ASC").Find(&characters).Error
	return characters, err
}

// FindAllPaginated returns all characters with pagination
func (r *CharacterRepo) FindAllPaginated(page, limit int) ([]*models.Character, int64, error) {
	var characters []*models.Character
	var totalCount int64

	query := r.db.Model(&models.Character{})
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Race").Preload("Class").Preload("Archetype").
		Order("name ASC").Offset(offset).Limit(limit).Find(&characters).Error
	return characters, totalCount, err
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
// SkillProficiencyIDs populated, and Race.Traits.Options preloaded
func (r *CharacterRepo) FindByIDWithSkills(id uuid.UUID) (*models.Character, error) {
	var character models.Character
	if err := r.db.Preload("Race").Preload("Class").Preload("Archetype").
		Preload("Race.Traits.Options").Preload("Components.Component").
		Preload("SkillProficiencies").
		First(&character, "id = ?", id).Error; err != nil {
		return nil, err
	}
	for _, s := range character.SkillProficiencies {
		if s.Proficient {
			character.SkillProficiencyIDs = append(character.SkillProficiencyIDs, s.SkillID)
		}
	}
	return &character, nil
}

// FindByIDWithRelations returns a character by ID with relations needed for gameplay logic
func (r *CharacterRepo) FindByIDWithRelations(id uuid.UUID) (*models.Character, error) {
	var character models.Character
	if err := r.db.
		Preload("Race").
		Preload("Race.Traits.Options").
		Preload("Race.Components").
		Preload("Class").
		Preload("Class.Components").
		Preload("Class.Levels", func(db *gorm.DB) *gorm.DB {
			return db.Order("level ASC")
		}).
		Preload("Archetype").
		Preload("Party").
		Preload("Components.Component").
		Preload("SkillProficiencies").
		First(&character, "id = ?", id).Error; err != nil {
		return nil, err
	}
	for _, s := range character.SkillProficiencies {
		if s.Proficient {
			character.SkillProficiencyIDs = append(character.SkillProficiencyIDs, s.SkillID)
		}
	}
	return &character, nil
}

// FindByIDForSheet returns a character fully loaded for the character sheet
func (r *CharacterRepo) FindByIDForSheet(id uuid.UUID) (*models.Character, error) {
	var character models.Character
	if err := r.db.
		Preload("Race").
		Preload("Race.Traits.Options").
		Preload("Race.Components").
		Preload("Class").
		Preload("Class.Components").
		Preload("Class.Levels", func(db *gorm.DB) *gorm.DB {
			return db.Order("level ASC")
		}).
		Preload("Class.Levels.LevelFeatures", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Archetype").
		Preload("Party").
		Preload("Components.Component").
		Preload("CharacterWeapons.Weapon.Damages").
		Preload("CharacterWeapons.Modifiers").
		Preload("Items").
		Preload("EquippedArmor").
		Preload("EquippedShield").
		Preload("SkillProficiencies").
		First(&character, "id = ?", id).Error; err != nil {
		return nil, err
	}
	for _, s := range character.SkillProficiencies {
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

// FindByUserIDPaginated returns all characters for a user with pagination
func (r *CharacterRepo) FindByUserIDPaginated(userID uuid.UUID, page, limit int) ([]*models.Character, int64, error) {
	var characters []*models.Character
	var totalCount int64

	query := r.db.Model(&models.Character{}).Where("user_id = ?", userID)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Race").Preload("Class").Preload("Archetype").
		Order("updated_at DESC").Offset(offset).Limit(limit).Find(&characters).Error
	return characters, totalCount, err
}

// CharacterBelongsToUser returns true if the character belongs to the given user
func (r *CharacterRepo) CharacterBelongsToUser(characterID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Character{}).Where("id = ? AND user_id = ?", characterID, userID).Count(&count).Error
	return count > 0, err
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

		// Insert new proficiencies in a single batch
		skills := make([]models.CharacterSkill, len(skillIDs))
		for i, skillID := range skillIDs {
			skills[i] = models.CharacterSkill{
				CharacterID: characterID,
				SkillID:     skillID,
				Proficient:  true,
			}
		}
		if len(skills) > 0 {
			if err := tx.Create(&skills).Error; err != nil {
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

// AppendWeapon adds a weapon association to a character
func (r *CharacterRepo) AppendWeapon(characterID uuid.UUID, weaponID uuid.UUID) error {
	// For weapons, we usually create a CharacterWeapon entry because it has is_equipped etc.
	cw := models.CharacterWeapon{
		CharacterID: characterID,
		WeaponID:    weaponID,
		IsEquipped:  false,
	}
	return r.db.Create(&cw).Error
}

// AppendItem adds an item association to a character
func (r *CharacterRepo) AppendItem(characterID uuid.UUID, itemID uuid.UUID) error {
	return r.db.Exec("INSERT INTO character_items (character_id, item_id) VALUES (?, ?)", characterID, itemID).Error
}

// UpdateMoney updates only the money field for a character
func (r *CharacterRepo) UpdateMoney(id uuid.UUID, money int64) error {
	return r.db.Model(&models.Character{}).Where("id = ?", id).Update("money", money).Error
}

// UpdateHP updates only the current_hp field for a character
func (r *CharacterRepo) UpdateHP(id uuid.UUID, hp int) error {
	return r.db.Model(&models.Character{}).Where("id = ?", id).Update("current_hp", hp).Error
}
