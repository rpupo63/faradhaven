# Faradhaven

Faradhaven is a full-stack, tabletop-inspired RPG platform set in a grimy magitech world called **South Axiom**.  
Players create and manage characters, cast and synthesize spells, track resources and effects, run encounters, and collaborate through shared campaign tools.

The game fiction centers on:
- **The Vitalic Age** and industrialized arcana
- Corporate factions like the **DCC**, **ACC**, and **ECC**
- Constructed "Thesis" beings struggling with serum dependence, decay, and identity
- A city split between elite aerostat estates and dangerous undercity zones

The landing narrative and tone are driven by the home page at `frontend/src/pages/Index.tsx`.

## Tech Stack

- **Frontend:** React 18 + TypeScript + Vite + Tailwind + shadcn/ui
- **Backend:** Go (chi router + GORM)
- **Database:** PostgreSQL with pgvector extension support
- **State/data fetching:** TanStack Query
- **Auth model:** Token-based (Bearer token validated against user token in DB)
- **Realtime:** WebSocket hub for battle-map collaboration

## Project Structure

```text
faradhaven/
  backend/    # Go API server, data models, repositories, services, seeding, migrations
  frontend/   # React app, routes/pages, UI components, API client integration
```

## How The System Operates

### 1) Frontend app lifecycle

- `frontend/src/App.tsx` boots providers in this order:
  - `QueryClientProvider` (data fetching/cache)
  - `AuthProvider` (token, user, auth lifecycle)
  - `GameProvider`
  - `BrowserRouter` (page routing)
- Public and protected pages are enforced with route wrappers.
- The landing page (`/`) is public and introduces setting/lore.

### 2) Authentication flow

- Login/register calls backend `/api/auth/login` and `/api/auth/register`.
- `AuthContext` stores token and user id in localStorage.
- On app load, token validity is checked by fetching the user profile.
- Authenticated API requests send `Authorization: Bearer <token>`.
- Backend middleware validates the token per request.

### 3) API and domain flow

Backend follows a layered flow:

`HTTP Handler -> Repository -> Service -> Database`

Key domain areas exposed by routes include:
- Character creation, sheets, leveling, rest/hp/hit-dice, resources, effects
- Spell compendium and player spellbook operations (including synthesis/AI helpers)
- Beasts, monsters, corpses, harvesting/scavenging loops
- Party management and campaign utility endpoints
- Shared notes/bulletin features
- Battle maps, map tokens/elements, and initiative tracking
- WebSocket map room updates

### 4) Startup, migration, and seed behavior

`backend/main.go` controls startup:

- Loads `.env`
- Connects database
- Optional migration+seed behavior via `GENERATE_MODELS`:
  - `only`: run migrations/seeding and exit
  - `true`: run migrations/seeding, then start server
  - `false`/unset: skip migrations/seeding and start server
- Starts HTTP server (default port `8080`)

The seed system uses deterministic UUIDs and ordered seed stages so core game data remains stable across reseeds.

### 5) Frontend-backend communication

- Frontend dev server runs on `8081`
- Vite proxies `/api/*` to backend `http://localhost:8080`
- This keeps browser requests same-origin in development from the app's perspective

## Local Development

### Prerequisites

- Go `1.24+`
- Node.js `18+` (or current LTS)
- Docker (for local Postgres)

## 1. Start database

From `backend/`:

```bash
docker compose up -d
```

This starts a local Postgres container on port `5432`.

## 2. Configure environment

```bash
cp backend/.env.example backend/.env
```

At minimum, set:
- `DATABASE_URL`
- `PORT` (default `8080`)
- `ACCEPTED_ORIGINS`

Optional but commonly used:
- `AUTH_EMAIL`, `AUTH_PASSWORD`, `AUTH_NAME` (bootstrap first user)
- `GOOGLE_API_KEY` (AI-assisted spell features)
- S3 bucket variables for uploads

## 3. Run backend

From `backend/`:

```bash
go run main.go
```

Useful variants:

```bash
# Run migrations/seeding then start server
GENERATE_MODELS=true go run main.go

# Run migrations/seeding only, then exit
GENERATE_MODELS=only go run main.go
```

## 4. Run frontend

From `frontend/`:

```bash
npm install
npm run dev
```

Frontend will be available on `http://localhost:8081`.

## Common Commands

### Backend (`backend/`)

```bash
go build ./...
```

### Frontend (`frontend/`)

```bash
npm run lint
npm run test
npm run build
```

## Health and Smoke Checks

- Backend root: `http://localhost:8080/`
- Health check: `http://localhost:8080/healthcheck`
- Frontend app: `http://localhost:8081`

Quick API check after login:
- `GET /api/classes`
- `GET /api/races`
- Authenticated: `GET /api/characters`

## Notes For Contributors

- This repository currently has active in-progress changes across backend files; coordinate before broad refactors.
- Prefer additive changes around handlers/services/models and keep route contracts stable.
- For seed data, preserve deterministic ID patterns and seed ordering conventions.
