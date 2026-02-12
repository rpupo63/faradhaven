package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Corpse represents a dead creature that can be harvested or consumed.
// Used by Ironwright (Scavenge) and Lorewright (Visceral Psychometry).
type Corpse struct {
	ID    uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	MapID *uuid.UUID `json:"map_id,omitempty" gorm:"type:uuid;index"`

	// Creature info
	Name            string  `json:"name" gorm:"type:text;not null"`
	CreatureType    string  `json:"creature_type" gorm:"type:text"`                    // "Beast", "Humanoid", "Undead", etc.
	CreatureSize    string  `json:"creature_size" gorm:"type:text;default:'Medium'"`   // "Small", "Medium", "Large", etc.
	ChallengeRating float64 `json:"challenge_rating" gorm:"type:decimal(4,2);default:0"`

	// Position (for map integration)
	GridX *int `json:"grid_x,omitempty" gorm:"type:int"`
	GridY *int `json:"grid_y,omitempty" gorm:"type:int"`

	// Harvest/loot data
	AvailableComponents pq.StringArray `json:"available_components" gorm:"type:text[]"` // Component IDs that can be harvested
	ComponentYield      int            `json:"component_yield" gorm:"type:int;default:1"` // Number of generic components

	// State tracking
	HasBeenHarvested bool `json:"has_been_harvested" gorm:"default:false"` // Ironwright scavenge
	HasBeenConsumed  bool `json:"has_been_consumed" gorm:"default:false"`  // Lorewright psychometry

	// Timing
	DiedAt    time.Time  `json:"died_at" gorm:"type:timestamptz;not null;default:now()"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" gorm:"type:timestamptz"` // Corpse decay time

	// Source creature (if from Beast compendium)
	SourceBeastID *uuid.UUID `json:"source_beast_id,omitempty" gorm:"type:uuid;index"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Map         *GameMap `json:"-" gorm:"foreignKey:MapID;references:ID;constraint:OnDelete:SET NULL"`
	SourceBeast *Beast   `json:"-" gorm:"foreignKey:SourceBeastID;references:ID;constraint:OnDelete:SET NULL"`
}

// TableName specifies the table name for GORM
func (Corpse) TableName() string {
	return "corpses"
}

// IsExpired checks if the corpse has expired
func (c *Corpse) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// MinutesSinceDeath returns how many minutes have passed since death
func (c *Corpse) MinutesSinceDeath() int {
	return int(time.Since(c.DiedAt).Minutes())
}

// CanBeHarvested checks if the corpse can still be harvested (within 1 hour of death)
func (c *Corpse) CanBeHarvested() bool {
	if c.HasBeenHarvested {
		return false
	}
	return c.MinutesSinceDeath() <= 60
}

// CanBeConsumed checks if the corpse can be consumed for psychometry (within 1 hour)
func (c *Corpse) CanBeConsumed() bool {
	if c.HasBeenConsumed {
		return false
	}
	return c.MinutesSinceDeath() <= 60
}
