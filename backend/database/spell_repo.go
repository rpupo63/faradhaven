package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type SpellRepository interface {
	FindAll() ([]*models.Spell, error)
	FindByID(id uuid.UUID) (*models.Spell, error)
	FindByUserID(userID uuid.UUID) ([]*models.Spell, error)
	FindByCharacterID(characterID uuid.UUID) ([]*models.Spell, error)
	Add(spell *models.Spell, componentIDs []uuid.UUID) error
	Update(spell *models.Spell) error
	ReplaceComponents(spellID uuid.UUID, componentIDs []uuid.UUID) error
	Delete(id uuid.UUID) error
}

type SpellRepo struct {
	db *gorm.DB
}

func NewSpellRepo(db *gorm.DB) *SpellRepo {
	return &SpellRepo{db}
}

// FindAll returns all spells with components preloaded
func (r *SpellRepo) FindAll() ([]*models.Spell, error) {
	var spells []*models.Spell
	err := r.db.Preload("Components").Find(&spells).Error
	return spells, err
}

// FindByID returns a spell by ID with components preloaded
func (r *SpellRepo) FindByID(id uuid.UUID) (*models.Spell, error) {
	var spell models.Spell
	err := r.db.Preload("Components").First(&spell, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &spell, nil
}

// FindByUserID returns all spells for a user with components preloaded
func (r *SpellRepo) FindByUserID(userID uuid.UUID) ([]*models.Spell, error) {
	var spells []*models.Spell
	err := r.db.Preload("Components").Where("user_id = ?", userID).Find(&spells).Error
	return spells, err
}

// FindByCharacterID returns all spells prepared by a character with components preloaded
func (r *SpellRepo) FindByCharacterID(characterID uuid.UUID) ([]*models.Spell, error) {
	var spells []*models.Spell
	err := r.db.Preload("Components").Where("character_id = ?", characterID).Find(&spells).Error
	return spells, err
}

// Add inserts a new spell and links it to the given components
func (r *SpellRepo) Add(spell *models.Spell, componentIDs []uuid.UUID) error {
	now := time.Now()
	spell.CreatedAt = now
	spell.UpdatedAt = now
	if err := r.db.Create(spell).Error; err != nil {
		return err
	}
	if len(componentIDs) > 0 {
		return replaceComponentsTx(r.db, spell.ID, componentIDs)
	}
	return nil
}

// ReplaceComponents removes existing spell-component links and creates new ones
func (r *SpellRepo) ReplaceComponents(spellID uuid.UUID, componentIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("spell_id = ?", spellID).Delete(&models.SpellComponent{}).Error; err != nil {
			return err
		}
		if len(componentIDs) == 0 {
			return nil
		}
		return replaceComponentsTx(tx, spellID, componentIDs)
	})
}

func replaceComponentsTx(tx *gorm.DB, spellID uuid.UUID, componentIDs []uuid.UUID) error {
	for _, compID := range componentIDs {
		sc := models.SpellComponent{SpellID: spellID, ComponentID: compID}
		if err := tx.Create(&sc).Error; err != nil {
			return err
		}
	}
	return nil
}

// Update updates an existing spell
func (r *SpellRepo) Update(spell *models.Spell) error {
	spell.UpdatedAt = time.Now()
	return r.db.Save(spell).Error
}

// Delete removes a spell by ID
func (r *SpellRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Spell{}, "id = ?", id).Error
}
