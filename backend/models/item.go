package models

import (
	"time"

	"github.com/google/uuid"
)

// Item represents a non-weapon object (consumables, tools, shields, utility items).
type Item struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	UserID      *uuid.UUID `json:"user_id,omitempty" gorm:"type:uuid;index"` // Null for system items
	Name        string     `json:"name" gorm:"type:text;not null"`
	Description string     `json:"description" gorm:"type:text"`

	Category string `json:"category" gorm:"type:text;not null"`       // e.g. "Potion", "Shield", "Tool", "Utility", "Gear", "Armor"
	Rarity   string `json:"rarity" gorm:"type:text;default:'Common'"` // e.g. "Common", "Uncommon", "Rare"

	// Costs and Weight
	Cost   string `json:"cost" gorm:"type:text"`   // e.g. "50 gp"
	Weight string `json:"weight" gorm:"type:text"` // e.g. "1 lb."

	// Effects / Properties
	// For shields, this might include "AC +2". For potions, "Heals 2d4+2 HP".
	Effects string `json:"effects" gorm:"type:text"`

	// Armor-specific fields
	ArmorType           *string `json:"armor_type,omitempty" gorm:"type:text"`              // e.g., "Light", "Medium", "Heavy", "Shield"
	BaseAC              *int    `json:"base_ac,omitempty" gorm:"type:int"`                  // Base AC provided by the armor
	StrengthRequirement *int    `json:"strength_requirement,omitempty" gorm:"type:int"`     // Minimum strength to wear without penalty
	StealthDisadvantage *bool   `json:"stealth_disadvantage,omitempty" gorm:"type:boolean"` // Does this armor impose disadvantage on stealth checks?

	// Consumable specific
	IsConsumable bool `json:"is_consumable" gorm:"type:bool;default:false"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User *User `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}
