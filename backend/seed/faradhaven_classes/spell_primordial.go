package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

func primordialComponents() []ComponentSeed {
	return []ComponentSeed{
		{
			Name:        "Ignis",
			Symbol:      "Ig",
			Category:    models.CategoryPrimordial,
			Description: "Raw plasma, thermal energy, flame substance",
			Element:     "fire",
		},
		{
			Name:        "Aqua",
			Symbol:      "Aq",
			Category:    models.CategoryPrimordial,
			Description: "Liquid, fluid dynamics, dousing",
			Element:     "water",
		},
		{
			Name:        "Terra",
			Symbol:      "Te",
			Category:    models.CategoryPrimordial,
			Description: "Dirt, stone, physical mass",
			Element:     "earth",
		},
		{
			Name:        "Aer",
			Symbol:      "Ae",
			Category:    models.CategoryPrimordial,
			Description: "Gas, wind, pressure",
			Element:     "air",
		},
		{
			Name:        "Fulgur",
			Symbol:      "Fu",
			Category:    models.CategoryPrimordial,
			Description: "Electricity, charged particles, current, voltage",
			Element:     "lightning",
		},
		{
			Name:        "Ferrum",
			Symbol:      "Fe",
			Category:    models.CategoryPrimordial,
			Description: "Hardness, conductivity, magnetism",
			Element:     "metal",
		},
		{
			Name:        "Vita",
			Symbol:      "Vi",
			Category:    models.CategoryPrimordial,
			Description: "Plants, vines, wood, life force, organic growth",
			Element:     "nature",
		},
		{
			Name:        "Umbra",
			Symbol:      "Um",
			Category:    models.CategoryPrimordial,
			Description: "Darkness, absence of light, void, shadow substance",
			Element:     "shadow",
		},
		{
			Name:        "Lux",
			Symbol:      "Lx",
			Category:    models.CategoryPrimordial,
			Description: "Illumination, radiant energy, photons, pure light",
			Element:     "light",
		},
		{
			Name:        "Arcanum",
			Symbol:      "Ar",
			Category:    models.CategoryPrimordial,
			Description: "Pure mana, purple energy, force fields",
			Element:     "arcane",
		},
		{
			Name:        "Sonus",
			Symbol:      "Sn",
			Category:    models.CategoryPrimordial,
			Description: "Sound waves, vibration, thunder, shockwaves",
			Element:     "thunder",
		},
		{
			Name:        "Acidum",
			Symbol:      "Ac",
			Category:    models.CategoryPrimordial,
			Description: "Corrosive substances, acid, chemical dissolution",
			Element:     "acid",
		},
		{
			Name:        "Psi",
			Symbol:      "Ps",
			Category:    models.CategoryPrimordial,
			Description: "Psychic energy, mental force, thought-based attacks",
			Element:     "psychic",
		},
		{
			Name:        "Sanctus",
			Symbol:      "Sa",
			Category:    models.CategoryPrimordial,
			Description: "Divine radiance, holy energy, sacred power",
			Element:     "radiant",
		},
	}
}
