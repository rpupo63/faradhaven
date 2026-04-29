package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type TrapPropertiesRepo struct {
	db *gorm.DB
}

func NewTrapPropertiesRepo(db *gorm.DB) *TrapPropertiesRepo {
	return &TrapPropertiesRepo{db: db}
}

func (r *TrapPropertiesRepo) Create(props *models.TrapProperties) error {
	return r.db.Create(props).Error
}

func (r *TrapPropertiesRepo) GetByID(id uuid.UUID) (*models.TrapProperties, error) {
	var props models.TrapProperties
	if err := r.db.First(&props, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &props, nil
}

func (r *TrapPropertiesRepo) Update(props *models.TrapProperties) error {
	return r.db.Save(props).Error
}

func (r *TrapPropertiesRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.TrapProperties{}, "id = ?", id).Error
}
