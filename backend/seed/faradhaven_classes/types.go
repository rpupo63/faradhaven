package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// EquipmentOptionSeed defines a single option within a choice (e.g. "Scale Mail and Shield")
type EquipmentOptionSeed struct {
	Description string
	Items       []string // For generic items/names
	WeaponNames []string // Names of weapons to look up and link
	ItemNames   []string // Names of items to look up and link
}

// EquipmentChoiceSeed defines a decision point for equipment
type EquipmentChoiceSeed struct {
	Instruction string // e.g. "Choose one of the following"
	Options     []EquipmentOptionSeed
}

// ArchetypeSeed defines archetype/subclass data for seeding
type ArchetypeSeed struct {
	Name        string
	Description string
	Features    map[int]string // level -> feature text (path-specific features)
}

// ClassLevelSeed defines structured level progression data for seeding
// These fields map directly to models.ClassLevel for level-up mechanics
type ClassLevelSeed struct {
	CantripsKnown     *int // cantrips known at this level (nil = no change)
	SpellsKnown       *int // spells known/prepared at this level (nil = no change)
	MaxSpellPoints    int  // override for spell point pool (0 = use default formula)
	ExtraAttackCount  int  // 0, 1, 2, or 3 for multiattack progression
	SneakAttackDice   int  // e.g. 2 for 2d6 sneak attack (rogue-like)
	RageDamageBonus   int  // bonus melee damage when raging (barbarian-like)
	MartialArtsDie    int  // unarmed/martial arts die size: 4, 6, 8, 10 (monk-like)
	UnarmoredMovement int  // extra feet of movement when unarmored
	SuperiorityDice   int  // number of superiority dice (battlemaster-like)
	SuperiorityDie    int  // superiority die size: 6, 8, 10, 12
	BardicInspiration int  // die size for inspiration: 6, 8, 10, 12
}

// FaradhavenClassSeed defines the full class data for seeding
type FaradhavenClassSeed struct {
	Name           string
	Description    string // Short 2-sentence description of how this class uniquely plays
	HitDie         int
	PrimaryAbility string
	PhotoURL       string // URL to class artwork/photo
	Archetype      string // Legacy: thematic archetype description
	Concept        string
	ClassFeatures  []string
	DnDSkillFocus  []string
	Proficiencies  string
	SkillChoice    []string
	Tools          []string
	SavingThrows   []string
	AutomaticEquipNames  []string // Generic names
	AutomaticWeaponNames []string // Names of weapons to look up and link
	AutomaticItemNames   []string // Names of items to look up and link

	// EquipmentChoices defines the decision points for starting equipment
	EquipmentChoices []EquipmentChoiceSeed

	// LevelFeatures maps level (1-20) to features gained at that level.
	// Used for D&D-style level progression. Level 1 base info is always built from the seed;
	// LevelFeatures[1] appends to it; LevelFeatures[2-20] set Features for those levels.
	// For archetype levels, these are shared features that all archetypes get.
	LevelFeatures map[int]string

	// LevelProgression maps level (1-20) to structured level data for combat/spellcasting progression.
	// These values are applied to models.ClassLevel when seeding.
	LevelProgression map[int]ClassLevelSeed

	// ArchetypeLevel is the level at which players choose their archetype (nil = no archetypes)
	ArchetypeLevel *int

	// Archetypes defines the available subclass paths for this class
	Archetypes []ArchetypeSeed
}

// ComponentSeed defines component data for seeding (used in faradhaven_components.go)
type ComponentSeed struct {
	Name        string
	Symbol      string // 2-3 letter alchemical symbol for periodic table display (e.g., "Ig" for Ignis)
	Category    models.ComponentCategory
	Description string
	Element     string
}
