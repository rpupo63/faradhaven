package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type characterEffectHandler struct {
	effectService *services.EffectService
}

func newCharacterEffectHandler(effectService *services.EffectService) *characterEffectHandler {
	return &characterEffectHandler{effectService: effectService}
}

// GetActiveEffects returns all active effects on a character
func (h *characterEffectHandler) GetActiveEffects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		effects, err := h.effectService.GetActiveEffects(characterID)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterID.String()).Msg("Failed to get active effects")
			respondError(w, http.StatusInternalServerError, "Failed to get active effects")
			return
		}

		response := make([]CharacterEffectResponse, len(effects))
		for i, ce := range effects {
			response[i] = CharacterEffectResponse{
				ID:                ce.ID,
				CharacterID:       ce.CharacterID,
				EffectID:          ce.EffectID,
				EffectName:        ce.Effect.Name,
				EffectDescription: ce.Effect.Description,
				EffectMechanics:   ce.Effect.Mechanics,
				EffectCategory:    ce.Effect.Category,
				Duration:          ce.Duration,
				Source:            ce.Source,
				Stacks:            ce.Stacks,
				MaxStacks:         ce.MaxStacks,
				DurationRounds:    ce.DurationRounds,
				DurationMinutes:   ce.DurationMinutes,
				SourceCharacterID: ce.SourceCharacterID,
				IsConcentration:   ce.IsConcentration,
			}
		}

		respondJSON(w, http.StatusOK, response)
	}
}

// ApplyEffect applies an effect to a character
func (h *characterEffectHandler) ApplyEffect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		var req ApplyEffectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		opts := services.ApplyEffectOptions{
			Stacks:            req.Stacks,
			MaxStacks:         req.MaxStacks,
			DurationRounds:    req.DurationRounds,
			DurationMinutes:   req.DurationMinutes,
			Source:            req.Source,
			SourceCharacterID: req.SourceCharacterID,
			SourceSpellID:     req.SourceSpellID,
			IsConcentration:   req.IsConcentration,
			StackIfExists:     req.StackIfExists,
		}

		effect, err := h.effectService.ApplyEffect(characterID, req.EffectID, opts)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterID.String()).Msg("Failed to apply effect")
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := CharacterEffectResponse{
			ID:                effect.ID,
			CharacterID:       effect.CharacterID,
			EffectID:          effect.EffectID,
			EffectName:        effect.Effect.Name,
			EffectDescription: effect.Effect.Description,
			EffectMechanics:   effect.Effect.Mechanics,
			EffectCategory:    effect.Effect.Category,
			Duration:          effect.Duration,
			Source:            effect.Source,
			Stacks:            effect.Stacks,
			MaxStacks:         effect.MaxStacks,
			DurationRounds:    effect.DurationRounds,
			DurationMinutes:   effect.DurationMinutes,
			SourceCharacterID: effect.SourceCharacterID,
			IsConcentration:   effect.IsConcentration,
		}

		respondJSON(w, http.StatusCreated, response)
	}
}

// ModifyStacks modifies the stack count of an effect
func (h *characterEffectHandler) ModifyStacks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		effectInstanceID, err := uuid.Parse(chi.URLParam(r, "effectInstanceID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid effect instance ID")
			return
		}

		var req ModifyStacksRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		effect, err := h.effectService.ModifyStacks(effectInstanceID, req.Delta)
		if err != nil {
			log.Error().Err(err).Str("effectInstanceID", effectInstanceID.String()).Msg("Failed to modify stacks")
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Effect was deleted if stacks hit zero
		if effect == nil {
			respondJSON(w, http.StatusOK, map[string]string{"status": "removed", "reason": "stacks_depleted"})
			return
		}

		response := CharacterEffectResponse{
			ID:                effect.ID,
			CharacterID:       effect.CharacterID,
			EffectID:          effect.EffectID,
			EffectName:        effect.Effect.Name,
			EffectDescription: effect.Effect.Description,
			EffectMechanics:   effect.Effect.Mechanics,
			EffectCategory:    effect.Effect.Category,
			Duration:          effect.Duration,
			Source:            effect.Source,
			Stacks:            effect.Stacks,
			MaxStacks:         effect.MaxStacks,
			DurationRounds:    effect.DurationRounds,
			DurationMinutes:   effect.DurationMinutes,
			SourceCharacterID: effect.SourceCharacterID,
			IsConcentration:   effect.IsConcentration,
		}

		respondJSON(w, http.StatusOK, response)
	}
}

// RemoveEffect removes an effect from a character
func (h *characterEffectHandler) RemoveEffect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		effectInstanceID, err := uuid.Parse(chi.URLParam(r, "effectInstanceID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid effect instance ID")
			return
		}

		// Check for request body (optional)
		var req RemoveEffectRequest
		_ = json.NewDecoder(r.Body).Decode(&req) // Ignore error, use defaults

		err = h.effectService.RemoveEffect(effectInstanceID, req.RemoveAllStacks, req.StacksToRemove)
		if err != nil {
			log.Error().Err(err).Str("effectInstanceID", effectInstanceID.String()).Msg("Failed to remove effect")
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"status": "removed"})
	}
}

// TickDuration processes duration countdown for all effects
func (h *characterEffectHandler) TickDuration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		var req TickDurationRequest
		_ = json.NewDecoder(r.Body).Decode(&req) // Ignore error, use default

		rounds := req.Rounds
		if rounds <= 0 {
			rounds = 1
		}

		result, err := h.effectService.TickDuration(characterID, rounds)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterID.String()).Msg("Failed to tick duration")
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		expired := make([]ExpiredEffectInfo, len(result.Expired))
		for i, exp := range result.Expired {
			expired[i] = ExpiredEffectInfo{
				EffectID:   exp.EffectID,
				EffectName: exp.EffectName,
				Reason:     exp.Reason,
			}
		}

		respondJSON(w, http.StatusOK, TickDurationResponse{
			Expired: expired,
			Rounds:  result.RoundsTickerd,
		})
	}
}

// BreakConcentration breaks all concentration effects on a character
func (h *characterEffectHandler) BreakConcentration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		expired, err := h.effectService.BreakConcentration(characterID)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterID.String()).Msg("Failed to break concentration")
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		broken := make([]ExpiredEffectInfo, len(expired))
		for i, exp := range expired {
			broken[i] = ExpiredEffectInfo{
				EffectID:   exp.EffectID,
				EffectName: exp.EffectName,
				Reason:     exp.Reason,
			}
		}

		respondJSON(w, http.StatusOK, BreakConcentrationResponse{
			Broken: broken,
		})
	}
}
