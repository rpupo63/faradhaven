package models

import (
	"time"

	"github.com/google/uuid"
)

// Attack represents an attack action for a beast.
// Stored in DB with foreign key to Beast.
type Attack struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	BeastID     uuid.UUID  `json:"beast_id" gorm:"type:uuid;not null;index"`
	Name        string     `json:"name" gorm:"type:text;not null"`
	AttackBonus int        `json:"attack_bonus" gorm:"type:int;not null"`
	DamageType  DamageType `json:"damage_type" gorm:"type:text;not null"`
	DamageDice  string     `json:"damage_dice" gorm:"type:text;not null"` // e.g., "2d6+4"
	Reach       *string    `json:"reach" gorm:"type:text"`                // e.g., "5 ft." or "30/120 ft."
	Description *string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationship back to Beast
	Beast Beast `json:"-" gorm:"foreignKey:BeastID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
