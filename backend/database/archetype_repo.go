package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type ArchetypeRepository interface {
	FindByClassID(classID uuid.UUID) ([]*models.Archetype, error)
	FindByID(id uuid.UUID) (*models.Archetype, error)
}

type ArchetypeRepo struct {
	db *gorm.DB
}

func NewArchetypeRepo(db *gorm.DB) *ArchetypeRepo {
	return &ArchetypeRepo{db}
}

// FindByClassID returns all archetypes for a class, ordered by SortOrder
func (r *ArchetypeRepo) FindByClassID(classID uuid.UUID) ([]*models.Archetype, error) {
	var archetypes []*models.Archetype
	err := r.db.Where("class_id = ?", classID).Order("sort_order ASC").Find(&archetypes).Error
	return archetypes, err
}

// FindByID returns an archetype by ID
func (r *ArchetypeRepo) FindByID(id uuid.UUID) (*models.Archetype, error) {
	var archetype models.Archetype
	err := r.db.First(&archetype, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &archetype, nil
}
