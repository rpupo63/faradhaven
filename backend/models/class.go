package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Class represents a character class (high-level data)
type Class struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Name           string    `json:"name" gorm:"type:text;not null;uniqueIndex"`
	HitDie         int       `json:"hit_die" gorm:"type:int;not null"`          // e.g. 10 for d10
	PrimaryAbility string    `json:"primary_ability" gorm:"type:text;not null"` // e.g. "intelligence", "wisdom"
	PhotoURL       string    `json:"photo_url" gorm:"type:text"`                // URL to class artwork/photo

	// D&D-style proficiencies and starting data (stored as PostgreSQL arrays)
	Proficiencies    string         `json:"proficiencies" gorm:"type:text"`               // weapon/armor proficiencies (e.g. "Simple weapons, Light Armor")
	SkillFocus       pq.StringArray `json:"skill_focus" gorm:"type:text[]"`               // fixed skill proficiencies (e.g. Persuasion, Medicine)
	SkillChoice      pq.StringArray `json:"skill_choice" gorm:"type:text[]"`              // choose N from these skills
	SkillChoiceCount int            `json:"skill_choice_count" gorm:"type:int;default:2"` // how many to pick from SkillChoice
	Tools            pq.StringArray `json:"tools" gorm:"type:text[]"`                     // tool proficiencies
	SavingThrows     pq.StringArray `json:"saving_throws" gorm:"type:text[]"`             // ability saving throw proficiencies (ability names)
	StartingEquip    pq.StringArray `json:"starting_equipment" gorm:"type:text[]"`        // starting equipment list

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// ============================================================
	// RELATIONSHIPS: Loaded via Preload, stored via foreign keys
	// ============================================================
	Levels     []ClassLevel `json:"levels,omitempty" gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE"`
	Components []Component  `json:"components,omitempty" gorm:"many2many:class_components;foreignKey:ID;joinForeignKey:ClassID;References:ID;joinReferences:ComponentID"`
}
