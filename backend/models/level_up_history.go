package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// LevelUpHistory tracks each level-up event and stores a snapshot for reversal.
// This enables level-down by restoring the previous state.
type LevelUpHistory struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	CharacterID uuid.UUID `json:"character_id" gorm:"type:uuid;not null;index"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`

	// The level AFTER this level-up (e.g., 2 means "leveled up to level 2")
	Level int `json:"level" gorm:"type:int;not null"`

	// Choices made during this level-up
	SkillSelections   pq.StringArray `json:"skill_selections" gorm:"type:text[]"`           // new skills gained
	ASIAllocation     datatypes.JSON `json:"asi_allocation" gorm:"type:jsonb"`              // {"strength": 1, "dexterity": 1}
	SpellsLearned     pq.StringArray `json:"spells_learned" gorm:"type:text[]"`             // spell IDs added
	FeaturesGained    pq.StringArray `json:"features_gained" gorm:"type:text[]"`            // feature IDs/names
	ArchetypeSelected *uuid.UUID     `json:"archetype_selected,omitempty" gorm:"type:uuid"` // archetype chosen at this level (if any)

	// SNAPSHOT: Complete character state BEFORE this level-up (for reversion)
	CharacterSnapshot datatypes.JSON `json:"character_snapshot" gorm:"type:jsonb;not null"`
	SkillSnapshot     pq.StringArray `json:"skill_snapshot" gorm:"type:text[]"` // skill IDs before

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Character Character `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnDelete:CASCADE"`
	User      User      `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

// CharacterSnapshotData is the structure stored in CharacterSnapshot JSON
type CharacterSnapshotData struct {
	Level              int      `json:"level"`
	Strength           int      `json:"strength"`
	Dexterity          int      `json:"dexterity"`
	Constitution       int      `json:"constitution"`
	Intelligence       int      `json:"intelligence"`
	Wisdom             int      `json:"wisdom"`
	Charisma           int      `json:"charisma"`
	SpellbookIDs       []string `json:"spellbook_ids"`
	CurrentSpellPoints int      `json:"current_spell_points"`
	Notoriety          int      `json:"notoriety"`

	// HP and hit dice state for reversion
	CurrentHP   int `json:"current_hp"`
	MaxHP       int `json:"max_hp"`
	TempHP      int `json:"temp_hp"`
	HitDiceUsed int `json:"hit_dice_used"`

	// Archetype state for reversion (nil if no archetype was selected before this level)
	ArchetypeID *string `json:"archetype_id,omitempty"`
}
