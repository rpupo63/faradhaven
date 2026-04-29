package models

import (
	"time"

	"github.com/google/uuid"
)

// ElevationProperties stores specific properties for a MapElement of Type "elevation".
type ElevationProperties struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	MapElementID uuid.UUID `json:"map_element_id" gorm:"type:uuid;not null;uniqueIndex"` // One-to-one with MapElement
	Level        int       `json:"level" gorm:"type:int;not null"`                       // Elevation level
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	MapElement MapElement `json:"-" gorm:"foreignKey:MapElementID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM
func (ElevationProperties) TableName() string {
	return "map_element_elevation_properties"
}
