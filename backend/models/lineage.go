package models

import (
	"time"

	"github.com/google/uuid"
)

// Lineage handles sub-races or specific choices like Draconic Ancestry
type Lineage struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	RaceID      uuid.UUID `json:"race_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"type:text;not null"` // e.g., "Gold Dragon" or "Wood Elf"
	Description string    `json:"description" gorm:"type:text"`
	DamageType  string    `json:"damage_type" gorm:"type:text"` // e.g., Fire, Acid (for Dragonborn/Drow)

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships (stored in DB via foreign keys)
	Race *Race `json:"-" gorm:"foreignKey:RaceID;references:ID;constraint:OnDelete:CASCADE"`

	// Nested relationships (loaded via Preload)
	LineageTraits []Trait `json:"lineage_traits" gorm:"foreignKey:LineageID;constraint:OnDelete:CASCADE"`
}
