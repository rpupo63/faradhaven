package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type WallPropertiesRepo struct {
	db *gorm.DB
}

func NewWallPropertiesRepo(db *gorm.DB) *WallPropertiesRepo {
	return &WallPropertiesRepo{db: db}
}

func (r *WallPropertiesRepo) Create(props *models.WallProperties) error {
	return r.db.Create(props).Error
}

func (r *WallPropertiesRepo) GetByID(id uuid.UUID) (*models.WallProperties, error) {
	var props models.WallProperties
	if err := r.db.First(&props, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &props, nil
}

func (r *WallPropertiesRepo) Update(props *models.WallProperties) error {
	return r.db.Save(props).Error
}

func (r *WallPropertiesRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.WallProperties{}, "id = ?", id).Error
}
