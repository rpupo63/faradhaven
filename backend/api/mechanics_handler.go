package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type MechanicsHandler struct {
	db *gorm.DB
}

func NewMechanicsHandler(db *gorm.DB) *MechanicsHandler {
	return &MechanicsHandler{db: db}
}

// RollTable rolls on a specific effect table and applies it to the character
func (h *MechanicsHandler) RollTable(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid character ID")
		return
	}

	var req RollTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 1. Determine category based on table name
	var category string
	var source string
	switch req.TableName {
	case "mutagen_feral":
		category = "Mutagen Feral"
		source = "Feral Mode (Madness Check Failure)"
	case "lorewright_madness":
		category = "Lorewright Madness"
		source = "Echo Usage (Madness Roll 1)"
	default:
		respondError(w, http.StatusBadRequest, "Invalid table name")
		return
	}

	// 2. Fetch all effects in that category
	var effects []models.Effect
	if err := h.db.Where("category = ?", category).Find(&effects).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch effects")
		return
	}

	if len(effects) == 0 {
		respondError(w, http.StatusNotFound, "No effects found for this table")
		return
	}

	// 3. Roll a random effect
	// Note: In a real seeded environment, effects names are like "Feral Mutation (Roll 1)", "Feral Mutation (Roll 2)"
	// We can either parse the name or just pick randomly from the list.
	// Since the seed uses deterministic UUIDs and order might vary, random pick from the list is safer and simpler.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rng.Intn(len(effects))
	selectedEffect := effects[randomIndex]

	// Extract the "Roll X" number if possible for display
	rollVal := randomIndex + 1
	var rollStr int
	_, err = fmt.Sscanf(selectedEffect.Name, "%s %s (Roll %d)", &category, &source, &rollStr)
	if err == nil {
		rollVal = rollStr
	}

	// 4. Apply to character (create CharacterEffect)
	charEffect := models.CharacterEffect{
		CharacterID: characterID,
		EffectID:    selectedEffect.ID,
		Duration:    "1 Minute", // Default for these tables usually
		Source:      source,
	}

	if err := h.db.Create(&charEffect).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to apply effect")
		return
	}

	respondJSON(w, http.StatusOK, RollTableResponse{
		Effect: &selectedEffect,
		Roll:   rollVal,
	})
}

// MutagenCast handles casting for the Mutagen class with threshold-based madness logic.
// A madness save is required if: (1) component count > proficiency bonus, or
// (2) total casts today > proficiency bonus + CON modifier.
// When a save IS required, the Madness DC increases by +2 (or +1 at level 6+).
func (h *MechanicsHandler) MutagenCast(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid character ID")
		return
	}

	var req MutagenCastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var character models.Character
	if err := h.db.First(&character, "id = ?", characterID).Error; err != nil {
		respondError(w, http.StatusNotFound, "Character not found")
		return
	}

	// Fetch class level for proficiency bonus
	var classLevel models.ClassLevel
	if err := h.db.Where("class_id = ? AND level = ?", character.ClassID, character.Level).First(&classLevel).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch class level")
		return
	}

	profBonus := classLevel.ProficiencyBonus
	conMod := (character.Constitution - 10) / 2
	if character.Constitution < 10 && (character.Constitution-10)%2 != 0 {
		conMod--
	}

	// Increment cast count
	character.MadnessCastCount++

	// Determine if madness save is required
	requiresSave := false
	if req.ComponentCount > profBonus {
		requiresSave = true
	}
	if character.MadnessCastCount > profBonus+conMod {
		requiresSave = true
	}

	// If save is required, increment the DC
	if requiresSave {
		increment := 2
		if character.Level >= 6 {
			increment = 1
		}
		character.CurrentMadnessDC += increment
	}

	if err := h.db.Save(&character).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update character")
		return
	}

	respondJSON(w, http.StatusOK, MutagenCastResponse{
		MadnessCastCount: character.MadnessCastCount,
		CurrentDC:        character.CurrentMadnessDC,
		RequiresSave:     requiresSave,
		SaveDC:           character.CurrentMadnessDC,
	})
}

// GetActiveEffects lists all active effects on a character
func (h *MechanicsHandler) GetActiveEffects(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid character ID")
		return
	}

	var charEffects []models.CharacterEffect
	if err := h.db.Preload("Effect").Where("character_id = ?", characterID).Find(&charEffects).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch effects")
		return
	}

	response := make([]ActiveEffectResponse, len(charEffects))
	for i, ce := range charEffects {
		response[i] = ActiveEffectResponse{
			ID:          ce.ID,
			EffectID:    ce.EffectID,
			Name:        ce.Effect.Name,
			Description: ce.Effect.Description,
			Mechanics:   ce.Effect.Mechanics,
			Category:    ce.Effect.Category,
			Duration:    ce.Duration,
			Source:      ce.Source,
		}
	}

	respondJSON(w, http.StatusOK, response)
}

// RemoveEffect removes an active effect
func (h *MechanicsHandler) RemoveEffect(w http.ResponseWriter, r *http.Request) {
	effectID, err := uuid.Parse(chi.URLParam(r, "effectID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid effect ID")
		return
	}

	if err := h.db.Delete(&models.CharacterEffect{}, "id = ?", effectID).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to remove effect")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}