package seed

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_classes"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_effects"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_items"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_races"
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
func AllSeeds() []Seed {
	return []Seed{
		{
			Name: "01_components",
			Run:  faradhaven_classes.SeedComponents,
			HashData: func() (interface{}, int) {
				data := faradhaven_classes.SpellSystemComponents()
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
				components := faradhaven_classes.ClassComponentNames()
				return struct {
					Classes    interface{}
					Components interface{}
				}{classes, components}, len(classes)
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
	}
}
