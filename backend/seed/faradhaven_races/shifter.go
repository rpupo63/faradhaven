package faradhaven_races

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/seedmedia"
)

// Shifter returns the Shifter race seed
func Shifter() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:               "Shifter",
		PhotoURL:           seedmedia.URL("shifter.jpg"),
		Description:        "Shifters—sometimes called \"weretouched\"—descend from people who contracted full or partial lycanthropy. Humanoids with a bestial aspect, shifters can't change shape fully, but they can enhance their animalistic features temporarily in a process they call \"shifting.\" Shifters resemble humans in height and build but are typically more lithe and flexible. Their facial features have a bestial cast, often with large eyes and pointed ears; most shifters also possess prominent canine teeth. They grow fur-like hair on nearly every part of their bodies. While a shifter's appearance might remind an onlooker of an animal, the shifter remains clearly identifiable as a Humanoid even when at their most feral.",
		CreatureType:       "Humanoid",
		Size:               "Medium (about 4–7 feet tall) or Small (about 2–4 feet tall), chosen when you select this species",
		BaseSpeed:          30,
		Languages:          []string{"Common"},
		BonusLanguageCount: 1,
		Traits:             shifterTraits(),
	}
}

func shifterTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Bestial Instincts",
			Description: "Channeling the beast within, you gain proficiency in one of the following skills of your choice: Acrobatics, Athletics, Intimidation, or Survival.",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     shifterBestialInstinctsOptions(),
		},
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 60 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "60",
		},
		{
			Name:           "Shifting",
			Description:    "As a Bonus Action, you can shape-shift to assume a more bestial appearance. This transformation lasts for 1 minute or until you revert to your normal appearance as a Bonus Action. When you shift, you gain Temporary Hit Points equal to 2 times your Proficiency Bonus. You can shift a number of times equal to your Proficiency Bonus, and you regain all expended uses when you finish a Long Rest. Whenever you shift, you gain the benefit of one of the following options (choose when you select this species):",
			LevelReq:       1,
			ActionType:     "Bonus Action",
			UsesPerRest:    "Proficiency Bonus",
			ResetCondition: "Long Rest",
			Options:        shifterShiftingOptions(),
		},
	}
}

func shifterBestialInstinctsOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{Name: "Acrobatics", Description: "You gain proficiency in Acrobatics."},
		{Name: "Athletics", Description: "You gain proficiency in Athletics."},
		{Name: "Intimidation", Description: "You gain proficiency in Intimidation."},
		{Name: "Survival", Description: "You gain proficiency in Survival."},
	}
}

func shifterShiftingOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:        "Beasthide",
			Description: "You gain 1d6 additional Temporary Hit Points. While shifted, you have a +1 bonus to your Armor Class.",
		},
		{
			Name:        "Longtooth",
			Description: "When you shift and as a Bonus Action on your other turns while shifted, you can use your elongated fangs to make an Unarmed Strike. If you hit with this Unarmed Strike and deal damage, you can deal Piercing damage equal to 1d6 plus your Strength modifier, instead of the normal damage of an Unarmed Strike.",
		},
		{
			Name:        "Swiftstride",
			Description: "While you are shifted, your Speed increases by 10 feet. Additionally, you can move up to 10 feet as a Reaction when a creature ends its turn within 5 feet of you. This reactive movement doesn't provoke Opportunity Attack action.",
		},
		{
			Name:        "Wildhunt",
			Description: "While shifted, you have Advantage on Wisdom checks. Additionally, no creature within 30 feet of you can have Advantage on an attack roll against you unless you have the Incapacitated condition.",
		},
	}
}
