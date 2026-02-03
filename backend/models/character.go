package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Character represents a player character in the game
type Character struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	UserID  uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Name    string    `json:"name" gorm:"type:text;not null"`
	RaceID  uuid.UUID `json:"race_id" gorm:"type:uuid;not null;index"`
	ClassID uuid.UUID `json:"class_id" gorm:"type:uuid;not null;index"`
	Level   int       `json:"level" gorm:"type:int;not null;default:1"`

	// Stored spell IDs (persisted in DB as text array)
	SpellbookIDs pq.StringArray `json:"spellbook_ids" gorm:"column:spellbook;type:text[]"`

	// Ability scores (for HP, AC, Save DC calculations)
	Strength     int `json:"strength" gorm:"type:int;default:10"`
	Dexterity    int `json:"dexterity" gorm:"type:int;default:10"`
	Constitution int `json:"constitution" gorm:"type:int;default:10"`
	Intelligence int `json:"intelligence" gorm:"type:int;default:10"`
	Wisdom       int `json:"wisdom" gorm:"type:int;default:10"`
	Charisma     int `json:"charisma" gorm:"type:int;default:10"`

	// Spell point pool: current points (max comes from ClassLevel)
	CurrentSpellPoints int `json:"current_spell_points" gorm:"type:int;default:0"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// ============================================================
	// RELATIONSHIPS: Loaded via Preload, stored via foreign keys
	// ============================================================
	User               User             `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Race               Race             `json:"race" gorm:"foreignKey:RaceID;references:ID"`
	Class              Class            `json:"class" gorm:"foreignKey:ClassID;references:ID"`
	SkillProficiencies []CharacterSkill `json:"-" gorm:"foreignKey:CharacterID;constraint:OnDelete:CASCADE"`

	// ============================================================
	// COMPUTED SUB-MODELS: Not stored in DB, populated in Go code
	// ============================================================

	// Populated from SkillProficiencies when loading for sheet
	SkillProficiencyIDs []string `json:"skill_proficiencies" gorm:"-"`

	// Populated from SpellbookIDs when loading full spell objects
	Spellbook []Spell `json:"spellbook,omitempty" gorm:"-"`

	// Current class level data (populated from ClassLevel table based on Level)
	CurrentClassLevel *ClassLevel `json:"current_class_level,omitempty" gorm:"-"`

	// Computed stats (HP, AC, etc.) - calculated in backend, not stored
	// See character_computed.go for CharacterComputedStats and ComputeStats()
	ComputedStats *CharacterComputedStats `json:"computed_stats,omitempty" gorm:"-"`
}
