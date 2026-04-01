package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account for authentication and game data ownership
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Name         string    `json:"name" gorm:"type:text;not null"`
	Email        string    `json:"email" gorm:"type:text;not null;unique"`
	PasswordHash string    `json:"-" gorm:"type:text;not null"`    // bcrypt hash; checked at login
	Token        *string   `json:"-" gorm:"type:text;uniqueIndex"` // session token; checked on each request (Bearer)
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	IsAdmin bool `json:"is_admin" gorm:"type:boolean;default:false"` // Added IsAdmin field

	// Active character for shop purchases and quick access
	ActiveCharacterID *uuid.UUID `json:"active_character_id,omitempty" gorm:"type:uuid;index"`

	// Dice appearance preferences (global defaults for this user)
	DiceTheme      string `json:"dice_theme" gorm:"type:text;default:'default'"`
	DiceThemeColor string `json:"dice_theme_color" gorm:"type:text;default:'#7A201C'"`
	DiceFontColor  string `json:"dice_font_color" gorm:"type:text;default:'#B8860B'"`

	// Relationships (HasMany: children hold FK; constraint ensures ON DELETE CASCADE)
	Characters      []Character `json:"characters,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Spells          []Spell     `json:"spells,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Beasts          []Beast     `json:"beasts,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	ActiveCharacter *Character  `json:"active_character,omitempty" gorm:"foreignKey:ActiveCharacterID;references:ID"`
}
