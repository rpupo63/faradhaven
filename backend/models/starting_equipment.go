package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ClassStartingEquipmentChoice represents a decision point for starting equipment (e.g., "Choice 1: Weapon A or Weapon B").
type ClassStartingEquipmentChoice struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ClassID     uuid.UUID `json:"class_id" gorm:"type:uuid;not null;index"`
	Instruction string    `json:"instruction" gorm:"type:text;not null"` // e.g., "Choose one of the following"
	SortOrder   int       `json:"sort_order" gorm:"type:int;default:0"`  // for ordering choices (Choice 1, Choice 2, etc.)
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Class   Class                          `json:"-" gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
	Options []ClassStartingEquipmentOption `json:"options" gorm:"foreignKey:ChoiceID;constraint:OnDelete:CASCADE"`
}

// ClassStartingEquipmentOption represents a valid selection within a choice.
type ClassStartingEquipmentOption struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ChoiceID    uuid.UUID      `json:"choice_id" gorm:"type:uuid;not null;index"`
	Description string         `json:"description" gorm:"type:text;not null"` // e.g., "Scale Mail and a Shield"
	Items       pq.StringArray `json:"items" gorm:"type:text[]"`              // e.g., ["Scale Mail", "Shield"] (Legacy/Automatic names)
	WeaponIDs   pq.StringArray `json:"weapon_ids" gorm:"type:text[]"`         // Array of Weapon UUID strings
	ItemIDs     pq.StringArray `json:"item_ids" gorm:"type:text[]"`           // Array of Item UUID strings
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	Choice ClassStartingEquipmentChoice `json:"-" gorm:"foreignKey:ChoiceID;references:ID;constraint:OnDelete:CASCADE"`
}
