package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type SavedSpellRepository interface {
	FindByCharacterID(characterID uuid.UUID) ([]*models.SavedSpell, error)
	FindByCharacterAndSlot(characterID uuid.UUID, slotIndex int) (*models.SavedSpell, error)
	FindByID(id uuid.UUID) (*models.SavedSpell, error)
	Add(spell *models.SavedSpell) error
	Update(spell *models.SavedSpell) error
	Delete(id uuid.UUID) error
	DeleteByCharacterID(characterID uuid.UUID) error
}

type SavedSpellRepo struct {
	db *gorm.DB
}

func NewSavedSpellRepo(db *gorm.DB) *SavedSpellRepo {
	return &SavedSpellRepo{db}
}

// FindByCharacterID returns all saved spells for a character
func (r *SavedSpellRepo) FindByCharacterID(characterID uuid.UUID) ([]*models.SavedSpell, error) {
	var spells []*models.SavedSpell
	err := r.db.Where("character_id = ?", characterID).Order("slot_index asc").Find(&spells).Error
	return spells, err
}

// FindByCharacterAndSlot returns a saved spell in a specific slot
func (r *SavedSpellRepo) FindByCharacterAndSlot(characterID uuid.UUID, slotIndex int) (*models.SavedSpell, error) {
	var spell models.SavedSpell
	err := r.db.Where("character_id = ? AND slot_index = ?", characterID, slotIndex).First(&spell).Error
	if err != nil {
		return nil, err
	}
	return &spell, nil
}

// FindByID returns a single saved spell by ID
func (r *SavedSpellRepo) FindByID(id uuid.UUID) (*models.SavedSpell, error) {
	var spell models.SavedSpell
	err := r.db.First(&spell, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &spell, nil
}

// Add creates a new saved spell
func (r *SavedSpellRepo) Add(spell *models.SavedSpell) error {
	now := time.Now()
	spell.CreatedAt = now
	spell.UpdatedAt = now
	return r.db.Create(spell).Error
}

// Update updates an existing saved spell
func (r *SavedSpellRepo) Update(spell *models.SavedSpell) error {
	spell.UpdatedAt = time.Now()
	return r.db.Save(spell).Error
}

// Delete removes a saved spell by ID
func (r *SavedSpellRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.SavedSpell{}, "id = ?", id).Error
}

// DeleteByCharacterID removes all saved spells for a character
func (r *SavedSpellRepo) DeleteByCharacterID(characterID uuid.UUID) error {
	return r.db.Where("character_id = ?", characterID).Delete(&models.SavedSpell{}).Error
}

// ResetUsesOnRest resets uses_remaining for all spells that restore on rest
func (r *SavedSpellRepo) ResetUsesOnRest(characterID uuid.UUID) error {
	return r.db.Model(&models.SavedSpell{}).
		Where("character_id = ? AND max_uses IS NOT NULL", characterID).
		Update("uses_remaining", gorm.Expr("max_uses")).Error
}
