package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// SavedSpell represents a pre-crafted spell combination saved in a "Speed Dial" slot.
// Used primarily by Powder Mage for quick-casting prepared component combinations.
type SavedSpell struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	CharacterID uuid.UUID `json:"character_id" gorm:"type:uuid;not null;index"`

	// Spell configuration
	Name         string         `json:"name" gorm:"type:text;not null"`
	ComponentIDs pq.StringArray `json:"component_ids" gorm:"type:text[]"` // Ordered list of component UUIDs
	SlotIndex    int            `json:"slot_index" gorm:"type:int;not null"`

	// Usage tracking (nil = unlimited until rest)
	UsesRemaining *int `json:"uses_remaining,omitempty" gorm:"type:int"`
	MaxUses       *int `json:"max_uses,omitempty" gorm:"type:int"`

	// Computed spell properties (cached for quick display)
	TotalCost   int     `json:"total_cost" gorm:"type:int;default:0"`   // Total resource cost
	DamageType  *string `json:"damage_type,omitempty" gorm:"type:text"` // Primary damage type
	Description *string `json:"description,omitempty" gorm:"type:text"` // User notes

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Character Character `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM
func (SavedSpell) TableName() string {
	return "saved_spells"
}
