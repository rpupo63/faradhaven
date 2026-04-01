package faradhaven_classes

var three = 3

// PowderMage returns the Powder Mage class seed
func PowderMage() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Powder Mage",
		Description:    "The 'Continuous Ignition' Powder Mage. Start inputting components as soon as you declare the Cast action; whatever you chain together before the timer runs out is what happens. You aren't fighting to make the spell work, but to pack as much power into the casting window as possible.",
		HitDie:         8,
		PrimaryAbility: "dexterity",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/powder_mage.jpg",
		Archetype:      "Speed Caster / Skirmisher",
		Concept:        "Continuous Ignition: Real-time rapid casting under pressure.",
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
		LevelFeatures:    powderMageLevelFeatures(),
		LevelProgression: powderMageLevelProgression(),
		ArchetypeLevel:   &three,
		Archetypes: []ArchetypeSeed{
			{
				Name:        "The Demolitionist",
				Description: "You specialize in wide-area destruction and powerful, explosive effects.",
				Features: map[int][]FeatureSeed{
					3: {{Name: "Volatile Mixture", Description: "At 3rd level, when you cast a spell using 3 or more components, you can choose to apply a delayed explosive charge. The spell's area of effect becomes an unstable zone for 1 round. Any creature entering or starting its turn in the zone takes 1d6 fire damage. This can be used once per short rest."}},
				},
			},
			{
				Name:        "The Duelist",
				Description: "You focus on precise, single-target application of your magical munitions, using recoil to outmaneuver foes.",
				Features: map[int][]FeatureSeed{
					3: {{Name: "Precision Recoil", Description: "At 3rd level, when you use Kinetic Recoil, you can choose to move an enemy within 5 feet of you 5 feet in a direction of your choice (Dexterity saving throw to resist). This movement does not provoke opportunity attacks. You can use this ability a number of times equal to your Proficiency Bonus per long rest."}},
				},
			},
		},
		ComponentPool: []string{
			// Forma (Shape)
			"Projectile", "Nova", "Cone", "Beam", "Zone",
			// Scopus (Targeting)
			"Target", "Ground", "Chain",
			// Essentia (Matter & Energy)
			"Ignis", "Fulgur", "Aer", "Sonus", "Arcanum", "Chronos",
			// Actio (Verbs)
			"Push", "Pierce", "Spin",
			// Magnitudo (Modifiers)
			"Increase", "Strong", "Extreme",
			// Logica (sequential links)
			"If", "Then", "Therefore",
		},
		ResourceDefinitions: []ResourceDefinitionSeed{
			{Key: "timer_duration", DisplayName: "Casting Timer", Category: "modifier", Description: "Seconds available for real-time component input", DisplayOrder: 1},
			{Key: "speed_dial_slots", DisplayName: "Speed Dial Slots", Category: "slot_count", Description: "Pre-saved component strings for instant casting. Click a slot to mark it used; restores on rest.", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 2},
			{Key: "cascade_bonus", DisplayName: "Cascade Bonus", Category: "pool", Description: "Stacking +1d6 damage bonus from multi-component spells (max 2 stacks). Resets on next cast.", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 3},
		},
	}
}

func powderMageLevelProgression() map[int]ClassLevelSeed {
	// Powder Mage gains movement speed bonuses and timer/speed dial upgrades
	return map[int]ClassLevelSeed{
		1:  {Resources: map[string]int{"timer_duration": 2, "speed_dial_slots": 0, "cascade_bonus": 0}},
		2:  {UnarmoredMovement: 10, Resources: map[string]int{"timer_duration": 2, "speed_dial_slots": 0, "cascade_bonus": 0}},
		3:  {UnarmoredMovement: 10, Resources: map[string]int{"timer_duration": 2, "speed_dial_slots": 1, "cascade_bonus": 0}},
		4:  {UnarmoredMovement: 10, Resources: map[string]int{"timer_duration": 2, "speed_dial_slots": 1, "cascade_bonus": 0}},
		5:  {UnarmoredMovement: 10, Resources: map[string]int{"timer_duration": 3, "speed_dial_slots": 1, "cascade_bonus": 0}},
		6:  {UnarmoredMovement: 15, Resources: map[string]int{"timer_duration": 3, "speed_dial_slots": 1, "cascade_bonus": 0}},
		7:  {UnarmoredMovement: 15, Resources: map[string]int{"timer_duration": 3, "speed_dial_slots": 1, "cascade_bonus": 0}},
		8:  {UnarmoredMovement: 15, Resources: map[string]int{"timer_duration": 3, "speed_dial_slots": 1, "cascade_bonus": 0}},
		9:  {UnarmoredMovement: 15, Resources: map[string]int{"timer_duration": 3, "speed_dial_slots": 2, "cascade_bonus": 0}},
		10: {UnarmoredMovement: 15, Resources: map[string]int{"timer_duration": 3, "speed_dial_slots": 2, "cascade_bonus": 2}},
		11: {UnarmoredMovement: 20, Resources: map[string]int{"timer_duration": 4, "speed_dial_slots": 2, "cascade_bonus": 2}},
		12: {UnarmoredMovement: 20, Resources: map[string]int{"timer_duration": 4, "speed_dial_slots": 2, "cascade_bonus": 2}},
		13: {UnarmoredMovement: 20, Resources: map[string]int{"timer_duration": 4, "speed_dial_slots": 2, "cascade_bonus": 2}},
		14: {UnarmoredMovement: 20, Resources: map[string]int{"timer_duration": 4, "speed_dial_slots": 2, "cascade_bonus": 2}},
		15: {UnarmoredMovement: 20, Resources: map[string]int{"timer_duration": 4, "speed_dial_slots": 3, "cascade_bonus": 2}},
		16: {UnarmoredMovement: 20, Resources: map[string]int{"timer_duration": 4, "speed_dial_slots": 3, "cascade_bonus": 2}},
		17: {UnarmoredMovement: 20, Resources: map[string]int{"timer_duration": 5, "speed_dial_slots": 3, "cascade_bonus": 2}},
		18: {UnarmoredMovement: 30, Resources: map[string]int{"timer_duration": 5, "speed_dial_slots": 3, "cascade_bonus": 2}},
		19: {UnarmoredMovement: 30, Resources: map[string]int{"timer_duration": 5, "speed_dial_slots": 3, "cascade_bonus": 2}},
		20: {UnarmoredMovement: 30, Resources: map[string]int{"timer_duration": 5, "speed_dial_slots": 3, "cascade_bonus": 2}},
	}
}

func powderMageLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "Flash-Point Casting", Description: "You cast spells in real-time. When you take the Cast action, a timer starts (base 2 seconds). You must input spell components rapidly. Components combine into a single damage type determined by the DM based on the combination. Each damage-dealing component adds 1d8 damage. Base range: 30ft. Save DC: 8 + Prof + Dex."},
			{Name: "Kinetic Recoil", Description: "When you finish casting a spell, you are automatically displaced 5 feet in a direction of your choice. This displacement happens before reactions can be taken against you and does not consume your movement."},
			{Name: "Residual Heat", Description: "Outside of combat, full timer-casts require a 1-minute cooldown between casts. Casting again before the cooldown expires deals 1d6 fire damage per component used to you. Speed Dial cantrips are unaffected and remain free to cast at any time."},
		},
		2:  {{Name: "Powder Sprint & Disengage", Description: "Your walking speed increases by 10 feet. You can use Disengage or Dash as a bonus action."}},
		3:  {{Name: "Speed Dial (The 'Pre-Mix')", Description: "Save one string of up to 3 components. You can cast this string as an action without using the timer. Speed Dial casts function as cantrips — they are free to use and can be cast an unlimited number of times."}},
		4:  {{Name: "Powder Sense", Description: "You have an instinctive awareness of volatile substances. You gain advantage on Investigation and Perception checks to detect explosives, alchemical substances, and traps within 30 feet of you."}},
		5:  {{Name: "Kinetic Recoil (Traversal)", Description: "When you finish a spell, you can move 10 feet in any direction (including upwards) without provoking opportunity attacks. Your Casting Timer increases to 3 seconds."}},
		6:  {{Name: "Recoil Mastery", Description: "You can move along vertical surfaces and across liquids during your Kinetic Recoil movement without falling."}, {Name: "Chain Reaction", Description: "When you cast a spell that includes at least two components of the same element, you can add your Dexterity modifier to the damage of that spell. If you include four or more components of the same element, you add double your Dexterity modifier to the damage instead."}},
		7:  {{Name: "Evasion", Description: "When you make a Dexterity save to take only half damage, you instead take no damage on a success, and only half damage on a failure."}},
		8:  {{Name: "Quick Reload", Description: "You can reload firearms and interact with components as a free action (no action or object interaction required). You can use Dexterity instead of Strength for firearm attack rolls."}},
		9:  {{Name: "Speed Dial Upgrade", Description: "You can now save 2 different component strings."}},
		10: {{Name: "Cascade Priming", Description: "When you cast a spell containing 3 or more components of the same element, your next spell that includes that element deals an additional +1d6 damage of that element's type. This bonus can stack up to 2 times (max +2d6) and resets after the next spell is cast or after 1 minute."}},
		11: {{Name: "Room Clearer", Description: "As an action, spend your full movement to charge through the battlefield. You move up to your full speed; each enemy you pass within 5 feet of is hit by a Speed Dial or quick single-component spell (input the component as you pass). Affected enemies cannot take opportunity attacks against you during this movement. This costs both your action and your movement for the turn."}},
		12: {{Name: "Controlled Detonation", Description: "Your expertise with volatile materials extends beyond combat. You gain Expertise in Investigation (double proficiency bonus), proficiency with Thieves' Tools, and advantage on checks to disable traps. When you successfully disarm a trap or explosive device, you can salvage components from it (DM discretion on type and quantity)."}},
		13: {{Name: "Sonic Boom", Description: "When you use Kinetic Recoil, you can deal 3d6 thunder damage to all creatures within 5 feet of your starting space (Dexterity save for half)."}},
		14: {{Name: "Component Mastery", Description: "You can inject a Speed Dial string into an active timer cast as a hotkey. When casting with the timer, you may use a bonus action to instantly add the components of one of your saved Speed Dial strings to the current spell without consuming the Speed Dial slot."}},
		15: {{Name: "Speed Dial Upgrade", Description: "You can now save 3 different component strings."}},
		16: {{Name: "Flash Performance", Description: "You can use your components to create dazzling non-combat spectacles — fireworks, light shows, smoke displays, and other feats of alchemical showmanship. You gain advantage on Performance checks when using components as part of the performance."}},
		17: {{Name: "Timer Upgrade (5 Seconds)", Description: "The Avalanche: Your Casting Timer increases to 5 seconds, allowing for complex, multi-component 'Paragraph' spells."}},
		18: {{Name: "Sound Barrier", Description: "You gain a flying speed equal to your walking speed, and you generate a field of Acoustic Dampening, rendering you and your immediate surroundings silent while moving."}},
		19: {{Name: "Combat Reload", Description: "At the start of each combat (when you roll initiative), you regain 1 expended Speed Dial slot. Additionally, you can use a bonus action to refresh all Speed Dial slots. You can use this bonus action refresh once per short rest."}},
		20: {{Name: "Manual Overdrive", Description: "Once per long rest, you can ignore the Casting Timer for 1 minute. You can input as many components as you can for each spell cast during this duration."}},
	}
}
