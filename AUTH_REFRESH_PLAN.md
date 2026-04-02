# Access & Refresh Token System Implementation Plan

## Goal
Implement a robust authentication system using stateless JWTs for Access Tokens and database-backed Opaque strings for Refresh Tokens, employing a hybrid storage approach (LocalStorage for Access Tokens, HttpOnly Cookies for Refresh Tokens).

## 1. Backend Changes

### 1.1 Dependency Updates
- Add `github.com/golang-jwt/jwt/v5` to `go.mod`.

### 1.2 Model Updates (`backend/models/user.go`)
- Remove the `Token *string` field.
- Add `RefreshToken *string` (`gorm:"type:text;uniqueIndex"`).
- Add `RefreshTokenExpiresAt *time.Time` (`gorm:"type:timestamptz"`).

### 1.3 Database Repository Updates (`backend/database/user_repo.go`)
- Update `FindByToken` method to `FindByRefreshToken(token string)`.

### 1.4 Auth Handlers & Routes (`backend/api/auth_handler.go`, `backend/api/routes.go`)
- Add `JWT_SECRET` environment variable handling in config/bootstrap.
- Implement `generateAccessToken(user *models.User)` generating a 15-minute JWT.
- Implement `generateRefreshToken()` generating a secure hex string.
- Update `login` and `register` endpoints to:
  - Generate both tokens.
  - Set the Refresh Token and expiration in the database.
  - Set the Refresh Token in an HttpOnly, Secure, SameSite=Lax cookie.
  - Return only the Access Token and User ID in the JSON response.
- Add `/api/auth/refresh` endpoint (POST):
  - Read the refresh token cookie.
  - Validate it against the database and check expiration.
  - Issue a new Access Token and a rotated Refresh Token (save new to DB, set new cookie).
- Add `/api/auth/logout` endpoint (POST):
  - Read refresh token cookie, clear it from the DB.
  - Expire the cookie on the client side.

### 1.5 Middleware Updates (`backend/api/middleware.go`)
- Update `authenticate` middleware to parse the Bearer JWT token from the `Authorization` header instead of hitting the database.
- Validate JWT signature and expiration.
- Extract `user_id` and attach it to the request context.

## 2. Frontend Changes

### 2.1 API Interceptor & Fetch Wrapper (`frontend/src/lib/api/base.ts`)
- Update `makeApiRequest` to include `credentials: 'include'` so cookies are sent to the backend.
- Implement a response interceptor logic for 401s:
  - If a 401 response is received, check if we are already refreshing.
  - If not, set a lock, call `/api/auth/refresh` to get a new access token, update LocalStorage, and retry the failed request.
  - If a refresh is in progress, queue the request until the refresh completes.
  - If the refresh fails, clear auth state and redirect to login.

### 2.2 Auth API Updates (`frontend/src/lib/api/auth.ts`)
- Add `refreshAuthToken` API function.
- Add `logout` API function to hit the backend logout endpoint.

### 2.3 Auth Context (`frontend/src/context/AuthContext.tsx`)
- Update `logout` to call the backend logout endpoint so the refresh cookie is cleared.
- Ensure the app attempts to use the token natively; if it fails, the new interceptor will handle the refresh seamlessly.

## Validation Strategy
- Verify a user can log in and register, receiving an access token and an HttpOnly cookie.
- Verify that authenticated API requests use the JWT successfully.
- Verify that waiting 15 minutes (or artificially expiring the JWT) triggers the frontend to automatically request a new access token via the `/refresh` endpoint and retries the original request.
- Verify logging out clears the local storage and invalidates the refresh token in the backend and removes the cookie.
