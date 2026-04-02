import re

with open("backend/api/auth_handler.go", "r") as f:
    content = f.read()

# 1. Imports
content = content.replace('"gorm.io/gorm"', '"time"\n\t"github.com/golang-jwt/jwt/v5"\n\t"gorm.io/gorm"')

# 2. Add JWT secret and token generation
jwt_funcs = """
var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Warn().Msg("JWT_SECRET not set, using default insecure secret")
		return "super_secret_default_key"
	}
	return secret
}

func generateAccessToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	return token.SignedString(jwtSecret)
}

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
		SameSite: http.SameSiteLaxMode,
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
"""

content = content.replace('// generateSessionToken', jwt_funcs + '\n// generateSessionToken')

# 3. Replace token logic in login
login_old = """		token, err := generateSessionToken()
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("token creation failed"))
			return
		}
		u.Token = &token
		if err := h.userRepo.Update(u); err != nil {
			h.responder.WriteError(w, errs.NewInternalError("could not save token"))
			return
		}

		resp := LoginResponse{Token: token, UserID: &u.ID}"""

login_new = """		accessToken, err := generateAccessToken(u.ID)
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("access token creation failed"))
			return
		}

		refreshToken, err := generateRefreshToken()
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("refresh token creation failed"))
			return
		}

		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		u.RefreshToken = &refreshToken
		u.RefreshTokenExpiresAt = &expiresAt
		if err := h.userRepo.Update(u); err != nil {
			h.responder.WriteError(w, errs.NewInternalError("could not save token"))
			return
		}

		setRefreshTokenCookie(w, refreshToken, expiresAt)

		resp := LoginResponse{Token: accessToken, UserID: &u.ID}"""

content = content.replace(login_old, login_new)

# 4. Replace token logic in register
reg_old = """		token, err := generateSessionToken()
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

		resp := LoginResponse{Token: token, UserID: &newUser.ID}"""

reg_new = """		refreshToken, err := generateRefreshToken()
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

		resp := LoginResponse{Token: accessToken, UserID: &newUser.ID}"""

content = content.replace(reg_old, reg_new)

# 5. Add refresh and logout methods
refresh_logout = """
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

		accessToken, err := generateAccessToken(u.ID)
		if err != nil {
			h.responder.WriteError(w, errs.NewInternalError("access token creation failed"))
			return
		}

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
				u.RefreshToken = nil
				u.RefreshTokenExpiresAt = nil
				_ = h.userRepo.Update(u)
			}
		}

		clearRefreshTokenCookie(w)
		w.WriteHeader(http.StatusOK)
	}
}

"""

# Insert before hashPassword
content = content.replace('// hashPassword', refresh_logout + '// hashPassword')

with open("backend/api/auth_handler.go", "w") as f:
    f.write(content)
