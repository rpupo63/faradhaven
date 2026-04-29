package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/services"
	"github.com/rs/zerolog/log"
)

type harvestHandler struct {
	harvestingService *services.HarvestingService
	charRepo          database.CharacterRepository
}

func newHarvestHandler(
	harvestingService *services.HarvestingService,
	charRepo database.CharacterRepository,
) *harvestHandler {
	return &harvestHandler{
		harvestingService: harvestingService,
		charRepo:          charRepo,
	}
}

// getHarvestableAbilities handles GET /api/beasts/{beastID}/harvestable-abilities
func (h *harvestHandler) getHarvestableAbilities() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		beastIDStr := chi.URLParam(r, "beastID")
		beastID, err := uuid.Parse(beastIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid beast ID")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userID, err := uuid.Parse(authUserID)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		// For simplicity, let's just grab the first character for the user
		// A more robust solution would involve the user selecting their active character
		chars, err := h.charRepo.FindByUserID(userID)
		if err != nil || len(chars) == 0 {
			respondError(w, http.StatusBadRequest, "No active character found for user")
			return
		}
		characterID := chars[0].ID

		resp, err := h.harvestingService.GetHarvestableAbilities(characterID, beastID)
		if err != nil {
			log.Error().Err(err).Str("beastID", beastIDStr).Msg("Failed to get harvestable abilities")
			respondError(w, http.StatusInternalServerError, "Failed to get harvestable abilities")
			return
		}

		respondJSON(w, http.StatusOK, resp)
	}
}

// confirmHarvest handles POST /api/characters/{characterID}/harvest
func (h *harvestHandler) confirmHarvest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterIDStr := chi.URLParam(r, "characterID")
		characterID, err := uuid.Parse(characterIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		if _, err := ctxGetUserID(r.Context()); err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var req services.ConfirmHarvestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		req.CharacterID = characterID // Ensure character ID from URL path is used

		character, err := h.harvestingService.ConfirmHarvest(req)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to confirm harvest")
			respondError(w, http.StatusInternalServerError, "Failed to confirm harvest: "+err.Error())
			return
		}

		respondJSON(w, http.StatusOK, character)
	}
}
