package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// MonsterGenerationEvent captures diagnostics for AI generation quality.
type MonsterGenerationEvent struct {
	Base
	UserID        uuid.UUID      `json:"user_id"`
	MonsterID     *uuid.UUID     `json:"monster_id,omitempty"`
	EventType     string         `json:"event_type"` // preview, create, regenerate, variant
	PromptVersion string         `json:"prompt_version"`
	ModelName     string         `json:"model_name"`
	LatencyMs     int64          `json:"latency_ms"`
	RetryCount    int            `json:"retry_count"`
	Success       bool           `json:"success"`
	FailureType   string         `json:"failure_type,omitempty"`
	Metadata      datatypes.JSON `json:"metadata,omitempty" gorm:"type:jsonb"`
}
