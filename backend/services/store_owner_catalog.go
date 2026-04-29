package services

import (
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
)

// StoreOwnerSellsItem returns true if the owner's catalog rules allow selling this item or weapon.
func StoreOwnerSellsItem(owner *models.StoreOwner, itemID uuid.UUID, category, rarity string, isWeapon bool) bool {
	if owner == nil {
		return false
	}
	for i := range owner.CatalogRules {
		rule := &owner.CatalogRules[i]
		if rule.ItemID != nil && *rule.ItemID == itemID && !isWeapon {
			return true
		}
		if rule.WeaponID != nil && *rule.WeaponID == itemID && isWeapon {
			return true
		}
		if rule.Category != nil && category == *rule.Category {
			if len(rule.AllowedRarities) == 0 {
				return true
			}
			for _, ar := range rule.AllowedRarities {
				if ar == rarity {
					return true
				}
			}
		}
	}
	return false
}
