package models

import (
	"time"

	"github.com/google/uuid"
)

// Archetype represents a class path/subclass (e.g., "Path of the Warlord" for Lorewright)
type Archetype struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ClassID     uuid.UUID `json:"class_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"type:text;not null"`
	Description string    `json:"description" gorm:"type:text"`
	SortOrder   int       `json:"sort_order" gorm:"type:int;default:0"` // for ordering archetypes within a class
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Class Class `json:"-" gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}
