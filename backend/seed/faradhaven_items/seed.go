package faradhaven_items

import (
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed/batch"
	"github.com/rpupo63/unified-personal-site-backend/seed/uuids"
	"gorm.io/gorm"
)

// SeedFaradhavenItems creates Weapons and Items using batch operations.
// Uses deterministic UUIDs so reseeding doesn't break character references.
func SeedFaradhavenItems(tx *gorm.DB) error {
	// Step 1: Collect all weapons into a slice
	weaponSeeds := AllWeapons()
	weapons := make([]models.Weapon, 0, len(weaponSeeds))
	var allDamages []models.WeaponDamage

	for _, ws := range weaponSeeds {
		weaponID := uuids.WeaponUUID(ws.Name)

		weapons = append(weapons, models.Weapon{
			ID:                  weaponID,
			Name:                ws.Name,
			Description:         ws.Description,
			Category:            ws.Category,
			Rarity:              ws.Rarity,
			RangeType:           ws.RangeType,
			Cost:                ws.Cost,
			Weight:              ws.Weight,
			AttackModifier:      ws.AttackModifier,
			Properties:          pq.StringArray(ws.Properties),
			RangeNormal:         ws.RangeNormal,
			RangeLong:           ws.RangeLong,
			VersatileDamageDice: ws.VersatileDamageDice,
			SecondaryEffect:     ws.SecondaryEffect,
		})

		// Collect damages for this weapon
		for _, ds := range ws.Damages {
			damageID := uuids.WeaponDamageUUID(weaponID, ds.DamageDice, ds.DamageType, ds.DamageCategory)
			allDamages = append(allDamages, models.WeaponDamage{
				ID:             damageID,
				WeaponID:       weaponID,
				DamageDice:     ds.DamageDice,
				DamageType:     ds.DamageType,
				DamageCategory: ds.DamageCategory,
			})
		}
	}

	// Step 2: Clear weapon_damages table (will be replaced entirely)
	if err := tx.Exec("DELETE FROM weapon_damages").Error; err != nil {
		return fmt.Errorf("could not clear weapon_damages: %w", err)
	}

	// Step 3: Batch upsert weapons (ON CONFLICT DO UPDATE)
	if err := batch.UpsertBatchUpdateAll(tx, weapons, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d weapons", len(weapons))

	// Step 4: Batch insert damages
	if err := batch.InsertBatch(tx, allDamages, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d weapon damages", len(allDamages))

	// Step 5: Collect all items into a slice
	itemSeeds := AllItems()
	items := make([]models.Item, 0, len(itemSeeds))

	for _, is := range itemSeeds {
		itemID := uuids.ItemUUID(is.Name)
		items = append(items, models.Item{
			ID:           itemID,
			Name:         is.Name,
			Description:  is.Description,
			Category:     is.Category,
			Rarity:       is.Rarity,
			Cost:         is.Cost,
			Weight:       is.Weight,
			Effects:      is.Effects,
			IsConsumable: is.IsConsumable,
		})
	}

	// Step 6: Batch upsert items (ON CONFLICT DO UPDATE)
	if err := batch.UpsertBatchUpdateAll(tx, items, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d items", len(items))

	return nil
}
