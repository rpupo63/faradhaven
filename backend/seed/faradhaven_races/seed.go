package faradhaven_races

import (
	"encoding/json"
	"log"

	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed/batch"
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

// RaceComponentLink represents a race-component join table entry
type RaceComponentLink struct {
	RaceID      string `gorm:"type:uuid;primaryKey"`
	ComponentID string `gorm:"type:uuid;primaryKey"`
}

func (RaceComponentLink) TableName() string {
	return "race_components"
}

// SeedFaradhavenRaces creates Race, Trait, TraitOption, and RaceComponent rows using batch operations.
// Uses deterministic UUIDs so reseeding doesn't break character references.
func SeedFaradhavenRaces(tx *gorm.DB) error {
	raceSeeds := AllRaces()

	// Step 1: Collect all entities into slices
	races := make([]models.Race, 0, len(raceSeeds))
	var allTraits []models.Trait
	var allTraitOptions []models.TraitOption
	var allRaceComponents []RaceComponentLink

	for _, rs := range raceSeeds {
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

		races = append(races, models.Race{
			ID:                  raceID,
			Name:                rs.Name,
			Description:         rs.Description,
			CreatureType:        rs.CreatureType,
			Size:                rs.Size,
			BaseSpeed:           rs.BaseSpeed,
			PhotoURL:            rs.PhotoURL,
			AbilityScoreBonuses: bonusesJSON,
		})

		// Collect traits
		for _, ts := range rs.Traits {
			traitID := uuids.TraitUUID(rs.Name, ts.Name)
			levelReq := ts.LevelReq
			if levelReq == 0 {
				levelReq = 1
			}

			allTraits = append(allTraits, models.Trait{
				ID:             traitID,
				RaceID:         &raceID,
				Name:           ts.Name,
				Description:    ts.Description,
				LevelReq:       levelReq,
				ActionType:     ts.ActionType,
				UsesPerRest:    ts.UsesPerRest,
				ResetCondition: ts.ResetCondition,
				RangeValue:     ts.RangeValue,
				AreaOfEffect:   ts.AreaOfEffect,
				SaveAbility:    ts.SaveAbility,
			})

			// Collect trait options
			for _, opt := range ts.Options {
				optionID := uuids.TraitOptionUUID(traitID, opt.Name)

				var optBonusesJSON datatypes.JSON
				if opt.AbilityScoreBonuses != nil {
					b, err := json.Marshal(opt.AbilityScoreBonuses)
					if err != nil {
						return err
					}
					optBonusesJSON = datatypes.JSON(b)
				}

				allTraitOptions = append(allTraitOptions, models.TraitOption{
					ID:                  optionID,
					TraitID:             traitID,
					Name:                opt.Name,
					Description:         opt.Description,
					AbilityScoreBonuses: optBonusesJSON,
				})
			}
		}

		// Collect race-component links
		for _, compName := range rs.ComponentNames {
			componentID := uuids.ComponentUUID(compName)
			allRaceComponents = append(allRaceComponents, RaceComponentLink{
				RaceID:      raceID.String(),
				ComponentID: componentID.String(),
			})
		}
	}

	// Step 2: Clear child tables (will be fully replaced)
	// Order matters: delete children before parents due to foreign keys
	if err := tx.Exec("DELETE FROM trait_options").Error; err != nil {
		log.Printf("  Note: Could not clear trait_options: %v", err)
	}
	if err := tx.Exec("DELETE FROM traits").Error; err != nil {
		log.Printf("  Note: Could not clear traits: %v", err)
	}
	if err := tx.Exec("DELETE FROM race_components").Error; err != nil {
		log.Printf("  Note: Could not clear race_components: %v", err)
	}

	// Step 3: Batch upsert races (ON CONFLICT DO UPDATE)
	if err := batch.UpsertBatchUpdateAll(tx, races, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d races", len(races))

	// Step 4: Batch insert traits
	if err := batch.InsertBatch(tx, allTraits, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d traits", len(allTraits))

	// Step 5: Batch insert trait options
	if err := batch.InsertBatch(tx, allTraitOptions, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d trait options", len(allTraitOptions))

	// Step 6: Batch insert race-component links
	// Use raw insert to handle the join table
	if len(allRaceComponents) > 0 {
		if err := tx.CreateInBatches(allRaceComponents, batch.DefaultBatchSize).Error; err != nil {
			return err
		}
		log.Printf("Inserted %d race-component links", len(allRaceComponents))
	}

	return nil
}
