package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type WeaponRepo struct {
	db *gorm.DB
}

func NewWeaponRepo(db *gorm.DB) *WeaponRepo {
	return &WeaponRepo{db: db}
}

// GetAllWeapons returns all weapons, including their damage details.
func (r *WeaponRepo) GetAllWeapons() ([]models.Weapon, error) {
	var weapons []models.Weapon
	// Preload the Damages relationship
	err := r.db.Preload("Damages").Order("name ASC").Find(&weapons).Error
	return weapons, err
}

// GetWeaponByID returns a single weapon by ID with details.
func (r *WeaponRepo) GetWeaponByID(id uuid.UUID) (*models.Weapon, error) {
	var weapon models.Weapon
	err := r.db.Preload("Damages").First(&weapon, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &weapon, nil
}
