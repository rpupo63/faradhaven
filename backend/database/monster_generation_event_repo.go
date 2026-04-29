package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type MonsterGenerationEventRepository interface {
	Add(event *models.MonsterGenerationEvent) error
	SummaryByUser(userID uuid.UUID) (map[string]int64, error)
}

type MonsterGenerationEventRepo struct {
	db *gorm.DB
}

func NewMonsterGenerationEventRepo(db *gorm.DB) *MonsterGenerationEventRepo {
	return &MonsterGenerationEventRepo{db: db}
}

func (r *MonsterGenerationEventRepo) Add(event *models.MonsterGenerationEvent) error {
	return r.db.Create(event).Error
}

func (r *MonsterGenerationEventRepo) SummaryByUser(userID uuid.UUID) (map[string]int64, error) {
	type row struct {
		EventType string
		Count     int64
	}
	var rows []row
	err := r.db.Model(&models.MonsterGenerationEvent{}).
		Select("event_type, count(*) as count").
		Where("user_id = ?", userID).
		Group("event_type").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to summarize generation events: %w", err)
	}
	out := make(map[string]int64, len(rows))
	for _, r := range rows {
		out[r.EventType] = r.Count
	}
	return out, nil
}
