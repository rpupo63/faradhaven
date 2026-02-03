package models

import (
	"time"

	"github.com/google/uuid"
)

// Race represents a core species (e.g., Elf, Dragonborn, Kalashtar)
type Race struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Name         string    `json:"name" gorm:"type:text;not null;uniqueIndex"`
	Description  string    `json:"description" gorm:"type:text"`
	CreatureType string    `json:"creature_type" gorm:"type:text;default:'Humanoid'"` // e.g., Aberration
	Size         string    `json:"size" gorm:"type:text"`                             // "Medium", "Small or Medium"
	BaseSpeed    int       `json:"base_speed" gorm:"type:int;default:30"`
	PhotoURL     string    `json:"photo_url" gorm:"type:text"` // URL to race artwork/photo

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// ============================================================
	// RELATIONSHIPS: Loaded via Preload, stored via foreign keys
	// ============================================================
	Traits   []Trait   `json:"traits" gorm:"foreignKey:RaceID;constraint:OnDelete:CASCADE"`
	Lineages []Lineage `json:"lineages" gorm:"foreignKey:RaceID;constraint:OnDelete:CASCADE"` // e.g., Wood Elf, High Elf
}
