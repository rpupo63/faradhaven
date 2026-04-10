package api

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	gameMapHandler         *gameMapHandler
	mapTokenHandler        *mapTokenHandler
	mapElementHandler      *mapElementHandler
	mechanicsHandler       *MechanicsHandler
	corpseHandler          *corpseHandler
	linkHandler            *linkHandler
	harvestHandler         *harvestHandler
	madnessHandler         *madnessHandler
	lootHandler            *LootHandler
	monsterHandler         *monsterHandler
	partyHandler           *partyHandler // NEW: Party Handler
	abilityHandler         *abilityHandler
	storeOwnerHandler      *storeOwnerHandler
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
	CharacterID    *uuid.UUID          `json:"character_id,omitempty"`
	MonsterID      *uuid.UUID          `json:"monster_id,omitempty"`
	AssignedUserID *uuid.UUID          `json:"assigned_user_id,omitempty"`
	Name           string              `json:"name"`
	ImageURL       string              `json:"image_url"`
	TokenType      models.MapTokenType `json:"token_type"` // "pc" or "npc"; empty defaults to npc
	GridX          int                 `json:"grid_x"`
	GridY          int                 `json:"grid_y"`
	Size           int                 `json:"size"`
	Color          string              `json:"color"`
	Visible        bool                `json:"visible"`
}

type UpdateTokenRequest struct {
	GridX          *int       `json:"grid_x,omitempty"`
	GridY          *int       `json:"grid_y,omitempty"`
	Visible        *bool      `json:"visible,omitempty"`
	Size           *int       `json:"size,omitempty"`
	Color          *string    `json:"color,omitempty"`
	AssignedUserID *uuid.UUID `json:"assigned_user_id,omitempty"` // For DM to reassign
	CharacterID    *uuid.UUID `json:"character_id,omitempty"`     // DM re-link to character
	MonsterID      *uuid.UUID `json:"monster_id,omitempty"`       // DM re-link to monster
	Name           *string    `json:"name,omitempty"`
	ImageURL       *string    `json:"image_url,omitempty"`
	TokenType      *models.MapTokenType `json:"token_type,omitempty"`
}

type CreateMapElementRequest struct {
	Type  models.MapElementType `json:"type"`
	GridX int                   `json:"grid_x"`
	GridY int                   `json:"grid_y"`

	// Specific properties for different element types (nullable)
	TrapName              *string `json:"trap_name,omitempty"`
	TrapDescription       *string `json:"trap_description,omitempty"`
	TrapVisibleToPlayers  *bool   `json:"trap_visible_to_players,omitempty"`
	DifficultTerrainColor *string `json:"difficult_terrain_color,omitempty"`
	ElevationLevel        *int    `json:"elevation_level,omitempty"`
	WallDetails           *string `json:"wall_details,omitempty"`
}

type CreateMultipleMapElementsRequest struct {
	Elements []CreateMapElementRequest `json:"elements"`
}

type SetInitiativeRequest struct {
	Entries []InitiativeEntry `json:"entries"`
}

type InitiativeEntry struct {
	TokenID uuid.UUID `json:"token_id"`
	Order   int       `json:"order"`
}

type GenerateLootRequest struct {
	CharacterID uuid.UUID `json:"character_id"`
	Source      string    `json:"source"`
	Tier        string    `json:"tier"`
}

// Party Types
type CreatePartyRequest struct {
	Name        string     `json:"name"`
	CharacterID *uuid.UUID `json:"character_id,omitempty"`
}

type UpdatePartyRequest struct {
	Name *string `json:"name,omitempty"`
}

type AddCharacterToPartyRequest struct {
	CharacterID uuid.UUID `json:"character_id"`
}

type AddIdentifiedBeastRequest struct {
	BeastID uuid.UUID `json:"beast_id"`
}

type SetCharacterPartyRequest struct {
	PartyID *uuid.UUID `json:"party_id,omitempty"`
}

// Character Creation Types

type CreationOptionsResponse struct {
	Races     []models.Race  `json:"races"`
	Classes   []models.Class `json:"classes"`
	PointsMax int            `json:"points_max"`
}

type CreateCharacterRequest struct {
	UserID             uuid.UUID   `json:"user_id"`
	Name               string      `json:"name"`
	RaceID             uuid.UUID   `json:"race_id"`
	LineageID          *uuid.UUID  `json:"lineage_id,omitempty"`
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
	Languages          []string    `json:"languages"`
	SkillProficiencies []string    `json:"skill_proficiencies"`
	EquipmentChoices   []uuid.UUID `json:"equipment_choices"`
	PrimaryWeaponID    *uuid.UUID  `json:"primary_weapon_id,omitempty"`
}

// Class/Compendium Types

type ClassWithLevelsResponse struct {
	models.Class       `json:",inline"` // Embed models.Class directly
	Levels       []models.ClassLevel `json:"levels"`
	MadnessTable map[int]string      `json:"madness_table,omitempty"`
}

type ClassResourceResponse struct {
	Key          string `json:"key"`
	DisplayName  string `json:"display_name"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	Value        int    `json:"value"`
	IsTrackable  bool   `json:"is_trackable"`
	DisplayOrder int    `json:"display_order"`
	CurrentValue *int   `json:"current_value,omitempty"`
	MaxValue     *int   `json:"max_value,omitempty"`
}

// Beast Types

type CreateBeastRequest struct {
	UserID           uuid.UUID           `json:"user_id"`
	Name             string              `json:"name"`
	ImageURL         *string             `json:"image_url,omitempty"`
	Size             models.CreatureSize `json:"size"`
	Type             models.CreatureType `json:"type"`
	Alignment        string              `json:"alignment"`
	ArmorClass       int                 `json:"armor_class"`
	HitPoints        int                 `json:"hit_points"`
	HitDice          string              `json:"hit_dice"`
	Speed            string              `json:"speed"`
	Strength         int                 `json:"strength"`
	Dexterity        int                 `json:"dexterity"`
	Constitution     int                 `json:"constitution"`
	Intelligence     int                 `json:"intelligence"`
	Wisdom           int                 `json:"wisdom"`
	Charisma         int                 `json:"charisma"`
	ChallengeRating  string              `json:"challenge_rating"`
	Abilities        pq.StringArray      `json:"abilities"`
	Actions          pq.StringArray      `json:"actions"`
	LegendaryActions pq.StringArray      `json:"legendary_actions"`
	Description      string              `json:"description"`
	Attacks          []AttackRequest     `json:"attacks"`
}

type AttackRequest struct {
	Name        string            `json:"name"`
	AttackBonus int               `json:"attack_bonus"`
	DamageType  models.DamageType `json:"damage_type"`
	DamageDice  string            `json:"damage_dice"`
	Reach       *string           `json:"reach,omitempty"`
	Description *string           `json:"description,omitempty"`
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
	Abilities        pq.StringArray       `json:"abilities,omitempty"`
	Actions          pq.StringArray       `json:"actions,omitempty"`
	LegendaryActions pq.StringArray       `json:"legendary_actions,omitempty"`
	Description      *string              `json:"description,omitempty"`
}

// Character Effect Types

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

type ApplyEffectRequest struct {
	EffectID          uuid.UUID  `json:"effect_id"`
	Stacks            int        `json:"stacks"`
	MaxStacks         int        `json:"max_stacks"`
	DurationRounds    *int       `json:"duration_rounds,omitempty"`
	DurationMinutes   *int       `json:"duration_minutes,omitempty"`
	Source            string     `json:"source"`
	SourceCharacterID *uuid.UUID `json:"source_character_id,omitempty"`
	SourceSpellID     *uuid.UUID `json:"source_spell_id,omitempty"`
	IsConcentration   bool       `json:"is_concentration"`
	StackIfExists     bool       `json:"stack_if_exists"`
}

type ModifyStacksRequest struct {
	Delta int `json:"delta"`
}

type RemoveEffectRequest struct {
	RemoveAllStacks bool `json:"remove_all_stacks"`
	StacksToRemove  int  `json:"stacks_to_remove"`
}

type TickDurationRequest struct {
	Rounds int `json:"rounds"`
}

type ExpiredEffectInfo struct {
	EffectID   uuid.UUID `json:"effect_id"`
	EffectName string    `json:"effect_name"`
	Reason     string    `json:"reason"`
}

type TickDurationResponse struct {
	Expired []ExpiredEffectInfo `json:"expired"`
	Rounds  int                 `json:"rounds"`
}

type BreakConcentrationResponse struct {
	Broken []ExpiredEffectInfo `json:"broken"`
}

type UpdateMapElementRequest struct {
	Type  *models.MapElementType `json:"type,omitempty"`
	GridX *int                   `json:"grid_x,omitempty"`
	GridY *int                   `json:"grid_y,omitempty"`

	// Specific properties for different element types (nullable)
	TrapName              *string `json:"trap_name,omitempty"`
	TrapDescription       *string `json:"trap_description,omitempty"`
	TrapVisibleToPlayers  *bool   `json:"trap_visible_to_players,omitempty"`
	DifficultTerrainColor *string `json:"difficult_terrain_color,omitempty"`
	ElevationLevel        *int    `json:"elevation_level,omitempty"`
	WallDetails           *string `json:"wall_details,omitempty"`
}

// Character handler request/response types (character_handler.go, character_sheet_handler.go)

// UpdateCharacterRequest is the request body for PATCH /characters/:id
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
	SkillProficiencies []string   `json:"skill_proficiencies,omitempty"`
	DiceTheme          *string    `json:"dice_theme,omitempty"`
	DiceThemeColor     *string    `json:"dice_theme_color,omitempty"`
	DiceFontColor      *string    `json:"dice_font_color,omitempty"`
	ClearDiceTheme     bool       `json:"clear_dice_theme,omitempty"`
	ClearDiceColors    bool       `json:"clear_dice_colors,omitempty"`
}

// UpdateBackstoryRequest is the request body for updating a character's backstory
type UpdateBackstoryRequest struct {
	Backstory         string  `json:"backstory"`
	BackstoryHexColor *string `json:"backstory_hex_color,omitempty"`
}

// PurchaseItemRequest is the request body for purchasing an item or weapon
type PurchaseItemRequest struct {
	ItemType      string     `json:"item_type"` // "weapon" or "item"
	ItemID        uuid.UUID  `json:"item_id"`
	StoreOwnerID  *uuid.UUID `json:"store_owner_id,omitempty"` // optional vendor; applies exchange_rate when valid
}

// CastSpellRequest is the request body for casting a spell
type CastSpellRequest struct {
	SpellID     *uuid.UUID `json:"spell_id,omitempty"`
	SpellLevel  int        `json:"spell_level"`
	ResourceKey *string    `json:"resource_key,omitempty"`
	Components  []string   `json:"components,omitempty"`
}

// ExtractRequest is the request body for Sanguine Extraction (Sanguinist)
type ExtractRequest struct {
	Source string `json:"source"` // e.g. "siphon", "bite"
}

// UpdateNotorietyRequest is the request body for manual notoriety delta update
type UpdateNotorietyRequest struct {
	Delta int `json:"delta"`
}

// UpdateSanguineNotorietyRequest is the request body for updating Sanguine MP/BR points
type UpdateSanguineNotorietyRequest struct {
	MPChange int `json:"mp_change"`
	BRChange int `json:"br_change"`
}

// CharacterWeaponResponse is a weapon in the character sheet inventory
type CharacterWeaponResponse struct {
	CharacterWeaponID string                   `json:"character_weapon_id"`
	Weapon            models.Weapon            `json:"weapon"`
	IsPrimary         bool                     `json:"is_primary"`
	IsEquipped        bool                     `json:"is_equipped"`
	CustomName        *string                  `json:"custom_name,omitempty"`
	ActiveModifiers   []WeaponModifierResponse `json:"active_modifiers"`
}

// WeaponModifierResponse is a modifier on a character weapon in the sheet
type WeaponModifierResponse struct {
	ModifierType string                  `json:"modifier_type"`
	IsActive     bool                    `json:"is_active"`
	BonusDamage  []BonusDamageInfo       `json:"bonus_damage"`
	Metadata     models.ModifierMetadata `json:"metadata,omitempty"`
}

// BonusDamageInfo describes bonus damage (dice + type) for a weapon modifier
type BonusDamageInfo struct {
	Dice       string `json:"dice"`
	DamageType string `json:"damage_type"`
}

// CharacterSheetResponse is the response for GET character sheet
type CharacterSheetResponse struct {
	Character                *models.Character           `json:"character"`
	Class                    *models.Class               `json:"class"`
	ClassLevel               *models.ClassLevel          `json:"class_level"`
	MaxHP                    int                         `json:"max_hp"`
	CurrentHP                int                         `json:"current_hp"`
	TempHP                   int                         `json:"temp_hp"`
	AC                       int                         `json:"ac"`
	SaveDC                   int                         `json:"save_dc"`
	MaxSpellPoints           int                         `json:"max_spell_points"`
	CurrentSpellPoints       int                         `json:"current_spell_points"`
	SavingThrowProficiencies []string                    `json:"saving_throw_proficiencies"`
	// Class + race pool from live Class.Components / Race.Components (not character_components).
	AvailableComponents      []models.Component          `json:"available_components"`
	HitDiceTotal             int                         `json:"hit_dice_total"`
	HitDiceRemaining         int                         `json:"hit_dice_remaining"`
	HitDie                   int                         `json:"hit_die"`
	Money                    int64                       `json:"money"`
	MeleeAttackBonus         int                         `json:"melee_attack_bonus"`
	RangedAttackBonus        int                         `json:"ranged_attack_bonus"`
	RaceTraits               []models.Trait              `json:"race_traits"`
	Lineage                  *models.Lineage             `json:"lineage,omitempty"`
	InventoryWeapons         []CharacterWeaponResponse   `json:"inventory_weapons"`
	InventoryItems           []models.Item               `json:"inventory_items"`
	// Expendable / acquired component counts (character_components). Spell pool is available_components.
	Components               []models.CharacterComponent `json:"components"`
	HarvestedAbilities       models.HarvestedAbilities   `json:"harvested_abilities"`
	ClassResources           []ClassResourceResponse     `json:"class_resources"`
	MadnessTable             map[int]string              `json:"madness_table,omitempty"`
	TraitUseStates           map[string]int              `json:"trait_use_states,omitempty"` // traitID → current uses
	TraitMaxUses             map[string]int              `json:"trait_max_uses,omitempty"`   // traitID → max uses
}

// ConsumeCorpseRequest is the request body for consuming a corpse (Lorewright Visceral Psychometry)
type ConsumeCorpseRequest struct {
	CharacterID uuid.UUID `json:"character_id"`
}

// CreateCorpseRequest is the request body for creating a corpse
type CreateCorpseRequest struct {
	MapID               *uuid.UUID `json:"map_id,omitempty"`
	Name                string     `json:"name"`
	CreatureType        string     `json:"creature_type"` // parsed via [models.ParseCreatureType]
	CreatureSize        string     `json:"creature_size,omitempty"`
	ChallengeRating     float64    `json:"challenge_rating,omitempty"`
	GridX               *int       `json:"grid_x,omitempty"`
	GridY               *int       `json:"grid_y,omitempty"`
	// When empty or all IDs unknown, server picks 1–4 random distinct components from the current catalog.
	AvailableComponents []string `json:"available_components,omitempty"`
	ComponentYield      int        `json:"component_yield,omitempty"`
	SourceBeastID       *uuid.UUID `json:"source_beast_id,omitempty"`
	ExpiresInMinutes    *int       `json:"expires_in_minutes,omitempty"`
}

// HarvestCorpseRequest is the request body for harvesting a corpse (Ironwright Scavenge)
type HarvestCorpseRequest struct {
	CharacterID uuid.UUID `json:"character_id"`
}

// RestResponse is the response for short/long rest endpoints
type RestResponse struct {
	CurrentHP          int                     `json:"current_hp"`
	MaxHP              int                     `json:"max_hp"`
	TempHP             int                     `json:"temp_hp"`
	CurrentSpellPoints int                     `json:"current_spell_points"`
	MaxSpellPoints     int                     `json:"max_spell_points"`
	HitDiceRemaining   int                     `json:"hit_dice_remaining"`
	HitDiceTotal       int                     `json:"hit_dice_total"`
	ClassResources     []ClassResourceResponse `json:"class_resources"`
}

// ScavengeComponentsRequest is the request body for Lorewright component scavenging from a corpse
type ScavengeComponentsRequest struct {
	CharacterID uuid.UUID `json:"character_id"`
}

// SetTempHPRequest is the request body for PUT temp-hp
type SetTempHPRequest struct {
	TempHP int `json:"temp_hp"`
}

// UpdateHPRequest is the request body for PATCH hp
type UpdateHPRequest struct {
	Delta    int        `json:"delta"`
	Source   *string    `json:"source,omitempty"`
	TargetID *uuid.UUID `json:"target_id,omitempty"`
}

// UseHitDiceRequest is the request body for POST hit-dice
type UseHitDiceRequest struct {
	Rolls []int `json:"rolls"`
}

// UseHitDiceResponse is the response for POST hit-dice
type UseHitDiceResponse struct {
	CurrentHP        int   `json:"current_hp"`
	MaxHP            int   `json:"max_hp"`
	HPHealed         int   `json:"hp_healed"`
	DiceUsed         int   `json:"dice_used"`
	DiceResults      []int `json:"dice_results"`
	HitDiceRemaining int   `json:"hit_dice_remaining"`
}

// Mechanics handler types (mechanics_handler.go)

// ActiveEffectResponse is the response for GET active effects
type ActiveEffectResponse struct {
	ID          uuid.UUID `json:"id"`
	EffectID    uuid.UUID `json:"effect_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Mechanics   string    `json:"mechanics"`
	Category    string    `json:"category"`
	Duration    string    `json:"duration"`
	Source      string    `json:"source"`
}

// MutagenCastRequest is the request body for Mutagen cast (madness logic)
type MutagenCastRequest struct {
	ComponentCount int `json:"component_count"`
}

// MutagenCastResponse is the response for Mutagen cast
type MutagenCastResponse struct {
	MadnessCastCount int  `json:"madness_cast_count"`
	CurrentDC        int  `json:"current_dc"`
	RequiresSave     bool `json:"requires_save"`
	SaveDC           int  `json:"save_dc"`
}

// RollTableRequest is the request body for rolling on an effect table
type RollTableRequest struct {
	TableName string `json:"table_name"`
}

// RollTableResponse is the response for rolling on an effect table
type RollTableResponse struct {
	Effect *models.Effect `json:"effect"`
	Roll   int            `json:"roll"`
}

// Minion handler types (minion_handler.go)

// CreateMinionRequest is the request body for creating a minion (construct, drone, or echo)
type CreateMinionRequest struct {
	MinionType         string     `json:"minion_type"` // "construct", "drone", "echo"
	TemplateKey        string     `json:"template_key,omitempty"`
	CustomName         string     `json:"custom_name,omitempty"`
	ComponentID        *uuid.UUID `json:"component_id,omitempty"`
	SlotIndex          int        `json:"slot_index,omitempty"`
	AbilityName        string     `json:"ability_name,omitempty"`
	AbilityDescription string     `json:"ability_description,omitempty"`
	AbilityType        string     `json:"ability_type,omitempty"`
	SourceCreatureType string     `json:"source_creature_type,omitempty"`
}

// UpdateMinionHPRequest is the request body for updating a minion's HP
type UpdateMinionHPRequest struct {
	Delta int `json:"delta"`
}

// Resource handler types (resource_handler.go)

// CreateResourceRequest is the request body for creating a character resource
type CreateResourceRequest struct {
	ResourceKey        string `json:"resource_key"`
	ResourceName       string `json:"resource_name"`
	CurrentValue       int    `json:"current_value"`
	MaxValue           *int   `json:"max_value,omitempty"`
	RestoreOnShortRest bool   `json:"restore_on_short_rest"`
	RestoreOnLongRest  bool   `json:"restore_on_long_rest"`
	RestoreAmount      *int   `json:"restore_amount,omitempty"`
	DecaysOnLongRest   bool   `json:"decays_on_long_rest"`
}

// ResourceDeltaRequest is the request body for spending or gaining a resource
type ResourceDeltaRequest struct {
	Amount int `json:"amount"`
}

// Spell Types

// RetryAIFieldRequest requests a fresh AI recommendation for one specific spell field.
// Set ApplyRecommendation to true to also overwrite the original spell field with the new recommendation.
type RetryAIFieldRequest struct {
	Field               string `json:"field"`                // one of: name, description, type, range, duration, damage_dice, damage_type, save_attr
	ApplyRecommendation bool   `json:"apply_recommendation"` // true = overwrite original field with recommendation
}

type CreateSpellRequest struct {
	UserID             uuid.UUID             `json:"user_id"`
	CharacterID        *uuid.UUID            `json:"character_id,omitempty"`
	Name               string                `json:"name"`
	Description        string                `json:"description"`
	Type               models.SpellType      `json:"type"`
	Range              *int                  `json:"range,omitempty"`
	Duration           *string               `json:"duration,omitempty"`
	Concentration      bool                  `json:"concentration"`
	SaveAttr           *models.SaveAttribute `json:"save_attr,omitempty"`
	DamageDiceCount    *int                  `json:"damage_dice_count,omitempty"`
	DamageDieSize      *int                  `json:"damage_die_size,omitempty"`
	// DamageType is optional; unknown or empty strings are treated as omitted (not all spells deal damage).
	DamageType *string `json:"damage_type,omitempty"`
	AddModifier        bool                  `json:"add_modifier"`
	ComponentIDs       []uuid.UUID           `json:"component_ids"`
}

type UpdateSpellRequest struct {
	Name              *string               `json:"name,omitempty"`
	Description       *string               `json:"description,omitempty"`
	Type              *models.SpellType     `json:"type,omitempty"`
	Range             *int                  `json:"range,omitempty"`
	Duration          *string               `json:"duration,omitempty"`
	Concentration     *bool                 `json:"concentration,omitempty"`
	SaveAttr          *models.SaveAttribute `json:"save_attr,omitempty"`
	DamageDiceCount   *int                  `json:"damage_dice_count,omitempty"`
	DamageDieSize     *int                  `json:"damage_die_size,omitempty"`
	// DamageType is optional; unknown or empty strings clear the field.
	DamageType *string `json:"damage_type,omitempty"`
	AddModifier       *bool                 `json:"add_modifier,omitempty"`
	Checked           *bool                 `json:"checked,omitempty"`
	ComponentIDs      []uuid.UUID           `json:"component_ids,omitempty"`
}

// User Types

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name              *string    `json:"name,omitempty"`
	Email             *string    `json:"email,omitempty"`
	ActiveCharacterID *uuid.UUID `json:"active_character_id,omitempty"`
	DiceTheme         *string    `json:"dice_theme,omitempty"`
	DiceThemeColor    *string    `json:"dice_theme_color,omitempty"`
	DiceFontColor     *string    `json:"dice_font_color,omitempty"`
}

type SetActiveCharacterRequest struct {
	CharacterID *uuid.UUID `json:"character_id"`
}
