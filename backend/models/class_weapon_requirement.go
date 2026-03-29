package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ClassWeaponRequirement defines when a class requires weapon selection during level-up.
// This triggers a weapon selection step in the level-up wizard.
type ClassWeaponRequirement struct {
	ID             uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ClassID        uuid.UUID    `json:"class_id" gorm:"type:uuid;not null;index"`
	SelectionLevel int          `json:"selection_level" gorm:"type:int;not null"` // Level at which selection is required (1 = character creation)
	ModifierType   ModifierType `json:"modifier_type" gorm:"type:text;not null"`  // e.g., "piston_core"
	Description    string       `json:"description" gorm:"type:text"`             // UI description for the selection

	// Weapon restrictions: if non-empty, only weapons in these categories are eligible
	AllowedCategories pq.StringArray `json:"allowed_categories" gorm:"type:text[]"` // e.g., ["Martial Melee", "Simple Melee"]

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Class *Class `json:"-" gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}
