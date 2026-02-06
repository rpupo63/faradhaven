package faradhaven_classes

// VaporBlade returns the Vapor Blade class seed
func VaporBlade() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Vapor Blade",
		Description:    "Coat your blades in venom and slip through shadows to strike before enemies know you're there. Master the art of assassination with poison damage over time, shadow teleportation, and lethal precision.",
		HitDie:         8,
		PrimaryAbility: "dexterity",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/vapor_blade.jpg",
		Archetype:      "Stealth / Assassin / Melee DPS",
		Concept:        "A killer who blends magic with blades. They coat weapons in venom, move through shadows, and strike before targets know they're there.",
		ClassFeatures: []string{
			"Venom Coating: You apply component effects to your blades. Venom deals poison damage over time; Shadow grants advantage on Stealth and muffles your presence.",
			"Shadow Step: When you are in dim light or darkness, you can spend spell points to slip between shadows and emerge unseen.",
		},
		DnDSkillFocus:    []string{"Stealth", "Deception"},
		Proficiencies:    "Simple weapons, Light armor, Daggers, Shortswords, Rapiers",
		SkillChoice:      []string{"Stealth", "Acrobatics", "Deception", "Perception"},
		Tools:            []string{"Poisoner's Kit", "Thieves' Tools"},
		SavingThrows:     []string{"Dexterity", "Intelligence"},
		AutomaticEquipNames: []string{"Dark cloak"},
		AutomaticItemNames:  []string{"Leather armor", "Poisoner's kit", "Thieves' Tools"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your blade",
				Options: []EquipmentOptionSeed{
					{Description: "Folding Blade (Concealable)", WeaponNames: []string{"Folding Blade"}},
					{Description: "Two Daggers and Smoke Pellet", WeaponNames: []string{"Dagger", "Dagger"}, ItemNames: []string{"Vapor-Smoke Pellet"}},
				},
			},
			{
				Instruction: "Choose your stealth gear",
				Options: []EquipmentOptionSeed{
					{Description: "Dark Cloak", ItemNames: []string{"Dark Cloak"}},
					{Description: "Burglar's Pack", ItemNames: []string{"Burglar's Pack"}},
				},
			},
		},
		LevelFeatures:    vaporBladeLevelFeatures(),
		LevelProgression: vaporBladeLevelProgression(),
		ResourceType:        "spell_points",
		ResourceName:        "Spell Points",
		ResourceRestoreType: "long_rest",
	}
}

func vaporBladeLevelProgression() map[int]ClassLevelSeed {
	// Vapor Blade gets Extra Attack at 5 and Sneak Strike damage scales
	// Sneak Strike: 2d6 at level 2, +2d6 with Death Mark at level 10, unlimited at 20
	return map[int]ClassLevelSeed{
		2:  {SneakAttackDice: 2},                      // Sneak Strike: 2d6
		5:  {ExtraAttackCount: 1, SneakAttackDice: 2}, // Extra Attack
		9:  {ExtraAttackCount: 1, SneakAttackDice: 3}, // Lethal Precision adds damage
		10: {ExtraAttackCount: 1, SneakAttackDice: 4}, // Death Mark: +2d6
		13: {ExtraAttackCount: 1, SneakAttackDice: 5}, // Venom Burst level
		17: {ExtraAttackCount: 1, SneakAttackDice: 6}, // Perfect Kill level
		20: {ExtraAttackCount: 1, SneakAttackDice: 7}, // Shadow-Stalker's Eclipse
	}
}

func vaporBladeLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Venom Coating, Shadow Step — You gain the core class features. When you hit with a finesse weapon, you can spend 2 spell points to apply Venom (target takes 1d4 poison damage at the start of their next turn) or Shadow (you have advantage on Stealth checks until the end of your next turn). As a bonus action in dim light or darkness, you can spend 5 spell points to teleport up to 30 feet to an unoccupied space you can see that is also in dim light or darkness.",
		2:  "Sneak Strike — When you attack a creature that hasn't taken a turn in combat yet, or when you have advantage on the attack roll, you can spend 3 spell points to add 2d6 damage to the hit. You can use this feature a number of times per short rest equal to your Dexterity modifier.",
		3:  "Silent Blade — When you apply the Silence modifier to an attack, the target cannot shout or cast spells with verbal components until the end of their next turn. You have advantage on Deception checks to pass unnoticed in crowds.",
		4:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Dexterity cannot exceed 20.",
		5:  "Extra Attack — When you take the Attack action, you can attack twice instead of once. Your second attack can apply a different component coating.",
		6:  "Venom Mastery — Your Venom damage increases to 2d4 and lasts until the end of the target's next turn. When you apply Venom to a creature that is already poisoned by you, they have disadvantage on the saving throw to end the effect.",
		7:  "Shadow Meld — While in dim light or darkness, you can spend 5 spell points as a bonus action to become invisible until the start of your next turn or until you attack or cast a spell. Moving more than half your speed ends the effect early.",
		8:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Dexterity cannot exceed 20.",
		9:  "Lethal Precision — When you score a critical hit or hit a surprised creature, you can spend 5 spell points to add the Lethal modifier: the target takes an extra 3d6 damage and must succeed on a Constitution save or be poisoned for 1 minute.",
		10: "Death Mark — Once per short rest, you can mark a creature you can see. For 1 minute, your Sneak Strike damage against that creature increases by 2d6, and you know their direction and distance from you while they are within 120 feet.",
		11: "Shadow Chain — When you use Shadow Step, you can spend 5 additional spell points to bring one willing creature of your size or smaller with you. They appear in an unoccupied space within 5 feet of your destination.",
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Dexterity cannot exceed 20.",
		13: "Venom Burst — When a creature affected by your Venom dies, you can spend 10 spell points as a reaction to cause their blood to spray. Each creature within 10 feet must succeed on a Constitution save or take 3d6 poison damage and be poisoned for 1 round.",
		14: "Ghost Walk — You can move through other creatures' spaces and treat them as difficult terrain. When you do so, you can spend 2 spell points per creature to leave them unable to take reactions until the end of their next turn.",
		15: "Shadow Step Mastery — Your Shadow Step costs 3 spell points and has a range of 60 feet. You can use it twice per round (once as an action, once as a bonus action), but only once per turn.",
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Dexterity cannot exceed 20.",
		17: "Perfect Kill — Once per long rest, when you hit a creature that is surprised or incapacitated with a melee weapon attack, you can spend 15 spell points to make the attack an automatic critical hit. If the target has 100 hit points or fewer, they must succeed on a Constitution save (DC 8 + your Dexterity modifier + your proficiency bonus) or die.",
		18: "Death Mark Mastery — Your Death Mark no longer requires a short rest. You can use it once per long rest. Additionally, marked creatures have disadvantage on Perception checks to detect you, and your Shadow Step has no range limit against them (you must still teleport to dim light or darkness).",
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Dexterity cannot exceed 20.",
		20: "Shadow-Stalker's Eclipse — You are one with the dark. You regain 20 spell points when you roll initiative and have none remaining. While in dim light or darkness, you are invisible to creatures that rely on darkvision to see you. Your Sneak Strike can be used on every hit (no limit per rest), and when you reduce a creature to 0 hit points with Sneak Strike, you can use Shadow Step as a free action.",
	}
}
