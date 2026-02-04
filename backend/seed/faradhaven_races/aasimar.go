package faradhaven_races

// Aasimar returns the Aasimar race seed
func Aasimar() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:           "Aasimar",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/aasimar.jpg",
		Description:    "Aasimar (pronounced AH-sih-mar) are mortals who carry a spark of the Upper Planes within their souls. Whether descended from an angelic being or infused with celestial power, they can fan that spark to bring light, healing, and heavenly fury. Aasimar can arise among any population of mortals. They resemble their parents, but they live for up to 160 years and have features that hint at their celestial heritage, such as metallic freckles, luminous eyes, a halo, or the skin color of an angel (silver, opalescent green, or coppery red). These features start subtle and become obvious when the aasimar learns to reveal their full celestial nature.",
		CreatureType:   "Humanoid",
		Size:           "Medium (4–7 ft) or Small (2–4 ft)",
		BaseSpeed:      30,
		Traits:         aasimarTraits(),
		ComponentNames: []string{"Lux", "Self", "Nova", "Umbra", "Fear", "Mend"},
	}
}

func aasimarTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Celestial Resistance",
			Description: "You have Resistance to Necrotic damage and Radiant damage.",
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
			Name:           "Healing Hands",
			Description:    "As a Magic action, you touch a creature and roll a number of d4s equal to your Proficiency Bonus. The creature regains a number of Hit Points equal to the total rolled. Once you use this trait, you can't use it again until you finish a Long Rest.",
			LevelReq:       1,
			ActionType:     "Magic Action",
			UsesPerRest:    "Proficiency Bonus",
			ResetCondition: "Long Rest",
		},
		{
			Name:        "Light Bearer",
			Description: "You innately know the Lux + Self component combination, allowing you to cause your body to shed bright light. Charisma is your spellcasting ability for it.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:           "Celestial Revelation",
			Description:    "When you reach character level 3, you can transform as a Bonus Action using one of the options below (choose the option each time you transform). The transformation lasts for 1 minute or until you end it (no action required). Once you transform, you can't do so again until you finish a Long Rest.\n\nOnce on each of your turns before the transformation ends, you can deal extra damage to one target when you deal damage to it with an attack or a spell. The extra damage equals your Proficiency Bonus, and the extra damage's type is either Necrotic for Necrotic Shroud or Radiant for Heavenly Wings and Inner Radiance.",
			LevelReq:       3,
			ActionType:     "Bonus Action",
			UsesPerRest:    "1",
			ResetCondition: "Long Rest",
			Options:        aasimarRevelationOptions(),
		},
	}
}

func aasimarRevelationOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:        "Heavenly Wings",
			Description: "Two spectral wings sprout from your back temporarily. Until the transformation ends, you have a Fly Speed equal to your Speed.",
		},
		{
			Name:        "Inner Radiance",
			Description: "Searing light temporarily radiates from your eyes and mouth, manifesting as Lux + Nova. For the duration, you shed Bright Light in a 10-foot radius and Dim Light for an additional 10 feet, and at the end of each of your turns, each creature within 10 feet of you takes Radiant damage equal to your Proficiency Bonus.",
		},
		{
			Name:        "Necrotic Shroud",
			Description: "Your eyes briefly become pools of darkness, manifesting as Umbra + Fear + Nova, and flightless wings sprout from your back temporarily. Creatures other than your allies within 10 feet of you must succeed on a Charisma saving throw (DC 8 plus your Charisma modifier and Proficiency Bonus) or have the Frightened condition until the end of your next turn.",
		},
	}
}
