package models

import (
	"time"

	"github.com/google/uuid"
)

// Party represents a group of characters.
type Party struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Name      string    `json:"name" gorm:"type:text;not null"`
	OwnerID   uuid.UUID `json:"owner_id" gorm:"type:uuid;not null;index"` // The user who created/manages the party
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Owner User `json:"-" gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE"`

	// Many2Many for IdentifiedBeasts (through PartyBeast)
	IdentifiedBeasts []*Beast `json:"identified_beasts" gorm:"many2many:party_beasts;references:ID;joinReferences:BeastID"`
	// Has many for Members
	Members []*Character `json:"members" gorm:"foreignKey:PartyID"`
}

// PartyBeast is the join table for Party and Beast to track which beasts a party has identified.
type PartyBeast struct {
	PartyID   uuid.UUID `json:"party_id" gorm:"type:uuid;primaryKey;not null"`
	BeastID   uuid.UUID `json:"beast_id" gorm:"type:uuid;primaryKey;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
}
