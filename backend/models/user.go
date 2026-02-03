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

	// Relationships (HasMany: children hold FK; constraint ensures ON DELETE CASCADE)
	Characters []Character `json:"characters,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Spells     []Spell     `json:"spells,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Beasts     []Beast     `json:"beasts,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}
