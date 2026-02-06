package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rs/zerolog/log"
)

// getCreationOptions returns all races and classes with full details for the creation wizard
func (h *characterHandler) getCreationOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		races, err := h.raceRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get races")
			respondError(w, http.StatusInternalServerError, "Failed to get races")
			return
		}

		classes, err := h.classRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get classes")
			respondError(w, http.StatusInternalServerError, "Failed to get classes")
			return
		}

		// Dereference pointers for response
		racesValue := make([]models.Race, len(races))
		for i, r := range races {
			racesValue[i] = *r
		}

		classesValue := make([]models.Class, len(classes))
		for i, c := range classes {
			classesValue[i] = *c
		}

		resp := CreationOptionsResponse{
			Races:     racesValue,
			Classes:   classesValue,
			PointsMax: 27, // Standard D&D 5e point buy
		}
		respondJSON(w, http.StatusOK, resp)
	}
}

// createCharacter creates a new character
func (h *characterHandler) createCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateCharacterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Name == "" {
			respondError(w, http.StatusBadRequest, "Name is required")
			return
		}

		// Validate race and class exist
		if _, err := h.raceRepo.FindByID(req.RaceID); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid race_id")
			return
		}
		class, err := h.classRepo.FindByID(req.ClassID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid class_id")
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

		character := &models.Character{
			UserID:             req.UserID,
			Name:               req.Name,
			RaceID:             req.RaceID,
			LineageID:          req.LineageID,
			ClassID:            req.ClassID,
			Level:              req.Level,
			SpellbookIDs:       pq.StringArray(req.Spellbook),
			Strength:           req.Strength,
			Dexterity:          req.Dexterity,
			Constitution:       req.Constitution,
			Intelligence:       req.Intelligence,
			Wisdom:             req.Wisdom,
			Charisma:           req.Charisma,
			CurrentSpellPoints: req.CurrentSpellPoints,
			Money:              req.Money,
		}

		if character.Level == 0 {
			character.Level = 1
		}

		// Validate Skill Proficiencies
		if len(req.SkillProficiencies) > 0 {
			validSkills := make(map[string]bool)
			for _, s := range class.SkillChoice {
				validSkills[strings.ToLower(s)] = true
			}

			count := 0
			for _, s := range req.SkillProficiencies {
				if validSkills[strings.ToLower(s)] {
					count++
				} else {
					respondError(w, http.StatusBadRequest, "Invalid skill proficiency: "+s)
					return
				}
			}

			maxSkills := class.SkillChoiceCount
			if maxSkills == 0 {
				maxSkills = 2 // Default D&D standard
			}

			if count > maxSkills {
				respondError(w, http.StatusBadRequest, "Too many skill proficiencies selected")
				return
			}
		}

		// Calculate initial HP for level 1: HitDie (max) + CON modifier
		conMod := abilityMod(character.Constitution)
		initialHP := class.HitDie + conMod
		if initialHP < 1 {
			initialHP = 1 // Minimum 1 HP
		}
		character.MaxHP = initialHP
		character.CurrentHP = initialHP

		// Initialize class-specific resources
		if class.ResourceType == "stability" {
			// Find level 1 max stability
			if cl, err := h.classRepo.FindLevelByClassAndLevel(class.ID, 1); err == nil {
				character.CurrentStability = cl.MaxStability
			}
		} else if class.ResourceType == "blood_ichor" {
			character.CurrentBloodIchor = h.resourceService.ComputeMaxBloodIchor(character)
		}

		// Process Equipment
		var inventory []string
		var weaponIDs []uuid.UUID
		var itemIDs []uuid.UUID

		// 1. Add automatic equipment from class
		if len(class.StartingEquip) > 0 {
			inventory = append(inventory, class.StartingEquip...)
		}
		for _, wID := range class.StartingWeaponIDs {
			if id, err := uuid.Parse(wID); err == nil {
				weaponIDs = append(weaponIDs, id)
			}
		}
		for _, iID := range class.StartingItemIDs {
			if id, err := uuid.Parse(iID); err == nil {
				itemIDs = append(itemIDs, id)
			}
		}

		// 2. Add selected equipment choices
		if len(req.EquipmentChoices) > 0 {
			options, err := h.classRepo.FindEquipmentOptionsByIDs(req.EquipmentChoices)
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch equipment options")
				respondError(w, http.StatusInternalServerError, "Failed to fetch equipment options")
				return
			}
			for _, opt := range options {
				if len(opt.Items) > 0 {
					inventory = append(inventory, opt.Items...)
				}
				for _, wID := range opt.WeaponIDs {
					if id, err := uuid.Parse(wID); err == nil {
						weaponIDs = append(weaponIDs, id)
					}
				}
				for _, iID := range opt.ItemIDs {
					if id, err := uuid.Parse(iID); err == nil {
						itemIDs = append(itemIDs, id)
					}
				}
			}
		}

		character.Inventory = pq.StringArray(inventory)

		// Fetch weapon and item objects to associate
		if len(weaponIDs) > 0 {
			var weapons []models.Weapon
			if err := h.characterRepo.GetDB().Where("id IN ?", weaponIDs).Find(&weapons).Error; err == nil {
				character.Weapons = weapons
			}
		}
		if len(itemIDs) > 0 {
			var items []models.Item
			if err := h.characterRepo.GetDB().Where("id IN ?", itemIDs).Find(&items).Error; err == nil {
				character.Items = items
			}
		}

		if err := h.characterRepo.Add(character); err != nil {
			log.Error().Err(err).Msg("Failed to create character")
			respondError(w, http.StatusInternalServerError, "Failed to create character")
			return
		}

		// Handle Primary Weapon and Modifiers (v2)
		if req.PrimaryWeaponID != nil {
			var cw models.CharacterWeapon
			if err := h.characterRepo.GetDB().Table("character_weapons_v2").
				Where("character_id = ? AND weapon_id = ?", character.ID, *req.PrimaryWeaponID).
				First(&cw).Error; err == nil {
				
				cw.IsPrimary = true
				h.characterRepo.GetDB().Table("character_weapons_v2").Save(&cw)

				// Check for level 1 weapon requirements
				weaponReq, _ := h.classRepo.FindWeaponRequirementByClassAndLevel(character.ClassID, 1)
				if weaponReq != nil {
					modifier := models.WeaponModifier{
						CharacterWeaponID: cw.ID,
						ModifierType:      weaponReq.ModifierType,
						IsPermanent:       true,
						IsActive:          true,
					}
					if weaponReq.ModifierType == models.ModifierTypePistonCore {
						modifier.Metadata = models.PistonCoreMetadata(character.Level)
					}
					h.characterRepo.GetDB().Create(&modifier)
				}
			}
		}

		if len(req.SkillProficiencies) > 0 {
			if err := h.characterRepo.ReplaceSkillProficiencies(character.ID, req.SkillProficiencies); err != nil {
				log.Error().Err(err).Msg("Failed to save skill proficiencies")
				respondError(w, http.StatusInternalServerError, "Failed to save skill proficiencies")
				return
			}
		}

		respondJSON(w, http.StatusCreated, character)
	}
}
