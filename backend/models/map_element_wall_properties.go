package models

import (
	"time"

	"github.com/google/uuid"
)

// WallProperties stores specific properties for a MapElement of Type "wall".
type WallProperties struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	MapElementID uuid.UUID `json:"map_element_id" gorm:"type:uuid;not null;uniqueIndex"` // One-to-one with MapElement
	// Add wall-specific properties here. As per frontend, it might be an x_coords (start x of wall segment).
	// For simplicity, let's just add a placeholder for now if no specific properties are defined beyond geometry.
	// Frontend suggests: setMapElementProperties({ x_coords: 0 }); // Placeholder for wall specific properties, e.g. direction
	// If walls are represented as segments, this needs more thought. For now, a single "Direction" might suffice
	// or more complex geometry data. Let's start simple with a text field for "details" if no specific fields are known.
	Details   string    `json:"details" gorm:"type:text"` // e.g., "North-South segment", "Blocking line"
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	MapElement MapElement `json:"-" gorm:"foreignKey:MapElementID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM
func (WallProperties) TableName() string {
	return "map_element_wall_properties"
}
