package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// TraitOption handles sub-choices (e.g., Necrotic Shroud vs Heavenly Wings)
type TraitOption struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	TraitID     uuid.UUID `json:"trait_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"type:text;not null"`
	Description string    `json:"description" gorm:"type:text"`
	SpellName   string    `json:"spell_name" gorm:"type:text"` // For Elven Lineage spells

	// Ability Score Bonuses (e.g., {"intelligence": 1})
	AbilityScoreBonuses datatypes.JSON `json:"ability_score_bonuses" gorm:"type:jsonb"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships (stored in DB via foreign keys)
	Trait Trait `json:"-" gorm:"foreignKey:TraitID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
