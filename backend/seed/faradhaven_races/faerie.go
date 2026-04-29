package faradhaven_races

import (
	"github.com/rpupo63/faradhaven/backend/seed/seedmedia"
)

// Faerie returns the Faerie race seed
func Faerie() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Faerie",
		PhotoURL:     seedmedia.URL("faerie.jpg"),
		Description:  "Known for their mischief, faeries resemble insects with humanoid features. Their size and shape may vary, but all have antennae, black eyes, chitinous skin, and insectoid legs and wings. Every faerie is born from a flower and possesses innate magic, which many use to play pranks. Some Lorwyn faeries serve Queen Oura, who has proclaimed herself the faeries' ruler. Not all faeries recognize Oura's authority, however. In Shadowmoor, faeries might instead worship Queen Maralen, the elf who overthrew the previous faerie queen and, some say, ushered in the clashing of Lorwyn and Shadowmoor.",
		CreatureType: "Fey",
		Size:         "Small (about 2–4 feet tall)",
		BaseSpeed:    30,
		Languages:    []string{"Common", "Sylvan"},
		Traits:       faerieTraits(),
	}
}

func faerieTraits() []TraitSeed {
	return []TraitSeed{
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
