package models

import (
	"time"

	"github.com/google/uuid"
)

// ClassLevel represents level-specific data for a class (levels 1-20).
// Comprehensive of D&D-style upgrades: HP, proficiency, spellcasting, class traits,
// resource pools, and structured combat upgrades (excluding subclasses).
type ClassLevel struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ClassID   uuid.UUID `json:"class_id" gorm:"type:uuid;not null;index:idx_class_level,unique"`
	Level     int       `json:"level" gorm:"type:int;not null;index:idx_class_level,unique"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// --- HP & survivability ---
	HpGain int `json:"hp_gain" gorm:"type:int;not null"` // average HP per level (after 1st)

	// --- Proficiency ---
	ProficiencyBonus int `json:"proficiency_bonus" gorm:"type:int;not null"` // e.g. 2 at level 1-4

	// --- Spellcasting ---
	CantripsKnown  *int `json:"cantrips_known,omitempty" gorm:"type:int"`  // cantrips known at this level
	SpellsKnown    *int `json:"spells_known,omitempty" gorm:"type:int"`    // spells known/prepared at this level
	MaxSpellPoints int  `json:"max_spell_points" gorm:"type:int;not null"` // spell point pool (50-100+ scale)

	// --- Class traits/features ---
	// Structured features with name and description
	LevelFeatures []LevelFeature `json:"level_features,omitempty" gorm:"foreignKey:ClassLevelID"`

	// --- Ability Score Improvement (ASI) ---
	// Points to add when leveling to this level: 0 = none, 2 = +2 to one ability or +1 to two (D&D standard at 4,8,12,16,19)
	AbilityScoreImprovement int `json:"ability_score_improvement" gorm:"type:int;default:0"`

	// --- Structured combat upgrades ---
	ExtraAttackCount  int `json:"extra_attack_count" gorm:"type:int;default:0"` // 0, 1, 2, or 3 for multiattack
	SneakAttackDice   int `json:"sneak_attack_dice" gorm:"type:int;default:0"`  // e.g. 2 for 2d6 (rogue)
	RageDamageBonus   int `json:"rage_damage_bonus" gorm:"type:int;default:0"`  // bonus melee damage when raging (barbarian)
	MartialArtsDie    int `json:"martial_arts_die" gorm:"type:int;default:0"`   // unarmed/martial arts die size: 4, 6, 8, 10 (monk)
	UnarmoredMovement int `json:"unarmored_movement" gorm:"type:int;default:0"` // extra feet of movement when unarmored (monk)
	SuperiorityDice   int `json:"superiority_dice" gorm:"type:int;default:0"`   // number of superiority dice (battlemaster)
	SuperiorityDie    int `json:"superiority_die" gorm:"type:int;default:0"`    // superiority die size: 6, 8, 10, 12
	BardicInspiration int `json:"bardic_inspiration" gorm:"type:int;default:0"` // die size for bardic inspiration: 6, 8, 10, 12

	// --- Shared D&D mechanics ---
	MaxSpellLevel int `json:"max_spell_level" gorm:"type:int;default:0"` // Spell level cap

	// --- Faradhaven Class Resources ---
	// Class-specific per-level resource values are stored in the ClassLevelResource table.
	// Query via ClassRepo.GetLevelResourceMap(classID, level) or GetLevelResourceValue(classID, level, key).

	Class Class `json:"-" gorm:"foreignKey:ClassID;references:ID"`
}
