package models

import (
	"time"

	"github.com/google/uuid"
)

// ConsumptionHistory tracks when a Lorewright harvests a creature, enabling features
// like the Warlord's Predator's Strike.
type ConsumptionHistory struct {
	ID           uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	CharacterID  uuid.UUID    `json:"character_id" gorm:"type:uuid;not null;index"`
	CreatureType CreatureType `json:"creature_type" gorm:"type:text;not null"`
	HarvestedAt  time.Time    `json:"harvested_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Character Character `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
