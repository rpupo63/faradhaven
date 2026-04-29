package api

import (
	"sort"

	"github.com/rpupo63/faradhaven/backend/models"
)

func spellHasLogicaComponent(spell *models.Spell) bool {
	if spell == nil {
		return false
	}
	for i := range spell.Components {
		if spell.Components[i].Category == models.CategoryLogica {
			return true
		}
	}
	return false
}

// powderMageSpellMatchesSavedSpeedDial returns true if the spell's component signature matches a saved Speed Dial entry.
// Same rules as the forge: ordered match when Logica is present, multiset otherwise.
func powderMageSpellMatchesSavedSpeedDial(spell *models.Spell, saved []*models.SavedSpell) bool {
	if spell == nil || len(spell.Components) == 0 {
		return false
	}
	for _, s := range saved {
		if powderMageSpellMatchesOneSaved(spell, s) {
			return true
		}
	}
	return false
}

func powderMageSpellMatchesOneSaved(spell *models.Spell, saved *models.SavedSpell) bool {
	if saved == nil || len(saved.ComponentIDs) == 0 {
		return false
	}
	ids := make([]string, len(spell.Components))
	for i := range spell.Components {
		ids[i] = spell.Components[i].ID.String()
	}
	savedArr := []string(saved.ComponentIDs)
	if len(ids) != len(savedArr) {
		return false
	}
	if spellHasLogicaComponent(spell) {
		for i := range ids {
			if ids[i] != savedArr[i] {
				return false
			}
		}
		return true
	}
	a := append([]string(nil), ids...)
	b := append([]string(nil), savedArr...)
	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
