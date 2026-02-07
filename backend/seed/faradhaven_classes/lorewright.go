package faradhaven_classes

// Lorewright returns the Lorewright class seed.
// A class that consumes creatures to absorb their memories and abilities,
// balancing power against the risk of madness. High wisdom and constitution,
// with echoes of past lives guiding their actions.
func Lorewright() FaradhavenClassSeed {
	archetypeLevel := 3
	return FaradhavenClassSeed{
		Name:           "The Lorewright",
		Description:    "Consume the flesh of your prey to absorb their memories and abilities. Balance the power you gain against the fracturing of your mind as you become the sum of all you have eaten.",
		HitDie:         10,
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
		ResourceType:        "echo_slots",
		ResourceName:        "Echo Slots",
		ResourceRestoreType: "none", // Echoes are replaced, not restored
	}
}

func lorewrightArchetypes() []ArchetypeSeed {
	return []ArchetypeSeed{
		{
			Name:        "Path of the Warlord",
			Description: "You channel the combat memories of fallen warriors, gaining martial prowess and the ability to strike with the fury of a hundred battles.",
			Features: map[int][]FeatureSeed{
				3:  {{Name: "Martial Legacy", Description: "You gain proficiency with heavy armor and one martial weapon of your choice. The echoes of warriors past guide your hands in combat."}},
				10: {{Name: "Predator's Strike", Description: "When you hit a creature type you have consumed, add your Wisdom modifier to damage. The memories of hunting that prey grant deadly precision."}},
				14: {{Name: "Battle Trance", Description: "When you enter combat, you can activate one Echo as a free action to gain its benefits immediately. Additionally, you can use Cycle of Rebirth twice per Long Rest."}},
			},
		},
		{
			Name:        "Path of the Sage",
			Description: "You absorb the knowledge and magical insights of your prey, becoming a repository of arcane secrets and ancestral wisdom.",
			Features: map[int][]FeatureSeed{
				3:  {{Name: "Scholar's Legacy", Description: "You gain expertise in History or Arcana (double your proficiency bonus), plus two additional languages of your choice. The minds you consume expand your understanding."}},
				10: {{Name: "Deep Insight", Description: "When you consume a creature, you can ask two questions instead of one with Phylogenetic Recall. Additionally, you learn one spell the creature could cast (if any) which you can cast once before your next Long Rest using Wisdom."}},
				14: {{Name: "Shared Consciousness", Description: "You can maintain two Collective Consciousness transfers simultaneously. When you transfer an Echo, you retain partial benefit (half the bonus or limited version of the trait)."}},
			},
		},
	}
}

func lorewrightLevelProgression() map[int]ClassLevelSeed {
	// BardicInspiration represents Madness Die scaling (d4 -> d6 -> d8 -> d10 -> d12 -> d20)
	// EchoSlots scale: 0 at 1, 1 at 2, 2 at 5, 3 at 9, 4 at 13, 5 at 17
	return map[int]ClassLevelSeed{
		1:  {EchoSlots: 0, BardicInspiration: 4},                       // Madness Die: d4, no Echo Slots yet
		2:  {EchoSlots: 1, BardicInspiration: 4},                       // 1 Echo Slot
		3:  {EchoSlots: 1, BardicInspiration: 4},
		4:  {EchoSlots: 1, BardicInspiration: 4},
		5:  {EchoSlots: 2, BardicInspiration: 6},                       // Madness Die: d6, 2 Echo Slots
		6:  {EchoSlots: 2, BardicInspiration: 6},
		7:  {EchoSlots: 2, BardicInspiration: 6},
		8:  {EchoSlots: 2, BardicInspiration: 6},
		9:  {EchoSlots: 3, BardicInspiration: 8},                       // Madness Die: d8, 3 Echo Slots
		10: {EchoSlots: 3, BardicInspiration: 8},
		11: {EchoSlots: 3, BardicInspiration: 8},
		12: {EchoSlots: 3, BardicInspiration: 8},
		13: {EchoSlots: 4, BardicInspiration: 10},                      // Madness Die: d10, 4 Echo Slots
		14: {EchoSlots: 4, BardicInspiration: 10},
		15: {EchoSlots: 4, BardicInspiration: 10},
		16: {EchoSlots: 4, BardicInspiration: 10},
		17: {EchoSlots: 5, BardicInspiration: 12},                      // Madness Die: d12, 5 Echo Slots
		18: {EchoSlots: 5, BardicInspiration: 12},
		19: {EchoSlots: 5, BardicInspiration: 12},
		20: {EchoSlots: 5, BardicInspiration: 20},                      // Madness Die: d20, 5 Echo Slots (doubled capacity)
	}
}

func lorewrightLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		// Level 1: Core features
		// Echo Slots: 0, Madness Die: d4
		1: {
			{Name: "Visceral Psychometry", Description: "Consume the liver of a freshly killed creature (died within 1 hour) over 1 minute to absorb its memories. This is your primary method of learning: you gain access to the creature's recent experiences and knowledge."},
			{Name: "The Fracture", Description: "When using Visceral Psychometry, make a Wisdom Save against a DC determined by the creature's power (DC 10 for CR 0-4, DC 15 for CR 5-10, DC 20 for CR 11+). On failure, gain 1 Trauma. 1 Trauma: Disadvantage on Charisma checks. 2 Trauma: Gain a temporary flaw. 3 Trauma: Confusion. Remove 1 Trauma per Long Rest."},
			{Name: "Component Scavenging", Description: "You can harvest Components from the bodies of magical creatures. By spending 10 minutes harvesting specific organs, you can extract the creature's essence as a usable Component (e.g., 'Ignis' from a Fire Elemental, 'Flight' from a Roc). You can add these collected components to your repertoire to be used with your abilities or slotted into Echoes."},
			{Name: "Madness Die: d4", Description: "You have no Echo Slots yet. High-risk phase: stick to small game to avoid madness while learning your craft."},
		},

		// Level 2: Echo Slots: 1, Madness Die: d4
		2: {{Name: "Somatic Echoes", Description: "You gain 1 Echo Slot. You can fill this slot with a Skill, Tool Proficiency, Language, or Component you have scavenged. To add a component or trait, simply declare you are slotting it after a harvest. Replacing an Echo requires consuming a new source."}},

		// Level 3
		3: {{Name: "Past Life Archetype", Description: "Choose your archetype at this level."}},

		// Level 4
		4: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1. Wisdom and Constitution are recommended."}},

		// Level 5: Echo Slots: 2, Madness Die: d6
		5: {{Name: "Expanded Mind", Description: "Your Madness Die increases to d6, and you gain 2 Echo Slots. Your capacity to hold foreign memories expands."}},

		// Level 8
		8: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},

		// Level 9: Echo Slots: 3, Madness Die: d8
		9: {{Name: "Deep Storage", Description: "Your Madness Die increases to d8, and you gain 3 Echo Slots."}},

		// Level 12
		12: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},

		// Level 13: Echo Slots: 4, Madness Die: d10
		13: {{Name: "Resilient Psyche", Description: "Your Madness Die increases to d10, and you gain 4 Echo Slots."}},

		// Level 16
		16: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},

		// Level 17: Echo Slots: 5, Madness Die: d12
		17: {{Name: "Apex Predator's Soul", Description: "Your Madness Die increases to d12, and you gain 5 Echo Slots. As a bonus action, enter a heightened state for 1 minute. You can use any Action available to a creature stored in your Echo Slots. If an action deals damage, use your Wisdom modifier for the attack and damage rolls. If an action has a Recharge, you can use it once per transformation. Once used, you cannot use this feature again until you finish a Long Rest."}},

		// Level 19
		19: {{Name: "Ability Score Improvement", Description: "Increase one ability score by 2, or two ability scores by 1."}},

		// Level 20
		20: {{Name: "The Omega Point", Description: "Your Madness Die becomes d20. You automatically succeed on the Saving Throw for The Fracture. You can hold two Echoes in a single slot, doubling your capacity."}},
	}
}

// LorewrightMadnessTable returns the 50 steampunk-related madness effects.
func LorewrightMadnessTable() map[int]string {
	return map[int]string{
		1:  "Your skin hardens into brass plates for 1 minute (+2 AC), but your speed is halved.",
		2:  "Steam vents open on your shoulders, creating a 10ft radius fog cloud centered on you for 1d4 rounds.",
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
		13: "You emit bright light in a 20ft radius (filament glow) for 10 minutes.",
		14: "Time slows down for you; you are under the effect of the Slow spell for 1 minute.",
		15: "Time speeds up; you are under the effect of the Haste spell for 1 round, then lethargic for 1 round.",
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
		28: "You teleport 20ft in a random direction with a bamf of steam.",
		29: "Illusory cogs float around your head for 1 minute.",
		30: "You lose the ability to lie; you state facts like a computer for 1 hour.",
		31: "Your internal temperature rises; you are hot to the touch. Ice melts in your hand.",
		32: "You function as a compass (always know North) but get dizzy if you spin.",
		33: "Your fingernails turn into screwdriver heads for 24 hours.",
		34: "You can only move in cardinal directions (grid movement) for 1 minute.",
		35: "A pressure gauge appears on your chest; it redlines when you are angry.",
		36: "You emit a localized EMP; technological items within 10ft flicker or stall.",
		37: "You forget the concept of 'softness'; everything feels hard as steel.",
		38: "Your veins glow neon blue (mana leak) for 1 hour.",
		39: "You produce a small wind-up toy every time you cast a spell for 1 minute.",
		40: "Your size fluctuates slightly (grow/shrink 1 inch) every round for 1 minute.",
		41: "You smell of ozone and sulfur for 24 hours.",
		42: "Mechanized wings sprout from your back but can't fly, just flutter. Lasts 1 hour.",
		43: "You hear radio transmissions from the future for 10 minutes.",
		44: "Your touch magnetizes non-ferrous metals for 1 minute.",
		45: "You sweat lubricant/grease. Advantage to escape grapples, disadvantage to hold items for 1 hour.",
		46: "A monocle fuses to your face for 24 hours.",
		47: "You see mathematical equations floating over enemies (DM tells you AC of one enemy).",
		48: "Your heart stops beating and starts ticking for 1 hour.",
		49: "Steam pressure builds up; if you don't release it (scream/cast spell) in 1 round, you explode for 1d6 fire damage.",
		50: "System Reboot: You fall unconscious until the start of your next turn, then wake up with full HP or a spell slot restored.",
	}
}

// CheckLorewrightMadness determines if the Lorewright must roll for madness when using a creature ability.
// It returns true if a roll is needed, and the DC for the check.
// Logic:
// - As the Lorewright levels up, they can safely harness more complex creature memories.
// - Safe Threshold: Creature CR <= (Level / 3). If below this, no madness check is required.
// - If above the threshold, a Wisdom Save is required.
// - DC = 10 + (Creature CR - Safe Threshold). The harder the creature is relative to your level, the harder the check.
// - Minimum DC is 10.
func CheckLorewrightMadness(level int, creatureCR float64) (bool, int) {
	// Calculate the CR threshold the character can handle safely.
	// Level 1: CR 0.33 (Safe with CR 0, 1/8, 1/4)
	// Level 3: CR 1 (Safe with CR 1)
	// Level 6: CR 2
	// Level 20: CR 6.6
	safeThreshold := float64(level) / 3.0

	if creatureCR <= safeThreshold {
		return false, 0
	}

	// Calculate DC
	// Base DC 10 + how much the creature exceeds the safe threshold (rounded up)
	excess := creatureCR - safeThreshold
	dc := 10 + int(excess+0.99) // Ceiling

	return true, dc
}