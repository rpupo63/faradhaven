# Deprecation and Dead Code Report

This report details deprecated or dead code found within the codebase.

## 1. `TotalHP` Field

*   **Location:**
    *   Backend: `backend/api/types.go` (in `CharacterSheetResponse`)
    *   Frontend: `frontend/src/types/game/api.ts` (as `total_hp`)
*   **Description:**
    A legacy `TotalHP` field was calculated on the backend as a fallback for characters that did not have `max_hp` and `current_hp` persisted in the database. The frontend used this `total_hp` field as a fallback if `max_hp` was not present.
*   **Status:** Removed. `max_hp` and `current_hp` are now directly used from the database.
*   **Action Taken:**
    *   Removed `TotalHP` from `CharacterSheetResponse` struct in `backend/api/types.go`.
    *   Removed `totalHP` calculation logic from `getCharacterSheet` in `backend/api/character_sheet_handler.go`.
    *   Removed `total_hp` from `ApiCharacterSheet` interface in `frontend/src/types/game/api.ts`.
    *   Removed `total_hp` from `CharacterSheetComputed` and `NormalizedCharacterSheet` interfaces in `frontend/src/types/game/character.ts`.
    *   Updated `computeCharacterSheetFromApi` and `normalizeApiSheet` in `frontend/src/lib/classLevelData.ts` to no longer reference `total_hp`.

## 2. Legacy Class Resource System (`ResourceType`, `ResourceName`, `ResourceRestoreType`)

This section outlines the plan to fully remove the deprecated `ResourceType`, `ResourceName`, and `ResourceRestoreType` fields from the class seed files located in `backend/seed/faradhaven_classes/`.

These fields are already superseded by the `ResourceDefinitions` and `LevelProgression.Resources` maps within each class seed. The new system provides a more structured and extensible way to define class resources, and the seeder logic in `backend/seed/faradhaven_classes/seed.go` is already set up to process it.

### What will be removed?

In each of the class seed files (e.g., `ironwright.go`, `sanguinist.go`), the following three lines will be deleted:

```go
ResourceType:        "...",
ResourceName:        "...",
ResourceRestoreType: "...",
```

### Class-by-Class Migration Plan

Here is an analysis of each class and confirmation that their resources are correctly captured by the new system, or a plan to migrate them if not.

*   **Ironwright**
    *   **Deprecated:** `ResourceType: "component_items"`
    *   **New System:** `ResourceDefinitions` for `concurrency_limit` and `yield_die` are already present.
    *   **Action:** No further migration needed. The legacy fields can be safely removed.

*   **Lorewright**
    *   **Deprecated:** `ResourceType: "harvest_slots"`
    *   **New System:** The `ResourceDefinitions` section correctly defines `echo_slots`.
    *   **Action:** No further migration needed. The legacy fields can be safely removed.

*   **Mutagen**
    *   **Deprecated:** `ResourceType: "madness"`
    *   **New System:** The `ResourceDefinitions` section correctly defines `madness_base_dc` and `feral_bonus`.
    *   **Action:** No further migration needed. The legacy fields can be safely removed.

*   **Piston Brawler**
    *   **Deprecated:** `ResourceType: "stability"`
    *   **New System:** The `ResourceDefinitions` section correctly defines `max_stability`.
    *   **Action:** No further migration needed. The legacy fields can be safely removed.

*   **Powder Mage**
    *   **Deprecated:** `ResourceType: "timer"`
    *   **New System:** The `ResourceDefinitions` for `timer_duration`, `speed_dial_slots`, and `max_spell_length` are already present.
    *   **Action:** No further migration needed. The legacy fields can be safely removed.

*   **Sanguinist**
    *   **Deprecated:** `ResourceType: "blood_ichor"`
    *   **New System:** The `ResourceDefinitions` section correctly defines `max_blood_ichor`.
    *   **Action:** No further migration needed. The legacy fields can be safely removed.

*   **Rift Weaver**
    *   **Deprecated:** `ResourceType: "spell_points"`
    *   **New System:** The `ResourceDefinitions` section is currently empty. The class uses the default `maxSpellPointsByLevel` calculation.
    *   **Action:** To migrate, we will add a `ResourceDefinition` for `spell_points` to `rift_weaver.go`.
        ```go
        ResourceDefinitions: []ResourceDefinitionSeed{
            {
                Key: "spell_points",
                DisplayName: "Spell Points",
                Category: "pool",
                Description: "A pool of magical energy for casting spells.",
                IsTrackable: true,
                RestoreOnLongRest: true,
                DisplayOrder: 1,
            },
        },
        ```
        And update the `LevelProgression` to use this key. Since it uses a default formula, no values need to be added to the `Resources` map in the level progression.

*   **Vapor Blade**
    *   **Deprecated:** `ResourceType: "shadow_points"`
    *   **New System:** The `ResourceDefinitions` section is empty. The class seems to be intended to use a resource called "Shadow Points".
    *   **Action:** To migrate, we will add a `ResourceDefinition` for `shadow_points` to `vapor_blade.go`. We will also need to define its progression in `LevelProgression`. Based on the features, it seems to be a pool that restores on long rest. We'll assume a progression similar to other classes, like Proficiency Bonus + Dexterity modifier.
        ```go
        // In vapor_blade.go
        ResourceDefinitions: []ResourceDefinitionSeed{
            {
                Key: "shadow_points",
                DisplayName: "Shadow Points",
                Category: "pool",
                Description: "A pool of shadowy energy for fueling abilities.",
                IsTrackable: true,
                RestoreOnLongRest: true,
                DisplayOrder: 1,
            },
        },

        // In vaporBladeLevelProgression()
        // Assuming a formula of (Proficiency Bonus + Dex Mod).
        // For a character with +3 Dex mod:
        // Level 1: 2+3=5, Level 5: 3+3=6, Level 9: 4+3=7, etc.
        ```