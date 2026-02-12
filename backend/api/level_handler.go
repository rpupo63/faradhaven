package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type levelHandler struct {
	levelUpService         *services.LevelUpService
	classRepo              *database.ClassRepo
	characterResourceRepo  database.CharacterResourceRepository
	beastRepo              database.BeastRepository
	consumptionHistoryRepo database.ConsumptionHistoryRepository
}

func newLevelHandler(
	levelUpService *services.LevelUpService,
	classRepo *database.ClassRepo,
	characterResourceRepo database.CharacterResourceRepository,
	beastRepo database.BeastRepository,
	consumptionHistoryRepo database.ConsumptionHistoryRepository,
) *levelHandler {
	return &levelHandler{
		levelUpService:         levelUpService,
		classRepo:              classRepo,
		characterResourceRepo:  characterResourceRepo,
		beastRepo:              beastRepo,
		consumptionHistoryRepo: consumptionHistoryRepo,
	}
}

// buildClassResources aggregates resource definitions, level values, and character state
// into a response-ready slice of ClassResourceResponse.
func (h *levelHandler) buildClassResources(classID uuid.UUID, level int, characterID uuid.UUID) []ClassResourceResponse {
	defs, err := h.classRepo.FindResourceDefinitionsByClassID(classID)
	if err != nil || len(defs) == 0 {
		return nil
	}

	resourceMap, _ := h.classRepo.GetLevelResourceMap(classID, level)
	charResources, _ := h.characterResourceRepo.FindByCharacterID(characterID)
	charResMap := make(map[string]*models.CharacterResource, len(charResources))
	for _, cr := range charResources {
		charResMap[cr.ResourceKey] = cr
	}

	result := make([]ClassResourceResponse, 0, len(defs))
	for _, def := range defs {
		resp := ClassResourceResponse{
			Key:          def.ResourceKey,
			DisplayName:  def.DisplayName,
			Category:     def.Category,
			Description:  def.Description,
			Value:        resourceMap[def.ResourceKey],
			IsTrackable:  def.IsTrackable,
			DisplayOrder: def.DisplayOrder,
		}
		if def.IsTrackable {
			if cr, ok := charResMap[def.ResourceKey]; ok {
				resp.CurrentValue = &cr.CurrentValue
				resp.MaxValue = cr.MaxValue
			}
		}
		result = append(result, resp)
	}
	return result
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

		var source string
		if req.Source != nil {
			source = *req.Source
		}

		modifiedDelta := req.Delta

		// --- Lorewright: Predator's Strike Logic ---
		// Applies only if this is damage (delta < 0), a target is specified, and character is a Lorewright Warlord
		if modifiedDelta < 0 && req.TargetID != nil {
			// Fetch character with all relations to check class and archetype
			character, err := h.levelUpService.CharacterRepo.FindByIDWithRelations(characterID)
			if err != nil {
				log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to get character for Predator's Strike")
				// Continue without Predator's Strike if character fetch fails
			} else if character.Class.Name == "The Lorewright" && character.Archetype != nil && character.Archetype.Name == "Path of the Warlord" {
				// Fetch target beast
				targetBeast, err := h.beastRepo.FindByID(*req.TargetID)
				if err != nil {
					log.Error().Err(err).Stringer("targetBeastID", *req.TargetID).Msg("Failed to get target beast for Predator's Strike")
					// Continue without Predator's Strike if beast fetch fails
				} else {
					// Check consumption history for the last 24 hours
					since := time.Now().Add(-24 * time.Hour)
					_, err := h.consumptionHistoryRepo.FindRecentByCharacterAndType(character.ID, targetBeast.Type, since)
					if err == nil { // Entry found, meaning creature type was harvested recently
						wisdomMod := models.AbilityModifier(character.Wisdom)
						modifiedDelta -= wisdomMod // Delta is negative, so subtracting makes it more negative (more damage)
						log.Debug().Int("originalDelta", req.Delta).Int("modifiedDelta", modifiedDelta).Int("wisdomMod", wisdomMod).Msg("Predator's Strike applied")
					}
				}
			}
		}

		character, err := h.levelUpService.UpdateHP(userID, characterID, modifiedDelta, source)
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

		// Get class level for max spell points
		classLevel, _ := h.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
		maxSP := 0
		if classLevel != nil {
			maxSP = classLevel.MaxSpellPoints
		}

		resp := RestResponse{
			CurrentHP:          character.CurrentHP,
			MaxHP:              character.MaxHP,
			TempHP:             character.TempHP,
			CurrentSpellPoints: character.CurrentSpellPoints,
			MaxSpellPoints:     maxSP,
			HitDiceRemaining:   character.Level - character.HitDiceUsed,
			HitDiceTotal:       character.Level,
			ClassResources:     h.buildClassResources(character.ClassID, character.Level, character.ID),
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

		// Get class level for max spell points
		classLevel, _ := h.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
		maxSP := 0
		if classLevel != nil {
			maxSP = classLevel.MaxSpellPoints
		}

		resp := RestResponse{
			CurrentHP:          character.CurrentHP,
			MaxHP:              character.MaxHP,
			TempHP:             character.TempHP,
			CurrentSpellPoints: character.CurrentSpellPoints,
			MaxSpellPoints:     maxSP,
			HitDiceRemaining:   character.Level - character.HitDiceUsed,
			HitDiceTotal:       character.Level,
			ClassResources:     h.buildClassResources(character.ClassID, character.Level, character.ID),
		}
		respondJSON(w, http.StatusOK, resp)
	}
}
