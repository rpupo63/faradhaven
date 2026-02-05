package api

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rs/zerolog/log"
)

// getCharacterSheet returns the fully calculated character sheet (Class + ClassLevel joined)
func (h *characterHandler) getCharacterSheet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByIDWithSkills(id)
		if err != nil {
			log.Error().Err(err).Str("characterID", idStr).Msg("Failed to get character")
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		// Preload weapons and items for the sheet
		if err := h.characterRepo.GetDB().Preload("Weapons.Damages").Preload("Items").First(character, "id = ?", id).Error; err != nil {
			log.Error().Err(err).Str("characterID", idStr).Msg("Failed to preload equipment")
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != character.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		// Fetch class with components for spell crafting
		class, err := h.classRepo.FindByIDWithLevels(character.ClassID)
		if err != nil {
			log.Error().Err(err).Str("class_id", character.ClassID.String()).Msg("Class not found")
			respondError(w, http.StatusNotFound, "Class not found for character")
			return
		}

		// Fetch race with components for spell crafting
		race, err := h.raceRepo.FindByID(character.RaceID)
		if err != nil {
			log.Error().Err(err).Str("race_id", character.RaceID.String()).Msg("Race not found")
			respondError(w, http.StatusNotFound, "Race not found for character")
			return
		}

		classLevel, err := h.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
		if err != nil {
			log.Error().Err(err).Int("level", character.Level).Msg("ClassLevel not found")
			respondError(w, http.StatusNotFound, "Level data not found for class")
			return
		}

		conMod := abilityMod(character.Constitution)
		dexMod := abilityMod(character.Dexterity)
		primaryMod := primaryAbilityMod(character, class.PrimaryAbility)

		// TotalHP = BaseHP + (AvgHitDie * (Level - 1)) + (ConMod * Level); BaseHP = HitDie (max at level 1)
		avgHitDie := (class.HitDie + 1) / 2
		baseHP := class.HitDie
		totalHP := baseHP + (avgHitDie * (character.Level - 1)) + (conMod * character.Level)
		if character.Level < 1 {
			totalHP = baseHP + conMod
		}

		ac := 8 + classLevel.ProficiencyBonus + dexMod
		saveDC := 8 + classLevel.ProficiencyBonus + primaryMod

		maxSP := classLevel.MaxSpellPoints
		currentSP := character.CurrentSpellPoints
		if currentSP > maxSP {
			currentSP = maxSP
		}

		// Map class.SavingThrows (ability names) to ability ids for frontend
		saveProfs := make([]string, 0, len(class.SavingThrows))
		for _, name := range class.SavingThrows {
			saveProfs = append(saveProfs, strings.ToLower(strings.TrimSpace(name)))
		}

		// Combine class and race components for spell crafting (deduplicated by ID)
		componentMap := make(map[uuid.UUID]models.Component)
		for _, comp := range class.Components {
			componentMap[comp.ID] = comp
		}
		for _, comp := range race.Components {
			componentMap[comp.ID] = comp
		}
		availableComponents := make([]models.Component, 0, len(componentMap))
		for _, comp := range componentMap {
			availableComponents = append(availableComponents, comp)
		}

		// Use persisted HP values if set, otherwise fall back to computed
		maxHP := character.MaxHP
		currentHP := character.CurrentHP
		if maxHP == 0 {
			// First-time calculation for existing characters without HP
			maxHP = totalHP
			currentHP = totalHP
		}

		profBonus := classLevel.ProficiencyBonus
		strMod := abilityMod(character.Strength)

		// Race Traits
		raceTraits, err := h.raceRepo.FindByIDWithTraits(character.RaceID)
		var traits []models.Trait
		if err == nil {
			traits = raceTraits.Traits
		}

		// Lineage Traits
		var lineage *models.Lineage
		if character.LineageID != nil {
			lineage, err = h.raceRepo.FindLineageByIDWithTraits(*character.LineageID)
			if err == nil && lineage != nil {
				traits = append(traits, lineage.LineageTraits...)
			}
		}

		sheet := CharacterSheetResponse{
			Character:                character,
			Class:                    class,
			ClassLevel:               classLevel,
			TotalHP:                  maxHP, // Keep for backwards compatibility
			MaxHP:                    maxHP,
			CurrentHP:                currentHP,
			TempHP:                   character.TempHP,
			AC:                       ac,
			SaveDC:                   saveDC,
			MaxSpellPoints:           maxSP,
			CurrentSpellPoints:       currentSP,
			SavingThrowProficiencies: saveProfs,
			AvailableComponents:      availableComponents,
			HitDiceTotal:             character.Level,
			HitDiceRemaining:         character.Level - character.HitDiceUsed,
			HitDie:                   class.HitDie,
			MeleeAttackBonus:         profBonus + strMod,
			RangedAttackBonus:        profBonus + dexMod,
			RaceTraits:               traits,
			Lineage:                  lineage,
			InventoryWeapons:         character.Weapons,
			InventoryItems:           character.Items,
		}
		respondJSON(w, http.StatusOK, sheet)
	}
}
