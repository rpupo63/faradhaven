package models

import (
	"time"

	"github.com/google/uuid"
)

// Spell represents a crafted spell in the game
type Spell struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	CharacterID *uuid.UUID `json:"character_id,omitempty" gorm:"type:uuid;index"` // Character who created/prepared this spell (nil = user-level spell)
	Name        string     `json:"name" gorm:"type:text;not null"`
	Description string     `json:"description" gorm:"type:text"`
	SlotLevel   int        `json:"slot_level" gorm:"type:int;not null;default:1"`
	ManaCost    int        `json:"mana_cost" gorm:"type:int;not null;default:0"`
	BaseTier    int        `json:"base_tier" gorm:"type:int;not null;default:1"`

	// Spell Mechanics
	Type          string      `json:"type" gorm:"type:text;default:'Utility'"` // Attack, Save, Effect, Healing, Utility
	Range         *string     `json:"range,omitempty" gorm:"type:text"`
	Duration      *string     `json:"duration,omitempty" gorm:"type:text"`
	Concentration bool        `json:"concentration" gorm:"default:false"`
	SaveAttr      *string     `json:"save_attr,omitempty" gorm:"type:text"`   // e.g., "DEX", "WIS"
	DamageDice    *string     `json:"damage_dice,omitempty" gorm:"type:text"` // e.g., "2d8"
	DamageType    *DamageType `json:"damage_type,omitempty" gorm:"type:text"` // e.g., "Fire"
	AddModifier   bool        `json:"add_modifier" gorm:"default:false"`      // Add spellcasting mod to damage/healing

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User       User        `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Character  *Character  `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnDelete:SET NULL"`
	Components []Component `json:"components" gorm:"many2many:spell_components;foreignKey:ID;joinForeignKey:SpellID;References:ID;joinReferences:ComponentID"`
}
