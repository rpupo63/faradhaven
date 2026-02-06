package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type characterHandler struct {
	characterRepo   database.CharacterRepository
	raceRepo        database.RaceRepository
	classRepo       database.ClassRepository
	itemRepo        database.ItemRepository
	weaponRepo      database.WeaponRepository
	resourceService *services.ResourceService
}

func newCharacterHandler(
	characterRepo database.CharacterRepository,
	raceRepo database.RaceRepository,
	classRepo database.ClassRepository,
	itemRepo database.ItemRepository,
	weaponRepo database.WeaponRepository,
	resourceService *services.ResourceService,
) *characterHandler {
	return &characterHandler{
		characterRepo:   characterRepo,
		raceRepo:        raceRepo,
		classRepo:       classRepo,
		itemRepo:        itemRepo,
		weaponRepo:      weaponRepo,
		resourceService: resourceService,
	}
}

// abilityMod returns (score - 10) / 2 rounded down
func abilityMod(score int) int {
	return int(math.Floor(float64(score-10) / 2))
}

// primaryAbilityMod returns the modifier for the class primary ability from the character
func primaryAbilityMod(c *models.Character, primaryAbility string) int {
	switch strings.ToLower(primaryAbility) {
	case "strength":
		return abilityMod(c.Strength)
	case "dexterity":
		return abilityMod(c.Dexterity)
	case "constitution":
		return abilityMod(c.Constitution)
	case "intelligence":
		return abilityMod(c.Intelligence)
	case "wisdom":
		return abilityMod(c.Wisdom)
	case "charisma":
		return abilityMod(c.Charisma)
	default:
		return abilityMod(c.Intelligence)
	}
}

// getAllCharacters returns all characters
func (h *characterHandler) getAllCharacters() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characters, err := h.characterRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get characters")
			respondError(w, http.StatusInternalServerError, "Failed to get characters")
			return
		}
		respondJSON(w, http.StatusOK, characters)
	}
}

// getCharacter returns a character by ID
func (h *characterHandler) getCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByID(id)
		if err != nil {
			log.Error().Err(err).Str("characterID", idStr).Msg("Failed to get character")
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}
		respondJSON(w, http.StatusOK, character)
	}
}

// restSpellPoints resets CurrentSpellPoints to MaxSpellPoints (short/long rest)
func (h *characterHandler) restSpellPoints() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != character.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		classLevel, err := h.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
		if err != nil {
			respondError(w, http.StatusNotFound, "Level data not found")
			return
		}

		character.CurrentSpellPoints = classLevel.MaxSpellPoints
		if err := h.characterRepo.Update(character); err != nil {
			log.Error().Err(err).Msg("Failed to reset spell points")
			respondError(w, http.StatusInternalServerError, "Failed to reset spell points")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"current_spell_points": character.CurrentSpellPoints,
			"max_spell_points":     classLevel.MaxSpellPoints,
		})
	}
}

// getCharactersByUser returns all characters for a user
func (h *characterHandler) getCharactersByUser() http.HandlerFunc {
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

		characters, err := h.characterRepo.FindByUserID(userID)
		if err != nil {
			log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get characters")
			respondError(w, http.StatusInternalServerError, "Failed to get characters")
			return
		}
		respondJSON(w, http.StatusOK, characters)
	}
}

// updateCharacter updates an existing character
func (h *characterHandler) updateCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		// Verify user owns this character
		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != character.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		var req UpdateCharacterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Name != nil {
			character.Name = *req.Name
		}
		if req.RaceID != nil {
			if _, err := h.raceRepo.FindByID(*req.RaceID); err != nil {
				respondError(w, http.StatusBadRequest, "Invalid race_id")
				return
			}
			character.RaceID = *req.RaceID
		}
		if req.LineageID != nil {
			character.LineageID = req.LineageID
		}
		if req.ClassID != nil {
			if _, err := h.classRepo.FindByID(*req.ClassID); err != nil {
				respondError(w, http.StatusBadRequest, "Invalid class_id")
				return
			}
			character.ClassID = *req.ClassID
		}
		if req.Level != nil {
			character.Level = *req.Level
		}
		if req.Spellbook != nil {
			character.SpellbookIDs = req.Spellbook
		}
		if req.Strength != nil {
			character.Strength = *req.Strength
		}
		if req.Dexterity != nil {
			character.Dexterity = *req.Dexterity
		}
		if req.Constitution != nil {
			character.Constitution = *req.Constitution
		}
		if req.Intelligence != nil {
			character.Intelligence = *req.Intelligence
		}
		if req.Wisdom != nil {
			character.Wisdom = *req.Wisdom
		}
		if req.Charisma != nil {
			character.Charisma = *req.Charisma
		}
		if req.CurrentSpellPoints != nil {
			character.CurrentSpellPoints = *req.CurrentSpellPoints
		}
		if req.Money != nil {
			character.Money = *req.Money
		}
		if req.SkillProficiencies != nil {
			if err := h.characterRepo.ReplaceSkillProficiencies(character.ID, req.SkillProficiencies); err != nil {
				log.Error().Err(err).Msg("Failed to update skill proficiencies")
				respondError(w, http.StatusInternalServerError, "Failed to update skill proficiencies")
				return
			}
		}

		if err := h.characterRepo.Update(character); err != nil {
			log.Error().Err(err).Msg("Failed to update character")
			respondError(w, http.StatusInternalServerError, "Failed to update character")
			return
		}

		respondJSON(w, http.StatusOK, character)
	}
}

// deleteCharacter deletes a character
func (h *characterHandler) deleteCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		// Verify user owns this character
		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != character.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		if err := h.characterRepo.Delete(id); err != nil {
			log.Error().Err(err).Msg("Failed to delete character")
			respondError(w, http.StatusInternalServerError, "Failed to delete character")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Character deleted successfully"})
	}
}

// updateBackstory updates a character's backstory
func (h *characterHandler) updateBackstory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		// Verify user owns this character
		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != character.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		var req UpdateBackstoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		character.Backstory = req.Backstory
		if err := h.characterRepo.Update(character); err != nil {
			log.Error().Err(err).Msg("Failed to update backstory")
			respondError(w, http.StatusInternalServerError, "Failed to update backstory")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"message":   "Backstory updated successfully",
			"backstory": character.Backstory,
		})
	}
}

// getBackstory returns a character's backstory
func (h *characterHandler) getBackstory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"backstory": character.Backstory,
		})
	}
}

// parseCost converts a cost string (e.g. "50 gp") to copper pieces (int64)
func parseCost(costStr string) int64 {
	costStr = strings.TrimSpace(strings.ToLower(costStr))
	if costStr == "" {
		return 0
	}

	var value int64
	var unit string
	
	// Try to match "123 gp", "50 sp", etc.
	// Since we don't have regex readily available without importing regexp and compiling it,
	// we can do a simple scan if the format is consistent.
	// But regexp is safer. Let's use basic string splitting for now as it's faster and usually sufficient if data is clean.
	// "50 gp" -> ["50", "gp"]
	
	parts := strings.Fields(costStr)
	if len(parts) < 1 {
		return 0
	}

	// Parse number
	// valFloat, _ := map[string]interface{}{"val": 0}.(map[string]interface{}) // dummy
	// _ = valFloat
	
	// We'll use Sscanf
	n, err := fmt.Sscanf(parts[0], "%d", &value)
	if err != nil || n != 1 {
		return 0
	}
	
	if len(parts) > 1 {
		unit = parts[1]
	}

	switch unit {
	case "pp":
		return value * 1000
	case "gp":
		return value * 100
	case "sp":
		return value * 10
	case "cp":
		return value
	default:
		// Default to cp if no unit or unknown unit? Or GP? 
		// Assuming CP as base is safest, but usually things without unit might be GP in some contexts.
		// Given the parser in frontend defaults to value if no match, let's assume value (cp).
		return value
	}
}

// purchaseItem handles purchasing an item or weapon
func (h *characterHandler) purchaseItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		charID, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		// Verify user owns this character
		character, err := h.characterRepo.FindByID(charID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != character.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		var req PurchaseItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		var cost int64
		var itemName string

		if req.ItemType == "weapon" {
			weapon, err := h.weaponRepo.FindByID(req.ItemID)
			if err != nil {
				respondError(w, http.StatusNotFound, "Weapon not found")
				return
			}
			cost = parseCost(weapon.Cost)
			itemName = weapon.Name

			// Add weapon to character
			// We need to use database association handling.
			// character.Weapons = append(character.Weapons, *weapon) -- this works if we save, but GORM needs care.
			// Better to insert into join table directly or use association mode.
			// Since we don't have handy association methods in repo, let's use the DB instance from repo if possible
			// or assume Update works with associations if preloaded? No, Update usually updates fields.
			// We should use the association API.
			
			if character.Money < cost {
				respondError(w, http.StatusBadRequest, "Insufficient funds")
				return
			}

			// Deduct money
			character.Money -= cost
			if err := h.characterRepo.Update(character); err != nil {
				log.Error().Err(err).Msg("Failed to update character money")
				respondError(w, http.StatusInternalServerError, "Failed to update wallet")
				return
			}

			// Add association
			if err := h.characterRepo.GetDB().Model(character).Association("Weapons").Append(weapon); err != nil {
				log.Error().Err(err).Msg("Failed to add weapon to character")
				// Try to refund? For now, just error.
				respondError(w, http.StatusInternalServerError, "Failed to add weapon")
				return
			}

		} else if req.ItemType == "item" {
			item, err := h.itemRepo.FindByID(req.ItemID)
			if err != nil {
				respondError(w, http.StatusNotFound, "Item not found")
				return
			}
			cost = parseCost(item.Cost)
			itemName = item.Name

			if character.Money < cost {
				respondError(w, http.StatusBadRequest, "Insufficient funds")
				return
			}

			character.Money -= cost
			if err := h.characterRepo.Update(character); err != nil {
				log.Error().Err(err).Msg("Failed to update character money")
				respondError(w, http.StatusInternalServerError, "Failed to update wallet")
				return
			}

			if err := h.characterRepo.GetDB().Model(character).Association("Items").Append(item); err != nil {
				log.Error().Err(err).Msg("Failed to add item to character")
				respondError(w, http.StatusInternalServerError, "Failed to add item")
				return
			}

		} else {
			respondError(w, http.StatusBadRequest, "Invalid item type")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message":      fmt.Sprintf("Purchased %s for %d cp", itemName, cost),
			"money":        character.Money,
			"cost_deducted": cost,
		})
	}
}
