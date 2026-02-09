package models

import (
	"time"

	"github.com/google/uuid"
)

// SharedNote represents a community note or news item
type SharedNote struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	UserID      uuid.UUID `json:"userId" gorm:"type:uuid;not null"`
	User        User      `json:"-" gorm:"foreignKey:UserID"`
	Username    string    `json:"username" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"not null;default:now()"`
}
