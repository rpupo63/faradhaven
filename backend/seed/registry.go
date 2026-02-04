package seed

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_classes"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_items"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_races"
)

// AllSeeds returns all registered seeds in version order.
// Add new seeds here with a date prefix for proper ordering.
//
// Naming convention: YYYYMMDD_category_description
// Examples:
//   - 20240101_races_initial
//   - 20240115_classes_initial
//   - 20240201_items_weapons
//   - 20240215_races_add_new_race
func AllSeeds() []Seed {
	return []Seed{
		// Initial seeds - run once
		{
			Name: "20240101_races_initial",
			Run:  faradhaven_races.SeedFaradhavenRaces,
		},
		{
			Name: "20240102_classes_initial",
			Run:  faradhaven_classes.SeedFaradhavenClasses,
		},
		{
			Name: "20240103_items_weapons",
			Run:  faradhaven_items.SeedFaradhavenItems,
		},

		// Add new seeds below with incrementing dates
		// Example: adding a new race later
		// {
		// 	Name: "20240215_races_add_dragonkin",
		// 	Run: func(db *gorm.DB) error {
		// 		// seed just the new race
		// 	},
		// },
	}
}
