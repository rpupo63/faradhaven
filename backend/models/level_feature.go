package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// LevelFeature represents a single feature gained at a class level
type LevelFeature struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ClassLevelID uuid.UUID  `json:"class_level_id" gorm:"type:uuid;not null;index"`
	ArchetypeID  *uuid.UUID `json:"archetype_id,omitempty" gorm:"type:uuid;index"` // nil = all archetypes get this feature
	Name         string     `json:"name" gorm:"type:text;not null"`
	Description  string     `json:"description" gorm:"type:text"`
	SortOrder    int        `json:"sort_order" gorm:"type:int;default:0"` // for ordering multiple features at same level

	// Activation & Usage (mirrors models.Trait fields for class features)
	ActionType     string         `json:"action_type" gorm:"type:text"`               // "Action", "Bonus Action", "Reaction", ""
	UsesPerRest    string         `json:"uses_per_rest" gorm:"type:text"`             // "1", "Proficiency Bonus", ""
	ResetCondition string         `json:"reset_condition" gorm:"type:text"`           // "Long Rest", "Short Rest", ""
	ResourceCosts  datatypes.JSON `json:"resource_costs,omitempty" gorm:"type:jsonb"` // [{"key":"max_blood_ichor","amount":2}]
	ResourceGains  datatypes.JSON `json:"resource_gains,omitempty" gorm:"type:jsonb"` // [{"key":"max_blood_ichor","amount":1}]

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	ClassLevel ClassLevel `json:"-" gorm:"foreignKey:ClassLevelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Archetype  *Archetype `json:"archetype,omitempty" gorm:"foreignKey:ArchetypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
