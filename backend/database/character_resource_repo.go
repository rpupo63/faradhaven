package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type CharacterResourceRepository interface {
	FindByCharacterID(characterID uuid.UUID) ([]*models.CharacterResource, error)
	FindByCharacterAndKey(characterID uuid.UUID, key string) (*models.CharacterResource, error)
	Add(resource *models.CharacterResource) error
	Update(resource *models.CharacterResource) error
	Delete(id uuid.UUID) error
	DeleteByCharacterID(characterID uuid.UUID) error
	UpdateCurrentValue(characterID uuid.UUID, resourceKey string, delta int) error
}

type CharacterResourceRepo struct {
	db *gorm.DB
}

func NewCharacterResourceRepo(db *gorm.DB) *CharacterResourceRepo {
	return &CharacterResourceRepo{db}
}

// FindByCharacterID returns all resources for a character
func (r *CharacterResourceRepo) FindByCharacterID(characterID uuid.UUID) ([]*models.CharacterResource, error) {
	var resources []*models.CharacterResource
	err := r.db.Where("character_id = ?", characterID).Find(&resources).Error
	return resources, err
}

// FindByCharacterAndKey returns a specific resource by character and key
func (r *CharacterResourceRepo) FindByCharacterAndKey(characterID uuid.UUID, key string) (*models.CharacterResource, error) {
	var resource models.CharacterResource
	err := r.db.Where("character_id = ? AND resource_key = ?", characterID, key).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Add creates a new character resource
func (r *CharacterResourceRepo) Add(resource *models.CharacterResource) error {
	now := time.Now()
	resource.CreatedAt = now
	resource.UpdatedAt = now
	return r.db.Create(resource).Error
}

// Update updates an existing character resource
func (r *CharacterResourceRepo) Update(resource *models.CharacterResource) error {
	resource.UpdatedAt = time.Now()
	return r.db.Save(resource).Error
}

// Delete removes a character resource by ID
func (r *CharacterResourceRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.CharacterResource{}, "id = ?", id).Error
}

// DeleteByCharacterID removes all resources for a character
func (r *CharacterResourceRepo) DeleteByCharacterID(characterID uuid.UUID) error {
	return r.db.Where("character_id = ?", characterID).Delete(&models.CharacterResource{}).Error
}

// UpdateCurrentValue updates the CurrentValue of a CharacterResource, ensuring it stays within bounds.
func (r *CharacterResourceRepo) UpdateCurrentValue(characterID uuid.UUID, resourceKey string, delta int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		resource, err := r.FindByCharacterAndKey(characterID, resourceKey)
		if err != nil {
			return err
		}

		resource.CurrentValue += delta

		// Ensure CurrentValue doesn't go below 0
		if resource.CurrentValue < 0 {
			resource.CurrentValue = 0
		}

		// Ensure CurrentValue doesn't exceed MaxValue if it exists
		if resource.MaxValue != nil && resource.CurrentValue > *resource.MaxValue {
			resource.CurrentValue = *resource.MaxValue
		}

		return tx.Save(resource).Error
	})
}

// ProcessRestoration handles resource restoration during rests
func (r *CharacterResourceRepo) ProcessRestoration(characterID uuid.UUID, isLongRest bool) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var resources []*models.CharacterResource
		err := tx.Where("character_id = ?", characterID).Find(&resources).Error
		if err != nil {
			return err
		}

		for _, res := range resources {
			// Handle decay
			if isLongRest && res.DecaysOnLongRest {
				res.CurrentValue = 0
			}

			// Handle restoration
			shouldRestore := (isLongRest && res.RestoreOnLongRest) || (!isLongRest && res.RestoreOnShortRest)
			if shouldRestore && res.MaxValue != nil {
				if res.RestoreAmount != nil {
					res.CurrentValue += *res.RestoreAmount
					if res.CurrentValue > *res.MaxValue {
						res.CurrentValue = *res.MaxValue
					}
				} else {
					res.CurrentValue = *res.MaxValue
				}
			}

			res.UpdatedAt = time.Now()
			if err := tx.Save(res).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
