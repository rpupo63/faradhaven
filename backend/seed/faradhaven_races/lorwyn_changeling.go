package faradhaven_races

// LorwynChangeling returns the Lorwyn Changeling race seed
func LorwynChangeling() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:               "Lorwyn Changeling",
		PhotoURL:           "https://photos-for-apps.s3.us-east-2.amazonaws.com/lorwyn_changeling.jpg",
		Description:        "The changelings of Lorwyn are charismatic shapeshifters able to crudely mimic the forms of creatures and plants. Unlike most shapeshifters, however, their innate powers don't allow them to truly disguise themselves; rather, their transformations are purely surface level.\n\nNo matter their form, Lorwyn changelings keep their key traits: blue-green skin, tufts of tentacle-like fur, and bulbous yellow eyes with slitted pupils. The true changeling identity of a bug-eyed tree, teal elk, or furry boulder is always easy to spot. These counterfeits are so plain that many Lorwyn denizens find the effect inexplicably charming.\n\nChangelings from Lorwyn live in a vast, mystical cave called Velis Vel Grotto. Once per year, many changelings make a pilgrimage to the grotto to contemplate their past deeds.",
		CreatureType:       "Fey",
		Size:               "Medium (about 4–7 feet tall) or Small (about 2–4 feet tall), chosen when you select this species",
		BaseSpeed:          30,
		Languages:          []string{"Common", "Sylvan"},
		BonusLanguageCount: 1,
		Traits:             lorwynChangelingTraits(),
	}
}

func lorwynChangelingTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Shape Self",
			Description: "As an action, you can reshape your body to a two-legged Humanoid shape or to a four-legged Beast shape. While you have a Humanoid shape, you can wear clothing and armor made for a Humanoid of your size.",
			LevelReq:    1,
			ActionType:  "Action",
		},
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 120 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "120",
		},
		{
			Name:        "Delightful Imitator",
			Description: "You have proficiency in the Performance skill.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Unpredictable Movement",
			Description: "When you roll Initiative and don't have Disadvantage on that roll, you can take the Dash action as a Reaction.",
			LevelReq:    1,
			ActionType:  "Reaction",
		},
	}
}
