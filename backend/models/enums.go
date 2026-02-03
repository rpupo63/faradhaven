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
