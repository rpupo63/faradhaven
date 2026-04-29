package faradhaven_classes

import (
	"github.com/rpupo63/faradhaven/backend/seed/seedmedia"
)

// Elixirist returns the Elixirist class seed.
// A master of Emax alchemy who brews the liquid Vital Spark into dual-effect potions.
// Every brew carries a life-giving Elixir Effect and a destabilizing Flux Effect.
func Elixirist() FaradhavenClassSeed {
	archetypeLevel := 3

	return FaradhavenClassSeed{
		Name:                 "The Elixirist",
		Description:          "A master of Emax alchemy who bottles the liquid Vital Spark into potions with both healing and toxic effects. Every brew you craft carries a life-giving Elixir and a destabilizing Flux — you decide whether the price is worth paying.",
		HitDie:               8,
		PrimaryAbility:       "Intelligence",
		PhotoURL:             seedmedia.URL("elixirist.jpg"),
		Archetype:            "Alchemical Support / Dual-Effect Brewer",
		Concept:              "A trained Homan Mixologist, outcast human chemist, or Galvanized being who has learned to externalize the Vital Spark in liquid form. They navigate the volatile balance between the Elixir (benefit) and the Flux (cost) in every potion they craft.",
		DnDSkillFocus:        []string{"Intelligence", "Medicine"},
		Proficiencies:        "Light Armor, Simple Weapons",
		SkillChoice:          []string{"Arcana", "Medicine", "Nature", "Investigation", "Insight", "Perception"},
		Tools:                []string{"Alchemist's Supplies", "Poisoner's Kit"},
		SavingThrows:         []string{"Constitution", "Intelligence"},
		AutomaticEquipNames:  []string{"4 empty vials"},
		AutomaticWeaponNames: []string{},
		AutomaticItemNames:   []string{"Leather Armor", "Alchemist's Supplies"},
		SpellCastingComponent:   "material",
		SpellCastingDescription: "You enact formulae by measuring, mixing, and catalyzing reagents—vials, powders, and kit spread before you. Others see hands busy with apparatus and labeled components; the magic looks like chemistry brought to impossibly precise life.",
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your primary weapon",
				Options: []EquipmentOptionSeed{
					{Description: "A Dagger", WeaponNames: []string{"Dagger"}},
					{Description: "A Light Crossbow and 20 bolts", WeaponNames: []string{"Light Crossbow"}},
				},
			},
			{
				Instruction: "Choose your starting kit",
				Options: []EquipmentOptionSeed{
					{Description: "Healer's Kit and 2 Vital Tonics", ItemNames: []string{"Healer's Kit", "Vital Tonic", "Vital Tonic"}},
					{Description: "Poisoner's Kit and Explorer's Pack", ItemNames: []string{"Poisoner's Kit", "Explorer's Pack"}},
				},
			},
		},
		LevelFeatures:    elixiristLevelFeatures(),
		LevelProgression: elixiristLevelProgression(),
		ArchetypeLevel:   &archetypeLevel,
		ResourceDefinitions: []ResourceDefinitionSeed{
			{Key: "emax_reserves", DisplayName: "Emax Reserves", Category: "pool", Description: "Your pool of liquid Vital Spark available for brewing. Fully restored on Long Rest. During a Short Rest, use Vital Extraction to recover Emax via skill check.", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 1},
			{Key: "prepared_formulas", DisplayName: "Prepared Formulas", Category: "limit", Description: "How many potion formulas you have memorized and ready to synthesize on demand. Swap any or all during a Short or Long Rest.", DisplayOrder: 2},
			{Key: "flux_dampener_uses", DisplayName: "Flux Dampener", Category: "pool", Description: "Uses per Short Rest to fully suppress a potion's Flux Effect without halving the Elixir.", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 3},
			{Key: "overdrive_uses", DisplayName: "Overdrive", Category: "pool", Description: "Uses per Long Rest to double all dice on both the Elixir and Flux Effects of one potion.", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 4},
		},
		ComponentPool: []string{
			// Forma (Shape) — potions can be self-applied, touch-delivered, thrown, or splash
			"Self", "Aura", "Projectile", "Zone",
			// Scopus (Targeting)
			"Target",
			// Essentia — Vita (life), Ignis (distillation heat), Umbra (flux darkness), Odor (vapor/inhalation)
			"Vita", "Ignis", "Umbra", "Odor",
			// Essentia (Pathos) — emotional residue in the brew
			"Fear", "Sadness",
			// Actio (Verbs)
			"Bind", "Pierce", "Destroy", "Increase", "Decrease", "Pull",
			// Magnitudo (Modifiers)
			"Strong", "Weak", "Extreme",
			// Logica (Potion chains and conditionals)
			"If", "Then",
		},
		Archetypes: []ArchetypeSeed{
			{
				Name:        "The Purist",
				Description: "You refine the Flux out of your brews, becoming the finest field medic in the Underground. Your potions heal more and hurt less, and you learn to use the Vital Spark to mend wounds that others have given up on.",
				Features: map[int][]FeatureSeed{
					3: {{Name: "Refined Formula", Description: "Your process filters the harshest volatility from your brews. Flux Effect dice on your potions drop one size (d6→d4, d4→minimum d4). Additionally, your Vital Tonic Elixir heals for an additional 1d6 hit points."}},
					7: {
						{Name: "Curative Mastery", Description: "You have mastered the art of restorative brewing. By spending 3 Emax during Vitalic Alchemy, you can brew a Cure Draught. When administered, a creature within 5 feet regains 2d8 + your Intelligence modifier hit points. Flux: the creature takes 1d4 psychic damage as the Vital Spark briefly overloads their system."},
						{Name: "Conditioning Agent", Description: "When you use Flux Dampener to suppress a potion's Flux Effect, you may also remove one condition (Poisoned, Diseased, or Stunned) from the target creature."},
					},
					11: {{Name: "Stabilized Brew", Description: "Once per Long Rest, you may perform a 1-minute ritual to implant one of your prepared potions into a willing creature's bloodstream. When that creature drops to 0 hit points, the implanted potion activates automatically, applying only its Elixir Effect (no Flux). The creature is then stabilized. The implanted potion is consumed whether or not the creature dropped to 0 HP before your next Long Rest."}},
					15: {{Name: "Living Serum", Description: "Your Stabilized Brew ritual now requires only 1 action (no longer requires 1 minute), and you may maintain implanted potions in up to your Intelligence modifier willing creatures simultaneously. Additionally, using Flux Dampener on a Cure Draught no longer expends a use of Flux Dampener."}},
				},
			},
			{
				Name:        "The Flux Merchant",
				Description: "You have learned to weaponize the Flux — the destabilizing side of every brew. You split deliveries, redirect toxic overflow, and turn the price of every potion into a punishment for your enemies.",
				Features: map[int][]FeatureSeed{
					3: {{Name: "Toxic Splash", Description: "When you throw a potion as a ranged attack, you may choose which effect to split: trigger only the Flux Effect on the target (as normal), or trigger only the Elixir Effect on yourself or a willing ally within 5 feet of you (the bottle detonates Flux harmlessly). You no longer need the same creature to receive both effects when throwing potions."}},
					7: {{Name: "Reactive Flux", Description: "When a Flux Effect from one of your potions deals damage to any creature within 30 feet of you, you may immediately deal bonus necrotic damage equal to your Intelligence modifier to a creature of your choice within 5 feet of the Flux target. This requires no action."}},
					11: {
						{Name: "Overdose Protocol", Description: "Your Overdrive uses per Long Rest increase by 2. Additionally, when you Overdrive a thrown potion (triggering only the Flux Effect on an unwilling target), the target must succeed on a Constitution saving throw (DC 8 + your proficiency bonus + your Intelligence modifier) or be Poisoned until the end of their next turn."},
					},
					15: {{Name: "Toxic Cascade", Description: "When a creature fails a saving throw against one of your Flux Effects, you may spend 2 Emax as a reaction to jump the full Flux Effect to a second creature of your choice within 10 feet of the original target. The second creature takes the same Flux damage without a saving throw."}},
				},
			},
			{
				Name:        "The Equilibrist",
				Description: "You seek the perfect resonance between Elixir and Flux. When both effects land on the same creature, the Vital Spark achieves a harmony that neither pure healing nor pure toxin can match — generating unique effects beyond either track alone.",
				Features: map[int][]FeatureSeed{
					3:  {{Name: "Perfect Equilibrium", Description: "When a creature willingly accepts both the Elixir Effect and the Flux Effect of one of your potions (they do not use Flux Dampener and do not halve the Elixir by refusing the Flux), the Vital Spark achieves resonance. That creature gains a bonus equal to your Intelligence modifier on their next attack roll, ability check, or saving throw before the end of their next turn."}},
					7:  {{Name: "Synergetic Reaction", Description: "When you brew a potion using Vitalic Alchemy, you may merge its Elixir and Flux into a single Balanced Brew. When administered, the target simultaneously regains 1d8 + your Intelligence modifier hit points AND takes 1d8 + your Intelligence modifier acid damage. Both effects apply in the same instant — one action, one target, no choice to split."}},
					11: {{Name: "Emax Flow", Description: "Your Emax Reserves maximum increases by 4. Additionally, whenever you use Flux Dampener to suppress a Flux Effect, you regain 1 Emax (up to your reserve cap)."}},
					15: {{Name: "Grand Synthesis", Description: "Once per Long Rest, as an action, you shatter a specially prepared Grand Equilibrium Brew. Choose any number of creatures within 30 feet of you. Each chosen creature regains 3d8 + your Intelligence modifier hit points and suffers no Flux damage. Every other creature within 30 feet takes 2d8 + your Intelligence modifier acid damage (no saving throw)."}},
				},
			},
		},
	}
}

func elixiristLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "Emax Reserves", Description: "You maintain a reserve of Emax — the liquid form of the Vital Spark that keeps the Galvanized alive. Your reserve equals the `emax_reserves` value at your current level. You regain all expended Emax after a Long Rest. See Flux Siphon and Vital Extraction for how you reclaim it mid-fight and during rests."},
			{Name: "Potion Formulary", Description: "You have internalized a set of potion recipes so thoroughly you can synthesize them on demand — no pre-brewing required. You have a number of Prepared Formulas equal to your `prepared_formulas` value. During a Short Rest or Long Rest, you may swap any or all of your Prepared Formulas, choosing from any formulas you know. To use a Prepared Formula, spend 1 Emax as an action: the potion materializes in your hand and is immediately administered to a willing or unconscious creature within 5 feet, or thrown as a ranged attack (range 20/60 ft). A prepared formula can be used any number of times as long as you have Emax to spend — there is no per-use limit beyond the cost.\n\nTier I Formulas (known from Level 1, choose 4 to know initially):\n• Vital Tonic — Elixir: target regains 1d6 + INT mod HP. Flux: 1d4 necrotic. Thrown: Flux only, 1d6 acid.\n• Shock Vial — Elixir: target has advantage on their next attack roll before end of their next turn. Flux: 1d6 lightning. Thrown: Flux only, 1d6 lightning + target loses reactions until start of their next turn.\n• Acid Flask — Elixir: target's next weapon attack ignores resistance to acid damage. Flux: 1d6 acid. Thrown: Flux only, 1d8 acid + target's AC is reduced by 1 until end of your next turn.\n• Numbing Agent — Elixir: target suppresses one condition (Frightened, Poisoned, or one level of Exhaustion) until start of your next turn. Flux: 1d4 necrotic + target has disadvantage on their next saving throw.\n• Smoke Tonic — Elixir: 10ft radius heavily obscured cloud centered on the target lasts 1 minute (or until a strong wind disperses it). Flux: target and any creature that starts its turn in the cloud must make a CON save (DC 12) or be Blinded until end of their next turn."},
			{Name: "The Dual Brew", Description: "All formulas you use carry two effects: an Elixir Effect (beneficial) and a Flux Effect (detrimental). The Flux triggers immediately after the Elixir resolves. When administering to a willing or unconscious creature within 5 feet, that creature may refuse the Flux — but doing so also halves the Elixir. When thrown at an unwilling creature, only the Flux Effect triggers on a hit."},
			{Name: "Flux Siphon", Description: "You have learned to recapture the volatile energy released when a brew destabilizes. Whenever a Flux Effect from one of your formulas deals damage to any creature — ally or enemy — you automatically regain 1 Emax, up to your reserve cap. This triggers regardless of the amount of damage dealt. If the Flux Effect is fully suppressed (via Flux Dampener or any other means), the Siphon does not trigger. This is the core loop: letting allies take the Flux costs them some damage but keeps your Emax flowing for free. Suppressing Flux protects allies but drains your reserves."},
		},
		2: {
			{Name: "Vital Extraction", Description: "During a Short Rest, or within 10 minutes of an encounter ending, you may spend 1 minute working with your Alchemist's Supplies to extract residual Emax from the environment and fallen creatures. Make an Intelligence (Alchemist's Supplies) check. The DC equals 8 + half the CR of the most dangerous enemy defeated in the last hour (rounded down, minimum DC 8). On a success, regain 1d4 + your Intelligence modifier Emax (up to your reserve cap). On a failure, regain 1 Emax. Usable once per Short or Long Rest."},
			{Name: "Flux Dampener", Description: "Before a formula's Flux Effect triggers, you may use a bonus action to spend 1 use of Flux Dampener, fully suppressing the Flux for that use. The Elixir Effect is unmodified. Suppressing Flux also prevents your Flux Siphon from recovering Emax — you protect your ally but forfeit the refuel. Flux Dampener uses restore on a Short Rest.", ActionType: "Bonus Action", ResourceCosts: []ResourceCostSeed{{Key: "flux_dampener_uses", Amount: 1}}},
		},
		3: {
			{Name: "Elixirist's Path", Description: "You have developed a signature approach to the dual brew. Choose a path: The Purist, The Flux Merchant, or The Equilibrist. Each path grants features at levels 3, 7, 11, and 15."},
		},
		5: {
			{Name: "Tier II Formulas", Description: "Your formulary expands. You now know all Tier II formulas and can prepare one additional formula (see `prepared_formulas`). Additionally, all Tier I formula dice increase by one size (d4→d6, d6→d8). Tier II Formulas:\n• Accelerant Dose — Elixir: target gains +10ft movement and may Disengage as a bonus action on their next turn. Flux: 1d6 fire.\n• Vitalic Surge — Elixir: target regains 2d6 + INT mod HP. Flux: 1d6 necrotic + target has disadvantage on their first attack roll before end of their next turn.\n• Muscle Serum — Elixir: target's next weapon attack deals +1d8 extra damage. Flux: 1d6 necrotic + target has disadvantage on DEX saves until end of their next turn.\n• Ferrofluid — Elixir: target gains resistance to one damage type you choose when you prepare this formula until start of your next turn. Flux: 1d8 necrotic."},
			{Name: "Overdrive", Description: "When you use a Prepared Formula, you may spend a bonus action and 1 Overdrive use to double all dice on both its Elixir and Flux Effects. Overdrive uses restore on a Long Rest.", ActionType: "Bonus Action", ResourceCosts: []ResourceCostSeed{{Key: "overdrive_uses", Amount: 1}}},
		},
		6: {
			{Name: "Volatile Formula", Description: "When you use a component-crafted ability that requires an action, you may simultaneously use one of your Prepared Formulas on the same target (or yourself) as part of the same action. Both the component ability's effect and the formula's full Elixir and Flux apply. No extra action cost."},
		},
		9: {
			{Name: "Tier III Formulas", Description: "Your formulary reaches its advanced tier. You now know all Tier III formulas and can prepare one additional formula (see `prepared_formulas`). Tier III Formulas:\n• Paralytic Extract — Elixir (willing): target may voluntarily Restrain themselves until start of your next turn, gaining advantage on Strength (Athletics) checks and saves against forced movement during that time. Flux: 1d8 necrotic. Thrown (unwilling): target must succeed on a CON save (DC 8 + Prof + INT) or be Restrained until start of your next turn. Flux: 1d8 necrotic on hit regardless.\n• Coagulant — Elixir: target gains temporary HP equal to 2d8 + INT mod. Flux: 1d6 necrotic, dealt when those temporary HP are depleted.\n• Echo Serum — Elixir: until start of your next turn, target reduces all incoming damage by 1d4 (roll once on administration). Flux: 1d8 psychic.\n• Cascade Serum — Elixir: until end of target's next turn, the first time they would be healed by any source, double the amount healed. Flux: 2d6 necrotic."},
		},
		10: {
			{Name: "Practiced Hand", Description: "Once per Short Rest, you may use a Prepared Formula as a bonus action instead of an action.", ActionType: "Bonus Action"},
		},
		13: {
			{Name: "Master Reagents", Description: "Your formula work has become instinctive. You may use any Prepared Formula as a bonus action (instead of only once per Short Rest via Practiced Hand — that feature is now unlimited). Additionally, your Flux Siphon now triggers even on a partial resist: if a creature succeeds on a saving throw against a Flux Effect and takes half damage, you still regain 1 Emax. Finally, your Vital Extraction check is made with advantage.", ActionType: "Bonus Action"},
		},
		14: {
			{Name: "Cascade Reaction", Description: "When a creature takes Flux damage from one of your formulas, you may spend 2 Emax as a reaction to redirect half that Flux damage (rounded down) to a second creature of your choice within 10 feet of the first target. The second creature takes the damage without a saving throw.", ActionType: "Reaction", ResourceCosts: []ResourceCostSeed{{Key: "emax_reserves", Amount: 2}}},
		},
		17: {
			{Name: "Apex Distillation", Description: "All Tier II and Tier III formula dice increase by one size. Once per Long Rest, when you use a Prepared Formula, you may declare it an Apex Brew: apply the effects of two of your other Prepared Formulas simultaneously alongside it (one action, three formula effects, one Emax cost). All three Flux Effects trigger."},
		},
		18: {
			{Name: "Emax Mastery", Description: "Your Flux Siphon now recovers 2 Emax instead of 1. Additionally, once per Short Rest, you may use a Prepared Formula without spending any Emax."},
		},
		20: {
			{Name: "Perfect Formula", Description: "Your Intelligence score increases by 4, to a maximum of 24. Your `prepared_formulas` cap is ignored — you may have all formulas you know prepared simultaneously. Once per Long Rest, when you use a formula, you may declare it a Perfect Elixir: its Elixir Effect is doubled and no Flux Effect triggers."},
		},
	}
}

func elixiristLevelProgression() map[int]ClassLevelSeed {
	// prepared_formulas: starts at 2 (very tight), grows as tiers unlock.
	// emax_reserves: the Flux Siphon loop keeps this sustainable mid-combat;
	// the cap mainly limits big-cost features (Cascade Reaction, etc.).
	return map[int]ClassLevelSeed{
		1:  {Resources: map[string]int{"emax_reserves": 4, "prepared_formulas": 2, "flux_dampener_uses": 1, "overdrive_uses": 0}},
		2:  {Resources: map[string]int{"emax_reserves": 4, "prepared_formulas": 2, "flux_dampener_uses": 1, "overdrive_uses": 0}},
		3:  {Resources: map[string]int{"emax_reserves": 5, "prepared_formulas": 2, "flux_dampener_uses": 1, "overdrive_uses": 0}},
		4:  {Resources: map[string]int{"emax_reserves": 5, "prepared_formulas": 2, "flux_dampener_uses": 1, "overdrive_uses": 0}},
		5:  {Resources: map[string]int{"emax_reserves": 6, "prepared_formulas": 3, "flux_dampener_uses": 1, "overdrive_uses": 1}},
		6:  {Resources: map[string]int{"emax_reserves": 6, "prepared_formulas": 3, "flux_dampener_uses": 1, "overdrive_uses": 1}},
		7:  {Resources: map[string]int{"emax_reserves": 7, "prepared_formulas": 3, "flux_dampener_uses": 2, "overdrive_uses": 1}},
		8:  {Resources: map[string]int{"emax_reserves": 7, "prepared_formulas": 3, "flux_dampener_uses": 2, "overdrive_uses": 1}},
		9:  {Resources: map[string]int{"emax_reserves": 8, "prepared_formulas": 4, "flux_dampener_uses": 2, "overdrive_uses": 1}},
		10: {Resources: map[string]int{"emax_reserves": 8, "prepared_formulas": 4, "flux_dampener_uses": 2, "overdrive_uses": 2}},
		11: {Resources: map[string]int{"emax_reserves": 9, "prepared_formulas": 4, "flux_dampener_uses": 2, "overdrive_uses": 2}},
		12: {Resources: map[string]int{"emax_reserves": 9, "prepared_formulas": 4, "flux_dampener_uses": 2, "overdrive_uses": 2}},
		13: {Resources: map[string]int{"emax_reserves": 10, "prepared_formulas": 4, "flux_dampener_uses": 3, "overdrive_uses": 2}},
		14: {Resources: map[string]int{"emax_reserves": 10, "prepared_formulas": 4, "flux_dampener_uses": 3, "overdrive_uses": 2}},
		15: {Resources: map[string]int{"emax_reserves": 11, "prepared_formulas": 5, "flux_dampener_uses": 3, "overdrive_uses": 2}},
		16: {Resources: map[string]int{"emax_reserves": 11, "prepared_formulas": 5, "flux_dampener_uses": 3, "overdrive_uses": 2}},
		17: {Resources: map[string]int{"emax_reserves": 12, "prepared_formulas": 5, "flux_dampener_uses": 3, "overdrive_uses": 2}},
		18: {Resources: map[string]int{"emax_reserves": 12, "prepared_formulas": 5, "flux_dampener_uses": 3, "overdrive_uses": 2}},
		19: {Resources: map[string]int{"emax_reserves": 12, "prepared_formulas": 5, "flux_dampener_uses": 3, "overdrive_uses": 2}},
		20: {Resources: map[string]int{"emax_reserves": 12, "prepared_formulas": 5, "flux_dampener_uses": 3, "overdrive_uses": 2}},
	}
}
