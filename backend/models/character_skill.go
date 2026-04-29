package models

import (
	"time"

	"github.com/google/uuid"
)

// CharacterSkill represents a skill proficiency for a character.
// A row exists for each skill the character is proficient in.
// Proficient is always true when stored; the model supports future extension (e.g. expertise).
type CharacterSkill struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	CharacterID uuid.UUID `json:"character_id" gorm:"type:uuid;not null;uniqueIndex:idx_character_skill"`
	SkillID     string    `json:"skill_id" gorm:"type:text;not null;uniqueIndex:idx_character_skill"` // D&D 5e skill id (e.g. "persuasion", "stealth")
	Proficient  bool      `json:"proficient" gorm:"type:bool;not null;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	Character Character `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
