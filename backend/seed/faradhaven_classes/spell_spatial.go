package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

func spatialComponents() []ComponentSeed {
	return []ComponentSeed{
		{
			Name:        "Expand",
			Symbol:      "Ex",
			Category:    models.CategorySpatial,
			Description: "Increases the scale/hitbox of the object (growth)",
			Element:     "",
		},
		{
			Name:        "Shrink",
			Symbol:      "Sh",
			Category:    models.CategorySpatial,
			Description: "Decreases the scale/hitbox (diminish, good for concentrated density)",
			Element:     "",
		},
		{
			Name:        "Haste",
			Symbol:      "Ha",
			Category:    models.CategorySpatial,
			Description: "Increases the simulation speed of the target (accelerate)",
			Element:     "",
		},
		{
			Name:        "Slow",
			Symbol:      "Sl",
			Category:    models.CategorySpatial,
			Description: "Decreases the simulation speed of the target (decelerate)",
			Element:     "",
		},
		{
			Name:        "Teleport",
			Symbol:      "Tp",
			Category:    models.CategorySpatial,
			Description: "Instant displacement of matter (blink)",
			Element:     "",
		},
		{
			Name:        "Echo",
			Symbol:      "Ec",
			Category:    models.CategorySpatial,
			Description: "Causes the effect to trigger multiple times (repeat, DoT or multi-hit)",
			Element:     "",
		},
		{
			Name:        "Extreme",
			Symbol:      "Xt",
			Category:    models.CategorySpatial,
			Description: "Pushes another component to its absolute limit (Slow to Stop, Heat to Incinerate, Cool to Absolute Zero)",
			Element:     "",
		},
	}
}
