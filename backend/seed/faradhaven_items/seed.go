package faradhaven_items

import (
	"log"

	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// SeedFaradhavenItems creates Weapons and their related tables.
func SeedFaradhavenItems(db *gorm.DB) error {
	for _, ws := range AllWeapons() {
		var w models.Weapon
		err := db.Where("name = ?", ws.Name).First(&w).Error
		
		// If creating new or updating, we map fields
		weaponData := models.Weapon{
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

		if err == gorm.ErrRecordNotFound {
			// Create new
			if err := db.Create(&weaponData).Error; err != nil {
				return err
			}
			log.Printf("Created weapon: %s", ws.Name)
			// Re-fetch to get ID for child tables
			db.Where("name = ?", ws.Name).First(&w)
		} else if err == nil {
			// Update existing
			if err := db.Model(&w).Updates(weaponData).Error; err != nil {
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
	return nil
}
