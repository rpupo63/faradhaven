package models

import (
	"time"

	"github.com/google/uuid"
)

// CharacterEffect represents an active effect on a character.
// This is a join table between Character and Effect with additional instance data.
type CharacterEffect struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	CharacterID uuid.UUID `json:"character_id" gorm:"type:uuid;not null;index"`
	EffectID    uuid.UUID `json:"effect_id" gorm:"type:uuid;not null;index"`

	// Optional override or specific instance details
	Duration string `json:"duration,omitempty" gorm:"type:text"` // e.g. "1 minute", "Indefinite"
	Source   string `json:"source,omitempty" gorm:"type:text"`   // e.g. "Feral Table Roll", "Spell: Haste"

	// Stack tracking (for effects like Venom Intensity)
	Stacks    int `json:"stacks" gorm:"type:int;default:1"`
	MaxStacks int `json:"max_stacks" gorm:"type:int;default:1"`

	// Duration tracking (for turn-based countdown)
	DurationRounds  *int `json:"duration_rounds,omitempty" gorm:"type:int"`  // Rounds remaining (nil = indefinite)
	DurationMinutes *int `json:"duration_minutes,omitempty" gorm:"type:int"` // Minutes remaining

	// Source tracking (who applied this effect)
	SourceCharacterID *uuid.UUID `json:"source_character_id,omitempty" gorm:"type:uuid;index"`
	SourceSpellID     *uuid.UUID `json:"source_spell_id,omitempty" gorm:"type:uuid;index"`

	// Concentration tracking (for spells requiring concentration)
	IsConcentration bool `json:"is_concentration" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Character       Character  `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnDelete:CASCADE"`
	Effect          Effect     `json:"effect" gorm:"foreignKey:EffectID;references:ID;constraint:OnDelete:CASCADE"`
	SourceCharacter *Character `json:"-" gorm:"foreignKey:SourceCharacterID;references:ID;constraint:OnDelete:SET NULL"`
	SourceSpell     *Spell     `json:"-" gorm:"foreignKey:SourceSpellID;references:ID;constraint:OnDelete:SET NULL"`
}
