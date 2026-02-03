package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type RaceRepository interface {
	FindAll() ([]*models.Race, error)
	FindByID(id uuid.UUID) (*models.Race, error)
	FindByIDWithTraits(id uuid.UUID) (*models.Race, error)
	FindByName(name string) (*models.Race, error)
}

type RaceRepo struct {
	db *gorm.DB
}

func NewRaceRepo(db *gorm.DB) *RaceRepo {
	return &RaceRepo{db}
}

func (r *RaceRepo) FindAll() ([]*models.Race, error) {
	var races []*models.Race
	err := r.db.Find(&races).Error
	return races, err
}

func (r *RaceRepo) FindByID(id uuid.UUID) (*models.Race, error) {
	var race models.Race
	if err := r.db.First(&race, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &race, nil
}

// FindByIDWithTraits returns a race with Traits and TraitOptions preloaded (for compendium detail view)
func (r *RaceRepo) FindByIDWithTraits(id uuid.UUID) (*models.Race, error) {
	var race models.Race
	if err := r.db.Preload("Traits.Options").First(&race, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &race, nil
}

func (r *RaceRepo) FindByName(name string) (*models.Race, error) {
	var race models.Race
	if err := r.db.First(&race, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &race, nil
}
