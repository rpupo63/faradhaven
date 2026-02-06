package faradhaven_classes

// PistonBrawler returns the Piston Brawler class seed
func PistonBrawler() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Piston Brawler",
		Description:    "To accommodate any weapon type—be it a Greatsword, a pair of Daggers, or a Longbow—this class pivots identity from a specific item to the Piston Core. This mechanical engine is integrated into the weapon of your choice, converting kinetic energy into magical output.",
		HitDie:         10,
		PrimaryAbility: "Strength",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/piston_brawler.jpg",
		Archetype:      "Magical Striker / Tech Knight",
		Concept:        "A weapon-agnostic engineer who attaches a 'Piston Core' to their weapon of choice, managing a pool of 'Stability' to fuel magical strikes and utility spells.",
		ClassFeatures: []string{
			"Piston Core: A device attached to a weapon that converts kinetic energy into magical output. Adds Intelligence modifier to damage while Active.",
			"Stability System: A pool of points (10 + Level*2) that degrades on hit but fuels spells and abilities.",
			"Integrated Spell Logic: Instead of slots, the weapon has 'Fixed Spells' programmed into it, fueled by Stability.",
		},
		DnDSkillFocus:    []string{"Arcana", "Athletics"},
		Proficiencies:    "Simple Weapons, Martial Weapons, Light Armor, Medium Armor, Shields",
		SkillChoice:      []string{"Arcana", "Athletics", "Investigation", "Sleight of Hand", "Intimidation"},
		Tools:            []string{"Smith's Tools", "Tinker's Tools"},
		SavingThrows:     []string{"Constitution", "Intelligence"},
		AutomaticEquipNames: []string{"Piston Core assembly kit"},
		AutomaticItemNames:  []string{"Scale mail", "Explorer's pack", "Tinker's tools"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your primary Piston-weapon",
				Options: []EquipmentOptionSeed{
					{Description: "Atlas Gauntlets (Heavy Piston Fists)", WeaponNames: []string{"Atlas Gauntlets"}},
					{Description: "A Greatsword (Heavy Blade)", WeaponNames: []string{"Greatsword"}},
					{Description: "A Warhammer (Versatile Crusher)", WeaponNames: []string{"Warhammer"}},
				},
			},
			{
				Instruction: "Choose your maintenance gear",
				Options: []EquipmentOptionSeed{
					{Description: "A Shield and Pressure Gauge", ItemNames: []string{"Standard Shield", "Stability Pressure Gauge"}},
					{Description: "Piston Core Assembly Kit", ItemNames: []string{"Piston Core Assembly Kit"}},
				},
			},
		},
		LevelFeatures:       pistonBrawlerLevelFeatures(),
		LevelProgression:    pistonBrawlerLevelProgression(),
		ResourceType:        "stability",
		ResourceName:        "Stability",
		ResourceRestoreType: "special", // Reset by weapon/action, not rest
		WeaponRequirement: &WeaponRequirementSeed{
			SelectionLevel:    1,
			ModifierType:      "piston_core",
			Description:       "Choose a weapon to integrate your Piston Core. This mechanical device converts kinetic energy into magical output, permanently bonding with your chosen weapon.",
			AllowedCategories: []string{"Martial Melee", "Simple Melee", "Martial Ranged", "Simple Ranged"},
		},
	}
}

func pistonBrawlerLevelProgression() map[int]ClassLevelSeed {
	// Stability formula: 10 + (Level * 2), seeded per-level for display
	// MaxSpellLevel tracks the spell level cap for Fixed Spells
	return map[int]ClassLevelSeed{
		1:  {MaxStability: 12, MaxSpellLevel: 1},                       // 10 + (1*2) = 12
		5:  {MaxStability: 20, MaxSpellLevel: 2, ExtraAttackCount: 1},  // 10 + (5*2) = 20
		9:  {MaxStability: 28, MaxSpellLevel: 3, ExtraAttackCount: 1},  // 10 + (9*2) = 28
		13: {MaxStability: 36, MaxSpellLevel: 4, ExtraAttackCount: 1},  // 10 + (13*2) = 36
		17: {MaxStability: 44, MaxSpellLevel: 5, ExtraAttackCount: 1},  // 10 + (17*2) = 44
		20: {MaxStability: 50, MaxSpellLevel: 9, ExtraAttackCount: 1},  // Perfect Machine: no max (use 50 as display)
	}
}

func pistonBrawlerLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Piston Core & Stability — You attach a Piston Core to a weapon. Bonus Action to engage. While Active, add Int mod to damage. Stability Pool: 10 + (Level * 2). Light/Finesse lose 1 Stability/hit; Heavy/Two-Handed lose 2. At 0 Stability, weapon becomes 'Unruly'. Cast Fixed Spells (Max Level 1) by spending Stability equal to spell level.",
		2:  "Overdrive — When you hit with your weapon, you can expend 5 Stability to add 2d6 Force damage. This represents pushing the pistons past safety limits.",
		4:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Strength and Intelligence cannot exceed 20.",
		5:  "Durable Chassis & Extra Attack — If your weapon breaks (0 Stability), it functions as a non-magical weapon instead of becoming unusable. You can Attack twice instead of once when taking the Attack action. Max Spell Level increases to 2nd.",
		7:  "Kinetic Calibration — When Stability is < 50%, gain Advantage on Strength and Dexterity saving throws.",
		8:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Strength and Intelligence cannot exceed 20.",
		9:  "High-Pressure Output — Max Spell Level increases to 3rd (e.g., Haste, Fireball).",
		10: "Reinforced Valves — You only roll for Malfunction on every second hit while at 0 Stability.",
		11: "Feedback Loop — Casting a Fixed Spell restores Stability equal to spell level on your next hit. Overdrive damage increases to 3d6.",
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Strength and Intelligence cannot exceed 20.",
		13: "Advanced Manipulation — Max Spell Level increases to 4th (e.g., Staggering Smite).",
		15: "Emergency Venting — Reaction to prevent Malfunction, resetting Stability to 1.",
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Strength and Intelligence cannot exceed 20.",
		17: "Masterwork Engineering — Max Spell Level increases to 5th (e.g., Steel Wind Strike).",
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Strength and Intelligence cannot exceed 20.",
		20: "The Perfect Machine — No Stability maximum. Regain 5 Stability at start of turn. Only breaks if you 'Purge' (AoE damage, reset Stability to 0).",
	}
}
