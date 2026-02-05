package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

func physicalComponents() []ComponentSeed {
	return []ComponentSeed{
		{
			Name:        "Push",
			Symbol:      "Pu",
			Category:    models.CategoryPhysical,
			Description: "Applies force away from the origin (repel)",
			Element:     "force",
		},
		{
			Name:        "Pull",
			Symbol:      "Pl",
			Category:    models.CategoryPhysical,
			Description: "Applies force toward the origin (attract, implosion, magnetism)",
			Element:     "force",
		},
		{
			Name:        "Lift",
			Symbol:      "Li",
			Category:    models.CategoryPhysical,
			Description: "Negates gravity or applies upward force (levitate)",
			Element:     "force",
		},
		{
			Name:        "Crush",
			Symbol:      "Cr",
			Category:    models.CategoryPhysical,
			Description: "Increases gravity or applies downward force",
			Element:     "force",
		},
		{
			Name:        "Spin",
			Symbol:      "Sp",
			Category:    models.CategoryPhysical,
			Description: "Adds rotational force (creates tornadoes or drills)",
			Element:     "force",
		},
		{
			Name:        "Pierce",
			Symbol:      "Pi",
			Category:    models.CategoryPhysical,
			Description: "Ignores collision/armor for a set distance (penetrate)",
			Element:     "force",
		},
		{
			Name:        "Scatter",
			Symbol:      "Sc",
			Category:    models.CategoryPhysical,
			Description: "Spreads force over a wide area (disperse, shotgun effect)",
			Element:     "force",
		},
		{
			Name:        "Focus",
			Symbol:      "Fo",
			Category:    models.CategoryPhysical,
			Description: "Condenses force into a single point (concentrate, precision)",
			Element:     "force",
		},
		{
			Name:        "Homing",
			Symbol:      "Ho",
			Category:    models.CategoryPhysical,
			Description: "Causes projectiles to track and follow targets",
			Element:     "force",
		},
	}
}
