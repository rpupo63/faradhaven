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
		ClassFeatures: []string{
			"The Madness Mechanic: Instead of Spell Slots, you use Mutation Strains. Every time you cast a spell (Mutation), you must make a Madness Save (Constitution). The DC starts at 10 and increases by +2 with each subsequent cast. Reset: The DC resets to 10 after a Short or Long Rest.",
			"Feral Mode: On a failed Madness save, you enter a Feral State for 1 minute. You cannot cast spells, you lose the ability to distinguish friend from foe, and you must use your action to attack the nearest creature. You gain Resistance to Physical damage and deal extra damage (scaling with level) while Feral.",
			"Internalized Arcana: You can only cast spells with a range of 'Self' or 'Touch'.",
		},
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
				Features: map[int]string{
					3:  "Titan's Growth — You count as one size larger when determining carrying capacity and the weight you can push, drag, or lift. Your Feral Bonus Damage increases by +1.",
					7:  "Living Fortress — While in Feral Mode, you gain a +1 bonus to AC and temporary hit points equal to your Constitution modifier at the start of each of your turns.",
					11: "Siege Breaker — Your melee weapon attacks deal double damage to objects and structures. When you hit a creature with a melee weapon attack, you can push them up to 10 feet away.",
					15: "Colossal Form — When you enter Feral Mode, your size doubles in all dimensions, and your weight is multiplied by eight. You deal an extra 1d4 damage on all melee attacks.",
				},
			},
			{
				Name:        "The Phase-Shifter",
				Description: "Focuses on unstable molecular bonds, allowing for blink-steps, invisibility, and ethereal mutations.",
				Features: map[int]string{
					3:  "Blink Strike — When you take the Attack action, you can teleport up to 10 feet before each attack to an unoccupied space you can see.",
					7:  "Ethereal Step — As a bonus action, you can turn invisible until the start of your next turn or until you attack or cast a spell.",
					11: "Phase Dodge — When you are hit by an attack, you can use your reaction to roll 1d6. On a 4 or higher, the attack misses you as you phase out of reality.",
					15: "Reality Tear — You can cast the 'Blink' spell on yourself without expending a spell slot or increasing your Madness DC. While blinking, your Feral Bonus Damage applies to your first attack after reappearing.",
				},
			},
			{
				Name:        "The Plague-Bearer",
				Description: "Your mutations exude toxic clouds and life-steal mechanics. You are a walking biohazard.",
				Features: map[int]string{
					3:  "Toxic Aura — Creatures that start their turn within 5 feet of you take poison damage equal to your Constitution modifier.",
					7:  "Parasitic Drain — When you deal damage to a creature with a melee attack while in Feral Mode, you regain hit points equal to half the damage dealt.",
					11: "Contagion Carrier — You are immune to disease and poison damage. Your Toxic Aura range increases to 10 feet.",
					15: "Viral Explosion — When you are reduced to 0 hit points, your body releases a cloud of virulent gas. Each creature within 20 feet of you must make a Constitution saving throw (DC 8 + your Constitution modifier + your Proficiency Bonus), taking 8d6 poison damage on a failed save, or half as much on a successful one.",
				},
			},
		},
	}
}

func mutagenLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Internalized Arcana, Feral Mode — Feral Bonus Damage: +2. Madness Resistance (1/Day): You can choose to succeed on a failed Madness save.",
		2:  "Biological Adaptation — Choose one: Darkvision (60ft), Climbing Speed (equal to walking), or Natural Armor (+1 AC). Feral Bonus Damage: +2.",
		3:  "Mutagenic Archetype — Choose a path: The Behemoth, The Phase-Shifter, or The Plague-Bearer. Feral Bonus Damage: +2.",
		4:  "Ability Score Improvement. Feral Bonus Damage: +2.",
		5:  "Extra Attack. Madness Resistance (2/Day). Feral Bonus Damage: +3.",
		6:  "Stabilized Bloodstream — Madness DC increases by only +1 per spell cast. Feral Bonus Damage: +3.",
		7:  "Archetype Feature. Feral Bonus Damage: +3.",
		8:  "Ability Score Improvement. Feral Bonus Damage: +3.",
		9:  "Advanced Synthesis — You gain access to 3rd-level mutations (spells). Feral Bonus Damage: +4.",
		10: "Mind of the Hive — Immunity to Frightened condition. Feral Bonus Damage: +4.",
		11: "Archetype Feature. Feral Bonus Damage: +4.",
		12: "Ability Score Improvement. Feral Bonus Damage: +4.",
		13: "Rapid Cellular Repair — Regain hit points equal to 1d10 + Con modifier as a bonus action (Uses: Con modifier/Long Rest). Feral Bonus Damage: +5.",
		14: "Controlled Ferocity — In Feral Mode, you no longer attack allies and can choose targets (melee only). Feral Bonus Damage: +5.",
		15: "Archetype Feature. Feral Bonus Damage: +5.",
		16: "Ability Score Improvement. Feral Bonus Damage: +5.",
		17: "Apex Predator — Madness saves are made with Advantage permanently. Feral Bonus Damage: +6.",
		18: "Unstoppable Mutation — You stop aging and cannot be aged magically. Feral Bonus Damage: +6.",
		19: "Ability Score Improvement. Feral Bonus Damage: +6.",
		20: "Perfect Organism — No DC Increase for 1st level spells. Feral Bonus Damage: +7.",
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
