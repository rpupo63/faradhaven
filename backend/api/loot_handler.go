package api

import (
	"errors"
	"net/http"

	"encoding/json" // This import is used by json.NewDecoder(r.Body).Decode(&req)

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/errs" // Added errs import
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/services"
	"github.com/rs/zerolog"
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

// GenerateLoot handles preview generation of loot for a character (no pickup yet).
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

	preview, err := h.lootService.GenerateLootPreview(req.CharacterID, req.Source, req.RoomTheme, req.Location, req.LootLevel)
	if err != nil {
		if errors.Is(err, services.ErrInvalidLootOption) {
			h.responder.WriteError(w, errs.NewApiErr(http.StatusBadRequest, err.Error()))
			return
		}
		h.responder.WriteError(w, errs.NewApiErr(http.StatusInternalServerError, err.Error()))
		return
	}

	members := make([]PartyLootMember, 0, len(preview.PartyMembers))
	for _, m := range preview.PartyMembers {
		members = append(members, PartyLootMember{ID: m.ID, Name: m.Name})
	}
	h.responder.WriteJSON(w, GenerateLootPreviewResponse{
		SessionID:    preview.SessionID,
		Loot:         *preview.Loot,
		PartyMembers: members,
	})
}

func (h *LootHandler) ConfirmLootPickup(w http.ResponseWriter, r *http.Request) {
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
	var req ConfirmLootPickupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.WriteError(w, errs.NewApiErr(http.StatusBadRequest, "Invalid request payload"))
		return
	}
	assignments := make([]services.LootAssignmentRequest, 0, len(req.Assignments))
	for _, a := range req.Assignments {
		assignments = append(assignments, services.LootAssignmentRequest{
			DropIndex:   a.DropIndex,
			CharacterID: a.CharacterID,
		})
	}
	result, err := h.lootService.ConfirmLootPickup(characterID, req.SessionID, assignments)
	if err != nil {
		if errors.Is(err, services.ErrInvalidLootOption) {
			h.responder.WriteError(w, errs.NewApiErr(http.StatusBadRequest, err.Error()))
			return
		}
		h.responder.WriteError(w, errs.NewApiErr(http.StatusInternalServerError, err.Error()))
		return
	}
	h.responder.WriteJSON(w, result)
}

func (h *LootHandler) GetLootOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.responder.WriteError(w, errs.NewApiErr(http.StatusMethodNotAllowed, "Method not allowed"))
		return
	}
	resp := LootOptionsResponse{
		Themes:     make([]string, 0),
		Locations:  make([]string, 0),
		Sources:    make([]string, 0),
		LootLevels: make([]int, 0, 20),
	}
	for _, v := range models.AllLootThemes() {
		resp.Themes = append(resp.Themes, string(v))
	}
	for _, v := range models.AllLootLocations() {
		resp.Locations = append(resp.Locations, string(v))
	}
	for _, v := range models.AllLootSources() {
		resp.Sources = append(resp.Sources, string(v))
	}
	for i := 1; i <= 20; i++ {
		resp.LootLevels = append(resp.LootLevels, i)
	}
	h.responder.WriteJSON(w, resp)
}
