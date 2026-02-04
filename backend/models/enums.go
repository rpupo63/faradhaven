package models

// =============================================================================
// CREATURE ENUMS
// =============================================================================

// CreatureSize represents creature size categories (D&D 5e)
type CreatureSize string

const (
	SizeTiny       CreatureSize = "Tiny"
	SizeSmall      CreatureSize = "Small"
	SizeMedium     CreatureSize = "Medium"
	SizeLarge      CreatureSize = "Large"
	SizeHuge       CreatureSize = "Huge"
	SizeGargantuan CreatureSize = "Gargantuan"
)

// CreatureType represents creature type categories (D&D 5e)
type CreatureType string

const (
	TypeAberration  CreatureType = "Aberration"
	TypeBeast       CreatureType = "Beast"
	TypeCelestial   CreatureType = "Celestial"
	TypeConstruct   CreatureType = "Construct"
	TypeDragon      CreatureType = "Dragon"
	TypeElemental   CreatureType = "Elemental"
	TypeFey         CreatureType = "Fey"
	TypeFiend       CreatureType = "Fiend"
	TypeGiant       CreatureType = "Giant"
	TypeHumanoid    CreatureType = "Humanoid"
	TypeMonstrosity CreatureType = "Monstrosity"
	TypeOoze        CreatureType = "Ooze"
	TypePlant       CreatureType = "Plant"
	TypeUndead      CreatureType = "Undead"
)

// =============================================================================
// DAMAGE ENUMS
// =============================================================================

// DamageType represents damage type categories (D&D 5e)
type DamageType string

const (
	DamageSlashing    DamageType = "Slashing"
	DamagePiercing    DamageType = "Piercing"
	DamageBludgeoning DamageType = "Bludgeoning"
	DamageFire        DamageType = "Fire"
	DamageCold        DamageType = "Cold"
	DamageLightning   DamageType = "Lightning"
	DamageThunder     DamageType = "Thunder"
	DamagePoison      DamageType = "Poison"
	DamageAcid        DamageType = "Acid"
	DamageNecrotic    DamageType = "Necrotic"
	DamageRadiant     DamageType = "Radiant"
	DamageForce       DamageType = "Force"
	DamagePsychic     DamageType = "Psychic"
)

// =============================================================================
// COMPONENT ENUMS
// =============================================================================

// ComponentType indicates whether a component is base (core to the class) or modifier (enhancement)
type ComponentType string

const (
	ComponentTypeBase     ComponentType = "base"
	ComponentTypeModifier ComponentType = "modifier"
)

// ComponentCategory represents the magical school or category of a spell component
type ComponentCategory string

const (
	// Primordial Elements - base materials that give spells substance
	CategoryPrimordial ComponentCategory = "primordial"

	// Physical Forces - dictate how spells move or interact with physics
	CategoryPhysical ComponentCategory = "physical"

	// Thermodynamic & State Modifiers - change state of matter
	CategoryThermodynamic ComponentCategory = "thermodynamic"

	// Spatial & Temporal Forces - manipulate the fabric of the world
	CategorySpatial ComponentCategory = "spatial"

	// Life - interact with living things, life force, and vitality
	CategoryLife ComponentCategory = "life"

	// Spell Shapes - define how spells are delivered
	CategoryShape ComponentCategory = "shape"

	// Abjuration - the school of protection (blocking, banishing, protecting)
	CategoryAbjuration ComponentCategory = "abjuration"

	// Conjuration - the school of summoning (transporting objects, creating creatures)
	CategoryConjuration ComponentCategory = "conjuration"

	// Divination - the school of information (revealing secrets, predicting)
	CategoryDivination ComponentCategory = "divination"

	// Enchantment - the school of influence (affecting minds)
	CategoryEnchantment ComponentCategory = "enchantment"

	// Evocation - the school of energy (raw destructive power)
	CategoryEvocation ComponentCategory = "evocation"

	// Illusion - the school of deception (tricking the senses)
	CategoryIllusion ComponentCategory = "illusion"

	// Necromancy - the school of life & death (manipulating soul and body)
	CategoryNecromancy ComponentCategory = "necromancy"

	// Transmutation - the school of change (altering physical properties)
	CategoryTransmutation ComponentCategory = "transmutation"
)
