package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type ItemRepository interface {
	FindAll() ([]*models.Item, error)
	FindByID(id uuid.UUID) (*models.Item, error)
	Add(item *models.Item) error
	Update(item *models.Item) error
	Delete(id uuid.UUID) error
}

type ItemRepo struct {
	db *gorm.DB
}

func NewItemRepo(db *gorm.DB) *ItemRepo {
	return &ItemRepo{db}
}

func (r *ItemRepo) FindAll() ([]*models.Item, error) {
	var items []*models.Item
	err := r.db.Find(&items).Error
	return items, err
}

func (r *ItemRepo) FindByID(id uuid.UUID) (*models.Item, error) {
	var item models.Item
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepo) Add(item *models.Item) error {
	return r.db.Create(item).Error
}

func (r *ItemRepo) Update(item *models.Item) error {
	return r.db.Save(item).Error
}

func (r *ItemRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Item{}, "id = ?", id).Error
}
