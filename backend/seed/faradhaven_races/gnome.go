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
		Traits:         gnomeTraits(),
		ComponentNames: []string{"Phantom", "Vita", "Link", "Mend", "Arcanum", "Self"},
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
			Description: "You are part of a lineage that grants you supernatural abilities. Choose one of the following options; whichever one you choose, Intelligence, Wisdom, or Charisma is your spellcasting ability for the spells you cast with this trait (choose the ability when you select the lineage).",
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
			Description:         "You know the Phantom component (creating illusory images or sounds). You also always have Vita + Link prepared (communicating with beasts). You can cast it without a spell slot a number of times equal to your Proficiency Bonus, and you regain all expended uses when you finish a Long Rest. You can also use any spell slots you have to cast the spell.",
			AbilityScoreBonuses: map[string]int{"dexterity": 1},
		},
		{
			Name:                "Rock Gnome",
			Description:         "You know the Mend component (repairing objects) and Arcanum + Self combination (minor magical tricks). In addition, you can spend 10 minutes using Arcanum + Self to create a Tiny clockwork device (AC 5, 1 HP), such as a toy, fire starter, or music box. When you create the device, you determine its function by choosing one effect from the combination; the device produces that effect whenever you or another creature takes a Bonus Action to activate it with a touch. If the chosen effect has options within it, you choose one of those options for the device when you create it. For example, if you choose the ignite-extinguish effect, you determine whether the device ignites or extinguishes fire; the device doesn't do both. You can have three such devices in existence at a time, and each falls apart 8 hours after its creation or when you dismantle it with a touch as a Utilize action.",
			AbilityScoreBonuses: map[string]int{"constitution": 1},
		},
	}
}
