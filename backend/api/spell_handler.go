package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/errs"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

// parsePaginationParams parses page, limit, and level parameters from the request.
func parsePaginationParams(r *http.Request) (page, limit, levelFilter int, err error) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	levelStr := r.URL.Query().Get("level")

	page = 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return 0, 0, 0, errs.BadRequest("invalid page parameter")
		}
	}

	limit = 10 // Default limit
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return 0, 0, 0, errs.BadRequest("invalid limit parameter")
		}
	}

	levelFilter = 0 // 0 means no filter
	if levelStr != "" {
		levelFilter, err = strconv.Atoi(levelStr)
		if err != nil || levelFilter < 1 {
			return 0, 0, 0, errs.BadRequest("invalid level parameter")
		}
	}
	return page, limit, levelFilter, nil
}

// parseSpellbookPaginationParams is like parsePaginationParams but defaults limit to 20.
func parseSpellbookPaginationParams(r *http.Request) (page, limit, levelFilter int, err error) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	levelStr := r.URL.Query().Get("level")

	page = 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return 0, 0, 0, errs.BadRequest("invalid page parameter")
		}
	}

	limit = 20
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return 0, 0, 0, errs.BadRequest("invalid limit parameter")
		}
	}

	levelFilter = 0
	if levelStr != "" {
		levelFilter, err = strconv.Atoi(levelStr)
		if err != nil || levelFilter < 1 {
			return 0, 0, 0, errs.BadRequest("invalid level parameter")
		}
	}
	return page, limit, levelFilter, nil
}

func normalizeSpellDurationField(spell *models.Spell) {
	if spell.Duration == nil {
		return
	}
	t := strings.TrimSpace(*spell.Duration)
	if t == "" {
		spell.Duration = nil
		return
	}
	*spell.Duration = t
}

// validateSpellMechanicsOrRespond normalizes duration and checks type, range, and duration. Returns false if the response was already written.
func validateSpellMechanicsOrRespond(w http.ResponseWriter, spell *models.Spell) bool {
	normalizeSpellDurationField(spell)
	if spell.Type == "" {
		spell.Type = models.SpellTypeUtility
	}
	if !spell.Type.IsValid() {
		respondError(w, http.StatusBadRequest, "invalid spell type (use Attack, Save, Effect, Healing, or Utility)")
		return false
	}
	if spell.Range != nil && *spell.Range < 0 {
		respondError(w, http.StatusBadRequest, "range must be a non-negative integer (feet)")
		return false
	}
	if spell.Duration != nil && *spell.Duration != "" {
		if err := models.ValidateSpellDuration(*spell.Duration); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return false
		}
	}
	if spell.SaveAttr != nil && !spell.SaveAttr.IsValid() {
		respondError(w, http.StatusBadRequest, "invalid save_attr (use STR, DEX, CON, INT, WIS, or CHA)")
		return false
	}
	if spell.DamageType != nil && !spell.DamageType.IsValid() {
		respondError(w, http.StatusBadRequest, "invalid damage_type")
		return false
	}
	if err := models.ValidateSpellDamageDicePair(spell.DamageDiceCount, spell.DamageDieSize); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

const gmEmail = "rpupo63@gmail.com"

type spellHandler struct {
	spellRepo          database.SpellRepository
	characterRepo      database.CharacterRepository
	classRepo          database.ClassRepository
	raceRepo           database.RaceRepository
	userRepo           database.UserRepository
	synthesisService   *services.SpellSynthesisService
	interpreterService *services.ComponentInterpreterService
	spellAIService     *services.SpellAIService
}

func newSpellHandler(spellRepo database.SpellRepository, characterRepo database.CharacterRepository, classRepo database.ClassRepository, raceRepo database.RaceRepository, userRepo database.UserRepository, synthesisService *services.SpellSynthesisService, interpreterService *services.ComponentInterpreterService, spellAIService *services.SpellAIService) *spellHandler {
	return &spellHandler{
		spellRepo:          spellRepo,
		characterRepo:      characterRepo,
		classRepo:          classRepo,
		raceRepo:           raceRepo,
		userRepo:           userRepo,
		synthesisService:   synthesisService,
		interpreterService: interpreterService,
		spellAIService:     spellAIService,
	}
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

const (
	spellbookScopeMine          = "mine"
	spellbookScopeCastable      = "castable"
	spellbookScopeMineOrCastable = "mine_or_castable"
	spellbookScopeAll           = "all"
)

// isGM checks if the authenticated user is the game master
func (h *spellHandler) isGM(r *http.Request) (bool, error) {
	isAdmin, err := ctxGetIsAdmin(r.Context())
	if err == nil && isAdmin {
		return true, nil
	}
	// Fallback to email check if context fails or is false
	authUserID, _ := ctxGetUserID(r.Context())
	if authUserID == "" {
		return false, nil
	}
	id, err := uuid.Parse(authUserID)
	if err != nil {
		return false, err
	}
	user, err := h.userRepo.FindByID(id)
	if err != nil {
		return false, err
	}
	return user.Email == gmEmail, nil
}

// getUncheckedSpells returns all spells not yet reviewed by the GM (GM only)
func (h *spellHandler) getUncheckedSpells() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gm, err := h.isGM(r)
		if err != nil || !gm {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		spells, err := h.spellRepo.FindUnchecked()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get unchecked spells")
			respondError(w, http.StatusInternalServerError, "Failed to get unchecked spells")
			return
		}

		respondJSON(w, http.StatusOK, spells)
	}
}

// getCharacterSpellbook returns a paginated spell list for a character by scope.
func (h *spellHandler) getCharacterSpellbook() http.HandlerFunc {
	type paginatedSpellbookResponse struct {
		Spells     []SpellResponse `json:"spells"`
		TotalCount int64           `json:"total_count"`
		Page       int             `json:"page"`
		Limit      int             `json:"limit"`
	}

	respondSpellsPage := func(w http.ResponseWriter, spells []*models.Spell, totalCount int64, page, limit int) {
		out := make([]SpellResponse, len(spells))
		for i, spell := range spells {
			out[i] = SpellResponse{Spell: spell}
			if spell.Character != nil {
				out[i].CharacterName = &spell.Character.Name
				out[i].CharacterClass = &spell.Character.Class.Name
			}
		}
		respondJSON(w, http.StatusOK, paginatedSpellbookResponse{
			Spells:     out,
			TotalCount: totalCount,
			Page:       page,
			Limit:      limit,
		})
	}

	return func(w http.ResponseWriter, r *http.Request) {
		characterIDStr := chi.URLParam(r, "characterID")
		characterID, err := uuid.Parse(characterIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByIDWithRelations(characterID)
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

		scope := r.URL.Query().Get("scope")
		if scope == "" {
			scope = spellbookScopeMineOrCastable
		}
		switch scope {
		case spellbookScopeMine, spellbookScopeCastable, spellbookScopeMineOrCastable, spellbookScopeAll:
		default:
			respondError(w, http.StatusBadRequest, "invalid scope parameter")
			return
		}

		page, limit, levelFilter, err := parseSpellbookPaginationParams(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		switch scope {
		case spellbookScopeMine:
			spells, totalCount, err := h.spellRepo.FindByUserIDPaginated(character.UserID, page, limit, levelFilter)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get character spellbook (mine)")
				respondError(w, http.StatusInternalServerError, "Failed to get spells")
				return
			}
			respondSpellsPage(w, spells, totalCount, page, limit)

		case spellbookScopeAll:
			spells, totalCount, err := h.spellRepo.FindAll(page, limit, levelFilter)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get character spellbook (all)")
				respondError(w, http.StatusInternalServerError, "Failed to get spells")
				return
			}
			respondSpellsPage(w, spells, totalCount, page, limit)

		case spellbookScopeCastable, spellbookScopeMineOrCastable:
			allSpells, err := h.spellRepo.FindAllWithComponents()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get spells")
				respondError(w, http.StatusInternalServerError, "Failed to get spells")
				return
			}

			class := &character.Class
			race := &character.Race
			unlimitedComponentIDs := SpellPoolAllowlist(class, race)

			characterComponentCounts := make(map[uuid.UUID]int)
			for _, charComp := range character.Components {
				characterComponentCounts[charComp.ComponentID] = charComp.Count
			}

			var filtered []*models.Spell
			if scope == spellbookScopeCastable {
				for _, sp := range allSpells {
					if !spellCastableForCharacter(sp, unlimitedComponentIDs, characterComponentCounts) {
						continue
					}
					if levelFilter > 0 && sp.Level != levelFilter {
						continue
					}
					filtered = append(filtered, sp)
				}
			} else {
				seen := make(map[uuid.UUID]struct{})
				for _, sp := range allSpells {
					mine := sp.UserID == character.UserID
					cast := spellCastableForCharacter(sp, unlimitedComponentIDs, characterComponentCounts)
					if !mine && !cast {
						continue
					}
					if levelFilter > 0 && sp.Level != levelFilter {
						continue
					}
					if _, ok := seen[sp.ID]; ok {
						continue
					}
					seen[sp.ID] = struct{}{}
					filtered = append(filtered, sp)
				}
			}

			sort.Slice(filtered, func(i, j int) bool {
				if filtered[i].Level != filtered[j].Level {
					return filtered[i].Level < filtered[j].Level
				}
				return filtered[i].Name < filtered[j].Name
			})

			totalCount := int64(len(filtered))
			offset := (page - 1) * limit
			var pageSpells []*models.Spell
			if offset < len(filtered) {
				end := offset + limit
				if end > len(filtered) {
					end = len(filtered)
				}
				pageSpells = filtered[offset:end]
			}
			respondSpellsPage(w, pageSpells, totalCount, page, limit)
		}
	}
}

// getAllSpells returns all spells
func (h *spellHandler) getAllSpells() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, limit, levelFilter, err := parsePaginationParams(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		spells, totalCount, err := h.spellRepo.FindAll(page, limit, levelFilter)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get spells")
			respondError(w, http.StatusInternalServerError, "Failed to get spells")
			return
		}

		type PaginatedSpellsResponse struct {
			Spells     []SpellResponse `json:"spells"`
			TotalCount int64           `json:"total_count"`
			Page       int             `json:"page"`
			Limit      int             `json:"limit"`
		}

		responseSpells := make([]SpellResponse, len(spells))
		for i, spell := range spells {
			responseSpells[i] = SpellResponse{
				Spell: spell,
			}
			if spell.Character != nil {
				responseSpells[i].CharacterName = &spell.Character.Name
				responseSpells[i].CharacterClass = &spell.Character.Class.Name
			}
		}

		respondJSON(w, http.StatusOK, PaginatedSpellsResponse{
			Spells:     responseSpells,
			TotalCount: totalCount,
			Page:       page,
			Limit:      limit,
		})
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

		page, limit, levelFilter, err := parsePaginationParams(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		spells, totalCount, err := h.spellRepo.FindByUserIDPaginated(userID, page, limit, levelFilter)
		if err != nil {
			log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get spells")
			respondError(w, http.StatusInternalServerError, "Failed to get spells")
			return
		}

		type PaginatedSpellsResponse struct {
			Spells     []SpellResponse `json:"spells"`
			TotalCount int64           `json:"total_count"`
			Page       int             `json:"page"`
			Limit      int             `json:"limit"`
		}

		responseSpells := make([]SpellResponse, len(spells))
		for i, spell := range spells {
			responseSpells[i] = SpellResponse{
				Spell: spell,
			}
		}

		respondJSON(w, http.StatusOK, PaginatedSpellsResponse{
			Spells:     responseSpells,
			TotalCount: totalCount,
			Page:       page,
			Limit:      limit,
		})
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

		page, limit, levelFilter, err := parsePaginationParams(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		spells, totalCount, err := h.spellRepo.FindByCharacterID(characterID, page, limit, levelFilter)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterIDStr).Msg("Failed to get character spells")
			respondError(w, http.StatusInternalServerError, "Failed to get character spells")
			return
		}

		type PaginatedSpellsResponse struct {
			Spells     []SpellResponse `json:"spells"`
			TotalCount int64           `json:"total_count"`
			Page       int             `json:"page"`
			Limit      int             `json:"limit"`
		}

		responseSpells := make([]SpellResponse, len(spells))
		for i, spell := range spells {
			responseSpells[i] = SpellResponse{
				Spell: spell,
			}
			if spell.Character != nil {
				responseSpells[i].CharacterName = &spell.Character.Name
				responseSpells[i].CharacterClass = &spell.Character.Class.Name
			}
		}

		respondJSON(w, http.StatusOK, PaginatedSpellsResponse{
			Spells:     responseSpells,
			TotalCount: totalCount,
			Page:       page,
			Limit:      limit,
		})
	}
}

// retryAIField regenerates the AI recommendation for one specific spell field.
// If apply_recommendation is true, the recommendation is also written to the original field.
func (h *spellHandler) retryAIField() http.HandlerFunc {
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
		gm, _ := h.isGM(r)
		if !gm && authUserID != spell.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		var req RetryAIFieldRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		validFields := map[string]bool{
			"name": true, "description": true, "type": true, "range": true,
			"duration": true, "damage_dice": true, "damage_type": true, "save_attr": true,
		}
		if !validFields[req.Field] {
			respondError(w, http.StatusBadRequest, "field must be one of: name, description, type, range, duration, damage_dice, damage_type, save_attr")
			return
		}

		opinion, _, err := h.spellAIService.GetSpellOpinion(r.Context(), spell, spell.Components)
		if err != nil {
			log.Error().Err(err).Str("spellID", idStr).Msg("Failed to get AI opinion for retry")
			respondError(w, http.StatusInternalServerError, "Failed to get AI opinion")
			return
		}

		updates := map[string]interface{}{}

		switch req.Field {
		case "name":
			updates["ai_recommended_name"] = opinion.RecommendedName
			if req.ApplyRecommendation && opinion.RecommendedName != nil {
				updates["name"] = *opinion.RecommendedName
			}
		case "description":
			updates["ai_recommended_description"] = opinion.RecommendedDescription
			if req.ApplyRecommendation && opinion.RecommendedDescription != nil {
				updates["description"] = *opinion.RecommendedDescription
			}
		case "type":
			aiT, _, _ := models.NormalizeSpellAIRecommendations(opinion.RecommendedType, nil, nil)
			if aiT != nil {
				updates["ai_recommended_type"] = string(*aiT)
			} else {
				updates["ai_recommended_type"] = nil
			}
			if req.ApplyRecommendation && aiT != nil {
				updates["type"] = string(*aiT)
			}
		case "range":
			_, aiR, _ := models.NormalizeSpellAIRecommendations(nil, opinion.RecommendedRange, nil)
			if aiR != nil {
				updates["ai_recommended_range"] = *aiR
			} else {
				updates["ai_recommended_range"] = nil
			}
			if req.ApplyRecommendation && aiR != nil {
				updates["range"] = *aiR
			}
		case "duration":
			_, _, aiD := models.NormalizeSpellAIRecommendations(nil, nil, opinion.RecommendedDuration)
			if aiD != nil {
				updates["ai_recommended_duration"] = *aiD
			} else {
				updates["ai_recommended_duration"] = nil
			}
			if req.ApplyRecommendation && aiD != nil {
				updates["duration"] = *aiD
			}
		case "damage_dice":
			_, _, dc, ds := models.NormalizeSpellAIRecommendationsExtras(nil, nil, opinion.RecommendedDamageDiceCount, opinion.RecommendedDamageDieSize)
			updates["ai_recommended_damage_dice_count"] = dc
			updates["ai_recommended_damage_die_size"] = ds
			if req.ApplyRecommendation && dc != nil && ds != nil {
				updates["damage_dice_count"] = *dc
				updates["damage_die_size"] = *ds
			}
		case "damage_type":
			_, dt, _, _ := models.NormalizeSpellAIRecommendationsExtras(nil, opinion.RecommendedDamageType, nil, nil)
			updates["ai_recommended_damage_type"] = dt
			if req.ApplyRecommendation && dt != nil {
				updates["damage_type"] = *dt
			}
		case "save_attr":
			sa, _, _, _ := models.NormalizeSpellAIRecommendationsExtras(opinion.RecommendedSaveAttr, nil, nil, nil)
			updates["ai_recommended_save_attr"] = sa
			if req.ApplyRecommendation && sa != nil {
				updates["save_attr"] = *sa
			}
		}

		if err := h.spellRepo.UpdateFields(spell.ID, updates); err != nil {
			log.Error().Err(err).Str("spellID", idStr).Msg("Failed to save AI retry")
			respondError(w, http.StatusInternalServerError, "Failed to save AI retry")
			return
		}

		updated, err := h.spellRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch updated spell")
			return
		}
		respondJSON(w, http.StatusOK, updated)
	}
}

// getSpellOpinion returns an AI-generated opinion on a spell
func (h *spellHandler) getSpellOpinion() http.HandlerFunc {
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

		opinion, raw, err := h.spellAIService.GetSpellOpinion(r.Context(), spell, spell.Components)
		if err != nil {
			log.Error().Err(err).Str("spellID", idStr).Msg("Failed to get AI opinion")
			respondError(w, http.StatusInternalServerError, "Failed to get AI opinion")
			return
		}

		// Update spell with opinions and raw output
		spell.AIDescriptionOpinion = &opinion.DescriptionOpinion
		spell.AIDamageOpinion = &opinion.DamageOpinion
		spell.AIEffectOpinion = &opinion.EffectOpinion
		spell.AIOverallVerdict = &opinion.OverallVerdict
		spell.AIRawOutput = &raw

		// Save recommended edits (mechanics normalized to enum / feet / validated duration)
		spell.AIRecommendedName = opinion.RecommendedName
		spell.AIRecommendedDescription = opinion.RecommendedDescription
		models.ApplyNormalizedAIRecommendations(spell, opinion.RecommendedType, opinion.RecommendedRange, opinion.RecommendedDuration)
		models.ApplyNormalizedAIRecommendationsExtras(spell, opinion.RecommendedSaveAttr, opinion.RecommendedDamageType, opinion.RecommendedDamageDiceCount, opinion.RecommendedDamageDieSize)

		if err := h.spellRepo.Update(spell); err != nil {
			log.Error().Err(err).Str("spellID", idStr).Msg("Failed to update spell with opinions")
		}

		respondJSON(w, http.StatusOK, opinion)
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
			UserID:          req.UserID,
			CharacterID:     req.CharacterID,
			Name:            req.Name,
			Description:     req.Description,
			Type:            req.Type,
			Range:           req.Range,
			Duration:        req.Duration,
			Concentration:   req.Concentration,
			SaveAttr:        req.SaveAttr,
			DamageDiceCount: req.DamageDiceCount,
			DamageDieSize:   req.DamageDieSize,
			DamageType:      req.DamageType,
			AddModifier:     req.AddModifier,
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

			spell.Level = synthesis.Level
			spell.SuggestedDamageDiceCount = synthesis.SuggestedDamageDiceCount
			spell.SuggestedDamageDieSize = synthesis.SuggestedDamageDieSize

			// Hybrid fill: use request value if provided, else fall back to suggestion
			if spell.Type == "" {
				spell.Type = synthesis.SuggestedType
			}
			if spell.Range == nil && synthesis.SuggestedRange != nil {
				r := *synthesis.SuggestedRange
				spell.Range = &r
			}
			if spell.DamageType == nil && synthesis.SuggestedDamageType != nil {
				dt := models.DamageType(*synthesis.SuggestedDamageType)
				spell.DamageType = &dt
			}
			if spell.DamageDiceCount == nil && synthesis.SuggestedDamageDiceCount != nil {
				c := *synthesis.SuggestedDamageDiceCount
				spell.DamageDiceCount = &c
			}
			if spell.DamageDieSize == nil && synthesis.SuggestedDamageDieSize != nil {
				f := *synthesis.SuggestedDamageDieSize
				spell.DamageDieSize = &f
			}
			if spell.Duration == nil && synthesis.SuggestedDuration != nil {
				spell.Duration = synthesis.SuggestedDuration
			}
			if !spell.Concentration && synthesis.SuggestedConcentration {
				spell.Concentration = true
			}
		}

		if spell.Level == 0 {
			spell.Level = 1
		}
		if spell.Type == "" {
			spell.Type = models.SpellTypeUtility
		}

		normalizeSpellDurationField(spell)
		if !validateSpellMechanicsOrRespond(w, spell) {
			return
		}

		if err := h.spellRepo.Add(spell, req.ComponentIDs); err != nil {
			log.Error().Err(err).Msg("Failed to create spell")
			respondError(w, http.StatusInternalServerError, "Failed to create spell")
			return
		}

		// Trigger background AI review
		go h.triggerBackgroundAIReview(spell.ID)

		spell, _ = h.spellRepo.FindByID(spell.ID)
		respondJSON(w, http.StatusCreated, spell)
	}
}

// triggerBackgroundAIReview fetches a spell and its components, gets an AI opinion, and saves it.
func (h *spellHandler) triggerBackgroundAIReview(spellID uuid.UUID) {
	// 1. Fetch the spell with components
	spell, err := h.spellRepo.FindByID(spellID)
	if err != nil {
		log.Error().Err(err).Str("spellID", spellID.String()).Msg("Background AI review: failed to fetch spell")
		return
	}

	// 2. Get AI Opinion
	opinion, raw, err := h.spellAIService.GetSpellOpinion(context.Background(), spell, spell.Components)
	if err != nil {
		log.Error().Err(err).Str("spellID", spellID.String()).Msg("Background AI review: failed to get AI opinion")
		return
	}

	// 3. Update spell with opinions
	spell.AIDescriptionOpinion = &opinion.DescriptionOpinion
	spell.AIDamageOpinion = &opinion.DamageOpinion
	spell.AIEffectOpinion = &opinion.EffectOpinion
	spell.AIOverallVerdict = &opinion.OverallVerdict
	spell.AIRawOutput = &raw

	// Save recommended edits (mechanics normalized to enum / feet / validated duration)
	spell.AIRecommendedName = opinion.RecommendedName
	spell.AIRecommendedDescription = opinion.RecommendedDescription
	models.ApplyNormalizedAIRecommendations(spell, opinion.RecommendedType, opinion.RecommendedRange, opinion.RecommendedDuration)
	models.ApplyNormalizedAIRecommendationsExtras(spell, opinion.RecommendedSaveAttr, opinion.RecommendedDamageType, opinion.RecommendedDamageDiceCount, opinion.RecommendedDamageDieSize)

	if err := h.spellRepo.Update(spell); err != nil {
		log.Error().Err(err).Str("spellID", spellID.String()).Msg("Background AI review: failed to update spell with opinions")
	} else {
		log.Info().Str("spellID", spellID.String()).Msg("Background AI review: successfully updated spell with opinions")
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

		gm, _ := h.isGM(r)
		if !gm && authUserID != spell.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		var req UpdateSpellRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// If components change, re-run synthesis
		componentsChanged := false
		if len(req.ComponentIDs) > 0 {
			componentsChanged = true
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

			// Always recompute level and suggested damage from synthesis
			spell.Level = synthesis.Level
			spell.SuggestedDamageDiceCount = synthesis.SuggestedDamageDiceCount
			spell.SuggestedDamageDieSize = synthesis.SuggestedDamageDieSize

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
		if req.Type != nil {
			t := *req.Type
			if t == "" {
				spell.Type = models.SpellTypeUtility
			} else if !t.IsValid() {
				respondError(w, http.StatusBadRequest, "invalid spell type (use Attack, Save, Effect, Healing, or Utility)")
				return
			} else {
				spell.Type = t
			}
		}
		if req.Range != nil {
			if *req.Range < 0 {
				respondError(w, http.StatusBadRequest, "range must be a non-negative integer (feet)")
				return
			}
			r := *req.Range
			spell.Range = &r
		}
		if req.Duration != nil {
			spell.Duration = req.Duration
			normalizeSpellDurationField(spell)
		}
		if req.Concentration != nil {
			spell.Concentration = *req.Concentration
		}
		if req.SaveAttr != nil {
			spell.SaveAttr = req.SaveAttr
		}
		if req.DamageDiceCount != nil {
			spell.DamageDiceCount = req.DamageDiceCount
		}
		if req.DamageDieSize != nil {
			spell.DamageDieSize = req.DamageDieSize
		}
		if req.DamageType != nil {
			spell.DamageType = req.DamageType
		}
		if req.AddModifier != nil {
			spell.AddModifier = *req.AddModifier
		}
		if req.Checked != nil && gm {
			spell.Checked = *req.Checked
		}

		if !validateSpellMechanicsOrRespond(w, spell) {
			return
		}

		if err := h.spellRepo.Update(spell); err != nil {
			log.Error().Err(err).Msg("Failed to update spell")
			respondError(w, http.StatusInternalServerError, "Failed to update spell")
			return
		}

		// Trigger background AI review if components or description changed
		if componentsChanged || req.Description != nil {
			go h.triggerBackgroundAIReview(spell.ID)
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
