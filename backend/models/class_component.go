package models

import (
	"time"

	"github.com/google/uuid"
)

// ClassComponent links a class to the components available to it.
// A class may have many components, and a component may be available to many classes.
type ClassComponent struct {
	ClassID     uuid.UUID `json:"class_id" gorm:"type:uuid;primaryKey;not null"`
	ComponentID uuid.UUID `json:"component_id" gorm:"type:uuid;primaryKey;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	Class     Class     `json:"-" gorm:"foreignKey:ClassID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Component Component `json:"-" gorm:"foreignKey:ComponentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
