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
					3:  {{Name: "Titan's Growth", Description: "You count as one size larger for carrying capacity, and your reach with melee weapons increases by 5 feet. Your Feral Bonus Damage increases by an additional +1."}},
					7:  {{Name: "Living Fortress", Description: "While in Feral Mode, you gain a +1 bonus to AC. Additionally, you gain temporary hit points equal to your Constitution modifier at the start of each of your turns."}},
					11: {{Name: "Siege Breaker", Description: "Your melee weapon attacks deal double damage to objects and structures. When you hit a creature, you can use a bonus action to push them 10 feet away if they fail a Strength save (DC 8 + Prof + Con)."}},
					15: {{Name: "Colossal Form", Description: "When you enter Feral Mode, you can choose to become Large. If you do, your melee weapon attacks deal an extra 1d6 damage."}},
				},
			},
			{
				Name:        "The Phase-Shifter",
				Description: "Focuses on unstable molecular bonds, allowing for blink-steps, invisibility, and ethereal mutations.",
				Features: map[int][]FeatureSeed{
					3:  {{Name: "Blink Strike", Description: "When you take the Attack action, you can teleport up to 10 feet to an unoccupied space you can see before each attack."}},
					7:  {{Name: "Ethereal Step", Description: "As a bonus action, you can turn invisible until the start of your next turn or until you attack. You can use this a number of times equal to your Constitution modifier per long rest."}},
					11: {{Name: "Phase Dodge", Description: "When hit by an attack, use your reaction to roll 1d6. On a 4+, the attack misses as you phase out of reality."}},
					15: {{Name: "Reality Tear", Description: "You can cast the 'Blink' spell on yourself without increasing your Madness DC. While blinking, your first attack after reappearing deals an extra 2d6 force damage."}},
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
		1:  "Red Haze: You must attack the nearest creature (friend or foe). Your melee attacks deal an extra 1d6 damage.",
		2:  "Predatory Instinct: You must attack the creature with the lowest current hit points. You have advantage on these attacks.",
		3:  "Gluttonous Consumption: If you reduce a creature to 0 HP, you must spend your next action devouring it. You regain 2d10 HP, but lose your turn.",
		4:  "Bone Shard Eruption: Sharp bone spurs burst from your skin. You take 1d8 piercing damage, but any creature that hits you with a melee attack takes 1d8 piercing damage.",
		5:  "Hyper-Adrenaline: You gain the effects of the Haste spell for 2 rounds. After it ends, you are stunned for 1 round.",
		6:  "Toxic Bile: You vomit a 15ft cone of acid. Each creature must make a Con save (DC 8 + Prof + Con) or take 3d6 acid damage. You take 1d6 acid damage.",
		7:  "Muscle Spasm: Your Strength increases massively (+4 to hit/damage), but your movement speed is reduced to 5 feet.",
		8:  "Blind Agony: You are blinded by pain. Your attacks are at disadvantage, but you are immune to being Charmed or Frightened.",
		9:  "Leaping Hunter: Your jump distance is tripled. You must use your movement to leap towards the furthest enemy you can see and attack them.",
		10: "Acidic Blood: When you take damage, you spray acidic blood. All creatures within 5ft (including allies) take 1d4 acid damage.",
		11: "Hardened Carapace: Your skin turns into a thick shell. Your AC increases by 3, but you have disadvantage on all Dexterity saving throws.",
		12: "Territorial Roar: You must use your action to let out a blood-curdling roar. All creatures within 30ft must make a Wis save or be Frightened of you for 1 minute.",
		13: "Vampiric Fangs: Your jaw unhinges. You must use your action to make a Bite attack (+Prof + Con to hit, 1d10 + Con damage). You regain HP equal to the damage dealt.",
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
			{Name: "The Madness Mechanic", Description: "Instead of Spell Slots, you use Mutation Strains (spells with range Self or Touch). Every time you cast a Mutation, make a DC 10 Constitution saving throw. On failure, you enter Feral Mode. The DC increases by +2 (+1 starting at level 6) for each successful cast. Reset: DC 10 after a Short or Long Rest."},
			{Name: "Feral Mode", Description: "On a failed Madness save, you enter a Feral State for 1 minute. Roll on the Mutagen Feral Table to determine your specific mutation. While feral, you cannot cast Mutations, gain Resistance to Bludgeoning, Piercing, and Slashing damage, and deal extra damage based on your level."},
			{Name: "Internalized Arcana", Description: "Your body is the catalyst. You can only cast Mutations (spells) with a range of 'Self' or 'Touch'."},
			{Name: "Madness Resistance (1/Day)", Description: "Once per day, you can choose to succeed on a failed Madness saving throw."},
		},
		2:  {{Name: "Biological Adaptation", Description: "Choose one: Darkvision (60ft), Climbing Speed (equal to walking), or Natural Armor (+1 bonus to AC)."}},
		3:  {{Name: "Mutagenic Archetype", Description: "Choose a path: The Behemoth, The Phase-Shifter, or The Plague-Bearer. Each path grants unique mutations and physical changes."}},
		5:  {{Name: "Extra Attack, Madness Resistance (2/Day)", Description: "You can attack twice when you take the Attack action. You can now use Madness Resistance twice per day."}},
		6:  {{Name: "Stabilized Bloodstream", Description: "Your body adapts to the strain. The Madness DC now increases by only +1 for each successful Mutation cast."}},
		9:  {{Name: "Advanced Synthesis", Description: "You gain access to 3rd-level Mutations. Your Feral Bonus Damage increases to +4."}, {Name: "Mutation Potency", Description: "When you cast a Mutation, you can choose to increase the Madness DC by an additional +5 to cast it at one level higher than its base."}},
		10: {{Name: "Mind of the Hive", Description: "You are immune to the Frightened condition, and you have advantage on saving throws against being Charmed."}},
		13: {{Name: "Rapid Cellular Repair", Description: "As a bonus action, you regain hit points equal to 1d10 + your Constitution modifier. You can use this a number of times equal to your Constitution modifier per Long Rest."}},
		14: {{Name: "Controlled Ferocity", Description: "In Feral Mode, you no longer attack allies. You can choose your targets freely among hostile creatures (melee only)."}},
		17: {{Name: "Apex Predator", Description: "You gain advantage on all Madness saving throws. Your Feral Bonus Damage increases to +6."}, {Name: "Vicious Reflexes", Description: "While in Feral Mode, you can make one additional melee weapon attack as a bonus action."}},
		18: {{Name: "Unstoppable Mutation", Description: "You stop aging, cannot be aged magically, and are immune to exhaustion caused by casting Mutations."}},
		20: {{Name: "Perfect Organism", Description: "1st level Mutations no longer increase your Madness DC. Your Constitution score increases by 4, to a maximum of 24."}},
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
