package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type levelHandler struct {
	levelUpService *services.LevelUpService
}

func newLevelHandler(levelUpService *services.LevelUpService) *levelHandler {
	return &levelHandler{levelUpService: levelUpService}
}

// levelUp handles POST /api/character/{characterID}/level-up
func (h *levelHandler) levelUp() http.HandlerFunc {
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
		userID, err := uuid.Parse(authUserID)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		var req services.LevelUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		req.CharacterID = characterID

		resp, err := h.levelUpService.LevelUp(userID, req)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Level-up failed")
			switch err {
			case services.ErrMaxLevelReached:
				respondError(w, http.StatusBadRequest, err.Error())
			case services.ErrInsufficientASIPoints:
				respondError(w, http.StatusBadRequest, err.Error())
			case services.ErrUnauthorized:
				respondError(w, http.StatusForbidden, err.Error())
			default:
				respondError(w, http.StatusInternalServerError, "Level-up failed")
			}
			return
		}

		respondJSON(w, http.StatusOK, resp)
	}
}

// levelDown handles POST /api/character/{characterID}/level-down
func (h *levelHandler) levelDown() http.HandlerFunc {
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
		userID, err := uuid.Parse(authUserID)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		resp, err := h.levelUpService.LevelDown(userID, characterID)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Level-down failed")
			switch err {
			case services.ErrMinLevelReached, services.ErrNoHistoryForLevel:
				respondError(w, http.StatusBadRequest, err.Error())
			case services.ErrUnauthorized:
				respondError(w, http.StatusForbidden, err.Error())
			default:
				respondError(w, http.StatusInternalServerError, "Level-down failed")
			}
			return
		}

		respondJSON(w, http.StatusOK, resp)
	}
}

// getLevelHistory handles GET /api/character/{characterID}/level-history
func (h *levelHandler) getLevelHistory() http.HandlerFunc {
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
		userID, err := uuid.Parse(authUserID)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		history, err := h.levelUpService.GetLevelHistory(userID, characterID)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to get level history")
			if err == services.ErrUnauthorized {
				respondError(w, http.StatusForbidden, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to get level history")
			return
		}

		respondJSON(w, http.StatusOK, history)
	}
}

// getLevelUpPreview handles GET /api/character/{characterID}/level-up/preview
func (h *levelHandler) getLevelUpPreview() http.HandlerFunc {
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
		userID, err := uuid.Parse(authUserID)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		preview, err := h.levelUpService.GetLevelUpPreview(userID, characterID)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to get level-up preview")
			if err == services.ErrMaxLevelReached {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			if err == services.ErrUnauthorized {
				respondError(w, http.StatusForbidden, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to get preview")
			return
		}

		respondJSON(w, http.StatusOK, preview)
	}
}
