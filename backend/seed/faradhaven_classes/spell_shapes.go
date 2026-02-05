package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

func shapeComponents() []ComponentSeed {
	return []ComponentSeed{
		{
			Name:        "Projectile",
			Symbol:      "Pr",
			Category:    models.CategoryShape,
			Description: "A standard moving object (bolt)",
			Element:     "",
		},
		{
			Name:        "Beam",
			Symbol:      "Bm",
			Category:    models.CategoryShape,
			Description: "Instant line-of-sight connection (ray)",
			Element:     "",
		},
		{
			Name:        "Nova",
			Symbol:      "Nv",
			Category:    models.CategoryShape,
			Description: "Radial explosion from a point (area of effect)",
			Element:     "",
		},
		{
			Name:        "Wall",
			Symbol:      "Wa",
			Category:    models.CategoryShape,
			Description: "A stationary vertical plane (barrier)",
			Element:     "",
		},
		{
			Name:        "Zone",
			Symbol:      "Zo",
			Category:    models.CategoryShape,
			Description: "A lingering area on the ground (field, puddles/clouds)",
			Element:     "",
		},
		{
			Name:        "Self",
			Symbol:      "Se",
			Category:    models.CategoryShape,
			Description: "Applies directly to the caster (buff/cloak)",
			Element:     "",
		},
		{
			Name:        "Trap",
			Symbol:      "Tr",
			Category:    models.CategoryShape,
			Description: "Stationary object that triggers on contact (rune)",
			Element:     "",
		},
		{
			Name:        "Cone",
			Symbol:      "Cn",
			Category:    models.CategoryShape,
			Description: "A cone-shaped blast originating from the caster (breath weapon)",
			Element:     "",
		},
	}
}
