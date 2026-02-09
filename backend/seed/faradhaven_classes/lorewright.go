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
				3:  {{Name: "Martial Legacy", Description: "You gain proficiency with heavy armor and one martial weapon. Your melee weapon attacks deal an additional 1d4 damage of your weapon's type."}},
				10: {{Name: "Predator's Strike", Description: "When you hit a creature type you have consumed within the last 24 hours, add your Wisdom modifier to the damage roll."}},
				14: {{Name: "Battle Trance", Description: "When you enter combat, you can activate one Echo as a free action. Additionally, you can use Cycle of Rebirth twice per short rest."}},
			},
		},
		{
			Name:        "Path of the Sage",
			Description: "You absorb the knowledge and magical insights of your prey, becoming a repository of arcane secrets and ancestral wisdom.",
			Features: map[int][]FeatureSeed{
				3:  {{Name: "Scholar's Legacy", Description: "You gain expertise in History or Arcana (double your proficiency bonus), plus two additional languages. You gain one additional Component Scavenging slot."}},
				10: {{Name: "Deep Insight", Description: "When you use Visceral Psychometry, you can ask 5 questions instead of 3 with Phylogenetic Recall. You also learn one spell the creature could cast, which you can cast once using Wisdom (no slot required) before your next Long Rest."}},
				14: {{Name: "Shared Consciousness", Description: "You can maintain Collective Consciousness on up to three allies simultaneously. When you share an Echo, you retain the full benefit yourself."}},
			},
		},
	}
}

func lorewrightLevelProgression() map[int]ClassLevelSeed {
	// BardicInspiration repurposed as Madness Die size (d4 -> d6 -> d8 -> d10 -> d12 -> d20)
	//   Higher die = less chance of rolling a 1 = safer Echo usage
	// EchoSlots = number of harvested abilities that can be stored simultaneously
	//   0 at Lvl 1 (no slots yet), 1 at Lvl 2, 2 at Lvl 5, 3 at Lvl 9, 4 at Lvl 13, 5 at Lvl 17
	//   At Lvl 20, each slot holds 2 Echoes (Dual Imprint), effectively 10 stored abilities
	return map[int]ClassLevelSeed{
		1:  {EchoSlots: 0, BardicInspiration: 4},  // No Echo Slots — harvest & psychometry only
		2:  {EchoSlots: 1, BardicInspiration: 4},  // First Echo Slot: save 1 harvested skill/action
		3:  {EchoSlots: 1, BardicInspiration: 4},  // Archetype choice
		4:  {EchoSlots: 1, BardicInspiration: 4},  // ASI
		5:  {EchoSlots: 2, BardicInspiration: 6},  // +1 slot (2 total), Madness Die -> d6
		6:  {EchoSlots: 2, BardicInspiration: 6},
		7:  {EchoSlots: 2, BardicInspiration: 6},
		8:  {EchoSlots: 2, BardicInspiration: 6},  // ASI
		9:  {EchoSlots: 3, BardicInspiration: 8},  // +1 slot (3 total), Madness Die -> d8
		10: {EchoSlots: 3, BardicInspiration: 8},
		11: {EchoSlots: 3, BardicInspiration: 8},
		12: {EchoSlots: 3, BardicInspiration: 8},  // ASI
		13: {EchoSlots: 4, BardicInspiration: 10}, // +1 slot (4 total), Madness Die -> d10, Rapid Imprint
		14: {EchoSlots: 4, BardicInspiration: 10},
		15: {EchoSlots: 4, BardicInspiration: 10},
		16: {EchoSlots: 4, BardicInspiration: 10}, // ASI
		17: {EchoSlots: 5, BardicInspiration: 12}, // +1 slot (5 total, max), Madness Die -> d12, Predator's Trance
		18: {EchoSlots: 5, BardicInspiration: 12},
		19: {EchoSlots: 5, BardicInspiration: 12}, // ASI
		20: {EchoSlots: 5, BardicInspiration: 20}, // Madness Die -> d20, Dual Imprint (each slot holds 2)
	}
}

func lorewrightLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		// Level 1: Core harvest and psychometry (no Echo Slots yet)
		1: {
			{Name: "Visceral Psychometry", Description: "As an action, you consume the liver of a creature (died within 1 hour) over 1 minute. You absorb its memories, allowing you to ask the DM up to 3 questions about the creature's life, secrets, or surroundings (Phylogenetic Recall)."},
			{Name: "The Fracture", Description: "When using Visceral Psychometry, make a Wisdom Save (DC 10 for CR 0-4, DC 15 for CR 5-10, DC 20 for CR 11+). On failure, gain 1 Trauma. 1 Trauma: Disadvantage on Charisma checks. 2 Trauma: Gain a temporary flaw. 3 Trauma: Confusion (as the spell). Remove 1 Trauma per Long Rest."},
			{Name: "Component Scavenging", Description: "By spending 10 minutes harvesting a magical creature, you extract its essence as a usable Component (e.g., 'Ignis' from a Fire Elemental). Capacity: Wisdom modifier components."},
			{Name: "Madness Die: d4", Description: "When you activate an Echo Slot to use a creature's trait or action, you must roll your Madness Die (d4). On a 1, roll on the Lorewright Madness Table. The Madness Die scales with level. You do not yet have Echo Slots at this level — this rule applies once you gain them at level 2."},
		},

		// Level 2: Echo Slots: 1, Madness Die: d4
		2: {
			{Name: "Somatic Echoes (Echo Slot System)", Description: "You gain 1 Echo Slot. Echo Slots are your memory banks — each one holds a single ability harvested from a consumed creature. To save a new Echo: (1) Perform Visceral Psychometry on a creature. (2) Choose one of its Skills, Tool Proficiencies, Languages, or Actions to imprint into an empty Echo Slot. (3) You gain that ability as long as it remains slotted. Each slot holds exactly one Echo. To replace an Echo, you must perform a new harvest and overwrite the slot — the old Echo is lost permanently."},
			{Name: "Cycle of Rebirth", Description: "As a bonus action, you can swap which of your slotted Echoes is currently active. Only one Echo can be active at a time at this level. You can use Cycle of Rebirth once per short rest. Using an active Echo's action or trait triggers your Madness Die roll."},
		},

		// Level 3
		3: {{Name: "Past Life Archetype", Description: "Choose your archetype: Path of the Warlord or Path of the Sage. Your archetype reflects whose memories dominate your psyche."}},

		// Level 5: Echo Slots: 2, Madness Die: d6
		5: {
			{Name: "Expanded Mind (2 Echo Slots)", Description: "Your Echo Slot capacity increases to 2. You can now hold two different harvested abilities simultaneously — for example, a wolf's Keen Hearing (Perception proficiency) in one slot and a goblin's Nimble Escape (Disengage as bonus action) in another. To fill the new slot, perform another harvest and choose an ability to imprint. Your Madness Die increases to d6."},
			{Name: "Collective Consciousness", Description: "As an action, you can share the benefits of one of your active Echoes with a willing ally within 30 feet for 1 hour. The ally gains the slotted skill, proficiency, or ability for the duration. You retain the Echo while sharing it. Once per long rest."},
		},

		// Level 9: Echo Slots: 3, Madness Die: d8
		9: {
			{Name: "Deep Storage (3 Echo Slots)", Description: "Your Echo Slot capacity increases to 3. You can now maintain three harvested abilities at once. When you perform a harvest, you choose which slot to fill or overwrite — you are never forced to replace an existing Echo unless all slots are full and you want the new ability. Your Madness Die increases to d8."},
			{Name: "Collective Consciousness Upgrade", Description: "Collective Consciousness can now target up to two allies simultaneously. Each ally can receive a different Echo's benefit."},
		},

		// Level 13: Echo Slots: 4, Madness Die: d10
		13: {
			{Name: "Resilient Psyche (4 Echo Slots)", Description: "Your Echo Slot capacity increases to 4. Your mind has hardened against the strain of holding multiple foreign memories. Your Madness Die increases to d10. Additionally, you can now remove all Trauma when you finish a Long Rest (instead of only 1)."},
			{Name: "Rapid Imprint", Description: "When you perform Visceral Psychometry, you can now choose two abilities from the same creature to slot into two separate Echo Slots (if you have empty slots or choose to overwrite)."},
		},

		// Level 17: Echo Slots: 5, Madness Die: d12
		17: {
			{Name: "Apex Predator's Soul (5 Echo Slots)", Description: "Your Echo Slot capacity increases to 5 — the maximum a mortal mind can sustain. You can hold five distinct harvested abilities simultaneously. Your Madness Die increases to d12."},
			{Name: "Predator's Trance", Description: "As a bonus action, enter a heightened state for 1 minute. While in this trance, all of your Echo Slots are active simultaneously — you benefit from every slotted skill, proficiency, and trait at once, and can use any slotted Action on your turn. If an action deals damage, use Wisdom for attack and damage rolls. Each action used still triggers a Madness Die roll. Once per Long Rest."},
		},

		// Level 20
		20: {
			{Name: "The Omega Point", Description: "Your Madness Die becomes d20 (you almost never trigger madness). You automatically succeed on Saving Throws for The Fracture — consumption no longer risks Trauma."},
			{Name: "Dual Imprint", Description: "Each Echo Slot can now hold two different Echoes simultaneously, effectively doubling your capacity to 10 stored abilities. When you harvest a creature, you can save one ability per slot as normal, or stack a second ability into an existing slot. Both abilities in a dual-imprinted slot are active when that slot is active."},
		},
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