package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type ConsumptionHistoryRepository interface {
	Create(history *models.ConsumptionHistory) error
	FindRecentByCharacterAndType(characterID uuid.UUID, creatureType models.CreatureType, since time.Time) (*models.ConsumptionHistory, error)
}

type ConsumptionHistoryRepo struct {
	db *gorm.DB
}

func NewConsumptionHistoryRepo(db *gorm.DB) *ConsumptionHistoryRepo {
	return &ConsumptionHistoryRepo{db}
}

// Create inserts a new consumption history record
func (r *ConsumptionHistoryRepo) Create(history *models.ConsumptionHistory) error {
	history.HarvestedAt = time.Now()
	return r.db.Create(history).Error
}

// FindRecentByCharacterAndType finds if a character has harvested a specific creature type since a given time
func (r *ConsumptionHistoryRepo) FindRecentByCharacterAndType(characterID uuid.UUID, creatureType models.CreatureType, since time.Time) (*models.ConsumptionHistory, error) {
	var history models.ConsumptionHistory
	err := r.db.
		Where("character_id = ? AND creature_type = ? AND harvested_at >= ?", characterID, creatureType, since).
		First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}
