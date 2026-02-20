package faradhaven_races

// Goliath returns the Goliath race seed
func Goliath() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:           "Goliath",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/goliath.jpg",
		Description:    "Towering over most folk, goliaths are distant descendants of giants. Each goliath bears the favors of the first giants - favors that manifest in various supernatural boons, including the ability to quickly grow and temporarily approach the height of goliaths' gigantic kin.\n\nGoliaths have physical characteristics that are reminiscent of the giants in their family lines. For example, some goliaths look like stone giants, while others resemble fire giants. Whatever giants they count as kin, goliaths have forged their own path in the multiverse - unencumbered by the internecine conflicts that have ravaged giantkind for ages - and seek heights above those reached by their ancestors.",
		CreatureType:   "Humanoid",
		Size:           "Medium (about 7-8 feet tall)",
		BaseSpeed:      35,
		Languages:      []string{"Common", "Giant"},
		Traits:         goliathTraits(),
		ComponentNames: []string{"Spatium", "Ignis", "Create", "Decrease", "Crush", "Terra", "Strong", "Fulgur", "Increase", "Self"},
	}
}

func goliathTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:           "Giant Ancestry",
			Description:    "You are descended from Giants. Choose one of the following benefits - a supernatural boon from your ancestry; you can use the chosen benefit a number of times equal to your Proficiency Bonus, and you regain all expended uses when you finish a Long Rest.",
			LevelReq:       1,
			ActionType:     "Varies",
			UsesPerRest:    "Proficiency Bonus",
			ResetCondition: "Long Rest",
			Options:        goliathGiantAncestryOptions(),
		},
		{
			Name:           "Large Form",
			Description:    "Starting at character level 5, you can use Increase + Self to change your size to Large as a Bonus Action if you're in a big enough space. This transformation lasts for 10 minutes or until you end it (no action required). For that duration, you have Advantage on Strength checks, and your Speed increases by 10 feet. Once you use this trait, you can't use it again until you finish a Long Rest.",
			LevelReq:       5,
			ActionType:     "Bonus Action",
			UsesPerRest:    "1",
			ResetCondition: "Long Rest",
		},
		{
			Name:        "Powerful Build",
			Description: "You have Advantage on any ability check you make to end the Grappled condition. You also count as one size larger when determining your carrying capacity.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
	}
}

func goliathGiantAncestryOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:        "Cloud's Jaunt (Cloud Giant)",
			Description: "As a Bonus Action, you use Spatium to magically displace yourself up to 30 feet to an unoccupied space you can see.",
		},
		{
			Name:        "Fire's Burn (Fire Giant)",
			Description: "When you hit a target with an attack roll and deal damage to it, you can also use Ignis + Create to deal 1d10 Fire damage to that target.",
		},
		{
			Name:        "Frost's Chill (Frost Giant)",
			Description: "When you hit a target with an attack roll and deal damage to it, you can also use Decrease + Create to deal 1d6 Cold damage to that target and reduce its Speed by 10 feet until the start of your next turn.",
		},
		{
			Name:        "Hill's Tumble (Hill Giant)",
			Description: "When you hit a Large or smaller creature with an attack roll and deal damage to it, you can use Crush + Create to give that target the Prone condition.",
		},
		{
			Name:        "Stone's Endurance (Stone Giant)",
			Description: "When you take damage, you can take a Reaction to use Terra + Strong. Roll 1d12, add your Constitution modifier to the number rolled, and reduce the damage by that total.",
		},
		{
			Name:        "Storm's Thunder (Storm Giant)",
			Description: "When you take damage from a creature within 60 feet of you, you can take a Reaction to use Fulgur to deal 1d8 Thunder damage to that creature.",
		},
	}
}
