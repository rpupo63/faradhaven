# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Faradhaven is a D&D 5e-inspired fantasy RPG with a Go backend and React TypeScript frontend. The system features custom class/race/item seeding with deterministic UUIDs and a component-based spell crafting system.

## Commands

### Backend (from `backend/`)
```bash
docker compose up -d              # Start PostgreSQL (pgvector) on port 5432
go run main.go                    # Run server (port 8080)
go build ./...                    # Build all packages
GENERATE_MODELS=true go run main.go  # Run migrations then start server
GENERATE_MODELS=only go run main.go  # Run migrations only, then exit
```

### Frontend (from `frontend/`)
```bash
npm run dev          # Vite dev server (port 8081, proxies to backend)
npm run build        # Production build
npm run lint         # ESLint
npm run test         # Run tests (vitest)
npm run test:watch   # Watch mode
```

## Architecture

### Backend
- **Router**: go-chi/chi
- **ORM**: GORM with PostgreSQL (pgvector for AI embeddings)
- **Pattern**: Handler → Repository → Service

Key directories:
- `api/` - HTTP handlers and routes
- `database/` - Repository interfaces and implementations
- `models/` - GORM entities
- `services/` - Business logic (HP, leveling)
- `seed/` - Data seeding system

### Frontend
- **Build**: Vite + React 18 + TypeScript
- **UI**: shadcn/ui (Radix) + Tailwind CSS
- **Data**: TanStack Query
- **Forms**: React Hook Form + Zod

### Seed System

All game data (classes, races, items, components) uses a custom seeding system:

1. **Deterministic UUIDs** via `seed/uuids/uuids.go`:
   ```go
   uuids.ClassUUID("The Mutagen")  // Stable across reseeds
   uuids.RaceUUID("Human")
   ```

2. **Batch operations** via `seed/batch/upsert.go`:
   - `UpsertBatchUpdateAll` for parent tables (ON CONFLICT DO UPDATE)
   - `InsertBatch` for child tables (cleared before insert)

3. **Versioning** via `seed/versioning/version.go`:
   - SHA256 hash comparison skips unchanged seeds
   - Tracked in `models.SeedMetadata`

4. **Execution order** (alphabetical by prefix):
   - `01_components` → `02_items` → `03_races` → `04_classes`

5. **Transaction safety**: Each seed runs in its own `db.Transaction()`

### Entity Relationships

**Classes**: Class → Archetype, ClassLevel → LevelFeature, Class ↔ Component (many-to-many)

**Races**: Race → Trait → TraitOption, Race ↔ Component (many-to-many), Race → Lineage (subraces)

**Items**: Weapon → WeaponDamage, Item (standalone)

### Gotchas

- Join tables use `string` UUIDs, not `uuid.UUID`
- Clear child tables before parent upserts due to foreign keys
- LevelFeature needs `ArchetypeID` pointer (nil = shared by all archetypes)
- Component-based spell system: spells are composed of Components (Lux, Umbra, Nova, etc.)

## Environment

Copy `backend/.env.example` to `backend/.env`. Key variables:
- `DATABASE_URL` - PostgreSQL connection string
- `GENERATE_MODELS` - Set to `true` or `only` for migrations
- `PORT` - Backend port (default 8080)
- `ACCEPTED_ORIGINS` - CORS origins
