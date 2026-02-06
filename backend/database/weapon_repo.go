package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type WeaponRepository interface {
	FindAll() ([]models.Weapon, error)
	FindByID(id uuid.UUID) (*models.Weapon, error)
	FindByCategories(categories []string) ([]models.Weapon, error)
	FindByIDs(ids []uuid.UUID) ([]models.Weapon, error)
}

type WeaponRepo struct {
	db *gorm.DB
}

func NewWeaponRepo(db *gorm.DB) *WeaponRepo {
	return &WeaponRepo{db: db}
}

// FindAll returns all weapons, including their damage details.
func (r *WeaponRepo) FindAll() ([]models.Weapon, error) {
	var weapons []models.Weapon
	// Preload the Damages relationship
	err := r.db.Preload("Damages").Order("name ASC").Find(&weapons).Error
	return weapons, err
}

// FindByID returns a single weapon by ID with details.
func (r *WeaponRepo) FindByID(id uuid.UUID) (*models.Weapon, error) {
	var weapon models.Weapon
	err := r.db.Preload("Damages").First(&weapon, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &weapon, nil
}

// FindByCategories returns weapons matching any of the given categories.
// If categories is empty, returns all weapons.
func (r *WeaponRepo) FindByCategories(categories []string) ([]models.Weapon, error) {
	var weapons []models.Weapon
	query := r.db.Preload("Damages")
	if len(categories) > 0 {
		query = query.Where("category IN ?", categories)
	}
	err := query.Order("name ASC").Find(&weapons).Error
	return weapons, err
}

// FindByIDs returns weapons by their IDs.
func (r *WeaponRepo) FindByIDs(ids []uuid.UUID) ([]models.Weapon, error) {
	var weapons []models.Weapon
	if len(ids) == 0 {
		return weapons, nil
	}
	err := r.db.Preload("Damages").Where("id IN ?", ids).Find(&weapons).Error
	return weapons, err
}
