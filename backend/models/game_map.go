package models

import (
	"github.com/google/uuid"
)

type GameMap struct {
	ID            uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OwnerID       uuid.UUID    `gorm:"type:uuid;not null" json:"owner_id"`
	RoomCode      string       `gorm:"type:varchar(50);unique;not null" json:"room_code"`
	Name          string       `gorm:"type:varchar(255);not null" json:"name"`
	// BackgroundURL is the HTTPS URL of the map image (S3 maps/backgrounds/ when uploaded via API, or any image URL).
	BackgroundURL string       `gorm:"type:text" json:"background_url"`
	GridRows      int          `gorm:"default:20" json:"grid_rows"`
	GridCols      int          `gorm:"default:20" json:"grid_cols"`
	TileSize      int          `gorm:"default:50" json:"tile_size"`
	Tokens        []MapToken   `gorm:"foreignKey:MapID;constraint:OnDelete:CASCADE" json:"tokens"`
	Elements      []MapElement `gorm:"foreignKey:MapID;constraint:OnDelete:CASCADE" json:"elements"`

	Owner User `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"-"`
}
