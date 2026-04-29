package models

import (
	"time"

	"github.com/google/uuid"
)

// SpellComponent links a spell to the components used to craft it, in order.
// sort_order allows the same component_id to appear multiple times on one spell.
type SpellComponent struct {
	SpellID     uuid.UUID `json:"spell_id" gorm:"type:uuid;primaryKey;not null"`
	SortOrder   int       `json:"sort_order" gorm:"primaryKey;not null"`
	ComponentID uuid.UUID `json:"component_id" gorm:"type:uuid;not null;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	Spell     Spell     `json:"-" gorm:"foreignKey:SpellID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Component Component `json:"-" gorm:"foreignKey:ComponentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
