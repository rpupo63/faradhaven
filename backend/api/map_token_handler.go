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

type mapTokenHandler struct {
	tokenRepo   *database.MapTokenRepo
	gameMapRepo *database.GameMapRepo
}

func newMapTokenHandler(tokenRepo *database.MapTokenRepo, gameMapRepo *database.GameMapRepo) *mapTokenHandler {
	return &mapTokenHandler{tokenRepo: tokenRepo, gameMapRepo: gameMapRepo}
}

func (h *mapTokenHandler) addToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Invalid user ID in context")
			return
		}

		mapID, err := uuid.Parse(chi.URLParam(r, "mapID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		ownerID, err := h.gameMapRepo.GetOwnerID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if ownerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can add tokens")
			return
		}

		var req CreateTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if req.CharacterID != nil && req.MonsterID != nil {
			respondError(w, http.StatusBadRequest, "Token cannot be linked to both a character and a monster")
			return
		}

		tt := req.TokenType
		if tt == "" {
			tt = models.MapTokenNPC
		} else if !tt.IsValid() {
			if p, ok := models.ParseMapTokenType(string(req.TokenType)); ok {
				tt = p
			} else {
				respondError(w, http.StatusBadRequest, "invalid token_type")
				return
			}
		}

		token := models.MapToken{
			MapID:          mapID,
			CharacterID:    req.CharacterID,
			MonsterID:      req.MonsterID,
			AssignedUserID: req.AssignedUserID,
			Name:           req.Name,
			ImageURL:       req.ImageURL,
			TokenType:      tt,
			GridX:          req.GridX,
			GridY:          req.GridY,
			Size:           req.Size,
			Color:          req.Color,
			Visible:        req.Visible,
		}

		if err := h.tokenRepo.Create(&token); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to add token")
			return
		}

		respondJSON(w, http.StatusCreated, token)
	}
}

func (h *mapTokenHandler) updateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Invalid user ID in context")
			return
		}

		mapID, err := uuid.Parse(chi.URLParam(r, "mapID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}
		tokenID, err := uuid.Parse(chi.URLParam(r, "tokenID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid token ID")
			return
		}

		ownerID, err := h.gameMapRepo.GetOwnerID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		token, err := h.tokenRepo.GetByID(tokenID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Token not found")
			return
		}

		isDM := ownerID == userID
		isOwner := token.AssignedUserID != nil && *token.AssignedUserID == userID

		if !isDM && !isOwner {
			respondError(w, http.StatusForbidden, "You do not have permission to move this token")
			return
		}

		var req UpdateTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if req.Visible != nil {
			if !isDM {
				respondError(w, http.StatusForbidden, "Only DM can change visibility")
				return
			}
			token.Visible = *req.Visible
		}
		if req.Size != nil {
			if !isDM {
				respondError(w, http.StatusForbidden, "Only DM can change size")
				return
			}
			token.Size = *req.Size
		}
		if req.Color != nil {
			if !isDM {
				respondError(w, http.StatusForbidden, "Only DM can change color")
				return
			}
			token.Color = *req.Color
		}
		if req.AssignedUserID != nil {
			if !isDM {
				respondError(w, http.StatusForbidden, "Only DM can reassign tokens")
				return
			}
			token.AssignedUserID = req.AssignedUserID
		}
		if req.CharacterID != nil {
			if !isDM {
				respondError(w, http.StatusForbidden, "Only DM can change token entity")
				return
			}
			token.CharacterID = req.CharacterID
		}
		if req.MonsterID != nil {
			if !isDM {
				respondError(w, http.StatusForbidden, "Only DM can change token entity")
				return
			}
			token.MonsterID = req.MonsterID
		}
		if req.Name != nil {
			if !isDM {
				respondError(w, http.StatusForbidden, "Only DM can rename tokens")
				return
			}
			token.Name = *req.Name
		}
		if req.ImageURL != nil {
			if !isDM {
				respondError(w, http.StatusForbidden, "Only DM can change token image")
				return
			}
			token.ImageURL = *req.ImageURL
		}
		if req.TokenType != nil {
			if !isDM {
				respondError(w, http.StatusForbidden, "Only DM can change token type")
				return
			}
			tt := *req.TokenType
			if !tt.IsValid() {
				if p, ok := models.ParseMapTokenType(string(tt)); ok {
					tt = p
				} else {
					respondError(w, http.StatusBadRequest, "invalid token_type")
					return
				}
			}
			token.TokenType = tt
		}

		if req.GridX != nil {
			token.GridX = *req.GridX
		}
		if req.GridY != nil {
			token.GridY = *req.GridY
		}

		if err := h.tokenRepo.Update(token); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to update token")
			return
		}

		respondJSON(w, http.StatusOK, token)
	}
}

func (h *mapTokenHandler) deleteToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Invalid user ID in context")
			return
		}

		mapID, err := uuid.Parse(chi.URLParam(r, "mapID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}
		tokenID, err := uuid.Parse(chi.URLParam(r, "tokenID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid token ID")
			return
		}

		ownerID, err := h.gameMapRepo.GetOwnerID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if ownerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can delete tokens")
			return
		}

		if err := h.tokenRepo.Delete(tokenID); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to delete token")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Token deleted successfully"})
	}
}

func (h *mapTokenHandler) getInitiative() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mapID, err := uuid.Parse(chi.URLParam(r, "mapID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		tokens, err := h.tokenRepo.GetByInitiativeOrder(mapID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get initiative order")
			respondError(w, http.StatusInternalServerError, "Failed to get initiative order")
			return
		}

		respondJSON(w, http.StatusOK, tokens)
	}
}

func (h *mapTokenHandler) setInitiative() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Invalid user ID in context")
			return
		}

		mapID, err := uuid.Parse(chi.URLParam(r, "mapID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		ownerID, err := h.gameMapRepo.GetOwnerID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if ownerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can set initiative")
			return
		}

		var req SetInitiativeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		dbEntries := make([]database.InitiativeEntry, len(req.Entries))
		for i, e := range req.Entries {
			dbEntries[i] = database.InitiativeEntry{TokenID: e.TokenID, Order: e.Order}
		}
		if err := h.tokenRepo.BulkSetInitiativeOrder(dbEntries); err != nil {
			log.Error().Err(err).Msg("Failed to set initiative order")
			respondError(w, http.StatusInternalServerError, "Failed to set initiative")
			return
		}

		tokens, err := h.tokenRepo.GetByInitiativeOrder(mapID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to get updated initiative")
			return
		}

		respondJSON(w, http.StatusOK, tokens)
	}
}

func (h *mapTokenHandler) clearInitiative() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Invalid user ID in context")
			return
		}

		mapID, err := uuid.Parse(chi.URLParam(r, "mapID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		ownerID, err := h.gameMapRepo.GetOwnerID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if ownerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can clear initiative")
			return
		}

		if err := h.tokenRepo.ClearInitiativeByMapID(mapID); err != nil {
			log.Error().Err(err).Str("mapID", mapID.String()).Msg("Failed to clear initiative")
			respondError(w, http.StatusInternalServerError, "Failed to clear initiative")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"status": "initiative cleared"})
	}
}
