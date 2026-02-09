package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type levelHandler struct {
	levelUpService  *services.LevelUpService
	classRepo       *database.ClassRepo
	resourceService *services.ResourceService
}

func newLevelHandler(levelUpService *services.LevelUpService, classRepo *database.ClassRepo, resourceService *services.ResourceService) *levelHandler {
	return &levelHandler{
		levelUpService:  levelUpService,
		classRepo:       classRepo,
		resourceService: resourceService,
	}
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

// updateHP handles PATCH /api/character/{characterID}/hp
func (h *levelHandler) updateHP() http.HandlerFunc {
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

		var req UpdateHPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		character, err := h.levelUpService.UpdateHP(userID, characterID, req.Delta)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to update HP")
			if err == services.ErrUnauthorized {
				respondError(w, http.StatusForbidden, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to update HP")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"current_hp": character.CurrentHP,
			"max_hp":     character.MaxHP,
			"temp_hp":    character.TempHP,
		})
	}
}

// setTempHP handles PUT /api/character/{characterID}/temp-hp
func (h *levelHandler) setTempHP() http.HandlerFunc {
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

		var req SetTempHPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		character, err := h.levelUpService.SetTempHP(userID, characterID, req.TempHP)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to set temp HP")
			if err == services.ErrUnauthorized {
				respondError(w, http.StatusForbidden, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to set temp HP")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"current_hp": character.CurrentHP,
			"max_hp":     character.MaxHP,
			"temp_hp":    character.TempHP,
		})
	}
}

// useHitDice handles POST /api/character/{characterID}/hit-dice
func (h *levelHandler) useHitDice() http.HandlerFunc {
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

		var req UseHitDiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		result, err := h.levelUpService.UseHitDice(userID, characterID, req.Rolls)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to use hit dice")
			if err == services.ErrUnauthorized {
				respondError(w, http.StatusForbidden, err.Error())
				return
			}
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp := UseHitDiceResponse{
			CurrentHP:        result.Character.CurrentHP,
			MaxHP:            result.Character.MaxHP,
			HPHealed:         result.HPHealed,
			DiceUsed:         result.DiceUsed,
			DiceResults:      result.DiceResults,
			HitDiceRemaining: result.Character.Level - result.Character.HitDiceUsed,
		}
		respondJSON(w, http.StatusOK, resp)
	}
}

// shortRest handles POST /api/character/{characterID}/rest/short
func (h *levelHandler) shortRest() http.HandlerFunc {
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

		        character, err := h.levelUpService.ShortRest(userID, characterID)
				if err != nil {
					log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to short rest")
					if err == services.ErrUnauthorized {
						respondError(w, http.StatusForbidden, err.Error())
						return
					}
					respondError(w, http.StatusInternalServerError, "Failed to short rest")
					return
				}
		
				// Get class level for max values
				classLevel, _ := h.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
				maxSP := 0
				maxStability := 0
				if classLevel != nil {
					maxSP = classLevel.MaxSpellPoints
					maxStability = classLevel.MaxStability
				}
		
				resp := RestResponse{
					CurrentHP:          character.CurrentHP,
					MaxHP:              character.MaxHP,
					TempHP:             character.TempHP,
					CurrentSpellPoints: character.CurrentSpellPoints,
					MaxSpellPoints:     maxSP,
					HitDiceRemaining:   character.Level - character.HitDiceUsed,
					HitDiceTotal:       character.Level,
					CurrentStability:   character.CurrentStability,
					MaxStability:       maxStability,
					CurrentBloodIchor:  character.CurrentBloodIchor,
					MaxBloodIchor:      h.resourceService.ComputeMaxBloodIchor(character),
					SanguineMP:         character.SanguineMP,
					SanguineBR:         character.SanguineBR,
					MadnessCastCount:   character.MadnessCastCount,
				}
				respondJSON(w, http.StatusOK, resp)
	}
}

// longRest handles POST /api/character/{characterID}/rest/long
func (h *levelHandler) longRest() http.HandlerFunc {
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

		character, err := h.levelUpService.LongRest(userID, characterID)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to long rest")
			if err == services.ErrUnauthorized {
				respondError(w, http.StatusForbidden, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, "Failed to long rest")
			return
		}

		// Get class level for max values
		classLevel, _ := h.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
		maxSP := 0
		maxStability := 0
		if classLevel != nil {
			maxSP = classLevel.MaxSpellPoints
			maxStability = classLevel.MaxStability
		}

		resp := RestResponse{
			CurrentHP:          character.CurrentHP,
			MaxHP:              character.MaxHP,
			TempHP:             character.TempHP,
			CurrentSpellPoints: character.CurrentSpellPoints,
			MaxSpellPoints:     maxSP,
			HitDiceRemaining:   character.Level - character.HitDiceUsed,
			HitDiceTotal:       character.Level,
			CurrentStability:   character.CurrentStability,
			MaxStability:       maxStability,
			CurrentBloodIchor:  character.CurrentBloodIchor,
			MaxBloodIchor:      h.resourceService.ComputeMaxBloodIchor(character),
			MadnessCastCount:   character.MadnessCastCount,
		}
		respondJSON(w, http.StatusOK, resp)
	}
}
