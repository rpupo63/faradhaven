package faradhaven_classes

// PistonBrawler returns the Piston Brawler class seed
func PistonBrawler() FaradhavenClassSeed {
	archetypeLevel := 3
	return FaradhavenClassSeed{
		Name:           "The Piston Brawler",
		Description:    "A martial engineer who fuses mechanical power with weaponry. You manage Stability Charges to fuel explosive strikes and tactical steam-based magic.",
		HitDie:         10,
		PrimaryAbility: "Strength",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/piston_brawler.jpg",
		Archetype:      "Tech Knight / Combat Engineer",
		Concept:        "A weapon-agnostic engineer who attaches a 'Piston Core' to their weapon, managing 'Stability Charges' to fuel magical strikes and utility spells.",
		DnDSkillFocus:    []string{"Arcana", "Athletics"},
		Proficiencies:    "Simple Weapons, Martial Weapons, Light Armor, Medium Armor, Shields",
		SkillChoice:      []string{"Arcana", "Athletics", "Investigation", "Sleight of Hand", "Intimidation", "History", "Perception"},
		Tools:            []string{"Smith's Tools", "Tinker's Tools"},
		SavingThrows:     []string{"Constitution", "Intelligence"},
		AutomaticEquipNames: []string{"Piston Core assembly kit"},
		AutomaticItemNames:  []string{"Scale mail", "Explorer's pack", "Tinker's tools"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your primary Piston-weapon",
				Options: []EquipmentOptionSeed{
					{Description: "Atlas Gauntlets (Heavy Piston Fists)", WeaponNames: []string{"Atlas Gauntlets"}},
					{Description: "A Greatsword (Heavy Blade)", WeaponNames: []string{"Greatsword"}},
					{Description: "A Warhammer (Versatile Crusher)", WeaponNames: []string{"Warhammer"}},
				},
			},
			{
				Instruction: "Choose your maintenance gear",
				Options: []EquipmentOptionSeed{
					{Description: "A Shield and Pressure Gauge", ItemNames: []string{"Standard Shield", "Stability Pressure Gauge"}},
					{Description: "Piston Core Assembly Kit", ItemNames: []string{"Piston Core Assembly Kit"}},
				},
			},
		},
		LevelFeatures:       pistonBrawlerLevelFeatures(),
		LevelProgression:    pistonBrawlerLevelProgression(),
		ArchetypeLevel:      &archetypeLevel,
		Archetypes:          pistonBrawlerArchetypes(),
		ResourceType:        "stability",
		ResourceName:        "Stability Charges",
		ResourceRestoreType: "long_rest",
		WeaponRequirement: &WeaponRequirementSeed{
			SelectionLevel:    1,
			ModifierType:      "piston_core",
			Description:       "Choose a weapon to integrate your Piston Core. This mechanical device converts kinetic energy into magical output.",
			AllowedCategories: []string{"Martial Melee", "Simple Melee", "Martial Ranged", "Simple Ranged"},
		},
	}
}

func pistonBrawlerLevelProgression() map[int]ClassLevelSeed {
	// Stability charges: Proficiency Bonus + Intelligence Modifier
	return map[int]ClassLevelSeed{
		1:  {MaxStability: 5, MaxSpellLevel: 1},
		2:  {MaxStability: 5, MaxSpellLevel: 1},
		3:  {MaxStability: 5, MaxSpellLevel: 1},
		4:  {MaxStability: 6, MaxSpellLevel: 1},
		5:  {MaxStability: 7, MaxSpellLevel: 2, ExtraAttackCount: 1},
		6:  {MaxStability: 7, MaxSpellLevel: 2, ExtraAttackCount: 1},
		7:  {MaxStability: 7, MaxSpellLevel: 2, ExtraAttackCount: 1},
		8:  {MaxStability: 7, MaxSpellLevel: 2, ExtraAttackCount: 1},
		9:  {MaxStability: 8, MaxSpellLevel: 3, ExtraAttackCount: 1},
		10: {MaxStability: 8, MaxSpellLevel: 3, ExtraAttackCount: 1},
		11: {MaxStability: 8, MaxSpellLevel: 3, ExtraAttackCount: 1},
		12: {MaxStability: 9, MaxSpellLevel: 3, ExtraAttackCount: 1},
		13: {MaxStability: 10, MaxSpellLevel: 4, ExtraAttackCount: 1},
		14: {MaxStability: 10, MaxSpellLevel: 4, ExtraAttackCount: 1},
		15: {MaxStability: 10, MaxSpellLevel: 4, ExtraAttackCount: 1},
		16: {MaxStability: 10, MaxSpellLevel: 4, ExtraAttackCount: 1},
		17: {MaxStability: 11, MaxSpellLevel: 5, ExtraAttackCount: 1},
		18: {MaxStability: 11, MaxSpellLevel: 5, ExtraAttackCount: 1},
		19: {MaxStability: 11, MaxSpellLevel: 5, ExtraAttackCount: 1},
		20: {MaxStability: 20, MaxSpellLevel: 5, ExtraAttackCount: 1},
	}
}

func pistonBrawlerLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "Piston Core", Description: "You have integrated a mechanical core into your chosen weapon. You can use your Intelligence modifier, instead of Dexterity, for the attack and damage rolls of your Piston weapon if it has the Finesse property."},
			{Name: "Stability Charges", Description: "Your core generates Stability Charges equal to your Proficiency Bonus + Intelligence modifier. These are used to fuel your abilities and spells. You regain all charges on a Long Rest. Save DC = 8 + Proficiency + Intelligence modifier."},
			{Name: "Integrated Spell Logic", Description: "Your weapon has 'Fixed Spells' programmed into its core. You can cast these spells by spending Stability Charges equal to the spell's level."},
		},
		2: {{Name: "Overdrive", Description: "When you hit with your Piston weapon, you can expend 1 Stability Charge to add 2d6 Force damage to the strike."}},
		3: {{Name: "Brawler Specialization", Description: "Choose your engineering philosophy: The Vulcan (Heat/Magic) or The Piston (Kinetic/Impact)."}},
		4: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},
		5: {{Name: "Extra Attack", Description: "You can attack twice when taking the Attack action. Your Max Spell Level increases to 2nd."}},
		6: {{Name: "Engineer's Eye (Ribbon)", Description: "You gain proficiency in the Investigation skill. If you already have it, you gain expertise. Additionally, you can cast 'Identify' as a ritual, but only on mechanical objects."}},
		7: {{Name: "Kinetic Calibration", Description: "When you have at least 1 Stability Charge remaining, you have Advantage on Strength and Dexterity saving throws."}},
		8: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},
		9: {{Name: "High-Pressure Output", Description: "Max Spell Level increases to 3rd (e.g., Haste, Fireball)."}},
		10: {{Name: "Reinforced Valves", Description: "When you use Overdrive, you can choose to knock the target back 5 feet in addition to the extra damage (Strength save against your Stability DC to resist)."}},
		11: {{Name: "Feedback Loop", Description: "When you cast a Fixed Spell of 1st level or higher, your next weapon attack deals an extra 1d8 force damage."}, {Name: "Improved Overdrive", Description: "Overdrive damage increases to 3d6."}},
		12: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},
		13: {{Name: "Advanced Calibration (Ribbon)", Description: "Max Spell Level increases to 4th. You can spend 1 minute during a short rest to recalibrate your core, allowing you to swap one of your Fixed Spells for another appropriate for your level."}},
		14: {{Name: "Pressurized Armor (Ribbon)", Description: "While wearing medium or heavy armor, you can add your Intelligence modifier (max 2) to your AC instead of Dexterity."}},
		15: {{Name: "Emergency Venting", Description: "When you take damage, you can use your reaction and 1 Stability Charge to reduce that damage by an amount equal to your Intelligence modifier + your Level. You also gain resistance to that instance of damage."}},
		16: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},
		17: {{Name: "Masterwork Engineering", Description: "Max Spell Level increases to 5th. Your Stability Charges now also regain on a Short Rest (once per Long Rest)."}},
		18: {{Name: "Internal Lubrication (Ribbon)", Description: "The constant vibrations of your core make you immune to the paralyzed and restrained conditions."}},
		19: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},
		20: {{Name: "The Perfect Machine", Description: "You regain 1 Stability Charge at the start of each of your turns as long as you have at least 1 charge remaining. If you start your turn with 0 charges, this does not function."}},
	}
}

func pistonBrawlerArchetypes() []ArchetypeSeed {
	return []ArchetypeSeed{
		{
			Name:        "The Vulcan",
			Description: "Specializes in heat management and magical output, venting energy as scorching fire.",
			Features: map[int][]FeatureSeed{
				3: {
					{Name: "Thermal Integration", Description: "You add your Intelligence modifier to damage rolls made with your Piston weapon."},
					{Name: "Heat Sink", Description: "When you spend one or more Stability Charges on your turn, your next successful hit before the end of your next turn deals an extra 1d6 fire damage."},
				},
				7: {
					{Name: "Cinder Cloud", Description: "As a bonus action, you can spend 1 Stability Charge to emit a cloud of steam and ash. The area within 10 feet of you becomes heavily obscured until the start of your next turn. You can see through this cloud."},
				},
				15: {
					{Name: "Supercritical State", Description: "When you use Overdrive, the extra damage increases to 4d6 (Fire or Force), and all creatures within 5 feet of the target take fire damage equal to your Intelligence modifier."},
				},
			},
		},
		{
			Name:        "The Piston",
			Description: "Focuses on kinetic force and physical impact, using hydraulic pressure to dominate the battlefield.",
			Features: map[int][]FeatureSeed{
				3: {
					{Name: "Kinetic Integration", Description: "You add your Intelligence modifier to damage rolls made with your Piston weapon."},
					{Name: "Forceful Impact", Description: "When you hit a creature with your Piston weapon, you can spend 1 Stability Charge to push the creature up to 10 feet away from you or knock it prone."},
				},
				7: {
					{Name: "Hydraulic Leap", Description: "You can spend 1 Stability Charge as a bonus action to leap up to 30 feet to an unoccupied space you can see. This movement does not provoke opportunity attacks."},
				},
				15: {
					{Name: "Wall Breaker", Description: "Your Piston weapon attacks deal double damage to objects and structures. Additionally, when you use Forceful Impact, you can push the creature up to 20 feet. If they hit a solid surface, they take 2d6 bludgeoning damage and are stunned until the end of your next turn."},
				},
			},
		},
	}
}