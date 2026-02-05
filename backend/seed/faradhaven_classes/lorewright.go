package faradhaven_classes

// Lorewright returns the Lorewright class seed.
// A class that consumes creatures to absorb their memories and abilities,
// balancing power against the risk of madness. High wisdom and constitution,
// with echoes of past lives guiding their actions.
func Lorewright() FaradhavenClassSeed {
	archetypeLevel := 3
	return FaradhavenClassSeed{
		Name:           "The Lorewright",
		Description:    "Consume the flesh of your prey to absorb their memories and abilities. Balance the power you gain against the fracturing of your mind as you become the sum of all you have eaten.",
		HitDie:         8,
		PrimaryAbility: "wisdom",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/lorewright.jpg",
		Archetype:      "Knowledge / Survival / Transformation",
		Concept:        "A hunter-scholar who has learned that true understanding comes through consumption. You eat the livers of your prey to absorb their memories, skills, and eventually their physical traits. Max WIS and CON—your mind must be strong enough to hold the souls of beasts, and your body must process what you consume. Every meal is a gamble between power and madness.",
		ClassFeatures: []string{
			"Anatomical Insight: As a bonus action, make an Insight or Animal Handling check contested by a creature's Deception (intelligent) or Constitution (beast/monstrosity). On success, learn its Damage Immunities, Resistances, or Vulnerabilities, plus its HP relative to yours (Higher, Lower, or Equal).",
			"Visceral Psychometry: Consume the liver of a freshly killed creature (died within 1 hour) over 1 minute to absorb its memories from the last 24 hours with perfect clarity. At Level 7+, access deep memories and ancestral instincts.",
			"The Fracture: When using Visceral Psychometry, make a Wisdom Save (DC = 10 + half creature's CR, minimum 10). On failure, gain 1 Trauma. 1 Trauma: Disadvantage on Charisma checks. 2 Trauma: Gain a temporary flaw from Indefinite Madness table. 3 Trauma: Confusion spell for 1 minute when entering combat. Remove 1 Trauma per Long Rest in a safe sanctuary.",
		},
		DnDSkillFocus:    []string{"Insight", "Survival"},
		Proficiencies:    "Simple weapons, Martial weapons, Light armor, Medium armor, Shields",
		SkillChoice:      []string{"History", "Insight", "Medicine", "Nature", "Survival", "Animal Handling", "Perception"},
		Tools:            []string{"Herbalism Kit", "Cook's Utensils"},
		SavingThrows:     []string{"Wisdom", "Constitution"},
		AutomaticEquipNames: []string{"Preservation kit (salt, jars)", "Traveler's clothes", "Bone talisman"},
		AutomaticItemNames:  []string{"Leather armor", "Explorer's pack"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your hunting tool",
				Options: []EquipmentOptionSeed{
					{Description: "A Hunting Knife (counts as Dagger)", WeaponNames: []string{"Dagger"}},
					{Description: "A Handaxe", WeaponNames: []string{"Handaxe"}},
				},
			},
			{
				Instruction: "Choose your secondary weapon",
				Options: []EquipmentOptionSeed{
					{Description: "A Shortbow and 20 arrows", Items: []string{"Shortbow", "20 Arrows"}, WeaponNames: []string{"Shortbow"}},
					{Description: "A Spear", WeaponNames: []string{"Spear"}},
				},
			},
		},
		LevelFeatures:    lorewrightLevelFeatures(),
		LevelProgression: lorewrightLevelProgression(),
		ArchetypeLevel:   &archetypeLevel,
		Archetypes:       lorewrightArchetypes(),
	}
}

func lorewrightArchetypes() []ArchetypeSeed {
	return []ArchetypeSeed{
		{
			Name:        "Path of the Warlord",
			Description: "You channel the combat memories of fallen warriors, gaining martial prowess and the ability to strike with the fury of a hundred battles.",
			Features: map[int]string{
				3:  "Martial Legacy — You gain proficiency with heavy armor and one martial weapon of your choice. The echoes of warriors past guide your hands in combat.",
				10: "Predator's Strike — When you hit a creature type you have consumed, add your Wisdom modifier to damage. The memories of hunting that prey grant deadly precision.",
				14: "Battle Trance — When you enter combat, you can activate one Echo as a free action to gain its benefits immediately. Additionally, you can use Cycle of Rebirth twice per Long Rest.",
			},
		},
		{
			Name:        "Path of the Sage",
			Description: "You absorb the knowledge and magical insights of your prey, becoming a repository of arcane secrets and ancestral wisdom.",
			Features: map[int]string{
				3:  "Scholar's Legacy — You gain expertise in History or Arcana (double your proficiency bonus), plus two additional languages of your choice. The minds you consume expand your understanding.",
				10: "Deep Insight — When you consume a creature, you can ask two questions instead of one with Phylogenetic Recall. Additionally, you learn one spell the creature could cast (if any) which you can cast once before your next Long Rest using Wisdom.",
				14: "Shared Consciousness — You can maintain two Collective Consciousness transfers simultaneously. When you transfer an Echo, you retain partial benefit (half the bonus or limited version of the trait).",
			},
		},
	}
}

func lorewrightLevelProgression() map[int]ClassLevelSeed {
	// Lorewright gets Extra Attack at level 5
	// Using BardicInspiration to represent Madness Die scaling (d4 -> d6 -> d8 -> d10 -> d12 -> d20)
	return map[int]ClassLevelSeed{
		1:  {BardicInspiration: 4},                       // Madness Die: d4
		5:  {ExtraAttackCount: 1, BardicInspiration: 6},  // Extra Attack, Madness Die: d6
		9:  {ExtraAttackCount: 1, BardicInspiration: 8},  // Madness Die: d8
		13: {ExtraAttackCount: 1, BardicInspiration: 10}, // Madness Die: d10
		17: {ExtraAttackCount: 1, BardicInspiration: 12}, // Madness Die: d12
		20: {ExtraAttackCount: 1, BardicInspiration: 20}, // Madness Die: d20
	}
}

func lorewrightLevelFeatures() map[int]string {
	return map[int]string{
		// Level 1: Core features established in ClassFeatures
		// Echo Slots: 0, Madness Die: d4
		1: "Anatomical Insight, Visceral Psychometry, The Fracture — You gain all core class features. Your Madness Die is d4. You have no Echo Slots yet. High-risk phase: stick to small game (rats, wolves) to avoid madness while learning your craft.",

		// Level 2: Echo Slots: 1, Madness Die: d4
		2: "Somatic Echoes (Skills), Cycle of Rebirth — You gain 1 Echo Slot. When you consume a liver, you can store one Skill Proficiency, Tool Proficiency, or Language the creature possessed. You gain that proficiency while it occupies the slot. Replacing an Echo requires consuming a new liver. Additionally, once per Long Rest, when you fail an Attack Roll, Ability Check, or Saving Throw, roll your Madness Die (d4) and add it to the total—a past life momentarily takes over to correct your mistake.",

		// Level 3: Archetype choice level - shared feature text (archetype-specific features come from Archetypes)
		3: "Past Life Archetype — Choose your archetype at this level. Your choice determines the nature of the memories you absorb most deeply, granting you specialized abilities at levels 3, 10, and 14.",

		// Level 4: Echo Slots: 1, Madness Die: d4
		4: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Wisdom and Constitution are recommended.",

		// Level 5: Echo Slots: 2, Madness Die: d6
		5: "Extra Attack, Stabilized Mind — You can attack twice when you take the Attack action. Your Madness Die increases to d6, and you gain 2 Echo Slots. Your mind adapts to foreign memories: you gain proficiency in Wisdom Saving Throws (if already proficient, gain Charisma Save proficiency instead). The DC for The Fracture is reduced by your Proficiency Bonus. This is the stabilization point—you can now safely consume creatures of CR 5-8.",

		// Level 6: Echo Slots: 2, Madness Die: d6
		6: "Ancestral Guidance — Once per Short Rest, you can enter a trance for 1 minute to commune with the echoes within you. Ask one question about a creature type you have consumed; the DM provides truthful information about their habits, weaknesses, or typical lairs. Additionally, you have advantage on Survival checks to track creature types you have consumed.",

		// Level 7: Echo Slots: 2, Madness Die: d6
		7: "Phylogenetic Recall (Deep Memory) — Your Visceral Psychometry evolves. When you consume, you access the creature's entire lineage. You can ask the DM one question about the creature's pack location, lair layout, or the nature of the entity that created/spawned it. If you eat a Wolf, you don't just see the forest—you feel the path of the Alpha ancestor.",

		// Level 8: Echo Slots: 2, Madness Die: d6
		8: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Wisdom and Constitution are recommended.",

		// Level 9: Echo Slots: 3, Madness Die: d8
		9: "Iron Stomach, Somatic Echoes (Physical Traits) — Your Madness Die increases to d8, and you gain 3 Echo Slots. You can digest any organic matter without ill effect and are immune to ingested poisons. You can now fill an Echo Slot with a physical trait from a consumed creature: Darkvision (60ft), Amphibious (breathe air and water), Climbing/Swimming Speed (equal to walking speed), or Keen Senses (advantage on Perception checks using smell/hearing).",

		// Level 10: Archetype feature level - no shared feature (archetype-specific only)

		// Level 11: Echo Slots: 3, Madness Die: d8
		11: "Assimilated Resistance — You can spend an Echo Slot to gain Resistance to one damage type (Acid, Cold, Fire, Lightning, or Poison) possessed by a creature you consumed. This lasts until you finish a Long Rest or replace the Echo. Your body adapts to the threats your prey faced.",

		// Level 12: Echo Slots: 3, Madness Die: d8
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Wisdom and Constitution are recommended.",

		// Level 13: Echo Slots: 4, Madness Die: d10
		13: "Collective Consciousness — Your Madness Die increases to d10, and you gain 4 Echo Slots. As an action, touch a willing creature to transfer one of your Somatic Echoes to them for 1 hour. You lose access to the Echo for the duration; the ally gains the benefit. You can share your borrowed skills with those who fight beside you.",

		// Level 14: Archetype feature level - no shared feature (archetype-specific only)

		// Level 15: Echo Slots: 4, Madness Die: d10
		15: "Cognitive Bastion — You have mastered the noise in your head. You are immune to the Charmed and Frightened conditions. Your Trauma threshold increases: you do not suffer the debilitating effects of 3 Trauma until you reach 5 Trauma. The voices within are now allies, not threats.",

		// Level 16: Echo Slots: 4, Madness Die: d10
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Wisdom and Constitution are recommended.",

		// Level 17: Echo Slots: 5, Madness Die: d12
		17: "Apex Predator's Soul — Your Madness Die increases to d12, and you gain 5 Echo Slots. As a bonus action, enter a heightened state for 1 minute. During this time, you can use any Action or Bonus Action available to a creature currently stored in your Echo Slots (e.g., a wolf's Pack Tactics or Bite attack logic, a spider's Web, etc.). Once used, you cannot use this feature again until you finish a Long Rest.",

		// Level 18: Echo Slots: 5, Madness Die: d12
		18: "Timeless Body — Your body no longer ages, and you cannot be aged magically. You no longer need food or water (though you can still consume to use your abilities). You have advantage on Death Saving Throws as past lives pull you back from the brink. The countless souls within sustain your mortal form.",

		// Level 19: Echo Slots: 5, Madness Die: d12
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Wisdom and Constitution are recommended.",

		// Level 20: Echo Slots: 5, Madness Die: d20
		20: "The Omega Point — Your Madness Die becomes d20. You are the sum of all your lives and all your meals. You automatically succeed on the Saving Throw for The Fracture—madness can no longer touch you. Furthermore, you can hold two Echoes in a single slot, effectively doubling your capacity for borrowed skills and traits (10 total Echoes in 5 slots). When you die, you can choose to be reborn in a new body at a location of your choice within 1 mile, retaining all memories and Echoes. This rebirth takes 1d10 days.",
	}
}
