# Report on Backend Class Seeding Logic

This report details findings on redundant and confusing logic within the `backend/seed/faradhaven_classes/` directory and provides suggestions for improvement. The goal of these suggestions is to increase clarity, reduce redundancy, and make the class seeding process more maintainable.

## 1. Redundant and Unused `ComponentPool`

- **Issue:** The `FaradhavenClassSeed` struct in `types.go` contains a `ComponentPool []string` field. This field is populated in some class seed files (e.g., `piston_brawler.go`) but is empty in others (e.g., `lorewright.go`). The core seeding logic in `seed.go` completely ignores this field. Instead, it relies on `AllComponentClassMappings()` from `faradhaven_components.go`, which derives its data from the `ClassComponentNames()` map. This makes the `ComponentPool` field dead code, adding clutter and confusion for developers trying to understand how to add or modify class components.

- **Suggestion:** Remove the `ComponentPool` field from the `FaradhavenClassSeed` struct in `backend/seed/faradhaven_classes/types.go` and from all individual class seed files where it is used. The single source of truth for class-to-component mapping should be the `ClassComponentNames` map in `faradhaven_components.go`.

## 2. Inconsistent and Confusing Level Progression Logic

- **Issue:** There are multiple sources of truth for `ClassLevel` data.
    1.  The `LevelProgression` map in each class seed file (e.g., `lorewright.go`) defines class-specific values like `MaxStability` or `FeralBonus`.
    2.  The `SeedFaradhavenClasses` function in `seed.go` contains helper functions (`proficiencyByLevel`, `abilityScoreImprovementByLevel`, `maxSpellPointsByLevel`) that calculate standard D&D progression values.
    3.  The `seed.go` function then merges these two sources when creating `models.ClassLevel` objects.

    This division of logic is confusing. A developer needs to check both the class seed file and the main `seed.go` file to understand the full level progression for a class. It's unclear at a glance whether a value is hardcoded, calculated, or comes from the structured seed data.

- **Suggestion:** Make the `LevelProgression` map in each class seed file the definitive source for all level-specific data. The generic helper functions in `seed.go` (like `proficiencyByLevel`) should only be used as fallbacks if a value is *not* provided in the `LevelProgression` map.

    **Example Refactor:**
    In `seed.go`, when creating a `ClassLevel`:
    ```go
    // Inside the level creation loop...
    cl := models.ClassLevel{
        // ... basic fields
    }

    // Apply structured progression data from the seed file first
    if lp, ok := cs.LevelProgression[level]; ok {
        // Apply all fields from lp
        cl.MaxStability = lp.MaxStability
        cl.FeralBonus = lp.FeralBonus
        // ... etc.
    }

    // Apply calculated defaults ONLY if not set by the progression map
    if cl.ProficiencyBonus == 0 {
        cl.ProficiencyBonus = proficiencyByLevel(level)
    }
    if cl.AbilityScoreImprovement == 0 {
        cl.AbilityScoreImprovement = abilityScoreImprovementByLevel(level)
    }
    ```
    This approach prioritizes explicit data from the class seed files while still providing sensible defaults for common D&D mechanics, making the system more declarative and easier to understand.

## 3. Legacy `Features` Text Blob

- **Issue:** The `SeedFaradhavenClasses` function uses a helper `buildLevel1Features` to concatenate various class properties (proficiencies, skills, etc.) into a single unstructured string, which is then stored in the `ClassLevel.Features` field for level 1. For other levels, it concatenates the `Name` and `Description` of `FeatureSeed`s. The `level_features` table is already being populated correctly with structured data. This large text blob is a legacy approach that is difficult to parse, display on the frontend, and maintain.

- **Suggestion:**
    1.  **Deprecate `ClassLevel.Features`:** This field should be removed from the `models.ClassLevel` model.
    2.  **Rely on Structured Data:** The frontend should query and display the structured data from the `level_features` table, which is properly linked to class levels. Information currently in the level 1 blob (like proficiencies and saving throws) is already available as distinct fields on the `models.Class` object and should be queried from there.
    3.  **Remove `buildLevel1Features`:** This helper function in `seed.go` would no longer be needed and should be deleted.

By moving away from unstructured text blobs and relying on the relational data already being seeded, the system becomes more robust, easier to query, and simpler for frontend developers to consume.
