package faradhaven_races

// Dwarf returns the Dwarf race seed
func Dwarf() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Dwarf",
		PhotoURL:     "https://photos-for-apps.s3.us-east-2.amazonaws.com/dwarf.jpg",
		Description:  "Dwarves were raised from the earth in the elder days by a deity of the forge. Called by various names on different worlds - Moradin, Reorx, and others - that god gave dwarves an affinity for stone and metal and for living underground. The god also made them resilient like the mountains, with a life span of about 350 years. Squat and often bearded, the original dwarves carved cities and strongholds into mountainsides and under the earth. Their oldest legends tell of conflicts with the monsters of mountaintops and the Underdark, whether those monsters were towering giants or subterranean horrors. Inspired by those tales, dwarves of any culture often sing of valorous deeds - especially of the little overcoming the mighty. On some worlds in the multiverse, the first settlements of dwarves were built in hills or mountains, and the families who trace their ancestry to those settlements call themselves hill dwarves or mountain dwarves, respectively.",
		CreatureType: "Humanoid",
		Size:         "Medium (about 4-5 feet tall)",
		BaseSpeed:    30,
		AbilityScoreBonuses: map[string]int{
			"constitution": 2,
			"wisdom":       1,
		},
		Traits:       dwarfTraits(),
	}
}

func dwarfTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 120 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "120",
		},
		{
			Name:        "Dwarven Resilience",
			Description: "You have Resistance to Poison damage. You also have Advantage on saving throws you make to avoid or end the Poisoned condition.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Dwarven Toughness",
			Description: "Your Hit Point maximum increases by 1, and it increases by 1 again whenever you gain a level.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:           "Stonecunning",
			Description:    "As a Bonus Action, you gain Tremorsense with a range of 60 feet for 10 minutes. You must be on a stone surface or touching a stone surface to use this Tremorsense. The stone can be natural or worked. You can use this Bonus Action a number of times equal to your Proficiency Bonus, and you regain all expended uses when you finish a Long Rest.",
			LevelReq:       1,
			ActionType:     "Bonus Action",
			UsesPerRest:    "Proficiency Bonus",
			ResetCondition: "Long Rest",
			RangeValue:     "60",
		},
	}
}
