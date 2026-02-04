package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Weapon represents a weapon item (standard or homebrew).
// Can be a system standard weapon (UserID is null) or a user-created custom weapon.
type Weapon struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	UserID      *uuid.UUID `json:"user_id,omitempty" gorm:"type:uuid;index"` // Null for system weapons
	Name        string     `json:"name" gorm:"type:text;not null"`
	Description string     `json:"description" gorm:"type:text"`

	Category  string `json:"category" gorm:"type:text;not null"`   // e.g. "Simple Melee", "Martial Ranged"
	Rarity    string `json:"rarity" gorm:"type:text;default:'Common'"` // e.g. "Common", "Rare", "Legendary"
	RangeType string `json:"range_type" gorm:"type:text;not null"` // "Melee" or "Ranged"

	// Costs and Weight
	Cost   string `json:"cost" gorm:"type:text"`   // e.g. "15 gp"
	Weight string `json:"weight" gorm:"type:text"` // e.g. "2 lb."

	// Combat Stats
	// AttackModifier hints at which ability score to use (e.g., "Strength", "Dexterity")
	AttackModifier string `json:"attack_modifier" gorm:"type:text;default:'Strength'"`

	// Properties (e.g., "Finesse", "Light", "Thrown", "Two-Handed", "Reach")
	Properties pq.StringArray `json:"properties" gorm:"type:text[]"`

	// Range (in feet).
	// Melee: usually 5 (or 10 with Reach).
	// Ranged/Thrown: Normal/Long (e.g., 20/60).
	RangeNormal int `json:"range_normal" gorm:"type:int;default:5"`
	RangeLong   int `json:"range_long" gorm:"type:int;default:0"`

	// Versatile Variation
	// If the weapon has the "Versatile" property, this stores the damage dice when used with two hands.
	VersatileDamageDice *string `json:"versatile_damage_dice" gorm:"type:text"` // e.g. "1d10"

	// Secondary effects (text description of special rules, magical effects, etc.)
	SecondaryEffect string `json:"secondary_effect" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User    *User          `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Damages []WeaponDamage `json:"damages" gorm:"foreignKey:WeaponID;constraint:OnDelete:CASCADE"`
}
