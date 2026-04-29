package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type LevelUpHistoryRepository interface {
	FindByCharacterID(characterID uuid.UUID) ([]*models.LevelUpHistory, error)
	FindByCharacterAndLevel(characterID uuid.UUID, level int) (*models.LevelUpHistory, error)
	FindLatestByCharacter(characterID uuid.UUID) (*models.LevelUpHistory, error)
	Add(history *models.LevelUpHistory) error
	Delete(id uuid.UUID) error
	DeleteByCharacterAndLevel(characterID uuid.UUID, level int) error
}

type LevelUpHistoryRepo struct {
	db *gorm.DB
}

func NewLevelUpHistoryRepo(db *gorm.DB) *LevelUpHistoryRepo {
	return &LevelUpHistoryRepo{db}
}

// FindByCharacterID returns all level-up history for a character, ordered by level ASC
func (r *LevelUpHistoryRepo) FindByCharacterID(characterID uuid.UUID) ([]*models.LevelUpHistory, error) {
	var histories []*models.LevelUpHistory
	err := r.db.Where("character_id = ?", characterID).Order("level ASC").Find(&histories).Error
	return histories, err
}

// FindByCharacterAndLevel returns the history record for a specific character and level
func (r *LevelUpHistoryRepo) FindByCharacterAndLevel(characterID uuid.UUID, level int) (*models.LevelUpHistory, error) {
	var history models.LevelUpHistory
	err := r.db.First(&history, "character_id = ? AND level = ?", characterID, level).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// FindLatestByCharacter returns the most recent level-up history for a character
func (r *LevelUpHistoryRepo) FindLatestByCharacter(characterID uuid.UUID) (*models.LevelUpHistory, error) {
	var history models.LevelUpHistory
	err := r.db.Where("character_id = ?", characterID).Order("level DESC").First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// Add inserts a new level-up history record
func (r *LevelUpHistoryRepo) Add(history *models.LevelUpHistory) error {
	now := time.Now()
	history.CreatedAt = now
	history.UpdatedAt = now
	return r.db.Create(history).Error
}

// Delete removes a level-up history record by ID
func (r *LevelUpHistoryRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.LevelUpHistory{}, "id = ?", id).Error
}

// DeleteByCharacterAndLevel removes a level-up history record for a specific character and level
func (r *LevelUpHistoryRepo) DeleteByCharacterAndLevel(characterID uuid.UUID, level int) error {
	return r.db.Where("character_id = ? AND level = ?", characterID, level).Delete(&models.LevelUpHistory{}).Error
}
