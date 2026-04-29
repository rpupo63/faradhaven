package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
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

	// Ability Score Bonuses (e.g., {"dexterity": 2, "charisma": 1})
	AbilityScoreBonuses datatypes.JSON `json:"ability_score_bonuses" gorm:"type:jsonb"`

	// Languages known by this race (e.g., ["Common", "Elvish"])
	Languages pq.StringArray `json:"languages" gorm:"type:text[]"`
	// Number of bonus languages the player can choose during character creation
	BonusLanguageCount int `json:"bonus_language_count" gorm:"type:int;default:0"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// ============================================================
	// RELATIONSHIPS: Loaded via Preload, stored via foreign keys
	// ============================================================
	Traits     []Trait     `json:"traits" gorm:"foreignKey:RaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Lineages   []Lineage   `json:"lineages" gorm:"foreignKey:RaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // e.g., Wood Elf, High Elf
	Components []Component `json:"components,omitempty" gorm:"many2many:race_components;foreignKey:ID;joinForeignKey:RaceID;References:ID;joinReferences:ComponentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
