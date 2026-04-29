package faradhaven_effects

import (
	"fmt"

	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_classes"
)

// EffectSeed defines the data structure for seeding effects
type EffectSeed struct {
	Name        string
	Description string
	Category    string
	Mechanics   string
}

// AllEffects returns all effect seeds
func AllEffects() []EffectSeed {
	effects := []EffectSeed{
		// Custom Campaign Effects - General
		{
			Name:        "Madness",
			Description: "A state of mental fracture caused by eldritch contact or traumatic magic.",
			Category:    "Condition",
			Mechanics:   "Varies by level (Short-term, Long-term, Indefinite). See Madness table.",
		},
		{
			Name:        "Feral",
			Description: "A state of primal rage where the subject loses the ability to distinguish friend from foe.",
			Category:    "Condition",
			Mechanics:   "Cannot cast spells. Must attack nearest creature. Resistance to Physical damage.",
		},

		// Class Specific Conditions
		{
			Name:        "Sanguine Fatigue",
			Description: "Exhaustion caused by overspending Blood Ichor.",
			Category:    "Class Condition",
			Mechanics:   "Disadvantage on the next attack roll.",
		},

		// Standard D&D 5e Conditions
		{
			Name:        "Exhaustion",
			Description: "A cumulative state of fatigue measuring physical and mental depletion.",
			Category:    "Condition",
			Mechanics:   "1: Disadvantage on checks. 2: Speed halved. 3: Disadvantage on attacks/saves. 4: HP max halved. 5: Speed 0. 6: Death.",
		},
		{
			Name:        "Blinded",
			Description: "The creature cannot see.",
			Category:    "Condition",
			Mechanics:   "Auto-fail sight checks. Attackers have Advantage against you. Your attacks have Disadvantage.",
		},
		{
			Name:        "Charmed",
			Description: "The creature is charmed by another.",
			Category:    "Condition",
			Mechanics:   "Cannot attack charmer. Charmer has Advantage on social checks against you.",
		},
		{
			Name:        "Deafened",
			Description: "The creature cannot hear.",
			Category:    "Condition",
			Mechanics:   "Auto-fail hearing checks.",
		},
		{
			Name:        "Frightened",
			Description: "The creature is afraid of a source of fear.",
			Category:    "Condition",
			Mechanics:   "Disadvantage on checks/attacks while source is visible. Cannot move closer to source.",
		},
		{
			Name:        "Grappled",
			Description: "The creature is held by another.",
			Category:    "Condition",
			Mechanics:   "Speed 0. Ends if grappler is incapacitated or moved away.",
		},
		{
			Name:        "Incapacitated",
			Description: "The creature cannot take actions or reactions.",
			Category:    "Condition",
			Mechanics:   "No actions or reactions.",
		},
		{
			Name:        "Invisible",
			Description: "The creature is impossible to see without magical aid.",
			Category:    "Condition",
			Mechanics:   "Heavily obscured. Attackers have Disadvantage. Your attacks have Advantage.",
		},
		{
			Name:        "Paralyzed",
			Description: "The creature is incapacitated and can't move or speak.",
			Category:    "Condition",
			Mechanics:   "Auto-fail Str/Dex saves. Attackers have Advantage. Melee hits within 5ft are auto-crits.",
		},
		{
			Name:        "Petrified",
			Description: "The creature is transformed into a solid inanimate substance.",
			Category:    "Condition",
			Mechanics:   "Incapacitated, unaware, resistance to all damage, immune to poison/disease.",
		},
		{
			Name:        "Poisoned",
			Description: "The creature is affected by a toxin.",
			Category:    "Condition",
			Mechanics:   "Disadvantage on attack rolls and ability checks.",
		},
		{
			Name:        "Prone",
			Description: "The creature is lying on the ground.",
			Category:    "Condition",
			Mechanics:   "Crawling only. Melee attacks against you have Advantage (within 5ft). Ranged attacks have Disadvantage.",
		},
		{
			Name:        "Restrained",
			Description: "The creature is held in place by bonds or magic.",
			Category:    "Condition",
			Mechanics:   "Speed 0. Attackers have Advantage. Your attacks have Disadvantage. Disadvantage on Dex saves.",
		},
		{
			Name:        "Stunned",
			Description: "The creature is incapacitated and reeling.",
			Category:    "Condition",
			Mechanics:   "Auto-fail Str/Dex saves. Attackers have Advantage.",
		},
		{
			Name:        "Unconscious",
			Description: "The creature is knocked out.",
			Category:    "Condition",
			Mechanics:   "Incapacitated, unaware, drop held items. Auto-fail Str/Dex saves. Attackers have Advantage. Melee hits within 5ft are auto-crits.",
		},
	}

	// Add Mutagen Feral Effects
	for roll, desc := range faradhaven_classes.MutagenFeralTable() {
		effects = append(effects, EffectSeed{
			Name:        fmt.Sprintf("Feral Mutation (Roll %d)", roll),
			Description: desc,
			Category:    "Mutagen Feral",
			Mechanics:   desc,
		})
	}

	// Add Lorewright Madness Effects
	for roll, desc := range faradhaven_classes.LorewrightMadnessTable() {
		effects = append(effects, EffectSeed{
			Name:        fmt.Sprintf("Steampunk Madness (Roll %d)", roll),
			Description: desc,
			Category:    "Lorewright Madness",
			Mechanics:   desc,
		})
	}

	return effects
}
