package models

import (
	"time"

	"github.com/google/uuid"
)

// Effect represents a status effect, condition, or special state (e.g. Madness, Exhaustion, Feral).
type Effect struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Name        string    `json:"name" gorm:"type:text;not null;unique"`
	Description string    `json:"description" gorm:"type:text"`
	Category    string    `json:"category" gorm:"type:text"`  // e.g. "Condition", "Madness", "Class Feature"
	Mechanics   string    `json:"mechanics" gorm:"type:text"` // Detailed rules text

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
}
