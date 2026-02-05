package faradhaven_races

import (
	"encoding/json"
	"log"

	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed/uuids"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AllRaces returns all Faradhaven race seeds for seeding
func AllRaces() []FaradhavenRaceSeed {
	return []FaradhavenRaceSeed{
		Aasimar(),
		Boggart(),
		Changeling(),
		LorwynChangeling(),
		Dragonborn(),
		Dhampir(),
		Dwarf(),
		Flamekin(),
		Elf(),
		Faerie(),
		Gnome(),
		Goliath(),
		Halfling(),
		Human(),
		Kalashtar(),
		Khoravar(),
		Orc(),
		Rimekin(),
		Shifter(),
		Tiefling(),
		Warforged(),
	}
}

// SeedFaradhavenRaces creates Race, Trait, TraitOption, and RaceComponent rows for all Faradhaven races.
// Uses deterministic UUIDs so reseeding doesn't break character references.
func SeedFaradhavenRaces(db *gorm.DB) error {
	for _, rs := range AllRaces() {
		// Generate deterministic UUID for this race
		raceID := uuids.RaceUUID(rs.Name)

		// Marshal AbilityScoreBonuses to JSON
		var bonusesJSON datatypes.JSON
		if rs.AbilityScoreBonuses != nil {
			b, err := json.Marshal(rs.AbilityScoreBonuses)
			if err != nil {
				return err
			}
			bonusesJSON = datatypes.JSON(b)
		}

		var r models.Race
		err := db.Where("id = ?", raceID).First(&r).Error
		switch {
		case err == gorm.ErrRecordNotFound:
			r = models.Race{
				ID:                  raceID,
				Name:                rs.Name,
				Description:         rs.Description,
				CreatureType:        rs.CreatureType,
				Size:                rs.Size,
				BaseSpeed:           rs.BaseSpeed,
				PhotoURL:            rs.PhotoURL,
				AbilityScoreBonuses: bonusesJSON,
			}
			if err := db.Create(&r).Error; err != nil {
				return err
			}
			log.Printf("Created race: %s (ID: %s)", r.Name, r.ID)
		case err == nil:
			updates := map[string]interface{}{
				"name":                  rs.Name,
				"description":           rs.Description,
				"creature_type":         rs.CreatureType,
				"size":                  rs.Size,
				"base_speed":            rs.BaseSpeed,
				"photo_url":             rs.PhotoURL,
				"ability_score_bonuses": bonusesJSON,
			}
			if err := db.Model(&r).Updates(updates).Error; err != nil {
				return err
			}
			log.Printf("Updated race: %s", r.Name)
			if err := db.Where("race_id = ?", r.ID).Delete(&models.Trait{}).Error; err != nil {
				return err
			}
			// Clear existing race_components for re-seeding
			if err := db.Exec("DELETE FROM race_components WHERE race_id = ?", r.ID).Error; err != nil {
				return err
			}
		default:
			return err
		}

		for _, ts := range rs.Traits {
			t := models.Trait{
				RaceID:         &r.ID,
				Name:           ts.Name,
				Description:    ts.Description,
				LevelReq:       ts.LevelReq,
				ActionType:     ts.ActionType,
				UsesPerRest:    ts.UsesPerRest,
				ResetCondition: ts.ResetCondition,
				RangeValue:     ts.RangeValue,
				AreaOfEffect:   ts.AreaOfEffect,
				SaveAbility:    ts.SaveAbility,
			}
			if ts.LevelReq == 0 {
				t.LevelReq = 1
			}
			if err := db.Create(&t).Error; err != nil {
				return err
			}

			for _, opt := range ts.Options {
				// Marshal Option Bonuses
				var optBonusesJSON datatypes.JSON
				if opt.AbilityScoreBonuses != nil {
					b, err := json.Marshal(opt.AbilityScoreBonuses)
					if err != nil {
						return err
					}
					optBonusesJSON = datatypes.JSON(b)
				}

				o := models.TraitOption{
					TraitID:             t.ID,
					Name:                opt.Name,
					Description:         opt.Description,
					AbilityScoreBonuses: optBonusesJSON,
				}
				if err := db.Create(&o).Error; err != nil {
					return err
				}
			}
		}

		// Link components to race (for spell crafting)
		for _, compName := range rs.ComponentNames {
			// Use deterministic component UUID
			componentID := uuids.ComponentUUID(compName)

			// Verify component exists
			var component models.Component
			if err := db.Where("id = ?", componentID).First(&component).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					log.Printf("  Warning: Component '%s' not found for race %s", compName, rs.Name)
					continue
				}
				return err
			}

			// Insert into race_components join table
			if err := db.Exec("INSERT INTO race_components (race_id, component_id) VALUES (?, ?) ON CONFLICT DO NOTHING", r.ID, componentID).Error; err != nil {
				return err
			}
			log.Printf("  Linked component %s to race %s", compName, rs.Name)
		}
	}
	return nil
}
