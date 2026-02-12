package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/services"
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

		serviceReq := services.CreateCorpseRequest{
			MapID:               req.MapID,
			Name:                req.Name,
			CreatureType:        req.CreatureType,
			CreatureSize:        req.CreatureSize,
			ChallengeRating:     req.ChallengeRating,
			GridX:               req.GridX,
			GridY:               req.GridY,
			AvailableComponents: req.AvailableComponents,
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

		// Add harvested components to character's inventory
		componentNames := make([]string, 0)
		for _, compIDStr := range result.ComponentIDs {
			compID, err := uuid.Parse(compIDStr)
			if err != nil {
				continue
			}
			for i := 0; i < result.ComponentsYielded; i++ {
				if err := h.characterRepo.UpdateComponentCount(req.CharacterID, compID, 1); err != nil {
					log.Error().Err(err).Str("componentID", compIDStr).Msg("Failed to add component to inventory")
				}
			}
			// Look up component name for response
			comp, err := h.componentRepo.GetComponentByID(compID)
			if err == nil {
				componentNames = append(componentNames, comp.Name)
			}
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"components_yielded": result.ComponentsYielded,
			"component_ids":     result.ComponentIDs,
			"component_names":   componentNames,
			"message":           result.Message,
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
