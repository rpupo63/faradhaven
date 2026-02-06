package faradhaven_classes

// PowderMage returns the Powder Mage class seed
func PowderMage() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Powder Mage",
		Description:    "The 'Continuous Ignition' Powder Mage. Start speaking as soon as you declare the Cast action; whatever you say before the timer runs out is what happens. You aren't fighting to make the spell work, but to pack as much power into the casting window as possible.",
		HitDie:         8,
		PrimaryAbility: "dexterity",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/powder_mage.jpg",
		Archetype:      "Speed Caster / Skirmisher",
		Concept:        "Continuous Ignition: Real-time verbal casting under pressure.",
		ClassFeatures: []string{
			"Flash-Point Casting: Real-time casting window. Say words before the buzzer.",
			"Kinetic Recoil: Use unspent power for movement.",
			"Speed Dial: Instant cast of pre-saved component strings.",
		},
		DnDSkillFocus:    []string{"Performance", "Acrobatics"},
		Proficiencies:    "Light Armor, Simple Weapons",
		SkillChoice:      []string{"Performance", "Acrobatics", "Arcana", "Sleight of Hand"},
		Tools:            []string{"Alchemist's Supplies"},
		SavingThrows:     []string{"Dexterity", "Charisma"},
		AutomaticEquipNames: []string{"Component Pouch"},
		AutomaticItemNames:  []string{"Light Armor", "Explorer's Pack"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your firearm",
				Options: []EquipmentOptionSeed{
					{Description: "A Pepperbox Revolver", WeaponNames: []string{"Pepperbox Revolver"}},
					{Description: "A Steam-Pressure Rifle", WeaponNames: []string{"Steam-Pressure Rifle"}},
				},
			},
			{
				Instruction: "Choose your casting gear",
				Options: []EquipmentOptionSeed{
					{Description: "Clockwork Chronometer", ItemNames: []string{"Clockwork Chronometer"}},
					{Description: "Aetheric Component Pouch", ItemNames: []string{"Aetheric Component Pouch"}},
				},
			},
		},
		LevelFeatures:       powderMageLevelFeatures(),
		LevelProgression:    powderMageLevelProgression(),
		ResourceType:        "timer",
		ResourceName:        "Casting Timer",
		ResourceRestoreType: "none", // Timer is always available, not a consumable
	}
}

func powderMageLevelProgression() map[int]ClassLevelSeed {
	// Powder Mage gains movement speed bonuses and timer/speed dial upgrades
	return map[int]ClassLevelSeed{
		1:  {TimerDuration: 2},                                        // Flash-Point Casting: 2 seconds
		2:  {TimerDuration: 2, UnarmoredMovement: 10},                  // Powder Sprint: +10 feet
		3:  {TimerDuration: 2, UnarmoredMovement: 10, SpeedDialSlots: 1}, // Speed Dial: 1 saved string
		4:  {TimerDuration: 2, UnarmoredMovement: 10, SpeedDialSlots: 1},
		5:  {TimerDuration: 3, UnarmoredMovement: 10, SpeedDialSlots: 1}, // Timer increases to 3 seconds
		6:  {TimerDuration: 3, UnarmoredMovement: 15, SpeedDialSlots: 1}, // Recoil Mastery
		7:  {TimerDuration: 3, UnarmoredMovement: 15, SpeedDialSlots: 1},
		8:  {TimerDuration: 3, UnarmoredMovement: 15, SpeedDialSlots: 1},
		9:  {TimerDuration: 3, UnarmoredMovement: 15, SpeedDialSlots: 2}, // Speed Dial: 2 saved strings
		10: {TimerDuration: 3, UnarmoredMovement: 15, SpeedDialSlots: 2},
		11: {TimerDuration: 4, UnarmoredMovement: 20, SpeedDialSlots: 2}, // Timer increases to 4 seconds
		12: {TimerDuration: 4, UnarmoredMovement: 20, SpeedDialSlots: 2},
		13: {TimerDuration: 4, UnarmoredMovement: 20, SpeedDialSlots: 2},
		14: {TimerDuration: 4, UnarmoredMovement: 20, SpeedDialSlots: 2},
		15: {TimerDuration: 4, UnarmoredMovement: 20, SpeedDialSlots: 3}, // Speed Dial: 3 saved strings
		16: {TimerDuration: 4, UnarmoredMovement: 20, SpeedDialSlots: 3},
		17: {TimerDuration: 5, UnarmoredMovement: 20, SpeedDialSlots: 3}, // Timer increases to 5 seconds
		18: {TimerDuration: 5, UnarmoredMovement: 30, SpeedDialSlots: 3}, // Sound Barrier: flying speed
		19: {TimerDuration: 5, UnarmoredMovement: 30, SpeedDialSlots: 3},
		20: {TimerDuration: 5, UnarmoredMovement: 30, SpeedDialSlots: 3},
	}
}

func powderMageLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Flash-Point Casting (2 Seconds) — The DM starts a timer when you cast. You list components; the spell resolves based on what you finish saying before the buzzer. The timer acts as a potency cap. (The Snap-Fire).",
		2:  "Powder Sprint & Disengage — Movement speed increases by 10 feet. You can use Disengage or Dash as a bonus action.",
		3:  "Speed Dial (The 'Pre-Mix') — Save one specific string of components. Once per short rest, cast this saved string as an action without using the timer.",
		4:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1.",
		5:  "Kinetic Recoil (Traversal) — When you finish a spell, unspent powder creates a blast. Move 10 feet in any direction immediately after resolution (no opportunity attacks). Timer increases to 3 Seconds (The Compound).",
		6:  "Recoil Mastery — You can move along vertical surfaces during your Kinetic Recoil movement without falling during the move.",
		7:  "Evasion — When you are subjected to an effect that allows you to make a Dexterity saving throw to take only half damage, you instead take no damage if you succeed on the saving throw, and only half damage if you fail.",
		8:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1.",
		9:  "Speed Dial Upgrade — You can now save 2 component strings.",
		10: "Breath Control — You can hold your breath for up to 10 minutes. You have advantage on saving throws against inhaled gases and effects.",
		11: "Room Clearer — As an action, move up to full speed. Any creature you pass within 5ft can be targeted by a single-component spell (e.g., 'Fire') if you say the word as you pass. Timer increases to 4 Seconds (The Sequence).",
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1.",
		13: "Sonic Boom — When you use Kinetic Recoil, you can choose to deal 3d6 thunder damage to all creatures within 5 feet of your starting position (Dexterity save for half).",
		14: "Echoing Components — If you use the same component in two consecutive turns, that component is 'Active' automatically (no need to say it again), freeing up your timer.",
		15: "Speed Dial Upgrade — You can now save 3 component strings.",
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1.",
		17: "Timer Upgrade (5 Seconds) — The Avalanche: You can output massive, multi-component 'Paragraph' spells.",
		18: "Sound Barrier — You gain a flying speed equal to your walking speed, propelled by constant micro-blasts.",
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1.",
		20: "Manual Overdrive — Once per long rest, you ignore the timer for 1 minute.",
	}
}
