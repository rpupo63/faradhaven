package faradhaven_classes

// Ironwright returns the Ironwright class seed
func Ironwright() FaradhavenClassSeed {
	archetypeLevel := 3
	return FaradhavenClassSeed{
		Name:             "The Ironwright",
		Description:      "A resource-management class that scavenges distinct Components from fallen enemies to build independent constructs and fuel their abilities. Instead of spell slots, they manage an inventory of specific Components that function as 'Scrap' for their schematics. Their power scales through Concurrency (total available slots for constructs) and Yield (scavenging dice).",
		HitDie:           8,
		PrimaryAbility:   "intelligence",
		PhotoURL:         "https://photos-for-apps.s3.us-east-2.amazonaws.com/ironwright.jpg",
		Archetype:        "Resource Manager / Summoner",
		Concept:          "An engineer who realizes that machinery is just magic with better documentation. They rely on a looting-and-crafting economy, fueling their independent constructs with parts salvaged from the dead.",
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
		ComponentPool: []string{
			// Forma (Shape)
			"Projectile", "Beam", "Nova", "Wall", "Zone", "Cone", "Aura",
			// Scopus (Targeting)
			"Target", "Self", "Ground",
			// Essentia (Matter & Energy)
			"Ignis", "Terra", "Ferrum", "Fulgur", "Arcanum", "Spatium", "Sonus",
			// Actio (Verbs)
			"Push", "Pull", "Grab", "Lift", "Spin", "Crush", "Pierce", "Create", "Destroy", "Mutate", "Bind",
			// Magnitudo (Modifiers)
			"Increase", "Decrease", "Strong", "Weak", "Extreme",
		},
		ResourceDefinitions: []ResourceDefinitionSeed{
			{Key: "concurrency_limit", DisplayName: "Concurrency Limit", Category: "limit", Description: "Maximum number of active constructs", DisplayOrder: 1},
			{Key: "yield_die", DisplayName: "Scavenge Yield Die", Category: "die_size", Description: "Die size rolled when scavenging components from fallen creatures", DisplayOrder: 2},
		},
	}
}

func ironwrightLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "Scavenge", Description: "As an action, you harvest 1d4 distinct Components (e.g., 'Cog', 'Gear', 'Ignis') from a fallen creature (died within 1 hour, size Small or larger). This consumes the corpse. These Components are used as Scrap for your abilities."},
			{Name: "Construct Assembly", Description: "As an action, spend Components to build a construct in an unoccupied space within 10 feet. Constructs act immediately after your turn and last for (Component Cost × 10) minutes. Construct (Sentry): Cost: 2 Components (any type). Slot Cost: 1. AC: 13, HP: 10 + (2 * Level), Atk: + (Proficiency + INT), Dmg: 1d8 + INT force, Range: 60ft. Duration: 20 minutes."},
			{Name: "Concurrency", Description: "You can maintain a number of 'Live' constructs simultaneously. Max Concurrency Slots: 1."},
		},
		2: {
			{Name: "Efficiency (Reclaim Parts)", Description: "As a reaction when one of your constructs is reduced to 0 HP or its duration ends, you recover 50% of the Components used to build it (rounded down, minimum 1)."},
		},
		3: {
			{Name: "Ironwright Specialization", Description: "Choose your doctrine of engineering: The Automaton Forger, The Swarm Engineer, or The Recycler. Your specialization grants you unique abilities at levels 3, 7, and 10."},
		},
		4: {
			{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."},
		},
		5: {
			{Name: "Extra Attack, Striker Construct", Description: "You can attack twice when you take the Attack action. Construct (Striker): A melee combatant construct. Cost: 5 Components (any type). Slot Cost: 2. AC: 15, HP: 15 + (4 * Level), Atk: + (Proficiency + INT), Dmg: 1d10 + INT force. Duration: 50 minutes. Max Concurrency Slots: 2. Scavenge Yield: 1d6."},
		},
		6: {
			{Name: "Specialized Salvage", Description: "Scavenging CR 5+ enemies yields 1 'Primal Component' (Fire, Cold, Lightning, or Acid). As an action, you can expend this item to grant one of your constructs resistance to that damage type for 1 hour or add 1d6 of that damage type to its attacks for 1 minute."},
		},
		7: {
			{Name: "Mechanical Intuition", Description: "Your deep understanding of machinery grants you advantage on Investigation and Arcana checks related to machinery, traps, and constructs. When you spend 1 minute examining a mechanical device, you learn its purpose and how to operate or disable it."},
		},
		8: {
			{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."},
		},
		9: {
			{Name: "Construct (Titan)", Description: "Build a Large mountable construct. Cost: 10 Components (any type). Slot Cost: 3. AC: 18, HP: 30 + (6 * Level), Atk: + (Proficiency + INT), Dmg: 2d8 + INT force. Duration: 100 minutes. Max Concurrency Slots: 3. Scavenge Yield: 1d8."},
		},
		10: {
			{Name: "Fortified Assembly", Description: "Your constructs gain +1 AC and +5 max HP. Additionally, when one of your constructs is destroyed, you automatically recover 1 bonus Component. This is in addition to the Components recovered via your Efficiency feature."},
		},
		11: {
			{Name: "Rapid Assembly", Description: "You can summon a construct as a bonus action instead of an action. Additionally, as a bonus action, you can spend Components to execute Instant Schematics: 'Construct Triage' (Cost: 1 Mend Component — command a construct to repair an adjacent construct for 2d8 + INT HP), 'Kinetic Barrier' (Cost: 1 Push and 1 Zone Component — grant yourself +5 AC as a reaction until the start of your next turn), or 'Emergency Reboot' (Cost: 1 Fulgur, 1 Conduct, and 1 Mend Component, once per long rest — reconstruct all downed constructs to full HP). This feature also includes 'Salvage Protocol' (passive — your Efficiency feature now recovers 75% of a construct's Component cost, rounded down)."},
		},
		12: {
			{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."},
		},
		13: {
			{Name: "Overclock", Description: "As a bonus action, spend 5 Components (any type) to Overclock a construct for 1 minute. It gains +2 to AC, advantage on attack rolls, and deals an additional 1d8 force damage on hit. After 1 minute, the construct takes 3d8 force damage. Max Concurrency Slots: 4. Scavenge Yield: 1d10."},
		},
		14: {
			{Name: "Stockpile", Description: "You always keep a reserve of emergency parts. When you roll initiative, you instantly gain 5 temporary Components. These are consumed first and disappear at the end of combat."},
		},
		15: {
			{Name: "Automated Logistics", Description: "Finishing a long rest restores a baseline of Components equal to your Ironwright level. This is your daily resource floor."},
		},
		16: {
			{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."},
		},
		17: {
			{Name: "Swarm Protocol", Description: "Coordinate larger numbers of machines. Max Concurrency Slots: 4. Scavenge Yield: 1d12."},
		},
		18: {
			{Name: "Networked Repair", Description: "Constructs can use their action to repair themselves or an adjacent construct by spending 1 Component from your inventory. This repair restores 4d8 + Intelligence modifier HP."},
		},
		19: {
			{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."},
		},
		20: {
			{Name: "Grand Architect (The Foundry)", Description: "Max Concurrency Slots: 5. Infinite Loop: Constructs auto-scavenge kills they make, adding 1d4 Components to your pool instantly. Scavenge Yield: 1d12."},
		},
	}
}

func ironwrightLevelProgression() map[int]ClassLevelSeed {
	return map[int]ClassLevelSeed{
		1:  {Resources: map[string]int{"concurrency_limit": 1, "yield_die": 4}},
		2:  {Resources: map[string]int{"concurrency_limit": 1, "yield_die": 4}},
		3:  {Resources: map[string]int{"concurrency_limit": 1, "yield_die": 4}},
		4:  {Resources: map[string]int{"concurrency_limit": 1, "yield_die": 4}},
		5:  {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 2, "yield_die": 6}},
		6:  {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 2, "yield_die": 6}},
		7:  {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 2, "yield_die": 6}},
		8:  {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 2, "yield_die": 6}},
		9:  {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 3, "yield_die": 8}},
		10: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 3, "yield_die": 8}},
		11: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 3, "yield_die": 8}},
		12: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 3, "yield_die": 8}},
		13: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 4, "yield_die": 10}},
		14: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 4, "yield_die": 10}},
		15: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 4, "yield_die": 10}},
		16: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 4, "yield_die": 10}},
		17: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 4, "yield_die": 12}},
		18: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 4, "yield_die": 12}},
		19: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 4, "yield_die": 12}},
		20: {ExtraAttackCount: 1, Resources: map[string]int{"concurrency_limit": 5, "yield_die": 12}},
	}
}

func ironwrightArchetypes() []ArchetypeSeed {
	return []ArchetypeSeed{
		{
			Name:        "The Automaton Forger",
			Description: "You specialize in building heavily armored, durable constructs. Your machines are built to last and can withstand tremendous punishment.",
			Features: map[int][]FeatureSeed{
				3:  {{Name: "Reinforced Chassis", Description: "Your constructs gain bonus HP equal to your Intelligence modifier + your Ironwright level. They have advantage on saving throws against effects that would destroy or disable them."}},
				7:  {{Name: "Adaptive Plating", Description: "When you create a construct, it gains resistance to one damage type of your choice. You can change this resistance by spending 1 minute and 1 Component recalibrating the construct."}},
				10: {{Name: "Fortress Protocol", Description: "As a reaction when a construct within 30 feet takes damage, you can spend 2 Components to give it resistance to all damage from that source until the start of your next turn. Additionally, your constructs can use the Dodge action as a bonus action."}},
			},
		},
		{
			Name:        "The Swarm Engineer",
			Description: "You specialize in creating many small, expendable constructs that overwhelm enemies through numbers and coordinated attacks.",
			Features: map[int][]FeatureSeed{
				3:  {{Name: "Micro-Constructs", Description: "As an action, you can spend 1 Component to create a Drone construct. Drones have 1 HP, AC 12, and deal 1d4 force damage on a hit. These do not count against your Concurrency Slots. Your maximum number of active Drones is equal to your Intelligence modifier + your proficiency bonus."}},
				7:  {{Name: "Swarm Tactics", Description: "When two or more of your constructs are within 5 feet of the same target, each of their attacks deals an additional 1d6 force damage."}},
				10: {{Name: "Distributed Processing", Description: "Your constructs share a hive mind. If one construct sees an enemy, all your constructs have advantage on attacks against that enemy for 1 round. Once per short rest, when a construct is destroyed, you can immediately create a Drone as a free action."}},
			},
		},
		{
			Name:        "The Recycler",
			Description: "You specialize in efficiency and resource recovery. Nothing goes to waste in your workshop, and you can squeeze maximum value from every component.",
			Features: map[int][]FeatureSeed{
				3:  {{Name: "Salvage Expertise", Description: "When you use Scavenge, you add your proficiency bonus to the Component yield. Additionally, when you use your Efficiency feature, you recover 1 additional Component."}},
				7:  {{Name: "Modular Design", Description: "As an action, you can spend 1 Component to transform an existing construct into another type you can build (e.g., Sentry to Striker). The construct retains its current percentage of HP."}},
				10: {{Name: "Perpetual Motion", Description: "Once per long rest, you can designate one construct as 'Self-Sustaining' for 1 hour. It costs 0 Components to maintain and automatically repairs 1d6 HP at the start of each of its turns. If it is destroyed, you recover 100% of its Component cost."}},
			},
		},
	}
}
