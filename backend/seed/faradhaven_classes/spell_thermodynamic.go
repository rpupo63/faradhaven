package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

func thermodynamicComponents() []ComponentSeed {
	return []ComponentSeed{
		{
			Name:        "Heat",
			Symbol:      "He",
			Category:    models.CategoryThermodynamic,
			Description: "Increases temperature (exothermic). Turns Water to Steam, Earth to Magma",
			Element:     "",
		},
		{
			Name:        "Cool",
			Symbol:      "Co",
			Category:    models.CategoryThermodynamic,
			Description: "Decreases temperature (endothermic). Turns Water to Ice, Air to Liquid Nitrogen",
			Element:     "",
		},
		{
			Name:        "Conduct",
			Symbol:      "Cd",
			Category:    models.CategoryThermodynamic,
			Description: "Allows electricity/magic to flow through the material (charge)",
			Element:     "",
		},
		{
			Name:        "Insulate",
			Symbol:      "In",
			Category:    models.CategoryThermodynamic,
			Description: "Prevents energy transfer (block, creates shields)",
			Element:     "",
		},
	}
}
