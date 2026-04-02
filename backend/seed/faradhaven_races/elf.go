package faradhaven_races

// Elf returns the Elf race seed
func Elf() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Elf",
		PhotoURL:     "https://photos-for-apps.s3.us-east-2.amazonaws.com/elf.jpg",
		Description:  "Created by the god Corellon, the first elves could change their forms at will. They lost this ability when Corellon cursed them for plotting with the deity Lolth, who tried and failed to usurp Corellon's dominion. When Lolth was cast into the Abyss, most elves renounced her and earned Corellon's forgiveness, but that which Corellon had taken from them was lost forever. No longer able to shape-shift at will, the elves retreated to the Feywild, where their sorrow was deepened by that plane's influence. Over time, curiosity led many of them to explore other planes of existence, including worlds in the Material Plane. Elves have pointed ears and lack facial and body hair. They live for around 750 years, and they don't sleep but instead enter a trance when they need to rest. In that state, they remain aware of their surroundings while immersing themselves in memories and meditations. An environment subtly transforms elves after they inhabit it for a millennium or more, and it grants them certain kinds of magic. Drow, high elves, and wood elves are examples of elves who have been transformed thus.",
		CreatureType: "Humanoid",
		Size:         "Medium (about 5-6 feet tall)",
		BaseSpeed:    30,
		AbilityScoreBonuses: map[string]int{
			"dexterity": 2,
		},
		Languages:      []string{"Common", "Elvish"},
		Traits: elfTraits(),
	}
}

func elfTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 60 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "60",
		},
		{
			Name:        "Elven Lineage",
			Description: "Choose a lineage reflecting your culture and environment. You gain the level 1 benefit of that lineage. Higher-level magical effects use your class’s spell components.",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     elfLineageOptions(),
		},
		{
			Name:        "Fey Ancestry",
			Description: "You have Advantage on saving throws you make to avoid or end the Charmed condition.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Keen Senses",
			Description: "You have proficiency in the Insight, Perception, or Survival skill.",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     elfKeenSensesOptions(),
		},
		{
			Name:        "Trance",
			Description: "You don't need to sleep, and magic can't put you to sleep. You can finish a Long Rest in 4 hours if you spend those hours in a trancelike meditation, during which you retain consciousness.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
	}
}

func elfLineageOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:                "Drow",
			Description:         "The range of your Darkvision increases to 120 feet. +1 Charisma.",
			AbilityScoreBonuses: map[string]int{"charisma": 1},
		},
		{
			Name:                "High Elf",
			Description:         "+1 Intelligence. You favor arcane study; express wizardly flourishes through your class’s spell components.",
			AbilityScoreBonuses: map[string]int{"intelligence": 1},
		},
		{
			Name:                "Wood Elf",
			Description:         "Your Speed increases to 35 feet. +1 Wisdom.",
			AbilityScoreBonuses: map[string]int{"wisdom": 1},
		},
		{
			Name:                "Lorwyn Elf",
			Description:         "+1 Wisdom. You are attuned to Lorwyn’s wilds; primal effects use your class’s spell components.",
			AbilityScoreBonuses: map[string]int{"wisdom": 1},
		},
		{
			Name:                "Shadowmoor Elf",
			Description:         "The range of your Darkvision increases to 120 feet. +1 Charisma.",
			AbilityScoreBonuses: map[string]int{"charisma": 1},
		},
	}
}

func elfKeenSensesOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{Name: "Insight", Description: "You have proficiency in the Insight skill."},
		{Name: "Perception", Description: "You have proficiency in the Perception skill."},
		{Name: "Survival", Description: "You have proficiency in the Survival skill."},
	}
}
