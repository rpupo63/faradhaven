package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rs/zerolog/log"
)

type userHandler struct {
	userRepo database.UserRepository
}

func newUserHandler(userRepo database.UserRepository) *userHandler {
	return &userHandler{userRepo: userRepo}
}

// getAllUsers returns all users
func (h *userHandler) getAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.userRepo.FindAll()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get users")
			respondError(w, http.StatusInternalServerError, "Failed to get users")
			return
		}
		respondJSON(w, http.StatusOK, users)
	}
}

// getUser returns a user by ID
func (h *userHandler) getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "userID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != idStr {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		user, err := h.userRepo.FindByID(id)
		if err != nil {
			log.Error().Err(err).Str("userID", idStr).Msg("Failed to get user")
			respondError(w, http.StatusNotFound, "User not found")
			return
		}
		respondJSON(w, http.StatusOK, user)
	}
}

// getUserFull returns a user with all related data (characters, spells, beasts)
func (h *userHandler) getUserFull() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "userID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != idStr {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		user, err := h.userRepo.FindByIDWithAllRelations(id)
		if err != nil {
			log.Error().Err(err).Str("userID", idStr).Msg("Failed to get user with relations")
			respondError(w, http.StatusNotFound, "User not found")
			return
		}
		respondJSON(w, http.StatusOK, user)
	}
}

// createUser creates a new user
func (h *userHandler) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Name == "" || req.Email == "" {
			respondError(w, http.StatusBadRequest, "Name and email are required")
			return
		}

		user := &models.User{
			Name:  req.Name,
			Email: req.Email,
		}

		if err := h.userRepo.Add(user); err != nil {
			log.Error().Err(err).Msg("Failed to create user")
			respondError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		respondJSON(w, http.StatusCreated, user)
	}
}

// updateUser updates an existing user
func (h *userHandler) updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "userID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != idStr {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		user, err := h.userRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}

		var req UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Name != nil {
			user.Name = *req.Name
		}
		if req.Email != nil {
			user.Email = *req.Email
		}
		if req.ActiveCharacterID != nil {
			user.ActiveCharacterID = req.ActiveCharacterID
		}
		if req.DiceTheme != nil {
			user.DiceTheme = *req.DiceTheme
		}
		if req.DiceThemeColor != nil {
			user.DiceThemeColor = *req.DiceThemeColor
		}
		if req.DiceFontColor != nil {
			user.DiceFontColor = *req.DiceFontColor
		}

		if err := h.userRepo.Update(user); err != nil {
			log.Error().Err(err).Msg("Failed to update user")
			respondError(w, http.StatusInternalServerError, "Failed to update user")
			return
		}

		respondJSON(w, http.StatusOK, user)
	}
}

// setActiveCharacter sets the user's active character for shop purchases
func (h *userHandler) setActiveCharacter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "userID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != idStr {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		var req SetActiveCharacterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		user, err := h.userRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}

		// Update active character (can be nil to clear)
		user.ActiveCharacterID = req.CharacterID

		if err := h.userRepo.Update(user); err != nil {
			log.Error().Err(err).Msg("Failed to set active character")
			respondError(w, http.StatusInternalServerError, "Failed to set active character")
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message":             "Active character updated",
			"active_character_id": user.ActiveCharacterID,
		})
	}
}

// deleteUser deletes a user
func (h *userHandler) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "userID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != idStr {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		if err := h.userRepo.Delete(id); err != nil {
			log.Error().Err(err).Msg("Failed to delete user")
			respondError(w, http.StatusInternalServerError, "Failed to delete user")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
	}
}
