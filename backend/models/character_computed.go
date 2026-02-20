package models

// CharacterComputedStats holds calculated values for the character sheet.
// Not stored in DB - computed at runtime from base stats + modifiers.
type CharacterComputedStats struct {
	MaxHP             int `json:"max_hp"`
	CurrentHP         int `json:"current_hp"`
	TempHP            int `json:"temp_hp"`
	ArmorClass        int `json:"armor_class"`
	Initiative        int `json:"initiative"`
	ProficiencyBonus  int `json:"proficiency_bonus"`
	SpellSaveDC       int `json:"spell_save_dc"`
	SpellAttackBonus  int `json:"spell_attack_bonus"`
	MaxSpellPoints    int `json:"max_spell_points"`
	PassivePerception int `json:"passive_perception"`

	// Attack bonuses
	MeleeAttackBonus  int `json:"melee_attack_bonus"`
	RangedAttackBonus int `json:"ranged_attack_bonus"`

	// Hit dice tracking
	HitDiceTotal     int `json:"hit_dice_total"`     // = Level
	HitDiceRemaining int `json:"hit_dice_remaining"` // = Level - HitDiceUsed
	HitDie           int `json:"hit_die"`            // e.g., 10 for d10

	// Ability modifiers
	StrengthMod     int `json:"strength_mod"`
	DexterityMod    int `json:"dexterity_mod"`
	ConstitutionMod int `json:"constitution_mod"`
	IntelligenceMod int `json:"intelligence_mod"`
	WisdomMod       int `json:"wisdom_mod"`
	CharismaMod     int `json:"charisma_mod"`
}

// AbilityModifier calculates the modifier for an ability score (D&D 5e formula).
func AbilityModifier(score int) int {
	return (score - 10) / 2
}

// ComputeStats calculates all derived stats for a character.
// Call this after loading the character with relationships populated.
func (c *Character) ComputeStats() {
	stats := &CharacterComputedStats{}

	// Calculate ability modifiers
	stats.StrengthMod = AbilityModifier(c.Strength)
	stats.DexterityMod = AbilityModifier(c.Dexterity)
	stats.ConstitutionMod = AbilityModifier(c.Constitution)
	stats.IntelligenceMod = AbilityModifier(c.Intelligence)
	stats.WisdomMod = AbilityModifier(c.Wisdom)
	stats.CharismaMod = AbilityModifier(c.Charisma)

	// Proficiency bonus (D&D 5e: starts at +2, increases every 4 levels)
	stats.ProficiencyBonus = 2 + (c.Level-1)/4

	// Initiative = Dexterity modifier
	stats.Initiative = stats.DexterityMod

	// Passive Perception = 10 + Wisdom modifier (+ proficiency if proficient)
	stats.PassivePerception = 10 + stats.WisdomMod

	// HP calculation requires class data
	if c.CurrentClassLevel != nil {
		// Max spell points from class level
		stats.MaxSpellPoints = c.CurrentClassLevel.MaxSpellPoints
	}

	// Spell save DC and attack bonus depend on primary ability
	primaryMod := stats.IntelligenceMod // default
	if c.Class.PrimaryAbility == "wisdom" {
		primaryMod = stats.WisdomMod
	} else if c.Class.PrimaryAbility == "charisma" {
		primaryMod = stats.CharismaMod
	}
	stats.SpellSaveDC = 8 + stats.ProficiencyBonus + primaryMod
	stats.SpellAttackBonus = stats.ProficiencyBonus + primaryMod

	// --- Armor Class (AC) Calculation ---
	baseAC := 10
	dexModForAC := stats.DexterityMod
	shieldsAllowed := true // Default

	// Check for Unarmored Defense features
	// This method assumes that c.Class.LevelFeatures and c.Race.Traits are preloaded.
	unarmoredBaseAC, unarmoredAbilityMods, unarmoredShieldsAllowed := c.GetUnarmoredDefenseACBonus()

	if c.EquippedArmor == nil && unarmoredBaseAC > 0 {
		// Apply Unarmored Defense if no armor is equipped and a feature provides it
		baseAC = unarmoredBaseAC
		if unarmoredAbilityMods["dexterity"] {
			baseAC += stats.DexterityMod
		}
		if unarmoredAbilityMods["constitution"] {
			baseAC += stats.ConstitutionMod
		}
		if unarmoredAbilityMods["wisdom"] {
			baseAC += stats.WisdomMod
		}
		if unarmoredAbilityMods["charisma"] {
			baseAC += stats.CharismaMod
		}
		dexModForAC = 0 // Dex modifier already incorporated into baseAC if it was part of Unarmored Defense
		shieldsAllowed = unarmoredShieldsAllowed
	} else if c.EquippedArmor != nil {
		// Existing armor calculation logic if armor is equipped
		if c.EquippedArmor.BaseAC != nil {
			baseAC = *c.EquippedArmor.BaseAC
		}
		if c.EquippedArmor.ArmorType != nil {
			switch *c.EquippedArmor.ArmorType {
			case "Heavy":
				dexModForAC = 0 // Heavy armor gets no dexterity bonus
			case "Medium":
				if dexModForAC > 2 {
					dexModForAC = 2 // Medium armor has a max dex bonus of +2
				}
			case "Light":
				// No change to dexModForAC
			default:
				// "Unarmored" or other cases, do nothing special
			}
		}
	} else {
		// Default unarmored AC: 10 + Dex modifier if no armor and no specific Unarmored Defense feature
		baseAC = 10
		// dexModForAC already set to stats.DexterityMod
	}

	// Calculate final AC from base + any remaining dex bonus
	stats.ArmorClass = baseAC + dexModForAC

	// Add shield bonus if a shield is equipped and allowed by unarmored defense (if active)
	if c.EquippedShield != nil && c.EquippedShield.BaseAC != nil && shieldsAllowed {
		stats.ArmorClass += *c.EquippedShield.BaseAC
	}

	// Attack bonuses (proficiency + ability modifier)
	stats.MeleeAttackBonus = stats.ProficiencyBonus + stats.StrengthMod
	stats.RangedAttackBonus = stats.ProficiencyBonus + stats.DexterityMod

	// HP from character (persisted values)
	stats.MaxHP = c.MaxHP
	stats.CurrentHP = c.CurrentHP
	stats.TempHP = c.TempHP

	// Hit dice tracking
	stats.HitDiceTotal = c.Level
	stats.HitDiceRemaining = c.Level - c.HitDiceUsed
	if c.Class.HitDie > 0 {
		stats.HitDie = c.Class.HitDie
	}

	c.ComputedStats = stats
}
