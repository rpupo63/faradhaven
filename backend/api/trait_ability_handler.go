package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rs/zerolog/log"
)

type abilityHandler struct {
	characterRepo         database.CharacterRepository
	characterResourceRepo database.CharacterResourceRepository
}

func newAbilityHandler(
	characterRepo database.CharacterRepository,
	characterResourceRepo database.CharacterResourceRepository,
) *abilityHandler {
	return &abilityHandler{
		characterRepo:         characterRepo,
		characterResourceRepo: characterResourceRepo,
	}
}

// useTraitAbility handles POST /api/characters/{characterID}/traits/{traitID}/use
// Decrements the use counter for a limited-use racial trait.
func (h *abilityHandler) useTraitAbility() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterIDStr := chi.URLParam(r, "characterID")
		characterID, err := uuid.Parse(characterIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}
		traitIDStr := chi.URLParam(r, "traitID")
		traitID, err := uuid.Parse(traitIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid trait ID")
			return
		}

		authUserIDStr, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		authUserID, _ := uuid.Parse(authUserIDStr)

		belongs, err := h.characterRepo.CharacterBelongsToUser(characterID, authUserID)
		if err != nil || !belongs {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		character, err := h.characterRepo.FindByID(characterID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		// Verify trait belongs to character's race or lineage
		var trait models.Trait
		if err := h.characterRepo.GetDB().First(&trait, "id = ?", traitID).Error; err != nil {
			respondError(w, http.StatusNotFound, "Trait not found")
			return
		}

		isValid := (trait.RaceID != nil && *trait.RaceID == character.RaceID) ||
			(trait.LineageID != nil && character.LineageID != nil && *trait.LineageID == *character.LineageID)
		if !isValid {
			respondError(w, http.StatusForbidden, "Trait does not belong to this character")
			return
		}

		resourceKey := "trait_" + traitIDStr
		if err := h.characterResourceRepo.UpdateCurrentValue(characterID, resourceKey, -1); err != nil {
			if strings.Contains(err.Error(), "insufficient") {
				respondError(w, http.StatusBadRequest, "No uses remaining for this trait")
				return
			}
			log.Error().Err(err).Msg("Failed to decrement trait use")
			respondError(w, http.StatusInternalServerError, "Failed to use trait")
			return
		}

		resource, _ := h.characterResourceRepo.FindByCharacterAndKey(characterID, resourceKey)
		resp := map[string]interface{}{
			"trait_id": traitIDStr,
		}
		if resource != nil {
			resp["current_uses"] = resource.CurrentValue
			resp["max_uses"] = resource.MaxValue
		}
		respondJSON(w, http.StatusOK, resp)
	}
}

// useFeatureAbility handles POST /api/characters/{characterID}/features/{featureID}/use
// Spends resource costs for an active class feature.
func (h *abilityHandler) useFeatureAbility() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterIDStr := chi.URLParam(r, "characterID")
		characterID, err := uuid.Parse(characterIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}
		featureIDStr := chi.URLParam(r, "featureID")
		featureID, err := uuid.Parse(featureIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid feature ID")
			return
		}

		authUserIDStr, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		authUserID, _ := uuid.Parse(authUserIDStr)

		belongs, err := h.characterRepo.CharacterBelongsToUser(characterID, authUserID)
		if err != nil || !belongs {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		character, err := h.characterRepo.FindByID(characterID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		// Load feature with its class level
		var feature models.LevelFeature
		if err := h.characterRepo.GetDB().Preload("ClassLevel").First(&feature, "id = ?", featureID).Error; err != nil {
			respondError(w, http.StatusNotFound, "Feature not found")
			return
		}

		// Verify the feature belongs to this character's class and is within their level
		if feature.ClassLevel.ClassID != character.ClassID {
			respondError(w, http.StatusForbidden, "Feature does not belong to this character's class")
			return
		}
		if feature.ClassLevel.Level > character.Level {
			respondError(w, http.StatusForbidden, "Character has not yet reached this feature's level")
			return
		}

		// Parse and spend resource costs
		type resourceEntryJSON struct {
			Key    string `json:"key"`
			Amount int    `json:"amount"`
		}
		var costs []resourceEntryJSON
		if len(feature.ResourceCosts) > 0 {
			if jsonErr := json.Unmarshal(feature.ResourceCosts, &costs); jsonErr != nil {
				log.Error().Err(jsonErr).Msg("Failed to parse feature resource costs")
			}
		}

		for _, cost := range costs {
			if spendErr := h.characterResourceRepo.UpdateCurrentValue(characterID, cost.Key, -cost.Amount); spendErr != nil {
				if strings.Contains(spendErr.Error(), "insufficient") {
					respondError(w, http.StatusBadRequest, fmt.Sprintf("Insufficient resources: %s", spendErr.Error()))
					return
				}
				log.Error().Err(spendErr).Str("key", cost.Key).Msg("Failed to spend resource for feature")
				respondError(w, http.StatusInternalServerError, "Failed to spend resource")
				return
			}
		}

		// Apply resource gains (e.g. Siphon grants 1 Ichor)
		var gains []resourceEntryJSON
		if len(feature.ResourceGains) > 0 {
			if jsonErr := json.Unmarshal(feature.ResourceGains, &gains); jsonErr != nil {
				log.Error().Err(jsonErr).Msg("Failed to parse feature resource gains")
			}
		}

		for _, gain := range gains {
			if gainErr := h.characterResourceRepo.UpdateCurrentValue(characterID, gain.Key, gain.Amount); gainErr != nil {
				log.Error().Err(gainErr).Str("key", gain.Key).Msg("Failed to apply resource gain for feature")
				respondError(w, http.StatusInternalServerError, "Failed to apply resource gain")
				return
			}
		}

		// If feature has its own limited uses, also decrement the feature use counter
		if feature.UsesPerRest != "" {
			featureKey := "feature_" + featureIDStr
			if spendErr := h.characterResourceRepo.UpdateCurrentValue(characterID, featureKey, -1); spendErr != nil {
				if strings.Contains(spendErr.Error(), "insufficient") {
					respondError(w, http.StatusBadRequest, "No uses remaining for this feature")
					return
				}
				// If resource doesn't exist yet, that's OK — feature use tracking is optional
				log.Debug().Err(spendErr).Str("feature_key", featureKey).Msg("Feature use counter not found")
			}
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"feature_id": featureIDStr,
			"message":    "Feature used successfully",
		})
	}
}
