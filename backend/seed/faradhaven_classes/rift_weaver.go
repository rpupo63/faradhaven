package faradhaven_classes

// RiftWeaver returns the Rift Weaver class seed
func RiftWeaver() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Rift Weaver",
		HitDie:         8,
		PrimaryAbility: "intelligence",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/rift_weaver.jpg",
		Archetype:      "Full Caster / Evoker / Damage",
		Concept:        "A scholar who learned to open rifts to elemental planes. They draw power from fire, ice, lightning, and earth—commanding raw elemental forces for devastating evocation magic. You are the party's artillery.",
		ClassFeatures: []string{
			"Elemental Channeling: You do not have spell slots. Instead, you spend spell points to channel components drawn from elemental rifts. Your components manifest as fire, cold, lightning, or earth effects—blasts, cones, and area devastation.",
			"Elemental Sight: You can perceive elemental auras within 30 feet—creatures or objects with fire, cold, lightning, or earth affinity appear as a faint glow. You have advantage on Arcana checks to identify elemental effects.",
		},
		DnDSkillFocus: []string{"Arcana", "Nature"},
		Proficiencies: "Simple weapons, Daggers, Quarterstaffs, Light Armor",
		SkillChoice:   []string{"Arcana", "History", "Nature", "Religion"},
		Tools:         []string{"Arcane Focus (orb or wand)"},
		SavingThrows:  []string{"Intelligence", "Wisdom"},
		StartingEquip: []string{"Scholar's robes (light armor)", "Quarterstaff", "Orb or wand focus", "Spellbook (grimoire of elemental formulae)", "Explorer's pack"},
		LevelFeatures: riftWeaverLevelFeatures(),
	}
}

func riftWeaverLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Elemental Sense — You can detect the presence of elemental energy within 60 feet. You gain advantage on Arcana checks to identify elemental effects and resistances.",
		2:  "Ritual Attunement — You can cast any elemental component you know as a ritual, taking 10 minutes and spending no spell points. The effect's duration is halved when cast this way.",
		3:  "Overchannel — When you spend spell points on an elemental damage component, you can spend 3 additional spell points to add your Intelligence modifier to the damage. Once per short rest.",
		4:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Intelligence cannot exceed 20.",
		5:  "Elemental Surge — Once per short rest, when you spend spell points on an elemental component, you can double its range or area of effect. Alternatively, change the damage type to another element you know (fire, cold, lightning, or earth).",
		6:  "Potent Evocation — Your elemental damage components ignore resistance to one damage type of your choice (fire, cold, lightning, or earth). You can change this choice when you finish a long rest.",
		7:  "Elemental Resilience — You have resistance to one damage type of your choice (fire, cold, lightning, or earth). When you would take that damage type, you can spend 5 spell points as a reaction to reduce it by half (before applying resistance).",
		8:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Intelligence cannot exceed 20.",
		9:  "Sculpt Spells — When you cast an elemental component that affects an area, you can choose a number of creatures equal to 1 + your Intelligence modifier. Those creatures automatically succeed on the save and take no damage if they would otherwise be in the area.",
		10: "Elemental Amplification — When you apply an elemental modifier to an ally, the effect lasts 1 hour instead of its normal duration. When you deal elemental damage, you can spend 5 spell points to add 1d6 damage of the same type.",
		11: "Elemental Barrage — When you use an elemental damage component, you can spend 5 additional spell points to target a second creature within 30 feet of the first with the same effect.",
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Intelligence cannot exceed 20.",
		13: "Elemental Veil — You can spend 15 spell points as an action to wreathe yourself in elemental energy for 1 minute. Choose fire, cold, lightning, or earth. Creatures that touch you or hit you with a melee attack take 2d6 damage of that type. You have resistance to that damage type.",
		14: "Dual Element — When you cast an elemental damage component, you can spend 10 additional spell points to deal half the damage as a second element of your choice. Both damage types apply.",
		15: "Elemental Mastery — Your elemental components that force a saving throw impose disadvantage on the first save. You can maintain two elemental effects (e.g., a wall and a zone) simultaneously.",
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Intelligence cannot exceed 20.",
		17: "Elemental Overload — Once per long rest, you can spend all remaining spell points to create an elemental burst. Each creature of your choice within 30 feet takes damage of a type you choose (fire, cold, lightning, or earth) equal to the points spent. Creatures in the area make a Dexterity save for half.",
		18: "Rift Dominion — When an enemy fails a save against your elemental component, you can spend 10 spell points to extend the effect's duration by 1 minute or cause the damage to deal maximum dice. Once per target per long rest.",
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Intelligence cannot exceed 20.",
		20: "Archmage of the Rifts — You are one with the elements. You regain 20 spell points when you roll initiative and have none remaining. Your Elemental Sight range increases to 120 feet, and you automatically know the elemental resistances and vulnerabilities of creatures within that range. Your elemental damage ignores immunity once per long rest per creature.",
	}
}
