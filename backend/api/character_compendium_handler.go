package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_classes"
	"github.com/rs/zerolog/log"
)

// getAllRaces returns all races (for character creation dropdowns)
func (h *characterHandler) getAllRaces() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		races, err := h.raceRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get races")
			respondError(w, http.StatusInternalServerError, "Failed to get races")
			return
		}
		respondJSON(w, http.StatusOK, races)
	}
}

// getRaceByID returns a race with traits (for D&D Beyond-style compendium detail view)
func (h *characterHandler) getRaceByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "raceID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid race ID")
			return
		}

		race, err := h.raceRepo.FindByIDWithTraits(id)
		if err != nil {
			log.Error().Err(err).Str("raceID", idStr).Msg("Failed to get race")
			respondError(w, http.StatusNotFound, "Race not found")
			return
		}

		respondJSON(w, http.StatusOK, race)
	}
}

// getAllClasses returns all classes (for character creation dropdowns)
func (h *characterHandler) getAllClasses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		classes, err := h.classRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get classes")
			respondError(w, http.StatusInternalServerError, "Failed to get classes")
			return
		}
		respondJSON(w, http.StatusOK, classes)
	}
}

// getClassByID returns a class with all levels 1-20 (for D&D Beyond-style compendium)
func (h *characterHandler) getClassByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "classID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid class ID")
			return
		}

		class, err := h.classRepo.FindByIDWithLevels(id)
		if err != nil {
			log.Error().Err(err).Str("classID", idStr).Msg("Failed to get class")
			respondError(w, http.StatusNotFound, "Class not found")
			return
		}

		resp := ClassWithLevelsResponse{
			Class:  *class, // Assign the dereferenced models.Class
			Levels: class.Levels,
		}

		if class.Name == "The Lorewright" {
			resp.MadnessTable = faradhaven_classes.LorewrightMadnessTable()
		}

		respondJSON(w, http.StatusOK, resp)
	}
}
