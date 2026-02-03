package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type AttackRepository interface {
	FindByBeastID(beastID uuid.UUID) ([]*models.Attack, error)
	Add(attack *models.Attack) error
	Update(attack *models.Attack) error
	Delete(id uuid.UUID) error
}

type AttackRepo struct {
	db *gorm.DB
}

func NewAttackRepo(db *gorm.DB) *AttackRepo {
	return &AttackRepo{db}
}

// FindByBeastID returns all attacks for a beast
func (r *AttackRepo) FindByBeastID(beastID uuid.UUID) ([]*models.Attack, error) {
	var attacks []*models.Attack
	err := r.db.Where("beast_id = ?", beastID).Find(&attacks).Error
	return attacks, err
}

// Add inserts a new attack
func (r *AttackRepo) Add(attack *models.Attack) error {
	attack.CreatedAt = time.Now()
	return r.db.Create(attack).Error
}

// Update updates an existing attack
func (r *AttackRepo) Update(attack *models.Attack) error {
	return r.db.Save(attack).Error
}

// Delete removes an attack by ID
func (r *AttackRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Attack{}, "id = ?", id).Error
}
