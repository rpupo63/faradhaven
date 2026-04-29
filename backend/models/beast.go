package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Beast represents a bestiary entry (creature/monster).
// Enums (CreatureSize, CreatureType, DamageType) are in enums.go.
// Attack sub-model is in attack.go.
type Beast struct {
	ID         uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	UserID     uuid.UUID    `json:"user_id" gorm:"type:uuid;not null;index"`
	Name       string       `json:"name" gorm:"type:text;not null"`
	ImageURL   *string      `json:"image_url" gorm:"type:text"`
	Size       CreatureSize `json:"size" gorm:"type:text;not null"`
	Type       CreatureType `json:"type" gorm:"type:text;not null"`
	Alignment  string       `json:"alignment" gorm:"type:text;not null"`
	ArmorClass int          `json:"armor_class" gorm:"type:int;not null"`
	HitPoints  int          `json:"hit_points" gorm:"type:int;not null"`
	HitDice    string       `json:"hit_dice" gorm:"type:text;not null"`
	Speed      string       `json:"speed" gorm:"type:text;not null"`
	// Ability scores
	Strength     int `json:"strength" gorm:"type:int;not null"`
	Dexterity    int `json:"dexterity" gorm:"type:int;not null"`
	Constitution int `json:"constitution" gorm:"type:int;not null"`
	Intelligence int `json:"intelligence" gorm:"type:int;not null"`
	Wisdom       int `json:"wisdom" gorm:"type:int;not null"`
	Charisma     int `json:"charisma" gorm:"type:int;not null"`
	// Combat
	ChallengeRating  string         `json:"challenge_rating" gorm:"type:text;not null"`
	Abilities        pq.StringArray `json:"abilities" gorm:"type:text[]"`
	Actions          pq.StringArray `json:"actions" gorm:"type:text[]"`
	LegendaryActions pq.StringArray `json:"legendary_actions" gorm:"type:text[]"`
	Description      string         `json:"description" gorm:"type:text"`
	CreatedAt        time.Time      `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User    User         `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Attacks []Attack     `json:"attacks,omitempty" gorm:"foreignKey:BeastID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Skills  []BeastSkill `json:"skills,omitempty" gorm:"foreignKey:BeastID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Parties []Party      `json:"-" gorm:"many2many:party_beasts;references:ID;joinReferences:PartyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // NEW: Many2Many with Party
}

// ParseChallengeRating converts the string ChallengeRating to a float64.
// Handles formats like "1/2", "5", "20".
func (b *Beast) ParseChallengeRating() float64 {
	if b.ChallengeRating == "" {
		return 0.0
	}

	// Try to parse as a fraction
	parts := strings.Split(b.ChallengeRating, "/")
	if len(parts) == 2 {
		numerator, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0.0 // Or handle error appropriately
		}
		denominator, err := strconv.ParseFloat(parts[1], 64)
		if err != nil || denominator == 0 {
			return 0.0 // Or handle error appropriately
		}
		return numerator / denominator
	}

	// Try to parse as a whole number
	cr, err := strconv.ParseFloat(b.ChallengeRating, 64)
	if err != nil {
		return 0.0 // Or handle error appropriately
	}
	return cr
}
