package faradhaven_classes

// Mutagen returns the Mutagen class seed
func Mutagen() FaradhavenClassSeed {
	archetypeLevel := 3

	return FaradhavenClassSeed{
		Name:           "The Mutagen",
		Description:    "A volatile warrior who trades spell slots for mutation strains, risking madness to unleash primal power. You balance on the edge of a feral state, growing stronger as your humanity slips away.",
		HitDie:         10,
		PrimaryAbility: "Constitution",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/mutagen.jpg", // Keeping existing URL or placeholder
		Archetype:      "Biological Horror / Berserker",
		Concept:        "A survivor of experimental alchemy or plague who metabolizes magic into physical change. They struggle to maintain control against the beast within.",
		DnDSkillFocus:    []string{"Constitution", "Survival"},
		Proficiencies:    "Light & Medium Armor, Simple & Martial Weapons",
		SkillChoice:      []string{"Athletics", "Intimidation", "Survival", "Nature", "Medicine"},
		Tools:            []string{"Alchemist's Supplies"},
		SavingThrows:     []string{"Constitution", "Strength"},
		AutomaticEquipNames: []string{"Alchemist's Supplies"},
		AutomaticItemNames:  []string{"Scale Mail", "Explorer's Pack"},
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
		LevelFeatures:       mutagenLevelFeatures(),
		LevelProgression:    mutagenLevelProgression(),
		ArchetypeLevel:      &archetypeLevel,
		ResourceType:        "madness",
		ResourceName:        "Madness DC",
		ResourceRestoreType: "short_rest", // DC resets on short or long rest
		Archetypes: []ArchetypeSeed{
			{
				Name:        "The Behemoth",
				Description: "Focuses on raw physical power, size increases, and heavy armor plating. You become an unstoppable siege monster.",
				Features: map[int][]FeatureSeed{
					3:  {{Name: "Titan's Growth", Description: "You count as one size larger when determining carrying capacity and the weight you can push, drag, or lift. Your Feral Bonus Damage increases by +1."}},
					7:  {{Name: "Living Fortress", Description: "While in Feral Mode, you gain a +1 bonus to AC and temporary hit points equal to your Constitution modifier at the start of each of your turns."}},
					11: {{Name: "Siege Breaker", Description: "Your melee weapon attacks deal double damage to objects and structures. When you hit a creature with a melee weapon attack, you can push them up to 10 feet away."}},
					15: {{Name: "Colossal Form", Description: "When you enter Feral Mode, your size doubles in all dimensions, and your weight is multiplied by eight. You deal an extra 1d4 damage on all melee attacks."}},
				},
			},
			{
				Name:        "The Phase-Shifter",
				Description: "Focuses on unstable molecular bonds, allowing for blink-steps, invisibility, and ethereal mutations.",
				Features: map[int][]FeatureSeed{
					3:  {{Name: "Blink Strike", Description: "When you take the Attack action, you can teleport up to 10 feet before each attack to an unoccupied space you can see."}},
					7:  {{Name: "Ethereal Step", Description: "As a bonus action, you can turn invisible until the start of your next turn or until you attack or cast a spell."}},
					11: {{Name: "Phase Dodge", Description: "When you are hit by an attack, you can use your reaction to roll 1d6. On a 4 or higher, the attack misses you as you phase out of reality."}},
					15: {{Name: "Reality Tear", Description: "You can cast the 'Blink' spell on yourself without expending a spell slot or increasing your Madness DC. While blinking, your Feral Bonus Damage applies to your first attack after reappearing."}},
				},
			},
			{
				Name:        "The Plague-Bearer",
				Description: "Your mutations exude toxic clouds and life-steal mechanics. You are a walking biohazard.",
				Features: map[int][]FeatureSeed{
					3:  {{Name: "Toxic Aura", Description: "Creatures that start their turn within 5 feet of you take poison damage equal to your Constitution modifier."}},
					7:  {{Name: "Parasitic Drain", Description: "When you deal damage to a creature with a melee attack while in Feral Mode, you regain hit points equal to half the damage dealt."}},
					11: {{Name: "Contagion Carrier", Description: "You are immune to disease and poison damage. Your Toxic Aura range increases to 10 feet."}},
					15: {{Name: "Viral Explosion", Description: "When you are reduced to 0 hit points, your body releases a cloud of virulent gas. Each creature within 20 feet of you must make a Constitution saving throw (DC 8 + your Constitution modifier + your Proficiency Bonus), taking 8d6 poison damage on a failed save, or half as much on a successful one."}},
				},
			},
		},
	}
}

func mutagenLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "The Madness Mechanic", Description: "Instead of Spell Slots, you use Mutation Strains. Every time you cast a spell (Mutation), you must make a Madness Save (Constitution). The DC starts at 10 and increases by +2 with each subsequent cast. Reset: The DC resets to 10 after a Short or Long Rest."},
			{Name: "Feral Mode", Description: "On a failed Madness save, you enter a Feral State for 1 minute. You cannot cast spells, you lose the ability to distinguish friend from foe, and you must use your action to attack the nearest creature. You gain Resistance to Bludgeoning, Piercing, and Slashing damage and deal extra damage while Feral. Feral Bonus Damage: +2."},
			{Name: "Internalized Arcana", Description: "You can only cast spells with a range of 'Self' or 'Touch'."},
			{Name: "Madness Resistance (1/Day)", Description: "You can choose to succeed on a failed Madness save."},
		},
		2:  {{Name: "Biological Adaptation", Description: "Choose one: Darkvision (60ft), Climbing Speed (equal to walking), or Natural Armor (+1 AC). Feral Bonus Damage: +2."}},
		3:  {{Name: "Mutagenic Archetype", Description: "Choose a path: The Behemoth, The Phase-Shifter, or The Plague-Bearer. Feral Bonus Damage: +2."}},
		4:  {{Name: "Ability Score Improvement", Description: "Feral Bonus Damage: +2."}},
		5:  {{Name: "Extra Attack, Madness Resistance (2/Day)", Description: "Feral Bonus Damage: +3."}},
		6:  {{Name: "Stabilized Bloodstream", Description: "Madness DC increases by only +1 per spell cast. Feral Bonus Damage: +3."}},
		7:  {{Name: "Archetype Feature", Description: "Feral Bonus Damage: +3."}},
		8:  {{Name: "Ability Score Improvement", Description: "Feral Bonus Damage: +3."}},
		9:  {{Name: "Advanced Synthesis", Description: "You gain access to 3rd-level mutations (spells). Feral Bonus Damage: +4."}},
		10: {{Name: "Mind of the Hive", Description: "Immunity to Frightened condition. Feral Bonus Damage: +4."}},
		11: {{Name: "Archetype Feature", Description: "Feral Bonus Damage: +4."}},
		12: {{Name: "Ability Score Improvement", Description: "Feral Bonus Damage: +4."}},
		13: {{Name: "Rapid Cellular Repair", Description: "Regain hit points equal to 1d10 + Con modifier as a bonus action (Uses: Con modifier/Long Rest). Feral Bonus Damage: +5."}},
		14: {{Name: "Controlled Ferocity", Description: "In Feral Mode, you no longer attack allies and can choose targets (melee only). Feral Bonus Damage: +5."}},
		15: {{Name: "Archetype Feature", Description: "Feral Bonus Damage: +5."}},
		16: {{Name: "Ability Score Improvement", Description: "Feral Bonus Damage: +5."}},
		17: {{Name: "Apex Predator", Description: "Madness saves are made with Advantage permanently. Feral Bonus Damage: +6."}},
		18: {{Name: "Unstoppable Mutation", Description: "You stop aging and cannot be aged magically. Feral Bonus Damage: +6."}},
		19: {{Name: "Ability Score Improvement", Description: "Feral Bonus Damage: +6."}},
		20: {{Name: "Perfect Organism", Description: "No DC Increase for 1st level spells. Feral Bonus Damage: +7."}},
	}
}

func mutagenLevelProgression() map[int]ClassLevelSeed {
	// Mutagen gets Extra Attack at level 5
	// RageDamageBonus represents Feral Bonus Damage scaling
	// MadnessBaseDC is always 10, FeralBonus scales with level
	return map[int]ClassLevelSeed{
		1:  {MadnessBaseDC: 10, FeralBonus: 2, RageDamageBonus: 2},                      // Feral Bonus Damage: +2
		2:  {MadnessBaseDC: 10, FeralBonus: 2, RageDamageBonus: 2},
		3:  {MadnessBaseDC: 10, FeralBonus: 2, RageDamageBonus: 2},
		4:  {MadnessBaseDC: 10, FeralBonus: 2, RageDamageBonus: 2},
		5:  {MadnessBaseDC: 10, FeralBonus: 3, ExtraAttackCount: 1, RageDamageBonus: 3}, // Extra Attack, Feral Bonus Damage: +3
		6:  {MadnessBaseDC: 10, FeralBonus: 3, ExtraAttackCount: 1, RageDamageBonus: 3},
		7:  {MadnessBaseDC: 10, FeralBonus: 3, ExtraAttackCount: 1, RageDamageBonus: 3},
		8:  {MadnessBaseDC: 10, FeralBonus: 3, ExtraAttackCount: 1, RageDamageBonus: 3},
		9:  {MadnessBaseDC: 10, FeralBonus: 4, ExtraAttackCount: 1, RageDamageBonus: 4}, // Feral Bonus Damage: +4
		10: {MadnessBaseDC: 10, FeralBonus: 4, ExtraAttackCount: 1, RageDamageBonus: 4},
		11: {MadnessBaseDC: 10, FeralBonus: 4, ExtraAttackCount: 1, RageDamageBonus: 4},
		12: {MadnessBaseDC: 10, FeralBonus: 4, ExtraAttackCount: 1, RageDamageBonus: 4},
		13: {MadnessBaseDC: 10, FeralBonus: 5, ExtraAttackCount: 1, RageDamageBonus: 5}, // Feral Bonus Damage: +5
		14: {MadnessBaseDC: 10, FeralBonus: 5, ExtraAttackCount: 1, RageDamageBonus: 5},
		15: {MadnessBaseDC: 10, FeralBonus: 5, ExtraAttackCount: 1, RageDamageBonus: 5},
		16: {MadnessBaseDC: 10, FeralBonus: 5, ExtraAttackCount: 1, RageDamageBonus: 5},
		17: {MadnessBaseDC: 10, FeralBonus: 6, ExtraAttackCount: 1, RageDamageBonus: 6}, // Feral Bonus Damage: +6
		18: {MadnessBaseDC: 10, FeralBonus: 6, ExtraAttackCount: 1, RageDamageBonus: 6},
		19: {MadnessBaseDC: 10, FeralBonus: 6, ExtraAttackCount: 1, RageDamageBonus: 6},
		20: {MadnessBaseDC: 10, FeralBonus: 7, ExtraAttackCount: 1, RageDamageBonus: 7}, // Feral Bonus Damage: +7
	}
}
