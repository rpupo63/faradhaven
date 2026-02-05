# Character Creation System Upgrade Plan

This document outlines the plan to implement a robust, decision-based character creation process for Faradhaven. It covers backend schema changes, seeding updates, and frontend integration to support D&D 5e mechanics (Ability Scores, Proficiencies, Equipment, etc.).

## 1. Backend Architecture Changes

### 1.1. New Models
We will introduce new models to handle the complexity of starting equipment choices (e.g., "Choose a greataxe OR two handaxes").

#### `ClassStartingEquipmentChoice`
Represents a decision point for a class (e.g., "Weapon Choice 1").
*   `ID` (UUID, PK)
*   `ClassID` (UUID, FK)
*   `Instruction` (String): e.g., "Choose one of the following"
*   `SortOrder` (Int): For UI ordering.

#### `ClassStartingEquipmentOption`
Represents a valid selection within a choice.
*   `ID` (UUID, PK)
*   `ChoiceID` (UUID, FK -> `ClassStartingEquipmentChoice`)
*   `Description` (String): e.g., "Scale Mail and a Shield"
*   `Items` (String Array/JSON): The actual list of items granted (e.g., `["Scale Mail", "Shield"]`).

#### `Race` & `Lineage` Updates
To support ability score calculation:
*   Add `AbilityScoreBonuses` (JSONB) to `Race` and `Lineage` models.
    *   Example: `{"dexterity": 2, "charisma": 1}`

### 1.2. Seeding Strategy
We will update the `FaradhavenClassSeed` struct to include nested equipment choice definitions, moving away from the simple `StartingEquip` string array for choices (though keeping it for "automatic" items).

*   **File:** `@backend/seed/faradhaven_classes/types.go`
    *   Update `FaradhavenClassSeed` to include `EquipmentChoices []EquipmentChoiceSeed`.
*   **File:** `@backend/seed/faradhaven_classes/[classname].go`
    *   Refactor each class seed to define its specific equipment choices.
*   **File:** `@backend/seed/faradhaven_races/[racename].go`
    *   Ensure racial seeds include the correct stat bonuses.

## 2. Character Creation Logic (Service Layer)

We need a dedicated `CharacterCreationService` or extensions to `CharacterService` to handle the logic.

### 2.1. Ability Scores
Support three generation methods:
1.  **Standard Array:** Fixed set `[15, 14, 13, 12, 10, 8]` assigned to attributes by the user.
2.  **Point Buy:** Start with 8s, 27 points to spend. (Cost: 8->0, 13->5, 14->7, 15->9).
3.  **Manual Roll:** 4d6 drop lowest (repeated 6 times). *Note: Backend should probably validate the final numbers fall within reasonable bounds or trust the client for "Manual" mode.*

**Calculation:** `Base Score` + `Race Bonus` + `Lineage Bonus` = `Final Score`.

### 2.2. Hit Points (Level 1)
Logic: `Class.HitDie` (Max) + `Constitution Modifier`.
*   *Note: This must be calculated AFTER ability scores are finalized.*

### 2.3. Proficiencies
*   Validate that selected skills are within the `Class.SkillChoice` list.
*   Validate `SkillChoiceCount` (usually 2, sometimes more).
*   Add Background proficiencies (if Backgrounds are implemented in the future, currently just Class).

### 2.4. Starting Equipment
*   **Auto-Grant:** Items in `Class.StartingEquip` (renamed to `AutomaticEquipment` for clarity?) are added to the character.
*   **Choices:** User selections from `ClassStartingEquipmentChoice` are processed and added.
*   **Storage:** Since we don't have a full Inventory model yet, these will be stored as a `Text[]` array `Inventory` on the `Character` model (or a new `CharacterEquipment` relation if preferred).

## 3. Frontend Implementation

### 3.1. Character Creation Wizard
A multi-step React component (`CharacterCreationWizard.tsx`) in `@frontend/src/components/`.

**Steps:**
1.  **Race Selection:**
    *   Visual selection of Race & Lineage.
    *   Preview `Race.AbilityScoreBonuses`.
2.  **Class Selection:**
    *   Visual selection of Class.
    *   Preview `Class.HitDie`, `Class.PrimaryAbility`.
3.  **Ability Scores:**
    *   Tabs for "Standard Array", "Point Buy", "Roll".
    *   Interactive UI to assign scores to STR/DEX/CON/etc.
    *   Real-time modifier preview.
4.  **Class Options:**
    *   **Proficiencies:** Checkboxes limited by `Class.SkillChoiceCount`.
    *   **Equipment:** Radio groups for each `ClassStartingEquipmentChoice`.
5.  **Description:** Name, Alignment, Appearance.
6.  **Summary & Save:** Final review before POST.

### 3.2. Character Sheet View Updates
Update `@frontend/src/components/CharacterSheetView.tsx` to display the seeded data correctly.

*   **Class Features Table (Right Column):**
    *   **Action:** Fetch `LevelFeature`s for the character's class/level.
    *   **Display:** Populate the "Class Features" card which is currently a placeholder. Group by Level or "Active/Passive".
*   **Racial Traits:**
    *   Ensure `AbilityScoreBonuses` are visually accounted for (maybe show base + bonus tooltip).
*   **Equipment:**
    *   Replace "No equipment tracked yet" with a list of the starting equipment generated.

## 4. Execution Plan (Todo)

### Phase 1: Backend Models & Seeds
1.  Create `ClassStartingEquipmentChoice` and `ClassStartingEquipmentOption` structs in `@backend/models`.
2.  Add `AbilityScoreBonuses` to `Race` model.
3.  Update `AllModels()` in `models.go`.
4.  Update `FaradhavenClassSeed` type.
5.  Refactor `Ironwright`, `Sanguinist`, etc. seeds with equipment choices.
6.  Run seed migration.

### Phase 2: Creation Logic
1.  Create `POST /api/characters/create` endpoint.
2.  Implement validation logic (Skill counts, valid equipment options).
3.  Implement `CalculateMaxHP(class, con)` logic.

### Phase 3: Frontend Wizard
1.  Build `CharacterCreationWizard` component.
2.  Implement Point Buy/Standard Array logic in frontend state.
3.  Connect to backend for static data (Races, Classes, Options).

### Phase 4: Sheet Integration
1.  Update `CharacterSheetView` to fetch and render `LevelFeatures`.
2.  Update `CharacterSheetView` to render `Inventory`.

---
*Created by Gemini CLI*
