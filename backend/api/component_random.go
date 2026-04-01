package api

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// RandomDistinctComponentIDs returns up to n distinct component IDs from the live catalog.
// If n <= 0 or the catalog is empty, returns nil. If n > len(catalog), uses len(catalog).
func RandomDistinctComponentIDs(catalog []models.Component, n int) []uuid.UUID {
	if n <= 0 || len(catalog) == 0 {
		return nil
	}
	if n > len(catalog) {
		n = len(catalog)
	}
	perm := rand.Perm(len(catalog))
	out := make([]uuid.UUID, n)
	for i := 0; i < n; i++ {
		out[i] = catalog[perm[i]].ID
	}
	return out
}

// RandomDistinctComponentIDStrings is like [RandomDistinctComponentIDs] but returns UUID strings for JSON/corpse storage.
func RandomDistinctComponentIDStrings(catalog []models.Component, n int) []string {
	ids := RandomDistinctComponentIDs(catalog, n)
	if len(ids) == 0 {
		return nil
	}
	s := make([]string, len(ids))
	for i, id := range ids {
		s[i] = id.String()
	}
	return s
}

// RollCorpseComponentSlotCount returns a 1d4 count (1–4) for how many distinct component types a corpse may offer when none were specified.
func RollCorpseComponentSlotCount() int {
	return rand.Intn(4) + 1
}
