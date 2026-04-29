package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/services"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type storeOwnerHandler struct {
	storeOwnerRepo database.StoreOwnerRepository
	s3             *services.S3Service
}

func newStoreOwnerHandler(storeOwnerRepo database.StoreOwnerRepository, s3 *services.S3Service) *storeOwnerHandler {
	return &storeOwnerHandler{storeOwnerRepo: storeOwnerRepo, s3: s3}
}

func (h *storeOwnerHandler) enrichImageURL(o *models.StoreOwner) {
	if h.s3 == nil || o == nil {
		return
	}
	if u := h.s3.StoreownerPortraitURL(o.Name); u != "" {
		o.ImageURL = u
	}
}

func (h *storeOwnerHandler) getStoreOwners() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owners, err := h.storeOwnerRepo.FindAllWithRules()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get store owners")
			respondError(w, http.StatusInternalServerError, "Failed to get store owners")
			return
		}
		for i := range owners {
			h.enrichImageURL(&owners[i])
		}
		respondJSON(w, http.StatusOK, owners)
	}
}

// getStoreOwnerPortrait redirects to the public S3 object for this vendor's portrait
// (storeowners/{name}.png), so clients can use a stable backend URL if desired.
func (h *storeOwnerHandler) getStoreOwnerPortrait() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "storeOwnerID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid store owner ID")
			return
		}
		owner, err := h.storeOwnerRepo.FindByIDWithRules(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				respondError(w, http.StatusNotFound, "Store owner not found")
				return
			}
			log.Error().Err(err).Msg("Failed to load store owner for portrait")
			respondError(w, http.StatusInternalServerError, "Failed to load store owner")
			return
		}
		if h.s3 == nil {
			respondError(w, http.StatusServiceUnavailable, "Image URLs are not configured")
			return
		}
		u := h.s3.StoreownerPortraitURL(owner.Name)
		if u == "" {
			respondError(w, http.StatusNotFound, "Portrait not available")
			return
		}
		http.Redirect(w, r, u, http.StatusFound)
	}
}
