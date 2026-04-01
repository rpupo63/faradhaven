package faradhaven_classes

// Syllogist returns The Syllogist class seed — a null-caster conductor of action economy and Logica chains for allies (Faradhaven: no spell points; potency from Protocol Charges).
func Syllogist() FaradhavenClassSeed {
	archetypeLevel := 3
	return FaradhavenClassSeed{
		Name:        "The Syllogist",
		Description: "Trained tacticians who wield notation magic—a cognitivist tradition that treats battle as ordered proof, not evocation. You are a null-caster in Faradhaven terms: this class grants no spell points and no spell levels; your formulae cannot satisfy the spell forge alone (your pool omits Forma and Essentia). You instead spend Protocol Charges to donate Logica and Magnitudo, reorder allies’ chains, share components, and steal moments from the round.",
		HitDie:         8,
		PrimaryAbility: "Intelligence",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/syllogist.jpg",
		Archetype:      "Null-caster support / Action economy / Logica",
		Concept: "Power source: Notation—patterns of causality written in gesture, voice, and scraped formulae; it is not divine, not primal, and not spell-slot arcana. Origin: guild academies, military tribunals, and proof-houses that denied that only a full formula can change the world—you change whose sequence fires and when. You advance by study (Intelligence as your spellcasting-equivalent ability); you endure stress with clarity of mind (Wisdom saving throws).\n\n" +
			"Spell slot progression (5e frame): null-caster—potency comes from Protocol Charges and features, not Faradhaven spell points. By 2nd level your mechanical identity is fixed: you lend components (Lend Axiom) and react to Logica casts (Interleave); you are ally-focused permission and timing, not personal spell damage.",
		DnDSkillFocus: []string{"Arcana", "History"},
		Proficiencies: "Simple weapons, light armor",
		SkillChoice: []string{
			"Arcana", "History", "Insight", "Investigation", "Perception",
		},
		Tools:            []string{"Calligrapher's Supplies"},
		SavingThrows:     []string{"Intelligence", "Wisdom"},
		AutomaticItemNames: []string{"Scholar's robes (light armor)", "Scholar's pack", "Calligrapher's Supplies"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your simple weapon",
				Options: []EquipmentOptionSeed{
					{Description: "A quarterstaff", WeaponNames: []string{"Quarterstaff"}},
					{Description: "A dagger", WeaponNames: []string{"Dagger"}},
				},
			},
		},
		LevelFeatures:       syllogistLevelFeatures(),
		LevelProgression:    syllogistLevelProgression(),
		ArchetypeLevel:      &archetypeLevel,
		Archetypes:          syllogistArchetypes(),
		ResourceDefinitions: syllogistResourceDefinitions(),
		ComponentPool: []string{
			"If", "Then", "Therefore",
			"Increase", "Decrease", "Strong", "Weak", "Extreme",
		},
	}
}

func syllogistResourceDefinitions() []ResourceDefinitionSeed {
	return []ResourceDefinitionSeed{
		{Key: "protocol_charges", DisplayName: "Protocol Charges", Category: "pool", Description: "Your expendable notation budget for class features. You regain all expended charges when you finish a long rest.", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 1},
		{Key: "grant_hustle_uses", DisplayName: "Grant Hustle", Category: "pool", Description: "Uses of Grant Hustle per short rest (Proficiency Bonus).", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 2},
		{Key: "resonant_strong_uses", DisplayName: "Resonant Strong", Category: "pool", Description: "Uses of Resonant Strong per long rest (Proficiency Bonus).", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 3},
		{Key: "efficient_theorem_uses", DisplayName: "Efficient Theorem", Category: "pool", Description: "Uses of Efficient Theorem extra die per long rest (Proficiency Bonus).", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 4},
		{Key: "refuting_gesture_uses", DisplayName: "Refuting Gesture", Category: "pool", Description: "Uses of Refuting Gesture per long rest (Intelligence modifier).", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 5},
		{Key: "simultaneous_lemma_uses", DisplayName: "Simultaneous Lemma", Category: "pool", Description: "1/Long Rest: Action to grant an ally a bonus action spell cast.", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 6},
		{Key: "rapid_axiom_uses", DisplayName: "Rapid Axiom", Category: "pool", Description: "1/Short Rest: Use Lend Axiom as a bonus action.", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 7},
		{Key: "distributed_proof_uses", DisplayName: "Distributed Proof", Category: "pool", Description: "Uses of Distributed Proof per long rest (Proficiency Bonus).", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 8},
		{Key: "perfected_interleave_uses", DisplayName: "Perfected Interleave", Category: "pool", Description: "1/Short Rest: Use Interleave without spending a Protocol Charge.", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 9},
		{Key: "reductio_uses", DisplayName: "Reductio ad Absurdum", Category: "pool", Description: "1/Long Rest: Impose disruption or grant impulse after an ally casts a Logica spell.", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 10},
		{Key: "penultimate_proof_uses", DisplayName: "Penultimate Proof", Category: "pool", Description: "1/Long Rest: Regain 2 Protocol Charges when rolling initiative with 0 remaining.", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 11},
		{Key: "qed_uses", DisplayName: "Q.E.D.", Category: "pool", Description: "1/Long Rest: Regain charges based on support provided since last long rest.", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 12},
	}
}

func syllogistProtocolMax(level int) int {
	switch {
	case level <= 4:
		return 5
	case level <= 10:
		return 7
	case level <= 12:
		return 9
	case level <= 16:
		return 10
	case level <= 17:
		return 11
	case level <= 19:
		return 12
	default:
		return 14
	}
}

func syllogistLevelProgression() map[int]ClassLevelSeed {
	m := make(map[int]ClassLevelSeed, 20)
	for level := 1; level <= 20; level++ {
		pb := proficiencyByLevel(level)
		res := map[string]int{
			"protocol_charges":          syllogistProtocolMax(level),
			"grant_hustle_uses":         0,
			"resonant_strong_uses":      0,
			"efficient_theorem_uses":    0,
			"refuting_gesture_uses":     0,
			"simultaneous_lemma_uses":   0,
			"rapid_axiom_uses":          0,
			"distributed_proof_uses":    0,
			"perfected_interleave_uses": 0,
			"reductio_uses":             0,
			"penultimate_proof_uses":    0,
			"qed_uses":                  0,
		}

		if level >= 3 {
			res["grant_hustle_uses"] = pb
			res["resonant_strong_uses"] = pb
		}
		if level >= 5 {
			res["efficient_theorem_uses"] = pb
		}
		if level >= 7 {
			res["refuting_gesture_uses"] = 3 // baseline INT mod
		}
		if level >= 9 {
			res["simultaneous_lemma_uses"] = 1
		}
		if level >= 10 {
			res["rapid_axiom_uses"] = 1
		}
		if level >= 11 {
			res["distributed_proof_uses"] = pb
		}
		if level >= 14 {
			res["perfected_interleave_uses"] = 1
		}
		if level >= 17 {
			res["reductio_uses"] = 1
		}
		if level >= 18 {
			res["penultimate_proof_uses"] = 1
		}
		if level >= 20 {
			res["qed_uses"] = 1
		}

		z := 0
		m[level] = ClassLevelSeed{
			MaxSpellPoints: &z,
			MaxSpellLevel:  0,
			Resources:      res,
		}
	}
	return m
}

func syllogistArchetypes() []ArchetypeSeed {
	return []ArchetypeSeed{
		{
			Name:        "Chronologist",
			Description: "You specialize in tempo—freeing movement and compressing the round.",
			Features: map[int][]FeatureSeed{
				3: {
					{
						Name:        "Grant Hustle",
						Description: "As a bonus action, you can spend 1 Protocol Charge and choose one willing ally you can see within 30 feet. Until the start of your next turn, that ally can take the Dash or Disengage action once on their turn without using an action. You can use Grant Hustle a number of times equal to your proficiency bonus, and you regain all expended uses when you finish a short rest.",
						ActionType:     "Bonus Action",
						UsesPerRest:    "Proficiency Bonus",
						ResetCondition: "Short Rest",
						ResourceCosts:  []ResourceCostSeed{{Key: "protocol_charges", Amount: 1}, {Key: "grant_hustle_uses", Amount: 1}},
					},
				},
				7: {
					{
						Name:        "Borrowed Seconds",
						Description: "When you use Grant Hustle, the chosen ally also gains +10 feet to their walking speed until the end of their next turn.",
					},
				},
			},
		},
		{
			Name:        "Amplifier",
			Description: "You resonate Magnitudo others commit—especially Strong—without stealing their spotlight.",
			Features: map[int][]FeatureSeed{
				3: {
					{
						Name:        "Resonant Strong",
						Description: "When an ally you can see casts a spell that includes the Strong component, and this turn you have either spent Protocol Charges on a Syllogist feature or given that ally a benefit from Lend Axiom, you can treat that Strong as one magnitude step greater for that casting when the table negotiates area, duration, or damage escalation (the DM adjudicates what one step means in Faradhaven’s spell tables). You can do so a number of times equal to your proficiency bonus, and you regain all expended uses when you finish a long rest.",
						UsesPerRest:    "Proficiency Bonus",
						ResetCondition: "Long Rest",
						ResourceCosts:  []ResourceCostSeed{{Key: "resonant_strong_uses", Amount: 1}},
					},
				},
				7: {
					{
						Name:          "Harmonic Surge",
						Description:   "When an ally you can see finishes resolving a spell that included Strong and that benefited from Resonant Strong on the same turn, you can spend 1 Protocol Charge (no action) to let them reroll one damage die of that spell; they must use the new roll.",
						ResourceCosts: []ResourceCostSeed{{Key: "protocol_charges", Amount: 1}},
					},
				},
			},
		},
		{
			Name:        "Sympath",
			Description: "You bind multiple casters into one proof—sharing burden and sequence.",
			Features: map[int][]FeatureSeed{
				3: {
					{
						Name:        "Widened Lend Axiom",
						Description: "When you take the Lend Axiom action, you can spend 2 Protocol Charges instead of 1. If you do, you choose two willing allies within 30 feet instead of one. Each ally gains Lend Axiom’s benefit for the next eligible spell they cast before the start of your next turn (you still name one borrowed component per ally when you use the action).",
					},
				},
				7: {
					{
						Name:        "Linked Ledger",
						Description: "When two willing allies you can see are within 15 feet of each other, once per round when one of them casts a spell, you can spend 1 Protocol Charge (no action) to let the other donate one spell component they carry to that casting. The donor’s inventory is reduced as normal; the caster must still satisfy forge rules and their own class pool.",
						ResourceCosts: []ResourceCostSeed{{Key: "protocol_charges", Amount: 1}},
					},
				},
			},
		},
	}
}

func syllogistLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{
				Name: "No Independent Evocation",
				Description: "You cannot cast spells from this class: you have 0 spell points, 0 maximum spell level from this class, and your class component pool contains no Forma and no Essentia, so you cannot pass spell synthesis validation using only your class list.\n\n" +
					"You still track Protocol Charges and use your features to support allies. Any saving throw DC for a Syllogist feature that calls for a DC equals 8 + your proficiency bonus + your Intelligence modifier.",
			},
			{
				Name:        "Lend Axiom",
				Description: "As an action, you can spend 1 Protocol Charge and choose one willing ally within 30 feet who can see or hear you. Until the start of your next turn, the first time that ally casts a spell (if they do), they may treat one component named when you used this feature—If, Then, Therefore, or a Magnitudo you announce—as if it were on their class-and-race spell pool for that casting without expending that component from inventory.",
				ActionType:     "Action",
				ResourceCosts:  []ResourceCostSeed{{Key: "protocol_charges", Amount: 1}},
				UsesPerRest:    "",
				ResetCondition: "",
			},
		},
		2: {
			{
				Name:        "Interleave",
				Description: "When an ally you can see begins casting a spell that already includes at least one Logica component (If, Then, or Therefore), you can use your reaction and spend 1 Protocol Charge to assert how a Logica beat resolves relative to other beats (for example: a Then phase now vs after another phase). The DM adjudicates the result using your group’s spell rules.",
				ActionType:     "Reaction",
				ResourceCosts:  []ResourceCostSeed{{Key: "protocol_charges", Amount: 1}},
				UsesPerRest:    "",
				ResetCondition: "",
			},
		},
		3: {
			{
				Name:        "Syllogist Tradition",
				Description: "Choose Chronologist, Amplifier, or Sympath. You gain that tradition’s 3rd-level feature and its 7th-level feature when you reach that level.",
			},
		},
		4: {
			{
				Name:        "Ability Score Improvement",
				Description: "Increase one ability score by 2, or increase two ability scores by 1. As usual, you cannot raise a score above 20 with this feature unless another rule allows.",
			},
		},
		5: {
			{
				Name: "Efficient Theorem",
				Description: "Offensive support. Once on each of your turns when you spend Protocol Charges on a Syllogist feature, choose one ally you can see. The first time that ally deals damage with a spell or restores hit points with a spell before the end of their next turn, they roll one additional damage die or healing die for one of that spell’s rolls (their choice). Each spell can benefit only once. " +
					"You can grant this extra die a number of times equal to your proficiency bonus, and you regain all expended uses when you finish a long rest.",
				UsesPerRest:    "Proficiency Bonus",
				ResetCondition: "Long Rest",
				ResourceCosts:  []ResourceCostSeed{{Key: "efficient_theorem_uses", Amount: 1}},
			},
		},
		6: {
			{
				Name: "Chain Clause",
				Description: "When an ally casts a spell that includes Then or Therefore, you can spend 2 Protocol Charges (no action required unless the DM says timing matters) to let that spell’s narration include one extra resolution beat after the table agrees. That beat cannot add a new Essentia ingredient and cannot add an extra damage die unless another rule grants one; it may adjust Logica, Forma emphasis, or targeting only.",
				ResourceCosts:  []ResourceCostSeed{{Key: "protocol_charges", Amount: 2}},
				UsesPerRest:    "",
				ResetCondition: "",
			},
		},
		7: {
			{
				Name: "Refuting Gesture",
				Description: "Defense. When an ally within 30 feet who this round benefited from Lend Axiom or Interleave must make a saving throw against a spell, you can use your reaction to give them advantage on that save. You can use Refuting Gesture a number of times equal to your Intelligence modifier (minimum once), and you regain all expended uses when you finish a long rest.",
				ActionType:     "Reaction",
				UsesPerRest:    "Intelligence Modifier (min 1)",
				ResetCondition: "Long Rest",
				ResourceCosts:  []ResourceCostSeed{{Key: "refuting_gesture_uses", Amount: 1}},
			},
		},
		8: {
			{
				Name:        "Ability Score Improvement",
				Description: "Increase one ability score by 2, or increase two ability scores by 1. As usual, you cannot raise a score above 20 with this feature unless another rule allows.",
			},
		},
		9: {
			{
				Name: "Simultaneous Lemma",
				Description: "Utility / tempo spike. Once per long rest, you can use an action to choose two willing allies you can see within 60 feet. Until the start of your next turn, one of those allies may cast one spell they are allowed to cast with a casting time of 1 action by using a bonus action instead, provided that spell’s negotiated tier is no higher than 2nd for this rule. The second ally gains no automatic benefit from Simultaneous Lemma unless another feature targets them.",
				ActionType:     "Action",
				UsesPerRest:    "1",
				ResetCondition: "Long Rest",
				ResourceCosts:  []ResourceCostSeed{{Key: "simultaneous_lemma_uses", Amount: 1}},
			},
		},
		10: {
			{
				Name:        "Rapid Axiom",
				Description: "Tier upgrade (heroic tier). Once per short rest, you can use Lend Axiom by taking a bonus action instead of an action. You still spend Protocol Charges and follow all other requirements of Lend Axiom.",
				ActionType:     "Bonus Action",
				UsesPerRest:    "1",
				ResetCondition: "Short Rest",
				ResourceCosts:  []ResourceCostSeed{{Key: "rapid_axiom_uses", Amount: 1}},
			},
		},
		11: {
			{
				Name: "Distributed Proof",
				Description: "Utility. A number of times per long rest equal to your proficiency bonus, you can use an action to link yourself and one willing ally within 30 feet for 1 minute. While linked, when either of you casts a spell, non-pool spell components required by that spell may be deducted from either character’s inventory, provided that component is actually carried by one of you and you choose which inventory loses the piece at cast time.",
				ActionType:     "Action",
				UsesPerRest:    "Proficiency Bonus",
				ResetCondition: "Long Rest",
				ResourceCosts:  []ResourceCostSeed{{Key: "distributed_proof_uses", Amount: 1}},
			},
		},
		12: {
			{
				Name:        "Ability Score Improvement",
				Description: "Increase one ability score by 2, or increase two ability scores by 1. As usual, you cannot raise a score above 20 with this feature unless another rule allows.",
			},
		},
		13: {
			{
				Name: "Breather Lemma",
				Description: "Once per short rest, when you use Lend Axiom or Interleave, you can regain 1 Protocol Charge (you cannot exceed your maximum). Your Protocol Charges maximum increases at this level (see your class resources on the character sheet).",
			},
		},
		14: {
			{
				Name:        "Perfected Interleave",
				Description: "Once per short rest, when you use Interleave, you do not spend the Protocol Charge that Interleave would normally require.",
				UsesPerRest:    "1",
				ResetCondition: "Short Rest",
				ResourceCosts:  []ResourceCostSeed{{Key: "perfected_interleave_uses", Amount: 1}},
			},
		},
		15: {
			{
				Name: "Theorem Reserve",
				Description: "When you finish a short rest, you regain one expended use of Efficient Theorem’s extra die, provided you have fewer remaining uses than your maximum for that feature (your maximum is your proficiency bonus).",
			},
		},
		16: {
			{
				Name:        "Ability Score Improvement",
				Description: "Increase one ability score by 2, or increase two ability scores by 1. As usual, you cannot raise a score above 20 with this feature unless another rule allows.",
			},
		},
		17: {
			{
				Name: "Reductio ad Absurdum",
				Description: "Once per long rest, after an ally you can see finishes resolving a spell that included at least two Logica components, you choose one:\n\n" +
					"Disruption: You impose a harmless but tactically meaningful rider on one creature in the spell’s reach (no extra damage dice).\n\n" +
					"Impulse: Until the start of your next turn, that ally gains either one extra weapon attack when they take the Attack action, or they can take the Dash or Disengage action once without using an action—not both.",
				UsesPerRest:    "1",
				ResetCondition: "Long Rest",
				ResourceCosts:  []ResourceCostSeed{{Key: "reductio_uses", Amount: 1}},
			},
		},
		18: {
			{
				Name:        "Penultimate Proof",
				Description: "Your Protocol Charges maximum increases by 1 (already reflected in your level progression). In addition, once per long rest, when you roll initiative and you have no Protocol Charges remaining, you regain 2 Protocol Charges.",
				UsesPerRest:    "1",
				ResetCondition: "Long Rest",
				ResourceCosts:  []ResourceCostSeed{{Key: "penultimate_proof_uses", Amount: 1}},
			},
		},
		19: {
			{
				Name:        "Ability Score Improvement",
				Description: "Increase one ability score by 2, or increase two ability scores by 1. As usual, you cannot raise a score above 20 with this feature unless another rule allows.",
			},
		},
		20: {
			{
				Name: "Q.E.D.",
				Description: "Your Protocol Charges maximum increases by 2 (already reflected in your level progression). Once per long rest, after you finish a long rest, you regain Protocol Charges equal to half the Protocol Charges you spent only on features that targeted or directly benefited allies since your previous long rest (round down).",
				UsesPerRest:    "1",
				ResetCondition: "Long Rest",
				ResourceCosts:  []ResourceCostSeed{{Key: "qed_uses", Amount: 1}},
			},
		},
	}
}
