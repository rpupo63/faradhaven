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
				Instruction: "Choose your weapon to upgrade",
				Options: []EquipmentOptionSeed{
					{Description: "A Warhammer", WeaponNames: []string{"Warhammer"}},
					{Description: "A Longsword", WeaponNames: []string{"Longsword"}},
					{Description: "A Battleaxe", WeaponNames: []string{"Battleaxe"}},
				},
			},
			{
				Instruction: "Choose your shield or sidearm",
				Options: []EquipmentOptionSeed{
					{Description: "A Shield", ItemNames: []string{"Standard Shield"}},
					{Description: "Two Light Hammers", WeaponNames: []string{"Light Hammer", "Light Hammer"}},
				},
			},
		},
		LevelFeatures:    pistonBrawlerLevelFeatures(),
		LevelProgression: pistonBrawlerLevelProgression(),
	}
}

func pistonBrawlerLevelProgression() map[int]ClassLevelSeed {
	return map[int]ClassLevelSeed{
		5: {ExtraAttackCount: 1},
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
