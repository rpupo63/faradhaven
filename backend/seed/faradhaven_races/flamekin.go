package faradhaven_races

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/seedmedia"
)

// Flamekin returns the Flamekin race seed
func Flamekin() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Flamekin",
		PhotoURL:     seedmedia.URL("flamekin.jpg"),
		Description:  "Flamekin are people made from two key elements of creation: fire and stone. As a result, many flamekin feel a strong connection to the natural world. Flamekin's bodies radiate harmless magical flames, though they possess innate magic that allows them to create burning flames in a multitude of forms.\n\nFlamekin view self-discovery and self-expression the noblest of aspirations and believe that self-realization is the most important thing an individual can do with their life. Flamekin call this lifelong pursuit the Path of Flame.\n\nFlamekin dwell in either Lorwyn or Shadowmoor. Physically and culturally, they are similar in both lands.",
		CreatureType: "Humanoid",
		Size:         "Medium (about 4–7 feet tall) or Small (about 2–4 feet tall), chosen when you select this species",
		BaseSpeed:    30,
		Languages:    []string{"Common", "Primordial"},
		Traits:       flamekinTraits(),
	}
}

func flamekinTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Darkvision",
			Description: "You can see in dim light within 60 feet of you as if it were bright light and in darkness as if it were dim light. You discern colors in that darkness only as shades of gray.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "60",
		},
		{
			Name:        "Fire Resistance",
			Description: "You have resistance to fire damage.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
	}
}
