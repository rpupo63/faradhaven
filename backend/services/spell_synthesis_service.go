package services

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// SpellSynthesis holds auto-derived suggestions from a set of components.
type SpellSynthesis struct {
	ManaCost               int     `json:"mana_cost"`
	BaseTier               int     `json:"base_tier"`
	SlotLevel              int     `json:"slot_level"`
	SuggestedType          string  `json:"suggested_type"`
	SuggestedRange         *string `json:"suggested_range,omitempty"`
	SuggestedDamageType    *string `json:"suggested_damage_type,omitempty"`
	SuggestedDamageDice    *string `json:"suggested_damage_dice,omitempty"`
	SuggestedDuration      *string `json:"suggested_duration,omitempty"`
	SuggestedConcentration bool    `json:"suggested_concentration"`
	ValidationErrors       []string `json:"validation_errors"`
}

// SpellSynthesisService computes spell properties from components.
type SpellSynthesisService struct {
	componentRepo *database.ComponentRepo
}

// NewSpellSynthesisService creates a new SpellSynthesisService.
func NewSpellSynthesisService(componentRepo *database.ComponentRepo) *SpellSynthesisService {
	return &SpellSynthesisService{componentRepo: componentRepo}
}

// FetchComponents retrieves components by their IDs.
func (s *SpellSynthesisService) FetchComponents(ids []uuid.UUID) ([]models.Component, error) {
	return s.componentRepo.GetComponentsByIDs(ids)
}

// Validate checks component MECE rules and returns any violations.
func (s *SpellSynthesisService) Validate(components []models.Component) []string {
	var errors []string

	var formaCount int
	var essentiaCount int
	var hasIncrease, hasDecrease, hasStrong, hasWeak bool

	for _, c := range components {
		switch c.Category {
		case models.CategoryForma:
			formaCount++
		case models.CategoryEssentia:
			essentiaCount++
		case models.CategoryMagnitudo:
			switch c.Name {
			case "Increase":
				hasIncrease = true
			case "Decrease":
				hasDecrease = true
			case "Strong":
				hasStrong = true
			case "Weak":
				hasWeak = true
			}
		}
	}

	if formaCount == 0 {
		errors = append(errors, "Exactly 1 Forma (shape) component is required")
	} else if formaCount > 1 {
		errors = append(errors, "Only 1 Forma (shape) component is allowed")
	}

	if essentiaCount == 0 {
		errors = append(errors, "At least 1 Essentia (element/domain) component is required")
	}

	if hasIncrease && hasDecrease {
		errors = append(errors, "Cannot use both Increase and Decrease modifiers")
	}
	if hasStrong && hasWeak {
		errors = append(errors, "Cannot use both Strong and Weak modifiers")
	}

	return errors
}

// ComputeCost calculates mana cost and base tier from components.
func (s *SpellSynthesisService) ComputeCost(components []models.Component) (manaCost int, baseTier int) {
	baseTier = 1
	hasExtreme := false

	for _, c := range components {
		manaCost += c.Tier
		if c.Tier > baseTier {
			baseTier = c.Tier
		}
		if c.Name == "Extreme" {
			hasExtreme = true
		}
	}

	if hasExtreme {
		manaCost *= 2
	}

	return manaCost, baseTier
}

// Synthesize derives spell property suggestions from a set of components.
func (s *SpellSynthesisService) Synthesize(components []models.Component) *SpellSynthesis {
	syn := &SpellSynthesis{
		SuggestedType: "Utility",
	}

	syn.ValidationErrors = s.Validate(components)
	syn.ManaCost, syn.BaseTier = s.ComputeCost(components)

	// SlotLevel derived from base tier
	syn.SlotLevel = syn.BaseTier

	var formaName string
	var magnitudes []string
	var hasDamage bool

	for _, c := range components {
		switch c.Category {
		case models.CategoryEssentia:
			if c.Element != "" {
				dmgType := mapElementToDamageType(c.Element)
				syn.SuggestedDamageType = &dmgType
				hasDamage = true
			}

		case models.CategoryForma:
			formaName = c.Name
			r := mapFormaToRange(c.Name)
			syn.SuggestedRange = &r

			// Concentration for sustained shapes
			switch c.Name {
			case "Zone", "Wall", "Aura":
				syn.SuggestedConcentration = true
				dur := "1 minute"
				syn.SuggestedDuration = &dur
			}

		case models.CategoryActio:
			syn.SuggestedType = mapActioToSpellType(c.Name)

			// Concentration for Bind
			if c.Name == "Bind" {
				syn.SuggestedConcentration = true
				dur := "1 minute"
				syn.SuggestedDuration = &dur
			}

		case models.CategoryMagnitudo:
			magnitudes = append(magnitudes, c.Name)
		}
	}

	// Compute damage dice via Laws of Equivalency
	if hasDamage && syn.SuggestedDamageType != nil {
		pool := CalculateSpellEffect(syn.BaseTier, isAoEForma(formaName), magnitudes)
		dice := pool.String()
		syn.SuggestedDamageDice = &dice
	}

	return syn
}

// mapElementToDamageType maps component Element strings to D&D damage types.
func mapElementToDamageType(element string) string {
	switch element {
	case "Fire":
		return "Fire"
	case "Cold":
		return "Cold"
	case "Lightning":
		return "Lightning"
	case "Thunder":
		return "Thunder"
	case "Acid":
		return "Acid"
	case "Poison":
		return "Poison"
	case "Necrotic":
		return "Necrotic"
	case "Radiant":
		return "Radiant"
	case "Force":
		return "Force"
	case "Psychic":
		return "Psychic"
	case "Bludgeoning":
		return "Bludgeoning"
	case "Piercing":
		return "Piercing"
	case "Slashing":
		return "Slashing"
	default:
		return element
	}
}

// mapFormaToRange maps Forma component names to range strings.
func mapFormaToRange(name string) string {
	switch name {
	case "Projectile":
		return "120 ft"
	case "Beam":
		return "Self (100 ft line)"
	case "Nova":
		return "Self (15 ft radius)"
	case "Cone":
		return "Self (15 ft cone)"
	case "Wall":
		return "120 ft (30 ft wall)"
	case "Zone":
		return "90 ft (10 ft radius)"
	case "Aura":
		return "Self (10 ft radius)"
	default:
		return "60 ft"
	}
}

// mapActioToSpellType maps Actio component names to spell types.
func mapActioToSpellType(name string) string {
	switch name {
	case "Push", "Pull", "Crush", "Pierce":
		return "Attack"
	case "Create":
		return "Utility"
	case "Bind", "Grab":
		return "Save"
	case "Mutate":
		return "Effect"
	case "Destroy":
		return "Attack"
	case "Lift", "Spin":
		return "Save"
	default:
		return "Utility"
	}
}

