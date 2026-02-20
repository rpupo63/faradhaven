package models

import (
	"github.com/google/uuid"
)

type MapElementType string

const (
	MapElementTypeWall             MapElementType = "wall"
	MapElementTypeDifficultTerrain MapElementType = "difficult_terrain"
	MapElementTypeTrap             MapElementType = "trap"
	MapElementTypeElevation        MapElementType = "elevation"
)

type MapElement struct {
	ID    uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MapID uuid.UUID      `gorm:"type:uuid;not null;index" json:"map_id"`
	Type  MapElementType `gorm:"type:varchar(50);not null" json:"type"`
	GridX int            `gorm:"not null" json:"grid_x"`
	GridY int            `gorm:"not null" json:"grid_y"`

	Map GameMap `gorm:"foreignKey:MapID;constraint:OnDelete:CASCADE" json:"-"`

	// Has-one relationships via MapElementID on each property table
	TrapProperties             *TrapProperties             `json:"trap_properties,omitempty" gorm:"foreignKey:MapElementID"`
	DifficultTerrainProperties *DifficultTerrainProperties `json:"difficult_terrain_properties,omitempty" gorm:"foreignKey:MapElementID"`
	ElevationProperties        *ElevationProperties        `json:"elevation_properties,omitempty" gorm:"foreignKey:MapElementID"`
	WallProperties             *WallProperties             `json:"wall_properties,omitempty" gorm:"foreignKey:MapElementID"`
}
