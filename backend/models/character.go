package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"regexp"
)

// HarvestedAbilities defines the structure for storing Lorewright's harvested skills, attacks, and recipes.
type HarvestedAbilities struct {
	Skills  []string `json:"skills"`
	Attacks []string `json:"attacks"` // Storing as string IDs/names for now
	Recipes []string `json:"recipes"` // Storing as string IDs/names for now
}

// Character represents a player character in the game
type Character struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	Name        string     `json:"name" gorm:"type:text;not null"`
	RaceID      uuid.UUID  `json:"race_id" gorm:"type:uuid;not null;index"`
	LineageID   *uuid.UUID `json:"lineage_id,omitempty" gorm:"type:uuid;index"`
	ClassID     uuid.UUID  `json:"class_id" gorm:"type:uuid;not null;index"`
	ArchetypeID *uuid.UUID `json:"archetype_id,omitempty" gorm:"type:uuid;index"` // nil until archetype is chosen
	PartyID     *uuid.UUID `json:"party_id,omitempty" gorm:"type:uuid;index"`     // NEW: Foreign key to Party
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
	CurrentStability   int            `json:"current_stability" gorm:"type:int;default:0"`   // Piston Brawler
	CurrentBloodIchor  int            `json:"current_blood_ichor" gorm:"type:int;default:0"` // Sanguinist
	SanguineNotoriety  int            `json:"sanguine_notoriety" gorm:"type:int;default:0"`
	SanguineMP         int            `json:"sanguine_mp" gorm:"type:int;default:0"`          // Medical Prodigy points
	SanguineBR         int            `json:"sanguine_br" gorm:"type:int;default:0"`          // Blood Rage points
	EchoSlotsUsed      int            `json:"echo_slots_used" gorm:"type:int;default:0"`      // Lorewright (legacy)
	HarvestSkillsUsed  int            `json:"harvest_skills_used" gorm:"type:int;default:0"`  // Lorewright: skill slots used
	HarvestAttacksUsed int            `json:"harvest_attacks_used" gorm:"type:int;default:0"` // Lorewright: attack slots used
	HarvestRecipesUsed int            `json:"harvest_recipes_used" gorm:"type:int;default:0"` // Lorewright: recipe slots used
	HarvestedAbilities datatypes.JSON `json:"harvested_abilities" gorm:"type:jsonb"`          // Lorewright: stored abilities
	Trauma             int            `json:"trauma" gorm:"type:int;default:0"`               // Lorewright: psychological damage from harvesting
	LastComponentUsed  *string        `json:"last_component_used,omitempty" gorm:"type:text"` // Powder Mage: for Echoing Components

	// Mutagen Class Mechanics:
	// The following fields track the state for the Mutagen's Madness mechanic.
	// This state should be exposed to the backend API to be displayed in the Character Sheet UI.
	//
	// Formula for triggering a Madness Save:
	// A save is required when a Mutagen casts a component-ability if:
	//   `ability.Level >= 1` AND
	//   (`ability.ComponentCount > character.ProficiencyBonus` OR `MadnessCastCount + 1 > character.ProficiencyBonus + character.ConstitutionModifier`)
	//
	// State Update Logic (on cast of ANY component ability):
	// 1. Increment `MadnessCastCount` regardless of whether a save was triggered.
	// 2. If the cast required a save (per the formula above), increase `CurrentMadnessDC`.
	//    - Increase by +2 if Character.Level < 6.
	//    - Increase by +1 if Character.Level >= 6.
	//
	// Reset Logic (on Short or Long Rest):
	// - Set `CurrentMadnessDC` to 10.
	// - Set `MadnessCastCount` to 0.
	MadnessCastCount int `json:"madness_cast_count" gorm:"type:int;default:0"`  // Mutagen: casts since rest, used to determine if a save is needed.
	CurrentMadnessDC int `json:"current_madness_dc" gorm:"type:int;default:10"` // Mutagen: escalating DC for Madness saves.

	// Concentration tracking (for multi-concentration at high levels)
	ConcentrationSlots pq.StringArray `json:"concentration_slots" gorm:"type:text[]"` // CharacterEffect IDs being concentrated on

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

	// Equipment IDs
	EquippedArmorID  *uuid.UUID `json:"equipped_armor_id,omitempty" gorm:"type:uuid"`
	EquippedShieldID *uuid.UUID `json:"equipped_shield_id,omitempty" gorm:"type:uuid"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// ============================================================
	// RELATIONSHIPS: Loaded via Preload, stored via foreign keys
	// ============================================================
	User               User                 `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Race               Race                 `json:"race" gorm:"foreignKey:RaceID;references:ID"`
	Class              Class                `json:"class" gorm:"foreignKey:ClassID;references:ID"`
	Archetype          *Archetype           `json:"archetype,omitempty" gorm:"foreignKey:ArchetypeID;references:ID"`
	Party              *Party               `json:"party,omitempty" gorm:"foreignKey:PartyID;references:ID;constraint:OnDelete:SET NULL"` // NEW: Relationship to Party
	SkillProficiencies []CharacterSkill     `json:"-" gorm:"foreignKey:CharacterID;constraint:OnDelete:CASCADE"`
	Components         []CharacterComponent `json:"components,omitempty" gorm:"foreignKey:CharacterID;constraint:OnDelete:CASCADE"`
	ActiveEffects      []CharacterEffect    `json:"active_effects,omitempty" gorm:"foreignKey:CharacterID;constraint:OnDelete:CASCADE"`

	// Inventory Relationships
	CharacterWeapons []CharacterWeapon `json:"character_weapons,omitempty" gorm:"foreignKey:CharacterID;constraint:OnDelete:CASCADE"`
	Weapons          []Weapon          `json:"weapons,omitempty" gorm:"many2many:character_weapons_v2;constraint:OnDelete:CASCADE"`
	Items            []Item            `json:"items,omitempty" gorm:"many2many:character_items;constraint:OnDelete:CASCADE"`
	EquippedArmor    *Item             `json:"equipped_armor,omitempty" gorm:"foreignKey:EquippedArmorID;references:ID"`
	EquippedShield   *Item             `json:"equipped_shield,omitempty" gorm:"foreignKey:EquippedShieldID;references:ID"`

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

	// Carried Weight (calculated)
	CarriedWeight float64 `json:"carried_weight" gorm:"-"`
}

// GetUnarmoredDefenseACBonus retrieves the unarmored defense bonus from character features.
// It returns the base AC (e.g., 10), a map of ability modifiers to add (e.g., {"dexterity": true, "charisma": true}),
// and a boolean indicating if shields are allowed with this unarmored defense.
func (c *Character) GetUnarmoredDefenseACBonus() (int, map[string]bool, bool) {
	baseAC := 0
	abilityMods := make(map[string]bool)
	shieldsAllowed := true

	// Check Class Level Features (flattened from Class.Levels[].LevelFeatures)
	for _, level := range c.Class.Levels {
		if level.Level > c.Level {
			continue
		}
		for _, feature := range level.LevelFeatures {
			if feature.ArchetypeID == nil || (c.ArchetypeID != nil && *feature.ArchetypeID == *c.ArchetypeID) {
				if result, ok := parseUnarmoredDefense(feature.Name, feature.Description); ok {
					return result.baseAC, result.abilityMods, result.shieldsAllowed
				}
			}
		}
	}

	// Check Race Traits
	for _, trait := range c.Race.Traits {
		if trait.LevelReq <= c.Level {
			if result, ok := parseUnarmoredDefense(trait.Name, trait.Description); ok {
				return result.baseAC, result.abilityMods, result.shieldsAllowed
			}
		}
	}

	return baseAC, abilityMods, shieldsAllowed
}

type unarmoredDefenseResult struct {
	baseAC         int
	abilityMods    map[string]bool
	shieldsAllowed bool
}

// parseUnarmoredDefense checks if a feature/trait grants Unarmored Defense and extracts the AC formula.
func parseUnarmoredDefense(name, description string) (unarmoredDefenseResult, bool) {
	if !strings.Contains(description, "Unarmored Defense") && !strings.Contains(name, "Unarmored Defense") {
		return unarmoredDefenseResult{}, false
	}

	result := unarmoredDefenseResult{
		abilityMods:    make(map[string]bool),
		shieldsAllowed: true,
	}

	re := regexp.MustCompile(`AC equals (\d+) \+ (?:your )?(\w+) modifier(?: \+ (?:your )?(\w+) modifier)?`)
	matches := re.FindStringSubmatch(description)

	validAbilities := map[string]bool{"dexterity": true, "constitution": true, "wisdom": true, "charisma": true}

	if len(matches) > 0 {
		parsedBaseAC, err := strconv.Atoi(matches[1])
		if err == nil {
			result.baseAC = parsedBaseAC
			if key := strings.ToLower(matches[2]); validAbilities[key] {
				result.abilityMods[key] = true
			}
			if len(matches) > 3 && matches[3] != "" {
				if key := strings.ToLower(matches[3]); validAbilities[key] {
					result.abilityMods[key] = true
				}
			}
		}
	} else {
		// Fallback for descriptions the regex can't parse
		switch {
		case strings.Contains(description, "AC equals 10 + your Dexterity modifier + your Charisma modifier"):
			result.baseAC = 10
			result.abilityMods["dexterity"] = true
			result.abilityMods["charisma"] = true
		case strings.Contains(description, "AC equals 10 + your Dexterity modifier + your Constitution modifier"):
			result.baseAC = 10
			result.abilityMods["dexterity"] = true
			result.abilityMods["constitution"] = true
		case strings.Contains(description, "AC equals 10 + your Wisdom modifier + your Constitution modifier"):
			result.baseAC = 10
			result.abilityMods["wisdom"] = true
			result.abilityMods["constitution"] = true
		case strings.Contains(description, "AC equals 10 + your Dexterity modifier"):
			result.baseAC = 10
			result.abilityMods["dexterity"] = true
		}
	}

	if strings.Contains(description, "You cannot use a shield with this feature") {
		result.shieldsAllowed = false
	}

	return result, result.baseAC > 0
}

// COINS_PER_POUND defines how many coins (of any type) weigh one pound.
const COINS_PER_POUND = 50.0

// CalculateCarriedWeight calculates the total weight carried by the character.
// This includes money, items, and weapons.
func (c *Character) CalculateCarriedWeight() {
	var totalWeight float64

	// 1. Weight from Money (Copper Pieces)
	// Assuming 1 gp = 100 cp, and 50 coins weigh 1 pound.
	// So, (c.Money / 100) gives gold pieces. Each gold piece is effectively 1 coin.
	// Total coins = c.Money (assuming all currency is standardized to copper for weight calc, or convert to gold equiv)
	// A more accurate D&D 5e rule: 50 coins (of any denomination) weigh 1 lb.
	// Since Money is in copper pieces, we need to convert it to "number of coins".
	// For simplicity, let's assume 1 cp, 1 sp, 1 gp, 1 pp all count as 1 coin for weight purposes.
	// We'll approximate total coin weight by dividing total copper pieces by 100 to get an approximate "gold piece equivalent" count, then divide by COINS_PER_POUND.
	// This is a simplification, but better than nothing.
	goldPieces := float64(c.Money) / 100.0
	totalWeight += goldPieces / COINS_PER_POUND

	// 2. Weight from Items
	for _, item := range c.Items {
		totalWeight += ParseWeightString(item.Weight)
	}

	// 3. Weight from Weapons
	for _, weapon := range c.Weapons {
		totalWeight += ParseWeightString(weapon.Weight)
	}

	c.CarriedWeight = totalWeight
}
