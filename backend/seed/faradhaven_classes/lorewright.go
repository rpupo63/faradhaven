package faradhaven_classes

// Lorewright returns the Lorewright class seed.
// A class that harvests creatures to absorb their components, skills, attacks, and recipes,
// balancing power against the risk of madness. High wisdom and constitution,
// with harvested abilities guiding their actions.
func Lorewright() FaradhavenClassSeed {
	archetypeLevel := 3
	return FaradhavenClassSeed{
		Name:           "The Lorewright",
		Description:    "Harvest the flesh of your prey to absorb their components, skills, and abilities. Balance the power you gain against the fracturing of your mind as you become the sum of all you have consumed.",
		HitDie:         8,
		PrimaryAbility: "wisdom",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/lorewright.jpg",
		Archetype:      "Knowledge / Survival / Transformation",
		Concept:        "A hunter-scholar who has learned that true understanding comes through consumption. You eat the livers of your prey to absorb their memories, skills, and eventually their physical traits. Max WIS and CON—your mind must be strong enough to hold the souls of beasts, and your body must process what you consume. Every meal is a gamble between power and madness.",
		DnDSkillFocus:    []string{"Insight", "Survival"},
		Proficiencies:    "Simple weapons, Martial weapons, Light armor, Medium armor, Shields",
		SkillChoice:      []string{"History", "Insight", "Medicine", "Nature", "Survival", "Animal Handling", "Perception"},
		Tools:            []string{"Herbalism Kit", "Cook's Utensils"},
		SavingThrows:     []string{"Wisdom", "Intelligence"},
		AutomaticEquipNames: []string{"Preservation kit (salt, jars)", "Traveler's clothes", "Bone talisman"},
		AutomaticItemNames:  []string{"Leather armor", "Explorer's pack"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your harvesting tool",
				Options: []EquipmentOptionSeed{
					{Description: "Specimen Preservation Jar and Knife", ItemNames: []string{"Specimen Preservation Jar"}, WeaponNames: []string{"Dagger"}},
					{Description: "Anatomist's Kit", ItemNames: []string{"Anatomist's Kit"}},
				},
			},
			{
				Instruction: "Choose your weapon",
				Options: []EquipmentOptionSeed{
					{Description: "A Shortbow and 20 arrows", Items: []string{"Shortbow", "20 Arrows"}, WeaponNames: []string{"Shortbow"}},
					{Description: "A Spear", WeaponNames: []string{"Spear"}},
				},
			},
		},
		LevelFeatures:       lorewrightLevelFeatures(),
		LevelProgression:    lorewrightLevelProgression(),
		ArchetypeLevel:      &archetypeLevel,
		Archetypes:          lorewrightArchetypes(),
		ComponentPool: []string{
			// Forma (Shape)
			"Aura", "Self",
			// Scopus (Targeting)
			"Target", "Self",
			// Essentia (Matter & Energy)
			"Vita", "Mortis", "Odor", "Arcanum",
			// Essentia (Pathos)
			"Anger", "Sadness", "Fear", "Happiness", "Disgust", "Embarrassment",
			// Actio (Verbs)
			"Mutate", "Bind", "Crush", "Pierce", "Destroy",
			// Magnitudo (Modifiers)
			"Increase", "Decrease", "Strong", "Weak",
		},
		ResourceDefinitions: []ResourceDefinitionSeed{
			{Key: "echo_slots", DisplayName: "Harvest Slots", Category: "slot_count", Description: "Total harvest slots for storing absorbed abilities", DisplayOrder: 1},
		},
	}
}

func lorewrightArchetypes() []ArchetypeSeed {
	return []ArchetypeSeed{
		{
			Name:        "Path of the Warlord",
			Description: "You channel the combat memories of fallen warriors, gaining martial prowess and the ability to strike with the fury of a hundred battles.",
			Features: map[int][]FeatureSeed{
				3:  {{Name: "Martial Legacy", Description: "You gain proficiency with heavy armor and one martial weapon. Your melee weapon attacks deal an additional 1d4 damage of your weapon's type."}},
				10: {{Name: "Predator's Strike", Description: "When you hit a creature type you have harvested within the last 24 hours, add your Wisdom modifier to the damage roll."}},
				14: {{Name: "Battle Trance", Description: "When you enter combat, you can activate one harvested attack as a free action. Additionally, you gain one additional harvested attack use per short rest."}},
			},
		},
		{
			Name:        "Path of the Sage",
			Description: "You absorb the knowledge and magical insights of your prey, becoming a repository of arcane secrets and ancestral wisdom.",
			Features: map[int][]FeatureSeed{
				3:  {{Name: "Scholar's Legacy", Description: "You gain expertise in History or Arcana (double your proficiency bonus), plus two additional languages. You gain one additional Recipe Slot."}},
				10: {{Name: "Deep Insight", Description: "When you use Visceral Psychometry, you can ask 5 questions instead of 3 with Phylogenetic Recall. You also intuitively learn the Component recipe for one of the creature's abilities. You also harvest 1 new Component from the creature."}},
				14: {{Name: "Shared Consciousness", Description: "You can maintain Collective Consciousness on up to three allies simultaneously. When you share a harvested ability, you retain the full benefit yourself."}},
			},
		},
	}
}

func lorewrightLevelProgression() map[int]ClassLevelSeed {
	return map[int]ClassLevelSeed{
		1:  {Resources: map[string]int{"echo_slots": 0, "madness_die": 4}},
		2:  {Resources: map[string]int{"echo_slots": 1, "madness_die": 4}},
		3:  {Resources: map[string]int{"echo_slots": 1, "madness_die": 4}},
		4:  {Resources: map[string]int{"echo_slots": 1, "madness_die": 4}},
		5:  {Resources: map[string]int{"echo_slots": 2, "madness_die": 6}},
		6:  {Resources: map[string]int{"echo_slots": 2, "madness_die": 6}},
		7:  {Resources: map[string]int{"echo_slots": 2, "madness_die": 6}},
		8:  {Resources: map[string]int{"echo_slots": 2, "madness_die": 6}},
		9:  {Resources: map[string]int{"echo_slots": 3, "madness_die": 8}},
		10: {Resources: map[string]int{"echo_slots": 3, "madness_die": 8}},
		11: {Resources: map[string]int{"echo_slots": 3, "madness_die": 8}},
		12: {Resources: map[string]int{"echo_slots": 3, "madness_die": 8}},
		13: {Resources: map[string]int{"echo_slots": 4, "madness_die": 10}},
		14: {Resources: map[string]int{"echo_slots": 4, "madness_die": 10}},
		15: {Resources: map[string]int{"echo_slots": 4, "madness_die": 10}},
		16: {Resources: map[string]int{"echo_slots": 4, "madness_die": 10}},
		17: {Resources: map[string]int{"echo_slots": 5, "madness_die": 12}},
		18: {Resources: map[string]int{"echo_slots": 5, "madness_die": 12}},
		19: {Resources: map[string]int{"echo_slots": 5, "madness_die": 12}},
		20: {Resources: map[string]int{"echo_slots": 5, "madness_die": 20}},
	}
}

func lorewrightLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		// Level 1: Core harvest and psychometry
		1: {
			{Name: "Visceral Psychometry", Description: "As an action, you begin channeling to consume the liver of a creature (must have died within 1 hour). The process takes 1 minute and is interrupted if you take damage or move. You absorb its memories, allowing you to ask the DM up to 3 questions about the creature's life, secrets, or surroundings (Phylogenetic Recall)."},
			{Name: "Psychic Strain", Description: "When using Visceral Psychometry, you may suffer psychic backlash. Make a Wisdom Save with DC `10 + (Creature CR - (Your Level / 3))`. On a failure, you gain 1 Strain and must immediately roll on the Lorewright Madness Table. All accumulated Strain is removed when you finish a Long Rest."},
			{Name: "Component Scavenging", Description: "By spending 10 minutes harvesting a magical creature, you extract its essence as a usable Component (e.g., 'Ignis' from a Fire Elemental). Capacity: Wisdom modifier components."},
		},

		// Level 2: Harvest Bank begins
		2: {
			{Name: "Harvest Bank", Description: "You unlock the ability to store harvested abilities. You begin with 1 Harvest Slot, and your capacity increases as you level up. A Harvest Slot can hold any single ability you learn through Visceral Psychometry—be it a skill proficiency, a creature's attack, or a magical recipe. To replace an ability, you must overwrite the slot with a new one."},
			{Name: "Prey Instinct", Description: "Once per short rest, when you make an ability check using a harvested skill proficiency, you can gain advantage on the roll. Your predatory connection to the source creature sharpens your instincts."},
		},

		// Level 3: Archetype choice
		3: {{Name: "Past Life Archetype", Description: "Choose your archetype: Path of the Warlord or Path of the Sage. Your archetype reflects whose memories dominate your psyche."}},

		// Level 5: Expanded Harvest
		5: {
			{Name: "Expanded Harvest", Description: "Your Harvest Slot capacity increases to 2. You can now hold two different harvested abilities simultaneously."},
			{Name: "Collective Consciousness", Description: "As an action, you can share the benefits of one of your harvested abilities with a willing ally within 30 feet for 1 hour. You retain the ability while sharing it. Once per long rest."},
		},

		// Level 6: Predatory Techniques
		6: {
			{Name: "Predatory Techniques", Description: "You can use a harvested attack a number of times per short rest equal to your Wisdom modifier (minimum 1). When you use a harvested attack, you use your Wisdom for the attack and damage rolls."},
		},

		// Level 7: Essence Decanting
		7: {
			{Name: "Essence Decanting", Description: "You can spend 10 minutes converting a harvested component into a different component of the same category (e.g., one Primal component for another Primal component). Additionally, your Component Scavenging now yields +1 component per harvest (minimum 2 total)."},
		},

		// Level 9: Deep Storage
		9: {
			{Name: "Deep Storage", Description: "Your Harvest Slot capacity increases to 3. You can now maintain three total harvested abilities at once."},
			{Name: "Collective Consciousness Upgrade", Description: "Collective Consciousness can now target up to two allies simultaneously. Each ally can receive a different harvested ability's benefit."},
		},

		// Level 11: Spell Scribe
		11: {
			{Name: "Spell Scribe", Description: "You gain an intuitive understanding of magical biology. When you perform Visceral Psychometry on a creature with magical abilities, you can now choose to learn one of its abilities as a spell recipe, storing it in one of your Harvest Slots."},
		},

		// Level 13: Resilient Psyche
		13: {
			{Name: "Resilient Psyche", Description: "Your Harvest Slot capacity increases to 4. Your mind has hardened against the strain of holding multiple foreign memories. You can now remove all Strain when you finish a Long Rest."},
			{Name: "Rapid Imprint", Description: "When you perform Visceral Psychometry, you can now harvest two things from the same creature — any combination of skill, attack, or recipe — into two separate Harvest Slots."},
		},

		// Level 15: Apex Harvester
		15: {
			{Name: "Apex Harvester", Description: "Your harvested attacks can now be used a number of times per short rest equal to your proficiency bonus. When you cast a spell from a harvested recipe, the component cost is reduced by 1 (minimum 1 component)."},
		},

		// Level 17: Apex Predator's Soul
		17: {
			{Name: "Apex Predator's Soul", Description: "Your Harvest Slot capacity increases to 5."},
			{Name: "Predator's Trance", Description: "As a bonus action, enter a heightened state for 1 minute. While in this trance, you can use any of your harvested attacks without expending their per-rest uses. Once per Long Rest."},
		},

		// Level 18: Perfected Digestion
		18: {
			{Name: "Perfected Digestion", Description: "You no longer need to make a Psychic Strain saving throw when consuming creatures with a CR equal to or less than half your level (rounded down). Your body has become so efficient at processing foreign essences that weaker creatures pose no psychological risk."},
		},

		// Level 20: The Omega Point
		20: {
			{Name: "The Omega Point", Description: "You automatically succeed on all Psychic Strain saving throws — your mind has achieved perfect harmony with the lives you've consumed."},
			{Name: "Dual Imprint", Description: "Your total number of Harvest Slots is doubled."},
		},
	}
}

// LorewrightMadnessTable returns the 50 steampunk-related madness effects.
func LorewrightMadnessTable() map[int]string {
	return map[int]string{
		1:  "Your skin hardens into brass plates for 1 minute (+2 AC), but your speed is halved.",
		2:  "Steam vents open on your shoulders, creating a 10ft radius area of heavy obscurement (fog) centered on you for 1d4 rounds.",
		3:  "You speak only in ticker-tape code (clicks and whirs) for 1 minute.",
		4:  "Your eyes turn into camera lenses; you gain Darkvision 60ft but lose color vision for 1 hour.",
		5:  "Gravity reverses for you personally for 1 round.",
		6:  "You cough up a small, functional gear every time you speak for 1 hour.",
		7:  "Your blood turns into hot oil; take 1d4 fire damage but become immune to poison for 1 hour.",
		8:  "A spectral train whistle blows loudly from your location, audible out to 300ft.",
		9:  "You become magnetic; metal objects within 5ft are pulled towards you for 1 minute.",
		10: "Your hair turns into copper wires for 24 hours.",
		11: "You hear the constant ticking of a giant clock. Disadvantage on Perception checks for 10 minutes.",
		12: "One of your arms temporarily turns into a piston; Unarmed strikes deal 1d6 bludgeoning but you drop held items.",
		13: "You emit bright light in a 20ft radius and dim light for an additional 20ft (filament glow) for 10 minutes.",
		14: "Time slows down for you; you become Lethargic (Speed halved, -2 AC, no reactions, one action per turn) for 1 minute.",
		15: "Time speeds up; you become Accelerated (Speed doubled, +2 AC, additional action) for 1 round, then Lethargic for 1 round.",
		16: "You grow a vestigial mechanical limb that flails uselessly for 1 hour.",
		17: "Your voice amplifies like a phonograph (3x volume) for 10 minutes.",
		18: "You calculate probabilities aloud constantly. Disadvantage on Stealth checks for 1 hour.",
		19: "Your reflection in mirrors shows you as a clockwork automaton.",
		20: "Small lightning arcs between your fingers. Unarmed strikes deal +1 lightning damage, but you take 1 lightning damage per hit.",
		21: "You believe you are a construct and don't need to breathe (you do). Lasts 1 hour.",
		22: "Smoke belches from your mouth when you open it for 10 minutes.",
		23: "Your footsteps sound like heavy metal clanking for 1 hour.",
		24: "You perceive living creatures as skeletal structures with glowing hearts for 10 minutes.",
		25: "You gain a craving for coal or oil for 24 hours.",
		26: "Your shadow detaches and acts out your subconscious desires for 1 hour.",
		27: "A random non-magical item you touch turns to rusted scrap.",
		28: "You teleport 20ft to an unoccupied space in a random direction with a bamf of steam.",
		29: "Illusory cogs float around your head for 1 minute.",
		30: "You lose the ability to lie; you state facts like a computer for 1 hour.",
		31: "Your internal temperature rises; you are hot to the touch. Ice melts in your hand.",
		32: "You function as a compass (always know North) but get dizzy if you spin.",
		33: "Your fingernails turn into screwdriver heads for 24 hours.",
		34: "You can only move in cardinal directions (grid movement) for 1 minute.",
		35: "A pressure gauge appears on your chest; it redlines when you are angry.",
		36: "You emit a localized EMP; technological items within 10ft cease to function for 1 round.",
		37: "You forget the concept of 'softness'; everything feels hard as steel.",
		38: "Your veins glow neon blue (mana leak) for 1 hour.",
		39: "You produce a small wind-up toy every time you use a harvested component for 1 minute.",
		40: "Your size fluctuates slightly (grow/shrink 1 inch) every round for 1 minute.",
		41: "You smell of ozone and sulfur for 24 hours.",
		42: "Mechanized wings sprout from your back but can't fly, just flutter. Lasts 1 hour.",
		43: "You hear radio transmissions from the future for 10 minutes.",
		44: "Your touch magnetizes non-ferrous metals for 1 minute.",
		45: "You sweat lubricant/grease. Advantage to escape grapples, disadvantage to hold items for 1 hour.",
		46: "A monocle fuses to your face for 24 hours.",
		47: "You see mathematical equations floating over enemies (DM tells you AC of one enemy).",
		48: "Your heart stops beating and starts ticking for 1 hour.",
		49: "Steam pressure builds up; if you don't release it (scream/cast spell) in 1 round, you take 1d6 fire damage.",
		50: "System Reboot: You fall unconscious until the start of your next turn, then wake up with full HP or regain 1d4 harvested components of your choice.",
	}
}

// CalculatePsychicStrainDC determines if a Lorewright must make a saving throw against Psychic Strain
// when harvesting a creature. It returns true if a save is needed, and the DC for the check.
// Logic:
// - As the Lorewright levels up, they can safely harness more complex creature memories.
// - Safe Threshold: Creature CR <= (Level / 3). If below this, no save is required.
// - If above the threshold, a Wisdom Save is required.
// - DC = 10 + (Creature CR - Safe Threshold). The harder the creature is relative to your level, the harder the check.
// - Minimum DC is 10.
func CalculatePsychicStrainDC(level int, creatureCR float64) (bool, int) {
	// Calculate the CR threshold the character can handle safely.
	safeThreshold := float64(level) / 3.0

	if creatureCR <= safeThreshold {
		return false, 0
	}

	// Calculate DC
	excess := creatureCR - safeThreshold
	dc := 10 + int(excess+0.99) // Ceiling

	return true, dc
}
