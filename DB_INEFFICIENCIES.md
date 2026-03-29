# Database Inefficiency Report — Faradhaven Backend

_Updated: 2026-03-27 (original: 2026-03-22)_

Items marked ✅ have been resolved since the original report. Items marked 🆕 are new findings.

---

## Status of Original Issues

| # | Issue | Status |
|---|-------|--------|
| 1 | Hard-coded 10k spell limit | ✅ Fixed — `FindAllWithComponents()` used |
| 2 | N+1 component lookups in corpse handler | ✅ Fixed — batch `GetComponentsByIDs()` |
| 3 | N+1 map element writes | ✅ Partially — `CreateMultipleWithProperties()` added; token initiative loop remains |
| 4 | 6 queries for character sheet | ✅ Fixed — `FindByIDForSheet()` consolidates preloads |
| 5 | 3 queries for spellbook | ✅ Fixed — `FindByIDWithRelations()` used |
| 6 | `buildClassResources` duplicated | ✅ Fixed — moved to shared `api/class_resources.go` |
| 7 | Redundant party queries | 🔴 Open |
| 8 | Repeated resource key lookups | 🔴 Open |
| 9 | Raw DB access bypassing repos | 🔴 Open |
| 10 | Auth checks load full entities | 🔴 Open |
| 11 | `GetIdentifiedBeasts` loads entire party | 🔴 Open |
| 12 | Missing pagination on list endpoints | 🔴 Open |
| 13 | Missing transactions on multi-step writes | 🔴 Open |
| 14 | Manual skill ID extraction (second query) | 🔴 Open |
| 15 | Missing indexes | 🟡 Partial — GORM tags present, but compound index on `(character_id, resource_key)` still absent |

---

## Open Issues

---

### 1. Double Character Load on Every Sheet Request — **HIGH**

**File:** `backend/api/character_sheet_handler.go:26–56`

`getCharacterSheet()` calls `FindByID()` (with Race/Class/Archetype/Components preloads) purely for the ownership check, then immediately calls `getCharacterSheetData()` which calls `FindByIDForSheet()` — a second, heavier load of the same character. Every sheet view fires two full character queries.

```go
// Line 26: auth load
character, err := h.characterRepo.FindByID(id)
// Line 56: full sheet load (FindByIDForSheet internally)
sheet, err := h.getCharacterSheetData(id)
```

**Fix:** Call `FindByIDForSheet()` once, check `character.UserID` against the auth token, then proceed. No second query needed.

---

### 2. Party Data: 4 Round-Trips Per Character Sheet — **HIGH**

**File:** `backend/api/character_sheet_handler.go:62–81`, `backend/database/party_repo.go:142–148`

When a character belongs to a party, the sheet handler issues three calls:
1. `partyRepo.FindByID()` — loads party + owner
2. `partyRepo.GetMembers()` — `SELECT … WHERE party_id = ?`
3. `partyRepo.GetIdentifiedBeasts()` — internally calls `First(party)` **again** then returns `party.IdentifiedBeasts`

`GetIdentifiedBeasts()` re-fetches the entire `Party` row a second time internally:

```go
// party_repo.go:144
r.db.Preload("IdentifiedBeasts").First(&party, "id = ?", partyID)
return party.IdentifiedBeasts, nil
```

That's 4 DB hits (2 party row fetches + 1 members query + 1 join table query) to populate party data.

**Fix:** Add `FindByIDWithMembersAndBeasts()` that preloads both in one call. Replace the three-call pattern in the handler with the single preloaded result.

---

### 3. Map Token/Element Operations Load the Full Map — **HIGH**

**Files:** `backend/api/map_token_handler.go`, `backend/api/map_element_handler.go`

Every token write (add, update, delete) and every element write calls `h.gameMapRepo.GetByID(mapID)` for an ownership check. `GetByID` preloads **all tokens** and **all elements with their four property tables** (Trap, DifficultTerrain, Elevation, Wall):

```go
// game_map_repo.go:23–33
r.db.Preload("Tokens").
    Preload("Elements").
    Preload("Elements.TrapProperties").
    Preload("Elements.DifficultTerrainProperties").
    Preload("Elements.ElevationProperties").
    Preload("Elements.WallProperties").
    First(&gameMap, "id = ?", id)
```

For a map with 50 elements this fires 6 queries and returns megabytes of data just to verify `gameMap.OwnerID == userID`.

**Fix:** Add `GetOwnerID(mapID uuid.UUID) (uuid.UUID, error)` that runs `SELECT owner_id FROM game_maps WHERE id = ? LIMIT 1`. Use it for auth checks; keep the heavy `GetByID` only for the map-state read endpoint.

---

### 4. `BulkSetInitiativeOrder` Issues N UPDATEs — **HIGH**

**File:** `backend/database/map_token_repo.go:56–70`

Wrapped in a transaction but still fires one `UPDATE` per token:

```go
for _, e := range entries {
    tx.Model(&models.MapToken{}).Where("id = ?", e.TokenID).
        Update("initiative_order", &order)
}
```

For 10 tokens in initiative this is 10 round-trips inside the transaction.

**Fix:** Use a single `UPDATE … CASE WHEN` statement or PostgreSQL's `unnest` approach:
```sql
UPDATE map_tokens SET initiative_order = v.order
FROM (SELECT unnest($1::uuid[]) AS id, unnest($2::int[]) AS order) v
WHERE map_tokens.id = v.id
```

---

### 5. `ReplaceSkillProficiencies` N+1 Inserts — **MEDIUM**

**File:** `backend/database/character_repo.go:192–200`

Inserts one `CharacterSkill` per skill in a loop:

```go
for _, skillID := range skillIDs {
    tx.Create(&cs)
}
```

**Fix:** Build a slice and use `tx.CreateInBatches(skills, len(skills))` or a single `INSERT … VALUES (…),(…)`.

---

### 6. `replaceComponentsTx` N+1 Inserts — **MEDIUM**

**File:** `backend/database/spell_repo.go:135–142`

Same pattern — inserts one `SpellComponent` per component:

```go
for _, compID := range componentIDs {
    tx.Create(&sc)
}
```

**Fix:** Batch insert with `tx.CreateInBatches`.

---

### 7. `ProcessRestoration` N+1 Saves — **MEDIUM**

**File:** `backend/database/character_resource_repo.go:103–129`

Fetches all character resources (one query), then saves each one individually inside a loop:

```go
for _, res := range resources {
    // ... mutate ...
    tx.Save(res)  // One UPDATE per resource
}
```

For a class with 6 resources this is 7 queries (1 SELECT + 6 UPDATEs) on every rest.

**Fix:** Use `tx.Save` with a batch or issue a single `UPDATE … SET current_value = CASE …` statement, or use `tx.CreateInBatches`.

---

### 8. `UpdateCurrentValue` Fetches Outside Transaction — **MEDIUM**

**File:** `backend/database/character_resource_repo.go:74–91`

Inside a `r.db.Transaction(func(tx *gorm.DB) error {…})`, the fetch uses the outer receiver `r` instead of the transaction `tx`:

```go
r.db.Transaction(func(tx *gorm.DB) error {
    resource, err := r.FindByCharacterAndKey(...)  // uses r.db, not tx!
    // ...
    return tx.Save(resource).Error
})
```

The read and write are on different connections — not actually atomic. A concurrent update between the read and the save is possible.

**Fix:** Inline the fetch with `tx.Where(...).First(...)` inside the transaction closure.

---

### 9. Auth Checks Load Full Character Entity — **MEDIUM**

**File:** `backend/api/character_handler.go` (multiple handlers)

`updateCharacter`, `deleteCharacter`, `uploadProfilePicture`, `updateBackstory`, `restSpellPoints`, and others call `FindByID()` solely to read `character.UserID` for an ownership check. This loads Race, Class, Archetype, and Components needlessly.

**Fix:** Add `CharacterBelongsToUser(characterID, userID uuid.UUID) (bool, error)` running:
```sql
SELECT 1 FROM characters WHERE id = $1 AND user_id = $2 LIMIT 1
```

---

### 10. `AddIdentifiedBeast` — Two Wasted Lookups — **MEDIUM**

**File:** `backend/database/party_repo.go:109–118`

Fetches the full `Party` row and the full `Beast` row before calling `Association("IdentifiedBeasts").Append()`, even though only the IDs are needed:

```go
r.db.First(&party, "id = ?", partyID)
r.db.First(&beast, "id = ?", beastID)
r.db.Model(&party).Association("IdentifiedBeasts").Append(&beast)
```

**Fix:** Use a direct join-table insert:
```go
r.db.Exec("INSERT INTO party_beasts (party_id, beast_id) VALUES (?, ?) ON CONFLICT DO NOTHING", partyID, beastID)
```

---

### 11. `FindByUserID` for Spells Has No Pagination — **MEDIUM**

**File:** `backend/database/spell_repo.go:73–77`

`FindByUserID` returns every spell for a user with no limit. Over time a prolific spellcaster will see degrading response times.

**Fix:** Add `page`/`limit` parameters (same pattern as `FindAll`).

---

### 12. `FindByUserID` / `FindAll` for Characters Has No Pagination — **MEDIUM**

**File:** `backend/database/character_repo.go:41–44, 148–152`

`FindAll()` and `FindByUserID()` each load every character with full Race/Class/Archetype/Components preloads. No `LIMIT`.

**Fix:** Add pagination. For the character list (dashboard), most clients only need `id`, `name`, `level`, `race.name`, `class.name` — consider a dedicated lightweight list query.

---

### 13. `buildClassResources` Issues 3 Queries Per Rest — **MEDIUM**

**File:** `backend/api/class_resources.go`, called from `shortRest` and `longRest` in `level_handler.go`

Each rest response calls `buildClassResources()` which fires:
1. `FindResourceDefinitionsByClassID`
2. `GetLevelResourceMap` → `FindLevelResourcesByClassAndLevel` (raw SQL JOIN)
3. `FindByCharacterID` on character_resources

These could be combined into a single JOIN query or a single service method that batches all three.

---

### 14. `FindByIDWithSkills` Manual Second Query — **LOW**

**File:** `backend/database/character_repo.go:64–73`

After preloading the character, a second query fetches skills:

```go
var skills []models.CharacterSkill
r.db.Where("character_id = ?", id).Find(&skills)
```

This pattern is repeated in `FindByIDWithRelations`, `FindByIDForSheet` — three separate places.

**Fix:** Add `Preload("SkillProficiencies")` to the main query and derive `SkillProficiencyIDs` in a `BeforeSend` transform, eliminating the secondary query.

---

### 15. Raw DB Access Bypassing Repository Layer — **LOW**

**File:** `backend/api/character_handler.go`, `backend/api/character_creation_handler.go`

```go
h.characterRepo.GetDB().Where("id IN ?", weaponIDs).Find(&weapons)
h.characterRepo.GetDB().Model(character).Association("Weapons").Append(weapon)
```

Bypasses the repository abstraction; invisible to future caching/tracing middleware.

**Fix:** Add `weaponRepo.FindByIDs([]uuid.UUID)`, `itemRepo.FindByIDs([]uuid.UUID)`, and `characterRepo.AppendWeapon(characterID, weapon)`.

---

### 16. Missing Compound Index on `character_resources` — **HIGH IMPACT**

**File:** `backend/models/character_resource.go`

`FindByCharacterAndKey` is called on every resource read/write:
```go
r.db.Where("character_id = ? AND resource_key = ?", characterID, key).First(&resource)
```

No compound index exists on `(character_id, resource_key)`. Each call is a sequential scan filtered by two columns.

**Fix:** Add to the model:
```go
CharacterID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_char_resource_key"`
ResourceKey  string    `gorm:"type:text;not null;uniqueIndex:idx_char_resource_key"`
```

Or via raw migration:
```sql
CREATE UNIQUE INDEX IF NOT EXISTS idx_char_resource_key
  ON character_resources (character_id, resource_key);
```

---

### 17. `isGM` Queries DB on Every GM-Gated Request — **LOW**

**File:** `backend/api/spell_handler.go:92–102`

```go
user, err := h.userRepo.FindByID(id)
return user.IsAdmin || user.Email == gmEmail, nil
```

A full user row is fetched to check two fields. Called on every GM endpoint hit (unchecked spells, spell approval).

**Fix:** Include `is_admin` as a claim in the JWT at login so the middleware can extract it without a DB round-trip.

---

### 18. `gorm.Save()` Issues Full-Row UPDATE — **LOW**

**Files:** `character_repo.go`, `character_resource_repo.go`, `spell_repo.go`, `party_repo.go`, `map_token_repo.go`, `map_element_repo.go`

`db.Save(entity)` generates `UPDATE … SET col1=?, col2=?, … colN=?` for every column, even when only one field changed (e.g., `current_hp`, `initiative_order`).

For high-frequency updates like HP changes or token moves, this sends far more data than needed.

**Fix:** Replace with targeted updates where the changed fields are known:
```go
r.db.Model(&models.Character{}).Where("id = ?", id).
    Updates(map[string]interface{}{"current_hp": newHP, "updated_at": time.Now()})
```

---

## Summary

| # | Issue | Severity | Est. Query Reduction |
|---|-------|----------|----------------------|
| 1 | Double character load on sheet | High | 2 → 1 per sheet view |
| 2 | Party: 4 round-trips on sheet | High | 4 → 1 |
| 3 | Map load for token auth checks | High | 6 queries → 1 per token op |
| 4 | `BulkSetInitiativeOrder` N UPDATEs | High | N → 1 |
| 5 | `ReplaceSkillProficiencies` N inserts | Medium | N → 1 |
| 6 | `replaceComponentsTx` N inserts | Medium | N → 1 |
| 7 | `ProcessRestoration` N+1 saves | Medium | N+1 → 2 |
| 8 | `UpdateCurrentValue` reads outside tx | Medium | correctness |
| 9 | Auth checks load full character | Medium | full load → COUNT |
| 10 | `AddIdentifiedBeast` wasted lookups | Medium | 3 → 1 |
| 11 | Spell `FindByUserID` unbounded | Medium | safety |
| 12 | Character list unbounded | Medium | safety |
| 13 | `buildClassResources` 3 queries | Medium | 3 → 1 |
| 14 | Manual skill second query | Low | 2 → 1 |
| 15 | Raw DB bypasses repos | Low | code quality |
| 16 | Missing compound index on resources | High | full scan → index scan |
| 17 | `isGM` DB lookup per request | Low | 1 → 0 (JWT claim) |
| 18 | `Save()` full-row UPDATE | Low | bandwidth |

**Highest ROI fixes (implement these first):**
- **#1 + #2**: Reduce every character sheet load from 6+ queries to 2
- **#3**: Reduce every token/element write from 7+ queries to 2
- **#16**: Add compound index — zero code change, instant improvement on all resource operations
- **#4**: Single SQL for bulk initiative instead of N individual UPDATEs
