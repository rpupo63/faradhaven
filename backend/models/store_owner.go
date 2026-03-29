package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// StoreOwner represents a vendor NPC (name, location, trade behavior).
//
// WillingnessToPurchase is a multiplier for how much the vendor pays when buying
// from the player vs a baseline (1.0 = normal, 0.5 = stingy).
//
// ExchangeRate is a multiplier on listed/book price when the player buys from
// this vendor (1.0 = list price, 1.2 = 20% markup).
type StoreOwner struct {
	ID                    uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Name                  string         `json:"name" gorm:"type:text;not null;uniqueIndex"`
	Location              string         `json:"location" gorm:"type:text;not null"`
	Personality           string         `json:"personality" gorm:"type:text"`
	WillingnessToPurchase float64        `json:"willingness_to_purchase" gorm:"type:double precision;not null;default:1"`
	ExchangeRate          float64        `json:"exchange_rate" gorm:"type:double precision;not null;default:1"`
	CategoriesObtained    pq.StringArray `json:"categories_obtained" gorm:"type:text[]"`
	CreatedAt             time.Time      `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt             time.Time      `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	CatalogRules []StoreOwnerCatalogRule `json:"catalog_rules,omitempty" gorm:"foreignKey:StoreOwnerID;constraint:OnDelete:CASCADE"`
}

// StoreOwnerCatalogRule links a vendor to stock: a specific item, a specific weapon,
// or a category with optional rarity allow-list.
//
// Row semantics: exactly one of (1) ItemID set, (2) WeaponID set, or (3) Category set
// with both IDs nil. AllowedRarities empty means any rarity matches.
type StoreOwnerCatalogRule struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();not null"`
	StoreOwnerID    uuid.UUID      `json:"store_owner_id" gorm:"type:uuid;not null;index"`
	ItemID          *uuid.UUID     `json:"item_id,omitempty" gorm:"type:uuid;index"`
	WeaponID        *uuid.UUID     `json:"weapon_id,omitempty" gorm:"type:uuid;index"`
	Category        *string        `json:"category,omitempty" gorm:"type:text"`
	AllowedRarities pq.StringArray `json:"allowed_rarities" gorm:"type:text[]"`
	CreatedAt       time.Time      `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	StoreOwner *StoreOwner `json:"-" gorm:"foreignKey:StoreOwnerID;references:ID;constraint:OnDelete:CASCADE"`
	Item       *Item       `json:"item,omitempty" gorm:"foreignKey:ItemID;references:ID;constraint:OnDelete:CASCADE"`
	Weapon     *Weapon     `json:"weapon,omitempty" gorm:"foreignKey:WeaponID;references:ID;constraint:OnDelete:CASCADE"`
}
