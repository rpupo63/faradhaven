package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rs/zerolog/log"
)

type componentHandler struct {
	repo *database.ComponentRepo
}

func newComponentHandler(repo *database.ComponentRepo) *componentHandler {
	return &componentHandler{repo: repo}
}

// getAllComponents handles GET /api/components
// Returns all spell system components for the periodic table visualization.
func (h *componentHandler) getAllComponents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		components, err := h.repo.GetAllComponents()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get components")
			respondError(w, http.StatusInternalServerError, "Failed to get components")
			return
		}
		respondJSON(w, http.StatusOK, components)
	}
}

// getComponentByID handles GET /api/components/{componentID}
func (h *componentHandler) getComponentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "componentID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid component ID")
			return
		}

		component, err := h.repo.GetComponentByID(id)
		if err != nil {
			log.Error().Err(err).Str("id", idStr).Msg("Failed to get component")
			respondError(w, http.StatusNotFound, "Component not found")
			return
		}
		respondJSON(w, http.StatusOK, component)
	}
}

// getComponentsByCategory handles GET /api/components/category/{category}
func (h *componentHandler) getComponentsByCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryStr := chi.URLParam(r, "category")
		category := models.ComponentCategory(categoryStr)

		components, err := h.repo.GetComponentsByCategory(category)
		if err != nil {
			log.Error().Err(err).Str("category", categoryStr).Msg("Failed to get components by category")
			respondError(w, http.StatusInternalServerError, "Failed to get components")
			return
		}
		respondJSON(w, http.StatusOK, components)
	}
}
