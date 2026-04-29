package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ModifierType defines the type of weapon modifier
type ModifierType string

const (
	ModifierTypePistonCore    ModifierType = "piston_core"
	ModifierTypeVenomCoating  ModifierType = "venom_coating"
	ModifierTypeShadowCoating ModifierType = "shadow_coating"
	ModifierTypeLethalCoating ModifierType = "lethal_coating"
)

// WeaponModifier represents an active modifier on a character's weapon.
// Examples: Piston Core (permanent), Venom Coating (temporary per-attack).
type WeaponModifier struct {
	ID                uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	CharacterWeaponID uuid.UUID    `json:"character_weapon_id" gorm:"type:uuid;not null;index"`
	ModifierType      ModifierType `json:"modifier_type" gorm:"type:text;not null"`        // "piston_core", "venom_coating", etc.
	IsPermanent       bool         `json:"is_permanent" gorm:"type:boolean;default:false"` // Piston Core = true, Coatings = false
	IsActive          bool         `json:"is_active" gorm:"type:boolean;default:true"`     // Engaged/disengaged state

	// Type-specific data (e.g., Stability pool for Piston Core)
	Metadata ModifierMetadata `json:"metadata,omitempty" gorm:"type:jsonb"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	CharacterWeapon *CharacterWeapon `json:"-" gorm:"foreignKey:CharacterWeaponID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// ModifierMetadata holds type-specific data for weapon modifiers
type ModifierMetadata struct {
	// Piston Core fields
	CurrentStability *int `json:"current_stability,omitempty"`
	MaxStability     *int `json:"max_stability,omitempty"`
	IsEngaged        bool `json:"is_engaged,omitempty"`

	// Coating fields
	CoatingLevel   *int    `json:"coating_level,omitempty"` // For scaling coating effects by level
	CoatingCharges *int    `json:"coating_charges,omitempty"`
	BonusDamage    *string `json:"bonus_damage,omitempty"` // e.g., "1d4" or "2d6"
	DamageType     *string `json:"damage_type,omitempty"`  // e.g., "Poison", "Force"
}

// Value implements driver.Valuer for GORM
func (m ModifierMetadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan implements sql.Scanner for GORM
func (m *ModifierMetadata) Scan(value interface{}) error {
	if value == nil {
		*m = ModifierMetadata{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal ModifierMetadata value")
	}
	return json.Unmarshal(bytes, m)
}

// PistonCoreMetadata returns a ModifierMetadata configured for a Piston Core at the given level
func PistonCoreMetadata(level int) ModifierMetadata {
	maxStability := 10 + (level * 2)
	return ModifierMetadata{
		CurrentStability: &maxStability,
		MaxStability:     &maxStability,
		IsEngaged:        false,
	}
}
