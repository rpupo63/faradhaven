package api

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rs/zerolog/log"
)

type characterHandler struct {
	characterRepo database.CharacterRepository
	raceRepo      database.RaceRepository
	classRepo     database.ClassRepository
}

func newCharacterHandler(characterRepo database.CharacterRepository, raceRepo database.RaceRepository, classRepo database.ClassRepository) *characterHandler {
	return &characterHandler{characterRepo: characterRepo, raceRepo: raceRepo, classRepo: classRepo}
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

// getAllRaces returns all races (for character creation dropdowns)
func (h *characterHandler) getAllRaces() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		races, err := h.raceRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get races")
			respondError(w, http.StatusInternalServerError, "Failed to get races")
			return
		}
		respondJSON(w, http.StatusOK, races)
	}
}

// getRaceByID returns a race with traits (for D&D Beyond-style compendium detail view)
func (h *characterHandler) getRaceByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "raceID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid race ID")
			return
		}

		race, err := h.raceRepo.FindByIDWithTraits(id)
		if err != nil {
			log.Error().Err(err).Str("raceID", idStr).Msg("Failed to get race")
			respondError(w, http.StatusNotFound, "Race not found")
			return
		}

		respondJSON(w, http.StatusOK, race)
	}
}

// getAllClasses returns all classes (for character creation dropdowns)
func (h *characterHandler) getAllClasses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		classes, err := h.classRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get classes")
			respondError(w, http.StatusInternalServerError, "Failed to get classes")
			return
		}
		respondJSON(w, http.StatusOK, classes)
	}
}

// getClassByID returns a class with all levels 1-20 (for D&D Beyond-style compendium)
func (h *characterHandler) getClassByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "classID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid class ID")
			return
		}

		class, err := h.classRepo.FindByIDWithLevels(id)
		if err != nil {
			log.Error().Err(err).Str("classID", idStr).Msg("Failed to get class")
			respondError(w, http.StatusNotFound, "Class not found")
			return
		}

		resp := ClassWithLevelsResponse{
			Class:  class,
			Levels: class.Levels,
		}
		respondJSON(w, http.StatusOK, resp)
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

// getCharacterSheet returns the fully calculated character sheet (Class + ClassLevel joined)
// TotalHP = BaseHP + (AvgHitDie * (Level - 1)) + (ConMod * Level), BaseHP = HitDie + ConMod
// AC = 8 + ProficiencyBonus + DexMod, SaveDC = 8 + ProficiencyBonus + PrimaryAbilityMod
func (h *characterHandler) getCharacterSheet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByIDWithSkills(id)
		if err != nil {
			log.Error().Err(err).Str("characterID", idStr).Msg("Failed to get character")
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

		class, err := h.classRepo.FindByID(character.ClassID)
		if err != nil {
			log.Error().Err(err).Str("class_id", character.ClassID.String()).Msg("Class not found")
			respondError(w, http.StatusNotFound, "Class not found for character")
			return
		}

		classLevel, err := h.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
		if err != nil {
			log.Error().Err(err).Int("level", character.Level).Msg("ClassLevel not found")
			respondError(w, http.StatusNotFound, "Level data not found for class")
			return
		}

		conMod := abilityMod(character.Constitution)
		dexMod := abilityMod(character.Dexterity)
		primaryMod := primaryAbilityMod(character, class.PrimaryAbility)

		// TotalHP = BaseHP + (AvgHitDie * (Level - 1)) + (ConMod * Level); BaseHP = HitDie (max at level 1)
		avgHitDie := (class.HitDie + 1) / 2
		baseHP := class.HitDie
		totalHP := baseHP + (avgHitDie * (character.Level - 1)) + (conMod * character.Level)
		if character.Level < 1 {
			totalHP = baseHP + conMod
		}

		ac := 8 + classLevel.ProficiencyBonus + dexMod
		saveDC := 8 + classLevel.ProficiencyBonus + primaryMod

		maxSP := classLevel.MaxSpellPoints
		currentSP := character.CurrentSpellPoints
		if currentSP > maxSP {
			currentSP = maxSP
		}

		// Map class.SavingThrows (ability names) to ability ids for frontend
		saveProfs := make([]string, 0, len(class.SavingThrows))
		for _, name := range class.SavingThrows {
			saveProfs = append(saveProfs, strings.ToLower(strings.TrimSpace(name)))
		}

		sheet := CharacterSheetResponse{
			Character:                character,
			Class:                    class,
			ClassLevel:               classLevel,
			TotalHP:                  totalHP,
			AC:                       ac,
			SaveDC:                   saveDC,
			MaxSpellPoints:           maxSP,
			CurrentSpellPoints:       currentSP,
			SavingThrowProficiencies: saveProfs,
		}
		respondJSON(w, http.StatusOK, sheet)
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
		if _, err := h.classRepo.FindByID(req.ClassID); err != nil {
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
		}

		if character.Level == 0 {
			character.Level = 1
		}
		if character.Strength == 0 {
			character.Strength = 10
		}
		if character.Dexterity == 0 {
			character.Dexterity = 10
		}
		if character.Constitution == 0 {
			character.Constitution = 10
		}
		if character.Intelligence == 0 {
			character.Intelligence = 10
		}
		if character.Wisdom == 0 {
			character.Wisdom = 10
		}
		if character.Charisma == 0 {
			character.Charisma = 10
		}

		if err := h.characterRepo.Add(character); err != nil {
			log.Error().Err(err).Msg("Failed to create character")
			respondError(w, http.StatusInternalServerError, "Failed to create character")
			return
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
			character.SpellbookIDs = pq.StringArray(req.Spellbook)
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
