package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Character represents a player character in the game
type Character struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	Name        string     `json:"name" gorm:"type:text;not null"`
	RaceID      uuid.UUID  `json:"race_id" gorm:"type:uuid;not null;index"`
	LineageID   *uuid.UUID `json:"lineage_id,omitempty" gorm:"type:uuid;index"`
	ClassID     uuid.UUID  `json:"class_id" gorm:"type:uuid;not null;index"`
	ArchetypeID *uuid.UUID `json:"archetype_id,omitempty" gorm:"type:uuid;index"` // nil until archetype is chosen
	Level       int        `json:"level" gorm:"type:int;not null;default:1"`

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

	// Class-specific resource tracking (current values; max comes from ClassLevel)
	CurrentStability  int `json:"current_stability" gorm:"type:int;default:0"`   // Piston Brawler
	CurrentBloodIchor int `json:"current_blood_ichor" gorm:"type:int;default:0"` // Sanguinist
	SanguineNotoriety int `json:"sanguine_notoriety" gorm:"type:int;default:0"`  // Net value between -20 and +20
	SanguineMP        int `json:"sanguine_mp" gorm:"type:int;default:0"`         // Medical Prodigy points
	SanguineBR        int `json:"sanguine_br" gorm:"type:int;default:0"`         // Blood Rage points
	MadnessCastCount  int `json:"madness_cast_count" gorm:"type:int;default:0"`  // Mutagen: casts since rest
	EchoSlotsUsed     int `json:"echo_slots_used" gorm:"type:int;default:0"`     // Lorewright

	// HP tracking (persisted, not computed)
	CurrentHP int `json:"current_hp" gorm:"type:int;default:0"`
	MaxHP     int `json:"max_hp" gorm:"type:int;default:0"`
	TempHP    int `json:"temp_hp" gorm:"type:int;default:0"`

	// Hit dice tracking: how many hit dice have been spent (total = Level)
	HitDiceUsed int `json:"hit_dice_used" gorm:"type:int;default:0"`

	// Currency (in Copper Pieces)
	// 1 gp = 100 cp, 1 sp = 10 cp, 1 pp = 1000 cp
	Money int64 `json:"money" gorm:"type:bigint;default:0"`

	// Languages known by this character (e.g., ["Common", "Elvish", "Dwarvish"])
	Languages pq.StringArray `json:"languages" gorm:"type:text[]"`

	// Character backstory (rich text stored as HTML)
	Backstory string `json:"backstory" gorm:"type:text"`
	// Hex color for the backstory text
	BackstoryHexColor string `json:"backstory_hex_color" gorm:"type:text;default:'#000000'"`

	// Character profile picture URL
	ImageURL string `json:"image_url" gorm:"type:text"`

	// Character notes (miscellaneous notes for the player)
	Notes string `json:"notes" gorm:"type:text"`

	// Inventory (Text Array) - Stores names of automatic/legacy items
	Inventory pq.StringArray `json:"inventory" gorm:"type:text[]"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// ============================================================
	// RELATIONSHIPS: Loaded via Preload, stored via foreign keys
	// ============================================================
	User               User                 `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Race               Race                 `json:"race" gorm:"foreignKey:RaceID;references:ID"`
	Class              Class                `json:"class" gorm:"foreignKey:ClassID;references:ID"`
	Archetype          *Archetype           `json:"archetype,omitempty" gorm:"foreignKey:ArchetypeID;references:ID"`
	SkillProficiencies []CharacterSkill     `json:"-" gorm:"foreignKey:CharacterID;constraint:OnDelete:CASCADE"`
	Components         []CharacterComponent `json:"components,omitempty" gorm:"foreignKey:CharacterID;constraint:OnDelete:CASCADE"`

	// Inventory Relationships
	CharacterWeapons []CharacterWeapon `json:"character_weapons,omitempty" gorm:"foreignKey:CharacterID;constraint:OnDelete:CASCADE"`
	Weapons          []Weapon          `json:"weapons,omitempty" gorm:"many2many:character_weapons_v2;constraint:OnDelete:CASCADE"`
	Items            []Item            `json:"items,omitempty" gorm:"many2many:character_items;constraint:OnDelete:CASCADE"`

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
