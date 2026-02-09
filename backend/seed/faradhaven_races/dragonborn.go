package faradhaven_races

// Dragonborn returns the Dragonborn race seed
func Dragonborn() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Dragonborn",
		PhotoURL:     "https://photos-for-apps.s3.us-east-2.amazonaws.com/dragonborn.jpg",
		Description:  "The ancestors of dragonborn hatched from the eggs of chromatic and metallic dragons. One story holds that these eggs were blessed by the dragon gods Bahamut and Tiamat, who wanted to populate the multiverse with people created in their image. Another story claims that dragons created the first dragonborn without the gods' blessings. Whatever their origin, dragonborn have made homes for themselves on the Material Plane. Dragonborn look like wingless, bipedal dragons - scaly, bright-eyed, and thick-boned with horns on their heads - and their coloration and other features are reminiscent of their draconic ancestors.",
		CreatureType: "Humanoid",
		Size:         "Medium (about 5–7 feet tall)",
		BaseSpeed:    30,
		AbilityScoreBonuses: map[string]int{
			"strength": 2,
			"charisma": 1,
		},
		Languages:      []string{"Common", "Draconic"},
		Traits:         dragonbornTraits(),
		ComponentNames: []string{"Ignis", "Fulgur", "Cool", "Wither", "Aer", "Nova", "Channel"},
	}
}

func dragonbornTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Draconic Ancestry",
			Description: "Your lineage stems from a dragon progenitor. Choose the kind of dragon from the Draconic Ancestors table. Your choice affects your Breath Weapon and Damage Resistance traits as well as your appearance.",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     dragonbornAncestryOptions(),
		},
		{
			Name:           "Breath Weapon",
			Description:    "When you take the Attack action on your turn, you can replace one of your attacks with an exhalation of magical energy using your ancestral element component in either Nova (15-foot Cone) or Channel (30-foot Line that is 5 feet wide) shape, chosen each time. Each creature in that area must make a Dexterity saving throw (DC 8 plus your Constitution modifier and Proficiency Bonus). On a failed save, a creature takes 1d10 damage of the type determined by your Draconic Ancestry trait. On a successful save, a creature takes half as much damage. This damage increases by 1d10 when you reach character levels 5 (2d10), 11 (3d10), and 17 (4d10). You can use this Breath Weapon a number of times equal to your Proficiency Bonus, and you regain all expended uses when you finish a Long Rest.",
			LevelReq:       1,
			ActionType:     "Action",
			UsesPerRest:    "Proficiency Bonus",
			ResetCondition: "Long Rest",
			AreaOfEffect:   "15-foot Cone or 30-foot Line (5 ft wide)",
			SaveAbility:    "DEX",
		},
		{
			Name:        "Damage Resistance",
			Description: "You have Resistance to the damage type determined by your Draconic Ancestry trait.",
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
			Name:           "Draconic Flight",
			Description:    "When you reach character level 5, you can channel draconic magic to give yourself temporary flight. As a Bonus Action, you sprout spectral wings on your back that last for 10 minutes or until you retract the wings (no action required) or have the Incapacitated condition. During that time, you have a Fly Speed equal to your Speed. Your wings appear to be made of the same energy as your Breath Weapon. Once you use this trait, you can't use it again until you finish a Long Rest.",
			LevelReq:       5,
			ActionType:     "Bonus Action",
			UsesPerRest:    "1",
			ResetCondition: "Long Rest",
		},
	}
}

func dragonbornAncestryOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{Name: "Black", Description: "Wither component (Acid damage)."},
		{Name: "Blue", Description: "Fulgur component (Lightning damage)."},
		{Name: "Brass", Description: "Ignis component (Fire damage)."},
		{Name: "Bronze", Description: "Fulgur component (Lightning damage)."},
		{Name: "Copper", Description: "Wither component (Acid damage)."},
		{Name: "Gold", Description: "Ignis component (Fire damage)."},
		{Name: "Green", Description: "Wither + Aer components (Poison damage)."},
		{Name: "Red", Description: "Ignis component (Fire damage)."},
		{Name: "Silver", Description: "Cool component (Cold damage)."},
		{Name: "White", Description: "Cool component (Cold damage)."},
	}
}
