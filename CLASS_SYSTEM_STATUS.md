# Class system — unified status

Single reference for what used to be split across **CLASS_SEEDING_REFACTOR_REPORT.md**, **CONFUSING_LOGIC_AND_UI_REPORT.md**, **CLASS_SEEDING_AND_UX_NOTES.md**, **DEPRECATION_REPORT.md**, and **`backend/CASTER_RESOURCE_BALANCE.md`**.

---

## Completed (historical / shipped)

### API `TotalHP` / `total_hp`

Removed from `CharacterSheetResponse` and frontend sheet types; HP comes from persisted / computed values.

### Class seeds: old `ResourceType` / `ResourceName` / `ResourceRestoreType`

Removed from **`FaradhavenClassSeed`**. Seeding uses **`ResourceDefinitions`** + **`LevelProgression.Resources`**. API **`resource_name`** on character resources is the display label, not the deprecated trio.

### Seeding merge contract + optional overrides ([`seed.go`](backend/seed/faradhaven_classes/seed.go), [`types.go`](backend/seed/faradhaven_classes/types.go))

- **Documented** on `SeedFaradhavenClasses` and in `ClassLevelSeed` comments: defaults first (HP gain, proficiency, max spell points, ASI), then `LevelProgression[level]` overlays fields and writes `class_level_resources`.
- **`ProficiencyBonus *int`** and **`AbilityScoreImprovement *int`** on `ClassLevelSeed`: when non-nil, replace the helper values for that level (so explicit `0` is valid).

### Ironwright constructs vs component catalog ([`minion.go`](backend/models/minion.go), [`ironwright.go`](backend/seed/faradhaven_classes/ironwright.go))

- **`GetConstructTemplates()`** enforces **named** `RequiredComponents`; `CreateConstruct` deducts via `validateAndDeductComponents`.
- **Striker / Titan** recipes updated to use **only components that exist** in [`faradhaven_components/components.go`](backend/seed/faradhaven_components/components.go) (replaced invalid names like Imbue, Heat, Conduct, Mend).
- **Level feature text** aligned with templates (Sentry **Fulgur + Ferrum**, Striker five-part recipe, Titan ten-part recipe, Rapid Assembly schematics use Vita / Push / Zone / Fulgur / Ferrum).

### Lorewright + Sanguinist copy (seeds)

- **Scholar’s Legacy** references **Harvest Slots** instead of a separate “Recipe Slot.”
- **Sanguinist** concept and **Moral Seesaw** text match server behavior in [`hp_service.go`](backend/services/hp_service.go): **Backfire** = extra **1d8 to self** on offensive features when MP-heavy; **Ravenous** = **1d8 to self** before Blood Graft heal resolves when BR-heavy; **Balanced** called out as complication-free.

### Frontend normalization ([`classLevelData.ts`](frontend/src/lib/classLevelData.ts))

- Shared helpers: **`readAbilityScoresForSheet`**, **`abilityModsFromScores`**, **`buildSkillBonusMap`**, **`buildSavingThrowBonusMap`**, **`collectFeaturesThroughLevel`**, **`computeInitiativeModifier`** (Jack of All Trades), **`ensureSanguinistBiteWeapon`** / **`virtualSanguinistBiteWeapon`**.
- Vitest: [`classLevelData.test.ts`](frontend/src/lib/classLevelData.test.ts).

### Class resource UI registry

- Class-specific blocks live in [`ClassResourceClassSections.tsx`](frontend/src/components/character-sheet/ClassResourceClassSections.tsx); lookup in [`classResourceSectionRegistry.tsx`](frontend/src/components/character-sheet/classResourceSectionRegistry.tsx); props type in [`classResourceSectionTypes.ts`](frontend/src/components/character-sheet/classResourceSectionTypes.ts).
- [`ClassResourceDisplay.tsx`](frontend/src/components/character-sheet/ClassResourceDisplay.tsx) renders **`ClassResourceExtraSections`**.

---

## Backend: `backend/seed/faradhaven_classes/` (current behavior)

- **`ComponentPool`** + **`HeritageSpeciesComponents()`** → `class_components`.
- **`ClassLevel`** features → **`level_features`** only.
- **`ClassLevelSeed`**: shared combat/spell fields + **`Resources`** map; optional **`ProficiencyBonus` / `AbilityScoreImprovement`** pointers for per-level overrides.

---

## Optional follow-ups (not blocking)

| Topic | Notes |
|--------|--------|
| **Lorewright complexity** | Multiple subsystems (harvest, trauma, strain, madness) remain heavy at the table; further **mechanical** consolidation would be a separate design pass. |
| **Sanguinist** | Notoriety + ichor + unstable components still many knobs; rules text now matches **self**-damage backfire/ravenous. Changing to ally/enemy targeting would need **`hp_service`** updates. |
| **Normalizers** | API vs computed sheets still **two entry points**, but they share the same **math helpers**; extend tests when adding sheet fields. |
| **Overclock / generic “5 components”** | Narrative spends not tied to `CreateConstruct` may stay “any mixture” by design. |

---

## Deprecation / cleanup queue

**Empty.** Add rows here when retiring API fields or seed shapes.
