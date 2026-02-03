package models

import (
	"time"

	"github.com/google/uuid"
)

// Trait represents a specific ability (Healing Hands, Mind Link, Breath Weapon)
type Trait struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	RaceID      *uuid.UUID `json:"race_id" gorm:"type:uuid;index"`    // Null if it belongs to a Lineage
	LineageID   *uuid.UUID `json:"lineage_id" gorm:"type:uuid;index"` // Null if it belongs to a base Race
	Name        string     `json:"name" gorm:"type:text;not null"`
	Description string     `json:"description" gorm:"type:text"`
	LevelReq    int        `json:"level_req" gorm:"type:int;default:1"`

	// Activation & Usage
	ActionType     string `json:"action_type" gorm:"type:text"`     // "Action", "Bonus Action", "Reaction", "Passive"
	UsesPerRest    string `json:"uses_per_rest" gorm:"type:text"`   // "1", "Proficiency Bonus"
	ResetCondition string `json:"reset_condition" gorm:"type:text"` // "Long Rest", "Short Rest"

	// Scaling & Mechanics
	RangeValue   string `json:"range_value" gorm:"type:text"`  // e.g., "60", "10 * level" (Mind Link)
	AreaOfEffect string `json:"aoe" gorm:"type:text"`          // "15ft Cone", "30ft Line"
	SaveAbility  string `json:"save_ability" gorm:"type:text"` // "DEX", "CON", "CHA"

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships (stored in DB via foreign keys)
	Race    *Race    `json:"-" gorm:"foreignKey:RaceID;references:ID;constraint:OnDelete:CASCADE"`
	Lineage *Lineage `json:"-" gorm:"foreignKey:LineageID;references:ID;constraint:OnDelete:CASCADE"`

	// Nested Options (for Aasimar Revelations or Spell Selection)
	Options []TraitOption `json:"options" gorm:"foreignKey:TraitID;constraint:OnDelete:CASCADE"`
}
