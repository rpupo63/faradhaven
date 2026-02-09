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

func powderMageLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "Flash-Point Casting", Description: "You cast spells in real-time. When you take the Cast action, a timer starts (base 2 seconds). You must speak the names of spell components aloud. Each offensive component adds 1d8 damage of its type. Base range: 30ft. Save DC: 8 + Prof + Dex."},
			{Name: "Kinetic Recoil", Description: "When you finish casting a spell, you can move 5 feet immediately without using an action or provoking opportunity attacks."},
			{Name: "Speed Dial", Description: "You can save specific component strings to cast instantly (Unlocks at Level 3)."},
		},
		2:  {{Name: "Powder Sprint & Disengage", Description: "Your walking speed increases by 10 feet. You can use Disengage or Dash as a bonus action."}},
		3:  {{Name: "Speed Dial (The 'Pre-Mix')", Description: "Save one string of up to 3 components. Once per short rest, you can cast this string as an action without using the timer."}},
		5:  {{Name: "Kinetic Recoil (Traversal)", Description: "When you finish a spell, you can move 10 feet in any direction (including upwards) without provoking opportunity attacks. Your Casting Timer increases to 3 seconds."}},
		6:  {{Name: "Recoil Mastery", Description: "You can move along vertical surfaces and across liquids during your Kinetic Recoil movement without falling."}},
		7:  {{Name: "Evasion", Description: "When you make a Dexterity save to take only half damage, you instead take no damage on a success, and only half damage on a failure."}},
		9:  {{Name: "Speed Dial Upgrade", Description: "You can now save 2 different component strings of up to 4 components each."}},
		11: {{Name: "Room Clearer", Description: "As an action, move up to your full speed. Each creature you pass within 5 feet can be targeted by a single-component spell (e.g., 'Ignis') if you speak the word as you move past. Your Casting Timer increases to 4 seconds."}},
		13: {{Name: "Sonic Boom", Description: "When you use Kinetic Recoil, you can deal 3d6 thunder damage to all creatures within 5 feet of your starting space (Dexterity save for half)."}},
		14: {{Name: "Echoing Components", Description: "If you use the same component in two consecutive turns, it is automatically active on the second turn without needing to be spoken, freeing up time on your timer."}},
		15: {{Name: "Speed Dial Upgrade", Description: "You can now save 3 different component strings of up to 5 components each."}},
		17: {{Name: "Timer Upgrade (5 Seconds)", Description: "The Avalanche: Your Casting Timer increases to 5 seconds, allowing for complex, multi-component 'Paragraph' spells."}},
		18: {{Name: "Sound Barrier", Description: "You gain a flying speed equal to your walking speed, and you are always under the effects of the 'Silence' spell (self only) while moving."}},
		20: {{Name: "Manual Overdrive", Description: "Once per long rest, you can ignore the Casting Timer for 1 minute. You can speak as many components as you can for each spell cast during this duration."}},
	}
}
