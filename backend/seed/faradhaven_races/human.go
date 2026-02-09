package faradhaven_races

// Human returns the Human race seed
func Human() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Human",
		PhotoURL:     "https://photos-for-apps.s3.us-east-2.amazonaws.com/human.jpg",
		Description:  "Found throughout the multiverse, humans are as varied as they are numerous, and they endeavor to achieve as much as they can in the years they are given. Their ambition and resourcefulness are commended, respected, and feared on many worlds.\n\nHumans are as diverse in appearance as the people of Earth, and they have many gods. Scholars dispute the origin of humanity, but one of the earliest known human gatherings is said to have occurred in Sigil, the torus-shaped city at the center of the multiverse and the place where the Common language was born. From there, humans could have spread to every part of the multiverse, bringing the City of Doors' cosmopolitanism with them.",
		CreatureType: "Humanoid",
		Size:         "Medium (about 4-7 feet tall) or Small (about 2-4 feet tall), chosen when you select this species",
		BaseSpeed:    30,
		AbilityScoreBonuses: map[string]int{
			"strength":     1,
			"dexterity":    1,
			"constitution": 1,
			"intelligence": 1,
			"wisdom":       1,
			"charisma":     1,
		},
		Languages:          []string{"Common"},
		BonusLanguageCount: 1,
		Traits:             humanTraits(),
	}
}

func humanTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Resourceful",
			Description: "You gain Heroic Inspiration whenever you finish a Long Rest.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Skillful",
			Description: "You gain proficiency in one skill of your choice.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Versatile",
			Description: "You gain an Origin feat of your choice. Skilled is recommended.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
	}
}
