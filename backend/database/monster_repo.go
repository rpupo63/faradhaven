package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// MonsterRepository defines the interface for interacting with Monster data.
type MonsterRepository interface {
	Add(monster *models.Monster) error
	FindByID(id uuid.UUID) (*models.Monster, error)
	FindByUserID(userID uuid.UUID) ([]*models.Monster, error)
	FindAll() ([]*models.Monster, error)
	Update(monster *models.Monster) error
	Delete(id uuid.UUID) error
}

// MonsterRepo implements MonsterRepository using GORM.
type MonsterRepo struct {
	db *gorm.DB
}

// NewMonsterRepo creates a new MonsterRepo.
func NewMonsterRepo(db *gorm.DB) *MonsterRepo {
	return &MonsterRepo{db: db}
}

// Add creates a new Monster record.
func (r *MonsterRepo) Add(monster *models.Monster) error {
	return r.db.Create(monster).Error
}

// FindByID retrieves a Monster by its ID, preloading relations.
func (r *MonsterRepo) FindByID(id uuid.UUID) (*models.Monster, error) {
	var monster models.Monster
	err := r.db.Preload("Attacks.Weapon").Preload("Actions").First(&monster, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to find monster by ID: %w", err)
	}
	return &monster, nil
}

// FindByUserID retrieves all Monsters belonging to a specific user, preloading relations.
func (r *MonsterRepo) FindByUserID(userID uuid.UUID) ([]*models.Monster, error) {
	var monsters []*models.Monster
	err := r.db.Preload("Attacks.Weapon").Preload("Actions").Where("user_id = ?", userID).Find(&monsters).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find monsters by user ID: %w", err)
	}
	return monsters, nil
}

// FindAll retrieves all Monsters, preloading relations. (Admin/Debug only, or for system-wide monsters)
func (r *MonsterRepo) FindAll() ([]*models.Monster, error) {
	var monsters []*models.Monster
	err := r.db.Preload("Attacks.Weapon").Preload("Actions").Find(&monsters).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all monsters: %w", err)
	}
	return monsters, nil
}

// Update updates an existing Monster record.
func (r *MonsterRepo) Update(monster *models.Monster) error {
	// Use Save to handle both creation and updates, or selectively update fields.
	// For full object update:
	return r.db.Save(monster).Error
	// If only specific fields are updated, use r.db.Model(&monster).Updates(...)
}

// Delete removes a Monster record by its ID.
func (r *MonsterRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Monster{}, id).Error
}
