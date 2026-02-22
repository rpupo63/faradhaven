package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rs/zerolog/log"
)

type partyHandler struct {
	partyRepo     database.PartyRepository
	characterRepo database.CharacterRepository
	beastRepo     database.BeastRepository
}

func newPartyHandler(
	partyRepo database.PartyRepository,
	characterRepo database.CharacterRepository,
	beastRepo database.BeastRepository,
) *partyHandler {
	return &partyHandler{
		partyRepo:     partyRepo,
		characterRepo: characterRepo,
		beastRepo:     beastRepo,
	}
}

// createParty handles POST /api/parties
func (h *partyHandler) createParty() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerIDStr, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		ownerID, err := uuid.Parse(ownerIDStr)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Invalid user ID in context")
			return
		}

		var req CreatePartyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Name == "" {
			respondError(w, http.StatusBadRequest, "Party name is required")
			return
		}

		party := &models.Party{
			Name:    req.Name,
			OwnerID: ownerID,
		}

		if err := h.partyRepo.Add(party); err != nil {
			log.Error().Err(err).Msg("Failed to create party")
			respondError(w, http.StatusInternalServerError, "Failed to create party")
			return
		}

		respondJSON(w, http.StatusCreated, party)
	}
}

// getParty handles GET /api/parties/{partyID}
func (h *partyHandler) getParty() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyIDStr := chi.URLParam(r, "partyID")
		partyID, err := uuid.Parse(partyIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid party ID")
			return
		}

		party, err := h.partyRepo.FindByID(partyID)
		if err != nil {
			log.Error().Err(err).Str("partyID", partyIDStr).Msg("Failed to get party")
			respondError(w, http.StatusNotFound, "Party not found")
			return
		}

		// Ensure the user is either the owner or a member of the party
		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		authUUID, _ := uuid.Parse(authUserID)

		if party.OwnerID != authUUID {
			// Check if any character owned by auth user is in this party
			isMember := false
			characters, err := h.characterRepo.FindByUserID(authUUID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to check user characters for party membership")
				respondError(w, http.StatusInternalServerError, "Failed to authorize party access")
				return
			}
			for _, char := range characters {
				if char.PartyID != nil && *char.PartyID == partyID {
					isMember = true
					break
				}
			}
			if !isMember {
				respondError(w, http.StatusForbidden, "Forbidden")
				return
			}
		}

		// Preload members and identified beasts
		members, err := h.partyRepo.GetMembers(partyID)
		if err != nil {
			log.Error().Err(err).Str("partyID", partyIDStr).Msg("Failed to get party members")
			respondError(w, http.StatusInternalServerError, "Failed to retrieve party details")
			return
		}
		identifiedBeasts, err := h.partyRepo.GetIdentifiedBeasts(partyID)
		if err != nil {
			log.Error().Err(err).Str("partyID", partyIDStr).Msg("Failed to get identified beasts")
			respondError(w, http.StatusInternalServerError, "Failed to retrieve party details")
			return
		}

		// Create a response struct that includes members and identified beasts
		type PartyResponse struct {
			models.Party
			Members          []*models.Character `json:"members"`
			IdentifiedBeasts []*models.Beast     `json:"identified_beasts"`
		}

		resp := PartyResponse{
			Party:            *party,
			Members:          members,
			IdentifiedBeasts: identifiedBeasts,
		}

		respondJSON(w, http.StatusOK, resp)
	}
}

// getPartiesByOwner handles GET /api/users/{userID}/parties
func (h *partyHandler) getPartiesByOwner() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerIDStr := chi.URLParam(r, "userID")
		ownerID, err := uuid.Parse(ownerIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != ownerIDStr {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		parties, err := h.partyRepo.FindByOwnerID(ownerID)
		if err != nil {
			log.Error().Err(err).Str("ownerID", ownerIDStr).Msg("Failed to get parties by owner")
			respondError(w, http.StatusInternalServerError, "Failed to get parties")
			return
		}

		respondJSON(w, http.StatusOK, parties)
	}
}

// updateParty handles PUT /api/parties/{partyID}
func (h *partyHandler) updateParty() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyIDStr := chi.URLParam(r, "partyID")
		partyID, err := uuid.Parse(partyIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid party ID")
			return
		}

		party, err := h.partyRepo.FindByID(partyID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Party not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != party.OwnerID.String() {
			respondError(w, http.StatusForbidden, "Only the party owner can update the party")
			return
		}

		var req UpdatePartyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Name != nil {
			party.Name = *req.Name
		}

		if err := h.partyRepo.Update(party); err != nil {
			log.Error().Err(err).Msg("Failed to update party")
			respondError(w, http.StatusInternalServerError, "Failed to update party")
			return
		}

		respondJSON(w, http.StatusOK, party)
	}
}

// deleteParty handles DELETE /api/parties/{partyID}
func (h *partyHandler) deleteParty() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyIDStr := chi.URLParam(r, "partyID")
		partyID, err := uuid.Parse(partyIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid party ID")
			return
		}

		party, err := h.partyRepo.FindByID(partyID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Party not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != party.OwnerID.String() {
			respondError(w, http.StatusForbidden, "Only the party owner can delete the party")
			return
		}

		if err := h.partyRepo.Delete(partyID); err != nil {
			log.Error().Err(err).Msg("Failed to delete party")
			respondError(w, http.StatusInternalServerError, "Failed to delete party")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Party deleted successfully"})
	}
}

// addCharacterToParty handles POST /api/parties/{partyID}/members
func (h *partyHandler) addCharacterToParty() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyIDStr := chi.URLParam(r, "partyID")
		partyID, err := uuid.Parse(partyIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid party ID")
			return
		}

		party, err := h.partyRepo.FindByID(partyID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Party not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != party.OwnerID.String() {
			respondError(w, http.StatusForbidden, "Only the party owner can add members")
			return
		}

		var req AddCharacterToPartyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		character, err := h.characterRepo.FindByID(req.CharacterID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		// Check if character belongs to the authenticated user
		if character.UserID.String() != authUserID {
			respondError(w, http.StatusForbidden, "Cannot add character owned by another user")
			return
		}

		if err := h.partyRepo.AddMember(partyID, req.CharacterID); err != nil {
			log.Error().Err(err).Msg("Failed to add character to party")
			respondError(w, http.StatusInternalServerError, "Failed to add character to party")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Character added to party successfully"})
	}
}

// removeCharacterFromParty handles DELETE /api/parties/{partyID}/members/{characterID}
func (h *partyHandler) removeCharacterFromParty() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyIDStr := chi.URLParam(r, "partyID")
		partyID, err := uuid.Parse(partyIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid party ID")
			return
		}

		characterIDStr := chi.URLParam(r, "characterID")
		characterID, err := uuid.Parse(characterIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		party, err := h.partyRepo.FindByID(partyID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Party not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != party.OwnerID.String() {
			respondError(w, http.StatusForbidden, "Only the party owner can remove members")
			return
		}

		// Verify character is actually in this party
		character, err := h.characterRepo.FindByID(characterID)
		if err != nil || character.PartyID == nil || *character.PartyID != partyID {
			respondError(w, http.StatusBadRequest, "Character not found in this party")
			return
		}

		if err := h.partyRepo.RemoveMember(partyID, characterID); err != nil {
			log.Error().Err(err).Msg("Failed to remove character from party")
			respondError(w, http.StatusInternalServerError, "Failed to remove character from party")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Character removed from party successfully"})
	}
}

// addIdentifiedBeast handles POST /api/parties/{partyID}/identified-beasts
func (h *partyHandler) addIdentifiedBeast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyIDStr := chi.URLParam(r, "partyID")
		partyID, err := uuid.Parse(partyIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid party ID")
			return
		}

		party, err := h.partyRepo.FindByID(partyID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Party not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != party.OwnerID.String() {
			respondError(w, http.StatusForbidden, "Only the party owner can add identified beasts")
			return
		}

		var req AddIdentifiedBeastRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		_, err = h.beastRepo.FindByID(req.BeastID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Beast not found")
			return
		}

		if err := h.partyRepo.AddIdentifiedBeast(partyID, req.BeastID); err != nil {
			log.Error().Err(err).Msg("Failed to add identified beast to party")
			respondError(w, http.StatusInternalServerError, "Failed to add identified beast")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Beast added to identified list successfully"})
	}
}

// removeIdentifiedBeast handles DELETE /api/parties/{partyID}/identified-beasts/{beastID}
func (h *partyHandler) removeIdentifiedBeast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyIDStr := chi.URLParam(r, "partyID")
		partyID, err := uuid.Parse(partyIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid party ID")
			return
		}

		beastIDStr := chi.URLParam(r, "beastID")
		beastID, err := uuid.Parse(beastIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid beast ID")
			return
		}

		party, err := h.partyRepo.FindByID(partyID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Party not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != party.OwnerID.String() {
			respondError(w, http.StatusForbidden, "Only the party owner can remove identified beasts")
			return
		}

		if err := h.partyRepo.RemoveIdentifiedBeast(partyID, beastID); err != nil {
			log.Error().Err(err).Msg("Failed to remove identified beast from party")
			respondError(w, http.StatusInternalServerError, "Failed to remove identified beast")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Beast removed from identified list successfully"})
	}
}

// setCharacterParty handles PUT /api/characters/{characterID}/party
func (h *partyHandler) setCharacterParty() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterIDStr := chi.URLParam(r, "characterID")
		characterID, err := uuid.Parse(characterIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByID(characterID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		// Verify user owns this character
		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != character.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		var req SetCharacterPartyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.PartyID != nil {
			// Check if the party exists
			party, err := h.partyRepo.FindByID(*req.PartyID)
			if err != nil {
				respondError(w, http.StatusNotFound, "Party not found")
				return
			}
			// Only party owner can assign characters to their party
			if party.OwnerID.String() != authUserID {
				respondError(w, http.StatusForbidden, "Cannot assign character to a party not owned by you")
				return
			}
		}

		character.PartyID = req.PartyID
		if err := h.characterRepo.Update(character); err != nil {
			log.Error().Err(err).Msg("Failed to set character party")
			respondError(w, http.StatusInternalServerError, "Failed to set character party")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message":    "Character party updated successfully",
			"character_id": character.ID,
			"party_id":   character.PartyID,
		})
	}
}
