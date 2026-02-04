package models

import (
	"time"

	"github.com/google/uuid"
)

// Component represents a spell-like component available to a class.
// ComponentType enum is defined in enums.go.
// Classes use components instead of traditional spells (e.g., Chimera metabolizes mutagens,
// Cog-Weaver builds rigs, Vitalist uses blood as conduit).
type Component struct {
	ID          uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Name        string            `json:"name" gorm:"type:text;not null"`
	Type        ComponentType     `json:"type" gorm:"type:text;not null"`     // base or modifier
	Category    ComponentCategory `json:"category" gorm:"type:text;not null"` // primordial, physical, etc.
	Description string            `json:"description" gorm:"type:text"`
	Element     string            `json:"element,omitempty" gorm:"type:text"` // e.g. fire, lightning, necrotic
	CreatedAt   time.Time         `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
}
