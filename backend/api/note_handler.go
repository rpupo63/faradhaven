package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rs/zerolog/log"
)

type noteHandler struct {
	noteRepo database.NoteRepository
}

func newNoteHandler(noteRepo database.NoteRepository) *noteHandler {
	return &noteHandler{noteRepo: noteRepo}
}

func (h *noteHandler) getAllNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notes, err := h.noteRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get notes")
			respondError(w, http.StatusInternalServerError, "Failed to get notes")
			return
		}
		respondJSON(w, http.StatusOK, notes)
	}
}

func (h *noteHandler) createNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			UserID      string `json:"userId"`
			Username    string `json:"username"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		userUUID, err := uuid.Parse(req.UserID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		note := &models.SharedNote{
			Title:       req.Title,
			Description: req.Description,
			UserID:      userUUID,
			Username:    req.Username,
		}

		if err := h.noteRepo.Add(note); err != nil {
			log.Error().Err(err).Msg("Failed to create note")
			respondError(w, http.StatusInternalServerError, "Failed to create note")
			return
		}

		respondJSON(w, http.StatusCreated, note)
	}
}
