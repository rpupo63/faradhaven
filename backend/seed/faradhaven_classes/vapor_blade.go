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
		1:  {SneakAttackDice: 0},
		2:  {SneakAttackDice: 2},                      // Sneak Strike: 2d6
		3:  {SneakAttackDice: 2},
		4:  {SneakAttackDice: 2},
		5:  {ExtraAttackCount: 1, SneakAttackDice: 2}, // Extra Attack
		6:  {ExtraAttackCount: 1, SneakAttackDice: 2},
		7:  {ExtraAttackCount: 1, SneakAttackDice: 2},
		8:  {ExtraAttackCount: 1, SneakAttackDice: 2},
		9:  {ExtraAttackCount: 1, SneakAttackDice: 3}, // Lethal Precision adds damage
		10: {ExtraAttackCount: 1, SneakAttackDice: 4}, // Death Mark: +2d6
		11: {ExtraAttackCount: 1, SneakAttackDice: 4},
		12: {ExtraAttackCount: 1, SneakAttackDice: 4},
		13: {ExtraAttackCount: 1, SneakAttackDice: 5}, // Venom Burst level
		14: {ExtraAttackCount: 1, SneakAttackDice: 5},
		15: {ExtraAttackCount: 1, SneakAttackDice: 5},
		16: {ExtraAttackCount: 1, SneakAttackDice: 5},
		17: {ExtraAttackCount: 1, SneakAttackDice: 6}, // Perfect Kill level
		18: {ExtraAttackCount: 1, SneakAttackDice: 6},
		19: {ExtraAttackCount: 1, SneakAttackDice: 6},
		20: {ExtraAttackCount: 1, SneakAttackDice: 7}, // Shadow-Stalker's Eclipse
	}
}

func vaporBladeLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "Venom Coating", Description: "When you hit with a finesse weapon, you can spend 2 spell points to apply a coating. Venom: Target takes 1d4 poison damage at the start of its next turn. Shadow: You have advantage on Stealth checks until the end of your next turn."},
			{Name: "Shadow Step", Description: "As a bonus action while in dim light or darkness, spend 5 spell points to teleport up to 30 feet to another unoccupied space in dim light or darkness."},
		},
		2:  {{Name: "Sneak Strike", Description: "When you have advantage on an attack or the target is surprised, you can spend 3 spell points to deal an extra 2d6 damage. You can use this a number of times equal to your Dexterity modifier per short rest."}},
		3:  {{Name: "Silent Blade", Description: "When you apply a coating, you can spend 1 additional spell point to add the Silence effect. The target must succeed on a Wisdom save (DC 8 + Prof + Dex) or be unable to speak or cast verbal spells for 1 round."}},
		5:  {{Name: "Extra Attack", Description: "You can attack twice when you take the Attack action. Each hit can apply a different coating."}, {Name: "Blade Flurry", Description: "When you use Shadow Step, your next attack before the end of your turn deals an extra 1d6 damage."}},
		6:  {{Name: "Venom Mastery", Description: "Your Venom damage increases to 2d4 and lasts for 2 rounds. Targets poisoned by you have disadvantage on saving throws against your other Vapor Blade features."}},
		7:  {{Name: "Shadow Meld", Description: "As a bonus action in dim light or darkness, spend 5 spell points to become invisible for 1 minute or until you attack or cast a spell."}},
		9:  {{Name: "Lethal Precision", Description: "Critical hits and attacks against surprised targets deal an extra 3d6 damage of your weapon's type for 5 spell points."}, {Name: "Weak Point Analysis", Description: "You score a critical hit on a roll of 19 or 20."}},
		10: {{Name: "Death Mark", Description: "As a bonus action, mark a creature within 60 feet. For 1 minute, your Sneak Strike deals an extra 2d6 damage to it. Once per short rest."}},
		11: {{Name: "Shadow Chain", Description: "When you use Shadow Step, you can spend 5 additional spell points to bring one willing creature within 5 feet with you to the destination."}},
		13: {{Name: "Venom Burst", Description: "When a creature poisoned by you dies, you can spend 10 spell points as a reaction to make it explode. Each creature within 10 feet takes 3d6 poison damage (Constitution save for half)."}},
		14: {{Name: "Ghost Walk", Description: "You can move through creatures and objects as if they were difficult terrain. If you end your turn inside an object, you take 1d10 force damage. You can spend 2 spell points to avoid opportunity attacks while moving this way."}},
		15: {{Name: "Shadow Step Mastery", Description: "Shadow Step cost is reduced to 3 spell points, and its range increases to 60 feet. You can use it twice per turn (one action, one bonus action)."}},
		17: {{Name: "Perfect Kill", Description: "Once per long rest, when you hit a surprised creature, you can spend 15 spell points to force a Constitution save (DC 8 + Dex + Prof). On failure, the target is reduced to 0 HP instantly. On success, it takes 10d10 necrotic damage."}, {Name: "Executioner", Description: "Your critical hit range increases to 18-20."}},
		18: {{Name: "Death Mark Mastery", Description: "Death Mark no longer requires a short rest (it is always available). Marked creatures have disadvantage on all Perception checks to find you."}},
		20: {{Name: "Shadow-Stalker's Eclipse", Description: "You regain 20 spell points when you roll initiative and have none. While in dim light, you are invisible to creatures relying on darkvision. Sneak Strike no longer has a per-rest limit; you can use it on every hit as long as you have spell points."}},
	}
}
