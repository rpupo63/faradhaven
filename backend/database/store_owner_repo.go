package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// StoreOwnerRepository loads vendors and catalog rules for the shop API.
type StoreOwnerRepository interface {
	FindAllWithRules() ([]models.StoreOwner, error)
	FindByIDWithRules(id uuid.UUID) (*models.StoreOwner, error)
}

// StoreOwnerRepo implements StoreOwnerRepository.
type StoreOwnerRepo struct {
	db *gorm.DB
}

// NewStoreOwnerRepo creates a StoreOwnerRepo.
func NewStoreOwnerRepo(db *gorm.DB) *StoreOwnerRepo {
	return &StoreOwnerRepo{db: db}
}

// FindAllWithRules returns all store owners ordered by name with catalog rules preloaded.
func (r *StoreOwnerRepo) FindAllWithRules() ([]models.StoreOwner, error) {
	var owners []models.StoreOwner
	err := r.db.Preload("CatalogRules").Order("name ASC").Find(&owners).Error
	return owners, err
}

// FindByIDWithRules returns one store owner with rules or gorm.ErrRecordNotFound.
func (r *StoreOwnerRepo) FindByIDWithRules(id uuid.UUID) (*models.StoreOwner, error) {
	var owner models.StoreOwner
	err := r.db.Preload("CatalogRules").First(&owner, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &owner, nil
}
