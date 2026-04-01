package api

import (
	"sort"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// SpellPoolAllowlist returns component IDs a character always has from their **current**
// class and race definitions (join tables), independent of character_components rows.
func SpellPoolAllowlist(class *models.Class, race *models.Race) map[uuid.UUID]bool {
	m := make(map[uuid.UUID]bool)
	if class != nil {
		for _, c := range class.Components {
			m[c.ID] = true
		}
	}
	if race != nil {
		for _, c := range race.Components {
			m[c.ID] = true
		}
	}
	return m
}

// MergeSpellPoolComponents returns the union of class + race pool components for API
// responses (deduped by ID, stable sort). Always read from Class.Components and
// Race.Components, not from character inventory.
func MergeSpellPoolComponents(class *models.Class, race *models.Race) []models.Component {
	byID := make(map[uuid.UUID]models.Component)
	if class != nil {
		for _, c := range class.Components {
			byID[c.ID] = c
		}
	}
	if race != nil {
		for _, c := range race.Components {
			byID[c.ID] = c
		}
	}
	out := make([]models.Component, 0, len(byID))
	for _, c := range byID {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Category != out[j].Category {
			return out[i].Category < out[j].Category
		}
		return out[i].Name < out[j].Name
	})
	return out
}
