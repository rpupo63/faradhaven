package faradhaven_races

// Gnome returns the Gnome race seed
func Gnome() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Gnome",
		PhotoURL:     "https://photos-for-apps.s3.us-east-2.amazonaws.com/gnome.jpg",
		Description:  "Gnomes are magical folk created by gods of invention, illusions, and life underground. The earliest gnomes were seldom seen by other folk due to their secretive nature and their propensity for living in forests and burrows. What they lacked in size, they made up for in cleverness. They confounded predators with traps and labyrinthine tunnels. They also learned magic from gods like Garl Glittergold, Baervan Wildwanderer, and Baravar Cloakshadow, who visited them in disguise. That magic eventually created the lineages of forest gnomes and rock gnomes. Gnomes are petite folk with big eyes and pointed ears, who live around 425 years. Many gnomes like the feeling of a roof over their head, even if that \"roof\" is nothing more than a hat.",
		CreatureType: "Humanoid",
		Size:         "Small (about 3-4 feet tall)",
		BaseSpeed:    30,
		AbilityScoreBonuses: map[string]int{
			"intelligence": 2,
		},
		Languages:      []string{"Common", "Gnomish"},
		Traits: gnomeTraits(),
	}
}

func gnomeTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 60 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "60",
		},
		{
			Name:        "Gnomish Cunning",
			Description: "You have Advantage on Intelligence, Wisdom, and Charisma saving throws.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Gnomish Lineage",
			Description: "Choose forest or rock gnomish heritage.",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     gnomeLineageOptions(),
		},
	}
}

func gnomeLineageOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:                "Forest Gnome",
			Description:         "+1 Dexterity. You are at home among beasts, small illusions, and natural trickery; express such effects through your class’s spell components.",
			AbilityScoreBonuses: map[string]int{"dexterity": 1},
		},
		{
			Name:                "Rock Gnome",
			Description:         "+1 Constitution. You excel at tinkering, repair, and clever devices; express gadgetry and minor wonders through your class’s spell components.",
			AbilityScoreBonuses: map[string]int{"constitution": 1},
		},
	}
}
