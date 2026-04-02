package faradhaven_races

// Tiefling returns the Tiefling race seed
func Tiefling() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:        "Tiefling",
		PhotoURL:    "https://photos-for-apps.s3.us-east-2.amazonaws.com/tiefling.jpg",
		Description: "Tieflings are either born in the Lower Planes or have fiendish ancestors who originated there. A tiefling (pronounced TEE-fling) is linked by blood to a devil, a demon, or some other Fiend. This connection to the Lower Planes is the tiefling's fiendish legacy, which comes with the promise of power yet has no effect on the tiefling's moral outlook. A tiefling chooses whether to embrace or lament their fiendish legacy. Abyssal: The entropy of the Abyss, the chaos of Pandemonium, and the despair of Carceri call to tieflings who have the abyssal legacy. Horns, fur, tusks, and peculiar scents are common physical features of such tieflings, most of whom have the blood of demons coursing through their veins. Chthonic: Tieflings who have the chthonic legacy feel not only the tug of Carceri but also the greed of Gehenna and the gloom of Hades. Some of these tieflings look cadaverous. Others possess the unearthly beauty of a succubus, or they have physical features in common with a night hag, a yugoloth, or some other Neutral Evil fiendish ancestor. Infernal: The infernal legacy connects tieflings not only to Gehenna but also the Nine Hells and the raging battlefields of Acheron. Horns, spines, tails, golden eyes, and a faint odor of sulfur or smoke are common physical features of such tieflings, most of whom trace their ancestry to devils.",
		CreatureType: "Humanoid",
		Size:         "Medium (about 4–7 feet tall) or Small (about 3–4 feet tall), chosen when you select this species",
		BaseSpeed:    30,
		AbilityScoreBonuses: map[string]int{
			"charisma":     2,
			"intelligence": 1,
		},
		Languages: []string{"Common", "Infernal"},
		Traits:    tieflingTraits(),
	}
}

func tieflingTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 60 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "60",
		},
		{
			Name:        "Fiendish Legacy",
			Description: "Your fiendish blood confers elemental resistance. Choose one legacy.",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     tieflingFiendishLegacyOptions(),
		},
	}
}

func tieflingFiendishLegacyOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:        "Abyssal",
			Description: "You have Resistance to Poison damage.",
		},
		{
			Name:        "Chthonic",
			Description: "You have Resistance to Necrotic damage.",
		},
		{
			Name:        "Infernal",
			Description: "You have Resistance to Fire damage.",
		},
	}
}
