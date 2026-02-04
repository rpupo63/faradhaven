package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type ClassRepository interface {
	FindAll() ([]*models.Class, error)
	FindByID(id uuid.UUID) (*models.Class, error)
	FindByIDWithLevels(id uuid.UUID) (*models.Class, error)
	FindByName(name string) (*models.Class, error)
	FindLevelByClassAndLevel(classID uuid.UUID, level int) (*models.ClassLevel, error)
	Add(class *models.Class) error
	AddLevel(level *models.ClassLevel) error
}

type ClassRepo struct {
	db *gorm.DB
}

func NewClassRepo(db *gorm.DB) *ClassRepo {
	return &ClassRepo{db}
}

func (r *ClassRepo) FindAll() ([]*models.Class, error) {
	var classes []*models.Class
	err := r.db.Preload("Components").Preload("Levels", func(db *gorm.DB) *gorm.DB {
		return db.Where("level = ?", 1)
	}).Preload("Levels.LevelFeatures", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Find(&classes).Error
	return classes, err
}

func (r *ClassRepo) FindByID(id uuid.UUID) (*models.Class, error) {
	var c models.Class
	if err := r.db.First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClassRepo) FindByIDWithLevels(id uuid.UUID) (*models.Class, error) {
	var c models.Class
	if err := r.db.Preload("Levels", func(db *gorm.DB) *gorm.DB {
		return db.Order("level ASC")
	}).Preload("Levels.LevelFeatures", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Preload("Components").First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClassRepo) FindByName(name string) (*models.Class, error) {
	var c models.Class
	if err := r.db.First(&c, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClassRepo) FindLevelByClassAndLevel(classID uuid.UUID, level int) (*models.ClassLevel, error) {
	var cl models.ClassLevel
	if err := r.db.First(&cl, "class_id = ? AND level = ?", classID, level).Error; err != nil {
		return nil, err
	}
	return &cl, nil
}

func (r *ClassRepo) Add(class *models.Class) error {
	return r.db.Create(class).Error
}

func (r *ClassRepo) AddLevel(level *models.ClassLevel) error {
	return r.db.Create(level).Error
}
