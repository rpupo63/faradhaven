package models

import (
	"time"

	"github.com/google/uuid"
)

// Spell represents a crafted spell in the game
type Spell struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"type:text;not null"`
	Description string    `json:"description" gorm:"type:text"`
	SlotLevel   int       `json:"slot_level" gorm:"type:int;not null;default:1"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User       User        `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Components []Component `json:"components" gorm:"many2many:spell_components;foreignKey:ID;joinForeignKey:SpellID;References:ID;joinReferences:ComponentID"`
}
