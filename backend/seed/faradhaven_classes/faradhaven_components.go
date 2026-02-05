package faradhaven_classes

// ComponentClassMapping pairs a component with the class that has access to it.
// Supports many-to-many: a component could appear in multiple ClassNames in the future.
type ComponentClassMapping struct {
	Component ComponentSeed
	ClassName string
}

// ClassComponentNames maps class names to the spell system component names they have access to.
// These reference components defined in SpellSystemComponents().
// Each class has 15-60 components based on their magical proficiency.
// Advanced spellcasters (Rift Weaver, Powder Mage) have access to most components.
func ClassComponentNames() map[string][]string {
	return map[string][]string{
		// =============================================================================
		// The Rift Weaver - dimensional/psychic magic specialist (ADVANCED CASTER)
		// Master of space, time, mind, and planar magic - has access to almost everything
		// =============================================================================
		"The Rift Weaver": {
			// Primordial Elements (9/14) - energy and mind-based elements
			"Aqua", "Aer", "Fulgur", "Vita", "Umbra", "Lux", "Arcanum", "Sonus", "Psi",
			// Physical Forces (5/9) - telekinetic forces
			"Push", "Pull", "Lift", "Pierce", "Focus",
			// Thermodynamic (2/4) - energy states
			"Heat", "Cool",
			// Spatial & Temporal (7/7) - core domain, full mastery
			"Expand", "Shrink", "Haste", "Slow", "Teleport", "Echo", "Extreme",
			// Life (3/4) - life force manipulation
			"Mend", "Wither", "Revive",
			// Spell Shapes (7/8) - dimensional delivery
			"Projectile", "Beam", "Nova", "Wall", "Zone", "Self", "Cone",
			// Abjuration (4/4) - dimensional protection
			"Reflect", "Dispel", "Absorb", "Seal",
			// Conjuration (2/3) - summoning across planes
			"Summon", "Portal",
			// Divination (4/4) - seeing through dimensions
			"Sight", "Identify", "Predict", "Link",
			// Enchantment (3/4) - mental domination
			"Fear", "Frenzy", "Command",
			// Evocation (4/4) - raw energy
			"Channel", "Chain", "Imbue", "Rupture",
			// Illusion (5/5) - reality manipulation
			"Decoy", "Invisible", "Silence", "Disguise", "Phantom",
			// Necromancy (2/4) - soul manipulation
			"Soul", "Curse",
			// Transmutation (3/3) - matter alteration
			"Transmute", "Burrow", "Fuse",
		},

		// =============================================================================
		// The Powder Mage - ranged elemental combat specialist (ADVANCED CASTER)
		// Master of elements, projectiles, and destructive magic
		// =============================================================================
		"The Powder Mage": {
			// Primordial Elements (12/14) - elemental mastery
			"Ignis", "Aqua", "Terra", "Aer", "Fulgur", "Ferrum", "Umbra", "Lux", "Arcanum", "Sonus", "Acidum", "Psi",
			// Physical Forces (9/9) - projectile physics mastery
			"Push", "Pull", "Lift", "Crush", "Spin", "Pierce", "Scatter", "Focus", "Homing",
			// Thermodynamic (4/4) - temperature extremes
			"Heat", "Cool", "Conduct", "Insulate",
			// Spatial & Temporal (6/7) - combat timing
			"Expand", "Shrink", "Haste", "Slow", "Echo", "Extreme",
			// Life (2/4) - combat sustain
			"Mend", "Wither",
			// Spell Shapes (8/8) - full delivery mastery
			"Projectile", "Beam", "Nova", "Wall", "Zone", "Self", "Trap", "Cone",
			// Abjuration (3/4) - combat protection
			"Reflect", "Dispel", "Absorb",
			// Conjuration (2/3) - battlefield control
			"Summon", "Grease",
			// Divination (3/4) - targeting
			"Sight", "Identify", "Predict",
			// Enchantment (2/4) - combat control
			"Fear", "Taunt",
			// Evocation (4/4) - core damage domain
			"Channel", "Chain", "Imbue", "Rupture",
			// Illusion (3/5) - combat deception
			"Decoy", "Phantom", "Silence",
			// Transmutation (2/3) - material manipulation
			"Transmute", "Fuse",
		},

		// =============================================================================
		// The Lorewright - information/divination specialist (STRONG CASTER)
		// Scholar of all magic, extensive theoretical knowledge
		// =============================================================================
		"The Lorewright": {
			// Primordial Elements (8/14) - scholarly elements
			"Aqua", "Aer", "Fulgur", "Vita", "Lux", "Arcanum", "Psi", "Sanctus",
			// Physical Forces (4/9) - precise manipulation
			"Push", "Pull", "Focus", "Homing",
			// Thermodynamic (3/4) - energy theory
			"Heat", "Cool", "Conduct",
			// Spatial & Temporal (5/7) - time/space study
			"Expand", "Shrink", "Slow", "Echo", "Teleport",
			// Life (3/4) - life studies
			"Mend", "Wither", "Revive",
			// Spell Shapes (6/8) - common shapes
			"Projectile", "Beam", "Wall", "Zone", "Self", "Trap",
			// Abjuration (4/4) - full protective knowledge
			"Reflect", "Dispel", "Absorb", "Seal",
			// Conjuration (3/3) - summoning lore
			"Summon", "Portal", "Grease",
			// Divination (4/4) - core domain, full mastery
			"Sight", "Identify", "Predict", "Link",
			// Enchantment (4/4) - mind studies
			"Fear", "Taunt", "Frenzy", "Command",
			// Evocation (2/4) - energy channeling
			"Channel", "Chain",
			// Illusion (5/5) - perception manipulation
			"Decoy", "Invisible", "Silence", "Disguise", "Phantom",
			// Necromancy (2/4) - forbidden knowledge
			"Soul", "Curse",
			// Transmutation (3/3) - matter theory
			"Transmute", "Burrow", "Fuse",
		},

		// =============================================================================
		// The Sanguinist - life/death magic specialist (STRONG CASTER)
		// Master of blood, souls, and the boundary between life and death
		// =============================================================================
		"The Sanguinist": {
			// Primordial Elements (6/14) - life/death elements
			"Aqua", "Vita", "Umbra", "Arcanum", "Acidum", "Sanctus",
			// Physical Forces (4/9) - blood manipulation
			"Pull", "Crush", "Pierce", "Focus",
			// Thermodynamic (2/4) - body temperature
			"Heat", "Cool",
			// Spatial & Temporal (4/7) - life force timing
			"Slow", "Haste", "Echo", "Extreme",
			// Life (4/4) - core domain, full mastery
			"Mend", "Wither", "Revive", "Drain",
			// Spell Shapes (6/8) - ritualistic shapes
			"Beam", "Nova", "Wall", "Zone", "Self", "Cone",
			// Abjuration (3/4) - life protection
			"Absorb", "Seal", "Dispel",
			// Conjuration (2/3) - summoning undead
			"Summon", "Portal",
			// Divination (3/4) - life sensing
			"Sight", "Identify", "Link",
			// Enchantment (4/4) - terror and control
			"Fear", "Taunt", "Frenzy", "Command",
			// Evocation (3/4) - life energy
			"Channel", "Chain", "Rupture",
			// Necromancy (4/4) - core domain, full mastery
			"Bone", "Soul", "Curse", "Sacrifice",
			// Transmutation (2/3) - flesh shaping
			"Transmute", "Fuse",
		},

		// =============================================================================
		// The Ironwright - lightning/metal tech specialist (MODERATE CASTER)
		// Engineer-mage focused on constructs, electricity, and machinery
		// =============================================================================
		"The Ironwright": {
			// Primordial Elements (6/14) - tech elements
			"Terra", "Fulgur", "Ferrum", "Arcanum", "Sonus", "Lux",
			// Physical Forces (6/9) - mechanical forces
			"Push", "Pull", "Crush", "Spin", "Focus", "Homing",
			// Thermodynamic (4/4) - full thermal control for forging
			"Heat", "Cool", "Conduct", "Insulate",
			// Spatial & Temporal (3/7) - precision timing
			"Haste", "Slow", "Extreme",
			// Life (1/4) - repair
			"Mend",
			// Spell Shapes (7/8) - construct delivery
			"Projectile", "Beam", "Nova", "Wall", "Zone", "Trap", "Cone",
			// Abjuration (3/4) - tech shields
			"Reflect", "Absorb", "Seal",
			// Conjuration (2/3) - construct summoning
			"Summon", "Grease",
			// Divination (2/4) - scanning
			"Sight", "Identify",
			// Evocation (4/4) - energy channeling
			"Channel", "Chain", "Imbue", "Rupture",
			// Transmutation (3/3) - core domain for crafting
			"Transmute", "Burrow", "Fuse",
		},

		// =============================================================================
		// The Mutagen - biological transformation specialist (MODERATE CASTER)
		// Alchemist focused on body modification and organic matter
		// =============================================================================
		"The Mutagen": {
			// Primordial Elements (5/14) - organic elements
			"Aqua", "Vita", "Acidum", "Psi", "Sanctus",
			// Physical Forces (4/9) - body mechanics
			"Push", "Crush", "Pierce", "Focus",
			// Thermodynamic (2/4) - metabolism
			"Heat", "Cool",
			// Spatial & Temporal (5/7) - body alteration
			"Expand", "Shrink", "Haste", "Slow", "Extreme",
			// Life (4/4) - core domain, full mastery
			"Mend", "Wither", "Revive", "Drain",
			// Spell Shapes (1/8) - self-focused only
			"Self",
			// Abjuration (2/4) - natural resistance
			"Absorb", "Insulate",
			// Conjuration (1/3) - organic creation
			"Grease",
			// Divination (2/4) - biological sensing
			"Sight", "Identify",
			// Enchantment (2/4) - pheromones
			"Fear", "Frenzy",
			// Illusion (2/5) - camouflage
			"Disguise", "Invisible",
			// Necromancy (2/4) - life/death boundary
			"Curse", "Sacrifice",
			// Transmutation (3/3) - core domain
			"Transmute", "Burrow", "Fuse",
		},

		// =============================================================================
		// The Vapor Blade - stealth/poison specialist (FOCUSED CASTER)
		// Assassin-mage focused on shadows, toxins, and infiltration
		// =============================================================================
		"The Vapor Blade": {
			// Primordial Elements (5/14) - stealth elements
			"Aer", "Umbra", "Acidum", "Psi", "Sonus",
			// Physical Forces (5/9) - precise strikes
			"Push", "Spin", "Pierce", "Scatter", "Focus",
			// Thermodynamic (2/4) - poison states
			"Heat", "Cool",
			// Spatial & Temporal (4/7) - timing and escape
			"Haste", "Slow", "Teleport", "Shrink",
			// Life (2/4) - poison effects
			"Wither", "Drain",
			// Spell Shapes (5/8) - assassination delivery
			"Projectile", "Zone", "Self", "Trap", "Cone",
			// Abjuration (1/4) - evasion
			"Absorb",
			// Divination (2/4) - target awareness
			"Sight", "Predict",
			// Enchantment (2/4) - distraction
			"Fear", "Taunt",
			// Evocation (2/4) - weapon enhancement
			"Imbue", "Rupture",
			// Illusion (5/5) - core domain, full mastery
			"Decoy", "Invisible", "Silence", "Disguise", "Phantom",
			// Necromancy (1/4) - poison curses
			"Curse",
			// Transmutation (2/3) - escape routes
			"Burrow", "Transmute",
		},

		// =============================================================================
		// The Piston Brawler - physical force melee specialist (FOCUSED CASTER)
		// Warrior-mage focused on kinetic energy and physical enhancement
		// =============================================================================
		"The Piston Brawler": {
			// Primordial Elements (5/14) - physical elements
			"Terra", "Ferrum", "Arcanum", "Sonus", "Fulgur",
			// Physical Forces (9/9) - core domain, full mastery
			"Push", "Pull", "Lift", "Crush", "Spin", "Pierce", "Scatter", "Focus", "Homing",
			// Thermodynamic (2/4) - body heat
			"Heat", "Insulate",
			// Spatial & Temporal (4/7) - combat timing
			"Expand", "Haste", "Slow", "Extreme",
			// Life (2/4) - combat sustain
			"Mend", "Drain",
			// Spell Shapes (6/8) - melee delivery
			"Nova", "Wall", "Zone", "Self", "Trap", "Cone",
			// Abjuration (2/4) - physical defense
			"Reflect", "Absorb",
			// Conjuration (1/3) - battlefield hazard
			"Grease",
			// Enchantment (3/4) - combat psychology
			"Fear", "Taunt", "Frenzy",
			// Evocation (3/4) - physical enhancement
			"Channel", "Imbue", "Rupture",
			// Transmutation (1/3) - material hardening
			"Fuse",
		},
	}
}

// AllComponentClassMappings returns all spell components and their class associations.
// Components reference the universal spell system components by name.
// ClassComponent links are created from these mappings.
func AllComponentClassMappings() []ComponentClassMapping {
	// Build mappings by looking up components from the spell system
	spellComponents := SpellSystemComponents()
	componentMap := make(map[string]ComponentSeed)
	for _, c := range spellComponents {
		componentMap[c.Name] = c
	}

	var mappings []ComponentClassMapping
	for className, componentNames := range ClassComponentNames() {
		for _, compName := range componentNames {
			if comp, ok := componentMap[compName]; ok {
				mappings = append(mappings, ComponentClassMapping{
					Component: comp,
					ClassName: className,
				})
			}
		}
	}
	return mappings
}
