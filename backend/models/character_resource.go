package models

import (
	"time"

	"github.com/google/uuid"
)

// CharacterResource represents a dynamic resource tracked for a character.
// This is used for class-specific resources that don't fit into the standard fields,
// such as Ironwright components, extracted essences, or custom pools.
type CharacterResource struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	CharacterID uuid.UUID `json:"character_id" gorm:"type:uuid;not null;uniqueIndex:idx_char_resource_lookup,priority:1"`

	// Resource identification
	// The unique composite index on (character_id, resource_key) covers FindByCharacterAndKey queries.
	ResourceKey  string `json:"resource_key" gorm:"type:text;not null;uniqueIndex:idx_char_resource_lookup,priority:2"`
	ResourceName string `json:"resource_name" gorm:"type:text;not null"` // Display name

	// Current and max values
	CurrentValue int  `json:"current_value" gorm:"type:int;default:0"`
	MaxValue     *int `json:"max_value,omitempty" gorm:"type:int"` // nil = unlimited

	// Rest behavior
	RestoreOnShortRest bool `json:"restore_on_short_rest" gorm:"default:false"`
	RestoreOnLongRest  bool `json:"restore_on_long_rest" gorm:"default:false"`
	RestoreAmount      *int `json:"restore_amount,omitempty" gorm:"type:int"` // nil = full restore to max

	// Decay behavior (for temporary resources like extracted components)
	DecaysOnLongRest bool `json:"decays_on_long_rest" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Character Character `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM
func (CharacterResource) TableName() string {
	return "character_resources"
}
