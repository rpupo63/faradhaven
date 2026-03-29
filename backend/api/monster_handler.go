package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type monsterHandler struct {
	monsterRepo database.MonsterRepository
	monsterGenerationService *services.MonsterGenerationService
}

func newMonsterHandler(monsterRepo database.MonsterRepository, monsterGenerationService *services.MonsterGenerationService) *monsterHandler {
	return &monsterHandler{
		monsterRepo: monsterRepo,
		monsterGenerationService: monsterGenerationService,
	}
}

// CreateMonsterRequest represents the request body for creating a monster.
type CreateMonsterRequest struct {
	Description     string    `json:"description"`
	ChallengeRating string    `json:"challenge_rating"`
	UserID          uuid.UUID `json:"user_id"`
}

// CreateMonster creates a new monster using the generation service.
func (h *monsterHandler) createMonster() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateMonsterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Description == "" || req.ChallengeRating == "" {
			respondError(w, http.StatusBadRequest, "Description and Challenge Rating are required")
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

		monster, err := h.monsterGenerationService.GenerateMonsterFromPrompt(r.Context(), services.GenerateMonsterFromPromptRequest{
			UserID:          req.UserID,
			Description:     req.Description,
			ChallengeRating: req.ChallengeRating,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate monster")
			respondError(w, http.StatusInternalServerError, "Failed to generate monster")
			return
		}

		respondJSON(w, http.StatusCreated, monster)
	}
}

// GetMonster returns a monster by ID.
func (h *monsterHandler) getMonster() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "monsterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid monster ID")
			return
		}

		monster, err := h.monsterRepo.FindByID(id)
		if err != nil {
			log.Error().Err(err).Str("monsterID", idStr).Msg("Failed to get monster")
			respondError(w, http.StatusNotFound, "Monster not found")
			return
		}

		// Verify user owns this monster (or if it's a system monster)
		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if monster.UserID != uuid.Nil && authUserID != monster.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		respondJSON(w, http.StatusOK, monster)
	}
}

// GetMonstersByUser returns all monsters for a user.
func (h *monsterHandler) getMonstersByUser() http.HandlerFunc {
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

		monsters, err := h.monsterRepo.FindByUserID(userID)
		if err != nil {
			log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get monsters")
			respondError(w, http.StatusInternalServerError, "Failed to get monsters")
			return
		}
		respondJSON(w, http.StatusOK, monsters)
	}
}

// UpdateMonsterRequest represents the request body for updating a monster.
// All fields are pointers to allow partial updates.
type UpdateMonsterRequest struct {
	Name                *string               `json:"name"`
	Size                *models.CreatureSize  `json:"size"`
	Type                *models.CreatureType  `json:"type"`
	Alignment           *string    `json:"alignment"`
	ChallengeRating     *string    `json:"challenge_rating"`
	ProficiencyBonus    *int       `json:"proficiency_bonus"`
	ArmorClass          *int       `json:"armor_class"`
	HitPoints           *int       `json:"hit_points"`
	HitDice             *string    `json:"hit_dice"`
	Speed               *string    `json:"speed"`
	Strength            *int       `json:"strength"`
	Dexterity           *int       `json:"dexterity"`
	Constitution        *int       `json:"constitution"`
	Intelligence        *int       `json:"intelligence"`
	Wisdom              *int       `json:"wisdom"`
	Charisma            *int       `json:"charisma"`
	SavingThrows        *[]string  `json:"saving_throws"`
	Skills              *[]string  `json:"skills"`
	DamageImmunities    *[]string  `json:"damage_immunities"`
	DamageResistances   *[]string  `json:"damage_resistances"`
	ConditionImmunities *[]string  `json:"condition_immunities"`
	Senses              *string    `json:"senses"`
	Languages           *string    `json:"languages"`
	Notes               *string    `json:"notes"`
	ImageURL            *string    `json:"image_url"`
	SpecialTraits       *[]string  `json:"special_traits"`
	LegendaryActions    *[]string  `json:"legendary_actions"`
	Reactions           *[]string  `json:"reactions"`
	LairActions         *[]string  `json:"lair_actions"`
	Environments        *[]string  `json:"environments"`
	Source              *string    `json:"source"`
	VisualDescription   *string    `json:"visual_description"`
	// Attacks and Actions updates might require more complex logic (e.g., replace all, update specific)
	// For simplicity, let's assume they are fully replaced if provided.
	Attacks             *[]models.MonsterAttack `json:"attacks"`
	Actions             *[]models.MonsterAction `json:"actions"`
}

// UpdateMonster updates an existing monster.
func (h *monsterHandler) updateMonster() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "monsterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid monster ID")
			return
		}

		monster, err := h.monsterRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Monster not found")
			return
		}

		// Verify user owns this monster
		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != monster.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		var req UpdateMonsterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Apply updates
		if req.Name != nil { monster.Name = *req.Name }
		if req.Size != nil {
			sz, ok := models.ParseCreatureSize(strings.TrimSpace(string(*req.Size)))
			if !ok {
				respondError(w, http.StatusBadRequest, "invalid size")
				return
			}
			monster.Size = sz
		}
		if req.Type != nil {
			ty, ok := models.ParseCreatureType(strings.TrimSpace(string(*req.Type)))
			if !ok {
				respondError(w, http.StatusBadRequest, "invalid type")
				return
			}
			monster.Type = ty
		}
		if req.Alignment != nil { monster.Alignment = *req.Alignment }
		if req.ChallengeRating != nil { monster.ChallengeRating = *req.ChallengeRating }
		if req.ProficiencyBonus != nil { monster.ProficiencyBonus = *req.ProficiencyBonus }
		if req.ArmorClass != nil { monster.ArmorClass = *req.ArmorClass }
		if req.HitPoints != nil { monster.HitPoints = *req.HitPoints }
		if req.HitDice != nil { monster.HitDice = *req.HitDice }
		if req.Speed != nil { monster.Speed = *req.Speed }
		if req.Strength != nil { monster.Strength = *req.Strength }
		if req.Dexterity != nil { monster.Dexterity = *req.Dexterity }
		if req.Constitution != nil { monster.Constitution = *req.Constitution }
		if req.Intelligence != nil { monster.Intelligence = *req.Intelligence }
		if req.Wisdom != nil { monster.Wisdom = *req.Wisdom }
		if req.Charisma != nil { monster.Charisma = *req.Charisma }
		if req.SavingThrows != nil { monster.SavingThrows = *req.SavingThrows }
		if req.Skills != nil { monster.Skills = *req.Skills }
		if req.DamageImmunities != nil { monster.DamageImmunities = *req.DamageImmunities }
		if req.DamageResistances != nil { monster.DamageResistances = *req.DamageResistances }
		if req.ConditionImmunities != nil { monster.ConditionImmunities = *req.ConditionImmunities }
		if req.Senses != nil { monster.Senses = *req.Senses }
		if req.Languages != nil { monster.Languages = *req.Languages }
		if req.Notes != nil { monster.Notes = *req.Notes }
		if req.ImageURL != nil { monster.ImageURL = req.ImageURL }
		if req.SpecialTraits != nil { monster.SpecialTraits = *req.SpecialTraits }
		if req.LegendaryActions != nil { monster.LegendaryActions = *req.LegendaryActions }
		if req.Reactions != nil { monster.Reactions = *req.Reactions }
		if req.LairActions != nil { monster.LairActions = *req.LairActions }
		if req.Environments != nil { monster.Environments = *req.Environments }
		if req.Source != nil { monster.Source = *req.Source }
		if req.VisualDescription != nil { monster.VisualDescription = *req.VisualDescription }

		// Handle nested updates for Attacks and Actions (assuming full replacement for now)
		if req.Attacks != nil {
			monster.Attacks = *req.Attacks
		}
		if req.Actions != nil {
			monster.Actions = *req.Actions
		}

		if err := h.monsterRepo.Update(monster); err != nil {
			log.Error().Err(err).Msg("Failed to update monster")
			respondError(w, http.StatusInternalServerError, "Failed to update monster")
			return
		}

		respondJSON(w, http.StatusOK, monster)
	}
}

// DeleteMonster deletes a monster.
func (h *monsterHandler) deleteMonster() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "monsterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid monster ID")
			return
		}

		monster, err := h.monsterRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Monster not found")
			return
		}

		// Verify user owns this monster
		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != monster.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		if err := h.monsterRepo.Delete(id); err != nil {
			log.Error().Err(err).Msg("Failed to delete monster")
			respondError(w, http.StatusInternalServerError, "Failed to delete monster")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Monster deleted successfully"})
	}
}
