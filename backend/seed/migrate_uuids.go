package seed

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/seed/uuids"
	"gorm.io/gorm"
)

// MigrateToStableUUIDs migrates existing entities from random UUIDs to deterministic UUIDs.
// This is a ONE-TIME migration that should be run before switching to the new seeding system.
// It updates all foreign key references (characters, etc.) to point to the new UUIDs.
//
// Strategy: Add temp column with new UUID, update FKs via join, then update ID column.
// This avoids FK constraint violations by never having orphaned references.
//
// After running this migration, future reseeds will preserve the same UUIDs.
func MigrateToStableUUIDs(db *gorm.DB) error {
	log.Println("Migrating to stable/deterministic UUIDs...")

	// Migrate components first (no FK dependencies on other seeded tables)
	if err := migrateComponentUUIDs(db); err != nil {
		return fmt.Errorf("components: %w", err)
	}

	// Migrate races
	if err := migrateRaceUUIDs(db); err != nil {
		return fmt.Errorf("races: %w", err)
	}

	// Migrate classes (before archetypes since archetypes depend on classes)
	if err := migrateClassUUIDs(db); err != nil {
		return fmt.Errorf("classes: %w", err)
	}

	// Migrate archetypes (after classes)
	if err := migrateArchetypeUUIDs(db); err != nil {
		return fmt.Errorf("archetypes: %w", err)
	}

	// Migrate weapons
	if err := migrateWeaponUUIDs(db); err != nil {
		return fmt.Errorf("weapons: %w", err)
	}

	log.Println("UUID migration completed successfully")
	return nil
}

func migrateRaceUUIDs(db *gorm.DB) error {
	var races []models.Race
	if err := db.Find(&races).Error; err != nil {
		return err
	}

	// Check if any need migration
	needsMigration := false
	for _, race := range races {
		if race.ID != uuids.RaceUUID(race.Name) {
			needsMigration = true
			break
		}
	}
	if !needsMigration {
		log.Println("Races already have stable UUIDs")
		return nil
	}

	log.Println("Migrating race UUIDs...")

	// Build old->new ID mapping
	idMap := make(map[uuid.UUID]uuid.UUID)
	for _, race := range races {
		idMap[race.ID] = uuids.RaceUUID(race.Name)
		log.Printf("  Race %s: %s -> %s", race.Name, race.ID, idMap[race.ID])
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Drop ALL foreign key constraints referencing races
		if err := dropFKsReferencing(tx, "races"); err != nil {
			return err
		}

		// Update child table FKs first (while parent still has old IDs)
		for oldID, newID := range idMap {
			if err := tx.Exec(`UPDATE traits SET race_id = ? WHERE race_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating traits for race %s: %w", oldID, err)
			}
			if err := tx.Exec(`UPDATE lineages SET race_id = ? WHERE race_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating lineages for race %s: %w", oldID, err)
			}
			if err := tx.Exec(`UPDATE race_components SET race_id = ? WHERE race_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating race_components for race %s: %w", oldID, err)
			}
			if err := tx.Exec(`UPDATE characters SET race_id = ? WHERE race_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating characters for race %s: %w", oldID, err)
			}
		}

		// Now update the parent table IDs
		for oldID, newID := range idMap {
			if err := tx.Exec(`UPDATE races SET id = ? WHERE id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating race %s: %w", oldID, err)
			}
		}

		// Recreate FK constraints
		if err := tx.Exec(`ALTER TABLE traits ADD CONSTRAINT fk_races_traits FOREIGN KEY (race_id) REFERENCES races(id) ON DELETE CASCADE`).Error; err != nil {
			return fmt.Errorf("recreating fk_races_traits: %w", err)
		}
		if err := tx.Exec(`ALTER TABLE lineages ADD CONSTRAINT fk_races_lineages FOREIGN KEY (race_id) REFERENCES races(id) ON DELETE CASCADE`).Error; err != nil {
			return fmt.Errorf("recreating fk_races_lineages: %w", err)
		}
		if err := tx.Exec(`ALTER TABLE race_components ADD CONSTRAINT fk_race_components_race FOREIGN KEY (race_id) REFERENCES races(id) ON DELETE CASCADE`).Error; err != nil {
			return fmt.Errorf("recreating fk_race_components_race: %w", err)
		}
		if err := tx.Exec(`ALTER TABLE characters ADD CONSTRAINT fk_characters_race FOREIGN KEY (race_id) REFERENCES races(id)`).Error; err != nil {
			return fmt.Errorf("recreating fk_characters_race: %w", err)
		}

		log.Printf("Migrated %d race UUIDs", len(races))
		return nil
	})
}

// dropFKsReferencing drops all foreign key constraints that reference the given table
func dropFKsReferencing(tx *gorm.DB, tableName string) error {
	// Query to find all FK constraints referencing this table
	var constraints []struct {
		ConstraintName string
		TableName      string
	}
	if err := tx.Raw(`
		SELECT tc.constraint_name, tc.table_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.constraint_column_usage ccu 
			ON tc.constraint_name = ccu.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY' 
			AND ccu.table_name = ?
	`, tableName).Scan(&constraints).Error; err != nil {
		return fmt.Errorf("querying FK constraints for %s: %w", tableName, err)
	}

	log.Printf("  Found %d FK constraints referencing %s", len(constraints), tableName)
	for _, c := range constraints {
		log.Printf("  Dropping FK constraint %s on %s", c.ConstraintName, c.TableName)
		if err := tx.Exec(fmt.Sprintf(`ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s`, c.TableName, c.ConstraintName)).Error; err != nil {
			return fmt.Errorf("dropping constraint %s on %s: %w", c.ConstraintName, c.TableName, err)
		}
	}
	return nil
}

func migrateClassUUIDs(db *gorm.DB) error {
	var classes []models.Class
	if err := db.Find(&classes).Error; err != nil {
		return err
	}

	needsMigration := false
	for _, class := range classes {
		if class.ID != uuids.ClassUUID(class.Name) {
			needsMigration = true
			break
		}
	}
	if !needsMigration {
		log.Println("Classes already have stable UUIDs")
		return nil
	}

	log.Println("Migrating class UUIDs...")

	// Build old->new ID mapping
	idMap := make(map[uuid.UUID]uuid.UUID)
	for _, class := range classes {
		idMap[class.ID] = uuids.ClassUUID(class.Name)
		log.Printf("  Class %s: %s -> %s", class.Name, class.ID, idMap[class.ID])
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Drop ALL foreign key constraints referencing classes
		if err := dropFKsReferencing(tx, "classes"); err != nil {
			return err
		}

		// Update child table FKs first
		for oldID, newID := range idMap {
			if err := tx.Exec(`UPDATE class_levels SET class_id = ? WHERE class_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating class_levels for class %s: %w", oldID, err)
			}
			if err := tx.Exec(`UPDATE archetypes SET class_id = ? WHERE class_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating archetypes for class %s: %w", oldID, err)
			}
			if err := tx.Exec(`UPDATE class_components SET class_id = ? WHERE class_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating class_components for class %s: %w", oldID, err)
			}
			if err := tx.Exec(`UPDATE characters SET class_id = ? WHERE class_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating characters for class %s: %w", oldID, err)
			}
		}

		// Now update the parent table IDs
		for oldID, newID := range idMap {
			if err := tx.Exec(`UPDATE classes SET id = ? WHERE id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating class %s: %w", oldID, err)
			}
		}

		// Recreate FK constraints
		if err := tx.Exec(`ALTER TABLE class_levels ADD CONSTRAINT fk_classes_levels FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE`).Error; err != nil {
			return fmt.Errorf("recreating fk_classes_levels: %w", err)
		}
		if err := tx.Exec(`ALTER TABLE archetypes ADD CONSTRAINT fk_classes_archetypes FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE`).Error; err != nil {
			return fmt.Errorf("recreating fk_classes_archetypes: %w", err)
		}
		if err := tx.Exec(`ALTER TABLE class_components ADD CONSTRAINT fk_class_components_class FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE`).Error; err != nil {
			return fmt.Errorf("recreating fk_class_components_class: %w", err)
		}
		if err := tx.Exec(`ALTER TABLE characters ADD CONSTRAINT fk_characters_class FOREIGN KEY (class_id) REFERENCES classes(id)`).Error; err != nil {
			return fmt.Errorf("recreating fk_characters_class: %w", err)
		}

		log.Printf("Migrated %d class UUIDs", len(classes))
		return nil
	})
}

func migrateArchetypeUUIDs(db *gorm.DB) error {
	var archetypes []models.Archetype
	if err := db.Preload("Class").Find(&archetypes).Error; err != nil {
		return err
	}

	if len(archetypes) == 0 {
		log.Println("No archetypes to migrate")
		return nil
	}

	needsMigration := false
	for _, archetype := range archetypes {
		if archetype.Class.ID == uuid.Nil {
			continue
		}
		if archetype.ID != uuids.ArchetypeUUID(archetype.Class.Name, archetype.Name) {
			needsMigration = true
			break
		}
	}
	if !needsMigration {
		log.Println("Archetypes already have stable UUIDs")
		return nil
	}

	log.Println("Migrating archetype UUIDs...")

	// Build old->new ID mapping
	idMap := make(map[uuid.UUID]uuid.UUID)
	for _, archetype := range archetypes {
		if archetype.Class.ID == uuid.Nil {
			continue
		}
		idMap[archetype.ID] = uuids.ArchetypeUUID(archetype.Class.Name, archetype.Name)
		log.Printf("  Archetype %s (%s): %s -> %s", archetype.Name, archetype.Class.Name, archetype.ID, idMap[archetype.ID])
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Drop ALL foreign key constraints referencing archetypes
		if err := dropFKsReferencing(tx, "archetypes"); err != nil {
			return err
		}

		// Update child table FKs first
		for oldID, newID := range idMap {
			if err := tx.Exec(`UPDATE level_features SET archetype_id = ? WHERE archetype_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating level_features for archetype %s: %w", oldID, err)
			}
			if err := tx.Exec(`UPDATE characters SET archetype_id = ? WHERE archetype_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating characters for archetype %s: %w", oldID, err)
			}
		}

		// Now update the parent table IDs
		for oldID, newID := range idMap {
			if err := tx.Exec(`UPDATE archetypes SET id = ? WHERE id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating archetype %s: %w", oldID, err)
			}
		}

		// Recreate FK constraints
		if err := tx.Exec(`ALTER TABLE level_features ADD CONSTRAINT fk_level_features_archetype FOREIGN KEY (archetype_id) REFERENCES archetypes(id) ON DELETE SET NULL`).Error; err != nil {
			return fmt.Errorf("recreating fk_level_features_archetype: %w", err)
		}
		if err := tx.Exec(`ALTER TABLE characters ADD CONSTRAINT fk_characters_archetype FOREIGN KEY (archetype_id) REFERENCES archetypes(id)`).Error; err != nil {
			return fmt.Errorf("recreating fk_characters_archetype: %w", err)
		}

		log.Printf("Migrated %d archetype UUIDs", len(archetypes))
		return nil
	})
}

func migrateComponentUUIDs(db *gorm.DB) error {
	var components []models.Component
	if err := db.Find(&components).Error; err != nil {
		return err
	}

	if len(components) == 0 {
		log.Println("No components to migrate")
		return nil
	}

	needsMigration := false
	for _, component := range components {
		if component.ID != uuids.ComponentUUID(component.Name) {
			needsMigration = true
			break
		}
	}
	if !needsMigration {
		log.Println("Components already have stable UUIDs")
		return nil
	}

	log.Println("Migrating component UUIDs...")

	// Build old->new ID mapping
	idMap := make(map[uuid.UUID]uuid.UUID)
	for _, component := range components {
		idMap[component.ID] = uuids.ComponentUUID(component.Name)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Drop ALL foreign key constraints referencing components
		if err := dropFKsReferencing(tx, "components"); err != nil {
			return err
		}

		// Update child table FKs first
		for oldID, newID := range idMap {
			if err := tx.Exec(`UPDATE class_components SET component_id = ? WHERE component_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating class_components for component %s: %w", oldID, err)
			}
			if err := tx.Exec(`UPDATE race_components SET component_id = ? WHERE component_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating race_components for component %s: %w", oldID, err)
			}
			if err := tx.Exec(`UPDATE spell_components SET component_id = ? WHERE component_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating spell_components for component %s: %w", oldID, err)
			}
		}

		// Now update the parent table IDs
		for oldID, newID := range idMap {
			if err := tx.Exec(`UPDATE components SET id = ? WHERE id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating component %s: %w", oldID, err)
			}
		}

		// Recreate FK constraints
		if err := tx.Exec(`ALTER TABLE class_components ADD CONSTRAINT fk_class_components_component FOREIGN KEY (component_id) REFERENCES components(id) ON DELETE CASCADE`).Error; err != nil {
			return fmt.Errorf("recreating fk_class_components_component: %w", err)
		}
		if err := tx.Exec(`ALTER TABLE race_components ADD CONSTRAINT fk_race_components_component FOREIGN KEY (component_id) REFERENCES components(id) ON DELETE CASCADE`).Error; err != nil {
			return fmt.Errorf("recreating fk_race_components_component: %w", err)
		}
		if err := tx.Exec(`ALTER TABLE spell_components ADD CONSTRAINT fk_spell_components_component FOREIGN KEY (component_id) REFERENCES components(id) ON DELETE CASCADE`).Error; err != nil {
			return fmt.Errorf("recreating fk_spell_components_component: %w", err)
		}

		log.Printf("Migrated %d component UUIDs", len(components))
		return nil
	})
}

func migrateWeaponUUIDs(db *gorm.DB) error {
	var weapons []models.Weapon
	// Only migrate system weapons (UserID is null)
	if err := db.Where("user_id IS NULL").Find(&weapons).Error; err != nil {
		return err
	}

	if len(weapons) == 0 {
		log.Println("No weapons to migrate")
		return nil
	}

	needsMigration := false
	for _, weapon := range weapons {
		if weapon.ID != uuids.WeaponUUID(weapon.Name) {
			needsMigration = true
			break
		}
	}
	if !needsMigration {
		log.Println("Weapons already have stable UUIDs")
		return nil
	}

	log.Println("Migrating weapon UUIDs...")

	// Build old->new ID mapping
	idMap := make(map[uuid.UUID]uuid.UUID)
	for _, weapon := range weapons {
		idMap[weapon.ID] = uuids.WeaponUUID(weapon.Name)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Drop ALL foreign key constraints referencing weapons
		if err := dropFKsReferencing(tx, "weapons"); err != nil {
			return err
		}

		// Update child table FKs first
		for oldID, newID := range idMap {
			if err := tx.Exec(`UPDATE weapon_damages SET weapon_id = ? WHERE weapon_id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating weapon_damages for weapon %s: %w", oldID, err)
			}
		}

		// Now update the parent table IDs
		for oldID, newID := range idMap {
			if err := tx.Exec(`UPDATE weapons SET id = ? WHERE id = ?`, newID, oldID).Error; err != nil {
				return fmt.Errorf("updating weapon %s: %w", oldID, err)
			}
		}

		// Recreate FK constraints
		if err := tx.Exec(`ALTER TABLE weapon_damages ADD CONSTRAINT fk_weapons_damages FOREIGN KEY (weapon_id) REFERENCES weapons(id) ON DELETE CASCADE`).Error; err != nil {
			return fmt.Errorf("recreating fk_weapons_damages: %w", err)
		}

		log.Printf("Migrated %d weapon UUIDs", len(weapons))
		return nil
	})
}

// NeedsUUIDMigration checks if any entities still have non-deterministic UUIDs.
func NeedsUUIDMigration(db *gorm.DB) bool {
	// Check races
	var race models.Race
	if err := db.First(&race).Error; err == nil {
		if race.ID != uuids.RaceUUID(race.Name) {
			log.Println("Race UUIDs need migration")
			return true
		}
	}

	// Check classes
	var class models.Class
	if err := db.First(&class).Error; err == nil {
		if class.ID != uuids.ClassUUID(class.Name) {
			log.Println("Class UUIDs need migration")
			return true
		}
	}

	// Check components
	var component models.Component
	if err := db.First(&component).Error; err == nil {
		if component.ID != uuids.ComponentUUID(component.Name) {
			log.Println("Component UUIDs need migration")
			return true
		}
	}

	// Check archetypes
	var archetype models.Archetype
	if err := db.Preload("Class").First(&archetype).Error; err == nil {
		if archetype.Class.ID != uuid.Nil && archetype.ID != uuids.ArchetypeUUID(archetype.Class.Name, archetype.Name) {
			log.Println("Archetype UUIDs need migration")
			return true
		}
	}

	// Check weapons (only system weapons)
	var weapon models.Weapon
	if err := db.Where("user_id IS NULL").First(&weapon).Error; err == nil {
		if weapon.ID != uuids.WeaponUUID(weapon.Name) {
			log.Println("Weapon UUIDs need migration")
			return true
		}
	}

	log.Println("All UUIDs already stable")
	return false
}

// Helper to check if a UUID is the expected deterministic one
func isExpectedUUID(actual, expected uuid.UUID) bool {
	return actual == expected
}
