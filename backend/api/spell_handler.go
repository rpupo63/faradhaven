package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type spellHandler struct {
	spellRepo          database.SpellRepository
	synthesisService   *services.SpellSynthesisService
	interpreterService *services.ComponentInterpreterService
}

func newSpellHandler(spellRepo database.SpellRepository, synthesisService *services.SpellSynthesisService, interpreterService *services.ComponentInterpreterService) *spellHandler {
	return &spellHandler{spellRepo: spellRepo, synthesisService: synthesisService, interpreterService: interpreterService}
}

// SpellResponse is a custom response for a spell that includes character details
type SpellResponse struct {
	*models.Spell
	CharacterName  *string `json:"character_name,omitempty"`
	CharacterClass *string `json:"character_class,omitempty"`
}

// SynthesizeRequest is the request body for the synthesize endpoint
type SynthesizeRequest struct {
	ComponentIDs []uuid.UUID `json:"component_ids"`
}

// getAllSpells returns all spells
func (h *spellHandler) getAllSpells() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spells, err := h.spellRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get spells")
			respondError(w, http.StatusInternalServerError, "Failed to get spells")
			return
		}

		response := make([]SpellResponse, len(spells))
		for i, spell := range spells {
			response[i] = SpellResponse{
				Spell: spell,
			}
			if spell.Character != nil {
				response[i].CharacterName = &spell.Character.Name
				response[i].CharacterClass = &spell.Character.Class.Name
			}
		}

		respondJSON(w, http.StatusOK, response)
	}
}

// getSpell returns a spell by ID
func (h *spellHandler) getSpell() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "spellID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid spell ID")
			return
		}

		spell, err := h.spellRepo.FindByID(id)
		if err != nil {
			log.Error().Err(err).Str("spellID", idStr).Msg("Failed to get spell")
			respondError(w, http.StatusNotFound, "Spell not found")
			return
		}
		respondJSON(w, http.StatusOK, spell)
	}
}

// getSpellsByUser returns all spells for a user
func (h *spellHandler) getSpellsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "userID")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != userIDStr {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		spells, err := h.spellRepo.FindByUserID(userID)
		if err != nil {
			log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get spells")
			respondError(w, http.StatusInternalServerError, "Failed to get spells")
			return
		}
		respondJSON(w, http.StatusOK, spells)
	}
}

// getSpellsByCharacter returns all prepared spells for a character
func (h *spellHandler) getSpellsByCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterIDStr := chi.URLParam(r, "characterID")
		characterID, err := uuid.Parse(characterIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		spells, err := h.spellRepo.FindByCharacterID(characterID)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to get character spells")
			respondError(w, http.StatusInternalServerError, "Failed to get character spells")
			return
		}
		respondJSON(w, http.StatusOK, spells)
	}
}

// synthesizeSpell returns auto-derived spell properties from components
func (h *spellHandler) synthesizeSpell() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SynthesizeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if len(req.ComponentIDs) == 0 {
			respondError(w, http.StatusBadRequest, "At least one component ID is required")
			return
		}

		components, err := h.synthesisService.FetchComponents(req.ComponentIDs)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch components for synthesis")
			respondError(w, http.StatusInternalServerError, "Failed to fetch components")
			return
		}

		synthesis := h.synthesisService.Synthesize(components)
		respondJSON(w, http.StatusOK, synthesis)
	}
}

// createSpell creates a new spell
func (h *spellHandler) createSpell() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateSpellRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Name == "" {
			respondError(w, http.StatusBadRequest, "Name is required")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != req.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		spell := &models.Spell{
			UserID:        req.UserID,
			CharacterID:   req.CharacterID,
			Name:          req.Name,
			Description:   req.Description,
			SlotLevel:     req.SlotLevel,
			Type:          req.Type,
			Range:         req.Range,
			Duration:      req.Duration,
			Concentration: req.Concentration,
			SaveAttr:      req.SaveAttr,
			DamageDice:    req.DamageDice,
			DamageType:    req.DamageType,
			AddModifier:   req.AddModifier,
		}

		// Run synthesis if components are provided
		if len(req.ComponentIDs) > 0 {
			components, err := h.synthesisService.FetchComponents(req.ComponentIDs)
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch components for synthesis")
				respondError(w, http.StatusInternalServerError, "Failed to fetch components")
				return
			}

			synthesis := h.synthesisService.Synthesize(components)

			if len(synthesis.ValidationErrors) > 0 {
				respondJSON(w, http.StatusBadRequest, map[string]interface{}{
					"error":             "Component validation failed",
					"validation_errors": synthesis.ValidationErrors,
				})
				return
			}

			// Always set computed fields
			spell.ManaCost = synthesis.ManaCost
			spell.BaseTier = synthesis.BaseTier

			// Hybrid fill: use request value if provided, else fall back to suggestion
			if spell.SlotLevel == 0 {
				spell.SlotLevel = synthesis.SlotLevel
			}
			if spell.Type == "" {
				spell.Type = synthesis.SuggestedType
			}
			if spell.Range == nil && synthesis.SuggestedRange != nil {
				spell.Range = synthesis.SuggestedRange
			}
			if spell.DamageType == nil && synthesis.SuggestedDamageType != nil {
				dt := models.DamageType(*synthesis.SuggestedDamageType)
				spell.DamageType = &dt
			}
			if spell.DamageDice == nil && synthesis.SuggestedDamageDice != nil {
				spell.DamageDice = synthesis.SuggestedDamageDice
			}
			if spell.Duration == nil && synthesis.SuggestedDuration != nil {
				spell.Duration = synthesis.SuggestedDuration
			}
			if !spell.Concentration && synthesis.SuggestedConcentration {
				spell.Concentration = true
			}
		}

		if spell.SlotLevel == 0 {
			spell.SlotLevel = 1
		}
		if spell.Type == "" {
			spell.Type = "Utility"
		}

		if err := h.spellRepo.Add(spell, req.ComponentIDs); err != nil {
			log.Error().Err(err).Msg("Failed to create spell")
			respondError(w, http.StatusInternalServerError, "Failed to create spell")
			return
		}

		spell, _ = h.spellRepo.FindByID(spell.ID)
		respondJSON(w, http.StatusCreated, spell)
	}
}

// updateSpell updates an existing spell
func (h *spellHandler) updateSpell() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "spellID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid spell ID")
			return
		}

		spell, err := h.spellRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Spell not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != spell.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		var req UpdateSpellRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// If components change, re-run synthesis
		if len(req.ComponentIDs) > 0 {
			components, err := h.synthesisService.FetchComponents(req.ComponentIDs)
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch components for synthesis")
				respondError(w, http.StatusInternalServerError, "Failed to fetch components")
				return
			}

			synthesis := h.synthesisService.Synthesize(components)

			if len(synthesis.ValidationErrors) > 0 {
				respondJSON(w, http.StatusBadRequest, map[string]interface{}{
					"error":             "Component validation failed",
					"validation_errors": synthesis.ValidationErrors,
				})
				return
			}

			// Always recompute cost fields
			spell.ManaCost = synthesis.ManaCost
			spell.BaseTier = synthesis.BaseTier

			if err := h.spellRepo.ReplaceComponents(spell.ID, req.ComponentIDs); err != nil {
				log.Error().Err(err).Msg("Failed to update spell components")
				respondError(w, http.StatusInternalServerError, "Failed to update spell components")
				return
			}
		}

		if req.Name != nil {
			spell.Name = *req.Name
		}
		if req.Description != nil {
			spell.Description = *req.Description
		}
		if req.SlotLevel != nil {
			spell.SlotLevel = *req.SlotLevel
		}
		if req.Type != nil {
			spell.Type = *req.Type
		}
		if req.Range != nil {
			spell.Range = req.Range
		}
		if req.Duration != nil {
			spell.Duration = req.Duration
		}
		if req.Concentration != nil {
			spell.Concentration = *req.Concentration
		}
		if req.SaveAttr != nil {
			spell.SaveAttr = req.SaveAttr
		}
		if req.DamageDice != nil {
			spell.DamageDice = req.DamageDice
		}
		if req.DamageType != nil {
			spell.DamageType = req.DamageType
		}
		if req.AddModifier != nil {
			spell.AddModifier = *req.AddModifier
		}

		if err := h.spellRepo.Update(spell); err != nil {
			log.Error().Err(err).Msg("Failed to update spell")
			respondError(w, http.StatusInternalServerError, "Failed to update spell")
			return
		}

		spell, _ = h.spellRepo.FindByID(spell.ID)
		respondJSON(w, http.StatusOK, spell)
	}
}

// deleteSpell deletes a spell
func (h *spellHandler) deleteSpell() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "spellID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid spell ID")
			return
		}

		spell, err := h.spellRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Spell not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != spell.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		if err := h.spellRepo.Delete(id); err != nil {
			log.Error().Err(err).Msg("Failed to delete spell")
			respondError(w, http.StatusInternalServerError, "Failed to delete spell")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Spell deleted successfully"})
	}
}

// getSpellExecution returns the interpreted execution details for a spell
func (h *spellHandler) getSpellExecution() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "spellID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid spell ID")
			return
		}

		spell, err := h.spellRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Spell not found")
			return
		}

		execution, err := h.interpreterService.InterpretModels(spell.Components)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to interpret spell")
			return
		}

		respondJSON(w, http.StatusOK, execution)
	}
}
