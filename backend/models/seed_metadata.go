package models

import (
	"time"

	"github.com/google/uuid"
)

// SeedMetadata tracks seed execution history for versioning and skip optimization.
// By storing a content hash, we can skip unchanged seeds on subsequent runs.
type SeedMetadata struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	SeedName    string    `json:"seed_name" gorm:"type:text;uniqueIndex;not null"`
	ContentHash string    `json:"content_hash" gorm:"type:text;not null"`
	RecordCount int       `json:"record_count" gorm:"type:int;not null"`
	ExecutedAt  time.Time `json:"executed_at" gorm:"type:timestamptz;not null"`
	DurationMs  int       `json:"duration_ms" gorm:"type:int;not null"`
	Status      string    `json:"status" gorm:"type:text;not null"` // "success" or "failed"
}
