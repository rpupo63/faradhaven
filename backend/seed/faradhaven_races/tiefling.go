package faradhaven_races

// Tiefling returns the Tiefling race seed
func Tiefling() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:           "Tiefling",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/tiefling.jpg",
		Description:    "Tieflings are either born in the Lower Planes or have fiendish ancestors who originated there. A tiefling (pronounced TEE-fling) is linked by blood to a devil, a demon, or some other Fiend. This connection to the Lower Planes is the tiefling's fiendish legacy, which comes with the promise of power yet has no effect on the tiefling's moral outlook. A tiefling chooses whether to embrace or lament their fiendish legacy. Abyssal: The entropy of the Abyss, the chaos of Pandemonium, and the despair of Carceri call to tieflings who have the abyssal legacy. Horns, fur, tusks, and peculiar scents are common physical features of such tieflings, most of whom have the blood of demons coursing through their veins. Chthonic: Tieflings who have the chthonic legacy feel not only the tug of Carceri but also the greed of Gehenna and the gloom of Hades. Some of these tieflings look cadaverous. Others possess the unearthly beauty of a succubus, or they have physical features in common with a night hag, a yugoloth, or some other Neutral Evil fiendish ancestor. Infernal: The infernal legacy connects tieflings not only to Gehenna but also the Nine Hells and the raging battlefields of Acheron. Horns, spines, tails, golden eyes, and a faint odor of sulfur or smoke are common physical features of such tieflings, most of whom trace their ancestry to devils.",
		CreatureType:   "Humanoid",
		Size:           "Medium (about 4–7 feet tall) or Small (about 3–4 feet tall), chosen when you select this species",
		BaseSpeed:      30,
		Traits:         tieflingTraits(),
		ComponentNames: []string{"Arcanum", "Self", "Wither", "Projectile", "Beam", "Slow", "Extreme", "Soul", "Mend", "Curse", "Ignis", "Reflect", "Umbra", "Zone"},
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
			Description: "You are the recipient of a legacy that grants you supernatural abilities. Choose a legacy from the Fiendish Legacies table. You gain the level 1 benefit of the chosen legacy. When you reach character levels 3 and 5, you learn a higher-level spell, as shown on the table. You always have that spell prepared. You can cast it once without a spell slot, and you regain the ability to cast it in that way when you finish a Long Rest. You can also cast the spell using any spell slots you have of the appropriate level. Intelligence, Wisdom, or Charisma is your spellcasting ability for the spells you cast with this trait (choose the ability when you select the legacy).",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     tieflingFiendishLegacyOptions(),
		},
		{
			Name:        "Otherworldly Presence",
			Description: "You innately know the Arcanum + Self component combination, allowing you to manifest minor supernatural effects such as tremors, amplified voice, or flickering flames. When you cast it with this trait, the spell uses the same spellcasting ability you use for your Fiendish Legacy trait.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
	}
}

func tieflingFiendishLegacyOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:        "Abyssal",
			Description: "You have Resistance to Poison damage. You also know the Wither + Projectile component combination (a spray of decay). When you reach character level 3, you learn Wither + Beam (a ray of sickening energy). When you reach character level 5, you learn Slow + Extreme (complete stasis). You always have these component combinations prepared. You can cast each once without a spell slot, and you regain the ability to cast them in that way when you finish a Long Rest. You can also cast them using any spell slots you have of the appropriate level.",
		},
		{
			Name:        "Chthonic",
			Description: "You have Resistance to Necrotic damage. You also know the Soul + Projectile component combination (spectral cold touch). When you reach character level 3, you learn Mend + Self (bolstering your own vitality with false life). When you reach character level 5, you learn Curse + Beam (a ray that weakens your foes). You always have these component combinations prepared. You can cast each once without a spell slot, and you regain the ability to cast them in that way when you finish a Long Rest. You can also cast them using any spell slots you have of the appropriate level.",
		},
		{
			Name:        "Infernal",
			Description: "You have Resistance to Fire damage. You also know the Ignis + Projectile component combination (a bolt of fire). When you reach character level 3, you learn Ignis + Reflect (fiery retribution when struck). When you reach character level 5, you learn Umbra + Zone (an area of magical darkness). You always have these component combinations prepared. You can cast each once without a spell slot, and you regain the ability to cast them in that way when you finish a Long Rest. You can also cast them using any spell slots you have of the appropriate level.",
		},
	}
}
