package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type EffectRepository interface {
	FindAll() ([]*models.Effect, error)
	FindByID(id uuid.UUID) (*models.Effect, error)
	FindByName(name string) (*models.Effect, error)
}

type EffectRepo struct {
	db *gorm.DB
}

func NewEffectRepo(db *gorm.DB) *EffectRepo {
	return &EffectRepo{db}
}

func (r *EffectRepo) FindAll() ([]*models.Effect, error) {
	var effects []*models.Effect
	// Order by Name
	err := r.db.Order("name asc").Find(&effects).Error
	return effects, err
}

func (r *EffectRepo) FindByID(id uuid.UUID) (*models.Effect, error) {
	var effect models.Effect
	if err := r.db.First(&effect, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &effect, nil
}

func (r *EffectRepo) FindByName(name string) (*models.Effect, error) {
	var effect models.Effect
	if err := r.db.First(&effect, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &effect, nil
}
