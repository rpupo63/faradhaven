package models

import (
	"time"

	"github.com/google/uuid"
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
	ChallengeRating  string    `json:"challenge_rating" gorm:"type:text;not null"`
	Abilities        *string   `json:"abilities" gorm:"type:text"`
	Actions          *string   `json:"actions" gorm:"type:text"`
	LegendaryActions *string   `json:"legendary_actions" gorm:"type:text"`
	Description      string    `json:"description" gorm:"type:text"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User    User     `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Attacks []Attack `json:"attacks,omitempty" gorm:"foreignKey:BeastID;references:ID;constraint:OnDelete:CASCADE"`
}
