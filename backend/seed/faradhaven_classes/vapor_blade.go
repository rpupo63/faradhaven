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
		DnDSkillFocus:  []string{"Stealth", "Deception"},
		Proficiencies:  "Simple weapons, Light armor, Daggers, Shortswords, Rapiers",
		SkillChoice:    []string{"Stealth", "Acrobatics", "Deception", "Perception"},
		Tools:          []string{"Poisoner's Kit", "Thieves' Tools"},
		SavingThrows:   []string{"Dexterity", "Intelligence"},
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
		ComponentPool: []string{
			// Forma (Shape)
			"Aura", "Self", "Projectile", "Zone", "Cone",
			// Scopus (Targeting)
			"Target", "Self",
			// Essentia (Matter & Energy)
			"Umbra", "Mortis", "Odor", "Aer", "Sonus", "Fear",
			// Actio (Verbs)
			"Pierce", "Destroy", "Bind",
			// Magnitudo (Modifiers)
			"Decrease", "Weak",
			// Logica (sequential links)
			"If", "Then", "Therefore",
		},
		ResourceDefinitions: []ResourceDefinitionSeed{
			{Key: "shadow_points", DisplayName: "Shadow Points", Category: "pool", Description: "Pool size equals Dexterity modifier + proficiency bonus (minimum 1). Spent on Venom, Shadow Step, Sneak Strike, etc.", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 1},
			{Key: "sneak_strike_uses", DisplayName: "Sneak Strike Uses", Category: "pool", Description: "Uses per short rest (equal to Dex Mod, min 1)", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 2},
			{Key: "death_mark_uses", DisplayName: "Death Mark", Category: "pool", Description: "1/Short Rest — mark a target for +2d6 Sneak Strike damage", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 3},
			{Key: "perfect_kill_uses", DisplayName: "Perfect Kill", Category: "pool", Description: "1/Long Rest: Force a Con save (DC 8+Dex+Prof) against a surprised target. Fail = 0 HP instantly.", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 4},
		},
	}
}

func vaporBladeLevelProgression() map[int]ClassLevelSeed {
	// Vapor Blade gets Extra Attack at 5 and Sneak Strike damage scales
	// Sneak Strike: 2d6 at level 2, +2d6 with Death Mark at level 10, unlimited at 20
	// Shadow Points max = Dex mod + proficiency (computed at runtime in ResourceService)
	levelProgression := make(map[int]ClassLevelSeed)
	for i := 1; i <= 20; i++ {
		// sneak_strike_uses: 0 at L1 (not unlocked), 3 at L2-4, scaling +1 per PB tier after that
		// death_mark_uses: 0 until L10, then 1
		var sneakStrikeUses int
		var deathMarkUses int
		switch {
		case i < 2:
			sneakStrikeUses = 0
		case i < 5:
			sneakStrikeUses = 3
		case i < 9:
			sneakStrikeUses = 4
		case i < 13:
			sneakStrikeUses = 5
		case i < 17:
			sneakStrikeUses = 6
		default:
			sneakStrikeUses = 7
		}
		if i >= 10 {
			deathMarkUses = 1
		}
		perfectKillUses := 0
		if i >= 17 {
			perfectKillUses = 1
		}
		entry := ClassLevelSeed{
			SneakAttackDice: 0,
			Resources: map[string]int{
				"sneak_strike_uses": sneakStrikeUses,
				"death_mark_uses":   deathMarkUses,
				"perfect_kill_uses": perfectKillUses,
			},
		}

		// Adjust SneakAttackDice based on level
		switch {
		case i >= 17:
			entry.ExtraAttackCount = 1
			entry.SneakAttackDice = 6
		case i >= 13:
			entry.ExtraAttackCount = 1
			entry.SneakAttackDice = 5
		case i == 10:
			entry.ExtraAttackCount = 1
			entry.SneakAttackDice = 4 // Death Mark adds +2d6
		case i == 9:
			entry.ExtraAttackCount = 1
			entry.SneakAttackDice = 3
		case i >= 5:
			entry.ExtraAttackCount = 1
			entry.SneakAttackDice = 2
		case i >= 2:
			entry.SneakAttackDice = 2
		}
		if i == 20 {
			entry.SneakAttackDice = 7
		}
		levelProgression[i] = entry
	}
	return levelProgression
}
func vaporBladeLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "Venom Coating", Description: "When you hit with a finesse weapon, you can spend 2 Shadow Points to apply a coating. Venom: Target takes 1d4 poison damage at the start of its next turn. Stacks intensity up to 3 times. Shadow: You have advantage on Stealth checks until the end of your next turn."},
			{Name: "Shadow Step", Description: "As a bonus action while in dim light or darkness, spend 5 Shadow Points to teleport up to 30 feet to another unoccupied space in dim light or darkness."},
		},
		2: {{Name: "Sneak Strike", Description: "Trigger: Advantage OR Unaware Target OR Flanking. Cost: 3 Shadow Points. Deal an extra 2d6 damage. Limit: Uses equal to Dex Mod (min 1). Regain all uses on Short Rest."}},
		3: {{Name: "Silent Blade", Description: "When you apply a coating, you can spend 1 additional Shadow Point to apply the **Silenced** status (preventing Verbal actions and spell components). The target must succeed on a Wisdom save (DC 8 + Prof + Dex) or be Silenced for 1 round."}},
		5: {{Name: "Extra Attack", Description: "You can attack twice when you take the Attack action. Each hit can apply a different coating."}, {Name: "Blade Flurry", Description: "When you use Shadow Step, your next attack before the end of your turn deals an extra 1d6 damage."}},
		6: {{Name: "Venom Mastery", Description: "Your Venom damage die increases to 2d4. Consecutive applications stack intensity (up to 3x total damage). Targets poisoned by you have disadvantage on saving throws against your other Vapor Blade features."}},
		7: {{Name: "Shadow Meld", Description: "As a bonus action in dim light or darkness, spend 5 Shadow Points to become invisible for 1 minute or until you attack or cast a spell."}},
		9: {{Name: "Lethal Precision", Description: "Critical hits and attacks against surprised targets (who have not yet acted) deal an extra 3d6 damage of your weapon's type for 5 Shadow Points."}, {Name: "Weak Point Analysis", Description: "You score a critical hit on a roll of 19 or 20."}},
		10: {{Name: "Death Mark", Description: "Bonus Action. Cost: 0 Shadow Points (Limit: 1/Short Rest). Mark a creature within 60 feet. For 1 minute, your Sneak Strike deals an extra 2d6 damage to it."}},
		11: {{Name: "Shadow Chain", Description: "When you use Shadow Step, you can spend 5 additional Shadow Points to bring one willing creature within 5 feet with you to the destination."}},
		13: {{Name: "Venom Burst", Description: "Trigger: Poisoned creature dies. Reaction. Cost: 10 Shadow Points. The creature explodes. Each creature within 10 feet takes 3d6 poison damage (Constitution save for half)."}},
		14: {{Name: "Ghost Walk", Description: "You can move through creatures and objects as if they were difficult terrain. If you end your turn inside an object, you take 1d10 force damage. You can spend 2 Shadow Points to avoid opportunity attacks while moving this way."}},
		15: {{Name: "Shadow Step Mastery", Description: "Shadow Step cost is reduced to 3 Shadow Points, and its range increases to 60 feet. You can use it twice per turn (one action, one bonus action)."}},
		17: {{Name: "Perfect Kill", Description: "Trigger: Hit a Surprised creature (has not acted). Cost: 15 Shadow Points. Limit: 1/Long Rest. Force a Constitution save (DC 8 + Dex + Prof). On failure, the target is reduced to 0 HP instantly. On success, it takes 10d10 necrotic damage."}, {Name: "Executioner", Description: "Your critical hit range increases to 18-20."}},
		18: {{Name: "Death Mark Mastery", Description: "Death Mark no longer requires a short rest (it is always available). Marked creatures have disadvantage on all Perception checks to find you."}},
		20: {{Name: "Shadow-Stalker's Eclipse", Description: "You regain 20 Shadow Points when you roll initiative and have none. While in dim light, you are invisible to creatures relying on darkvision. Sneak Strike no longer has a per-rest limit; you can use it on every hit as long as you have Shadow Points."}},
	}
}