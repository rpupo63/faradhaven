package faradhaven_races

import (
	"log"

	"github.com/rpupo63/unified-personal-site-backend/models"
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

// SeedFaradhavenRaces creates Race, Trait, and TraitOption rows for all Faradhaven races.
// Skips races that already exist by name.
func SeedFaradhavenRaces(db *gorm.DB) error {
	for _, rs := range AllRaces() {
		var r models.Race
		err := db.Where("name = ?", rs.Name).First(&r).Error
		switch {
		case err == gorm.ErrRecordNotFound:
			r = models.Race{
				Name:         rs.Name,
				Description:  rs.Description,
				CreatureType: rs.CreatureType,
				Size:         rs.Size,
				BaseSpeed:    rs.BaseSpeed,
				PhotoURL:     rs.PhotoURL,
			}
			if err := db.Create(&r).Error; err != nil {
				return err
			}
			log.Printf("Created race: %s", r.Name)
		case err == nil:
			updates := map[string]interface{}{
				"description":   rs.Description,
				"creature_type": rs.CreatureType,
				"size":          rs.Size,
				"base_speed":    rs.BaseSpeed,
				"photo_url":     rs.PhotoURL,
			}
			if err := db.Model(&r).Updates(updates).Error; err != nil {
				return err
			}
			log.Printf("Updated race: %s", r.Name)
			if err := db.Where("race_id = ?", r.ID).Delete(&models.Trait{}).Error; err != nil {
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
				o := models.TraitOption{
					TraitID:     t.ID,
					Name:        opt.Name,
					Description: opt.Description,
				}
				if err := db.Create(&o).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}
