package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/services"
	"github.com/rs/zerolog/log"
)

type resourceHandler struct {
	resourceRepo    *database.CharacterResourceRepo
	characterRepo   database.CharacterRepository
	resourceService *services.ResourceService
}

func newResourceHandler(
	resourceRepo *database.CharacterResourceRepo,
	characterRepo database.CharacterRepository,
	resourceService *services.ResourceService,
) *resourceHandler {
	return &resourceHandler{
		resourceRepo:    resourceRepo,
		characterRepo:   characterRepo,
		resourceService: resourceService,
	}
}

// GetResources returns all resources for a character
func (h *resourceHandler) GetResources() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		resources, err := h.resourceRepo.FindByCharacterID(characterID)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterID.String()).Msg("Failed to get resources")
			respondError(w, http.StatusInternalServerError, "Failed to get resources")
			return
		}

		respondJSON(w, http.StatusOK, resources)
	}
}

// GetResource returns a specific resource by key
func (h *resourceHandler) GetResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		key := chi.URLParam(r, "key")
		if key == "" {
			respondError(w, http.StatusBadRequest, "Resource key required")
			return
		}

		resource, err := h.resourceRepo.FindByCharacterAndKey(characterID, key)
		if err != nil {
			respondError(w, http.StatusNotFound, "Resource not found")
			return
		}

		respondJSON(w, http.StatusOK, resource)
	}
}

// CreateResource creates a new resource for a character
func (h *resourceHandler) CreateResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		var req CreateResourceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Check if resource already exists
		existing, _ := h.resourceRepo.FindByCharacterAndKey(characterID, req.ResourceKey)
		if existing != nil {
			respondError(w, http.StatusConflict, "Resource with this key already exists")
			return
		}

		resource := &models.CharacterResource{
			CharacterID:        characterID,
			ResourceKey:        req.ResourceKey,
			ResourceName:       req.ResourceName,
			CurrentValue:       req.CurrentValue,
			MaxValue:           req.MaxValue,
			RestoreOnShortRest: req.RestoreOnShortRest,
			RestoreOnLongRest:  req.RestoreOnLongRest,
			RestoreAmount:      req.RestoreAmount,
			DecaysOnLongRest:   req.DecaysOnLongRest,
		}

		if err := h.resourceRepo.Add(resource); err != nil {
			log.Error().Err(err).Msg("Failed to create resource")
			respondError(w, http.StatusInternalServerError, "Failed to create resource")
			return
		}

		respondJSON(w, http.StatusCreated, resource)
	}
}

// SpendResource decrements a resource value
func (h *resourceHandler) SpendResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		key := chi.URLParam(r, "key")
		if key == "" {
			respondError(w, http.StatusBadRequest, "Resource key required")
			return
		}

		var req ResourceDeltaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if character, err := h.characterRepo.FindByID(characterID); err == nil && h.resourceService != nil {
			if ensureErr := h.resourceService.EnsureTrackableClassResources(character); ensureErr != nil {
				log.Warn().Err(ensureErr).Str("characterID", characterID.String()).Msg("EnsureTrackableClassResources before spend")
			}
		}

		resource, err := h.resourceRepo.FindByCharacterAndKey(characterID, key)
		if err != nil {
			respondError(w, http.StatusNotFound, "Resource not found")
			return
		}

		// Check if we have enough
		if resource.CurrentValue < req.Amount {
			respondError(w, http.StatusBadRequest, "Insufficient resource")
			return
		}

		resource.CurrentValue -= req.Amount

		if err := h.resourceRepo.Update(resource); err != nil {
			log.Error().Err(err).Msg("Failed to update resource")
			respondError(w, http.StatusInternalServerError, "Failed to update resource")
			return
		}

		respondJSON(w, http.StatusOK, resource)
	}
}

// GainResource increments a resource value
func (h *resourceHandler) GainResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		key := chi.URLParam(r, "key")
		if key == "" {
			respondError(w, http.StatusBadRequest, "Resource key required")
			return
		}

		var req ResourceDeltaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if character, err := h.characterRepo.FindByID(characterID); err == nil && h.resourceService != nil {
			if ensureErr := h.resourceService.EnsureTrackableClassResources(character); ensureErr != nil {
				log.Warn().Err(ensureErr).Str("characterID", characterID.String()).Msg("EnsureTrackableClassResources before gain")
			}
		}

		resource, err := h.resourceRepo.FindByCharacterAndKey(characterID, key)
		if err != nil {
			respondError(w, http.StatusNotFound, "Resource not found")
			return
		}

		resource.CurrentValue += req.Amount

		// Cap at max if set
		if resource.MaxValue != nil && resource.CurrentValue > *resource.MaxValue {
			resource.CurrentValue = *resource.MaxValue
		}

		if err := h.resourceRepo.Update(resource); err != nil {
			log.Error().Err(err).Msg("Failed to update resource")
			respondError(w, http.StatusInternalServerError, "Failed to update resource")
			return
		}

		respondJSON(w, http.StatusOK, resource)
	}
}

// DeleteResource removes a resource
func (h *resourceHandler) DeleteResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		key := chi.URLParam(r, "key")
		if key == "" {
			respondError(w, http.StatusBadRequest, "Resource key required")
			return
		}

		resource, err := h.resourceRepo.FindByCharacterAndKey(characterID, key)
		if err != nil {
			respondError(w, http.StatusNotFound, "Resource not found")
			return
		}

		if err := h.resourceRepo.Delete(resource.ID); err != nil {
			log.Error().Err(err).Msg("Failed to delete resource")
			respondError(w, http.StatusInternalServerError, "Failed to delete resource")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}
