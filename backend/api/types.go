package api

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// routeHandlers contains all the handlers for the API routes
type routeHandlers struct {
	authHandler            *authHandler
	userHandler            *userHandler
	characterHandler       *characterHandler
	spellHandler           *spellHandler
	beastHandler           *beastHandler
	levelHandler           *levelHandler
	weaponHandler          *weaponHandler
	itemHandler            *itemHandler
	componentHandler       *componentHandler
	effectHandler          *effectHandler
	characterEffectHandler *characterEffectHandler
	resourceHandler        *resourceHandler
	minionHandler          *minionHandler
	noteHandler            *noteHandler
	mapHandler             *mapHandler
	mechanicsHandler       *MechanicsHandler
	corpseHandler          *corpseHandler
	linkHandler            *linkHandler
	harvestHandler         *harvestHandler
	madnessHandler         *madnessHandler
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error   string `json:"error" example:"Internal Server Error"`
	Status  string `json:"status" example:"error"`
	Field   string `json:"field,omitempty" example:"title"`
	Details string `json:"details,omitempty" example:"Additional error details"`
	Cause   string `json:"cause,omitempty" example:"Underlying error cause"`
}

// Map Request/Response Types

type CreateMapRequest struct {
	Name          string `json:"name"`
	RoomCode      string `json:"room_code"`
	BackgroundURL string `json:"background_url"`
	GridRows      int    `json:"grid_rows"`
	GridCols      int    `json:"grid_cols"`
	TileSize      int    `json:"tile_size"`
}

type UpdateMapRequest struct {
	Name          *string `json:"name,omitempty"`
	BackgroundURL *string `json:"background_url,omitempty"`
	GridRows      *int    `json:"grid_rows,omitempty"`
	GridCols      *int    `json:"grid_cols,omitempty"`
	TileSize      *int    `json:"tile_size,omitempty"`
}

type CreateTokenRequest struct {
	CharacterID    *uuid.UUID `json:"character_id,omitempty"`
	AssignedUserID *uuid.UUID `json:"assigned_user_id,omitempty"`
	Name           string     `json:"name"`
	ImageURL       string     `json:"image_url"`
	TokenType      string     `json:"token_type"` // "pc" or "npc"
	GridX          int        `json:"grid_x"`
	GridY          int        `json:"grid_y"`
	Size           int        `json:"size"`
	Color          string     `json:"color"`
	Visible        bool       `json:"visible"`
}

type UpdateTokenRequest struct {
	GridX          *int       `json:"grid_x,omitempty"`
	GridY          *int       `json:"grid_y,omitempty"`
	Visible        *bool      `json:"visible,omitempty"`
	Size           *int       `json:"size,omitempty"`
	Color          *string    `json:"color,omitempty"`
	AssignedUserID *uuid.UUID `json:"assigned_user_id,omitempty"` // For DM to reassign
}

// User Request/Response Types

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name              *string    `json:"name,omitempty"`
	Email             *string    `json:"email,omitempty"`
	ActiveCharacterID *uuid.UUID `json:"active_character_id,omitempty"`
}

type SetActiveCharacterRequest struct {
	CharacterID *uuid.UUID `json:"character_id"` // nil to clear active character
}

// Character Request/Response Types

type CreateCharacterRequest struct {
	UserID             uuid.UUID   `json:"user_id"`
	Name               string      `json:"name"`
	RaceID             uuid.UUID   `json:"race_id"`
	LineageID          *uuid.UUID  `json:"lineage_id"`
	ClassID            uuid.UUID   `json:"class_id"`
	Level              int         `json:"level"`
	Spellbook          []string    `json:"spellbook"`
	Strength           int         `json:"strength"`
	Dexterity          int         `json:"dexterity"`
	Constitution       int         `json:"constitution"`
	Intelligence       int         `json:"intelligence"`
	Wisdom             int         `json:"wisdom"`
	Charisma           int         `json:"charisma"`
	CurrentSpellPoints int         `json:"current_spell_points"`
	Money              int64       `json:"money"`
	SkillProficiencies []string    `json:"skill_proficiencies"`         // D&D 5e skill ids (e.g. "persuasion", "stealth")
	Languages          []string    `json:"languages"`                   // Languages chosen during character creation
	EquipmentChoices   []uuid.UUID `json:"equipment_choices"`           // IDs of selected ClassStartingEquipmentOption
	PrimaryWeaponID    *uuid.UUID  `json:"primary_weapon_id,omitempty"` // selected primary weapon for signature items
}

type CreationOptionsResponse struct {
	Races     []models.Race  `json:"races"`
	Classes   []models.Class `json:"classes"`
	PointsMax int            `json:"points_max"` // e.g. 27 for point buy
}

type UpdateCharacterRequest struct {
	Name               *string    `json:"name,omitempty"`
	RaceID             *uuid.UUID `json:"race_id,omitempty"`
	LineageID          *uuid.UUID `json:"lineage_id,omitempty"`
	ClassID            *uuid.UUID `json:"class_id,omitempty"`
	Level              *int       `json:"level,omitempty"`
	Spellbook          []string   `json:"spellbook,omitempty"`
	Strength           *int       `json:"strength,omitempty"`
	Dexterity          *int       `json:"dexterity,omitempty"`
	Constitution       *int       `json:"constitution,omitempty"`
	Intelligence       *int       `json:"intelligence,omitempty"`
	Wisdom             *int       `json:"wisdom,omitempty"`
	Charisma           *int       `json:"charisma,omitempty"`
	CurrentSpellPoints *int       `json:"current_spell_points,omitempty"`
	Money              *int64     `json:"money,omitempty"`
	Notoriety          *int       `json:"notoriety,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
	SkillProficiencies []string   `json:"skill_proficiencies,omitempty"` // D&D 5e skill ids
}

type CastSpellRequest struct {
	SpellID     *uuid.UUID `json:"spell_id,omitempty"`
	SpellLevel  int        `json:"spell_level"`
	Components  []string   `json:"components,omitempty"` // For dynamic casting like Powder Mage
	ResourceKey *string    `json:"resource_key,omitempty"` // Key of the resource to deduct from (e.g., "spell_points", "shadow_points")
}

// UpdateNotorietyRequest is used to update notoriety manually
type UpdateNotorietyRequest struct {
	Delta int `json:"delta"`
}

// UpdateSanguineNotorietyRequest is used to update sanguine notoriety points from modular choices
type UpdateSanguineNotorietyRequest struct {
	MPChange int `json:"mp_change"`
	BRChange int `json:"br_change"`
}

// ClassWithLevelsResponse is a class with all levels 1-20 for the compendium/book view
type ClassWithLevelsResponse struct {
	*models.Class
	Levels       []models.ClassLevel `json:"levels"`
	MadnessTable map[int]string      `json:"madness_table,omitempty"`
}

// CharacterSheetResponse is the fully calculated character sheet (Class + ClassLevel joined)
type CharacterSheetResponse struct {
	Character                *models.Character           `json:"character"`
	Class                    *models.Class               `json:"class"`
	ClassLevel               *models.ClassLevel          `json:"class_level"`
	MadnessTable             map[int]string              `json:"madness_table,omitempty"`
	MaxHP                    int                         `json:"max_hp"`     // Character's maximum HP (persisted)
	CurrentHP                int                         `json:"current_hp"` // Character's current HP (persisted)
	TempHP                   int                         `json:"temp_hp"`    // Temporary HP
	AC                       int                         `json:"ac"`         // 8 + ProficiencyBonus + DexMod
	SaveDC                   int                         `json:"save_dc"`    // 8 + ProficiencyBonus + PrimaryAbilityMod
	MaxSpellPoints           int                         `json:"max_spell_points"`
	CurrentSpellPoints       int                         `json:"current_spell_points"`
	SavingThrowProficiencies []string                    `json:"saving_throw_proficiencies"`    // ability ids from class.SavingThrows
	AvailableComponents      []models.Component          `json:"available_components"`          // combined class + race components for spell crafting
	HitDiceTotal             int                         `json:"hit_dice_total"`                // Total hit dice = Level
	HitDiceRemaining         int                         `json:"hit_dice_remaining"`            // Available hit dice = Level - HitDiceUsed
	HitDie                   int                         `json:"hit_die"`                       // Class hit die (e.g., 10 for d10)
	Money                    int64                       `json:"money"`                         // Currency in copper pieces
	MeleeAttackBonus         int                         `json:"melee_attack_bonus"`            // Proficiency + STR mod
	RangedAttackBonus        int                         `json:"ranged_attack_bonus"`           // Proficiency + DEX mod
	RaceTraits               []models.Trait              `json:"race_traits"`                   // Race traits for the character's race
	Lineage                  *models.Lineage             `json:"lineage,omitempty"`             // Character's lineage (sub-race)
	InventoryWeapons         []CharacterWeaponResponse   `json:"inventory_weapons"`             // Detailed weapon objects with modifiers
	InventoryItems           []models.Item               `json:"inventory_items"`               // Detailed item objects
	Components               []models.CharacterComponent `json:"components"`                    // Character's current component inventory
	HarvestedAbilities       models.HarvestedAbilities   `json:"harvested_abilities,omitempty"` // Lorewright's harvested abilities

	// --- Generic Class Resources ---
	ClassResources []ClassResourceResponse `json:"class_resources,omitempty"` // Dynamic class resource array
}

// ClassResourceResponse represents a single class resource in the character sheet API.
// This is the generic system that replaces individual class-specific resource fields.
type ClassResourceResponse struct {
	Key          string `json:"key"`                     // machine identifier (e.g., "concurrency_limit")
	DisplayName  string `json:"display_name"`            // UI label (e.g., "Concurrency Limit")
	Category     string `json:"category"`                // rendering hint: "pool", "die_size", "limit", "slot_count", "modifier", "state"
	Description  string `json:"description,omitempty"`   // tooltip text
	Value        int    `json:"value"`                   // progression value at current level (from ClassLevelResource)
	CurrentValue *int   `json:"current_value,omitempty"` // mutable current value (from CharacterResource, nil if not trackable)
	MaxValue     *int   `json:"max_value,omitempty"`     // max value (from CharacterResource, nil if not trackable)
	IsTrackable  bool   `json:"is_trackable"`            // whether this has mutable state
	DisplayOrder int    `json:"display_order"`           // UI ordering
}

type CharacterWeaponResponse struct {
	CharacterWeaponID string                   `json:"character_weapon_id"`
	Weapon            models.Weapon            `json:"weapon"`
	IsPrimary         bool                     `json:"is_primary"`
	CustomName        *string                  `json:"custom_name,omitempty"`
	ActiveModifiers   []WeaponModifierResponse `json:"active_modifiers,omitempty"`
}

type WeaponModifierResponse struct {
	ModifierType string            `json:"modifier_type"`
	IsActive     bool              `json:"is_active"`
	BonusDamage  []BonusDamageInfo `json:"bonus_damage,omitempty"`
	Metadata     interface{}       `json:"metadata,omitempty"`
}

type BonusDamageInfo struct {
	Dice       string `json:"dice"`
	DamageType string `json:"damage_type"`
}

// Spell Request/Response Types

type CreateSpellRequest struct {
	UserID        uuid.UUID          `json:"user_id"`
	CharacterID   *uuid.UUID         `json:"character_id,omitempty"` // Optional: character who prepared this spell
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	ComponentIDs  []uuid.UUID        `json:"component_ids"`
	SlotLevel     int                `json:"slot_level"`
	Type          string             `json:"type"`
	Range         *string            `json:"range,omitempty"`
	Duration      *string            `json:"duration,omitempty"`
	Concentration bool               `json:"concentration"`
	SaveAttr      *string            `json:"save_attr,omitempty"`
	DamageDice    *string            `json:"damage_dice,omitempty"`
	DamageType    *models.DamageType `json:"damage_type,omitempty"`
	AddModifier   bool               `json:"add_modifier"`
}

type UpdateSpellRequest struct {
	Name          *string            `json:"name,omitempty"`
	Description   *string            `json:"description,omitempty"`
	ComponentIDs  []uuid.UUID        `json:"component_ids,omitempty"`
	SlotLevel     *int               `json:"slot_level,omitempty"`
	CharacterID   *uuid.UUID         `json:"character_id,omitempty"`
	Type          *string            `json:"type,omitempty"`
	Range         *string            `json:"range,omitempty"`
	Duration      *string            `json:"duration,omitempty"`
	Concentration *bool              `json:"concentration,omitempty"`
	SaveAttr      *string            `json:"save_attr,omitempty"`
	DamageDice    *string            `json:"damage_dice,omitempty"`
	DamageType    *models.DamageType `json:"damage_type,omitempty"`
	AddModifier   *bool              `json:"add_modifier,omitempty"`
}

// Purchase Request Types

type PurchaseItemRequest struct {
	ItemID   uuid.UUID `json:"item_id"`
	ItemType string    `json:"item_type"` // "item" or "weapon"
}

// UpdateBackstoryRequest is used to update a character's backstory
type UpdateBackstoryRequest struct {
	Backstory         string  `json:"backstory"`
	BackstoryHexColor *string `json:"backstory_hex_color,omitempty"`
}

// Beast Request/Response Types

type CreateBeastRequest struct {
	UserID           uuid.UUID             `json:"user_id"`
	Name             string                `json:"name"`
	ImageURL         *string               `json:"image_url,omitempty"`
	Size             models.CreatureSize   `json:"size"`
	Type             models.CreatureType   `json:"type"`
	Alignment        string                `json:"alignment"`
	ArmorClass       int                   `json:"armor_class"`
	HitPoints        int                   `json:"hit_points"`
	HitDice          string                `json:"hit_dice"`
	Speed            string                `json:"speed"`
	Strength         int                   `json:"strength"`
	Dexterity        int                   `json:"dexterity"`
	Constitution     int                   `json:"constitution"`
	Intelligence     int                   `json:"intelligence"`
	Wisdom           int                   `json:"wisdom"`
	Charisma         int                   `json:"charisma"`
	ChallengeRating  string                `json:"challenge_rating"`
	Abilities        []string              `json:"abilities,omitempty"`
	Actions          []string              `json:"actions,omitempty"`
	LegendaryActions []string              `json:"legendary_actions,omitempty"`
	Description      string                `json:"description"`
	Attacks          []CreateAttackRequest `json:"attacks,omitempty"`
}

type UpdateBeastRequest struct {
	Name             *string              `json:"name,omitempty"`
	ImageURL         *string              `json:"image_url,omitempty"`
	Size             *models.CreatureSize `json:"size,omitempty"`
	Type             *models.CreatureType `json:"type,omitempty"`
	Alignment        *string              `json:"alignment,omitempty"`
	ArmorClass       *int                 `json:"armor_class,omitempty"`
	HitPoints        *int                 `json:"hit_points,omitempty"`
	HitDice          *string              `json:"hit_dice,omitempty"`
	Speed            *string              `json:"speed,omitempty"`
	Strength         *int                 `json:"strength,omitempty"`
	Dexterity        *int                 `json:"dexterity,omitempty"`
	Constitution     *int                 `json:"constitution,omitempty"`
	Intelligence     *int                 `json:"intelligence,omitempty"`
	Wisdom           *int                 `json:"wisdom,omitempty"`
	Charisma         *int                 `json:"charisma,omitempty"`
	ChallengeRating  *string              `json:"challenge_rating,omitempty"`
	Abilities        []string             `json:"abilities,omitempty"`
	Actions          []string             `json:"actions,omitempty"`
	LegendaryActions []string             `json:"legendary_actions,omitempty"`
	Description      *string              `json:"description,omitempty"`
}

// Attack Request/Response Types

type CreateAttackRequest struct {
	Name        string            `json:"name"`
	AttackBonus int               `json:"attack_bonus"`
	DamageType  models.DamageType `json:"damage_type"`
	DamageDice  string            `json:"damage_dice"`
	Reach       *string           `json:"reach,omitempty"`
	Description *string           `json:"description,omitempty"`
}

// HP Management Request/Response Types

// UpdateHPRequest is used to apply damage or healing
type UpdateHPRequest struct {
	Delta    int        `json:"delta"` // Positive = heal, negative = damage
	Source   *string    `json:"source,omitempty"`
	TargetID *uuid.UUID `json:"target_id,omitempty"` // ID of the beast being attacked (for Lorewright's Predator's Strike)
	WeaponID *uuid.UUID `json:"weapon_id,omitempty"` // ID of the weapon used (for damage calculations)
}

// SetTempHPRequest sets temporary HP
type SetTempHPRequest struct {
	TempHP int `json:"temp_hp"`
}

// UseHitDiceRequest contains the rolled values for hit dice
type UseHitDiceRequest struct {
	Rolls []int `json:"rolls"` // Array of roll results (d{HitDie} + ConMod each)
}

// UseHitDiceResponse returns the result of using hit dice
type UseHitDiceResponse struct {
	CurrentHP        int   `json:"current_hp"`
	MaxHP            int   `json:"max_hp"`
	HPHealed         int   `json:"hp_healed"`
	DiceUsed         int   `json:"dice_used"`
	DiceResults      []int `json:"dice_results"`
	HitDiceRemaining int   `json:"hit_dice_remaining"`
}

// RestResponse returns HP state after a rest
type RestResponse struct {
	CurrentHP          int `json:"current_hp"`
	MaxHP              int `json:"max_hp"`
	TempHP             int `json:"temp_hp"`
	CurrentSpellPoints int `json:"current_spell_points"`
	MaxSpellPoints     int `json:"max_spell_points"`
	HitDiceRemaining   int `json:"hit_dice_remaining"`
	HitDiceTotal       int `json:"hit_dice_total"`

	// Generic class resources (updated after rest)
	ClassResources []ClassResourceResponse `json:"class_resources,omitempty"`
}

// Mechanics Request/Response Types

type RollTableRequest struct {
	TableName string `json:"table_name"` // "mutagen_feral", "lorewright_madness"
}

type RollTableResponse struct {
	Effect *models.Effect `json:"effect"`
	Roll   int            `json:"roll"`
}

type MutagenCastRequest struct {
	ComponentCount int `json:"component_count"`
}

type MutagenCastResponse struct {
	MadnessCastCount int  `json:"madness_cast_count"`
	CurrentDC        int  `json:"current_dc"`
	RequiresSave     bool `json:"requires_save"`
	SaveDC           int  `json:"save_dc"`
}

type ActiveEffectResponse struct {
	ID          uuid.UUID `json:"id"`
	EffectID    uuid.UUID `json:"effect_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Mechanics   string    `json:"mechanics"`
	Duration    string    `json:"duration"`
	Source      string    `json:"source"`
	Category    string    `json:"category"`
}

// Character Effect Management Request/Response Types

type ApplyEffectRequest struct {
	EffectID          uuid.UUID  `json:"effect_id"`
	Stacks            int        `json:"stacks,omitempty"`
	MaxStacks         int        `json:"max_stacks,omitempty"`
	DurationRounds    *int       `json:"duration_rounds,omitempty"`
	DurationMinutes   *int       `json:"duration_minutes,omitempty"`
	Source            string     `json:"source,omitempty"`
	SourceCharacterID *uuid.UUID `json:"source_character_id,omitempty"`
	SourceSpellID     *uuid.UUID `json:"source_spell_id,omitempty"`
	IsConcentration   bool       `json:"is_concentration,omitempty"`
	StackIfExists     bool       `json:"stack_if_exists,omitempty"`
}

type ModifyStacksRequest struct {
	Delta int `json:"delta"` // Positive = add stacks, negative = remove stacks
}

type RemoveEffectRequest struct {
	RemoveAllStacks bool `json:"remove_all_stacks,omitempty"`
	StacksToRemove  int  `json:"stacks_to_remove,omitempty"`
}

type TickDurationRequest struct {
	Rounds int `json:"rounds,omitempty"` // Default 1
}

type TickDurationResponse struct {
	Expired []ExpiredEffectInfo `json:"expired"`
	Rounds  int                 `json:"rounds"`
}

type ExpiredEffectInfo struct {
	EffectID   uuid.UUID `json:"effect_id"`
	EffectName string    `json:"effect_name"`
	Reason     string    `json:"reason"`
}

type BreakConcentrationResponse struct {
	Broken []ExpiredEffectInfo `json:"broken"`
}

type CharacterEffectResponse struct {
	ID                uuid.UUID  `json:"id"`
	CharacterID       uuid.UUID  `json:"character_id"`
	EffectID          uuid.UUID  `json:"effect_id"`
	EffectName        string     `json:"effect_name"`
	EffectDescription string     `json:"effect_description"`
	EffectMechanics   string     `json:"effect_mechanics"`
	EffectCategory    string     `json:"effect_category"`
	Duration          string     `json:"duration"`
	Source            string     `json:"source"`
	Stacks            int        `json:"stacks"`
	MaxStacks         int        `json:"max_stacks"`
	DurationRounds    *int       `json:"duration_rounds,omitempty"`
	DurationMinutes   *int       `json:"duration_minutes,omitempty"`
	SourceCharacterID *uuid.UUID `json:"source_character_id,omitempty"`
	IsConcentration   bool       `json:"is_concentration"`
}

// Character Resource Request/Response Types

type CreateResourceRequest struct {
	ResourceKey        string `json:"resource_key"`
	ResourceName       string `json:"resource_name"`
	CurrentValue       int    `json:"current_value"`
	MaxValue           *int   `json:"max_value,omitempty"`
	RestoreOnShortRest bool   `json:"restore_on_short_rest,omitempty"`
	RestoreOnLongRest  bool   `json:"restore_on_long_rest,omitempty"`
	RestoreAmount      *int   `json:"restore_amount,omitempty"`
	DecaysOnLongRest   bool   `json:"decays_on_long_rest,omitempty"`
}

type ResourceDeltaRequest struct {
	Amount int `json:"amount"`
}

// Minion Request/Response Types

type CreateMinionRequest struct {
	MinionType  string `json:"minion_type"` // "construct", "echo", "drone"
	TemplateKey string `json:"template_key,omitempty"`
	CustomName  string `json:"custom_name,omitempty"`
	// Drone-specific: which component to spend
	ComponentID *uuid.UUID `json:"component_id,omitempty"`
	// Echo-specific fields
	SlotIndex          int    `json:"slot_index,omitempty"`
	AbilityName        string `json:"ability_name,omitempty"`
	AbilityDescription string `json:"ability_description,omitempty"`
	AbilityType        string `json:"ability_type,omitempty"`
	SourceCreatureType string `json:"source_creature_type,omitempty"`
}

type UpdateMinionHPRequest struct {
	Delta int `json:"delta"`
}

// Corpse Request/Response Types

type CreateCorpseRequest struct {
	MapID               *uuid.UUID `json:"map_id,omitempty"`
	Name                string     `json:"name"`
	CreatureType        string     `json:"creature_type"`
	CreatureSize        string     `json:"creature_size,omitempty"`
	ChallengeRating     float64    `json:"challenge_rating,omitempty"`
	GridX               *int       `json:"grid_x,omitempty"`
	GridY               *int       `json:"grid_y,omitempty"`
	AvailableComponents []string   `json:"available_components,omitempty"`
	ComponentYield      int        `json:"component_yield,omitempty"`
	SourceBeastID       *uuid.UUID `json:"source_beast_id,omitempty"`
	ExpiresInMinutes    *int       `json:"expires_in_minutes,omitempty"`
}

type HarvestCorpseRequest struct {
	CharacterID uuid.UUID `json:"character_id"`
}

type ConsumeCorpseRequest struct {
	CharacterID uuid.UUID `json:"character_id"`
}

// Initiative Request/Response Types

type InitiativeEntry struct {
	TokenID uuid.UUID `json:"token_id"`
	Order   int       `json:"order"`
}

type SetInitiativeRequest struct {
	Entries []InitiativeEntry `json:"entries"`
}

// API Response Types

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
