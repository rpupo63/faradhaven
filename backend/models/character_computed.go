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

	// Base AC (unarmored) = 10 + Dex modifier
	stats.ArmorClass = 10 + stats.DexterityMod

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
