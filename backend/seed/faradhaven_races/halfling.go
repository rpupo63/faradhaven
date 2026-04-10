package faradhaven_races

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/seedmedia"
)

// Halfling returns the Halfling race seed
func Halfling() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Halfling",
		PhotoURL:     seedmedia.URL("halfling.jpg"),
		Description:  "Cherished and guided by gods who value life, home, and hearth, halflings gravitate toward bucolic havens where family and community help shape their lives. That said, many halflings possess a brave and adventurous spirit that leads them on journeys of discovery, affording them the chance to explore a bigger world and make new friends along the way. Their size—similar to that of a human child—helps them pass through crowds unnoticed and slip through tight spaces. Anyone who has spent time around halflings, particularly halfling adventurers, has likely witnessed the storied \"luck of the halflings\" in action. When a halfling is in mortal danger, an unseen force seems to intervene on the halfling's behalf. Many halflings believe in the power of luck, and they attribute their unusual gift to one or more of their benevolent gods, including Yondalla, Brandobaris, and Charmalaine. The same gift might contribute to their robust life spans (about 150 years). Halfling communities come in all varieties: sequestered shires tucked away in unspoiled parts of the world, crime syndicates, or territorial mobs. Halflings who prefer to live underground are sometimes called strongheart halflings or stouts. Nomadic halflings, as well as those who live among humans and other tall folk, are sometimes called lightfoot halflings or tallfellows.",
		CreatureType: "Humanoid",
		Size:         "Small (about 2–3 feet tall)",
		BaseSpeed:    30,
		AbilityScoreBonuses: map[string]int{
			"dexterity": 2,
			"charisma":  1,
		},
		Languages: []string{"Common", "Halfling"},
		Traits:    halflingTraits(),
	}
}

func halflingTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Brave",
			Description: "You have Advantage on saving throws you make to avoid or end the Frightened condition.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Halfling Nimbleness",
			Description: "You can move through the space of any creature that is a size larger than you, but you can't stop in the same space.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Luck",
			Description: "When you roll a 1 on the d20 of a D20 Test, you can reroll the die, and you must use the new roll.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Naturally Stealthy",
			Description: "You can take the Hide action even when you are obscured only by a creature that is at least one size larger than you.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
	}
}
