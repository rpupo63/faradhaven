package uuids

import (
	"fmt"

	"github.com/google/uuid"
)

// Namespace UUIDs for generating deterministic UUIDs for each entity type.
// These are fixed UUIDs that serve as namespaces for UUID v5 generation.
// Using UUID v5 means the same name always produces the same UUID.
var (
	NamespaceRace         = uuid.MustParse("a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d")
	NamespaceClass        = uuid.MustParse("b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e")
	NamespaceComponent    = uuid.MustParse("c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f")
	NamespaceArchetype    = uuid.MustParse("d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a")
	NamespaceWeapon       = uuid.MustParse("e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b")
	NamespaceLineage      = uuid.MustParse("f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c")
	NamespaceChoice       = uuid.MustParse("01234567-89ab-cdef-0123-456789abcdef")
	NamespaceOption       = uuid.MustParse("fedcba98-7654-3210-fedc-ba9876543210")
	NamespaceItem         = uuid.MustParse("12345678-1234-5678-1234-567812345678")
	NamespaceClassLevel   = uuid.MustParse("23456789-2345-6789-2345-678923456789")
	NamespaceLevelFeature = uuid.MustParse("34567890-3456-7890-3456-789034567890")
	NamespaceTrait        = uuid.MustParse("45678901-4567-8901-4567-890145678901")
	NamespaceTraitOption  = uuid.MustParse("56789012-5678-9012-5678-901256789012")
	NamespaceWeaponDamage       = uuid.MustParse("67890123-6789-0123-6789-012367890123")
	NamespaceWeaponRequirement  = uuid.MustParse("78901234-7890-1234-7890-123478901234")
	NamespaceEffect               = uuid.MustParse("89012345-8901-2345-8901-234589012345")
	NamespaceClassResourceDef     = uuid.MustParse("9a0b1c2d-3e4f-5a6b-7c8d-9e0f1a2b3c4d")
	NamespaceClassLevelResource   = uuid.MustParse("ab1c2d3e-4f5a-6b7c-8d9e-0f1a2b3c4d5e")
	NamespaceSystemUser         = uuid.MustParse("bc2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f")
)

// RaceUUID generates a deterministic UUID for a race by name.
func RaceUUID(name string) uuid.UUID {
	return uuid.NewSHA1(NamespaceRace, []byte(name))
}

// SystemUserUUID generates a deterministic UUID for a system user by name.
func SystemUserUUID(name string) uuid.UUID {
	return uuid.NewSHA1(NamespaceSystemUser, []byte(name))
}

// ClassUUID generates a deterministic UUID for a class by name.
func ClassUUID(name string) uuid.UUID {
	return uuid.NewSHA1(NamespaceClass, []byte(name))
}

// ComponentUUID generates a deterministic UUID for a component by name.
func ComponentUUID(name string) uuid.UUID {
	return uuid.NewSHA1(NamespaceComponent, []byte(name))
}

// ArchetypeUUID generates a deterministic UUID for an archetype by class name + archetype name.
func ArchetypeUUID(className, archetypeName string) uuid.UUID {
	return uuid.NewSHA1(NamespaceArchetype, []byte(className+":"+archetypeName))
}

// WeaponUUID generates a deterministic UUID for a weapon by name.
func WeaponUUID(name string) uuid.UUID {
	return uuid.NewSHA1(NamespaceWeapon, []byte(name))
}

// ItemUUID generates a deterministic UUID for an item by name.
func ItemUUID(name string) uuid.UUID {
	return uuid.NewSHA1(NamespaceItem, []byte(name))
}

// LineageUUID generates a deterministic UUID for a lineage by race name + lineage name.
func LineageUUID(raceName, lineageName string) uuid.UUID {
	return uuid.NewSHA1(NamespaceLineage, []byte(raceName+":"+lineageName))
}

// EquipmentChoiceUUID generates a deterministic UUID for an equipment choice.
// Key: ClassName + Instruction
func EquipmentChoiceUUID(className, instruction string) uuid.UUID {
	return uuid.NewSHA1(NamespaceChoice, []byte(className+":"+instruction))
}

// EquipmentOptionUUID generates a deterministic UUID for an equipment option.
// Key: ChoiceUUID string + Description
func EquipmentOptionUUID(choiceID uuid.UUID, description string) uuid.UUID {
	return uuid.NewSHA1(NamespaceOption, []byte(choiceID.String()+":"+description))
}

// ClassLevelUUID generates a deterministic UUID for a class level.
// Key: ClassName + Level
func ClassLevelUUID(className string, level int) uuid.UUID {
	return uuid.NewSHA1(NamespaceClassLevel, []byte(fmt.Sprintf("%s:%d", className, level)))
}

// LevelFeatureUUID generates a deterministic UUID for a level feature.
// Key: ClassLevelID + FeatureName + SortOrder (to handle duplicates)
func LevelFeatureUUID(classLevelID uuid.UUID, featureName string, sortOrder int) uuid.UUID {
	return uuid.NewSHA1(NamespaceLevelFeature, []byte(fmt.Sprintf("%s:%s:%d", classLevelID.String(), featureName, sortOrder)))
}

// TraitUUID generates a deterministic UUID for a race trait.
// Key: RaceName + TraitName
func TraitUUID(raceName, traitName string) uuid.UUID {
	return uuid.NewSHA1(NamespaceTrait, []byte(raceName+":"+traitName))
}

// TraitOptionUUID generates a deterministic UUID for a trait option.
// Key: TraitID + OptionName
func TraitOptionUUID(traitID uuid.UUID, optionName string) uuid.UUID {
	return uuid.NewSHA1(NamespaceTraitOption, []byte(traitID.String()+":"+optionName))
}

// WeaponDamageUUID generates a deterministic UUID for a weapon damage entry.
// Key: WeaponID + DamageDice + DamageType + DamageCategory
func WeaponDamageUUID(weaponID uuid.UUID, damageDice, damageType, damageCategory string) uuid.UUID {
	return uuid.NewSHA1(NamespaceWeaponDamage, []byte(fmt.Sprintf("%s:%s:%s:%s", weaponID.String(), damageDice, damageType, damageCategory)))
}

// WeaponRequirementUUID generates a deterministic UUID for a class weapon requirement.
// Key: ClassName + ModifierType + SelectionLevel
func WeaponRequirementUUID(className, modifierType string, selectionLevel int) uuid.UUID {
	return uuid.NewSHA1(NamespaceWeaponRequirement, []byte(fmt.Sprintf("%s:%s:%d", className, modifierType, selectionLevel)))
}

// EffectUUID generates a deterministic UUID for an effect by name.
func EffectUUID(name string) uuid.UUID {
	return uuid.NewSHA1(NamespaceEffect, []byte(name))
}

// ClassResourceDefUUID generates a deterministic UUID for a class resource definition.
// Key: ClassName + ResourceKey
func ClassResourceDefUUID(className, resourceKey string) uuid.UUID {
	return uuid.NewSHA1(NamespaceClassResourceDef, []byte(className+":"+resourceKey))
}

// ClassLevelResourceUUID generates a deterministic UUID for a class level resource value.
// Key: ClassLevelID + ResourceKey
func ClassLevelResourceUUID(classLevelID uuid.UUID, resourceKey string) uuid.UUID {
	return uuid.NewSHA1(NamespaceClassLevelResource, []byte(classLevelID.String()+":"+resourceKey))
}
