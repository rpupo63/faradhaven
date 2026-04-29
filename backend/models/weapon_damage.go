package models

import (
	"time"

	"github.com/google/uuid"
)

// WeaponDamage represents a specific damage component of a weapon.
// A weapon can have multiple damage entries (e.g., a magic sword dealing 1d8 Slashing + 1d6 Fire).
type WeaponDamage struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	WeaponID uuid.UUID `json:"weapon_id" gorm:"type:uuid;not null;index"`

	DamageDice string `json:"damage_dice" gorm:"type:text;not null"` // e.g. "1d8", "2d6"
	DamageType DamageType `json:"damage_type" gorm:"type:text;not null"`

	// Category distinguishes between the base damage and extra damage (e.g., from enchantments).
	DamageCategory WeaponDamageCategory `json:"damage_category" gorm:"type:text;default:'Base'"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	Weapon Weapon `json:"-" gorm:"foreignKey:WeaponID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
