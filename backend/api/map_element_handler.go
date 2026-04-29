package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
)

type mapElementHandler struct {
	elementRepo *database.MapElementRepo
	gameMapRepo *database.GameMapRepo
}

func newMapElementHandler(elementRepo *database.MapElementRepo, gameMapRepo *database.GameMapRepo) *mapElementHandler {
	return &mapElementHandler{elementRepo: elementRepo, gameMapRepo: gameMapRepo}
}

func (h *mapElementHandler) createMapElement() http.HandlerFunc {
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
			respondError(w, http.StatusForbidden, "Only the DM can add map elements")
			return
		}

		var req CreateMapElementRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		element := models.MapElement{
			MapID: mapID,
			Type:  req.Type,
			GridX: req.GridX,
			GridY: req.GridY,
		}

		switch req.Type {
		case models.MapElementTypeTrap:
			element.TrapProperties = &models.TrapProperties{
				Name:             derefString(req.TrapName),
				Description:      derefString(req.TrapDescription),
				VisibleToPlayers: derefBool(req.TrapVisibleToPlayers),
			}
		case models.MapElementTypeDifficultTerrain:
			element.DifficultTerrainProperties = &models.DifficultTerrainProperties{
				Color: derefString(req.DifficultTerrainColor),
			}
		case models.MapElementTypeElevation:
			element.ElevationProperties = &models.ElevationProperties{
				Level: derefInt(req.ElevationLevel),
			}
		case models.MapElementTypeWall:
			element.WallProperties = &models.WallProperties{
				Details: derefString(req.WallDetails),
			}
		}

		if err := h.elementRepo.CreateWithProperties(&element); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to create map element")
			return
		}

		respondJSON(w, http.StatusCreated, element)
	}
}

func (h *mapElementHandler) createMultipleMapElements() http.HandlerFunc {
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
			respondError(w, http.StatusForbidden, "Only the DM can add map elements")
			return
		}

		var reqs CreateMultipleMapElementsRequest
		if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		elements := make([]*models.MapElement, 0, len(reqs.Elements))
		for _, req := range reqs.Elements {
			element := &models.MapElement{
				MapID: mapID,
				Type:  req.Type,
				GridX: req.GridX,
				GridY: req.GridY,
			}
			switch req.Type {
			case models.MapElementTypeTrap:
				element.TrapProperties = &models.TrapProperties{
					Name:             derefString(req.TrapName),
					Description:      derefString(req.TrapDescription),
					VisibleToPlayers: derefBool(req.TrapVisibleToPlayers),
				}
			case models.MapElementTypeDifficultTerrain:
				element.DifficultTerrainProperties = &models.DifficultTerrainProperties{
					Color: derefString(req.DifficultTerrainColor),
				}
			case models.MapElementTypeElevation:
				element.ElevationProperties = &models.ElevationProperties{
					Level: derefInt(req.ElevationLevel),
				}
			case models.MapElementTypeWall:
				element.WallProperties = &models.WallProperties{
					Details: derefString(req.WallDetails),
				}
			}
			elements = append(elements, element)
		}

		if err := h.elementRepo.CreateMultipleWithProperties(elements); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to create map elements")
			return
		}

		createdElements := make([]models.MapElement, len(elements))
		for i, e := range elements {
			createdElements[i] = *e
		}
		respondJSON(w, http.StatusCreated, createdElements)
	}
}

func (h *mapElementHandler) updateMapElement() http.HandlerFunc {
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

		elementID, err := uuid.Parse(chi.URLParam(r, "elementID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid element ID")
			return
		}

		ownerID, err := h.gameMapRepo.GetOwnerID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if ownerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can update map elements")
			return
		}

		element, err := h.elementRepo.GetByID(elementID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map element not found")
			return
		}

		var req UpdateMapElementRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if req.Type != nil {
			element.Type = *req.Type
		}
		if req.GridX != nil {
			element.GridX = *req.GridX
		}
		if req.GridY != nil {
			element.GridY = *req.GridY
		}

		if err := h.elementRepo.Update(element); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to update map element")
			return
		}

		respondJSON(w, http.StatusOK, element)
	}
}

func (h *mapElementHandler) deleteMapElement() http.HandlerFunc {
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

		elementID, err := uuid.Parse(chi.URLParam(r, "elementID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid element ID")
			return
		}

		ownerID, err := h.gameMapRepo.GetOwnerID(mapID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Map not found")
			return
		}

		if ownerID != userID {
			respondError(w, http.StatusForbidden, "Only the DM can delete map elements")
			return
		}

		if err := h.elementRepo.Delete(elementID); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to delete map element")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Map element deleted successfully"})
	}
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

func derefInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}
