package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type madnessHandler struct {
	madnessService *services.MadnessService
}

func newMadnessHandler(madnessService *services.MadnessService) *madnessHandler {
	return &madnessHandler{
		madnessService: madnessService,
	}
}

// rollMadness handles POST /api/characters/{characterID}/madness/roll
func (h *madnessHandler) rollMadness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterIDStr := chi.URLParam(r, "characterID")
		characterID, err := uuid.Parse(characterIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		_ = authUserID // Marked as used
		// In a real application, you'd verify the user owns the character
		// For now, we proceed assuming user has access.

		result, err := h.madnessService.RollMadness(characterID)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to roll madness")
			if err.Error() == "character is not a Lorewright" { // Specific error for user feedback
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to roll madness")
			return
		}

		respondJSON(w, http.StatusOK, result)
	}
}
