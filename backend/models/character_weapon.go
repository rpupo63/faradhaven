package models

import (
	"time"

	"github.com/google/uuid"
)

// CharacterWeapon represents a weapon in a character's inventory with metadata.
// This is an explicit join table that tracks weapon ownership and state.
type CharacterWeapon struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	CharacterID uuid.UUID  `json:"character_id" gorm:"type:uuid;not null;index"`
	WeaponID    uuid.UUID  `json:"weapon_id" gorm:"type:uuid;not null;index"`
	IsPrimary   bool       `json:"is_primary" gorm:"type:boolean;default:false"` // For signature item attachment (e.g., Piston Core)
	IsEquipped  bool       `json:"is_equipped" gorm:"type:boolean;default:false"`
	CustomName  *string    `json:"custom_name,omitempty" gorm:"type:text"`       // Optional nickname for the weapon

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Character *Character       `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnDelete:CASCADE"`
	Weapon    Weapon           `json:"weapon" gorm:"foreignKey:WeaponID;references:ID"`
	Modifiers []WeaponModifier `json:"modifiers,omitempty" gorm:"foreignKey:CharacterWeaponID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM
func (CharacterWeapon) TableName() string {
	return "character_weapons_v2"
}
