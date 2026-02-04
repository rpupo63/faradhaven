package faradhaven_classes

// PistonBrawler returns the Piston Brawler class seed
func PistonBrawler() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Piston Brawler",
		Description:    "Wear massive steam-powered gauntlets and fight unarmored, relying on raw Strength and Constitution for defense. Punch through enemies with explosive force while shielding yourself with pressurized steam.",
		HitDie:         10,
		PrimaryAbility: "strength",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/piston_brawler.jpg",
		Archetype:      "Front-line Tank",
		Concept:        "A front-line tank that relies on 'unarmored defense.' Max STR and CON—you don't need armor when your fists and raw toughness are your shield.",
		ClassFeatures: []string{
			"Atlas Gauntlets: Massive steam-powered gloves that amplify your punches and provide shields. Your unarmed strikes count as martial weapons and deal enhanced damage.",
			"Unarmored Defense: When not wearing armor, your AC equals 10 + your Dexterity modifier + your Constitution modifier.",
		},
		DnDSkillFocus: []string{"Athletics", "Intimidation"},
		Proficiencies: "Simple Weapons, Martial Weapons (Atlas Gauntlets), Unarmored Defense",
		SkillChoice:   []string{"Athletics", "Acrobatics", "Survival", "Intimidation"},
		Tools:         []string{"Smith's Tools (for gauntlet maintenance)"},
		SavingThrows:  []string{"Strength", "Constitution"},
		StartingEquip: []string{"Atlas Gauntlets", "Traveler's clothes", "Prybar", "Repair kit"},
		LevelFeatures: pistonBrawlerLevelFeatures(),
	}
}

func pistonBrawlerLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Atlas Gauntlets, Unarmored Defense — You gain the core class features. Your unarmed strikes with the Atlas Gauntlets deal 1d8 bludgeoning damage (instead of 1) and count as magical. As a bonus action, you can spend 2 spell points to gain temporary hit points equal to 1d6 + your Constitution modifier (gauntlet shields). You have one use, regained on a short rest.",
		2:  "Steam-forged Grip — You channel spell points through your gauntlets. When you make an opportunity attack with your Atlas Gauntlets, you can spend 3 spell points to add your Strength modifier to the damage and prevent the target from taking reactions until the start of their next turn.",
		3:  "Heavy Stance — You plant your feet and absorb impact. As a reaction when you take bludgeoning, piercing, or slashing damage, you can spend 5 spell points to reduce the damage by 1d10 + your Constitution modifier. You have one use, regained on a short rest.",
		4:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Strength and Constitution cannot exceed 20.",
		5:  "Extra Attack — You can attack twice, instead of once, whenever you take the Attack action on your turn. Your gauntlet strikes define you as a juggernaut.",
		6:  "Gauntlet Shield Surge — When you use your gauntlet shields, you can spend 5 additional spell points to extend the temporary hit points to 2d6 + twice your Constitution modifier. While you have these temporary hit points, you have +1 AC.",
		7:  "Overpressure Punch — When you hit with your Atlas Gauntlets, you can spend 3 spell points to vent steam in a burst. The target must succeed on a Strength save or be pushed 10 feet and take 1d6 fire damage. Creatures within 5 feet of the target take half the fire damage.",
		8:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Strength and Constitution cannot exceed 20.",
		9:  "Explosive Feedback — When you hit with your Atlas Gauntlets, you can spend 8 spell points as a reaction to detonate the impact. Each creature within 10 feet of the target (including the target) takes 2d6 fire damage. You have resistance to this damage.",
		10: "Overpressure Surge — You overcharge your gauntlets' steam vents. Once per short rest, as an action, you release a 20-foot cone of scalding steam. Each creature in the area takes 4d6 fire damage and must succeed on a Constitution save or be incapacitated until the end of your next turn.",
		11: "Unstoppable Momentum — When you move at least 15 feet straight toward a creature before making a melee attack with your Atlas Gauntlets, you can spend 5 spell points to add 2d6 damage to the attack and push the target 15 feet on a hit. If they collide with a solid object, they take an additional 1d6 bludgeoning damage.",
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Strength and Constitution cannot exceed 20.",
		13: "Shattering Blow — When you hit with your Atlas Gauntlets, you can spend 5 spell points to reduce the target's AC by 2 for 1 minute. This stacks up to a maximum of -4.",
		14: "Extra Attack (2) — You can attack three times, instead of twice, whenever you take the Attack action on your turn.",
		15: "Reinforced Frame — Your body and gauntlets are one. You have resistance to fire damage. When you use your gauntlet shields, you can choose to also gain +2 AC until the start of your next turn.",
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Strength and Constitution cannot exceed 20.",
		17: "Devastating Impact — When you score a critical hit with your Atlas Gauntlets, you can spend 10 spell points to add an extra damage die and force the target to make a Constitution save or be knocked prone and lose their reaction until the end of their next turn.",
		18: "Gauntlet Mastery — Your gauntlet shields no longer require a short rest; you have two uses, regained on a long rest. Your Overpressure Surge can be used once per long rest. Your unarmed strikes with the Atlas Gauntlets now deal 1d10 bludgeoning damage.",
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Strength and Constitution cannot exceed 20.",
		20: "Hextech Brawler's Fury — You embody raw industrial might. Once per long rest, when you roll initiative, you can enter a state of fury for 1 minute. While in this state: you have advantage on Strength checks and saving throws; your Atlas Gauntlet attacks deal an extra 2d6 damage; your gauntlet shields grant double temporary hit points; and you cannot be knocked prone or moved against your will. When the minute ends, you gain one level of exhaustion.",
	}
}
