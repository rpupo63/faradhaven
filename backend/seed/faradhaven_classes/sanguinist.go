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

func sanguinistLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "Unarmored Defense", Description: "While you are not wearing armor, your AC equals 10 + your Dexterity modifier + your Charisma modifier. You cannot use a shield and gain this benefit."},
			{Name: "The Thirst", Description: "You possess a pool of Blood Ichor points equal to your Level + Charisma Modifier. You regain all expended Ichor when you finish a Long Rest. Bite: As an action, you can make a melee spell attack against a creature within 5 feet. On a hit, the target takes 1d6 necrotic damage, and you regain Ichor equal to your Proficiency Bonus."},
			{Name: "The Moral Seesaw", Description: "The Healer (Preservation): If you spend Ichor to heal a creature, the next time you deal damage with a Class Feature within 1 minute, that damage is increased by your Charisma modifier. The Predator (Consumption): As a bonus action, you can Siphon a willing or unconscious humanoid within 5 feet. They take 1d10 necrotic damage (which cannot be reduced), and you regain Ichor equal to half your Sanguinist Level + Charisma Modifier. A creature Siphoned this way cannot be Siphoned again until it finishes a Long Rest."},
		},
		2: {
			{Name: "Sanguine Extraction", Description: "You learn to harvest magical essences from the blood of others. When you use your Bite feature on a hostile creature or Siphon a willing ally, you extract one unstable component (determined by the DM or the Faradhaven Component Table).\n\nCapacity: You can carry a number of unstable components equal to your Proficiency Bonus. If you extract a new one while at max capacity, you must drop one.\n\nActivation: You can use a Bonus Action to consume a component and trigger its effect.\n\nDecay: All unstable components are destroyed when you finish a Long Rest."},
			{Name: "Blood Graft", Description: "As an action, you can expend 2 Ichor to touch a creature and knit their wounds. The creature regains hit points equal to 1d8 + your Charisma modifier."},
			{Name: "Shadow Mist", Description: "Expend 2 Ichor to create a 10ft radius cloud. Each creature in the area must succeed on a Constitution save (DC 8 + Prof + Cha) or take 2d6 necrotic damage and be blinded until the start of your next turn. Success halves damage and avoids blindness."},
		},
		3: {{Name: "Archetype", Description: "At 3rd level, you choose a path that defines your relationship with your curse: The Path of the Doctor or The Path of the Beast."}},
		4: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1. Charisma cannot exceed 20."}},
		5: {{Name: "Extra Attack", Description: "You can attack twice when you take the Attack action on your turn."}},
		6: {
			{Name: "Renfield's Devotion", Description: "You can bind a humanoid to your service as a 'Renfield' (hype man). They grant advantage on Persuasion/Intimidation and can use their reaction to grant you advantage on one attack or saving throw per round."},
			{Name: "Coagulation", Description: "Spend 3 Ichor to target a creature within 30 feet. The target must succeed on a Strength save (DC 8 + Prof + Cha) or be paralyzed for 1 minute. The target can repeat the save at the end of each of its turns."},
		},
		7: {
			{Name: "Blood Graft Surge", Description: "Blood Graft die increases to 2d8 + Charisma. Healing an ally below half health grants temporary HP."},
			{Name: "Enhanced Extraction", Description: "Your Sanguine Extraction yields more: 2 components from enemies (2 charges each) and 3 components from allies (2 charges each)."},
		},
		8:  {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},
		9:  {{Name: "Shadow Mist Expansion", Description: "Radius increases to 15 feet, deals 3d6 necrotic damage, and restrains targets."}},
		10: {{Name: "Eternal Covenant", Description: "Perform a 1-hour sacrifice on a downed creature to prevent their resurrection and gain a powerful blessing (+2 AC, +2 damage, or Save advantage) for 1d6 days."}},
		11: {{Name: "Ichor Reservoir", Description: "Maximum Blood Ichor increases by Charisma. Reducing a creature to 0 HP with a Bite restores all Ichor (once per long rest)."}},
		12: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},
		13: {{Name: "Blood Graft Mastery", Description: "Blood Graft die increases to 3d8. You can expend 1 extra Ichor to remove a condition (poisoned, diseased, or exhaustion)."}},
		14: {{Name: "Coagulation Mastery", Description: "Coagulation costs 2 Ichor and targets have disadvantage on the first save to end it."}},
		15: {
			{Name: "Predator's Resolve", Description: "Siphon advantage lasts 10 minutes. Bite damage increases to 2d6, and you regain double proficiency in Ichor."},
			{Name: "Master Extraction", Description: "Your Sanguine Extraction is perfected: 3 components from enemies (3 charges) and 5 components from allies (especially your Renfield) with 3 charges each."},
		},
		16: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},
		17: {{Name: "Vital Drain", Description: "When you heal an ally with Blood Graft, you can deal necrotic damage equal to half the amount healed to an enemy within 30 feet. Alternatively, dealing necrotic damage with class features restores 1 Ichor for every 10 damage dealt (minimum 1)."}},
		18: {{Name: "Ichor Overflow", Description: "Regain 1 Ichor at the start of your turn if you have none and took damage. Moral Seesaw benefits are enhanced: Preservation discount is 2; Consumption doesn't exhaust willing allies."}},
		19: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},
		20: {{Name: "Ekon Ascension", Description: "Regain half max Ichor when rolling initiative if low. Bite can be used as a bonus action against targets below half health."}},
	}
}

func sanguinistLevelProgression() map[int]ClassLevelSeed {
	return map[int]ClassLevelSeed{
		1:  {BiteDamageDice: 1},
		2:  {BiteDamageDice: 1},
		3:  {BiteDamageDice: 1},
		4:  {BiteDamageDice: 1},
		5:  {BiteDamageDice: 1, ExtraAttackCount: 1},
		6:  {BiteDamageDice: 1, ExtraAttackCount: 1},
		7:  {BiteDamageDice: 1, ExtraAttackCount: 1},
		8:  {BiteDamageDice: 1, ExtraAttackCount: 1},
		9:  {BiteDamageDice: 1, ExtraAttackCount: 1},
		10: {BiteDamageDice: 1, ExtraAttackCount: 1},
		11: {BiteDamageDice: 1, ExtraAttackCount: 1},
		12: {BiteDamageDice: 1, ExtraAttackCount: 1},
		13: {BiteDamageDice: 1, ExtraAttackCount: 1},
		14: {BiteDamageDice: 1, ExtraAttackCount: 1},
		15: {BiteDamageDice: 2, ExtraAttackCount: 1},
		16: {BiteDamageDice: 2, ExtraAttackCount: 1},
		17: {BiteDamageDice: 2, ExtraAttackCount: 1},
		18: {BiteDamageDice: 2, ExtraAttackCount: 1},
		19: {BiteDamageDice: 2, ExtraAttackCount: 1},
		20: {BiteDamageDice: 2, ExtraAttackCount: 1},
	}
}

func sanguinistArchetypes() []ArchetypeSeed {
	return []ArchetypeSeed{
		{
			Name:        "Path of the Doctor",
			Description: "You embrace the role of healer and protector. Your blood magic mends wounds and shields allies, though the weight of preservation takes a physical toll.",
			Features: map[int][]FeatureSeed{
				3:  {{Name: "Medical Prodigy & Fragile Vessel", Description: "Pro: Gain proficiency in Medicine and Insight. Blood Graft heals nearby allies for half. Con: Disadvantage on Strength checks/saves and max HP reduced by 1 per level."}},
				14: {{Name: "Transfusion & Empathic Shield", Description: "Pro: Bonus to Passive Perception/Insight. When you take damage, you can use a reaction to spend 1 Ichor and heal an ally for half that damage."}},
			},
		},
		{
			Name:        "Path of the Beast",
			Description: "You give in to the predator within. Your vampiric powers grow stronger through consumption, transforming you into a terrifying avatar of shadow and hunger.",
			Features: map[int][]FeatureSeed{
				3:  {{Name: "Blood Rage & The Red Mist", Description: "Pro: Darkvision (60ft) and +2 Perception (super hearing and eyesight). Kills grant 'Primal Hunter' stacks (+2 necrotic damage). Con: At 3 stacks, must save vs Wisdom or enter a frenzy attacking nearest creature."}},
				14: {{Name: "Abyss & Unquenchable Thirst", Description: "Pro: Transformation into a shadow avatar. Con: Must save vs Charisma in social scenes with wounded NPCs or be forced to Bite."}},
			},
		},
	}
}
