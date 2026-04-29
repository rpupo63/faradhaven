package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type DifficultTerrainPropertiesRepo struct {
	db *gorm.DB
}

func NewDifficultTerrainPropertiesRepo(db *gorm.DB) *DifficultTerrainPropertiesRepo {
	return &DifficultTerrainPropertiesRepo{db: db}
}

func (r *DifficultTerrainPropertiesRepo) Create(props *models.DifficultTerrainProperties) error {
	return r.db.Create(props).Error
}

func (r *DifficultTerrainPropertiesRepo) GetByID(id uuid.UUID) (*models.DifficultTerrainProperties, error) {
	var props models.DifficultTerrainProperties
	if err := r.db.First(&props, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &props, nil
}

func (r *DifficultTerrainPropertiesRepo) Update(props *models.DifficultTerrainProperties) error {
	return r.db.Save(props).Error
}

func (r *DifficultTerrainPropertiesRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.DifficultTerrainProperties{}, "id = ?", id).Error
}
