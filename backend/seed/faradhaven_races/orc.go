package faradhaven_races

// Orc returns the Orc race seed
func Orc() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Orc",
		PhotoURL:     "https://photos-for-apps.s3.us-east-2.amazonaws.com/orc.jpg",
		Description:  "Orcs trace their creation to Gruumsh, a powerful god who roamed the wide open spaces of the Material Plane. Gruumsh equipped his children with gifts to help them wander great plains, vast caverns, and churning seas and to face the monsters that lurk there. Even when they turn their devotion to other gods, orcs retain Gruumsh's gifts: endurance, determination, and the ability to see in darkness.\n\nOrcs are, on average, tall and broad. They have gray skin, ears that are sharply pointed, and prominent lower canines that resemble small tusks. Orc youths on some worlds are told about their ancestors' great travels and travails. Inspired by those tales, many of those orcs wonder when Gruumsh will call on them to match the heroic deeds of old and if they will prove worthy of his favor. Other orcs are happy to leave old tales in the past and find their own way.",
		CreatureType: "Humanoid",
		Size:         "Medium (about 6-7 feet tall)",
		BaseSpeed:    30,
		Traits:       orcTraits(),
	}
}

func orcTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:           "Adrenaline Rush",
			Description:    "You can take the Dash action as a Bonus Action. When you do so, you gain a number of Temporary Hit Points equal to your Proficiency Bonus.\n\nYou can use this trait a number of times equal to your Proficiency Bonus, and you regain all expended uses when you finish a Short or Long Rest.",
			LevelReq:       1,
			ActionType:     "Bonus Action",
			UsesPerRest:    "Proficiency Bonus",
			ResetCondition: "Short or Long Rest",
		},
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 120 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "120",
		},
		{
			Name:           "Relentless Endurance",
			Description:    "When you are reduced to 0 Hit Points but not killed outright, you can drop to 1 Hit Point instead. Once you use this trait, you can't do so again until you finish a Long Rest.",
			LevelReq:       1,
			ActionType:     "Reaction",
			UsesPerRest:    "1",
			ResetCondition: "Long Rest",
		},
	}
}
