package models

import (
	"github.com/google/uuid"
)

type GameMap struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OwnerID       uuid.UUID  `gorm:"type:uuid;not null" json:"owner_id"`
	RoomCode      string     `gorm:"type:varchar(50);unique;not null" json:"room_code"`
	Name          string     `gorm:"type:varchar(255);not null" json:"name"`
	BackgroundURL string     `gorm:"type:text" json:"background_url"`
	GridRows      int        `gorm:"default:20" json:"grid_rows"`
	GridCols      int        `gorm:"default:20" json:"grid_cols"`
	TileSize      int        `gorm:"default:50" json:"tile_size"`
	Tokens        []MapToken `gorm:"foreignKey:MapID;constraint:OnDelete:CASCADE" json:"tokens"`
}

type MapToken struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MapID          uuid.UUID  `gorm:"type:uuid;not null" json:"map_id"`
	CharacterID    *uuid.UUID `gorm:"type:uuid" json:"character_id"`
	AssignedUserID *uuid.UUID `gorm:"type:uuid" json:"assigned_user_id"`
	Name           string     `gorm:"type:varchar(255);not null" json:"name"`
	ImageURL       string     `gorm:"type:text" json:"image_url"`
	TokenType      string     `gorm:"type:varchar(50);default:'npc'" json:"token_type"` // "pc" or "npc"
	GridX          int        `gorm:"not null" json:"grid_x"`
	GridY          int        `gorm:"not null" json:"grid_y"`
	Size           int        `gorm:"default:1" json:"size"`
	Color          string     `gorm:"type:varchar(50)" json:"color"`
	Visible        bool       `gorm:"default:true" json:"visible"`
	InitiativeOrder *int      `gorm:"type:int" json:"initiative_order,omitempty"`
}