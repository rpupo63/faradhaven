package seed

import (
	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_classes"
	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_components"
	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_effects"
	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_items"
	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_races"
	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_storeowners"
)

// AllSeeds returns all registered seeds in dependency order.
// Seeds run in alphabetical order by name, so prefixes control execution order.
// When GENERATE_MODELS=true, all data is cleared and reseeded from scratch.
//
// Dependency order:
// 1. Components - independent, used by classes and races
// 2. Items - weapons and items, referenced by classes for starting equipment
// 3. Races - references components
// 4. Classes - references components and items
// 5. Effects
// 6. Existing spells
// 7. Store owners - references seeded items and weapons
func AllSeeds() []Seed {
	return []Seed{
		{
			Name: "01_components",
			Run:  faradhaven_components.SeedFaradhavenComponents,
			HashData: func() (interface{}, int) {
				data := faradhaven_components.AllComponents()
				return data, len(data)
			},
		},
		{
			Name: "02_items",
			Run:  faradhaven_items.SeedFaradhavenItems,
			HashData: func() (interface{}, int) {
				weapons := faradhaven_items.AllWeapons()
				items := faradhaven_items.AllItems()
				return struct {
					Weapons interface{}
					Items   interface{}
				}{weapons, items}, len(weapons) + len(items)
			},
		},
		{
			Name: "03_races",
			Run:  faradhaven_races.SeedFaradhavenRaces,
			HashData: func() (interface{}, int) {
				data := faradhaven_races.AllRaces()
				return data, len(data)
			},
		},
		{
			Name: "04_classes",
			Run:  faradhaven_classes.SeedFaradhavenClasses,
			HashData: func() (interface{}, int) {
				classes := faradhaven_classes.AllClasses()
				return classes, len(classes)
			},
		},
		{
			Name: "05_effects",
			Run:  faradhaven_effects.SeedFaradhavenEffects,
			HashData: func() (interface{}, int) {
				data := faradhaven_effects.AllEffects()
				return data, len(data)
			},
		},
		{
			Name: "06_existing_spells",
			Run:  SeedExistingSpells,
		},
		{
			Name: "07_store_owners",
			Run:  faradhaven_storeowners.SeedFaradhavenStoreOwners,
			HashData: func() (interface{}, int) {
				data := faradhaven_storeowners.AllStoreOwnerSeeds()
				n := 0
				for _, v := range data {
					n += 1 + len(v.Rules)
				}
				return data, n
			},
		},
	}
}
