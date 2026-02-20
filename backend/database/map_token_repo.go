package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type MapTokenRepo struct {
	db *gorm.DB
}

func NewMapTokenRepo(db *gorm.DB) *MapTokenRepo {
	return &MapTokenRepo{db: db}
}

func (r *MapTokenRepo) Create(token *models.MapToken) error {
	return r.db.Create(token).Error
}

func (r *MapTokenRepo) GetByID(id uuid.UUID) (*models.MapToken, error) {
	var token models.MapToken
	if err := r.db.First(&token, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *MapTokenRepo) Update(token *models.MapToken) error {
	return r.db.Save(token).Error
}

func (r *MapTokenRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.MapToken{}, "id = ?", id).Error
}

func (r *MapTokenRepo) GetByInitiativeOrder(mapID uuid.UUID) ([]models.MapToken, error) {
	var tokens []models.MapToken
	err := r.db.Where("map_id = ? AND initiative_order IS NOT NULL", mapID).
		Order("initiative_order ASC").Find(&tokens).Error
	return tokens, err
}

func (r *MapTokenRepo) SetInitiativeOrder(tokenID uuid.UUID, order *int) error {
	return r.db.Model(&models.MapToken{}).Where("id = ?", tokenID).
		Update("initiative_order", order).Error
}
