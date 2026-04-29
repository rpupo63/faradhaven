package models

import "github.com/google/uuid"

// BeastSkill represents a structured skill proficiency for a Beast.
type BeastSkill struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	BeastID uuid.UUID `json:"beast_id" gorm:"type:uuid;not null;index"`
	Name    string    `json:"name" gorm:"type:text;not null"` // e.g., "Stealth", "Perception"
	Value   int       `json:"value" gorm:"type:int;not null"` // e.g., 6, 4

	Beast Beast `json:"-" gorm:"foreignKey:BeastID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
