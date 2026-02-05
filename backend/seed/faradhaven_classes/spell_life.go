package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

func lifeComponents() []ComponentSeed {
	return []ComponentSeed{
		{
			Name:        "Mend",
			Symbol:      "Me",
			Category:    models.CategoryLife,
			Description: "Restores structural integrity or HP (heal)",
			Element:     "",
		},
		{
			Name:        "Wither",
			Symbol:      "Wi",
			Category:    models.CategoryLife,
			Description: "Adds damage over time or rust/corrosion (decay)",
			Element:     "",
		},
		{
			Name:        "Revive",
			Symbol:      "Rv",
			Category:    models.CategoryLife,
			Description: "Brings a destroyed object/entity back to a base state (resurrect)",
			Element:     "",
		},
		{
			Name:        "Drain",
			Symbol:      "Dr",
			Category:    models.CategoryLife,
			Description: "Transfers energy/HP from target to caster (leech)",
			Element:     "",
		},
	}
}
