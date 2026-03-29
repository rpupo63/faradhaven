package api

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type noteHandler struct {
	noteRepo   database.NoteRepository
	s3Service  *services.S3Service
}

func newNoteHandler(noteRepo database.NoteRepository, s3Service *services.S3Service) *noteHandler {
	return &noteHandler{noteRepo: noteRepo, s3Service: s3Service}
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
		authUserIDStr, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		contentType := r.Header.Get("Content-Type")
		var title, description, userIDStr, username string
		var pdfFile multipart.File
		var pdfHeader *multipart.FileHeader

		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(20 << 20); err != nil {
				respondError(w, http.StatusBadRequest, "Invalid multipart form")
				return
			}
			title = strings.TrimSpace(r.FormValue("title"))
			description = strings.TrimSpace(r.FormValue("description"))
			userIDStr = r.FormValue("userId")
			username = strings.TrimSpace(r.FormValue("username"))

			f, hdr, ferr := r.FormFile("pdf")
			if ferr != nil {
				if !errors.Is(ferr, http.ErrMissingFile) {
					respondError(w, http.StatusBadRequest, "Invalid file")
					return
				}
			} else {
				pdfFile = f
				pdfHeader = hdr
				defer f.Close()
			}
		} else {
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
			title = strings.TrimSpace(req.Title)
			description = strings.TrimSpace(req.Description)
			userIDStr = req.UserID
			username = strings.TrimSpace(req.Username)
		}

		if title == "" || description == "" {
			respondError(w, http.StatusBadRequest, "Title and description are required")
			return
		}
		if userIDStr != authUserIDStr {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		note := &models.SharedNote{
			Title:       title,
			Description: description,
			UserID:      userUUID,
			Username:    username,
		}

		if pdfFile != nil {
			if h.s3Service == nil {
				respondError(w, http.StatusInternalServerError, "File upload is not available")
				return
			}
			url, err := h.s3Service.UploadBulletinPDF(pdfFile, pdfHeader)
			if err != nil {
				log.Error().Err(err).Msg("Failed to upload bulletin PDF")
				status := http.StatusInternalServerError
				msg := "Failed to upload PDF"
				if strings.HasPrefix(err.Error(), "invalid") {
					status = http.StatusBadRequest
					msg = err.Error()
				}
				respondError(w, status, msg)
				return
			}
			note.PdfURL = &url
		}

		if err := h.noteRepo.Add(note); err != nil {
			log.Error().Err(err).Msg("Failed to create note")
			respondError(w, http.StatusInternalServerError, "Failed to create note")
			return
		}

		respondJSON(w, http.StatusCreated, note)
	}
}
