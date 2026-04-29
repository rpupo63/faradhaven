package models

import (
	"github.com/google/uuid"
)

type MapToken struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MapID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"map_id"`
	CharacterID     *uuid.UUID `gorm:"type:uuid;index" json:"character_id"`
	MonsterID       *uuid.UUID `gorm:"type:uuid;index" json:"monster_id"`
	AssignedUserID  *uuid.UUID `gorm:"type:uuid;index" json:"assigned_user_id"`
	Name            string     `gorm:"type:varchar(255);not null" json:"name"`
	ImageURL        string     `gorm:"type:text" json:"image_url"`
	TokenType       MapTokenType `gorm:"type:varchar(50);default:'npc'" json:"token_type"`
	GridX           int        `gorm:"not null" json:"grid_x"`
	GridY           int        `gorm:"not null" json:"grid_y"`
	Size            int        `gorm:"default:1" json:"size"`
	Color           string     `gorm:"type:varchar(50)" json:"color"`
	Visible         bool       `gorm:"default:true" json:"visible"`
	InitiativeOrder *int       `gorm:"type:int" json:"initiative_order,omitempty"`

	Map          GameMap    `gorm:"foreignKey:MapID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Character    *Character `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Monster      *Monster   `gorm:"foreignKey:MonsterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	AssignedUser *User      `gorm:"foreignKey:AssignedUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
}
