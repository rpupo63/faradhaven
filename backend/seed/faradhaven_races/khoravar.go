package faradhaven_races

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/seedmedia"
)

// Khoravar returns the Khoravar race seed
func Khoravar() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Khoravar",
		PhotoURL:     seedmedia.URL("khoravar.jpg"),
		Description:  "Over the course of centuries, those descended from both humans and elves have developed their own communities and traditions in Khorvaire. The rise of House Lyrandar and House Medani has strengthened this identity. Members of these communities call themselves Khoravar, an Elvish term meaning \"children of Khorvaire,\" as they dislike the term \"half-elf.\"\n\nMany Khoravar espouse the idea of being \"the bridge between,\" believing they are called to facilitate communication and cooperation between members of different cultures or species. Khoravar who follow this philosophy often become bards, diplomats, mediators, or translators. Others, fascinated by their distant connection to the Fey, seek to build bridges between the Material Plane and the Feywild of Thelanis. These Khoravar often become druids or warlocks with archfey patrons.",
		CreatureType: "Humanoid",
		Size:         "Medium (about 4–6 feet tall) or Small (about 2–4 feet tall), chosen when you select this species",
		BaseSpeed:    30,
		Languages:    []string{"Common", "Elvish"},
		Traits:       khoravarTraits(),
	}
}

func khoravarTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 60 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "60",
		},
		{
			Name:        "Fey Ancestry",
			Description: "You have Advantage on saving throws you make to avoid or end the Charmed condition.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:           "Lethargy Resilience",
			Description:    "When you fail a saving throw to avoid or end the Unconscious condition, you can succeed instead. Once you use this trait, you can't do so again until you finish 1d4 Long Rests.",
			LevelReq:       1,
			ActionType:     "Passive",
			UsesPerRest:    "1",
			ResetCondition: "1d4 Long Rests",
		},
		{
			Name:        "Skill Versatility",
			Description: "You gain proficiency in one skill or with one tool of your choice. Whenever you finish a Long Rest, you can replace it with another skill or tool proficiency.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
	}
}
