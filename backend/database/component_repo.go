package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type ComponentRepo struct {
	db *gorm.DB
}

func NewComponentRepo(db *gorm.DB) *ComponentRepo {
	return &ComponentRepo{db: db}
}

// GetAllComponents returns all components ordered by category then name.
func (r *ComponentRepo) GetAllComponents() ([]models.Component, error) {
	var components []models.Component
	err := r.db.Order("category ASC, name ASC").Find(&components).Error
	return components, err
}

// GetComponentByID returns a single component by ID.
func (r *ComponentRepo) GetComponentByID(id uuid.UUID) (*models.Component, error) {
	var component models.Component
	err := r.db.First(&component, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &component, nil
}

// GetComponentsByCategory returns all components in a specific category.
func (r *ComponentRepo) GetComponentsByCategory(category models.ComponentCategory) ([]models.Component, error) {
	var components []models.Component
	err := r.db.Where("category = ?", category).Order("name ASC").Find(&components).Error
	return components, err
}
