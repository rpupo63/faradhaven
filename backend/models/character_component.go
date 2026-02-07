package models

import (
	"time"

	"github.com/google/uuid"
)

// CharacterComponent tracks the quantity of components a character currently possesses.
type CharacterComponent struct {
	CharacterID uuid.UUID `json:"character_id" gorm:"type:uuid;primaryKey;not null"`
	ComponentID uuid.UUID `json:"component_id" gorm:"type:uuid;primaryKey;not null"`
	Count       int       `json:"count" gorm:"type:int;not null;default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	Character Character `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnDelete:CASCADE"`
	Component Component `json:"component" gorm:"foreignKey:ComponentID;references:ID;constraint:OnDelete:CASCADE"`
}
