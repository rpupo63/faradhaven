package faradhaven_races

// Flamekin returns the Flamekin race seed
func Flamekin() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:           "Flamekin",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/flamekin.jpg",
		Description:    "Flamekin are people made from two key elements of creation: fire and stone. As a result, many flamekin feel a strong connection to the natural world. Flamekin's bodies radiate harmless magical flames, though they possess innate magic that allows them to create burning flames in a multitude of forms.\n\nFlamekin view self-discovery and self-expression the noblest of aspirations and believe that self-realization is the most important thing an individual can do with their life. Flamekin call this lifelong pursuit the Path of Flame.\n\nFlamekin dwell in either Lorwyn or Shadowmoor. Physically and culturally, they are similar in both lands.",
		CreatureType:   "Humanoid",
		Size:           "Medium (about 4–7 feet tall) or Small (about 2–4 feet tall), chosen when you select this species",
		BaseSpeed:      30,
		Traits:         flamekinTraits(),
		ComponentNames: []string{"Ignis", "Self", "Nova", "Imbue"},
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
		{
			Name:           "Reach to the Blaze",
			Description:    "You know the Ignis + Self component combination (producing flames in your hand). Starting at 3rd level, you can cast Ignis + Nova (a cone of burning fire). Starting at 5th level, you can also cast Ignis + Imbue (a blade of pure flame), without requiring a material component. Once you cast Ignis + Nova or Ignis + Imbue with this trait, you can't cast that combination with it again until you finish a long rest. You can also cast either of those combinations using any spell slots you have of the appropriate level. Intelligence, Wisdom, or Charisma is your spellcasting ability for these spells when you cast them with this trait (choose when you select this race).",
			LevelReq:       1,
			ActionType:     "Passive",
			UsesPerRest:    "1 per spell",
			ResetCondition: "Long Rest",
			Options:        flamekinSpellcastingOptions(),
		},
	}
}

func flamekinSpellcastingOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:        "Intelligence",
			Description: "Intelligence is your spellcasting ability for Ignis + Self, Ignis + Nova, and Ignis + Imbue when you cast them with this trait.",
		},
		{
			Name:        "Wisdom",
			Description: "Wisdom is your spellcasting ability for Ignis + Self, Ignis + Nova, and Ignis + Imbue when you cast them with this trait.",
		},
		{
			Name:        "Charisma",
			Description: "Charisma is your spellcasting ability for Ignis + Self, Ignis + Nova, and Ignis + Imbue when you cast them with this trait.",
		},
	}
}
