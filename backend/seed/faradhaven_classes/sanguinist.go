package faradhaven_classes

// Sanguinist returns the Sanguinist class seed.
// Inspired by Vampyr: a high-stakes vampire class centered on Blood Ichor,
// moral choice (Healer vs Predator), and the tension between saving allies and consuming them for power.
func Sanguinist() FaradhavenClassSeed {
	archetypeLevel := 3
	return FaradhavenClassSeed{
		Name:           "The Sanguinist",
		Description:    "Manage Blood Ichor to fuel both healing miracles and devastating vampiric attacks. Every encounter forces a moral choice—save your allies with blood magic or siphon them for overwhelming predatory power.",
		HitDie:         10,
		PrimaryAbility: "charisma",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/sanguinist.jpg",
		Archetype:      "Healer / Predator / High-Stakes Controller",
		Concept:        "A vampire-like being who must manage Blood Ichor to fuel both life-saving medical abilities and devastating supernatural attacks. Every encounter becomes a moral choice: heal your allies or siphon them for power.",
		ClassFeatures: []string{
			"The Thirst: You possess a pool of Blood Ichor points equal to your Level + Charisma Modifier. Expend these to use class features. Regain all Ichor after a Long Rest. As an action, you can Bite a creature that is Grappled, Incapacitated, or Restrained: deal 1d6 piercing damage and regain Ichor equal to your proficiency bonus.",
			"The Moral Seesaw: The Healer (Preservation)—if you use Ichor to heal an ally, your next Killing ability costs 1 less Ichor (minimum 1). The Predator (Consumption)—you can Siphon a willing or unconscious ally: they take 1d10 necrotic damage (cannot be reduced), you gain maximum Ichor, you have advantage on all attack rolls for 1 minute, but the ally gains one level of Exhaustion.",
		},
		DnDSkillFocus:    []string{"Medicine", "Stealth"},
		Proficiencies:    "Light armor, Simple weapons, Martial finesse weapons",
		SkillChoice:      []string{"Deception", "Insight", "Persuasion", "Stealth"},
		Tools:            []string{"Poisoner's Kit", "Healer's Kit"},
		SavingThrows:     []string{"Constitution", "Charisma"},
		AutomaticEquipNames: []string{"3 blood vials (empty)"},
		AutomaticItemNames:  []string{"Leather armor", "Healer's kit"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your surgical weapon",
				Options: []EquipmentOptionSeed{
					{Description: "A Saw Cleaver (Trick Weapon)", WeaponNames: []string{"Saw Cleaver"}},
					{Description: "A Rapier (Precision)", WeaponNames: []string{"Rapier"}},
				},
			},
			{
				Instruction: "Choose your harvesting kit",
				Options: []EquipmentOptionSeed{
					{Description: "Sanguine Extraction Pump", ItemNames: []string{"Sanguine Extraction Pump"}},
					{Description: "Medical Kit and 3 Blood Vials", ItemNames: []string{"Healer's Kit", "Empty Blood Vial", "Empty Blood Vial", "Empty Blood Vial"}},
				},
			},
		},
		LevelFeatures:       sanguinistLevelFeatures(),
		LevelProgression:    sanguinistLevelProgression(),
		ArchetypeLevel:      &archetypeLevel,
		Archetypes:          sanguinistArchetypes(),
		ResourceType:        "blood_ichor",
		ResourceName:        "Blood Ichor",
		ResourceRestoreType: "long_rest", // Regain all Ichor after a Long Rest
	}
}

func sanguinistLevelFeatures() map[int]string {
	return map[int]string{
		1:  "The Thirst, Moral Seesaw — You gain your Blood Ichor pool (Level + Cha modifier) and the core moral choice mechanic. You can Bite grappled/incapacitated/restrained creatures to regain Ichor, and choose between Preservation (heal to discount Killing abilities) or Consumption (Siphon allies for power and advantage).",
		2:  "Blood Graft, Shadow Mist — Blood Graft: As an action, expend 2 Ichor to touch a creature. They regain 1d8 + your Charisma modifier hit points (die scales at higher levels). Shadow Mist: As an action, expend 2 Ichor to create a shadow cloud at a point within 30 feet. Creatures in a 10-foot radius must succeed on a Constitution save or take 2d6 necrotic damage and be Blinded until the end of their next turn.",
		3:  "Blood Path — Choose your path at this level. Your choice determines whether you lean toward preservation or predation, granting you specialized abilities at levels 3 and 10.",
		4:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Charisma cannot exceed 20.",
		5:  "Extra Attack — You can attack twice, instead of once, whenever you take the Attack action on your turn.",
		6:  "Coagulation — You can freeze the blood in a target's veins. Spend 3 Ichor to force a target to make a Strength save. On failure, they are Paralyzed for 1 minute. They can repeat the save at the end of each of their turns.",
		7:  "Blood Graft Surge — Your Blood Graft die increases to 2d8 + Charisma modifier. When you use Blood Graft on an ally at or below half their hit point maximum, they gain temporary hit points equal to your Charisma modifier.",
		8:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Charisma cannot exceed 20.",
		9:  "Shadow Mist Expansion — Your Shadow Mist radius increases to 15 feet and deals 3d6 necrotic damage. Creatures that fail the save are also Restrained until the end of their next turn.",
		11: "Ichor Reservoir — Your maximum Blood Ichor increases by your Charisma modifier (in addition to Level + Cha). Once per long rest, when you reduce a creature to 0 hit points with a Bite, you regain all spent Ichor.",
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Charisma cannot exceed 20.",
		13: "Blood Graft Mastery — Your Blood Graft die increases to 3d8 + Charisma modifier. You can expend 1 additional Ichor to remove one condition (poisoned, diseased, or one level of exhaustion) from the target when you heal them.",
		14: "Coagulation Mastery — Coagulation now costs 2 Ichor. Creatures that fail the save have disadvantage on the first save they make at the end of their turn. You can target a second creature within 30 feet of the first by spending 2 extra Ichor.",
		15: "Predator's Resolve — When you Siphon an ally, the advantage on attack rolls lasts for 10 minutes instead of 1. Your Bite damage increases to 2d6, and you regain Ichor equal to twice your proficiency bonus when you use it.",
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Charisma cannot exceed 20.",
		17: "Vital Drain — When you use Blood Graft to heal an ally, you can choose to also deal necrotic damage equal to the healing to one hostile creature within 30 feet of the target. Alternatively, when you deal necrotic damage with Shadow Mist or Coagulation, you regain 1 Ichor per creature that fails its save.",
		18: "Ichor Overflow — At the start of each of your turns, if you have no Ichor remaining and have taken damage since your last turn, you regain 1 Ichor. Your Moral Seesaw benefits are enhanced: Preservation reduces the next Killing ability cost by 2 (minimum 1); Consumption no longer grants the ally Exhaustion when they are willing.",
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Charisma cannot exceed 20.",
		20: "Ekon Ascension — You have mastered the blood arts. You regain half your maximum Ichor (rounded up) when you roll initiative and have fewer than 3 remaining. Your Bite can target a creature within 5 feet without Grappling them if they are below half their hit point maximum. When you use your subclass's level 10 feature, you can maintain it for an additional minute by spending 3 Ichor as a bonus action.",
	}
}

func sanguinistLevelProgression() map[int]ClassLevelSeed {
	return map[int]ClassLevelSeed{
		1:  {BiteDamageDice: 1},
		5:  {BiteDamageDice: 1, ExtraAttackCount: 1},
		15: {BiteDamageDice: 2, ExtraAttackCount: 1},
	}
}

func sanguinistArchetypes() []ArchetypeSeed {
	return []ArchetypeSeed{
		{
			Name:        "Path of the Doctor",
			Description: "You embrace the role of healer and protector. Your blood magic mends wounds and shields allies, though the temptation to consume remains ever-present.",
			Features: map[int]string{
				3:  "Medical Prodigy — You gain proficiency in Medicine if you don't already have it. When you use Blood Graft to heal an ally, they also gain temporary hit points equal to your Charisma modifier.",
				10: "Transfusion — When you take damage, you can use your reaction to spend 1 Ichor and grant an ally within 30 feet healing equal to half the damage you took. This healing cannot exceed your maximum Ichor pool.",
			},
		},
		{
			Name:        "Path of the Beast",
			Description: "You give in to the predator within. Your vampiric powers grow stronger through consumption, transforming you into a terrifying avatar of shadow and hunger.",
			Features: map[int]string{
				3:  "Blood Rage — When you Siphon an ally, your weapon attacks deal an extra 1d8 necrotic damage for the duration of the advantage buff. This damage increases to 2d8 at level 11.",
				10: "Abyss — As an action, spend 6 Ichor to transform into a shadow avatar for 1 minute. You gain resistance to all damage except radiant and can move through spaces as small as 1 inch without squeezing.",
			},
		},
	}
}
