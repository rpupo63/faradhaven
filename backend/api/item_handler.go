package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rs/zerolog/log"
)

type itemHandler struct {
	itemRepo database.ItemRepository
}

func newItemHandler(itemRepo database.ItemRepository) *itemHandler {
	return &itemHandler{itemRepo: itemRepo}
}

func (h *itemHandler) getAllItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := h.itemRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get items")
			respondError(w, http.StatusInternalServerError, "Failed to get items")
			return
		}
		respondJSON(w, http.StatusOK, items)
	}
}

func (h *itemHandler) getItemByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "itemID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid item ID")
			return
		}

		item, err := h.itemRepo.FindByID(id)
		if err != nil {
			log.Error().Err(err).Str("itemID", idStr).Msg("Failed to get item")
			respondError(w, http.StatusNotFound, "Item not found")
			return
		}

		respondJSON(w, http.StatusOK, item)
	}
}
