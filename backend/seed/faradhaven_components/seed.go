package faradhaven_components

import (
	"log"

	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/seed/batch"
	"github.com/rpupo63/faradhaven/backend/seed/uuids"
	"gorm.io/gorm"
)

// SeedFaradhavenComponents creates Component rows using batch operations.
// Uses deterministic UUIDs so reseeding doesn't duplicate.
func SeedFaradhavenComponents(tx *gorm.DB) error {
	componentSeeds := AllComponents()
	components := make([]models.Component, 0, len(componentSeeds))

	for _, cs := range componentSeeds {
		componentID := uuids.ComponentUUID(cs.Name)
		components = append(components, models.Component{
			ID:             componentID,
			Name:           cs.Name,
			Symbol:         cs.Symbol,
			RpgAwesomeIcon: cs.RpgAwesomeIcon,
			Category:       cs.Category,
			Description:    cs.Description,
			Element:        cs.Element,
			Tier:           cs.Tier,
		})
	}

	if err := batch.UpsertBatchUpdateAll(tx, components, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d components", len(components))

	return nil
}
