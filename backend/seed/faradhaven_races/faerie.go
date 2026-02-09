package faradhaven_races

// Faerie returns the Faerie race seed
func Faerie() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:           "Faerie",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/faerie.jpg",
		Description:    "Known for their mischief, faeries resemble insects with humanoid features. Their size and shape may vary, but all have antennae, black eyes, chitinous skin, and insectoid legs and wings. Every faerie is born from a flower and possesses innate magic, which many use to play pranks. Some Lorwyn faeries serve Queen Oura, who has proclaimed herself the faeries' ruler. Not all faeries recognize Oura's authority, however. In Shadowmoor, faeries might instead worship Queen Maralen, the elf who overthrew the previous faerie queen and, some say, ushered in the clashing of Lorwyn and Shadowmoor.",
		CreatureType:   "Fey",
		Size:           "Small (about 2–4 feet tall)",
		BaseSpeed:      30,
		Languages:      []string{"Common", "Sylvan"},
		Traits:         faerieTraits(),
		ComponentNames: []string{"Vita", "Self", "Lux", "Zone", "Expand", "Shrink"},
	}
}

func faerieTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Fairy Magic",
			Description: "You know the Vita + Self component combination (minor nature magic). Starting at 3rd level, you can cast Lux + Zone (revealing light that outlines creatures). Starting at 5th level, you can also cast Expand or Shrink (altering the size of a creature or object). Once you cast Lux + Zone or Expand/Shrink with this trait, you can't cast that combination with it again until you finish a long rest. You can also cast either of those combinations using any spell slots you have of the appropriate level. Intelligence, Wisdom, or Charisma is your spellcasting ability for these spells when you cast them with this trait (choose when you select this race).",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     faerieSpellcastingAbilityOptions(),
		},
		{
			Name:        "Flight",
			Description: "Because of your wings, you have a flying speed equal to your walking speed. You can't use this flying speed if you're wearing medium or heavy armor.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Faerie Origin",
			Description: "Faeries hail from different realms. Choose one of the following options.",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     faerieOriginOptions(),
		},
	}
}

func faerieSpellcastingAbilityOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{Name: "Intelligence", Description: "Intelligence is your spellcasting ability for the spells you cast with the Fairy Magic trait."},
		{Name: "Wisdom", Description: "Wisdom is your spellcasting ability for the spells you cast with the Fairy Magic trait."},
		{Name: "Charisma", Description: "Charisma is your spellcasting ability for the spells you cast with the Fairy Magic trait."},
	}
}

func faerieOriginOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:        "Lorwyn Faerie",
			Description: "You have no additional benefit from this origin.",
		},
		{
			Name:        "Shadowmoor Faerie",
			Description: "You have Darkvision with a range of 120 feet.",
		},
	}
}
