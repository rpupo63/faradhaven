package faradhaven_races

// Kalashtar returns the Kalashtar race seed
func Kalashtar() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:               "Kalashtar",
		PhotoURL:           "https://photos-for-apps.s3.us-east-2.amazonaws.com/kalashtar.jpg",
		Description:        "Kalashtar (pronounced kal-ASH-tar) are created from the union of humanity and renegade spirits called quori from the plane of dreams. Kalashtar appear human, but their spiritual connection affects them in a variety of ways. They have symmetrical, slightly angular features, and their eyes often glow when they are concentrating or expressing strong emotions. Kalashtar can't communicate directly with their quori spirits. Rather, kalashtar might experience them as a source of instinct and inspiration, drawing on the spirits' memories when the kalashtar sleep. This connection grants kalashtar minor psionic abilities, as well as protection from psionic attacks.",
		CreatureType:       "Aberration",
		Size:               "Medium (about 6–7 feet tall)",
		BaseSpeed:          30,
		Languages:          []string{"Common", "Quori"},
		BonusLanguageCount: 1,
		Traits:             kalashtarTraits(),
	}
}

func kalashtarTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Dual Mind",
			Description: "You have Advantage on Wisdom and Charisma saving throws.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Mental Discipline",
			Description: "You have Resistance to Psychic damage.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Mind Link",
			Description: "You have telepathy with a range in feet equal to 10 times your level. When you're using this trait to speak telepathically to a creature, you can take a Magic action to give that creature the ability to speak telepathically with you for 1 hour or until you take another Magic action to end this effect.",
			LevelReq:    1,
			ActionType:  "Magic Action",
			RangeValue:  "10 × level",
		},
		{
			Name:           "Severed from Dreams",
			Description:    "You can't be the target of the Dream spell. In addition, when you finish a Long Rest, you gain proficiency in one skill of your choice. This proficiency lasts until you finish another Long Rest.",
			LevelReq:       1,
			ActionType:     "Passive",
			ResetCondition: "Long Rest",
		},
	}
}
