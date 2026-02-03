package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type BeastRepository interface {
	FindAll() ([]*models.Beast, error)
	FindByID(id uuid.UUID) (*models.Beast, error)
	FindByIDWithAttacks(id uuid.UUID) (*models.Beast, error)
	FindByUserID(userID uuid.UUID) ([]*models.Beast, error)
	FindByUserIDWithAttacks(userID uuid.UUID) ([]*models.Beast, error)
	Add(beast *models.Beast) error
	Update(beast *models.Beast) error
	Delete(id uuid.UUID) error
}

type BeastRepo struct {
	db *gorm.DB
}

func NewBeastRepo(db *gorm.DB) *BeastRepo {
	return &BeastRepo{db}
}

// FindAll returns all beasts
func (r *BeastRepo) FindAll() ([]*models.Beast, error) {
	var beasts []*models.Beast
	err := r.db.Find(&beasts).Error
	return beasts, err
}

// FindByID returns a beast by ID
func (r *BeastRepo) FindByID(id uuid.UUID) (*models.Beast, error) {
	var beast models.Beast
	err := r.db.First(&beast, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &beast, nil
}

// FindByIDWithAttacks returns a beast by ID with its attacks
func (r *BeastRepo) FindByIDWithAttacks(id uuid.UUID) (*models.Beast, error) {
	var beast models.Beast
	err := r.db.Preload("Attacks").First(&beast, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &beast, nil
}

// FindByUserID returns all beasts for a user
func (r *BeastRepo) FindByUserID(userID uuid.UUID) ([]*models.Beast, error) {
	var beasts []*models.Beast
	err := r.db.Where("user_id = ?", userID).Find(&beasts).Error
	return beasts, err
}

// FindByUserIDWithAttacks returns all beasts for a user with their attacks
func (r *BeastRepo) FindByUserIDWithAttacks(userID uuid.UUID) ([]*models.Beast, error) {
	var beasts []*models.Beast
	err := r.db.Preload("Attacks").Where("user_id = ?", userID).Find(&beasts).Error
	return beasts, err
}

// Add inserts a new beast with its attacks
func (r *BeastRepo) Add(beast *models.Beast) error {
	now := time.Now()
	beast.CreatedAt = now
	beast.UpdatedAt = now
	for i := range beast.Attacks {
		beast.Attacks[i].CreatedAt = now
	}
	return r.db.Create(beast).Error
}

// Update updates an existing beast
func (r *BeastRepo) Update(beast *models.Beast) error {
	beast.UpdatedAt = time.Now()
	return r.db.Save(beast).Error
}

// Delete removes a beast by ID (attacks are deleted via cascade)
func (r *BeastRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Beast{}, "id = ?", id).Error
}
