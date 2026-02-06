package api

import (
	"net/http"

	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rs/zerolog/log"
)

type weaponHandler struct {
	repo *database.WeaponRepo
}

func newWeaponHandler(repo *database.WeaponRepo) *weaponHandler {
	return &weaponHandler{repo: repo}
}

// getAllWeapons handles GET /api/weapons
func (h *weaponHandler) getAllWeapons() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		weapons, err := h.repo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get weapons")
			respondError(w, http.StatusInternalServerError, "Failed to get weapons")
			return
		}
		respondJSON(w, http.StatusOK, weapons)
	}
}
