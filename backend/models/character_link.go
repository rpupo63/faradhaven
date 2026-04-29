package models

import (
	"time"

	"github.com/google/uuid"
)

// LinkType represents the type of bond between characters
type LinkType string

const (
	LinkTypeRenfield      LinkType = "renfield"      // Sanguinist bonded ally
	LinkTypeConsciousness LinkType = "consciousness" // Lorewright collective consciousness
	LinkTypeEther         LinkType = "ether"         // Generic magical bond
	LinkTypeFamiliar      LinkType = "familiar"      // Familiar bond
)

// CharacterLink represents a bond between two characters.
// Used for Sanguinist's Renfield devotion, Lorewright's collective consciousness, etc.
type CharacterLink struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	SourceCharacterID uuid.UUID `json:"source_character_id" gorm:"type:uuid;not null;index"`
	TargetCharacterID uuid.UUID `json:"target_character_id" gorm:"type:uuid;not null;index"`

	LinkType LinkType `json:"link_type" gorm:"type:text;not null"`
	IsActive bool     `json:"is_active" gorm:"default:true"`

	// Optional expiration
	ExpiresAt *time.Time `json:"expires_at,omitempty" gorm:"type:timestamptz"`

	// Link-specific data
	SharedEffectID *uuid.UUID `json:"shared_effect_id,omitempty" gorm:"type:uuid"` // For shared buffs/echoes
	BonusType      *string    `json:"bonus_type,omitempty" gorm:"type:text"`       // e.g., "advantage_on_attacks"
	Notes          *string    `json:"notes,omitempty" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	SourceCharacter Character `json:"-" gorm:"foreignKey:SourceCharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TargetCharacter Character `json:"-" gorm:"foreignKey:TargetCharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM
func (CharacterLink) TableName() string {
	return "character_links"
}

// IsExpired checks if the link has expired
func (l *CharacterLink) IsExpired() bool {
	if l.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*l.ExpiresAt)
}
