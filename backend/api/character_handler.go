package api

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type characterHandler struct {
	characterRepo         database.CharacterRepository
	raceRepo              database.RaceRepository
	classRepo             database.ClassRepository
	characterResourceRepo database.CharacterResourceRepository
	itemRepo              database.ItemRepository
	weaponRepo            database.WeaponRepository
	spellRepo             database.SpellRepository
	resourceService       *services.ResourceService
	notorietyService      services.NotorietyService
	s3Service             *services.S3Service
	componentInterpreter  *services.ComponentInterpreterService
}

func newCharacterHandler(
	characterRepo database.CharacterRepository,
	raceRepo database.RaceRepository,
	classRepo database.ClassRepository,
	characterResourceRepo database.CharacterResourceRepository,
	itemRepo database.ItemRepository,
	weaponRepo database.WeaponRepository,
	spellRepo database.SpellRepository,
	resourceService *services.ResourceService,
	notorietyService services.NotorietyService,
	s3Service *services.S3Service,
	componentInterpreter *services.ComponentInterpreterService,
) *characterHandler {
	return &characterHandler{
		characterRepo:         characterRepo,
		raceRepo:              raceRepo,
		classRepo:             classRepo,
		characterResourceRepo: characterResourceRepo,
		itemRepo:              itemRepo,
		weaponRepo:            weaponRepo,
		spellRepo:             spellRepo,
		resourceService:       resourceService,
		notorietyService:      notorietyService,
		s3Service:             s3Service,
		componentInterpreter:  componentInterpreter,
	}
}

// buildClassResources aggregates resource definitions, level values, and character state
// into a response-ready slice of ClassResourceResponse.
func (h *characterHandler) buildClassResources(classID uuid.UUID, level int, characterID uuid.UUID) []ClassResourceResponse {
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

// uploadProfilePicture handles uploading a character image to S3
func (h *characterHandler) uploadProfilePicture() http.HandlerFunc {
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

		// Parse multipart form
		// Limit upload size to 10MB
		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("image")
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid file")
			return
		}
		defer file.Close()

		// Validate content type
		contentType := handler.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
			respondError(w, http.StatusBadRequest, "Invalid file type. Only JPEG, PNG, and WebP are allowed.")
			return
		}

		url, err := h.s3Service.UploadFile(file, handler)
		if err != nil {
			log.Error().Err(err).Msg("Failed to upload image")
			respondError(w, http.StatusInternalServerError, "Failed to upload image")
			return
		}

		character.ImageURL = url
		if err := h.characterRepo.Update(character); err != nil {
			log.Error().Err(err).Msg("Failed to update character image URL")
			respondError(w, http.StatusInternalServerError, "Failed to update character")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"image_url": url,
			"message":   "Profile picture updated successfully",
		})
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
		if req.Notoriety != nil {
			character.SanguineNotoriety = *req.Notoriety
		}
		if req.Notes != nil {
			character.Notes = *req.Notes
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
		if req.BackstoryHexColor != nil {
			character.BackstoryHexColor = *req.BackstoryHexColor
		}

		if err := h.characterRepo.Update(character); err != nil {
			log.Error().Err(err).Msg("Failed to update backstory")
			respondError(w, http.StatusInternalServerError, "Failed to update backstory")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"message":             "Backstory updated successfully",
			"backstory":           character.Backstory,
			"backstory_hex_color": character.BackstoryHexColor,
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
			"backstory":           character.Backstory,
			"backstory_hex_color": character.BackstoryHexColor,
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
			"message":       fmt.Sprintf("Purchased %s for %d cp", itemName, cost),
			"money":         character.Money,
			"cost_deducted": cost,
		})
	}
}

// castSpell handles deducting spell points and updating class-specific resources (like Madness)
func (h *characterHandler) castSpell() http.HandlerFunc {
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

		var req CastSpellRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Preload class with component pool for knowledge validation
		class, err := h.classRepo.FindByIDWithLevels(character.ClassID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to load class data")
			return
		}

		// Start transaction
		tx := h.characterRepo.GetDB().Begin()
		if tx.Error != nil {
			respondError(w, http.StatusInternalServerError, "Database error")
			return
		}
		defer tx.Rollback()

		// Build class component pool set for O(1) lookup
		classPoolIDs := make(map[uuid.UUID]bool, len(class.Components))
		for _, c := range class.Components {
			classPoolIDs[c.ID] = true
		}

		// 0. Handle Spell Components
		cost := 0
		if req.SpellID != nil {
			spell, err := h.spellRepo.FindByID(*req.SpellID)
			if err != nil {
				respondError(w, http.StatusNotFound, "Spell not found")
				return
			}

			// Validate: pool components pass freely; non-pool must be in inventory
			for _, comp := range spell.Components {
				if classPoolIDs[comp.ID] {
					continue
				}
				found := false
				for _, charComp := range character.Components {
					if charComp.ComponentID == comp.ID {
						if charComp.Count < 1 {
							respondError(w, http.StatusBadRequest, fmt.Sprintf("Not enough %s", comp.Name))
							return
						}
						found = true
						break
					}
				}
				if !found {
					respondError(w, http.StatusBadRequest, fmt.Sprintf("Character does not have component: %s", comp.Name))
					return
				}
			}

			// Consume only non-pool components from inventory
			for _, comp := range spell.Components {
				if classPoolIDs[comp.ID] {
					continue
				}
				for i, charComp := range character.Components {
					if charComp.ComponentID == comp.ID {
						character.Components[i].Count--
						if err := tx.Save(&character.Components[i]).Error; err != nil {
							log.Error().Err(err).Msg("Failed to update character component count")
							respondError(w, http.StatusInternalServerError, "Failed to update components")
							return
						}
						break
					}
				}
			}

			// Calculate Cost
			if class.Name == "The Rift Weaver" {
				// Rift Weaver: Cost is 2 SP per component
				cost = len(spell.Components) * 2
			} else {
				// Default: Cost is Spell Level
				if req.SpellLevel == 0 {
					cost = spell.SlotLevel
				} else {
					cost = req.SpellLevel
				}
			}
		} else {
			// Manual cast without ID
			cost = req.SpellLevel
		}

		// 1. Deduct Resources based on Resource Definitions
		var resourceKeyToDeduct string
		var finalCost int
		skipResourceDeduction := false

		// Determine which resource key to deduct from.
		if req.ResourceKey != nil && *req.ResourceKey != "" {
			resourceKeyToDeduct = *req.ResourceKey
		} else {
			switch class.Name {
			case "The Rift Weaver":
				resourceKeyToDeduct = "spell_points"
			case "The Piston Brawler":
				resourceKeyToDeduct = "max_stability"
			case "The Sanguinist":
				resourceKeyToDeduct = "max_blood_ichor"
			case "The Vapor Blade":
				resourceKeyToDeduct = "shadow_points"
			case "The Lorewright", "The Ironwright":
				// Components ARE the cost — no pool to deduct from
				skipResourceDeduction = true
			default:
				resourceKeyToDeduct = "spell_points"
			}
		}
		finalCost = cost

		// Handle Mutagen and Powder Mage as special cases
		if class.Name == "The Mutagen" {
			character.MadnessCastCount++
		} else if class.Name == "The Powder Mage" {
			spellResult, err := h.componentInterpreter.Interpret(req.Components)
			if err != nil {
				respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to interpret components: %v", err))
				return
			}
			tx.Commit()
			respondJSON(w, http.StatusOK, spellResult)
			return
		} else if !skipResourceDeduction && resourceKeyToDeduct != "" {
			err = h.resourceService.DeductResource(character.ID, resourceKeyToDeduct, finalCost)
			if err != nil {
				log.Error().Err(err).Str("resource_key", resourceKeyToDeduct).Msg("Failed to deduct resource")
				if strings.Contains(err.Error(), "insufficient") {
					respondError(w, http.StatusBadRequest, err.Error())
				} else {
					respondError(w, http.StatusInternalServerError, "Failed to deduct resource")
				}
				return
			}
		}

		// Update character using tx
		character.UpdatedAt = time.Now()
		if err := tx.Save(character).Error; err != nil {
			log.Error().Err(err).Msg("Failed to update character after casting")
			respondError(w, http.StatusInternalServerError, "Failed to update character")
			return
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			log.Error().Err(err).Msg("Failed to commit transaction")
			respondError(w, http.StatusInternalServerError, "Database commit failed")
			return
		}

		// Read updated resource values from the CharacterResource table (not stale model fields)
		response := map[string]interface{}{
			"current_spell_points": character.CurrentSpellPoints,
			"madness_cast_count":   character.MadnessCastCount,
		}
		if resourceKeyToDeduct != "" && !skipResourceDeduction {
			response["current_resource_value"] = h.resourceService.GetResourceValue(character.ID, resourceKeyToDeduct)
			response["resource_key"] = resourceKeyToDeduct
		}

		respondJSON(w, http.StatusOK, response)
	}
}

// consumeComponent handles manual consumption of a single component
func (h *characterHandler) consumeComponent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		charID, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		compIDStr := chi.URLParam(r, "componentID")
		compID, err := uuid.Parse(compIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid component ID")
			return
		}

		character, err := h.characterRepo.FindByID(charID)
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

		if err := h.characterRepo.UpdateComponentCount(charID, compID, -1); err != nil {
			log.Error().Err(err).Msg("Failed to consume component")
			respondError(w, http.StatusInternalServerError, "Failed to consume component")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Component consumed"})
	}
}

// extractComponents handles Sanguinist Sanguine Extraction — gains random class components
func (h *characterHandler) extractComponents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		charID, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

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

		if character.Class.Name != "The Sanguinist" {
			respondError(w, http.StatusBadRequest, "Only Sanguinists can use Sanguine Extraction")
			return
		}

		var req ExtractRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Load class with components
		class, err := h.classRepo.FindByIDWithLevels(character.ClassID)
		if err != nil || len(class.Components) == 0 {
			respondError(w, http.StatusInternalServerError, "Failed to load class components")
			return
		}

		// Calculate yield based on level
		yield := 1
		if character.Level >= 5 && character.Level <= 10 {
			yield = 2
		} else if character.Level >= 11 {
			yield = 3
		}

		// Level 15+: Master Extractions — 5 from allies (siphon), 3 from enemies (bite)
		if character.Level >= 15 && req.Source == "siphon" {
			yield = 5
		}

		// Pick N random components from the class component pool
		gained := make([]map[string]interface{}, 0, yield)
		for i := 0; i < yield; i++ {
			comp := class.Components[rand.Intn(len(class.Components))]
			if err := h.characterRepo.UpdateComponentCount(charID, comp.ID, 1); err != nil {
				log.Error().Err(err).Str("componentID", comp.ID.String()).Msg("Failed to add extracted component")
				continue
			}
			gained = append(gained, map[string]interface{}{
				"id":   comp.ID,
				"name": comp.Name,
			})
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message":    fmt.Sprintf("Extracted %d components via %s", len(gained), req.Source),
			"components": gained,
			"yield":      yield,
		})
	}
}

// updateNotoriety handles manual updates to notoriety by the DM/User
func (h *characterHandler) updateNotoriety() http.HandlerFunc {
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

		var req UpdateNotorietyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		character.SanguineNotoriety += req.Delta
		// Cap notoriety between -20 and 20
		if character.SanguineNotoriety > 20 {
			character.SanguineNotoriety = 20
		} else if character.SanguineNotoriety < -20 {
			character.SanguineNotoriety = -20
		}

		if err := h.characterRepo.Update(character); err != nil {
			log.Error().Err(err).Msg("Failed to update notoriety")
			respondError(w, http.StatusInternalServerError, "Failed to update notoriety")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message":   "Notoriety updated successfully",
			"notoriety": character.SanguineNotoriety,
		})
	}
}

func (h *characterHandler) updateSanguineNotorietyPoints() http.HandlerFunc {
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

		var req UpdateSanguineNotorietyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.notorietyService.UpdateNotoriety(id, req.MPChange, req.BRChange); err != nil {
			log.Error().Err(err).Msg("Failed to update sanguine notoriety points")
			respondError(w, http.StatusInternalServerError, "Failed to update notoriety points")
			return
		}

		// Refetch character to get updated values
		updatedCharacter, err := h.characterRepo.FindByID(id)
		if err != nil {
			log.Error().Err(err).Msg("Failed to refetch character after updating notoriety")
			respondError(w, http.StatusInternalServerError, "Failed to fetch updated character")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message":            "Sanguine notoriety points updated successfully",
			"sanguine_mp":        updatedCharacter.SanguineMP,
			"sanguine_br":        updatedCharacter.SanguineBR,
			"sanguine_notoriety": updatedCharacter.SanguineNotoriety,
		})
	}
}

// UpdateEquipmentRequest defines the structure for equipping/unequipping items.
type UpdateEquipmentRequest struct {
	ItemID   uuid.UUID `json:"item_id"`
	IsWeapon bool      `json:"is_weapon"`
	Equip    bool      `json:"equip"`
	Slot     string    `json:"slot"` // "armor", "shield", "weapon"
}

// weaponHandCost returns how many hands a weapon requires (2 for Two-Handed, 1 otherwise).
func weaponHandCost(properties []string) int {
	for _, p := range properties {
		if p == "Two-Handed" {
			return 2
		}
	}
	return 1
}

// countEquippedHands counts how many hand slots are in use by equipped weapons and shield.
func (h *characterHandler) countEquippedHands(charID uuid.UUID, hasShield bool) (int, error) {
	var equippedWeapons []models.CharacterWeapon
	if err := h.characterRepo.GetDB().
		Preload("Weapon").
		Where("character_id = ? AND is_equipped = true", charID).
		Find(&equippedWeapons).Error; err != nil {
		return 0, err
	}

	hands := 0
	for _, cw := range equippedWeapons {
		hands += weaponHandCost(cw.Weapon.Properties)
	}
	if hasShield {
		hands++
	}
	return hands, nil
}

// updateEquipment handles equipping and unequipping items and weapons.
// Enforces: 1 armor at a time, 2-hand limit for weapons + shield.
func (h *characterHandler) updateEquipment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		charIDStr := chi.URLParam(r, "characterID")
		charID, err := uuid.Parse(charIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

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

		var req UpdateEquipmentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.IsWeapon {
			if req.Equip {
				// Look up the weapon to determine hand cost
				weapon, err := h.weaponRepo.FindByID(req.ItemID)
				if err != nil {
					respondError(w, http.StatusNotFound, "Weapon not found")
					return
				}
				cost := weaponHandCost(weapon.Properties)

				usedHands, err := h.countEquippedHands(charID, character.EquippedShieldID != nil)
				if err != nil {
					log.Error().Err(err).Msg("Failed to count equipped hands")
					respondError(w, http.StatusInternalServerError, "Failed to check equipment")
					return
				}

				// Don't double-count this weapon if it's already equipped
				var alreadyEquipped models.CharacterWeapon
				if err := h.characterRepo.GetDB().
					Preload("Weapon").
					Where("character_id = ? AND weapon_id = ? AND is_equipped = true", charID, req.ItemID).
					First(&alreadyEquipped).Error; err == nil {
					usedHands -= weaponHandCost(alreadyEquipped.Weapon.Properties)
				}

				if usedHands+cost > 2 {
					respondError(w, http.StatusBadRequest, fmt.Sprintf(
						"Not enough free hands (need %d, have %d free)", cost, 2-usedHands))
					return
				}
			}

			err = h.characterRepo.GetDB().Model(&models.CharacterWeapon{}).
				Where("character_id = ? AND weapon_id = ?", charID, req.ItemID).
				Update("is_equipped", req.Equip).Error
			if err != nil {
				log.Error().Err(err).Msg("Failed to update weapon equipment status")
				respondError(w, http.StatusInternalServerError, "Failed to update weapon")
				return
			}
		} else {
			switch req.Slot {
			case "armor":
				if req.Equip {
					character.EquippedArmorID = &req.ItemID
				} else {
					character.EquippedArmorID = nil
				}

			case "shield":
				if req.Equip {
					usedHands, err := h.countEquippedHands(charID, false)
					if err != nil {
						log.Error().Err(err).Msg("Failed to count equipped hands")
						respondError(w, http.StatusInternalServerError, "Failed to check equipment")
						return
					}
					if usedHands+1 > 2 {
						respondError(w, http.StatusBadRequest, fmt.Sprintf(
							"Not enough free hands to equip shield (have %d free)", 2-usedHands))
						return
					}
					character.EquippedShieldID = &req.ItemID
				} else {
					character.EquippedShieldID = nil
				}

			default:
				respondError(w, http.StatusBadRequest, "Invalid equipment slot")
				return
			}

			if err := h.characterRepo.Update(character); err != nil {
				log.Error().Err(err).Msg("Failed to update character equipment")
				respondError(w, http.StatusInternalServerError, "Failed to update equipment")
				return
			}
		}

		// Return the updated character sheet
		sheet, err := h.getCharacterSheetData(character.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get updated character sheet")
			respondError(w, http.StatusInternalServerError, "Failed to retrieve updated character sheet")
			return
		}
		respondJSON(w, http.StatusOK, sheet)
	}
}
