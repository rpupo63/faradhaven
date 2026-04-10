package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/seedmedia"
)

// Mutagen returns the Mutagen class seed
func Mutagen() FaradhavenClassSeed {
	archetypeLevel := 3

	return FaradhavenClassSeed{
		Name:                 "The Mutagen",
		Description:          "A volatile warrior who channels component-crafted abilities through their own body, risking madness to unleash primal power. You balance on the edge of a feral state, growing stronger as your humanity slips away.",
		HitDie:               10,
		PrimaryAbility:       "Strength",
		PhotoURL:             seedmedia.URL("mutagen.jpg"),
		Archetype:            "Biological Horror / Berserker",
		Concept:              "A survivor of experimental alchemy or plague who metabolizes magic into physical change. They struggle to maintain control against the beast within.",
		DnDSkillFocus:        []string{"Strength", "Survival"},
		Proficiencies:        "Light & Medium Armor, Simple & Martial Weapons",
		SkillChoice:          []string{"Athletics", "Intimidation", "Survival", "Nature", "Medicine"},
		Tools:                []string{"Alchemist's Supplies"},
		SavingThrows:         []string{"Constitution", "Strength"},
		AutomaticEquipNames:  []string{"Alchemist's Supplies"},
		AutomaticWeaponNames: []string{"Mutagen Bite"},
		AutomaticItemNames:   []string{"Scale Mail", "Explorer's Pack"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your heavy weapon",
				Options: []EquipmentOptionSeed{
					{Description: "A Greataxe", WeaponNames: []string{"Greataxe"}},
					{Description: "A Greatsword", WeaponNames: []string{"Greatsword"}},
				},
			},
			{
				Instruction: "Choose your survival gear",
				Options: []EquipmentOptionSeed{
					{Description: "Stabilizer Serum and Dagger", ItemNames: []string{"Stabilizer Serum"}, WeaponNames: []string{"Dagger"}},
					{Description: "Explorer's Pack", ItemNames: []string{"Explorer's Pack"}},
				},
			},
		},
		LevelFeatures:    mutagenLevelFeatures(),
		LevelProgression: mutagenLevelProgression(),
		ArchetypeLevel:   &archetypeLevel,
		ResourceDefinitions: []ResourceDefinitionSeed{
			{Key: "madness_base_dc", DisplayName: "Madness Base DC", Category: "state", Description: "Starting DC for Madness saving throws (resets on rest)", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 1},
			{Key: "feral_bonus", DisplayName: "Feral Bonus", Category: "modifier", Description: "Extra melee damage while in Feral Mode", DisplayOrder: 2},
			{Key: "madness_resistance_uses", DisplayName: "Madness Resistance", Category: "pool", Description: "Uses per day to auto-succeed a failed Madness save", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 3},
			{Key: "rapid_repair_uses", DisplayName: "Rapid Repair", Category: "pool", Description: "Uses per Long Rest to regain HP as a bonus action (CON modifier).", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 4},
			{Key: "ethereal_step_uses", DisplayName: "Ethereal Step", Category: "pool", Description: "Uses per Long Rest to turn invisible (CON modifier).", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 5},
			{Key: "undying_mutation_uses", DisplayName: "Undying Mutation", Category: "pool", Description: "1/Long Rest: Drop to 1 HP instead of 0 when reduced to 0 HP (not killed outright).", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 6},
		},
		ComponentPool: []string{
			// Forma (Shape)
			"Self", "Aura",
			// Essentia (Matter & Energy)
			"Vita", "Mortis", "Ignis", "Terra", "Anger", "Fear",
			// Actio (Verbs)
			"Mutate", "Crush", "Pierce", "Push", "Grab", "Destroy",
			// Magnitudo (Modifiers)
			"Increase", "Decrease", "Strong", "Weak", "Extreme",
		},
		Archetypes: []ArchetypeSeed{
			{
				Name:        "The Behemoth",
				Description: "Focuses on raw physical power, size increases, and heavy armor plating. You become an unstoppable siege monster.",
				Features: map[int][]FeatureSeed{
					3:  {{Name: "Titan's Growth", Description: "You count as one size larger for carrying capacity. As part of your Attack action, you can take a 5-foot step before any attack without provoking opportunity attacks. Your Feral Bonus Damage increases by an additional +1."}},
					7:  {{Name: "Living Fortress", Description: "While in Feral Mode, you gain a +1 bonus to AC. Additionally, you gain temporary hit points equal to your Constitution modifier at the start of each of your turns."}},
					11: {{Name: "Siege Breaker", Description: "Your melee weapon attacks deal double damage to objects and structures. When you hit a creature, you can use a bonus action to push them 10 feet away if they fail a Strength save (DC 8 + your proficiency bonus + your Strength modifier)."}},
					15: {{Name: "Colossal Form", Description: "When you enter Feral Mode, you can choose to channel a colossal mutation. If you do, your body temporarily swells with power, and your melee weapon attacks deal an extra 1d6 damage for the duration of the Feral Mode."}},
				},
			},
			{
				Name:        "The Phase-Shifter",
				Description: "Focuses on unstable molecular bonds, allowing for blink-steps, invisibility, and ethereal mutations.",
				Features: map[int][]FeatureSeed{
					3:  {{Name: "Blink Strike", Description: "When you take the Attack action, you can teleport up to 10 feet to an unoccupied space you can see before each attack."}},
					7:  {{Name: "Ethereal Step", Description: "As a bonus action, you can turn invisible until the start of your next turn or until you attack.", ActionType: "Bonus Action", ResourceCosts: []ResourceCostSeed{{Key: "ethereal_step_uses", Amount: 1}}}},
					11: {{Name: "Phase Dodge", Description: "When hit by an attack, use your reaction to roll 1d6. On a 4+, the attack misses as you phase out of reality."}},
					15: {{Name: "Reality Tear", Description: "You shift between planes. At the end of each of your turns, roll a d20. On a 11 or higher, you vanish into the Ethereal Plane. While there, you can't be affected by anything on the material plane. At the start of your next turn, you reappear in an unoccupied space within 10 feet of where you vanished. Your first attack after reappearing deals an extra 2d6 force damage."}},
				},
			},
			{
				Name:        "The Plague-Bearer",
				Description: "Your mutations exude toxic clouds and life-steal mechanics. You are a walking biohazard.",
				Features: map[int][]FeatureSeed{
					3:  {{Name: "Toxic Aura", Description: "While in Feral Mode, creatures that start their turn within 5 feet of you take poison damage equal to your Constitution modifier (minimum 1)."}},
					7:  {{Name: "Parasitic Drain", Description: "Once per turn when you deal damage with a melee attack in Feral Mode, you regain hit points equal to your Constitution modifier."}},
					11: {{Name: "Contagion Carrier", Description: "You are immune to disease and poison damage. Your Toxic Aura range increases to 10 feet and ignores resistance to poison damage."}},
					15: {{Name: "Viral Explosion", Description: "When you are reduced to 0 hit points, you release a 20ft radius cloud of gas. Each creature must make a Constitution save (DC 8 + Prof + Con), taking 8d6 poison damage on a failure, or half on a success."}},
				},
			},
		},
	}
}

// MutagenFeralTable returns 20 unique effects for the Mutagen's feral state.
func MutagenFeralTable() map[int]string {
	return map[int]string{
		1:  "Red Haze: You enter a Berserk state. You must use your action to attack the nearest creature (friend or foe). Your melee attacks deal an extra 1d6 damage.",
		2:  "Predatory Instinct: You must attack the creature with the lowest current hit points. You have advantage on these attacks.",
		3:  "Gluttonous Consumption: If you reduce a creature to 0 HP, you must spend your next action devouring it. You regain 2d10 HP, but lose your turn.",
		4:  "Bone Shard Eruption: Sharp bone spurs burst from your skin. You take 1d8 piercing damage, but any creature that hits you with a melee attack takes 1d8 piercing damage.",
		5:  "Hyper-Adrenaline: Your speed is doubled, you gain a +2 bonus to AC, and you gain an additional action on each of your turns. This state lasts for 2 rounds. When it ends, you are lethargic and cannot move or take actions for 1 round.",
		6:  "Toxic Bile: You vomit a 15ft cone of acid. Each creature must make a Con save (DC 8 + Prof + Con) or take 3d6 acid damage. You take 1d6 acid damage.",
		7:  "Muscle Spasm: Your Strength increases massively (+4 to hit/damage), but your movement speed is reduced to 5 feet.",
		8:  "Blind Agony: You are blinded by pain. Your attacks are at disadvantage, but you are immune to being Charmed or Frightened.",
		9:  "Leaping Hunter: Your jump distance is tripled. You must use your movement to leap towards the furthest enemy you can see and attack them.",
		10: "Acidic Blood: When you take damage, you spray acidic blood. All creatures within 5ft (including allies) take 1d4 acid damage.",
		11: "Hardened Carapace: Your skin turns into a thick shell. Your AC increases by 3, but you have disadvantage on all Dexterity saving throws.",
		12: "Territorial Roar: You must use your action to let out a blood-curdling roar. All creatures within 30ft must make a Wis save or be Frightened of you for 1 minute.",
		13: "Vampiric Fangs: Your jaw unhinges. You must use your action to make a Bite attack using your Mutagen Bite weapon (+Prof + Str to hit, 1d10 + Str modifier piercing damage). You regain HP equal to the damage dealt.",
		14: "Unstable Growth: You grow one size category larger. You deal an extra 1d4 damage, but you take 1d4 force damage at the start of each of your turns as your cells rupture.",
		15: "Sensory Overload: You have advantage on Perception checks, but you are vulnerable to Thunder and Radiant damage.",
		16: "Numbing Cold: Your body temperature drops. You gain resistance to Fire damage, but your movement speed is reduced by 10 feet.",
		17: "Corrosive Touch: Any non-magical weapon you hold begins to melt (-1 damage per turn). Your unarmed strikes deal an extra 1d6 acid damage.",
		18: "Rejuvenating Slumber: You fall unconscious for 1 round as your body rapidly repairs itself. You regain 4d8 HP.",
		19: "Thickened Sinews: You gain a +5 bonus to Strength checks and saves, but your Intelligence and Charisma scores drop to 3 while feral.",
		20: "Apex Fury: You gain an extra action each turn, but you take 2d6 psychic damage at the end of each turn as your mind fractures.",
	}
}

func mutagenLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "The Madness Mechanic", Description: "Your mutations are crafted from components, using the Self or Touch shape exclusively. Each component-crafted ability of 1st level or higher requires a Madness save under these conditions: (1) the spell uses more components than your proficiency bonus, or (2) you have cast more spells today than your proficiency bonus + your Constitution modifier. When a Madness save is required, make a Constitution saving throw against your current Madness DC (starts at 10). On failure, you enter Feral Mode. The DC increases by +2 (+1 starting at level 6) for each cast that required a save, whether you passed or failed. Your Madness DC resets to 10 after a Short or Long Rest."},
			{Name: "Feral Mode", Description: "On a failed Madness save, you enter a Feral State for 1 minute. Roll on the Mutagen Feral Table to determine your specific mutation effect. While feral, you gain Resistance to Bludgeoning, Piercing, and Slashing damage, and your melee attacks deal extra damage equal to your Feral Bonus. You may still use component-crafted abilities while in Feral Mode."},
			{Name: "Internalized Arcana", Description: "Your body is the catalyst. You can only craft abilities using the Self or Touch shape components. All other shapes are inaccessible to you."},
			{Name: "Madness Resistance (1/Day)", Description: "Once per day, you can choose to succeed on a failed Madness saving throw. The Madness DC still increases as normal, since the ability was still cast."},
		},
		2:  {{Name: "Biological Adaptation", Description: "Choose one: Darkvision (60ft), Climbing Speed (equal to walking), or Natural Armor (+1 bonus to AC)."}},
		3:  {{Name: "Mutagenic Archetype", Description: "Choose a path: The Behemoth, The Phase-Shifter, or The Plague-Bearer. Each path grants unique mutations and physical changes."}},
		5:  {{Name: "Extra Attack, Madness Resistance (2/Day)", Description: "You can attack twice when you take the Attack action. You can now use Madness Resistance twice per day."}},
		6:  {{Name: "Stabilized Bloodstream", Description: "Your body adapts to the strain. The Madness DC now increases by only +1 for each component-crafted ability that required a save. Additionally, when you enter Feral Mode, you can roll two d20s for the Feral Table and choose which result applies."}},
		7:  {{Name: "Adaptive Immunity", Description: "Your body has developed resistance to foreign agents. You have advantage on saving throws against poison and disease, and you gain resistance to poison damage."}},
		9:  {{Name: "Advanced Synthesis", Description: "You gain access to 3rd-level component abilities. Your Feral Bonus Damage increases to +4."}, {Name: "Mutation Potency", Description: "When you cast a component-crafted ability, you can choose to increase the Madness DC by an additional +5 to cast it at one level higher than its base."}},
		10: {{Name: "Mind of the Hive", Description: "You are immune to the Frightened condition, and you have advantage on saving throws against being Charmed."}},
		11: {{Name: "Undying Mutation", Description: "When you are reduced to 0 hit points but not killed outright, you can drop to 1 hit point instead. Once you use this feature, you can't use it again until you finish a Long Rest."}},
		13: {{Name: "Rapid Cellular Repair", Description: "As a bonus action, you regain hit points equal to 1d10 + your Constitution modifier.", ActionType: "Bonus Action", ResourceCosts: []ResourceCostSeed{{Key: "rapid_repair_uses", Amount: 1}}}},
		14: {{Name: "Controlled Ferocity", Description: "In Feral Mode, you no longer attack allies. You can choose your targets freely among hostile creatures (melee only)."}},
		15: {{Name: "Evolved Metabolism", Description: "Your body has fully adapted to minor mutations. Component-crafted abilities of 2nd level or lower no longer increase your Madness DC."}},
		17: {{Name: "Apex Predator", Description: "You gain advantage on all Madness saving throws. Your Feral Bonus Damage increases to +6."}, {Name: "Vicious Reflexes", Description: "While in Feral Mode, you can make one additional melee weapon attack as a bonus action."}},
		18: {{Name: "Unstoppable Mutation", Description: "You stop aging, cannot be aged magically, and are immune to exhaustion caused by casting component-crafted abilities."}},
		20: {{Name: "Perfect Organism", Description: "1st-level component abilities no longer increase your Madness DC. Your Constitution score increases by 4, to a maximum of 24."}},
	}
}

func mutagenLevelProgression() map[int]ClassLevelSeed {
	// Mutagen gets Extra Attack at level 5
	// MadnessBaseDC is always 10, FeralBonus scales with level
	// undying_mutation_uses unlocks at level 11
	return map[int]ClassLevelSeed{
		1:  {Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 2, "madness_resistance_uses": 1, "rapid_repair_uses": 0, "ethereal_step_uses": 0, "undying_mutation_uses": 0}},
		2:  {Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 2, "madness_resistance_uses": 1, "rapid_repair_uses": 0, "ethereal_step_uses": 0, "undying_mutation_uses": 0}},
		3:  {Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 2, "madness_resistance_uses": 1, "rapid_repair_uses": 0, "ethereal_step_uses": 0, "undying_mutation_uses": 0}},
		4:  {Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 2, "madness_resistance_uses": 1, "rapid_repair_uses": 0, "ethereal_step_uses": 0, "undying_mutation_uses": 0}},
		5:  {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 3, "madness_resistance_uses": 2, "rapid_repair_uses": 0, "ethereal_step_uses": 0, "undying_mutation_uses": 0}},
		6:  {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 3, "madness_resistance_uses": 2, "rapid_repair_uses": 0, "ethereal_step_uses": 0, "undying_mutation_uses": 0}},
		7:  {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 3, "madness_resistance_uses": 2, "rapid_repair_uses": 0, "ethereal_step_uses": 3, "undying_mutation_uses": 0}},
		8:  {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 3, "madness_resistance_uses": 2, "rapid_repair_uses": 0, "ethereal_step_uses": 3, "undying_mutation_uses": 0}},
		9:  {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 4, "madness_resistance_uses": 2, "rapid_repair_uses": 0, "ethereal_step_uses": 3, "undying_mutation_uses": 0}},
		10: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 4, "madness_resistance_uses": 2, "rapid_repair_uses": 0, "ethereal_step_uses": 3, "undying_mutation_uses": 0}},
		11: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 4, "madness_resistance_uses": 2, "rapid_repair_uses": 0, "ethereal_step_uses": 3, "undying_mutation_uses": 1}},
		12: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 4, "madness_resistance_uses": 2, "rapid_repair_uses": 0, "ethereal_step_uses": 3, "undying_mutation_uses": 1}},
		13: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 5, "madness_resistance_uses": 2, "rapid_repair_uses": 3, "ethereal_step_uses": 3, "undying_mutation_uses": 1}},
		14: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 5, "madness_resistance_uses": 2, "rapid_repair_uses": 3, "ethereal_step_uses": 3, "undying_mutation_uses": 1}},
		15: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 5, "madness_resistance_uses": 2, "rapid_repair_uses": 3, "ethereal_step_uses": 3, "undying_mutation_uses": 1}},
		16: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 5, "madness_resistance_uses": 2, "rapid_repair_uses": 3, "ethereal_step_uses": 3, "undying_mutation_uses": 1}},
		17: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 6, "madness_resistance_uses": 2, "rapid_repair_uses": 3, "ethereal_step_uses": 3, "undying_mutation_uses": 1}},
		18: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 6, "madness_resistance_uses": 2, "rapid_repair_uses": 3, "ethereal_step_uses": 3, "undying_mutation_uses": 1}},
		19: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 6, "madness_resistance_uses": 2, "rapid_repair_uses": 3, "ethereal_step_uses": 3, "undying_mutation_uses": 1}},
		20: {ExtraAttackCount: 1, Resources: map[string]int{"madness_base_dc": 10, "feral_bonus": 7, "madness_resistance_uses": 2, "rapid_repair_uses": 3, "ethereal_step_uses": 3, "undying_mutation_uses": 1}},
	}
}
