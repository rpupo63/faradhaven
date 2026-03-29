package api

import (
	"net/http"

	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rs/zerolog/log"
)

type storeOwnerHandler struct {
	storeOwnerRepo database.StoreOwnerRepository
}

func newStoreOwnerHandler(storeOwnerRepo database.StoreOwnerRepository) *storeOwnerHandler {
	return &storeOwnerHandler{storeOwnerRepo: storeOwnerRepo}
}

func (h *storeOwnerHandler) getStoreOwners() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owners, err := h.storeOwnerRepo.FindAllWithRules()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get store owners")
			respondError(w, http.StatusInternalServerError, "Failed to get store owners")
			return
		}
		respondJSON(w, http.StatusOK, owners)
	}
}
