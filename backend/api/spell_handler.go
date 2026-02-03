package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rs/zerolog/log"
)

type spellHandler struct {
	spellRepo database.SpellRepository
}

func newSpellHandler(spellRepo database.SpellRepository) *spellHandler {
	return &spellHandler{spellRepo: spellRepo}
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
		respondJSON(w, http.StatusOK, spells)
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

		// Verify user is accessing their own data
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

		// Verify user is creating for themselves
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
			UserID:      req.UserID,
			Name:        req.Name,
			Description: req.Description,
			SlotLevel:   req.SlotLevel,
		}

		if spell.SlotLevel == 0 {
			spell.SlotLevel = 1
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

		// Verify user owns this spell
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

		if req.Name != nil {
			spell.Name = *req.Name
		}
		if req.Description != nil {
			spell.Description = *req.Description
		}
		if req.ComponentIDs != nil {
			if err := h.spellRepo.ReplaceComponents(spell.ID, req.ComponentIDs); err != nil {
				log.Error().Err(err).Msg("Failed to update spell components")
				respondError(w, http.StatusInternalServerError, "Failed to update spell components")
				return
			}
		}
		if req.SlotLevel != nil {
			spell.SlotLevel = *req.SlotLevel
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

		// Verify user owns this spell
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
