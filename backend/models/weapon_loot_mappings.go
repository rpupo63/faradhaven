package models

import "github.com/google/uuid"

type WeaponLootTheme struct {
	WeaponID uuid.UUID `json:"weapon_id" gorm:"type:uuid;primaryKey;not null"`
	Theme    LootTheme `json:"theme" gorm:"type:text;primaryKey;not null"`
	Weight   float64   `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Weapon   *Weapon   `json:"-" gorm:"foreignKey:WeaponID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WeaponLootTheme) TableName() string { return "weapon_loot_themes" }

type WeaponLootLocation struct {
	WeaponID uuid.UUID    `json:"weapon_id" gorm:"type:uuid;primaryKey;not null"`
	Location LootLocation `json:"location" gorm:"type:text;primaryKey;not null"`
	Weight   float64      `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Weapon   *Weapon      `json:"-" gorm:"foreignKey:WeaponID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WeaponLootLocation) TableName() string { return "weapon_loot_locations" }

type WeaponLootSource struct {
	WeaponID uuid.UUID  `json:"weapon_id" gorm:"type:uuid;primaryKey;not null"`
	Source   LootSource `json:"source" gorm:"type:text;primaryKey;not null"`
	Weight   float64    `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Weapon   *Weapon    `json:"-" gorm:"foreignKey:WeaponID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WeaponLootSource) TableName() string { return "weapon_loot_sources" }

type WeaponLootTier struct {
	WeaponID uuid.UUID `json:"weapon_id" gorm:"type:uuid;primaryKey;not null"`
	Tier     LootTier  `json:"tier" gorm:"type:text;primaryKey;not null"`
	Weight   float64   `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Weapon   *Weapon   `json:"-" gorm:"foreignKey:WeaponID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WeaponLootTier) TableName() string { return "weapon_loot_tiers" }

type WeaponLootRewardAmount struct {
	WeaponID     uuid.UUID        `json:"weapon_id" gorm:"type:uuid;primaryKey;not null"`
	RewardAmount LootRewardAmount `json:"reward_amount" gorm:"type:text;primaryKey;not null"`
	Weight       float64          `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Weapon       *Weapon          `json:"-" gorm:"foreignKey:WeaponID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WeaponLootRewardAmount) TableName() string { return "weapon_loot_reward_amounts" }

type WeaponLootLevelBand struct {
	WeaponID  uuid.UUID     `json:"weapon_id" gorm:"type:uuid;primaryKey;not null"`
	LevelBand LootLevelBand `json:"level_band" gorm:"type:text;primaryKey;not null"`
	Weight    float64       `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Weapon    *Weapon       `json:"-" gorm:"foreignKey:WeaponID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (WeaponLootLevelBand) TableName() string { return "weapon_loot_level_bands" }
