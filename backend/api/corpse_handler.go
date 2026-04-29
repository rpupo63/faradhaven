package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/services"
	"github.com/rs/zerolog/log"
)

type corpseHandler struct {
	corpseService *services.CorpseService
	characterRepo *database.CharacterRepo
	componentRepo *database.ComponentRepo
}

func newCorpseHandler(
	corpseService *services.CorpseService,
	characterRepo *database.CharacterRepo,
	componentRepo *database.ComponentRepo,
) *corpseHandler {
	return &corpseHandler{
		corpseService: corpseService,
		characterRepo: characterRepo,
		componentRepo: componentRepo,
	}
}

// GetCorpses returns all corpses, optionally filtered by map_id
func (h *corpseHandler) GetCorpses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mapID *uuid.UUID
		mapIDStr := r.URL.Query().Get("map_id")
		if mapIDStr != "" {
			id, err := uuid.Parse(mapIDStr)
			if err != nil {
				respondError(w, http.StatusBadRequest, "Invalid map_id")
				return
			}
			mapID = &id
		}

		corpses, err := h.corpseService.GetCorpses(mapID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get corpses")
			respondError(w, http.StatusInternalServerError, "Failed to get corpses")
			return
		}

		respondJSON(w, http.StatusOK, corpses)
	}
}

// GetCorpse returns a single corpse by ID
func (h *corpseHandler) GetCorpse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		corpseID, err := uuid.Parse(chi.URLParam(r, "corpseID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid corpse ID")
			return
		}

		corpse, err := h.corpseService.GetCorpse(corpseID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Corpse not found")
			return
		}

		respondJSON(w, http.StatusOK, corpse)
	}
}

// CreateCorpse creates a new corpse (DM creates when creature dies)
func (h *corpseHandler) CreateCorpse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateCorpseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ct, ok := models.ParseCreatureType(strings.TrimSpace(req.CreatureType))
		if !ok || strings.TrimSpace(req.CreatureType) == "" {
			respondError(w, http.StatusBadRequest, "invalid creature_type")
			return
		}
		sz := models.SizeMedium
		if strings.TrimSpace(req.CreatureSize) != "" {
			var okSz bool
			sz, okSz = models.ParseCreatureSize(strings.TrimSpace(req.CreatureSize))
			if !okSz {
				respondError(w, http.StatusBadRequest, "invalid creature_size")
				return
			}
		}

		avail := req.AvailableComponents
		if len(avail) > 0 {
			parsed := make([]uuid.UUID, 0, len(avail))
			for _, s := range avail {
				id, err := uuid.Parse(s)
				if err != nil {
					continue
				}
				parsed = append(parsed, id)
			}
			comps, err := h.componentRepo.GetComponentsByIDs(parsed)
			if err != nil || len(comps) == 0 {
				avail = nil
			} else {
				avail = make([]string, len(comps))
				for i, c := range comps {
					avail[i] = c.ID.String()
				}
			}
		}
		if len(avail) == 0 {
			all, err := h.componentRepo.GetAllComponents()
			if err != nil {
				log.Error().Err(err).Msg("Failed to load components for corpse loot")
				respondError(w, http.StatusInternalServerError, "Failed to create corpse")
				return
			}
			avail = RandomDistinctComponentIDStrings(all, RollCorpseComponentSlotCount())
		}

		serviceReq := services.CreateCorpseRequest{
			MapID:               req.MapID,
			Name:                req.Name,
			CreatureType:        ct,
			CreatureSize:        sz,
			ChallengeRating:     req.ChallengeRating,
			GridX:               req.GridX,
			GridY:               req.GridY,
			AvailableComponents: avail,
			ComponentYield:      req.ComponentYield,
			SourceBeastID:       req.SourceBeastID,
			ExpiresInMinutes:    req.ExpiresInMinutes,
		}

		corpse, err := h.corpseService.CreateCorpse(serviceReq)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create corpse")
			respondError(w, http.StatusInternalServerError, "Failed to create corpse")
			return
		}

		respondJSON(w, http.StatusCreated, corpse)
	}
}

// HarvestCorpse harvests components from a corpse (Ironwright Scavenge)
func (h *corpseHandler) HarvestCorpse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		corpseID, err := uuid.Parse(chi.URLParam(r, "corpseID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid corpse ID")
			return
		}

		var req HarvestCorpseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		result, err := h.corpseService.HarvestCorpse(corpseID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to harvest corpse")
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse component IDs; keep only rows that exist (stale corpse UUIDs after reseed)
		compIDs := make([]uuid.UUID, 0, len(result.ComponentIDs))
		for _, compIDStr := range result.ComponentIDs {
			compID, err := uuid.Parse(compIDStr)
			if err != nil {
				continue
			}
			compIDs = append(compIDs, compID)
		}
		fetchedComps, err := h.componentRepo.GetComponentsByIDs(compIDs)
		if err != nil {
			log.Error().Err(err).Msg("Failed to resolve corpse components")
			respondError(w, http.StatusInternalServerError, "Failed to harvest corpse")
			return
		}
		validCompIDs := make([]uuid.UUID, 0, len(fetchedComps))
		for _, c := range fetchedComps {
			validCompIDs = append(validCompIDs, c.ID)
		}
		if len(validCompIDs) == 0 {
			all, err := h.componentRepo.GetAllComponents()
			if err != nil {
				log.Error().Err(err).Msg("Failed to load components for harvest fallback")
				respondError(w, http.StatusInternalServerError, "Failed to harvest corpse")
				return
			}
			validCompIDs = RandomDistinctComponentIDs(all, RollCorpseComponentSlotCount())
			fetchedComps, err = h.componentRepo.GetComponentsByIDs(validCompIDs)
			if err != nil {
				log.Error().Err(err).Msg("Failed to resolve fallback harvest components")
				respondError(w, http.StatusInternalServerError, "Failed to harvest corpse")
				return
			}
		}
		compNameByID := make(map[uuid.UUID]string, len(fetchedComps))
		for _, c := range fetchedComps {
			compNameByID[c.ID] = c.Name
		}

		// Update inventory counts and collect names
		componentNames := make([]string, 0, len(validCompIDs))
		for _, compID := range validCompIDs {
			for i := 0; i < result.ComponentsYielded; i++ {
				if err := h.characterRepo.UpdateComponentCount(req.CharacterID, compID, 1); err != nil {
					log.Error().Err(err).Str("componentID", compID.String()).Msg("Failed to add component to inventory")
				}
			}
			if name, ok := compNameByID[compID]; ok {
				componentNames = append(componentNames, name)
			}
		}

		grantedIDs := make([]string, len(validCompIDs))
		for i, id := range validCompIDs {
			grantedIDs[i] = id.String()
		}
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"components_yielded": result.ComponentsYielded,
			"component_ids":      grantedIDs,
			"component_names":    componentNames,
			"message":            result.Message,
		})
	}
}

// ConsumeCorpse consumes a corpse for Lorewright Visceral Psychometry
func (h *corpseHandler) ConsumeCorpse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		corpseID, err := uuid.Parse(chi.URLParam(r, "corpseID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid corpse ID")
			return
		}

		var req ConsumeCorpseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Get character level for trauma DC calculation
		character, err := h.characterRepo.FindByID(req.CharacterID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Character not found")
			return
		}

		result, err := h.corpseService.ConsumeCorpse(corpseID, character.Level)
		if err != nil {
			log.Error().Err(err).Msg("Failed to consume corpse")
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, result)
	}
}

// DeleteCorpse removes a corpse
func (h *corpseHandler) DeleteCorpse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		corpseID, err := uuid.Parse(chi.URLParam(r, "corpseID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid corpse ID")
			return
		}

		if err := h.corpseService.DeleteCorpse(corpseID); err != nil {
			log.Error().Err(err).Msg("Failed to delete corpse")
			respondError(w, http.StatusInternalServerError, "Failed to delete corpse")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// ScavengeComponents handles Lorewright component scavenging from a corpse
func (h *corpseHandler) ScavengeComponents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		corpseID, err := uuid.Parse(chi.URLParam(r, "corpseID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid corpse ID")
			return
		}

		var req ScavengeComponentsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		character, err := h.characterRepo.FindByID(req.CharacterID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Character not found")
			return
		}

		if character.Class.Name != "The Lorewright" {
			respondError(w, http.StatusBadRequest, "Only Lorewrights can scavenge components")
			return
		}

		// Check component capacity: total count vs Wisdom modifier
		wisMod := int(math.Floor(float64(character.Wisdom-10) / 2))
		if wisMod < 1 {
			wisMod = 1
		}
		currentTotal := 0
		for _, cc := range character.Components {
			currentTotal += cc.Count
		}
		if currentTotal >= wisMod {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Component capacity full (%d/%d, Wisdom modifier)", currentTotal, wisMod))
			return
		}

		corpse, err := h.corpseService.GetCorpse(corpseID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Corpse not found")
			return
		}

		if !corpse.CanBeScavenged() {
			if corpse.HasBeenScavenged {
				respondError(w, http.StatusBadRequest, "Corpse has already been scavenged")
				return
			}
			respondError(w, http.StatusBadRequest, "Corpse is too old to scavenge (must be within 1 hour of death)")
			return
		}

		// Calculate yield
		yield := corpse.ComponentYield
		// Level 7+ (Essence Decanting): +1 component per harvest, min 2
		if character.Level >= 7 {
			yield++
			if yield < 2 {
				yield = 2
			}
		}

		// Cap yield by remaining capacity
		remaining := wisMod - currentTotal
		if yield > remaining {
			yield = remaining
		}

		// Resolve corpse component IDs against the live catalog; fall back to random if none match
		parsedCorpse := make([]uuid.UUID, 0, len(corpse.AvailableComponents))
		for _, compIDStr := range corpse.AvailableComponents {
			compID, err := uuid.Parse(compIDStr)
			if err != nil {
				continue
			}
			parsedCorpse = append(parsedCorpse, compID)
		}
		scavComps, err := h.componentRepo.GetComponentsByIDs(parsedCorpse)
		if err != nil {
			log.Error().Err(err).Msg("Failed to resolve scavenged components")
			respondError(w, http.StatusInternalServerError, "Failed to scavenge corpse")
			return
		}
		scavCompIDs := make([]uuid.UUID, 0, len(scavComps))
		for _, c := range scavComps {
			scavCompIDs = append(scavCompIDs, c.ID)
		}
		if len(scavCompIDs) == 0 {
			all, err := h.componentRepo.GetAllComponents()
			if err != nil {
				log.Error().Err(err).Msg("Failed to load components for scavenge fallback")
				respondError(w, http.StatusInternalServerError, "Failed to scavenge corpse")
				return
			}
			scavCompIDs = RandomDistinctComponentIDs(all, RollCorpseComponentSlotCount())
			scavComps, err = h.componentRepo.GetComponentsByIDs(scavCompIDs)
			if err != nil {
				log.Error().Err(err).Msg("Failed to resolve fallback scavenge components")
				respondError(w, http.StatusInternalServerError, "Failed to scavenge corpse")
				return
			}
		}
		scavNameByID := make(map[uuid.UUID]string, len(scavComps))
		for _, c := range scavComps {
			scavNameByID[c.ID] = c.Name
		}

		// Add components to character inventory
		componentNames := make([]string, 0, len(scavCompIDs))
		for _, compID := range scavCompIDs {
			for i := 0; i < yield; i++ {
				if err := h.characterRepo.UpdateComponentCount(req.CharacterID, compID, 1); err != nil {
					log.Error().Err(err).Str("componentID", compID.String()).Msg("Failed to add scavenged component")
				}
			}
			if name, ok := scavNameByID[compID]; ok {
				componentNames = append(componentNames, name)
			}
		}

		// Mark corpse as scavenged
		corpse.HasBeenScavenged = true
		h.corpseService.UpdateCorpse(corpse)

		nSlots := len(scavCompIDs)
		totalAdded := yield * nSlots
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"components_yielded": yield,
			"component_names":    componentNames,
			"capacity_used":      currentTotal + totalAdded,
			"capacity_max":       wisMod,
			"message":            fmt.Sprintf("Scavenged %d components from %s", totalAdded, corpse.Name),
		})
	}
}
