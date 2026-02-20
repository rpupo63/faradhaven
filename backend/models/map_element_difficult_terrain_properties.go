package models

import (
	"time"

	"github.com/google/uuid"
)

// DifficultTerrainProperties stores specific properties for a MapElement of Type "difficult_terrain".
type DifficultTerrainProperties struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	MapElementID uuid.UUID `json:"map_element_id" gorm:"type:uuid;not null;uniqueIndex"` // One-to-one with MapElement
	Color        string    `json:"color" gorm:"type:text;not null"`                     // e.g., "#6B8E2380" for transparent olive green
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	MapElement MapElement `json:"-" gorm:"foreignKey:MapElementID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM
func (DifficultTerrainProperties) TableName() string {
	return "map_element_difficult_terrain_properties"
}
