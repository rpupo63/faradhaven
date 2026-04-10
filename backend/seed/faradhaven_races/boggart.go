package faradhaven_races

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/seedmedia"
)

// Boggart returns the Boggart race seed
func Boggart() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Boggart",
		PhotoURL:     seedmedia.URL("boggart.jpg"),
		Description:  "Boggarts are small, squat goblinoids from the realm of Lorwyn-Shadowmoor. They have bestial features such as horns and animal-like snouts, and their appearance varies widely: one might resemble a hedgehog with spiky fur, while another has the snout and fleshy ears of a swine. Boggarts love crafting potions and often gravitate toward specific areas of expertise.\n\nIn Lorwyn, boggarts are born into communal warrens where laws and hierarchies are loose suggestions. The oldest and most powerful boggarts, known as aunties, serve as respected leaders who keep the peace. Lorwyn boggarts value sharing knowledge and past experiences, and many are willing to brave great dangers in pursuit of new experiences.\n\nIn Shadowmoor, boggarts are sharper-featured and may have extra sets of horns or body modifications like riveted armored plates. Their society is chaotic and decentralized, with small communities in isolated or dangerous places. Aunties in Shadowmoor wander the land, greeting those they meet with arbitrary tests, boons, or curses.",
		CreatureType: "Humanoid",
		Size:         "Small (about 2-4 feet tall)",
		BaseSpeed:    30,
		Languages:    []string{"Common", "Goblin"},
		Traits:       boggartTraits(),
	}
}

func boggartTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Goblinoid Heritage",
			Description: "You are also considered a goblinoid for any prerequisite or effect that requires you to be a goblinoid.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 60 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "60",
		},
		{
			Name:        "Fey Ancestry",
			Description: "You have Advantage on saving throws you make to avoid or end the Charmed condition on yourself.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:           "Fury of the Small",
			Description:    "When you damage a creature with an attack or a spell and the creature's size is larger than yours, you can cause the attack or spell to deal extra damage to the creature. The extra damage equals your Proficiency Bonus. You can use this trait a number of times equal to your Proficiency Bonus, and you regain all expended uses when you finish a Long Rest. You can use it no more than once per turn.",
			LevelReq:       1,
			ActionType:     "Passive",
			UsesPerRest:    "Proficiency Bonus",
			ResetCondition: "Long Rest",
		},
		{
			Name:        "Nimble Escape",
			Description: "You can take the Disengage or Hide action as a Bonus Action on each of your turns.",
			LevelReq:    1,
			ActionType:  "Bonus Action",
		},
	}
}
