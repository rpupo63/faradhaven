package models

import (
	"time"

	"github.com/google/uuid"
)

// ClassLevelResource stores the progression value of a resource at a specific class level.
// This replaces the Faradhaven-specific columns on ClassLevel (e.g., ConcurrencyLimit, MaxStability).
// One row per resource per level (e.g., Ironwright level 5: concurrency_limit=2, yield_die=6).
type ClassLevelResource struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ClassLevelID uuid.UUID `json:"class_level_id" gorm:"type:uuid;not null;index:idx_class_level_resource,unique"`
	ResourceKey  string    `json:"resource_key" gorm:"type:text;not null;index:idx_class_level_resource,unique"` // matches ClassResourceDefinition.ResourceKey
	Value        int       `json:"value" gorm:"type:int;not null"`                                               // the progression value at this level

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	ClassLevel ClassLevel `json:"-" gorm:"foreignKey:ClassLevelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
