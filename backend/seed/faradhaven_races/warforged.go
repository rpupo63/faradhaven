package faradhaven_races

// Warforged returns the Warforged race seed
func Warforged() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Warforged",
		PhotoURL:     "https://photos-for-apps.s3.us-east-2.amazonaws.com/warforged.jpg",
		Description:  "Warforged are mechanical beings built as weapons to fight in the Last War. An unexpected breakthrough produced sentient beings made from wood and metal that nevertheless can feel pain and emotion. Warforged comprise a blend of organic and inorganic materials. Rootlike cords infused with alchemical fluids serve as their muscles, wrapped around a framework of steel, darkwood, or stone. Armored plates form a protective outer shell and reinforce joints. The more a warforged cultivates their individuality, the more likely they are to modify their body, seeking out an artificer to customize the look of their face, limbs, and plating.",
		CreatureType: "Construct",
		Size:         "Medium (about 6–8 feet tall) or Small (about 3–4 feet tall), chosen when you select this species",
		BaseSpeed:    30,
		Traits:       warforgedTraits(),
	}
}

func warforgedTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Construct Resilience",
			Description: "You have Resistance to Poison damage. You also have Advantage on saving throws to avoid or end the Poisoned condition.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Integrated Protection",
			Description: "You gain a +1 bonus to your Armor Class. In addition, armor you have donned can't be removed against your will while you're alive.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Sentry's Rest",
			Description: "You don't need to sleep, and magic can't put you to sleep. You can finish a Long Rest in 6 hours if you spend those hours in an inactive, motionless state. During this time, you appear inert but remain conscious.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Specialized Design",
			Description: "You gain one skill proficiency and one tool proficiency of your choice.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Tireless",
			Description: "You don't gain Exhaustion levels from dehydration, malnutrition, or suffocation.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
	}
}
