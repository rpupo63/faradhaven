package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type CorpseRepository interface {
	FindAll() ([]*models.Corpse, error)
	FindByMapID(mapID uuid.UUID) ([]*models.Corpse, error)
	FindByID(id uuid.UUID) (*models.Corpse, error)
	Add(corpse *models.Corpse) error
	Update(corpse *models.Corpse) error
	Delete(id uuid.UUID) error
	DeleteExpired() (int64, error)
}

type CorpseRepo struct {
	db *gorm.DB
}

func NewCorpseRepo(db *gorm.DB) *CorpseRepo {
	return &CorpseRepo{db}
}

// FindAll returns all corpses
func (r *CorpseRepo) FindAll() ([]*models.Corpse, error) {
	var corpses []*models.Corpse
	err := r.db.Order("died_at desc").Find(&corpses).Error
	return corpses, err
}

// FindByMapID returns all corpses on a specific map
func (r *CorpseRepo) FindByMapID(mapID uuid.UUID) ([]*models.Corpse, error) {
	var corpses []*models.Corpse
	err := r.db.Where("map_id = ?", mapID).Order("died_at desc").Find(&corpses).Error
	return corpses, err
}

// FindByID returns a single corpse by ID
func (r *CorpseRepo) FindByID(id uuid.UUID) (*models.Corpse, error) {
	var corpse models.Corpse
	err := r.db.First(&corpse, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &corpse, nil
}

// FindHarvestable returns corpses that can still be harvested (within 1 hour of death)
func (r *CorpseRepo) FindHarvestable() ([]*models.Corpse, error) {
	var corpses []*models.Corpse
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	err := r.db.Where("has_been_harvested = false AND died_at > ?", oneHourAgo).Find(&corpses).Error
	return corpses, err
}

// FindConsumable returns corpses that can still be consumed (within 1 hour of death)
func (r *CorpseRepo) FindConsumable() ([]*models.Corpse, error) {
	var corpses []*models.Corpse
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	err := r.db.Where("has_been_consumed = false AND died_at > ?", oneHourAgo).Find(&corpses).Error
	return corpses, err
}

// Add creates a new corpse
func (r *CorpseRepo) Add(corpse *models.Corpse) error {
	corpse.CreatedAt = time.Now()
	if corpse.DiedAt.IsZero() {
		corpse.DiedAt = time.Now()
	}
	return r.db.Create(corpse).Error
}

// Update updates an existing corpse
func (r *CorpseRepo) Update(corpse *models.Corpse) error {
	return r.db.Save(corpse).Error
}

// Delete removes a corpse by ID
func (r *CorpseRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Corpse{}, "id = ?", id).Error
}

// DeleteExpired removes all expired corpses and returns the count
func (r *CorpseRepo) DeleteExpired() (int64, error) {
	result := r.db.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Delete(&models.Corpse{})
	return result.RowsAffected, result.Error
}

// DeleteByMapID removes all corpses from a map
func (r *CorpseRepo) DeleteByMapID(mapID uuid.UUID) error {
	return r.db.Where("map_id = ?", mapID).Delete(&models.Corpse{}).Error
}
