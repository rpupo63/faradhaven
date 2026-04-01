package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

type gameMapHandler struct {
	gameMapRepo *database.GameMapRepo
	hub         *Hub
}

func newGameMapHandler(gameMapRepo *database.GameMapRepo, hub *Hub) *gameMapHandler {
	return &gameMapHandler{gameMapRepo: gameMapRepo, hub: hub}
}

func (h *gameMapHandler) createMap() http.HandlerFunc {
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

		if err := h.gameMapRepo.Create(&gameMap); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to create map")
			return
		}

		respondJSON(w, http.StatusCreated, gameMap)
	}
}

func (h *gameMapHandler) getMap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mapID, err := uuid.Parse(chi.URLParam(r, "mapID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		gameMap, err := h.gameMapRepo.GetByID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		respondJSON(w, http.StatusOK, gameMap)
	}
}

func (h *gameMapHandler) getMapByRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomCode := chi.URLParam(r, "roomCode")

		gameMap, err := h.gameMapRepo.GetByRoomCode(roomCode)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		respondJSON(w, http.StatusOK, gameMap)
	}
}

func (h *gameMapHandler) getUserMaps() http.HandlerFunc {
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

		owned, err := h.gameMapRepo.GetByOwner(userID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to retrieve maps")
			return
		}

		joined, err := h.gameMapRepo.GetByAssignedUser(userID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to retrieve joined maps")
			return
		}

		// Merge, deduplicating by ID (owned takes precedence)
		seen := make(map[uuid.UUID]struct{}, len(owned))
		result := make([]models.GameMap, 0, len(owned)+len(joined))
		for _, m := range owned {
			seen[m.ID] = struct{}{}
			result = append(result, m)
		}
		for _, m := range joined {
			if _, ok := seen[m.ID]; !ok {
				result = append(result, m)
			}
		}

		respondJSON(w, http.StatusOK, result)
	}
}

func (h *gameMapHandler) updateMap() http.HandlerFunc {
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

		gameMap, err := h.gameMapRepo.GetByID(mapID)
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

		if err := h.gameMapRepo.Update(gameMap); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to update map")
			return
		}

		h.hub.BroadcastMapUpdate(mapID, "MAP_UPDATED", gameMap)

		respondJSON(w, http.StatusOK, gameMap)
	}
}

func (h *gameMapHandler) deleteMap() http.HandlerFunc {
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

		gameMap, err := h.gameMapRepo.GetByID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if gameMap.OwnerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can delete this map")
			return
		}

		if err := h.gameMapRepo.Delete(mapID); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to delete map")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Map deleted successfully"})
	}
}
