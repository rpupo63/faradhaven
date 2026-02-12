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

type mapHandler struct {
	mapRepo *database.MapRepo
}

func newMapHandler(mapRepo *database.MapRepo) *mapHandler {
	return &mapHandler{mapRepo: mapRepo}
}

// createMap handles creating a new game map
func (h *mapHandler) createMap() http.HandlerFunc {
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

		var req CreateMapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		gameMap := models.GameMap{
			OwnerID:       userID,
			RoomCode:      req.RoomCode,
			Name:          req.Name,
			BackgroundURL: req.BackgroundURL,
			GridRows:      req.GridRows,
			GridCols:      req.GridCols,
			TileSize:      req.TileSize,
		}

		if err := h.mapRepo.CreateMap(&gameMap); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to create map")
			return
		}

		respondJSON(w, http.StatusCreated, gameMap)
	}
}

// getMap retrieves a map by its ID
func (h *mapHandler) getMap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mapIDStr := chi.URLParam(r, "mapID")
		mapID, err := uuid.Parse(mapIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		gameMap, err := h.mapRepo.GetMapByID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		respondJSON(w, http.StatusOK, gameMap)
	}
}

// getMapByRoom retrieves a map by its room code
func (h *mapHandler) getMapByRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomCode := chi.URLParam(r, "roomCode")

		gameMap, err := h.mapRepo.GetMapByRoomCode(roomCode)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		respondJSON(w, http.StatusOK, gameMap)
	}
}

// getUserMaps retrieves all maps owned by the current user
func (h *mapHandler) getUserMaps() http.HandlerFunc {
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

		maps, err := h.mapRepo.GetMapsByOwner(userID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to retrieve maps")
			return
		}

		respondJSON(w, http.StatusOK, maps)
	}
}

// updateMap updates a map (DM only)
func (h *mapHandler) updateMap() http.HandlerFunc {
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

		mapIDStr := chi.URLParam(r, "mapID")
		mapID, err := uuid.Parse(mapIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		gameMap, err := h.mapRepo.GetMapByID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if gameMap.OwnerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can update this map")
			return
		}

		var req UpdateMapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if req.Name != nil {
			gameMap.Name = *req.Name
		}
		if req.BackgroundURL != nil {
			gameMap.BackgroundURL = *req.BackgroundURL
		}
		if req.GridRows != nil {
			gameMap.GridRows = *req.GridRows
		}
		if req.GridCols != nil {
			gameMap.GridCols = *req.GridCols
		}
		if req.TileSize != nil {
			gameMap.TileSize = *req.TileSize
		}

		if err := h.mapRepo.UpdateMap(gameMap); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to update map")
			return
		}

		respondJSON(w, http.StatusOK, gameMap)
	}
}

// deleteMap deletes a map (DM only)
func (h *mapHandler) deleteMap() http.HandlerFunc {
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

		mapIDStr := chi.URLParam(r, "mapID")
		mapID, err := uuid.Parse(mapIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		gameMap, err := h.mapRepo.GetMapByID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if gameMap.OwnerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can delete this map")
			return
		}

		if err := h.mapRepo.DeleteMap(mapID); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to delete map")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Map deleted successfully"})
	}
}

// addToken adds a token to the map (DM only)
func (h *mapHandler) addToken() http.HandlerFunc {
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

		mapIDStr := chi.URLParam(r, "mapID")
		mapID, err := uuid.Parse(mapIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		gameMap, err := h.mapRepo.GetMapByID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if gameMap.OwnerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can add tokens")
			return
		}

		var req CreateTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		token := models.MapToken{
			MapID:          mapID,
			CharacterID:    req.CharacterID,
			AssignedUserID: req.AssignedUserID,
			Name:           req.Name,
			ImageURL:       req.ImageURL,
			TokenType:      req.TokenType,
			GridX:          req.GridX,
			GridY:          req.GridY,
			Size:           req.Size,
			Color:          req.Color,
			Visible:        req.Visible,
		}

		if err := h.mapRepo.AddToken(&token); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to add token")
			return
		}

		respondJSON(w, http.StatusCreated, token)
	}
}

// updateToken updates a token (DM or Assigned Owner)
func (h *mapHandler) updateToken() http.HandlerFunc {
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

		mapIDStr := chi.URLParam(r, "mapID")
		mapID, err := uuid.Parse(mapIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}
		tokenIDStr := chi.URLParam(r, "tokenID")
		tokenID, err := uuid.Parse(tokenIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid token ID")
			return
		}

		gameMap, err := h.mapRepo.GetMapByID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		token, err := h.mapRepo.GetTokenByID(tokenID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Token not found")
			return
		}

		// Authorization check: DM or Token Owner
		isDM := gameMap.OwnerID == userID
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

		// DM-only updates
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

		// Allowed for Owner & DM: Movement
		if req.GridX != nil {
			token.GridX = *req.GridX
		}
		if req.GridY != nil {
			token.GridY = *req.GridY
		}

		if err := h.mapRepo.UpdateToken(token); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to update token")
			return
		}

		respondJSON(w, http.StatusOK, token)
	}
}

// deleteToken deletes a token (DM only)
func (h *mapHandler) deleteToken() http.HandlerFunc {
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

		mapIDStr := chi.URLParam(r, "mapID")
		mapID, err := uuid.Parse(mapIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}
		tokenIDStr := chi.URLParam(r, "tokenID")
		tokenID, err := uuid.Parse(tokenIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid token ID")
			return
		}

		gameMap, err := h.mapRepo.GetMapByID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if gameMap.OwnerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can delete tokens")
			return
		}

		if err := h.mapRepo.DeleteToken(tokenID); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to delete token")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Token deleted successfully"})
	}
}

// getInitiative returns the initiative order for a map
func (h *mapHandler) getInitiative() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mapID, err := uuid.Parse(chi.URLParam(r, "mapID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		tokens, err := h.mapRepo.GetTokensByInitiativeOrder(mapID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get initiative order")
			respondError(w, http.StatusInternalServerError, "Failed to get initiative order")
			return
		}

		respondJSON(w, http.StatusOK, tokens)
	}
}

// setInitiative sets the initiative order for tokens on a map (DM only)
func (h *mapHandler) setInitiative() http.HandlerFunc {
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

		gameMap, err := h.mapRepo.GetMapByID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if gameMap.OwnerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can set initiative")
			return
		}

		var req SetInitiativeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		for _, entry := range req.Entries {
			order := entry.Order
			if err := h.mapRepo.SetTokenInitiativeOrder(entry.TokenID, &order); err != nil {
				log.Error().Err(err).Str("tokenID", entry.TokenID.String()).Msg("Failed to set initiative")
			}
		}

		// Return updated initiative order
		tokens, err := h.mapRepo.GetTokensByInitiativeOrder(mapID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to get updated initiative")
			return
		}

		respondJSON(w, http.StatusOK, tokens)
	}
}

// clearInitiative clears initiative order for all tokens on a map (DM only)
func (h *mapHandler) clearInitiative() http.HandlerFunc {
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

		gameMap, err := h.mapRepo.GetMapByID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if gameMap.OwnerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can clear initiative")
			return
		}

		// Clear initiative for all tokens on this map
		for _, token := range gameMap.Tokens {
			_ = h.mapRepo.SetTokenInitiativeOrder(token.ID, nil)
		}

		respondJSON(w, http.StatusOK, map[string]string{"status": "initiative cleared"})
	}
}
