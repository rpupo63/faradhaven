package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"github.com/rs/zerolog/log"
)

type minionHandler struct {
	minionService *services.MinionService
}

func newMinionHandler(minionService *services.MinionService) *minionHandler {
	return &minionHandler{minionService: minionService}
}

// GetMinions returns all minions for a character
func (h *minionHandler) GetMinions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		// Optional type filter
		var minionType *models.MinionType
		typeStr := r.URL.Query().Get("type")
		if typeStr != "" {
			mt := models.MinionType(typeStr)
			minionType = &mt
		}

		minions, err := h.minionService.GetMinions(characterID, minionType)
		if err != nil {
			log.Error().Err(err).Str("characterID", characterID.String()).Msg("Failed to get minions")
			respondError(w, http.StatusInternalServerError, "Failed to get minions")
			return
		}

		respondJSON(w, http.StatusOK, minions)
	}
}

// GetMinion returns a single minion by ID
func (h *minionHandler) GetMinion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minionID, err := uuid.Parse(chi.URLParam(r, "minionID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid minion ID")
			return
		}

		minion, err := h.minionService.GetMinion(minionID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Minion not found")
			return
		}

		respondJSON(w, http.StatusOK, minion)
	}
}

// CreateMinion creates a new minion (construct, echo, or drone)
func (h *minionHandler) CreateMinion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		var req CreateMinionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		var minion *models.Minion

		switch models.MinionType(req.MinionType) {
		case models.MinionTypeConstruct:
			constructReq := services.CreateConstructRequest{
				TemplateKey: req.TemplateKey,
				CustomName:  req.CustomName,
			}
			minion, err = h.minionService.CreateConstruct(characterID, constructReq)
		case models.MinionTypeDrone:
			if req.ComponentID == nil {
				respondError(w, http.StatusBadRequest, "component_id is required for drone creation")
				return
			}
			droneReq := services.CreateDroneRequest{
				ComponentID: *req.ComponentID,
				CustomName:  req.CustomName,
			}
			minion, err = h.minionService.CreateDrone(characterID, droneReq)
		case models.MinionTypeEcho:
			echoReq := services.StoreEchoRequest{
				SlotIndex:          req.SlotIndex,
				AbilityName:        req.AbilityName,
				AbilityDescription: req.AbilityDescription,
				AbilityType:        req.AbilityType,
				SourceCreatureType: req.SourceCreatureType,
			}
			minion, err = h.minionService.StoreEcho(characterID, echoReq)
		default:
			respondError(w, http.StatusBadRequest, "Invalid minion type")
			return
		}

		if err != nil {
			log.Error().Err(err).Msg("Failed to create minion")
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondJSON(w, http.StatusCreated, minion)
	}
}

// UpdateMinionHP updates a minion's HP
func (h *minionHandler) UpdateMinionHP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minionID, err := uuid.Parse(chi.URLParam(r, "minionID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid minion ID")
			return
		}

		var req UpdateMinionHPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		minion, err := h.minionService.UpdateMinionHP(minionID, req.Delta)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update minion HP")
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, minion)
	}
}

// ActivateMinion activates an echo
func (h *minionHandler) ActivateMinion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minionID, err := uuid.Parse(chi.URLParam(r, "minionID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid minion ID")
			return
		}

		minion, err := h.minionService.ActivateEcho(minionID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to activate minion")
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, minion)
	}
}

// DeactivateMinion deactivates an echo
func (h *minionHandler) DeactivateMinion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minionID, err := uuid.Parse(chi.URLParam(r, "minionID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid minion ID")
			return
		}

		if err := h.minionService.DeactivateEcho(minionID); err != nil {
			log.Error().Err(err).Msg("Failed to deactivate minion")
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"status": "deactivated"})
	}
}

// DeleteMinion removes a minion
func (h *minionHandler) DeleteMinion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minionID, err := uuid.Parse(chi.URLParam(r, "minionID"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid minion ID")
			return
		}

		// Check for reclaim parameter (for constructs/drones)
		reclaim := r.URL.Query().Get("reclaim") == "true"

		minion, err := h.minionService.GetMinion(minionID)
		if err != nil {
			respondError(w, http.StatusNotFound, "Minion not found")
			return
		}

		if (minion.MinionType == models.MinionTypeConstruct || minion.MinionType == models.MinionTypeDrone) && reclaim {
			reclaimed, err := h.minionService.DestroyConstruct(minionID, true)
			if err != nil {
				log.Error().Err(err).Msg("Failed to destroy construct")
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"status":               "destroyed",
				"components_reclaimed": reclaimed,
			})
			return
		}

		if err := h.minionService.DeleteMinion(minionID); err != nil {
			log.Error().Err(err).Msg("Failed to delete minion")
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// GetConstructTemplates returns available construct templates
func (h *minionHandler) GetConstructTemplates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates := models.GetConstructTemplates()
		respondJSON(w, http.StatusOK, templates)
	}
}
