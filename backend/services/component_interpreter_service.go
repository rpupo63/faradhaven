package services

import (
	"fmt"

	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// DamageInstance represents a single instance of damage to be dealt.
type DamageInstance struct {
	DamageType string // e.g., "Fire", "Slashing", "Thunder"
	DiceCount  int    // e.g., 2 for 2d8
	DiceValue  int    // e.g., 8 for 2d8
	Modifier   int    // A flat bonus or penalty to the damage roll
}

// SpellResult is the deterministic outcome of a component combination.
type SpellResult struct {
	Damage             []DamageInstance
	Range              int
	AreaOfEffect       string // e.g., "5ft radius sphere", "30ft cone"
	Effects            []models.Effect
	SavingThrow        string // e.g., "dexterity", "wisdom"
	RequiresAttackRoll bool
	Description        string // A human-readable description of the spell created
}

// ComponentInterpreterService is responsible for parsing component strings into SpellResults.
type ComponentInterpreterService struct {
	componentRepo *database.ComponentRepo
	effectRepo    *database.EffectRepo
}

// NewComponentInterpreterService creates a new instance of the service.
func NewComponentInterpreterService(componentRepo *database.ComponentRepo, effectRepo *database.EffectRepo) *ComponentInterpreterService {
	return &ComponentInterpreterService{
		componentRepo: componentRepo,
		effectRepo:    effectRepo,
	}
}

// isEffectCategory checks if a component category should be treated as an effect.
func isEffectCategory(category models.ComponentCategory) bool {
	switch category {
	case models.CategoryAbjuration,
		models.CategoryConjuration,
		models.CategoryDivination,
		models.CategoryEnchantment,
		models.CategoryEvocation,
		models.CategoryIllusion,
		models.CategoryNecromancy,
		models.CategoryTransmutation,
		models.CategorySpatial,
		models.CategoryLife,
		models.CategoryThermodynamic:
		return true
	default:
		return false
	}
}

// Interpret takes a slice of component names and returns a structured SpellResult.
// This is where the core logic for combining components will be implemented.
func (s *ComponentInterpreterService) Interpret(components []string) (*SpellResult, error) {
	if len(components) == 0 {
		return nil, fmt.Errorf("no components provided to interpret")
	}

	// Fetch all component details from the database at once.
	componentDetails, err := s.componentRepo.GetComponentsByNames(components)
	if err != nil {
		return nil, fmt.Errorf("failed to get component details: %w", err)
	}

	// Create a map for quick lookups.
	componentMap := make(map[string]models.Component)
	for _, comp := range componentDetails {
		componentMap[comp.Name] = comp
	}

	// Initialize the result with defaults.
	result := &SpellResult{
		Description: "A dynamically crafted spell.",
		Range:       30, // Default range, can be modified by components.
	}
	damageMap := make(map[string]*DamageInstance)

	// Iterate through the original component list to preserve order and handle duplicates.
	var shapeSet bool
	for _, compName := range components {
		comp, ok := componentMap[compName]
		if !ok {
			// This could be a feature or a bug, depending on design. For now, we'll ignore unknown components.
			continue
		}

		// --- Rule 1: Damage Aggregation ---
		// If a component has an element, it contributes to damage.
		if comp.Element != "" {
			if instance, exists := damageMap[comp.Element]; exists {
				instance.DiceCount++
			} else {
				damageMap[comp.Element] = &DamageInstance{
					DamageType: comp.Element,
					DiceCount:  1,
					DiceValue:  8, // Base damage die for Powder Mage components
				}
			}
		}

		// --- Rule 2: Shape & Range ---
		// The first shape component found determines the spell's geometry.
		if !shapeSet && comp.Category == models.CategoryShape {
			switch comp.Name {
			case "Nova":
				result.AreaOfEffect = "15ft radius sphere"
				result.Range = 30 // Typically originates at a point within this range
			case "Cone":
				result.AreaOfEffect = "15ft cone"
				result.Range = 0 // Originates from self
			case "Wall":
				result.AreaOfEffect = "30ft long, 10ft high wall"
				result.Range = 60
			case "Zone":
				result.AreaOfEffect = "10ft radius zone"
				result.Range = 60
			case "Self":
				result.AreaOfEffect = "Self"
				result.Range = 0
			case "Touch":
				result.AreaOfEffect = "Touch"
				result.Range = 5
			case "Beam":
				result.AreaOfEffect = "100ft line"
				result.Range = 100
			case "Projectile":
				// This is the default, no change needed unless other components modify it.
				result.Range = 60
			case "Trap":
				result.AreaOfEffect = "5ft radius trap"
				result.Range = 5 // Placed within touch range
			}
			shapeSet = true // Ensure only one shape component is processed.
		}

		// --- Rule 3: Effects ---
		// If a component is from an effect category, find and add the corresponding effect.
		if isEffectCategory(comp.Category) {
			effect, err := s.effectRepo.FindByName(comp.Name)
			if err == nil && effect != nil {
				result.Effects = append(result.Effects, *effect)
			}
			// Note: We are ignoring errors here for now. In a real scenario,
			// you might want to log if an effect for a component is not found.
		}
	}

	// Convert the damage map into the result slice.
	for _, instance := range damageMap {
		result.Damage = append(result.Damage, *instance)
	}

	return result, nil
}
