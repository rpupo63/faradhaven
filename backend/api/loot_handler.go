package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/errs" // Added errs import
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog"
	"encoding/json" // This import is used by json.NewDecoder(r.Body).Decode(&req)
)

type LootHandler struct {
	lootService *services.LootService
	responder   Responder
}

// newLootHandler creates a new instance of LootHandler.
// This is a private constructor used within the api package.
func newLootHandler(lootService *services.LootService, logger zerolog.Logger) *LootHandler {
	return &LootHandler{
		lootService: lootService,
		responder:   NewResponder(logger),
	}
}

// GenerateLoot handles the generation of loot for a character.
func (h *LootHandler) GenerateLoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.responder.WriteError(w, errs.NewApiErr(http.StatusMethodNotAllowed, "Method not allowed"))
		return
	}

	characterIDStr := chi.URLParam(r, "characterID")
	characterID, err := uuid.Parse(characterIDStr)
	if err != nil {
		h.responder.WriteError(w, errs.NewApiErr(http.StatusBadRequest, "Invalid character ID format"))
		return
	}

	var req GenerateLootRequest // Use the struct defined in types.go
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.WriteError(w, errs.NewApiErr(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	// Override req.CharacterID with the one from the URL path to ensure consistency
	req.CharacterID = characterID

	lootResult, err := h.lootService.GenerateLoot(req.CharacterID, req.Source, req.Tier)
	if err != nil {
		h.responder.WriteError(w, errs.NewApiErr(http.StatusInternalServerError, err.Error()))
		return
	}

	h.responder.WriteJSON(w, lootResult)
}