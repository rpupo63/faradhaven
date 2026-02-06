package seed

import (
	"fmt"
	"log"

	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// MigrateCharacterWeaponsV2 migrates data from the legacy many2many character_weapons table
// to the new character_weapons_v2 explicit join table.
func MigrateCharacterWeaponsV2(db *gorm.DB) error {
	log.Println("Migrating character weapons to v2...")

	// 1. Ensure the new table exists
	if err := db.AutoMigrate(&models.CharacterWeapon{}); err != nil {
		return fmt.Errorf("failed to migrate CharacterWeapon model: %w", err)
	}

	// 2. Check if legacy table exists
	var exists bool
	err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'character_weapons')").Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check for legacy character_weapons table: %w", err)
	}

	if !exists {
		log.Println("Legacy character_weapons table does not exist, skipping migration.")
		return nil
	}

	// 3. Migrate data
	return db.Transaction(func(tx *gorm.DB) error {
		// Insert data from legacy table into new table if not already present
		// We use a subquery to avoid duplicates if migration is run multiple times
		err := tx.Exec(`
			INSERT INTO character_weapons_v2 (id, character_id, weapon_id, is_primary, created_at, updated_at)
			SELECT gen_random_uuid(), character_id, weapon_id, false, now(), now()
			FROM character_weapons
			WHERE NOT EXISTS (
				SELECT 1 FROM character_weapons_v2 
				WHERE character_weapons_v2.character_id = character_weapons.character_id 
				AND character_weapons_v2.weapon_id = character_weapons.weapon_id
			)
		`).Error

		if err != nil {
			return fmt.Errorf("failed to migrate data from character_weapons to character_weapons_v2: %w", err)
		}

		log.Println("Character weapons migration to v2 completed successfully.")
		return nil
	})
}
