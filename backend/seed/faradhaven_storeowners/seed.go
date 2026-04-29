package faradhaven_storeowners

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/lib/pq"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/seed/batch"
	"github.com/rpupo63/faradhaven/backend/seed/uuids"
	"gorm.io/gorm"
)

// SeedFaradhavenStoreOwners upserts store owners and replaces catalog rules.
func SeedFaradhavenStoreOwners(tx *gorm.DB) error {
	seeds := AllStoreOwnerSeeds()
	owners := make([]models.StoreOwner, 0, len(seeds))
	var rules []models.StoreOwnerCatalogRule

	for _, vs := range seeds {
		ownerID := uuids.StoreOwnerUUID(vs.Name)
		owners = append(owners, models.StoreOwner{
			ID:                    ownerID,
			Name:                  vs.Name,
			Location:              vs.Location,
			Personality:           vs.Personality,
			WillingnessToPurchase: vs.WillingnessToPurchase,
			ExchangeRate:          vs.ExchangeRate,
			CategoriesObtained:    pq.StringArray(vs.CategoriesObtained),
		})

		for _, rs := range vs.Rules {
			if err := validateRuleSeed(rs); err != nil {
				return fmt.Errorf("vendor %q: %w", vs.Name, err)
			}
			disc := ruleDiscriminator(rs)
			rule := models.StoreOwnerCatalogRule{
				ID:              uuids.StoreOwnerCatalogRuleUUID(vs.Name, disc),
				StoreOwnerID:    ownerID,
				AllowedRarities: pq.StringArray(rs.AllowedRarities),
			}
			switch {
			case rs.ItemName != "":
				id := uuids.ItemUUID(rs.ItemName)
				rule.ItemID = &id
			case rs.WeaponName != "":
				id := uuids.WeaponUUID(rs.WeaponName)
				rule.WeaponID = &id
			default:
				cat := rs.Category
				rule.Category = &cat
			}
			rules = append(rules, rule)
		}
	}

	if err := tx.Exec("DELETE FROM store_owner_catalog_rules").Error; err != nil {
		return fmt.Errorf("could not clear store_owner_catalog_rules: %w", err)
	}

	if err := batch.UpsertBatchUpdateAll(tx, owners, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d store owners", len(owners))

	if err := batch.InsertBatch(tx, rules, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d store owner catalog rules", len(rules))
	return nil
}

func validateRuleSeed(rs CatalogRuleSeed) error {
	nSet := 0
	if rs.ItemName != "" {
		nSet++
	}
	if rs.WeaponName != "" {
		nSet++
	}
	if rs.Category != "" {
		nSet++
	}
	if nSet != 1 {
		return fmt.Errorf("catalog rule must set exactly one of item_name, weapon_name, or category (got %d)", nSet)
	}
	return nil
}

func ruleDiscriminator(rs CatalogRuleSeed) string {
	switch {
	case rs.ItemName != "":
		return "item:" + rs.ItemName
	case rs.WeaponName != "":
		return "weapon:" + rs.WeaponName
	default:
		r := slices.Clone(rs.AllowedRarities)
		slices.Sort(r)
		return "cat:" + rs.Category + ":" + strings.Join(r, ",")
	}
}
