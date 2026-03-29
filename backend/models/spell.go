package models

import (
	"time"

	"github.com/google/uuid"
)

// Spell represents a crafted spell in the game
type Spell struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	CharacterID *uuid.UUID `json:"character_id,omitempty" gorm:"type:uuid;index"` // Character who created/prepared this spell (nil = user-level spell)
	Name        string     `json:"name" gorm:"type:text;not null"`
	Description string     `json:"description" gorm:"type:text"`
	Level       int        `json:"level" gorm:"type:int;not null;default:1"`

	// Spell Mechanics
	Type                SpellType   `json:"type" gorm:"type:text;not null;default:'Utility'"` // Attack, Save, Effect, Healing, Utility
	Range               *int        `json:"range,omitempty" gorm:"type:integer"`
	Duration            *string     `json:"duration,omitempty" gorm:"type:text"` // see ValidateSpellDuration
	Concentration         bool           `json:"concentration" gorm:"default:false"`
	SaveAttr              *SaveAttribute `json:"save_attr,omitempty" gorm:"type:text"`
	DamageDiceCount       *int           `json:"damage_dice_count,omitempty" gorm:"type:integer"`
	DamageDieSize         *int           `json:"damage_die_size,omitempty" gorm:"type:integer"` // faces: 6 = d6
	DamageType            *DamageType    `json:"damage_type,omitempty" gorm:"type:text"`
	SuggestedDamageDiceCount *int        `json:"suggested_damage_dice_count,omitempty" gorm:"type:integer"`
	SuggestedDamageDieSize   *int        `json:"suggested_damage_die_size,omitempty" gorm:"type:integer"`
	AddModifier         bool        `json:"add_modifier" gorm:"default:false"`                // Add spellcasting mod to damage/healing

	Checked bool `json:"checked" gorm:"type:boolean;not null;default:false"` // GM has reviewed and approved this spell

	// AI Review Opinions
	AIDescriptionOpinion *string `json:"ai_description_opinion,omitempty" gorm:"column:ai_description_opinion;type:text"`
	AIDamageOpinion      *string `json:"ai_damage_opinion,omitempty" gorm:"column:ai_damage_opinion;type:text"`
	AIEffectOpinion      *string `json:"ai_effect_opinion,omitempty" gorm:"column:ai_effect_opinion;type:text"`
	AIOverallVerdict     *string `json:"ai_overall_verdict,omitempty" gorm:"column:ai_overall_verdict;type:text"`
	AIRawOutput          *string `json:"ai_raw_output,omitempty" gorm:"column:ai_raw_output;type:text"`

	// AI Recommended Edits (mechanics align with Type / Range / Duration on this model)
	AIRecommendedName        *string `json:"ai_recommended_name,omitempty" gorm:"column:ai_recommended_name;type:text"`
	AIRecommendedDescription *string `json:"ai_recommended_description,omitempty" gorm:"column:ai_recommended_description;type:text"`
	AIRecommendedType        *SpellType `json:"ai_recommended_type,omitempty" gorm:"column:ai_recommended_type;type:text"`
	AIRecommendedRange       *int       `json:"ai_recommended_range,omitempty" gorm:"column:ai_recommended_range;type:integer"`
	AIRecommendedDuration       *string      `json:"ai_recommended_duration,omitempty" gorm:"column:ai_recommended_duration;type:text"` // see ValidateSpellDuration
	AIRecommendedDamageDiceCount *int        `json:"ai_recommended_damage_dice_count,omitempty" gorm:"column:ai_recommended_damage_dice_count;type:integer"`
	AIRecommendedDamageDieSize   *int        `json:"ai_recommended_damage_die_size,omitempty" gorm:"column:ai_recommended_damage_die_size;type:integer"`
	AIRecommendedDamageType      *DamageType `json:"ai_recommended_damage_type,omitempty" gorm:"column:ai_recommended_damage_type;type:text"`
	AIRecommendedSaveAttr        *SaveAttribute `json:"ai_recommended_save_attr,omitempty" gorm:"column:ai_recommended_save_attr;type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User       User        `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Character  *Character  `json:"-" gorm:"foreignKey:CharacterID;references:ID;constraint:OnDelete:SET NULL"`
	Components []Component `json:"components" gorm:"many2many:spell_components;foreignKey:ID;joinForeignKey:SpellID;References:ID;joinReferences:ComponentID"`
}
