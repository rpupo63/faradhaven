package models

import (
	"time"

	"github.com/google/uuid"
)

// TrapProperties stores specific properties for a MapElement of Type "trap".
type TrapProperties struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	MapElementID       uuid.UUID `json:"map_element_id" gorm:"type:uuid;not null;uniqueIndex"` // One-to-one with MapElement
	Name               string    `json:"name" gorm:"type:text;not null"`
	Description        string    `json:"description" gorm:"type:text"`
	VisibleToPlayers   bool      `json:"visible_to_players" gorm:"type:boolean;default:false"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	MapElement MapElement `json:"-" gorm:"foreignKey:MapElementID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM
func (TrapProperties) TableName() string {
	return "map_element_trap_properties"
}
