package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rs/zerolog/log"
)

// profBonusByLevel returns the proficiency bonus for a given character level (D&D standard).
func profBonusByLevel(level int) int {
	switch {
	case level <= 4:
		return 2
	case level <= 8:
		return 3
	case level <= 12:
		return 4
	case level <= 16:
		return 5
	default:
		return 6
	}
}

// initTraitUseResources creates or updates CharacterResource rows for traits with limited uses.
// Called after character creation and on level-up to keep max uses in sync with proficiency bonus.
func initTraitUseResources(
	characterID uuid.UUID,
	traits []models.Trait,
	profBonus int,
	characterLevel int,
	repo database.CharacterResourceRepository,
) error {
	for _, trait := range traits {
		if trait.UsesPerRest == "" {
			continue
		}
		var maxUses int
		switch trait.UsesPerRest {
		case "Proficiency Bonus":
			maxUses = profBonus
		case "1 per spell":
			if characterLevel >= 5 {
				maxUses = 2
			} else if characterLevel >= 3 {
				maxUses = 1
			} else {
				continue // not yet unlocked
			}
		default:
			n, err := strconv.Atoi(trait.UsesPerRest)
			if err != nil || n <= 0 {
				continue // Skip non-integer values
			}
			maxUses = n
		}

		key := "trait_" + trait.ID.String()
		existing, err := repo.FindByCharacterAndKey(characterID, key)
		if err != nil {
			// Doesn't exist — create it
			mv := maxUses
			resource := &models.CharacterResource{
				CharacterID:        characterID,
				ResourceKey:        key,
				ResourceName:       trait.Name + " Uses",
				CurrentValue:       maxUses,
				MaxValue:           &mv,
				RestoreOnLongRest:  strings.Contains(trait.ResetCondition, "Long"),
				RestoreOnShortRest: strings.Contains(trait.ResetCondition, "Short"),
			}
			if addErr := repo.Add(resource); addErr != nil {
				return fmt.Errorf("failed to init trait resource %s: %w", trait.Name, addErr)
			}
		} else {
			// Exists — update max if it changed (e.g. profBonus went up on level-up)
			if existing.MaxValue == nil || *existing.MaxValue != maxUses {
				mv := maxUses
				existing.MaxValue = &mv
				if updateErr := repo.Update(existing); updateErr != nil {
					return fmt.Errorf("failed to update trait resource %s: %w", trait.Name, updateErr)
				}
			}
		}
	}
	return nil
}

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
			Languages:          pq.StringArray(req.Languages),
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

		// Class-specific trackable resources are initialized via
		// resourceService.InitializeCharacterResources() after the character is persisted.

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

		// Initialize generic CharacterResource rows from class resource definitions
		if err := h.resourceService.InitializeCharacterResources(character); err != nil {
			log.Error().Err(err).Msg("Failed to initialize character resources")
			// Non-fatal: character was created successfully, resources can be backfilled later
		}

		// Align spell point column with class level max (CharacterResource.spell_points is the spend pool for Rift Weaver)
		if cl, err := h.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level); err == nil && cl != nil {
			character.CurrentSpellPoints = cl.MaxSpellPoints
			if err := h.characterRepo.Update(character); err != nil {
				log.Error().Err(err).Msg("Failed to sync current_spell_points after resource init")
			}
		}

		// Initialize trait use resources for race traits with limited uses
		race, raceErr := h.raceRepo.FindByIDWithTraits(character.RaceID)
		if raceErr == nil {
			profBonus := profBonusByLevel(character.Level)
			allTraits := race.Traits
			if character.LineageID != nil {
				if lineage, lineageErr := h.raceRepo.FindLineageByIDWithTraits(*character.LineageID); lineageErr == nil {
					allTraits = append(allTraits, lineage.LineageTraits...)
				}
			}
			if traitErr := initTraitUseResources(character.ID, allTraits, profBonus, character.Level, h.characterResourceRepo); traitErr != nil {
				log.Error().Err(traitErr).Msg("Failed to initialize trait use resources")
				// Non-fatal
			}
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
