package faradhaven_races

// Elf returns the Elf race seed
func Elf() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:           "Elf",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/elf.jpg",
		Description:    "Created by the god Corellon, the first elves could change their forms at will. They lost this ability when Corellon cursed them for plotting with the deity Lolth, who tried and failed to usurp Corellon's dominion. When Lolth was cast into the Abyss, most elves renounced her and earned Corellon's forgiveness, but that which Corellon had taken from them was lost forever. No longer able to shape-shift at will, the elves retreated to the Feywild, where their sorrow was deepened by that plane's influence. Over time, curiosity led many of them to explore other planes of existence, including worlds in the Material Plane. Elves have pointed ears and lack facial and body hair. They live for around 750 years, and they don't sleep but instead enter a trance when they need to rest. In that state, they remain aware of their surroundings while immersing themselves in memories and meditations. An environment subtly transforms elves after they inhabit it for a millennium or more, and it grants them certain kinds of magic. Drow, high elves, and wood elves are examples of elves who have been transformed thus.",
		CreatureType:   "Humanoid",
		Size:           "Medium (about 5-6 feet tall)",
		BaseSpeed:      30,
		Traits:         elfTraits(),
		ComponentNames: []string{"Lux", "Summon", "Zone", "Umbra", "Arcanum", "Self", "Sight", "Teleport", "Vita", "Haste", "Invisible", "Pull", "Projectile", "Command", "Silence", "Mend", "Soul", "Slow", "Extreme"},
	}
}

func elfTraits() []TraitSeed {
	return []TraitSeed{
		{
			Name:        "Darkvision",
			Description: "You have Darkvision with a range of 60 feet.",
			LevelReq:    1,
			ActionType:  "Passive",
			RangeValue:  "60",
		},
		{
			Name:        "Elven Lineage",
			Description: "You are part of a lineage that grants you supernatural abilities. Choose a lineage from the Elven Lineages table. You gain the level 1 benefit of that lineage. When you reach character levels 3 and 5, you learn a higher-level spell, as shown on the table. You always have that spell prepared. You can cast it once without a spell slot, and you regain the ability to cast it in that way when you finish a Long Rest. You can also cast the spell using any spell slots you have of the appropriate level. Intelligence, Wisdom, or Charisma is your spellcasting ability for the spells you cast with this trait (choose the ability when you select the lineage).",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     elfLineageOptions(),
		},
		{
			Name:        "Fey Ancestry",
			Description: "You have Advantage on saving throws you make to avoid or end the Charmed condition.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
		{
			Name:        "Keen Senses",
			Description: "You have proficiency in the Insight, Perception, or Survival skill.",
			LevelReq:    1,
			ActionType:  "Passive",
			Options:     elfKeenSensesOptions(),
		},
		{
			Name:        "Trance",
			Description: "You don't need to sleep, and magic can't put you to sleep. You can finish a Long Rest in 4 hours if you spend those hours in a trancelike meditation, during which you retain consciousness.",
			LevelReq:    1,
			ActionType:  "Passive",
		},
	}
}

func elfLineageOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{
			Name:        "Drow",
			Description: "Level 1: The range of your Darkvision increases to 120 feet. You also know the Lux + Summon component combination (dancing motes of light). Level 3: Lux + Zone (revealing light that outlines creatures). Level 5: Umbra + Zone (an area of magical darkness).",
		},
		{
			Name:        "High Elf",
			Description: "Level 1: You know the Arcanum + Self component combination (minor magical tricks). Whenever you finish a Long Rest, you can replace it with a different cantrip-tier component combination from the Wizard spell list. Level 3: Sight (detecting magical auras). Level 5: Teleport (short-range displacement through mist).",
		},
		{
			Name:        "Wood Elf",
			Description: "Level 1: Your Speed increases to 35 feet. You also know the Vita + Self component combination (minor nature magic). Level 3: Haste + Self (enhanced movement speed). Level 5: Invisible + Zone (you and nearby allies leave no trace of passage).",
		},
		{
			Name:        "Lorwyn Elf",
			Description: "Level 1: You know the Vita + Pull + Projectile component combination (a thorny vine that pulls targets toward you). Whenever you finish a Long Rest, you can replace it with a different cantrip-tier component combination from the Druid spell list. Level 3: Command (forcing a single action). Level 5: Silence + Zone (an area where sound cannot exist).",
		},
		{
			Name:        "Shadowmoor Elf",
			Description: "Level 1: The range of your Darkvision increases to 120 feet. You also know the Lux + Projectile component combination (a mote of starlight). Level 3: Mend + Self (bolstering courage and vitality). Level 5: Soul + Slow + Extreme (preserving the dead from decay).",
		},
	}
}

func elfKeenSensesOptions() []TraitOptionSeed {
	return []TraitOptionSeed{
		{Name: "Insight", Description: "You have proficiency in the Insight skill."},
		{Name: "Perception", Description: "You have proficiency in the Perception skill."},
		{Name: "Survival", Description: "You have proficiency in the Survival skill."},
	}
}
