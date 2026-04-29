package models

import "github.com/google/uuid"

type ItemLootTheme struct {
	ItemID uuid.UUID `json:"item_id" gorm:"type:uuid;primaryKey;not null"`
	Theme  LootTheme `json:"theme" gorm:"type:text;primaryKey;not null"`
	Weight float64   `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Item   *Item     `json:"-" gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ItemLootTheme) TableName() string { return "item_loot_themes" }

type ItemLootLocation struct {
	ItemID   uuid.UUID    `json:"item_id" gorm:"type:uuid;primaryKey;not null"`
	Location LootLocation `json:"location" gorm:"type:text;primaryKey;not null"`
	Weight   float64      `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Item     *Item        `json:"-" gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ItemLootLocation) TableName() string { return "item_loot_locations" }

type ItemLootSource struct {
	ItemID uuid.UUID  `json:"item_id" gorm:"type:uuid;primaryKey;not null"`
	Source LootSource `json:"source" gorm:"type:text;primaryKey;not null"`
	Weight float64    `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Item   *Item      `json:"-" gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ItemLootSource) TableName() string { return "item_loot_sources" }

type ItemLootTier struct {
	ItemID uuid.UUID `json:"item_id" gorm:"type:uuid;primaryKey;not null"`
	Tier   LootTier  `json:"tier" gorm:"type:text;primaryKey;not null"`
	Weight float64   `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Item   *Item     `json:"-" gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ItemLootTier) TableName() string { return "item_loot_tiers" }

type ItemLootRewardAmount struct {
	ItemID       uuid.UUID        `json:"item_id" gorm:"type:uuid;primaryKey;not null"`
	RewardAmount LootRewardAmount `json:"reward_amount" gorm:"type:text;primaryKey;not null"`
	Weight       float64          `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Item         *Item            `json:"-" gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ItemLootRewardAmount) TableName() string { return "item_loot_reward_amounts" }

type ItemLootLevelBand struct {
	ItemID    uuid.UUID     `json:"item_id" gorm:"type:uuid;primaryKey;not null"`
	LevelBand LootLevelBand `json:"level_band" gorm:"type:text;primaryKey;not null"`
	Weight    float64       `json:"weight" gorm:"type:double precision;default:1.0;not null"`
	Item      *Item         `json:"-" gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ItemLootLevelBand) TableName() string { return "item_loot_level_bands" }
