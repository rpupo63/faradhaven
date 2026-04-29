package models

import "time"

// LootLevelBudgetProfile stores normal-distribution parameters for each loot level.
type LootLevelBudgetProfile struct {
	LootLevel int       `json:"loot_level" gorm:"primaryKey;not null"`
	MeanGP    float64   `json:"mean_gp" gorm:"type:double precision;not null"`
	SigmaGP   float64   `json:"sigma_gp" gorm:"type:double precision;not null"`
	MinGP     float64   `json:"min_gp" gorm:"type:double precision;default:0;not null"`
	MaxGP     float64   `json:"max_gp" gorm:"type:double precision;default:0;not null"` // 0 means uncapped
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
}

