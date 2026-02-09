package faradhaven_races

// Dhampir returns the Dhampir race seed
func Dhampir() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:               "Dhampir",
		PhotoURL:           "https://photos-for-apps.s3.us-east-2.amazonaws.com/dhampir.jpg",
		Description:        "Dhampirs are living people who possess vampiric prowess but are cursed with macabre hunger. Most dhampirs thirst for blood, but some gain sustenance from dreams, life energy, or other vital sources. Dhampirs must choose whether to fight to control their hunger or give in to predatory urges.\n\nDhampirs often arise from encounters with vampires; some are the descendants of a powerful vampire, while others are partially transformed by a vampire's bite. All manner of macabre bargains and necromantic influences might also give rise to a dhampir. Regardless of their origins, dhampirs exhibit their vampiric nature in various ways, including increased speed and a life-draining bite.",
		CreatureType:       "Humanoid",
		Size:               "Medium (about 4–7 feet tall) or Small (about 2–4 feet tall), chosen when you select this species",
		BaseSpeed:          35,
		Languages:          []string{"Common"},
		BonusLanguageCount: 1,
		Traits:             dhampirTraits(),
	}
}

func dhampirTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 60 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "60",
		},
		{
			Name:        "Spider Climb",
			Description: "You have a Climb Speed equal to your Speed. When you reach character level 3, you can move up, down, and across vertical surfaces and along ceilings while leaving your hands free.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Trace of Undeath",
			Description: "You have Resistance to Necrotic damage.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:           "Vampiric Bite",
			Description:    "When you use your Unarmed Strike and deal damage, you can choose to bite with your fangs. You deal Piercing damage equal to 1d4 plus your Constitution modifier instead of the normal damage of an Unarmed Strike.\n\nIn addition, when you deal this damage to a creature that isn't a Construct or an Undead, you can empower yourself in one of the following ways. You can empower yourself with this trait a number of times equal to your Proficiency Bonus, and you regain all expended uses when you finish a Long Rest.",
			LevelReq:       1,
			ActionType:     "Passive",
			UsesPerRest:    "Proficiency Bonus",
			ResetCondition: "Long Rest",
			Options:        dhampirVampiricBiteOptions(),
		},
	}
}

func dhampirVampiricBiteOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:        "Drain",
			Description: "You regain Hit Points equal to the Piercing damage dealt.",
		},
		{
			Name:        "Strengthen",
			Description: "You gain a bonus to the next ability check or attack roll you make within the next minute; the bonus is equal to the Piercing damage dealt.",
		},
	}
}
