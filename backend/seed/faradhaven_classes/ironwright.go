package faradhaven_classes

// Ironwright returns the Ironwright class seed
func Ironwright() FaradhavenClassSeed {
	archetypeLevel := 3
	return FaradhavenClassSeed{
		Name:             "The Ironwright",
		Description:      "A resource-management class that scavenges parts from fallen enemies to build temporary constructs. Instead of spell slots, they manage an economy of Components found on the battlefield, scaling through Concurrency (active bots) and Yield (scavenging dice).",
		HitDie:           8,
		PrimaryAbility:   "intelligence",
		PhotoURL:         "https://photos-for-apps.s3.us-east-2.amazonaws.com/ironwright.jpg",
		Archetype:        "Resource Manager / Summoner",
		Concept:          "An engineer who realizes that machinery is just magic with better documentation. They rely on a looting-and-crafting economy, fueling their constructs with parts salvaged from the dead.",
		ClassFeatures:    []string{"Scavenge: Farm Components from fallen creatures instead of using spell slots.", "Construct Assembly: Build Sentries, Strikers, and Titans using harvested components.", "Concurrency: Manage a limit of active 'Live' constructs that scales with level."},
		DnDSkillFocus:    []string{"Investigation", "Arcana"},
		Proficiencies:    "Simple weapons, Light Armor, Hand Crossbows",
		SkillChoice:      []string{"Sleight of Hand", "History", "Medicine", "Investigation", "Survival"},
		Tools:            []string{"Tinker's Tools", "Smith's Tools"},
		SavingThrows:     []string{"Intelligence", "Constitution"},
		AutomaticEquipNames: []string{"Leather apron (light armor)", "Protective goggles", "Bag of gears and springs"},
		AutomaticItemNames:  []string{"Tinker's Tools"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your engineering tool",
				Options: []EquipmentOptionSeed{
					{Description: "A Heavy Wrench (Warhammer)", WeaponNames: []string{"Warhammer"}},
					{Description: "Mechanized Oil Can and Hammer", ItemNames: []string{"Mechanized Oil Can"}, WeaponNames: []string{"Light Hammer"}},
				},
			},
			{
				Instruction: "Choose your salvage gear",
				Options: []EquipmentOptionSeed{
					{Description: "Bag of Gears and Springs", ItemNames: []string{"Bag of Gears and Springs"}},
					{Description: "Smelter's Gloves", ItemNames: []string{"Smelter's Gloves"}},
				},
			},
		},
		LevelFeatures:       ironwrightLevelFeatures(),
		LevelProgression:    ironwrightLevelProgression(),
		ArchetypeLevel:      &archetypeLevel,
		Archetypes:          ironwrightArchetypes(),
		ResourceType:        "components",
		ResourceName:        "Components",
		ResourceRestoreType: "special", // gathered from fallen creatures, not restored on rest
	}
}

func ironwrightLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Scavenge — You harvest Components from fallen creatures (1d4 yield). Construct (Sentry): Build a stationary ranged turret (Cost: 2 Components). Max Active Constructs: 1.",
		2:  "Efficiency (Reclaim Parts) — Use a reaction to recover 50% of components when a construct dies or spell ends.",
		3:  "Ironwright Specialization — Choose your doctrine of engineering. Your specialization grants you unique abilities at levels 3, 7, and 10.",
		4:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1.",
		5:  "Extra Attack, Striker Construct — You can attack twice when you take the Attack action. Construct (Striker): Build a melee combatant construct (Cost: 5 Components). Max Active Constructs: 2. Scavenge Yield: 1d6.",
		6:  "Specialized Salvage — Scavenging CR 5+ enemies yields 'Primal Components' for elemental resistances/damage.",
		8:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1.",
		9:  "Construct (Titan) — Build a Large mountable construct with high AC (Cost: 10 Components). Max Active Constructs: 3. Scavenge Yield: 1d8.",
		11: "Rapid Assembly — Spend components to cast spells or summon constructs as a Bonus Action.",
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1.",
		13: "Overclock — Temporarily boost construct performance. Max Active Constructs: 4. Scavenge Yield: 1d10.",
		14: "Master of Scraps — Level 1 & 2 spells cost 0 components if you have at least 1. Only 'Heavy' bots cost resources now.",
		15: "Automated Logistics — Long Rest restores baseline components (5 x Level). Daily refresh floor.",
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1.",
		17: "Swarm Protocol — Coordinate larger numbers of machines. Max Active Constructs: 5. Scavenge Yield: 1d12.",
		18: "Eternal Engine — Your constructs become self-sustaining engines of war.",
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1.",
		20: "Grand Architect (The Foundry) — Max Active Constructs: 6. Infinite Loop: Constructs auto-scavenge kills instantly. Scavenge Yield: 2d12.",
	}
}

func ironwrightLevelProgression() map[int]ClassLevelSeed {
	return map[int]ClassLevelSeed{
		1:  {ConcurrencyLimit: 1, YieldDie: 4},
		2:  {ConcurrencyLimit: 1, YieldDie: 4},
		3:  {ConcurrencyLimit: 1, YieldDie: 4},
		4:  {ConcurrencyLimit: 1, YieldDie: 4},
		5:  {ExtraAttackCount: 1, ConcurrencyLimit: 2, YieldDie: 6},
		6:  {ExtraAttackCount: 1, ConcurrencyLimit: 2, YieldDie: 6},
		7:  {ExtraAttackCount: 1, ConcurrencyLimit: 2, YieldDie: 6},
		8:  {ExtraAttackCount: 1, ConcurrencyLimit: 2, YieldDie: 6},
		9:  {ExtraAttackCount: 1, ConcurrencyLimit: 3, YieldDie: 8},
		10: {ExtraAttackCount: 1, ConcurrencyLimit: 3, YieldDie: 8},
		11: {ExtraAttackCount: 1, ConcurrencyLimit: 3, YieldDie: 8},
		12: {ExtraAttackCount: 1, ConcurrencyLimit: 3, YieldDie: 8},
		13: {ExtraAttackCount: 1, ConcurrencyLimit: 4, YieldDie: 10},
		14: {ExtraAttackCount: 1, ConcurrencyLimit: 4, YieldDie: 10},
		15: {ExtraAttackCount: 1, ConcurrencyLimit: 4, YieldDie: 10},
		16: {ExtraAttackCount: 1, ConcurrencyLimit: 4, YieldDie: 10},
		17: {ExtraAttackCount: 1, ConcurrencyLimit: 5, YieldDie: 12},
		18: {ExtraAttackCount: 1, ConcurrencyLimit: 5, YieldDie: 12},
		19: {ExtraAttackCount: 1, ConcurrencyLimit: 5, YieldDie: 12},
		20: {ExtraAttackCount: 1, ConcurrencyLimit: 6, YieldDie: 12}, // 2d12 at 20
	}
}

func ironwrightArchetypes() []ArchetypeSeed {
	return []ArchetypeSeed{
		{
			Name:        "The Automaton Forger",
			Description: "You specialize in building heavily armored, durable constructs. Your machines are built to last and can withstand tremendous punishment.",
			Features: map[int]string{
				3:  "Reinforced Chassis — Your constructs gain bonus hit points equal to your Intelligence modifier + your level. They have advantage on saving throws against effects that would destroy or disable them.",
				7:  "Adaptive Plating — Your constructs gain resistance to one damage type of your choice when you create them. You can change this resistance when you finish a long rest.",
				10: "Fortress Protocol — As a reaction when a construct within 30 feet of you would take damage, you can spend 2 components to give it resistance to that damage. Additionally, your constructs can use the Dodge action as a bonus action.",
			},
		},
		{
			Name:        "The Swarm Engineer",
			Description: "You specialize in creating many small, expendable constructs that overwhelm enemies through numbers and coordinated attacks.",
			Features: map[int]string{
				3:  "Micro-Constructs — You can create Drone constructs (Cost: 1 Component) that have 1 HP but deal 1d4 damage. These do not count against your construct limit. You can have up to your Intelligence modifier active at once.",
				7:  "Swarm Tactics — When two or more of your constructs attack the same target in a round, each attack after the first deals an additional 1d6 damage.",
				10: "Distributed Processing — Your constructs share a hive mind. When one construct sees an enemy, all your constructs can target that enemy without line of sight for 1 round. Once per short rest, when a construct dies, you can immediately create a Drone as a free action.",
			},
		},
		{
			Name:        "The Recycler",
			Description: "You specialize in efficiency and resource recovery. Nothing goes to waste in your workshop, and you can squeeze maximum value from every component.",
			Features: map[int]string{
				3:  "Salvage Expertise — When you use Scavenge, add your proficiency bonus to the yield. When you recover components from your own constructs, you recover 75% instead of 50%.",
				7:  "Modular Design — You can spend 1 component to swap the type of an existing construct (Sentry to Striker, etc.) as an action. The construct retains its current HP.",
				10: "Perpetual Motion — Once per long rest, you can designate one construct as 'Self-Sustaining.' It costs 0 components to maintain and automatically repairs 1d6 HP at the start of each of its turns. When it dies, you recover 100% of its component cost.",
			},
		},
	}
}
