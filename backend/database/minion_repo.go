package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type MinionRepository interface {
	FindByCharacterID(characterID uuid.UUID) ([]*models.Minion, error)
	FindByCharacterAndType(characterID uuid.UUID, minionType models.MinionType) ([]*models.Minion, error)
	FindByID(id uuid.UUID) (*models.Minion, error)
	Add(minion *models.Minion) error
	Update(minion *models.Minion) error
	Delete(id uuid.UUID) error
	DeleteByCharacterID(characterID uuid.UUID) error
	CountByCharacterAndType(characterID uuid.UUID, minionType models.MinionType) (int64, error)
}

type MinionRepo struct {
	db *gorm.DB
}

func NewMinionRepo(db *gorm.DB) *MinionRepo {
	return &MinionRepo{db}
}

// FindByCharacterID returns all minions for a character
func (r *MinionRepo) FindByCharacterID(characterID uuid.UUID) ([]*models.Minion, error) {
	var minions []*models.Minion
	err := r.db.Where("character_id = ?", characterID).Order("created_at asc").Find(&minions).Error
	return minions, err
}

// FindByCharacterAndType returns minions of a specific type for a character
func (r *MinionRepo) FindByCharacterAndType(characterID uuid.UUID, minionType models.MinionType) ([]*models.Minion, error) {
	var minions []*models.Minion
	err := r.db.Where("character_id = ? AND minion_type = ?", characterID, minionType).Order("created_at asc").Find(&minions).Error
	return minions, err
}

// FindByID returns a single minion by ID
func (r *MinionRepo) FindByID(id uuid.UUID) (*models.Minion, error) {
	var minion models.Minion
	err := r.db.First(&minion, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &minion, nil
}

// Add creates a new minion
func (r *MinionRepo) Add(minion *models.Minion) error {
	now := time.Now()
	minion.CreatedAt = now
	minion.UpdatedAt = now
	return r.db.Create(minion).Error
}

// Update updates an existing minion
func (r *MinionRepo) Update(minion *models.Minion) error {
	minion.UpdatedAt = time.Now()
	return r.db.Save(minion).Error
}

// Delete removes a minion by ID
func (r *MinionRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Minion{}, "id = ?", id).Error
}

// DeleteByCharacterID removes all minions for a character
func (r *MinionRepo) DeleteByCharacterID(characterID uuid.UUID) error {
	return r.db.Where("character_id = ?", characterID).Delete(&models.Minion{}).Error
}

// CountByCharacterAndType counts minions of a specific type for concurrency limits
func (r *MinionRepo) CountByCharacterAndType(characterID uuid.UUID, minionType models.MinionType) (int64, error) {
	var count int64
	err := r.db.Model(&models.Minion{}).Where("character_id = ? AND minion_type = ?", characterID, minionType).Count(&count).Error
	return count, err
}

// FindActiveByCharacterAndType returns only active minions of a type
func (r *MinionRepo) FindActiveByCharacterAndType(characterID uuid.UUID, minionType models.MinionType) ([]*models.Minion, error) {
	var minions []*models.Minion
	err := r.db.Where("character_id = ? AND minion_type = ? AND is_active = true", characterID, minionType).Order("created_at asc").Find(&minions).Error
	return minions, err
}

// FindEchoBySlot finds an echo in a specific slot
func (r *MinionRepo) FindEchoBySlot(characterID uuid.UUID, slotIndex int) (*models.Minion, error) {
	var minion models.Minion
	err := r.db.Where("character_id = ? AND minion_type = ? AND slot_index = ?", characterID, models.MinionTypeEcho, slotIndex).First(&minion).Error
	if err != nil {
		return nil, err
	}
	return &minion, nil
}
