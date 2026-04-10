package faradhaven_races

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/seedmedia"
)

// Rimekin returns the Rimekin race seed
func Rimekin() FaradhavenRaceSeed {
	return FaradhavenRaceSeed{
		Name:         "Rimekin",
		PhotoURL:     seedmedia.URL("rimekin.jpg"),
		Description:  "Rimekin hail from both Lorwyn and Shadowmoor, though the first rimekin arose from flamekin during the Phyrexian invasion. These flamekin approached their problems with cold logic and rejected reactionary responses. As a result, the magical flames that engulfed their bodies took on a frigid air, and they became rimekin.\n\nLike flamekin, rimekin possess innate magic, but the flames they conjure burn icy blue rather than red hot. Further, these \"flames\" emanate a chilling cold rather than blazing heat. This effect extends, superficially, to the items rimekin touch. Any armor worn or weapons wielded by rimekin become coated in layers of spiky yet harmless hoarfrost.",
		CreatureType: "Humanoid",
		Size:         "Medium (about 4–7 feet tall) or Small (about 2–4 feet tall), chosen when you select this species",
		BaseSpeed:    30,
		Languages:    []string{"Common", "Primordial"},
		Traits:       rimekinTraits(),
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
	}
}
