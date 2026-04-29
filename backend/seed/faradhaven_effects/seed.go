package faradhaven_effects

import (
	"log"

	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/seed/batch"
	"github.com/rpupo63/faradhaven/backend/seed/uuids"
	"gorm.io/gorm"
)

// SeedFaradhavenEffects creates Effect rows using batch operations.
// Uses deterministic UUIDs so reseeding doesn't duplicate.
func SeedFaradhavenEffects(tx *gorm.DB) error {
	effectSeeds := AllEffects()
	effects := make([]models.Effect, 0, len(effectSeeds))

	for _, es := range effectSeeds {
		effectID := uuids.EffectUUID(es.Name)
		effects = append(effects, models.Effect{
			ID:          effectID,
			Name:        es.Name,
			Description: es.Description,
			Category:    es.Category,
			Mechanics:   es.Mechanics,
		})
	}

	// Batch upsert effects (ON CONFLICT DO UPDATE)
	if err := batch.UpsertBatchUpdateAll(tx, effects, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d effects", len(effects))

	return nil
}
