package seed

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_classes"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_items"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_races"
)

// AllSeeds returns all registered seeds.
// Seeds run in alphabetical order by name.
// When GENERATE_MODELS=true, all data is cleared and reseeded from scratch.
func AllSeeds() []Seed {
	return []Seed{
		{
			Name: "01_classes",
			Run:  faradhaven_classes.SeedFaradhavenClasses,
		},
		{
			Name: "02_races",
			Run:  faradhaven_races.SeedFaradhavenRaces,
		},
		{
			Name: "03_items",
			Run:  faradhaven_items.SeedFaradhavenItems,
		},
	}
}
