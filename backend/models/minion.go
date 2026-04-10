package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// MinionType represents the type of minion
type MinionType string

const (
	MinionTypeConstruct MinionType = "construct" // Ironwright constructs
	MinionTypeEcho      MinionType = "echo"      // Lorewright stored abilities (legacy)
	MinionTypeHarvest   MinionType = "harvest"   // Lorewright harvested abilities
	MinionTypeRenfield  MinionType = "renfield"  // Sanguinist bonded ally
	MinionTypeSummon    MinionType = "summon"    // Generic summoned creature
	MinionTypeDrone     MinionType = "drone"     // Swarm Engineer micro-constructs
)

// Minion represents a summoned or created entity controlled by a character.
// This includes Ironwright constructs, Lorewright echoes, and other minion types.
type Minion struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	CharacterID uuid.UUID  `json:"character_id" gorm:"type:uuid;not null;index"`
	Name        string     `json:"name" gorm:"type:text;not null"`
	MinionType  MinionType `json:"minion_type" gorm:"type:text;not null"`

	// Combat stats (primarily for constructs/summons)
	ArmorClass int `json:"armor_class" gorm:"type:int;default:10"`
	MaxHP      int `json:"max_hp" gorm:"type:int;default:1"`
	CurrentHP  int `json:"current_hp" gorm:"type:int;default:1"`
	Speed      int `json:"speed" gorm:"type:int;default:30"` // feet per round

	// Construct-specific fields (Ironwright)
	ComponentCost int     `json:"component_cost" gorm:"type:int;default:0"`
	TemplateKey   *string `json:"template_key,omitempty" gorm:"type:text"` // "sentry", "striker", "titan"

	// Echo-specific fields (Lorewright)
	AbilityName        *string `json:"ability_name,omitempty" gorm:"type:text"`
	AbilityDescription *string `json:"ability_description,omitempty" gorm:"type:text"`
	AbilityType        *string `json:"ability_type,omitempty" gorm:"type:text"` // "skill", "action", "trait", "language"
	SourceCreatureType *string `json:"source_creature_type,omitempty" gorm:"type:text"`
	SlotIndex          *int    `json:"slot_index,omitempty" gorm:"type:int"` // Which echo slot this occupies

	// State
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" gorm:"type:timestamptz"`

	// Extensible metadata (for actions, abilities, etc.)
	Metadata datatypes.JSON `json:"metadata,omitempty" gorm:"type:jsonb"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Character Character `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM
func (Minion) TableName() string {
	return "minions"
}

// ConstructTemplate defines the base stats for construct types
type ConstructTemplate struct {
	Key                string
	Name               string
	Description        string
	BaseAC             int
	BaseHP             int // Added to level-based calculation
	BaseSpeed          int
	Size               string // "Small", "Medium", "Large", "Tiny"
	Actions            []ConstructAction
	RequiredComponents []string `json:"required_components"` // Component names needed to build
}

// ConstructAction represents an action a construct can take
type ConstructAction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DamageDice  string `json:"damage_dice,omitempty"`
	DamageType  string `json:"damage_type,omitempty"`
	Range       string `json:"range,omitempty"`
	ActionType  string `json:"action_type"` // "action", "bonus_action", "reaction"
}

// GetConstructTemplates returns the available construct templates for Ironwright
func GetConstructTemplates() map[string]ConstructTemplate {
	return map[string]ConstructTemplate{
		"sentry": {
			Key:                "sentry",
			Name:               "Sentry Construct",
			Description:        "A ranged construct that provides covering fire.",
			BaseAC:             12,
			BaseHP:             5,
			BaseSpeed:          20,
			Size:               "Small",
			RequiredComponents: []string{"Fulgur", "Ferrum"},
			Actions: []ConstructAction{
				{
					Name:        "Bolt Shot",
					Description: "Ranged attack against a target within 60 feet.",
					DamageDice:  "1d8",
					DamageType:  "Piercing",
					Range:       "60 ft",
					ActionType:  "action",
				},
			},
		},
		"striker": {
			Key:                "striker",
			Name:               "Striker Construct",
			Description:        "A melee construct built for close combat.",
			BaseAC:             14,
			BaseHP:             8,
			BaseSpeed:          30,
			Size:               "Medium",
			RequiredComponents: []string{"Ferrum", "Ignis", "Crush", "Push", "Arcanum"},
			Actions: []ConstructAction{
				{
					Name:        "Slam",
					Description: "Melee attack against an adjacent target.",
					DamageDice:  "1d10",
					DamageType:  "Bludgeoning",
					Range:       "5 ft",
					ActionType:  "action",
				},
			},
		},
		"titan": {
			Key:                "titan",
			Name:               "Titan Construct",
			Description:        "A large construct that dominates the battlefield.",
			BaseAC:             16,
			BaseHP:             15,
			BaseSpeed:          25,
			Size:               "Large",
			RequiredComponents: []string{"Fulgur", "Ferrum", "Ignis", "Push", "Crush", "Zone", "Arcanum", "Terra", "Vita", "Ward"},
			Actions: []ConstructAction{
				{
					Name:        "Crushing Blow",
					Description: "Melee attack against an adjacent target.",
					DamageDice:  "2d8",
					DamageType:  "Bludgeoning",
					Range:       "10 ft",
					ActionType:  "action",
				},
				{
					Name:        "Fortress Stance",
					Description: "Grants half cover to allies within 5 feet.",
					ActionType:  "bonus_action",
				},
			},
		},
		"drone": {
			Key:                "drone",
			Name:               "Drone",
			Description:        "A tiny expendable micro-construct (Swarm Engineer only).",
			BaseAC:             12,
			BaseHP:             1,
			BaseSpeed:          30,
			Size:               "Tiny",
			RequiredComponents: []string{}, // Any 1 component
			Actions: []ConstructAction{
				{
					Name:        "Force Strike",
					Description: "Ranged attack against a target within 30 feet.",
					DamageDice:  "1d4",
					DamageType:  "Force",
					Range:       "30 ft",
					ActionType:  "action",
				},
			},
		},
	}
}

// CalculateConstructHP calculates the HP for a construct based on level
// Formula: BaseHP + (2 * Level)
func CalculateConstructHP(template ConstructTemplate, level int) int {
	return template.BaseHP + (2 * level)
}
