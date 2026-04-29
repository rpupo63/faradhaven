package models

import (
	"time"

	"github.com/google/uuid"
)

// ClassResourceDefinition defines a resource type for a class.
// This replaces the single ResourceType/ResourceName/ResourceRestoreType fields on Class,
// allowing each class to declare multiple resources with their metadata and behavior.
type ClassResourceDefinition struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ClassID uuid.UUID `json:"class_id" gorm:"type:uuid;not null;index:idx_class_resource_def,unique"`

	// Resource identification
	ResourceKey string `json:"resource_key" gorm:"type:text;not null;index:idx_class_resource_def,unique"` // e.g., "concurrency_limit", "max_stability"
	DisplayName string `json:"display_name" gorm:"type:text;not null"`                                     // e.g., "Concurrency Limit", "Stability"

	// Rendering/category metadata
	Category    ClassResourceCategory `json:"category" gorm:"type:text;not null"`
	Description string `json:"description" gorm:"type:text"`       // tooltip/explanation text

	// UI ordering
	DisplayOrder int `json:"display_order" gorm:"type:int;default:0"` // controls rendering order on character sheet

	// Behavior configuration
	IsTrackable        bool `json:"is_trackable" gorm:"default:false"`          // true if character has mutable current value (pools, counters)
	RestoreOnShortRest bool `json:"restore_on_short_rest" gorm:"default:false"` // restore current_value on short rest
	RestoreOnLongRest  bool `json:"restore_on_long_rest" gorm:"default:false"`  // restore current_value on long rest

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Class Class `json:"-" gorm:"foreignKey:ClassID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
