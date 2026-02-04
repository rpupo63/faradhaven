package faradhaven_races

// Rimekin returns the Rimekin race seed
func Rimekin() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:           "Rimekin",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/rimekin.jpg",
		Description:    "Rimekin hail from both Lorwyn and Shadowmoor, though the first rimekin arose from flamekin during the Phyrexian invasion. These flamekin approached their problems with cold logic and rejected reactionary responses. As a result, the magical flames that engulfed their bodies took on a frigid air, and they became rimekin.\n\nLike flamekin, rimekin possess innate magic, but the flames they conjure burn icy blue rather than red hot. Further, these \"flames\" emanate a chilling cold rather than blazing heat. This effect extends, superficially, to the items rimekin touch. Any armor worn or weapons wielded by rimekin become coated in layers of spiky yet harmless hoarfrost.",
		CreatureType:   "Humanoid",
		Size:           "Medium (about 4–7 feet tall) or Small (about 2–4 feet tall), chosen when you select this species",
		BaseSpeed:      30,
		Traits:         rimekinTraits(),
		ComponentNames: []string{"Cool", "Beam", "Aqua", "Projectile", "Imbue"},
	}
}

func rimekinTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 60 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "60",
		},
		{
			Name:        "Cold Resistance",
			Description: "You have resistance to Cold damage.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:           "Cold Fire Magic",
			Description:    "You know the Cool + Beam component combination (a ray of freezing cold). When you reach character levels 3 and 5, you learn Aqua + Cool + Projectile (an exploding shard of ice) and Cool + Imbue (a blade of frigid cold), respectively. You always have those component combinations prepared. You can cast each once without a spell slot, and you regain the ability to cast these spells in this way when you finish a Long Rest. You can also cast the spells using any spell slots you have of the appropriate level. Intelligence, Wisdom, or Charisma is your spellcasting ability for these spells (choose the ability when you select this species).",
			LevelReq:       1,
			ActionType:     "Passive",
			UsesPerRest:    "1 per spell",
			ResetCondition: "Long Rest",
			Options:        rimekinSpellcastingOptions(),
		},
	}
}

func rimekinSpellcastingOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:        "Intelligence",
			Description: "Intelligence is your spellcasting ability for Cool + Beam, Aqua + Cool + Projectile, and Cool + Imbue when you cast them with this trait.",
		},
		{
			Name:        "Wisdom",
			Description: "Wisdom is your spellcasting ability for Cool + Beam, Aqua + Cool + Projectile, and Cool + Imbue when you cast them with this trait.",
		},
		{
			Name:        "Charisma",
			Description: "Charisma is your spellcasting ability for Cool + Beam, Aqua + Cool + Projectile, and Cool + Imbue when you cast them with this trait.",
		},
	}
}
