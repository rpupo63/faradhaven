package uuids

import (
	"github.com/google/uuid"
)

// Namespace UUIDs for generating deterministic UUIDs for each entity type.
// These are fixed UUIDs that serve as namespaces for UUID v5 generation.
// Using UUID v5 means the same name always produces the same UUID.
var (
	NamespaceRace      = uuid.MustParse("a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d")
	NamespaceClass     = uuid.MustParse("b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e")
	NamespaceComponent = uuid.MustParse("c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f")
	NamespaceArchetype = uuid.MustParse("d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a")
	NamespaceWeapon    = uuid.MustParse("e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b")
	NamespaceLineage   = uuid.MustParse("f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c")
	NamespaceChoice    = uuid.MustParse("01234567-89ab-cdef-0123-456789abcdef")
	NamespaceOption    = uuid.MustParse("fedcba98-7654-3210-fedc-ba9876543210")
	NamespaceItem      = uuid.MustParse("12345678-1234-5678-1234-567812345678")
)

// RaceUUID generates a deterministic UUID for a race by name.
func RaceUUID(name string) uuid.UUID {
	return uuid.NewSHA1(NamespaceRace, []byte(name))
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
