package faradhaven_classes

// Ironwright returns the Ironwright class seed
func Ironwright() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Ironwright",
		HitDie:         8,
		PrimaryAbility: "intelligence",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/ironwright.jpg",
		Archetype:      "Summoner / Utility Support",
		Concept:        "An engineer who realizes that machinery is just magic with better documentation. They build temporary automatons and gadgets on the fly.",
		ClassFeatures:  []string{"Rapid Assembly: You use components to build small Rig drones (spiders, copters, turrets) that deliver spell effects.", "Scrap Titan: Higher levels allow you to magnetize scrap metal into a protective exosuit."},
		DnDSkillFocus:  []string{"Investigation", "Arcana"},
		Proficiencies:  "Simple weapons, Light Armor, Hand Crossbows",
		SkillChoice:    []string{"Sleight of Hand", "History", "Medicine", "Investigation"},
		Tools:          []string{"Tinker's Tools", "Smith's Tools"},
		SavingThrows:   []string{"Intelligence", "Dexterity"},
		StartingEquip:  []string{"Wrench (counts as mace)", "Leather apron (light armor)", "Protective goggles", "Bag of gears and springs"},
		LevelFeatures:  ironwrightLevelFeatures(),
	}
}

func ironwrightLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Rapid Assembly (basic): You can spend spell points and components to assemble one Rig drone as a bonus action. Available Rig types: Spider (ground, delivers touch-range effects), Copter (flying, delivers ranged effects). Rigs last until destroyed or dismissed. You can maintain 1 Rig at a time.",
		2:  "Improved Assembly: Your Rigs gain +2 HP per level. Spider Rigs can now climb difficult surfaces without penalty. Copter Rigs gain 10 ft. fly speed.",
		3:  "Turret Rig: You unlock the Turret Rig type—a stationary construct that can deliver one component effect per round within 30 ft. Turrets have higher AC but cannot move.",
		4:  "Ability Score Improvement: Increase one ability score by 2, or two ability scores by 1.",
		5:  "Dual Rig Control: You can maintain 2 Rigs simultaneously. Rigs you create gain resistance to non-magical bludgeoning, piercing, and slashing damage.",
		6:  "Clockwork Precision: Your Rigs add your proficiency bonus to their attack rolls and saving throw DCs when delivering component effects.",
		7:  "Reinforced Frame: Spider and Copter Rigs gain +1 AC. Turret Rigs gain temporary HP equal to your Intelligence modifier when created.",
		8:  "Ability Score Improvement: Increase one ability score by 2, or two ability scores by 1.",
		9:  "Scrap Titan (emerging): You can magnetize nearby scrap metal into light plating. As an action, gain temporary HP equal to 2 × your Cog-Weaver level. This plating lasts 1 minute or until the temp HP is depleted. Usable once per short rest.",
		10: "Triple Rig Control: You can maintain 3 Rigs simultaneously. When you create a Rig, you may imbue it with a second component (both effects trigger on delivery, using one spell point pool).",
		11: "Scrap Titan (improved): Your scrap plating now grants resistance to lightning and fire damage while active. Temp HP increases to 3 × your level.",
		12: "Ability Score Improvement: Increase one ability score by 2, or two ability scores by 1.",
		13: "Siege Turret: You unlock the Siege Turret—a larger Turret Rig with extended range (60 ft.) that can deliver area effects. Costs 2× normal spell points to create.",
		14: "Scrap Titan (full exosuit): Your scrap plating becomes a full exosuit. Gain +2 AC while the temp HP remains. The exosuit can absorb damage directed at an adjacent ally (you choose when to redirect).",
		15: "Quad Rig Control: You can maintain 4 Rigs simultaneously. Your Rigs can now deliver effects with the Clockwork Trigger modifier without consuming an extra action.",
		16: "Ability Score Improvement: Increase one ability score by 2, or two ability scores by 1.",
		17: "Scrap Titan (master): When an ally within 30 ft. takes damage, you can use your reaction to have your exosuit absorb half of it (rounded down) as long as you have temp HP remaining. Your exosuit temp HP increases to 4 × your level.",
		18: "Overclock: Once per long rest, you can command all your active Rigs to act twice in the same round. Each Rig may deliver one component effect per action.",
		19: "Ability Score Improvement: Increase one ability score by 2, or two ability scores by 1.",
		20: "Master Cog-Weaver: You can maintain 5 Rigs simultaneously. Scrap Titan no longer requires a short rest—you can reform the exosuit as an action when it is depleted. Your Rigs are considered magical for the purpose of overcoming resistance. When you create a Rig, you may apply one additional modifier component at no extra cost.",
	}
}
