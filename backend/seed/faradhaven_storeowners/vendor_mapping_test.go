package faradhaven_storeowners

import (
	"testing"

	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_items"
)

// TestVendorItemWeaponNamesExist ensures every explicit catalog rule references a seeded item or weapon name.
func TestVendorItemWeaponNamesExist(t *testing.T) {
	itemNames := map[string]struct{}{}
	for _, it := range faradhaven_items.AllItems() {
		itemNames[it.Name] = struct{}{}
	}
	weaponNames := map[string]struct{}{}
	for _, w := range faradhaven_items.AllWeapons() {
		weaponNames[w.Name] = struct{}{}
	}

	for _, v := range AllStoreOwnerSeeds() {
		for _, r := range v.Rules {
			if r.ItemName != "" {
				if _, ok := itemNames[r.ItemName]; !ok {
					t.Errorf("vendor %q: unknown ItemName %q (not in faradhaven_items.AllItems)", v.Name, r.ItemName)
				}
			}
			if r.WeaponName != "" {
				if _, ok := weaponNames[r.WeaponName]; !ok {
					t.Errorf("vendor %q: unknown WeaponName %q (not in faradhaven_items.AllWeapons)", v.Name, r.WeaponName)
				}
			}
		}
	}
}
