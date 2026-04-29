package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rs/zerolog/log"
)

type effectHandler struct {
	effectRepo database.EffectRepository
}

func newEffectHandler(effectRepo database.EffectRepository) *effectHandler {
	return &effectHandler{effectRepo: effectRepo}
}

func (h *effectHandler) getAllEffects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		effects, err := h.effectRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get effects")
			respondError(w, http.StatusInternalServerError, "Failed to get effects")
			return
		}
		respondJSON(w, http.StatusOK, effects)
	}
}

func (h *effectHandler) getEffectByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "effectID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid effect ID")
			return
		}

		effect, err := h.effectRepo.FindByID(id)
		if err != nil {
			log.Error().Err(err).Str("effectID", idStr).Msg("Failed to get effect")
			respondError(w, http.StatusNotFound, "Effect not found")
			return
		}

		respondJSON(w, http.StatusOK, effect)
	}
}
