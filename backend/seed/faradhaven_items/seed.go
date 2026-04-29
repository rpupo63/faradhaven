package faradhaven_items

import (
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/seed/batch"
	"github.com/rpupo63/faradhaven/backend/seed/uuids"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SeedFaradhavenItems creates Weapons and Items using batch operations.
// Uses deterministic UUIDs so reseeding doesn't break character references.
func SeedFaradhavenItems(tx *gorm.DB) error {
	// Step 1: Collect all weapons into a slice
	weaponSeeds := AllWeapons()
	weapons := make([]models.Weapon, 0, len(weaponSeeds))
	var allDamages []models.WeaponDamage
	weaponThemes := make([]models.WeaponLootTheme, 0)
	weaponLocations := make([]models.WeaponLootLocation, 0)
	weaponSources := make([]models.WeaponLootSource, 0)
	weaponTiers := make([]models.WeaponLootTier, 0)
	weaponRewards := make([]models.WeaponLootRewardAmount, 0)
	weaponBands := make([]models.WeaponLootLevelBand, 0)

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
			dt, ok := models.ParseDamageType(ds.DamageType)
			if !ok {
				return fmt.Errorf("weapon %q: invalid damage type %q", ws.Name, ds.DamageType)
			}
			dc, ok := models.ParseWeaponDamageCategory(ds.DamageCategory)
			if !ok {
				return fmt.Errorf("weapon %q: invalid damage category %q", ws.Name, ds.DamageCategory)
			}
			damageID := uuids.WeaponDamageUUID(weaponID, ds.DamageDice, strings.TrimSpace(ds.DamageType), strings.TrimSpace(ds.DamageCategory))
			allDamages = append(allDamages, models.WeaponDamage{
				ID:             damageID,
				WeaponID:       weaponID,
				DamageDice:     ds.DamageDice,
				DamageType:     dt,
				DamageCategory: dc,
			})
		}
		if ws.LootTags.Weight <= 0 {
			ws.LootTags.Weight = 1.0
		}
		for _, v := range ws.LootTags.Themes {
			if theme, ok := models.ParseLootTheme(v); ok {
				weaponThemes = append(weaponThemes, models.WeaponLootTheme{WeaponID: weaponID, Theme: theme, Weight: ws.LootTags.Weight})
			}
		}
		for _, v := range ws.LootTags.Locations {
			if loc, ok := models.ParseLootLocation(v); ok {
				weaponLocations = append(weaponLocations, models.WeaponLootLocation{WeaponID: weaponID, Location: loc, Weight: ws.LootTags.Weight})
			}
		}
		for _, v := range ws.LootTags.Sources {
			if source, ok := models.ParseLootSource(v); ok {
				weaponSources = append(weaponSources, models.WeaponLootSource{WeaponID: weaponID, Source: source, Weight: ws.LootTags.Weight})
			}
		}
		for _, v := range ws.LootTags.Tiers {
			if tier, ok := models.ParseLootTier(v); ok {
				weaponTiers = append(weaponTiers, models.WeaponLootTier{WeaponID: weaponID, Tier: tier, Weight: ws.LootTags.Weight})
			}
		}
		for _, v := range ws.LootTags.RewardAmounts {
			if reward, ok := models.ParseLootRewardAmount(v); ok {
				weaponRewards = append(weaponRewards, models.WeaponLootRewardAmount{WeaponID: weaponID, RewardAmount: reward, Weight: ws.LootTags.Weight})
			}
		}
		for _, v := range ws.LootTags.LevelBands {
			if band, ok := models.ParseLootLevelBand(v); ok {
				weaponBands = append(weaponBands, models.WeaponLootLevelBand{WeaponID: weaponID, LevelBand: band, Weight: ws.LootTags.Weight})
			}
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
	itemThemes := make([]models.ItemLootTheme, 0)
	itemLocations := make([]models.ItemLootLocation, 0)
	itemSources := make([]models.ItemLootSource, 0)
	itemTiers := make([]models.ItemLootTier, 0)
	itemRewards := make([]models.ItemLootRewardAmount, 0)
	itemBands := make([]models.ItemLootLevelBand, 0)

	for _, is := range itemSeeds {
		itemID := uuids.ItemUUID(is.Name)
		items = append(items, models.Item{
			ID:                  itemID,
			Name:                is.Name,
			Description:         is.Description,
			Category:            is.Category,
			Rarity:              is.Rarity,
			Cost:                is.Cost,
			Weight:              is.Weight,
			Effects:             is.Effects,
			IsConsumable:        is.IsConsumable,
			ArmorType:           is.ArmorType,
			BaseAC:              is.BaseAC,
			StrengthRequirement: is.StrengthRequirement,
			StealthDisadvantage: is.StealthDisadvantage,
		})
		if is.LootTags.Weight <= 0 {
			is.LootTags.Weight = 1.0
		}
		for _, v := range is.LootTags.Themes {
			if theme, ok := models.ParseLootTheme(v); ok {
				itemThemes = append(itemThemes, models.ItemLootTheme{ItemID: itemID, Theme: theme, Weight: is.LootTags.Weight})
			}
		}
		for _, v := range is.LootTags.Locations {
			if loc, ok := models.ParseLootLocation(v); ok {
				itemLocations = append(itemLocations, models.ItemLootLocation{ItemID: itemID, Location: loc, Weight: is.LootTags.Weight})
			}
		}
		for _, v := range is.LootTags.Sources {
			if source, ok := models.ParseLootSource(v); ok {
				itemSources = append(itemSources, models.ItemLootSource{ItemID: itemID, Source: source, Weight: is.LootTags.Weight})
			}
		}
		for _, v := range is.LootTags.Tiers {
			if tier, ok := models.ParseLootTier(v); ok {
				itemTiers = append(itemTiers, models.ItemLootTier{ItemID: itemID, Tier: tier, Weight: is.LootTags.Weight})
			}
		}
		for _, v := range is.LootTags.RewardAmounts {
			if reward, ok := models.ParseLootRewardAmount(v); ok {
				itemRewards = append(itemRewards, models.ItemLootRewardAmount{ItemID: itemID, RewardAmount: reward, Weight: is.LootTags.Weight})
			}
		}
		for _, v := range is.LootTags.LevelBands {
			if band, ok := models.ParseLootLevelBand(v); ok {
				itemBands = append(itemBands, models.ItemLootLevelBand{ItemID: itemID, LevelBand: band, Weight: is.LootTags.Weight})
			}
		}
	}

	// Step 6: Batch upsert items (ON CONFLICT DO UPDATE)
	if err := batch.UpsertBatchUpdateAll(tx, items, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d items", len(items))

	// Step 7: Reset loot mapping join tables and insert rows from tagged seeds.
	for _, table := range []string{
		"item_loot_themes",
		"item_loot_locations",
		"item_loot_sources",
		"item_loot_tiers",
		"item_loot_reward_amounts",
		"item_loot_level_bands",
		"weapon_loot_themes",
		"weapon_loot_locations",
		"weapon_loot_sources",
		"weapon_loot_tiers",
		"weapon_loot_reward_amounts",
		"weapon_loot_level_bands",
	} {
		if err := tx.Exec("DELETE FROM " + table).Error; err != nil {
			return fmt.Errorf("could not clear %s: %w", table, err)
		}
	}

	if len(itemThemes) > 0 {
		if err := batch.InsertBatch(tx, itemThemes, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(itemLocations) > 0 {
		if err := batch.InsertBatch(tx, itemLocations, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(itemSources) > 0 {
		if err := batch.InsertBatch(tx, itemSources, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(itemTiers) > 0 {
		if err := batch.InsertBatch(tx, itemTiers, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(itemRewards) > 0 {
		if err := batch.InsertBatch(tx, itemRewards, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(itemBands) > 0 {
		if err := batch.InsertBatch(tx, itemBands, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(weaponThemes) > 0 {
		if err := batch.InsertBatch(tx, weaponThemes, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(weaponLocations) > 0 {
		if err := batch.InsertBatch(tx, weaponLocations, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(weaponSources) > 0 {
		if err := batch.InsertBatch(tx, weaponSources, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(weaponTiers) > 0 {
		if err := batch.InsertBatch(tx, weaponTiers, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(weaponRewards) > 0 {
		if err := batch.InsertBatch(tx, weaponRewards, batch.DefaultBatchSize); err != nil {
			return err
		}
	}
	if len(weaponBands) > 0 {
		if err := batch.InsertBatch(tx, weaponBands, batch.DefaultBatchSize); err != nil {
			return err
		}
	}

	// Step 8: Upsert level-based budget profiles (1..20).
	profiles := make([]models.LootLevelBudgetProfile, 0, 20)
	meansByLevel := []float64{
		8, 10, 12, 15, 18,
		22, 26, 30, 35, 40, 46, 52,
		60, 70, 82, 95, 110, 126, 143, 161,
	}
	sigmasByLevel := []float64{
		5.5, 6.0, 6.5, 7.2, 8.0,
		9.0, 10.0, 11.5, 13.0, 14.5, 16.0, 17.5,
		19.5, 22.0, 24.5, 27.0, 30.0, 33.0, 36.0, 40.0,
	}
	for level := 1; level <= 20; level++ {
		mean := meansByLevel[level-1]
		sigma := sigmasByLevel[level-1]
		profiles = append(profiles, models.LootLevelBudgetProfile{
			LootLevel: level,
			MeanGP:    mean,
			SigmaGP:   sigma,
			MinGP:     0,
			MaxGP:     mean + sigma*2.75,
		})
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "loot_level"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"mean_gp", "sigma_gp", "min_gp", "max_gp", "updated_at",
		}),
	}).Create(&profiles).Error; err != nil {
		return fmt.Errorf("failed to upsert loot level budget profiles: %w", err)
	}

	return nil
}
