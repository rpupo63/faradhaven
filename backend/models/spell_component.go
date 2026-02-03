package models

import (
	"time"

	"github.com/google/uuid"
)

// SpellComponent links a spell to the components used to craft it.
// When a user creates a spell, they select components from their class; this table stores that association.
type SpellComponent struct {
	SpellID     uuid.UUID `json:"spell_id" gorm:"type:uuid;primaryKey;not null"`
	ComponentID uuid.UUID `json:"component_id" gorm:"type:uuid;primaryKey;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	Spell     Spell     `json:"-" gorm:"foreignKey:SpellID;references:ID;constraint:OnDelete:CASCADE"`
	Component Component `json:"-" gorm:"foreignKey:ComponentID;references:ID;constraint:OnDelete:CASCADE"`
}
