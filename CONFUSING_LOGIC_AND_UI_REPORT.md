# Report on Confusing Logic and Implementation

This report details confusing mechanics in the `backend/seed/faradhaven_classes/` directory and their unclear implementation in the `frontend/src/pages/CharacterSheetPage.tsx`. Each section identifies an issue and provides suggestions for improvement.

---

### 1. The Lorewright: Overly Complex Resource Management

**Backend Logic Issue (`lorewright.go`)**:
The Lorewright class juggles three separate and overlapping resource/status systems, creating significant cognitive load for the player:
1.  **Harvest Bank**: This is fragmented into `SkillSlots`, `AttackSlots`, and `RecipeSlots`, each with a different progression rate. The presence of a legacy `EchoSlots` field adds to the confusion.
2.  **The Fracture & Trauma**: This system requires a complex DC calculation (`10 + (Creature CR - (Your Level / 3))`) to avoid `Trauma` points, which apply staged debuffs.
3.  **Madness Die & Table**: A separate `d4` to `d20` die is rolled on *every use* of a harvested ability, creating random, high-stakes outcomes from a 50-entry table. This constant rolling interrupts gameplay flow.

**Frontend Implementation Issue (`CharacterSheetPage.tsx`)**:
While the frontend has a `LorewrightFeaturesCard` and `HarvestModal`, this is a bespoke implementation that highlights a larger issue: the main character sheet is not equipped to handle such complex, class-specific mechanics.
*   The generic `CharacterSheetView` has no apparent way to display the Harvest Bank slots, Trauma points, or the current Madness Die size. This information is critical to the player but is hidden away in a side component.
*   The `HarvestModal` flow is incomplete. It doesn't allow the player to see *what* abilities can be harvested from a creature before committing, removing strategic choice from the core class mechanic.

**Suggestions**:
*   **Backend**:
    *   **Unify Mental State Mechanics**: Merge "The Fracture" and the "Madness Die" into a single system. For instance, failing a harvest save could add points to a "Madness" track. At certain thresholds, this track could trigger effects from the madness table, making the outcome less frequent and more predictable.
    *   **Simplify the Harvest Bank**: Replace the three separate slot types with a single pool of "Harvest Slots". Allow any type of ability (skill, attack, recipe) to be slotted, simplifying tracking.
*   **Frontend**:
    *   **Create a Dynamic Resource Section**: The `CharacterSheetPage` should feature a dedicated, dynamically-rendered section for class-specific resources. For the Lorewright, this would clearly display harvested abilities, the current Madness level, and its effects.
    *   **Enhance the HarvestModal**: Add a step to the modal that shows the player the potential skills, attacks, and recipes of the creature *before* they choose to harvest.

---

### 2. The Sanguinist: Punitive and Unclear Notoriety System

**Backend Logic Issue (`sanguinist.go`)**:
The "Moral Seesaw (Notoriety System)" is an innovative but convoluted mechanic.
*   It forces constant mental arithmetic to track the balance between "Medical Prodigy" (MP) and "Blood Rage" (BR) points.
*   The feedback loop is punitive. Being a "specialist" (e.g., `MP > BR + 2`) results in a negative side effect (`Sanguine Backfire` healing enemies) rather than a strategic trade-off. This discourages specialization.
*   The class also manages `Blood Ichor` and `Unstable Components`, making for three resource pools to track (Ichor, Components, Notoriety).

**Frontend Implementation Issue (`CharacterSheetPage.tsx`)**:
There is no specific UI mentioned or visible for the Sanguinist.
*   The `CharacterSheetView` cannot display the MP/BR notoriety balance.
*   The UI provides no warning when a player enters the "Overloaded Healer" or "Starving Predator" state, leading to surprise negative outcomes that feel unfair.

**Suggestions**:
*   **Backend**:
    *   **Refactor to Trade-offs**: Instead of punishment, make specialization a choice with clear trade-offs. A "Starving Predator" could gain a bonus to `Bite` damage but have reduced `Blood Graft` healing potency, rather than having the heal harm allies.
    *   **Consolidate Resources**: Consider removing `Unstable Components`. `Sanguine Extraction` could simply restore `Blood Ichor` directly, simplifying the resource loop.
*   **Frontend**:
    *   **Add a Notoriety Visual**: Implement a visual component like a slider or scale on the character sheet to show the MP vs. BR balance.
    *   **Provide Status Warnings**: When a character enters a feedback loop state, display a clear status icon (e.g., "Ravenous") with a tooltip explaining the effect. The UI should warn the player *before* they confirm an action that will be negatively affected.

---

### 3. The Ironwright: Underdeveloped Scavenging Loop

**Backend Logic Issue (`ironwright.go`)**:
The core gameplay loop of scavenging components lacks strategic depth.
*   The class can scavenge many *distinct* component types ("Cog", "Gear", "Ignis").
*   However, core abilities like "Construct (Sentry)" cost "2 Components (any type)". This makes the act of collecting different types of components feel meaningless. If any component will do, why differentiate them?

**Frontend Implementation Issue (`CharacterSheetPage.tsx`)**:
The character sheet has no visible support for a summoner/pet class.
*   A player cannot see their inventory of scavenged components.
*   There is no way to track active constructs, their current HP, or their duration. This information is essential for a class focused on managing minions.

**Suggestions**:
*   **Backend**:
    *   **Add Strategic Depth**: Make constructs require *specific* components (e.g., a Sentry requires 1 `Ferrum` and 1 `Fulgur`). This would make the choice of what to scavenge from which enemy a meaningful tactical decision.
*   **Frontend**:
    *   **Implement Pet/Resource UI**: Add a dedicated "Ironwright" section to the sheet that includes:
        1.  A **Component Inventory** listing each component type and its quantity.
        2.  A **Constructs** panel showing the `Concurrency` limit (e.g., "Slots: 2/3") and a list of active constructs with their individual stats (HP, AC, remaining duration).

---

### 4. General Issue: Brittle and Overly Complex Data Structures

**Backend Logic Issue (Multiple Files)**:
The `ClassLevelSeed` struct in `types.go` is bloated with over a dozen unique, class-specific progression fields (e.g., `ConcurrencyLimit`, `YieldDie`, `TimerDuration`, `FeralBonus`, `SkillSlots`, `AttackSlots`).
*   This makes the `seed.go` file incredibly complex, as it has to manually handle each of these fields for every class at every level.
*   It makes the class data brittle and hard to maintain or scale. Adding a new class with a unique resource requires changing this central data structure and all associated seed logic.

**Frontend Implementation Issue (`CharacterSheetPage.tsx`)**:
The frontend is forced to handle this complexity, leading to divergent logic.
*   The code has to use two different normalization functions (`normalizeApiSheet` and `normalizeComputedSheet`) to handle data depending on its source.
*   It's not scalable to build a unique UI for every one of these custom resources.

**Suggestions**:
*   **Backend**:
    *   **Standardize Progression**: Where possible, unify resources. Can some of these be derived from a formula (e.g., `Proficiency Bonus + INT modifier`) instead of being hardcoded in a giant map for 20 levels?
    *   **Abstract Class Resources**: Instead of adding new fields to `ClassLevel` for every new class, create a more generic `ClassResource` model that can be associated with a class, defining its name, progression, and mechanics.
*   **Frontend**:
    *   **Create a Generic Resource Component**: Design a single, flexible React component that can display any class-specific resource. It would take properties like `name`, `currentValue`, `maxValue`, and `description` from the `ClassResource` data sent by the backend. The character sheet would then dynamically render these components based on the character's class.
