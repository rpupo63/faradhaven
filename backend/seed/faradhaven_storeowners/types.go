package faradhaven_storeowners

// CatalogRuleSeed describes one row in store_owner_catalog_rules after items/weapons exist.
type CatalogRuleSeed struct {
	ItemName        string   // if non-empty, rule targets this item by seed name
	WeaponName      string   // if non-empty, rule targets this weapon by seed name
	Category        string   // if non-empty with no item/weapon, category-wide rule
	AllowedRarities []string // empty = any rarity; otherwise item/weapon rarity must be listed
}

// StoreOwnerSeed is declarative vendor data for hashing and seeding.
type StoreOwnerSeed struct {
	Name                  string
	Location              string
	Personality           string
	WillingnessToPurchase float64
	ExchangeRate          float64
	CategoriesObtained    []string
	Rules                 []CatalogRuleSeed
}
