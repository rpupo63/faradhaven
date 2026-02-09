package faradhaven_classes

// Sanguinist returns the Sanguinist class seed.
// Restructured into a Modular Progression Class with a Notoriety System.
func Sanguinist() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Sanguinist",
		Description:    "A modular class balancing the Medical Prodigy and the Blood Rage tracks. Manage Ichor to fuel healing miracles or devastating vampiric attacks, but beware the Instability Feedback Loop.",
		HitDie:         10,
		PrimaryAbility: "charisma",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/sanguinist.jpg",
		Archetype:      "Modular Notoriety (Doctor vs Predator)",
		Concept:        "A vampire-like being who manages Blood Ichor. Your choices at key levels grant Notoriety Points (MP or BR), creating mechanical tension and punishing players who lean too hard into one side without balance.",
		DnDSkillFocus:    []string{"Medicine", "Stealth"},
		Proficiencies:    "Light armor, Simple weapons, Martial finesse weapons",
		SkillChoice:      []string{"Deception", "Insight", "Persuasion", "Stealth"},
		Tools:            []string{"Poisoner's Kit", "Healer's Kit"},
		SavingThrows:     []string{"Constitution", "Charisma"},
		AutomaticEquipNames: []string{"3 blood vials (empty)"},
		AutomaticWeaponNames: []string{"Bite"},
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
		ArchetypeLevel:      nil, // Modular progression instead of rigid archetype
		Archetypes:          nil,
		ResourceType:        "blood_ichor",
		ResourceName:        "Blood Ichor",
		ResourceRestoreType: "long_rest",
	}
}

func sanguinistLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "Unarmored Defense", Description: "While not wearing armor, your AC equals 10 + your Dexterity modifier + your Charisma modifier. You cannot use a shield with this feature."},
			{Name: "The Thirst", Description: "You have a pool of Blood Ichor equal to Level + Charisma modifier. Bite: As an action, make a melee spell attack (+ Prof + Cha). On hit, deal 2d8 necrotic damage and regain Ichor equal to your Proficiency Bonus. Yields 1 Unstable Component (scales at higher levels)."},
			{Name: "The Moral Seesaw (Notoriety System)", Description: "Your choices at levels 3, 7, 11, 15, and 18 grant Notoriety Points in Medical Prodigy (MP) or Blood Rage (BR). Instability Feedback Loop: If your MP is 3+ points higher than BR (Overloaded Healer), your damage-dealing features trigger 'Sanguine Backfires' (heal a nearby enemy for 1d8 or take that damage yourself). If your BR is 3+ points higher than MP (Starving Predator), your healing becomes 'Ravenous' (target takes 1d8 necrotic damage before receiving the heal)."},
			{Name: "Ichor Exhaustion Limit", Description: "If you spend more than 50% of your Max Ichor in a single turn, you must succeed on a DC 13 Constitution Save or gain 1 level of Sanguine Fatigue (Disadvantage on the next attack)."},
		},
		2: {
			{Name: "Sanguine Extraction", Description: "When you Bite a hostile creature or Siphon an ally, you extract unstable Components. Yields: 1 (Lvl 1-4), 2 (Lvl 5-10), 3 (Lvl 11+). Components decay on a Long Rest. As a bonus action, you can expend one Component to regain 1d4 Ichor or gain Advantage on your next saving throw. At Lvl 11+, you also gain 1 Ichor Regen per turn for 2 turns after an extraction."},
			{Name: "Blood Graft", Description: "As an action, spend 2 Ichor to heal a creature within 5ft for 1d8 + Charisma modifier hit points."},
			{Name: "Shadow Mist", Description: "As an action, spend 2 Ichor to create a 10ft radius cloud centered on you. Creatures must succeed on a Constitution save (DC 8 + Prof + Cha) or take 2d6 necrotic damage and be blinded until the start of your next turn."},
		},
		3: {
			{Name: "Modular Choice: Level 3", Description: "Choose one: (A) Doctor's Intuition (+1 MP) or (B) Primal Senses (+1 BR)."},
			{Name: "Choice A: Doctor's Intuition (+1 MP)", Description: "Blood Graft can now be used as a Bonus Action if the target is at 0 HP."},
			{Name: "Choice B: Primal Senses (+1 BR)", Description: "Gain Darkvision 60ft. Your Bite grants a +10ft movement speed buff for the current turn."},
		},
		5: {
			{Name: "Extra Attack", Description: "You can attack twice when you take the Attack action."},
			{Name: "Sanguine Recovery", Description: "You regain 2 Ichor when you finish a Short Rest."},
		},
		6: {
			{Name: "Renfield's Devotion", Description: "You can bind one humanoid as your 'Renfield'. While they are within 30 feet, you have advantage on Persuasion and Intimidation. Once per turn, they can use their reaction to grant you advantage on one attack roll or saving throw."},
			{Name: "Coagulation", Description: "As an action, spend 3 Ichor to target a creature within 30ft. It must succeed on a Strength save (DC 8 + Prof + Cha) or be paralyzed for 1 minute (save repeats at end of turn)."},
		},
		7: {
			{Name: "Modular Choice: Level 7", Description: "Choose one: (A) Field Medic (+1 MP) or (B) Ichor Lash (+1 BR)."},
			{Name: "Choice A: Field Medic (+1 MP)", Description: "Requires 1 MP. Spend 4 Ichor to cast Mass Blood Graft. Heal up to 3 allies for 2d8 + Charisma."},
			{Name: "Choice B: Ichor Lash (+1 BR)", Description: "Requires 1 BR. Spend 3 Ichor to extend your Bite reach to 15ft using a whip of blood. If it hits, the target is pulled 10ft toward you."},
		},
		9: {{Name: "Shadow Mist Expansion", Description: "Shadow Mist radius increases to 15 feet. Damage increases to 3d6 necrotic, and failed saves also result in the target being Restrained."}},
		10: {{Name: "Eternal Covenant", Description: "Perform a 1-hour ritual on a deceased creature. You gain a +2 bonus to AC and advantage on all saving throws for 24 hours. The creature cannot be resurrected by any means short of a 'Wish' spell."}},
		11: {
			{Name: "Modular Choice: Level 11", Description: "Choose one: (A) Panacea (+1 MP) or (B) Heart-Seeker (+1 BR)."},
			{Name: "Choice A: Panacea (+1 MP)", Description: "Blood Graft now removes the Poisoned, Diseased, or Stunned condition from the target."},
			{Name: "Choice B: Heart-Seeker (+1 BR)", Description: "Your Bite deals an extra 1d8 damage if the target is below its hit point maximum."},
		},
		13: {
			{Name: "Blood Graft Mastery", Description: "Blood Graft healing increases to 3d8 + Charisma."},
			{Name: "Empowered Siphon", Description: "As a bonus action, Siphon a creature within 5ft. They take 2d10 necrotic damage, and you regain double the normal Ichor amount."},
		},
		14: {{Name: "Coagulation Mastery", Description: "Coagulation cost is reduced to 2 Ichor. The target has disadvantage on its first saving throw against the effect."}},
		15: {
			{Name: "Modular Choice: Level 15", Description: "Choose one: (A) Transfusion Mastery (+1 MP) or (B) Predator's Resolve (+1 BR)."},
			{Name: "Choice A: Transfusion Mastery (+1 MP)", Description: "Requires 2 MP. When you use Blood Graft, you can choose to take 1d10 psychic damage to heal the target for an additional 2d10 hit points."},
			{Name: "Choice B: Predator's Resolve (+1 BR)", Description: "Requires 2 BR. After using Siphon, you have advantage on all attack rolls and ability checks for 10 minutes."},
			{Name: "Master Extractions", Description: "Extraction yields 3 components from enemies and 5 from allies. The feedback loop intensifies: If your MP and BR are within 1 point of each other (Balanced), you no longer suffer the Instability Feedback Loop. If you have 4+ points in one track (Specialist), the Backfire/Ravenous dice increases to 2d10."},
		},
		17: {{Name: "Vital Drain", Description: "When you heal with Blood Graft, you can deal necrotic damage equal to half the amount healed to a creature within 30 feet. Additionally, dealing necrotic damage with Bite restores 1 Ichor for every 10 damage dealt."}},
		18: {
			{Name: "Modular Choice: Level 18", Description: "Choose one: (A) Greater Panacea (+1 MP) or (B) Blood-Starved Beast (+1 BR)."},
			{Name: "Choice A: Greater Panacea (+1 MP)", Description: "Requires 3 MP. You can use Blood Graft on a creature that has been dead for no more than 1 minute. They return to life with 1 hit point."},
			{Name: "Choice B: Blood-Starved Beast (+1 BR)", Description: "Requires 3 BR. When you use Bite, you can make one additional melee weapon attack as part of the same action."},
			{Name: "Ichor Overflow", Description: "At the start of your turn, if you have 0 Ichor and have taken damage since your last turn, you regain 2 Ichor."},
		},
		20: {{Name: "Ekon Ascension", Description: "When you roll initiative and have less than half your max Ichor, you regain Ichor until you are at half. You can use Bite as a bonus action. Your Charisma increases by 4, to a maximum of 24."}},
	}
}

func sanguinistLevelProgression() map[int]ClassLevelSeed {
	return map[int]ClassLevelSeed{
		1:  {MaxBloodIchor: 4, BiteDamageDice: 2}, // Level 1 + 3 Cha
		2:  {MaxBloodIchor: 5, BiteDamageDice: 2},
		3:  {MaxBloodIchor: 6, BiteDamageDice: 2},
		4:  {MaxBloodIchor: 7, BiteDamageDice: 2},
		5:  {MaxBloodIchor: 8, BiteDamageDice: 2, ExtraAttackCount: 1},
		6:  {MaxBloodIchor: 9, BiteDamageDice: 2, ExtraAttackCount: 1},
		7:  {MaxBloodIchor: 10, BiteDamageDice: 2, ExtraAttackCount: 1},
		8:  {MaxBloodIchor: 11, BiteDamageDice: 2, ExtraAttackCount: 1},
		9:  {MaxBloodIchor: 12, BiteDamageDice: 2, ExtraAttackCount: 1},
		10: {MaxBloodIchor: 13, BiteDamageDice: 2, ExtraAttackCount: 1},
		11: {MaxBloodIchor: 14, BiteDamageDice: 2, ExtraAttackCount: 1},
		12: {MaxBloodIchor: 15, BiteDamageDice: 2, ExtraAttackCount: 1},
		13: {MaxBloodIchor: 16, BiteDamageDice: 2, ExtraAttackCount: 1},
		14: {MaxBloodIchor: 17, BiteDamageDice: 2, ExtraAttackCount: 1},
		15: {MaxBloodIchor: 18, BiteDamageDice: 2, ExtraAttackCount: 1},
		16: {MaxBloodIchor: 19, BiteDamageDice: 2, ExtraAttackCount: 1},
		17: {MaxBloodIchor: 20, BiteDamageDice: 2, ExtraAttackCount: 1},
		18: {MaxBloodIchor: 21, BiteDamageDice: 2, ExtraAttackCount: 1},
		19: {MaxBloodIchor: 22, BiteDamageDice: 2, ExtraAttackCount: 1},
		20: {MaxBloodIchor: 23, BiteDamageDice: 2, ExtraAttackCount: 1},
	}
}