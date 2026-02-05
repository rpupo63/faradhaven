package faradhaven_items

import (
	"log"

	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed/uuids"
	"gorm.io/gorm"
)

// SeedFaradhavenItems creates Weapons and their related tables.
// Uses deterministic UUIDs so reseeding doesn't break character references.
func SeedFaradhavenItems(db *gorm.DB) error {
	for _, ws := range AllWeapons() {
		// Generate deterministic UUID for this weapon
		weaponID := uuids.WeaponUUID(ws.Name)

		var w models.Weapon
		err := db.Where("id = ?", weaponID).First(&w).Error

		if err == gorm.ErrRecordNotFound {
			// Create new with deterministic UUID
			w = models.Weapon{
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
			}
			if err := db.Create(&w).Error; err != nil {
				return err
			}
			log.Printf("Created weapon: %s (ID: %s)", ws.Name, weaponID)
		} else if err == nil {
			// Update existing
			updates := map[string]interface{}{
				"name":                  ws.Name,
				"description":           ws.Description,
				"category":              ws.Category,
				"rarity":                ws.Rarity,
				"range_type":            ws.RangeType,
				"cost":                  ws.Cost,
				"weight":                ws.Weight,
				"attack_modifier":       ws.AttackModifier,
				"properties":            pq.StringArray(ws.Properties),
				"range_normal":          ws.RangeNormal,
				"range_long":            ws.RangeLong,
				"versatile_damage_dice": ws.VersatileDamageDice,
				"secondary_effect":      ws.SecondaryEffect,
			}
			if err := db.Model(&w).Updates(updates).Error; err != nil {
				return err
			}
			log.Printf("Updated weapon: %s", ws.Name)

			// Clear existing damages to re-seed
			if err := db.Where("weapon_id = ?", w.ID).Delete(&models.WeaponDamage{}).Error; err != nil {
				return err
			}
		} else {
			return err
		}

		// Seed Damages
		for _, ds := range ws.Damages {
			damage := models.WeaponDamage{
				WeaponID:       w.ID,
				DamageDice:     ds.DamageDice,
				DamageType:     ds.DamageType,
				DamageCategory: ds.DamageCategory,
			}
			if err := db.Create(&damage).Error; err != nil {
				return err
			}
		}
	}

	for _, is := range AllItems() {
		// Generate deterministic UUID for this item
		itemID := uuids.ItemUUID(is.Name)

		var it models.Item
		err := db.Where("id = ?", itemID).First(&it).Error

		if err == gorm.ErrRecordNotFound {
			// Create new with deterministic UUID
			it = models.Item{
				ID:           itemID,
				Name:         is.Name,
				Description:  is.Description,
				Category:     is.Category,
				Rarity:       is.Rarity,
				Cost:         is.Cost,
				Weight:       is.Weight,
				Effects:      is.Effects,
				IsConsumable: is.IsConsumable,
			}
			if err := db.Create(&it).Error; err != nil {
				return err
			}
			log.Printf("Created item: %s (ID: %s)", is.Name, itemID)
		} else if err == nil {
			// Update existing
			updates := map[string]interface{}{
				"name":          is.Name,
				"description":   is.Description,
				"category":      is.Category,
				"rarity":        is.Rarity,
				"cost":          is.Cost,
				"weight":        is.Weight,
				"effects":       is.Effects,
				"is_consumable": is.IsConsumable,
			}
			if err := db.Model(&it).Updates(updates).Error; err != nil {
				return err
			}
			log.Printf("Updated item: %s", is.Name)
		} else {
			return err
		}
	}

	return nil
}
