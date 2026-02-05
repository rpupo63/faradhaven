package models

import (
	"time"

	"github.com/google/uuid"
)

// LevelFeature represents a single feature gained at a class level
type LevelFeature struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ClassLevelID uuid.UUID  `json:"class_level_id" gorm:"type:uuid;not null;index"`
	ArchetypeID  *uuid.UUID `json:"archetype_id,omitempty" gorm:"type:uuid;index"` // nil = all archetypes get this feature
	Name         string     `json:"name" gorm:"type:text;not null"`
	Description  string     `json:"description" gorm:"type:text"`
	SortOrder    int        `json:"sort_order" gorm:"type:int;default:0"` // for ordering multiple features at same level
	CreatedAt    time.Time  `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	ClassLevel ClassLevel `json:"-" gorm:"foreignKey:ClassLevelID;references:ID"`
	Archetype  *Archetype `json:"archetype,omitempty" gorm:"foreignKey:ArchetypeID;references:ID"`
}
