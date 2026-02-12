package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type CharacterEffectRepository interface {
	FindByCharacterID(characterID uuid.UUID) ([]*models.CharacterEffect, error)
	FindByID(id uuid.UUID) (*models.CharacterEffect, error)
	Add(effect *models.CharacterEffect) error
	Update(effect *models.CharacterEffect) error
	Delete(id uuid.UUID) error
	DeleteByCharacterID(characterID uuid.UUID) error
	FindConcentrationEffects(characterID uuid.UUID) ([]*models.CharacterEffect, error)
	FindBySourceCharacter(sourceCharacterID uuid.UUID) ([]*models.CharacterEffect, error)
}

type CharacterEffectRepo struct {
	db *gorm.DB
}

func NewCharacterEffectRepo(db *gorm.DB) *CharacterEffectRepo {
	return &CharacterEffectRepo{db}
}

// FindByCharacterID returns all active effects for a character
func (r *CharacterEffectRepo) FindByCharacterID(characterID uuid.UUID) ([]*models.CharacterEffect, error) {
	var effects []*models.CharacterEffect
	err := r.db.Preload("Effect").Where("character_id = ?", characterID).Find(&effects).Error
	return effects, err
}

// FindByID returns a single character effect by ID
func (r *CharacterEffectRepo) FindByID(id uuid.UUID) (*models.CharacterEffect, error) {
	var effect models.CharacterEffect
	err := r.db.Preload("Effect").First(&effect, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &effect, nil
}

// Add inserts a new character effect
func (r *CharacterEffectRepo) Add(effect *models.CharacterEffect) error {
	now := time.Now()
	effect.CreatedAt = now
	effect.UpdatedAt = now
	return r.db.Create(effect).Error
}

// Update updates an existing character effect
func (r *CharacterEffectRepo) Update(effect *models.CharacterEffect) error {
	effect.UpdatedAt = time.Now()
	return r.db.Save(effect).Error
}

// Delete removes a character effect by ID
func (r *CharacterEffectRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.CharacterEffect{}, "id = ?", id).Error
}

// DeleteByCharacterID removes all effects from a character
func (r *CharacterEffectRepo) DeleteByCharacterID(characterID uuid.UUID) error {
	return r.db.Where("character_id = ?", characterID).Delete(&models.CharacterEffect{}).Error
}

// FindConcentrationEffects returns all concentration effects for a character
func (r *CharacterEffectRepo) FindConcentrationEffects(characterID uuid.UUID) ([]*models.CharacterEffect, error) {
	var effects []*models.CharacterEffect
	err := r.db.Preload("Effect").Where("character_id = ? AND is_concentration = true", characterID).Find(&effects).Error
	return effects, err
}

// FindBySourceCharacter returns all effects applied by a specific character
func (r *CharacterEffectRepo) FindBySourceCharacter(sourceCharacterID uuid.UUID) ([]*models.CharacterEffect, error) {
	var effects []*models.CharacterEffect
	err := r.db.Preload("Effect").Where("source_character_id = ?", sourceCharacterID).Find(&effects).Error
	return effects, err
}

// FindByEffectAndCharacter finds an existing effect instance on a character
func (r *CharacterEffectRepo) FindByEffectAndCharacter(characterID, effectID uuid.UUID) (*models.CharacterEffect, error) {
	var effect models.CharacterEffect
	err := r.db.Preload("Effect").Where("character_id = ? AND effect_id = ?", characterID, effectID).First(&effect).Error
	if err != nil {
		return nil, err
	}
	return &effect, nil
}

// DecrementDurationRounds decrements duration_rounds for all effects with duration and returns expired ones
func (r *CharacterEffectRepo) DecrementDurationRounds(characterID uuid.UUID, rounds int) ([]*models.CharacterEffect, error) {
	// First, decrement all durations
	err := r.db.Model(&models.CharacterEffect{}).
		Where("character_id = ? AND duration_rounds IS NOT NULL AND duration_rounds > 0", characterID).
		Update("duration_rounds", gorm.Expr("duration_rounds - ?", rounds)).Error
	if err != nil {
		return nil, err
	}

	// Find and return expired effects (duration_rounds <= 0)
	var expired []*models.CharacterEffect
	err = r.db.Preload("Effect").
		Where("character_id = ? AND duration_rounds IS NOT NULL AND duration_rounds <= 0", characterID).
		Find(&expired).Error
	return expired, err
}
