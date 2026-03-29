package models

import "strings"

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

// IsValid reports whether s is one of the defined creature sizes.
func (s CreatureSize) IsValid() bool {
	switch s {
	case SizeTiny, SizeSmall, SizeMedium, SizeLarge, SizeHuge, SizeGargantuan:
		return true
	default:
		return false
	}
}

// ParseCreatureSize normalizes text (e.g. "medium", "MEDIUM") to a [CreatureSize].
func ParseCreatureSize(s string) (CreatureSize, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	lower := strings.ToLower(s)
	switch lower {
	case "tiny":
		return SizeTiny, true
	case "small":
		return SizeSmall, true
	case "medium":
		return SizeMedium, true
	case "large":
		return SizeLarge, true
	case "huge":
		return SizeHuge, true
	case "gargantuan":
		return SizeGargantuan, true
	}
	cs := CreatureSize(s)
	if cs.IsValid() {
		return cs, true
	}
	return "", false
}

// IsValid reports whether t is one of the defined creature types.
func (t CreatureType) IsValid() bool {
	switch t {
	case TypeAberration, TypeBeast, TypeCelestial, TypeConstruct, TypeDragon,
		TypeElemental, TypeFey, TypeFiend, TypeGiant, TypeHumanoid,
		TypeMonstrosity, TypeOoze, TypePlant, TypeUndead:
		return true
	default:
		return false
	}
}

// ParseCreatureType normalizes text (e.g. "humanoid", "Humanoid") to a [CreatureType].
func ParseCreatureType(s string) (CreatureType, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	lower := strings.ToLower(s)
	switch lower {
	case "aberration":
		return TypeAberration, true
	case "beast":
		return TypeBeast, true
	case "celestial":
		return TypeCelestial, true
	case "construct":
		return TypeConstruct, true
	case "dragon":
		return TypeDragon, true
	case "elemental":
		return TypeElemental, true
	case "fey":
		return TypeFey, true
	case "fiend":
		return TypeFiend, true
	case "giant":
		return TypeGiant, true
	case "humanoid":
		return TypeHumanoid, true
	case "monstrosity":
		return TypeMonstrosity, true
	case "ooze":
		return TypeOoze, true
	case "plant":
		return TypePlant, true
	case "undead":
		return TypeUndead, true
	}
	ct := CreatureType(s)
	if ct.IsValid() {
		return ct, true
	}
	return "", false
}

// =============================================================================
// MAP TOKEN ENUMS
// =============================================================================

// MapTokenType is the token category on the battle map (PC vs NPC).
type MapTokenType string

const (
	MapTokenPC  MapTokenType = "pc"
	MapTokenNPC MapTokenType = "npc"
)

// IsValid reports whether t is a valid map token type.
func (t MapTokenType) IsValid() bool {
	switch t {
	case MapTokenPC, MapTokenNPC:
		return true
	default:
		return false
	}
}

// ParseMapTokenType normalizes token text to [MapTokenType].
func ParseMapTokenType(s string) (MapTokenType, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "pc":
		return MapTokenPC, true
	case "npc":
		return MapTokenNPC, true
	default:
		return "", false
	}
}

// =============================================================================
// CLASS RESOURCE ENUMS
// =============================================================================

// ClassResourceCategory matches frontend ApiClassResource.category.
type ClassResourceCategory string

const (
	ClassResourceCategoryPool      ClassResourceCategory = "pool"
	ClassResourceCategoryDieSize     ClassResourceCategory = "die_size"
	ClassResourceCategoryLimit       ClassResourceCategory = "limit"
	ClassResourceCategorySlotCount   ClassResourceCategory = "slot_count"
	ClassResourceCategoryModifier    ClassResourceCategory = "modifier"
	ClassResourceCategoryState       ClassResourceCategory = "state"
)

// IsValid reports whether c is a valid class resource category.
func (c ClassResourceCategory) IsValid() bool {
	switch c {
	case ClassResourceCategoryPool, ClassResourceCategoryDieSize, ClassResourceCategoryLimit,
		ClassResourceCategorySlotCount, ClassResourceCategoryModifier, ClassResourceCategoryState:
		return true
	default:
		return false
	}
}

// ParseClassResourceCategory parses API/seed category strings.
func ParseClassResourceCategory(s string) (ClassResourceCategory, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "pool":
		return ClassResourceCategoryPool, true
	case "die_size":
		return ClassResourceCategoryDieSize, true
	case "limit":
		return ClassResourceCategoryLimit, true
	case "slot_count":
		return ClassResourceCategorySlotCount, true
	case "modifier":
		return ClassResourceCategoryModifier, true
	case "state":
		return ClassResourceCategoryState, true
	default:
		return "", false
	}
}

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

// IsValid reports whether d is one of the defined damage types.
func (d DamageType) IsValid() bool {
	switch d {
	case DamageSlashing, DamagePiercing, DamageBludgeoning,
		DamageFire, DamageCold, DamageLightning, DamageThunder,
		DamagePoison, DamageAcid, DamageNecrotic, DamageRadiant,
		DamageForce, DamagePsychic:
		return true
	default:
		return false
	}
}

// ParseDamageType normalizes user/LLM text (e.g. "fire", "FIRE") to a [DamageType].
func ParseDamageType(s string) (DamageType, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	// Title-case single words; multi-word types match const values exactly.
	lower := strings.ToLower(s)
	switch lower {
	case "slashing":
		return DamageSlashing, true
	case "piercing":
		return DamagePiercing, true
	case "bludgeoning":
		return DamageBludgeoning, true
	case "fire":
		return DamageFire, true
	case "cold":
		return DamageCold, true
	case "lightning":
		return DamageLightning, true
	case "thunder":
		return DamageThunder, true
	case "poison":
		return DamagePoison, true
	case "acid":
		return DamageAcid, true
	case "necrotic":
		return DamageNecrotic, true
	case "radiant":
		return DamageRadiant, true
	case "force":
		return DamageForce, true
	case "psychic":
		return DamagePsychic, true
	}
	d := DamageType(s)
	if d.IsValid() {
		return d, true
	}
	return "", false
}

// WeaponDamageCategory distinguishes base weapon damage from bonus damage (e.g. enchantments).
type WeaponDamageCategory string

const (
	WeaponDamageCategoryBase  WeaponDamageCategory = "Base"
	WeaponDamageCategoryBonus WeaponDamageCategory = "Bonus"
)

// IsValid reports whether c is a defined weapon damage category.
func (c WeaponDamageCategory) IsValid() bool {
	switch c {
	case WeaponDamageCategoryBase, WeaponDamageCategoryBonus:
		return true
	default:
		return false
	}
}

// ParseWeaponDamageCategory normalizes seed/API text to [WeaponDamageCategory].
func ParseWeaponDamageCategory(s string) (WeaponDamageCategory, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	switch strings.ToLower(s) {
	case "base":
		return WeaponDamageCategoryBase, true
	case "bonus":
		return WeaponDamageCategoryBonus, true
	}
	w := WeaponDamageCategory(s)
	if w.IsValid() {
		return w, true
	}
	return "", false
}

// =============================================================================
// SPELL ENUMS
// =============================================================================

// SpellType classifies how a crafted spell is resolved (attack roll, save, etc.).
type SpellType string

const (
	SpellTypeAttack   SpellType = "Attack"
	SpellTypeSave     SpellType = "Save"
	SpellTypeEffect   SpellType = "Effect"
	SpellTypeHealing  SpellType = "Healing"
	SpellTypeUtility  SpellType = "Utility"
)

// IsValid reports whether t is one of the allowed spell types.
func (t SpellType) IsValid() bool {
	switch t {
	case SpellTypeAttack, SpellTypeSave, SpellTypeEffect, SpellTypeHealing, SpellTypeUtility:
		return true
	default:
		return false
	}
}

// SaveAttribute is the ability used for a saving throw (D&D-style abbreviation).
type SaveAttribute string

const (
	SaveAttrSTR SaveAttribute = "STR"
	SaveAttrDEX SaveAttribute = "DEX"
	SaveAttrCON SaveAttribute = "CON"
	SaveAttrINT SaveAttribute = "INT"
	SaveAttrWIS SaveAttribute = "WIS"
	SaveAttrCHA SaveAttribute = "CHA"
)

// IsValid reports whether a is one of the six ability save codes.
func (a SaveAttribute) IsValid() bool {
	switch a {
	case SaveAttrSTR, SaveAttrDEX, SaveAttrCON, SaveAttrINT, SaveAttrWIS, SaveAttrCHA:
		return true
	default:
		return false
	}
}

// ParseSaveAttribute normalizes user/LLM text (e.g. "dex", "Dexterity") to a [SaveAttribute].
func ParseSaveAttribute(s string) (SaveAttribute, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	lower := strings.ToLower(s)
	switch lower {
	case "str", "strength":
		return SaveAttrSTR, true
	case "dex", "dexterity":
		return SaveAttrDEX, true
	case "con", "constitution":
		return SaveAttrCON, true
	case "int", "intelligence":
		return SaveAttrINT, true
	case "wis", "wisdom":
		return SaveAttrWIS, true
	case "cha", "charisma":
		return SaveAttrCHA, true
	}
	a := SaveAttribute(strings.ToUpper(s))
	if a.IsValid() {
		return a, true
	}
	return "", false
}

// =============================================================================
// COMPONENT ENUMS
// =============================================================================

// ComponentCategory represents the magical school or category of a spell component
type ComponentCategory string

const (
	CategoryForma     ComponentCategory = "Forma"     // Shape — physical manifestation/delivery
	CategoryScopus    ComponentCategory = "Scopus"    // Targeting — anchor point
	CategoryEssentia  ComponentCategory = "Essentia"  // Domain/Matter — elemental substance
	CategoryActio     ComponentCategory = "Actio"     // Kinetic Verbs — what the spell does
	CategoryMagnitudo ComponentCategory = "Magnitudo" // Scale/Modifiers — power dials
)
