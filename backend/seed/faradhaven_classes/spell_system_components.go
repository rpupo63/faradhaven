package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// SpellSystemComponents returns all universal spell system components.
// These are the building blocks of the magic system, organized into categories:
// - Primordial Elements ("Nouns") - base materials that give spells substance
// - Physical Forces ("Verbs") - dictate how spells move or interact with physics
// - Thermodynamic & State Modifiers ("Adjectives") - change state of matter
// - Spatial & Temporal Forces - manipulate the fabric of the world
// - Biological & Transmutation - interact with living things or change materials
// - Spell Shapes ("Container") - define how spells are delivered
func SpellSystemComponents() []ComponentSeed {
	return []ComponentSeed{
		// =============================================================================
		// I. THE PRIMORDIAL ELEMENTS (The "Nouns")
		// Base materials - a spell usually requires at least one to give it "substance"
		// =============================================================================
		{
			Name:        "Ignis",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Raw plasma, thermal energy, flame substance",
			Element:     "fire",
		},
		{
			Name:        "Aqua",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Liquid, fluid dynamics, dousing",
			Element:     "water",
		},
		{
			Name:        "Terra",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Dirt, stone, physical mass",
			Element:     "earth",
		},
		{
			Name:        "Aer",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Gas, wind, pressure",
			Element:     "air",
		},
		{
			Name:        "Fulgur",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Electricity, charged particles, current, voltage",
			Element:     "lightning",
		},
		{
			Name:        "Ferrum",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Hardness, conductivity, magnetism",
			Element:     "metal",
		},
		{
			Name:        "Vita",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Plants, vines, wood, life force, organic growth",
			Element:     "nature",
		},
		{
			Name:        "Umbra",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Darkness, absence of light, void, shadow substance",
			Element:     "shadow",
		},
		{
			Name:        "Lux",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Illumination, radiant energy, photons, pure light",
			Element:     "light",
		},
		{
			Name:        "Arcanum",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Pure mana, purple energy, force fields",
			Element:     "arcane",
		},
		{
			Name:        "Sonus",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Sound waves, vibration, thunder, shockwaves",
			Element:     "thunder",
		},
		{
			Name:        "Acidum",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Corrosive substances, acid, chemical dissolution",
			Element:     "acid",
		},
		{
			Name:        "Psi",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Psychic energy, mental force, thought-based attacks",
			Element:     "psychic",
		},
		{
			Name:        "Sanctus",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryPrimordial,
			Description: "Divine radiance, holy energy, sacred power",
			Element:     "radiant",
		},

		// =============================================================================
		// II. PHYSICAL FORCES (The "Verbs")
		// Dictate how the spell moves or interacts with the physics engine
		// =============================================================================
		{
			Name:        "Push",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryPhysical,
			Description: "Applies force away from the origin (repel)",
			Element:     "force",
		},
		{
			Name:        "Pull",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryPhysical,
			Description: "Applies force toward the origin (attract, implosion, magnetism)",
			Element:     "force",
		},
		{
			Name:        "Lift",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryPhysical,
			Description: "Negates gravity or applies upward force (levitate)",
			Element:     "force",
		},
		{
			Name:        "Crush",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryPhysical,
			Description: "Increases gravity or applies downward force",
			Element:     "force",
		},
		{
			Name:        "Spin",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryPhysical,
			Description: "Adds rotational force (creates tornadoes or drills)",
			Element:     "force",
		},
		{
			Name:        "Pierce",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryPhysical,
			Description: "Ignores collision/armor for a set distance (penetrate)",
			Element:     "force",
		},
		{
			Name:        "Scatter",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryPhysical,
			Description: "Spreads force over a wide area (disperse, shotgun effect)",
			Element:     "force",
		},
		{
			Name:        "Focus",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryPhysical,
			Description: "Condenses force into a single point (concentrate, precision)",
			Element:     "force",
		},
		{
			Name:        "Homing",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryPhysical,
			Description: "Causes projectiles to track and follow targets",
			Element:     "force",
		},

		// =============================================================================
		// III. THERMODYNAMIC & STATE MODIFIERS (The "Adjectives")
		// Change the state of matter of elements or targets
		// =============================================================================
		{
			Name:        "Heat",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryThermodynamic,
			Description: "Increases temperature (exothermic). Turns Water to Steam, Earth to Magma",
			Element:     "",
		},
		{
			Name:        "Cool",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryThermodynamic,
			Description: "Decreases temperature (endothermic). Turns Water to Ice, Air to Liquid Nitrogen",
			Element:     "",
		},
		{
			Name:        "Conduct",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryThermodynamic,
			Description: "Allows electricity/magic to flow through the material (charge)",
			Element:     "",
		},
		{
			Name:        "Insulate",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryThermodynamic,
			Description: "Prevents energy transfer (block, creates shields)",
			Element:     "",
		},

		// =============================================================================
		// IV. SPATIAL & TEMPORAL FORCES
		// Manipulate the fabric of the game world itself
		// =============================================================================
		{
			Name:        "Expand",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategorySpatial,
			Description: "Increases the scale/hitbox of the object (growth)",
			Element:     "",
		},
		{
			Name:        "Shrink",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategorySpatial,
			Description: "Decreases the scale/hitbox (diminish, good for concentrated density)",
			Element:     "",
		},
		{
			Name:        "Haste",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategorySpatial,
			Description: "Increases the simulation speed of the target (accelerate)",
			Element:     "",
		},
		{
			Name:        "Slow",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategorySpatial,
			Description: "Decreases the simulation speed of the target (decelerate)",
			Element:     "",
		},
		{
			Name:        "Teleport",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategorySpatial,
			Description: "Instant displacement of matter (blink)",
			Element:     "",
		},
		{
			Name:        "Echo",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategorySpatial,
			Description: "Causes the effect to trigger multiple times (repeat, DoT or multi-hit)",
			Element:     "",
		},
		{
			Name:        "Extreme",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategorySpatial,
			Description: "Pushes another component to its absolute limit (Slow to Stop, Heat to Incinerate, Cool to Absolute Zero)",
			Element:     "",
		},

		// =============================================================================
		// V. LIFE
		// Forces that interact with living things and life force
		// =============================================================================
		{
			Name:        "Mend",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryLife,
			Description: "Restores structural integrity or HP (heal)",
			Element:     "",
		},
		{
			Name:        "Wither",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryLife,
			Description: "Adds damage over time or rust/corrosion (decay)",
			Element:     "",
		},
		{
			Name:        "Revive",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryLife,
			Description: "Brings a destroyed object/entity back to a base state (resurrect)",
			Element:     "",
		},
		{
			Name:        "Drain",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryLife,
			Description: "Transfers energy/HP from target to caster (leech)",
			Element:     "",
		},

		// =============================================================================
		// VI. SPELL SHAPES (The "Container")
		// Define how the spell is delivered in the game engine
		// =============================================================================
		{
			Name:        "Projectile",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryShape,
			Description: "A standard moving object (bolt)",
			Element:     "",
		},
		{
			Name:        "Beam",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryShape,
			Description: "Instant line-of-sight connection (ray)",
			Element:     "",
		},
		{
			Name:        "Nova",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryShape,
			Description: "Radial explosion from a point (area of effect)",
			Element:     "",
		},
		{
			Name:        "Wall",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryShape,
			Description: "A stationary vertical plane (barrier)",
			Element:     "",
		},
		{
			Name:        "Zone",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryShape,
			Description: "A lingering area on the ground (field, puddles/clouds)",
			Element:     "",
		},
		{
			Name:        "Self",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryShape,
			Description: "Applies directly to the caster (buff/cloak)",
			Element:     "",
		},
		{
			Name:        "Trap",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryShape,
			Description: "Stationary object that triggers on contact (rune)",
			Element:     "",
		},
		{
			Name:        "Cone",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryShape,
			Description: "A cone-shaped blast originating from the caster (breath weapon)",
			Element:     "",
		},

		// =============================================================================
		// VII. ABJURATION (The School of Protection)
		// Focuses on blocking, banishing, and protecting
		// =============================================================================
		{
			Name:        "Reflect",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryAbjuration,
			Description: "Bounces incoming projectiles back at the source (mirror)",
			Element:     "",
		},
		{
			Name:        "Dispel",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryAbjuration,
			Description: "Removes active status effects or destroys enemy traps/zones (nullify)",
			Element:     "",
		},
		{
			Name:        "Absorb",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryAbjuration,
			Description: "Converts incoming damage into Mana for the player",
			Element:     "",
		},
		{
			Name:        "Seal",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryAbjuration,
			Description: "Permanently closes a door, chest, or portal (lock)",
			Element:     "",
		},

		// =============================================================================
		// VIII. CONJURATION (The School of Summoning)
		// Focuses on transporting objects and creating creatures
		// =============================================================================
		{
			Name:        "Summon",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryConjuration,
			Description: "Spawns a friendly entity (AI companion, minion)",
			Element:     "",
		},
		{
			Name:        "Portal",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryConjuration,
			Description: "Creates two linked points; walking in one exits the other (gate)",
			Element:     "",
		},
		{
			Name:        "Grease",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryConjuration,
			Description: "Covers the floor in a zero-friction substance (slick, flammable based on other components)",
			Element:     "",
		},

		// =============================================================================
		// IX. DIVINATION (The School of Information)
		// Focuses on revealing secrets and predicting the future
		// =============================================================================
		{
			Name:        "Sight",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryDivination,
			Description: "Allows seeing enemies/loot through walls (x-ray)",
			Element:     "",
		},
		{
			Name:        "Identify",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryDivination,
			Description: "Reveals the HP, weaknesses, and stats of a target (lore)",
			Element:     "",
		},
		{
			Name:        "Predict",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryDivination,
			Description: "Shows the path of enemy projectiles before they fire (trajectory)",
			Element:     "",
		},
		{
			Name:        "Link",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryDivination,
			Description: "Connects two targets; damage dealt to one is felt by the other (share)",
			Element:     "",
		},

		// =============================================================================
		// X. ENCHANTMENT (The School of Influence)
		// Focuses on affecting the minds of others
		// =============================================================================
		{
			Name:        "Fear",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryEnchantment,
			Description: "Forces enemies to run directly away from the caster (panic)",
			Element:     "",
		},
		{
			Name:        "Taunt",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryEnchantment,
			Description: "Forces enemies to attack a specific target (aggro, good for decoys)",
			Element:     "",
		},
		{
			Name:        "Frenzy",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryEnchantment,
			Description: "Increases enemy damage but causes them to attack their own allies (berserk)",
			Element:     "",
		},
		{
			Name:        "Command",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryEnchantment,
			Description: "A generic trigger component that forces actions (obey, e.g., Command + Jump)",
			Element:     "",
		},

		// =============================================================================
		// XI. EVOCATION (The School of Energy)
		// Focuses on raw destructive power - energy shapes and delivery
		// =============================================================================
		{
			Name:        "Channel",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryEvocation,
			Description: "A continuous flow of energy like a flamethrower that drains mana over time (stream)",
			Element:     "",
		},
		{
			Name:        "Chain",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryEvocation,
			Description: "On hit, the spell jumps to a nearby enemy (arc)",
			Element:     "",
		},
		{
			Name:        "Imbue",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryEvocation,
			Description: "Instead of firing a spell, adds the element to your sword/melee attack (enchant weapon)",
			Element:     "",
		},
		{
			Name:        "Rupture",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryEvocation,
			Description: "A delayed explosion that triggers seconds after the initial hit (aftershock)",
			Element:     "",
		},

		// =============================================================================
		// XII. ILLUSION (The School of Deception)
		// Focuses on tricking the senses
		// =============================================================================
		{
			Name:        "Decoy",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryIllusion,
			Description: "Creates a fake copy of the player that attracts attacks but deals no damage (clone)",
			Element:     "",
		},
		{
			Name:        "Invisible",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryIllusion,
			Description: "Makes the target unseeable by AI (stealth)",
			Element:     "",
		},
		{
			Name:        "Silence",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryIllusion,
			Description: "Prevents the target from casting spells or calling for help (mute)",
			Element:     "",
		},
		{
			Name:        "Disguise",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryIllusion,
			Description: "Makes the player look like an enemy faction (mask, prevents aggro)",
			Element:     "",
		},
		{
			Name:        "Phantom",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryIllusion,
			Description: "Creates a visual spell effect that does no damage but triggers Fear or Dodge AI (fake)",
			Element:     "",
		},

		// =============================================================================
		// XIII. NECROMANCY (The School of Life & Death)
		// Focuses on manipulating the soul and the body
		// =============================================================================
		{
			Name:        "Bone",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryNecromancy,
			Description: "Uses physical remains to build walls or armor (structure)",
			Element:     "",
		},
		{
			Name:        "Soul",
			Type:        models.ComponentTypeBase,
			Category:    models.CategoryNecromancy,
			Description: "A component required for targeting ghosts or interacting with the dead (spirit)",
			Element:     "",
		},
		{
			Name:        "Curse",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryNecromancy,
			Description: "A debuff that lowers enemy stats (Attack/Defense) without doing damage (hex)",
			Element:     "",
		},
		{
			Name:        "Sacrifice",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryNecromancy,
			Description: "Casts a spell using HP instead of Mana (blood)",
			Element:     "",
		},

		// =============================================================================
		// XIV. TRANSMUTATION (The School of Change)
		// Focuses on altering physical properties
		// =============================================================================
		{
			Name:        "Transmute",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryTransmutation,
			Description: "The 'Make X' command. Turn Flesh to Stone, Wood to Gold (change)",
			Element:     "",
		},
		{
			Name:        "Burrow",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryTransmutation,
			Description: "Allows movement through the ground/terrain (phase)",
			Element:     "",
		},
		{
			Name:        "Fuse",
			Type:        models.ComponentTypeModifier,
			Category:    models.CategoryTransmutation,
			Description: "Permanently attaches two objects together physically (combine)",
			Element:     "",
		},
	}
}
