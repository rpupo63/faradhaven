package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Warn().Msg("JWT_SECRET not set, using default insecure secret")
		return "super_secret_default_key" // In production, this should fail or generate a random key
	}
	return secret
}

// generateAccessToken creates a short-lived JWT token
func generateAccessToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	return token.SignedString(jwtSecret)
}

// generateRefreshToken returns a cryptographically random token for API auth
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func setRefreshTokenCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // Adjust if cross-origin needs change
	})
}

func clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
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

		accessToken, err := generateAccessToken(u.ID)
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("access token creation failed"))
			return
		}

		refreshToken, err := generateRefreshToken()
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("refresh token creation failed"))
			return
		}

		expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days
		u.RefreshToken = &refreshToken
		u.RefreshTokenExpiresAt = &expiresAt
		if err := h.userRepo.Update(u); err != nil {
			h.responder.WriteError(w, errs.NewInternalError("could not save token"))
			return
		}

		setRefreshTokenCookie(w, refreshToken, expiresAt)

		resp := LoginResponse{Token: accessToken, UserID: &u.ID}

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

		refreshToken, err := generateRefreshToken()
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("token generation failed"))
			return
		}
		
		expiresAt := time.Now().Add(7 * 24 * time.Hour)

		newUser := &models.User{
			Name:                  req.Name,
			Email:                 req.Email,
			PasswordHash:          hash,
			RefreshToken:          &refreshToken,
			RefreshTokenExpiresAt: &expiresAt,
		}

		if err := h.userRepo.Add(newUser); err != nil {
			h.responder.WriteError(w, errs.NewInternalError("failed to create user"))
			return
		}
		
		accessToken, err := generateAccessToken(newUser.ID)
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("access token generation failed"))
			return
		}

		setRefreshTokenCookie(w, refreshToken, expiresAt)

		resp := LoginResponse{Token: accessToken, UserID: &newUser.ID}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func (h *authHandler) refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			h.responder.WriteError(w, errs.Unauthorized)
			return
		}

		tokenString := cookie.Value
		u, err := h.userRepo.FindByRefreshToken(tokenString)
		if err != nil || u == nil {
			h.responder.WriteError(w, errs.Unauthorized)
			return
		}

		if u.RefreshTokenExpiresAt == nil || time.Now().After(*u.RefreshTokenExpiresAt) {
			h.responder.WriteError(w, errs.Unauthorized)
			return
		}

		// Generate new access token
		accessToken, err := generateAccessToken(u.ID)
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("access token creation failed"))
			return
		}

		// Rotate refresh token
		newRefreshToken, err := generateRefreshToken()
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("refresh token creation failed"))
			return
		}

		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		u.RefreshToken = &newRefreshToken
		u.RefreshTokenExpiresAt = &expiresAt
		if err := h.userRepo.Update(u); err != nil {
			h.responder.WriteError(w, errs.NewInternalError("could not save token"))
			return
		}

		setRefreshTokenCookie(w, newRefreshToken, expiresAt)

		resp := LoginResponse{Token: accessToken, UserID: &u.ID}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func (h *authHandler) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err == nil && cookie.Value != "" {
			u, err := h.userRepo.FindByRefreshToken(cookie.Value)
			if err == nil && u != nil {
				// Clear the refresh token in the DB
				u.RefreshToken = nil
				u.RefreshTokenExpiresAt = nil
				_ = h.userRepo.Update(u)
			}
		}

		clearRefreshTokenCookie(w)
		w.WriteHeader(http.StatusOK)
	}
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
