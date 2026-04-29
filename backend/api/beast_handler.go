package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rs/zerolog/log"
)

type beastHandler struct {
	beastRepo  database.BeastRepository
	attackRepo database.AttackRepository
}

func newBeastHandler(beastRepo database.BeastRepository, attackRepo database.AttackRepository) *beastHandler {
	return &beastHandler{beastRepo: beastRepo, attackRepo: attackRepo}
}

// getAllBeasts returns all beasts
func (h *beastHandler) getAllBeasts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		beasts, err := h.beastRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get beasts")
			respondError(w, http.StatusInternalServerError, "Failed to get beasts")
			return
		}
		respondJSON(w, http.StatusOK, beasts)
	}
}

// getBeast returns a beast by ID with its attacks
func (h *beastHandler) getBeast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "beastID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid beast ID")
			return
		}

		beast, err := h.beastRepo.FindByIDWithAttacks(id)
		if err != nil {
			log.Error().Err(err).Str("beastID", idStr).Msg("Failed to get beast")
			respondError(w, http.StatusNotFound, "Beast not found")
			return
		}
		respondJSON(w, http.StatusOK, beast)
	}
}

// getBeastsByUser returns all beasts for a user with their attacks
func (h *beastHandler) getBeastsByUser() http.HandlerFunc {
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

		beasts, err := h.beastRepo.FindByUserIDWithAttacks(userID)
		if err != nil {
			log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get beasts")
			respondError(w, http.StatusInternalServerError, "Failed to get beasts")
			return
		}
		respondJSON(w, http.StatusOK, beasts)
	}
}

// createBeast creates a new beast with its attacks
func (h *beastHandler) createBeast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateBeastRequest
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

		beast := &models.Beast{
			UserID:           req.UserID,
			Name:             req.Name,
			ImageURL:         req.ImageURL,
			Size:             req.Size,
			Type:             req.Type,
			Alignment:        req.Alignment,
			ArmorClass:       req.ArmorClass,
			HitPoints:        req.HitPoints,
			HitDice:          req.HitDice,
			Speed:            req.Speed,
			Strength:         req.Strength,
			Dexterity:        req.Dexterity,
			Constitution:     req.Constitution,
			Intelligence:     req.Intelligence,
			Wisdom:           req.Wisdom,
			Charisma:         req.Charisma,
			ChallengeRating:  req.ChallengeRating,
			Abilities:        req.Abilities,
			Actions:          req.Actions,
			LegendaryActions: req.LegendaryActions,
			Description:      req.Description,
		}

		// Convert attack requests to attack models
		for _, attackReq := range req.Attacks {
			attack := models.Attack{
				Name:        attackReq.Name,
				AttackBonus: attackReq.AttackBonus,
				DamageType:  attackReq.DamageType,
				DamageDice:  attackReq.DamageDice,
				Reach:       attackReq.Reach,
				Description: attackReq.Description,
			}
			beast.Attacks = append(beast.Attacks, attack)
		}

		if err := h.beastRepo.Add(beast); err != nil {
			log.Error().Err(err).Msg("Failed to create beast")
			respondError(w, http.StatusInternalServerError, "Failed to create beast")
			return
		}

		respondJSON(w, http.StatusCreated, beast)
	}
}

// updateBeast updates an existing beast
func (h *beastHandler) updateBeast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "beastID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid beast ID")
			return
		}

		beast, err := h.beastRepo.FindByIDWithAttacks(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Beast not found")
			return
		}

		// Verify user owns this beast
		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != beast.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		var req UpdateBeastRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Name != nil {
			beast.Name = *req.Name
		}
		if req.ImageURL != nil {
			beast.ImageURL = req.ImageURL
		}
		if req.Size != nil {
			beast.Size = *req.Size
		}
		if req.Type != nil {
			beast.Type = *req.Type
		}
		if req.Alignment != nil {
			beast.Alignment = *req.Alignment
		}
		if req.ArmorClass != nil {
			beast.ArmorClass = *req.ArmorClass
		}
		if req.HitPoints != nil {
			beast.HitPoints = *req.HitPoints
		}
		if req.HitDice != nil {
			beast.HitDice = *req.HitDice
		}
		if req.Speed != nil {
			beast.Speed = *req.Speed
		}
		if req.Strength != nil {
			beast.Strength = *req.Strength
		}
		if req.Dexterity != nil {
			beast.Dexterity = *req.Dexterity
		}
		if req.Constitution != nil {
			beast.Constitution = *req.Constitution
		}
		if req.Intelligence != nil {
			beast.Intelligence = *req.Intelligence
		}
		if req.Wisdom != nil {
			beast.Wisdom = *req.Wisdom
		}
		if req.Charisma != nil {
			beast.Charisma = *req.Charisma
		}
		if req.ChallengeRating != nil {
			beast.ChallengeRating = *req.ChallengeRating
		}
		if req.Abilities != nil {
			beast.Abilities = req.Abilities
		}
		if req.Actions != nil {
			beast.Actions = req.Actions
		}
		if req.LegendaryActions != nil {
			beast.LegendaryActions = req.LegendaryActions
		}
		if req.Description != nil {
			beast.Description = *req.Description
		}

		if err := h.beastRepo.Update(beast); err != nil {
			log.Error().Err(err).Msg("Failed to update beast")
			respondError(w, http.StatusInternalServerError, "Failed to update beast")
			return
		}

		respondJSON(w, http.StatusOK, beast)
	}
}

// deleteBeast deletes a beast
func (h *beastHandler) deleteBeast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "beastID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid beast ID")
			return
		}

		beast, err := h.beastRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Beast not found")
			return
		}

		// Verify user owns this beast
		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != beast.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		if err := h.beastRepo.Delete(id); err != nil {
			log.Error().Err(err).Msg("Failed to delete beast")
			respondError(w, http.StatusInternalServerError, "Failed to delete beast")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Beast deleted successfully"})
	}
}
