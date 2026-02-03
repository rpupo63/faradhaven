package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/errs"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authHandler struct {
	userRepo  database.UserRepository
	responder Responder
}

func newAuthHandler(userRepo database.UserRepository) *authHandler {
	return &authHandler{
		userRepo:  userRepo,
		responder: NewResponder(log.With().Str("handler", "auth").Logger()),
	}
}

// LoginRequest is the body for POST /api/auth/login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is returned on successful login
type LoginResponse struct {
	Token  string     `json:"token"`
	UserID *uuid.UUID `json:"user_id,omitempty"` // optional first resume user for convenience
}

// RegisterRequest is the body for POST /api/auth/register
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *authHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			h.responder.WriteError(w, errs.BadRequest("method not allowed"))
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.responder.WriteError(w, errs.Malformed("body"))
			return
		}
		if req.Email == "" || req.Password == "" {
			h.responder.WriteError(w, errs.NewMissingRequiredFieldError("email and password"))
			return
		}

		u, err := h.userRepo.FindByEmail(req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				h.responder.WriteError(w, errs.Unauthorized)
				return
			}
			h.responder.WriteError(w, errs.NewInternalError("auth lookup failed"))
			return
		}
		if u.PasswordHash == "" {
			h.responder.WriteError(w, errs.Unauthorized)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
			h.responder.WriteError(w, errs.Unauthorized)
			return
		}

		token, err := generateSessionToken()
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("token creation failed"))
			return
		}
		u.Token = &token
		if err := h.userRepo.Update(u); err != nil {
			h.responder.WriteError(w, errs.NewInternalError("could not save token"))
			return
		}

		resp := LoginResponse{Token: token, UserID: &u.ID}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func (h *authHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.responder.WriteError(w, errs.Malformed("body"))
			return
		}

		if req.Email == "" || req.Password == "" || req.Name == "" {
			h.responder.WriteError(w, errs.NewMissingRequiredFieldError("name, email, and password"))
			return
		}

		// Check if user exists
		existing, err := h.userRepo.FindByEmail(req.Email)
		if err == nil && existing != nil {
			h.responder.WriteError(w, errs.NewConflictError("email already in use"))
			return
		}

		hash, err := hashPassword(req.Password)
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("password hashing failed"))
			return
		}

		token, err := generateSessionToken()
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("token generation failed"))
			return
		}

		newUser := &models.User{
			Name:         req.Name,
			Email:        req.Email,
			PasswordHash: hash,
			Token:        &token,
		}

		if err := h.userRepo.Add(newUser); err != nil {
			h.responder.WriteError(w, errs.NewInternalError("failed to create user"))
			return
		}

		resp := LoginResponse{Token: token, UserID: &newUser.ID}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// generateSessionToken returns a cryptographically random token for API auth (Bearer).
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashPassword hashes a password with bcrypt for storage.
func hashPassword(password string) (string, error) {
	const cost = bcrypt.DefaultCost
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// EnsureFirstUserExists creates the first user from AUTH_EMAIL and AUTH_PASSWORD if the users table is empty.
// Call once after DB migration (e.g. from main after AutoMigrate).
func EnsureFirstUserExists(userRepo database.UserRepository) {
	users, err := userRepo.FindAll()
	if err != nil || len(users) > 0 {
		return
	}
	email := os.Getenv("AUTH_EMAIL")
	password := os.Getenv("AUTH_PASSWORD")
	if email == "" || password == "" {
		return
	}
	hash, err := hashPassword(password)
	if err != nil {
		log.Warn().Err(err).Msg("could not hash AUTH_PASSWORD, skipping bootstrap")
		return
	}
	name := os.Getenv("AUTH_NAME")
	if name == "" {
		name = "User"
	}
	u := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	}
	if err := userRepo.Add(u); err != nil {
		log.Warn().Err(err).Msg("could not create initial user")
		return
	}
	log.Info().Str("email", email).Msg("created initial user from AUTH_EMAIL/AUTH_PASSWORD")
}
