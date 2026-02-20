package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type ElevationPropertiesRepo struct {
	db *gorm.DB
}

func NewElevationPropertiesRepo(db *gorm.DB) *ElevationPropertiesRepo {
	return &ElevationPropertiesRepo{db: db}
}

func (r *ElevationPropertiesRepo) Create(props *models.ElevationProperties) error {
	return r.db.Create(props).Error
}

func (r *ElevationPropertiesRepo) GetByID(id uuid.UUID) (*models.ElevationProperties, error) {
	var props models.ElevationProperties
	if err := r.db.First(&props, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &props, nil
}

func (r *ElevationPropertiesRepo) Update(props *models.ElevationProperties) error {
	return r.db.Save(props).Error
}

func (r *ElevationPropertiesRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ElevationProperties{}, "id = ?", id).Error
}
