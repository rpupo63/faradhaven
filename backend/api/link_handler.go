package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type linkHandler struct {
	service services.LinkService
}

func newLinkHandler(service services.LinkService) *linkHandler {
	return &linkHandler{service}
}

type CreateLinkRequest struct {
	TargetCharacterID uuid.UUID       `json:"target_character_id"`
	LinkType          models.LinkType `json:"link_type"`
	Notes             *string         `json:"notes,omitempty"`
}

func (h *linkHandler) createLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sourceCharacterIDStr := chi.URLParam(r, "characterID")
		sourceCharacterID, err := uuid.Parse(sourceCharacterIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid source character ID")
			return
		}

		var req CreateLinkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		link, err := h.service.CreateLink(sourceCharacterID, req.TargetCharacterID, req.LinkType, req.Notes)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create link")
			respondError(w, http.StatusInternalServerError, "Failed to create link")
			return
		}

		respondJSON(w, http.StatusCreated, link)
	}
}

func (h *linkHandler) removeLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		linkIDStr := chi.URLParam(r, "linkID")
		linkID, err := uuid.Parse(linkIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid link ID")
			return
		}

		err = h.service.RemoveLink(linkID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to remove link")
			respondError(w, http.StatusInternalServerError, "Failed to remove link")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Link removed successfully"})
	}
}

func (h *linkHandler) getLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterIDStr := chi.URLParam(r, "characterID")
		characterID, err := uuid.Parse(characterIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		links, err := h.service.GetLinks(characterID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get links")
			respondError(w, http.StatusInternalServerError, "Failed to get links")
			return
		}

		respondJSON(w, http.StatusOK, links)
	}
}
