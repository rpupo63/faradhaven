package services

import (
	"fmt"
	"strings"

	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// =============================================================================
// Legacy types (kept for backward compatibility)
// =============================================================================

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

// =============================================================================
// New SpellExecution types
// =============================================================================

// SpellExecution is the top-level output of the full interpretation engine.
type SpellExecution struct {
	Geometry      GeometryInfo       `json:"geometry"`
	Damage        []DamageLayer      `json:"damage"`
	StatusEffects []StatusEffectInfo `json:"status_effects"`
	Mechanics     MechanicsInfo      `json:"mechanics"`
	Modifiers     ModifiersInfo      `json:"modifiers"`
	Description   string             `json:"description"`
}

// GeometryInfo describes the shape, size, and origin of a spell.
type GeometryInfo struct {
	Shape    string `json:"shape"`
	Size     string `json:"size"`
	Origin   string `json:"origin"`
	MaxRange string `json:"max_range"`
}

// DamageLayer represents a single layer of elemental damage.
type DamageLayer struct {
	DamageType string `json:"damage_type"`
	BaseDice   string `json:"base_dice"`
	Source     string `json:"source"`
}

// StatusEffectInfo represents a non-elemental status effect from essentia.
type StatusEffectInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

// MechanicsInfo describes force, actions, attack requirements, and saving throws.
type MechanicsInfo struct {
	ForceDirection string   `json:"force_direction"`
	Actions        []string `json:"actions"`
	RequiresAttack bool     `json:"requires_attack"`
	SavingThrow    string   `json:"saving_throw"`
}

// ModifiersInfo describes power scaling and cost adjustments.
type ModifiersInfo struct {
	PowerScale     string `json:"power_scale"`
	PropertyShift  string `json:"property_shift"`
	CostMultiplier int    `json:"cost_multiplier"`
}

// =============================================================================
// Service
// =============================================================================

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
	return category == models.CategoryEssentia
}

// Interpret takes a slice of component names and returns a structured SpellResult.
func (s *ComponentInterpreterService) Interpret(components []string) (*SpellResult, error) {
	if len(components) == 0 {
		return nil, fmt.Errorf("no components provided to interpret")
	}

	componentDetails, err := s.componentRepo.GetComponentsByNames(components)
	if err != nil {
		return nil, fmt.Errorf("failed to get component details: %w", err)
	}

	componentMap := make(map[string]models.Component)
	for _, comp := range componentDetails {
		componentMap[comp.Name] = comp
	}

	result := &SpellResult{
		Description: "A dynamically crafted spell.",
		Range:       30,
	}
	damageMap := make(map[string]*DamageInstance)

	var shapeSet bool
	for _, compName := range components {
		comp, ok := componentMap[compName]
		if !ok {
			continue
		}

		if comp.Element != "" {
			if instance, exists := damageMap[comp.Element]; exists {
				instance.DiceCount++
			} else {
				damageMap[comp.Element] = &DamageInstance{
					DamageType: comp.Element,
					DiceCount:  1,
					DiceValue:  8,
				}
			}
		}

		if !shapeSet && comp.Category == models.CategoryForma {
			switch comp.Name {
			case "Nova":
				result.AreaOfEffect = "15ft radius sphere"
				result.Range = 30
			case "Cone":
				result.AreaOfEffect = "15ft cone"
				result.Range = 0
			case "Wall":
				result.AreaOfEffect = "30ft long, 10ft high wall"
				result.Range = 60
			case "Zone":
				result.AreaOfEffect = "10ft radius zone"
				result.Range = 60
			case "Beam":
				result.AreaOfEffect = "100ft line"
				result.Range = 100
			case "Projectile":
				result.Range = 60
			case "Aura":
				result.AreaOfEffect = "10ft radius aura"
				result.Range = 0
			}
			shapeSet = true
		}

		if isEffectCategory(comp.Category) {
			effect, err := s.effectRepo.FindByName(comp.Name)
			if err == nil && effect != nil {
				result.Effects = append(result.Effects, *effect)
			}
		}
	}

	for _, instance := range damageMap {
		result.Damage = append(result.Damage, *instance)
	}

	return result, nil
}

// =============================================================================
// New interpretation engine
// =============================================================================

// InterpretModels takes a slice of Component models and returns a full SpellExecution.
func (s *ComponentInterpreterService) InterpretModels(components []models.Component) (*SpellExecution, error) {
	if len(components) == 0 {
		return nil, fmt.Errorf("no components provided to interpret")
	}

	exec := &SpellExecution{
		Damage:        []DamageLayer{},
		StatusEffects: []StatusEffectInfo{},
		Mechanics: MechanicsInfo{
			Actions: []string{},
		},
		Modifiers: ModifiersInfo{
			CostMultiplier: 1,
		},
	}

	var formaNames []string
	var essentiaNames []string
	var actioNames []string
	var magnitudes []string
	var scopusNames []string

	for _, comp := range components {
		switch comp.Category {

		// -----------------------------------------------------------------
		// Forma → Geometry
		// -----------------------------------------------------------------
		case models.CategoryForma:
			formaNames = append(formaNames, comp.Name)
			switch comp.Name {
			case "Projectile":
				exec.Geometry = GeometryInfo{
					Shape:    "sphere",
					Size:     "point",
					Origin:   "ranged",
					MaxRange: "120 ft",
				}
			case "Beam":
				exec.Geometry = GeometryInfo{
					Shape:    "line",
					Size:     "100 ft long, 5 ft wide",
					Origin:   "self",
					MaxRange: "Self",
				}
			case "Nova":
				exec.Geometry = GeometryInfo{
					Shape:    "sphere",
					Size:     "15 ft radius",
					Origin:   "self",
					MaxRange: "Self",
				}
			case "Cone":
				exec.Geometry = GeometryInfo{
					Shape:    "cone",
					Size:     "15 ft",
					Origin:   "self",
					MaxRange: "Self",
				}
			case "Wall":
				exec.Geometry = GeometryInfo{
					Shape:    "wall",
					Size:     "30 ft long, 10 ft high, 1 ft thick",
					Origin:   "ranged",
					MaxRange: "120 ft",
				}
			case "Zone":
				exec.Geometry = GeometryInfo{
					Shape:    "sphere",
					Size:     "10 ft radius",
					Origin:   "ranged",
					MaxRange: "90 ft",
				}
			case "Aura":
				exec.Geometry = GeometryInfo{
					Shape:    "sphere",
					Size:     "10 ft radius",
					Origin:   "self",
					MaxRange: "Self",
				}
			}

		// -----------------------------------------------------------------
		// Scopus → modifies origin
		// -----------------------------------------------------------------
		case models.CategoryScopus:
			scopusNames = append(scopusNames, comp.Name)
			switch comp.Name {
			case "Target":
				exec.Geometry.Origin = "target"
			case "Self":
				exec.Geometry.Origin = "self"
			case "Ground":
				exec.Geometry.Origin = "point"
			case "Chain":
				exec.Mechanics.Actions = append(exec.Mechanics.Actions, "Chain (bounces to 3 targets)")
			}

		// -----------------------------------------------------------------
		// Essentia → DamageLayer per element OR StatusEffect
		// -----------------------------------------------------------------
		case models.CategoryEssentia:
			essentiaNames = append(essentiaNames, comp.Name)
			if comp.Element != "" {
				exec.Damage = append(exec.Damage, DamageLayer{
					DamageType: comp.Element,
					BaseDice:   "", // computed after loop via CalculateSpellEffect
					Source:     comp.Name,
				})
			} else {
				// Non-elemental essentia → status effects
				se := StatusEffectInfo{Source: comp.Name}
				switch comp.Name {
				case "Spatium":
					se.Name = "Spatial Distortion"
					se.Description = "Warps the fabric of space around the target"
				case "Chronos":
					se.Name = "Time Warp"
					se.Description = "Distorts the flow of time around the target"
				case "Odor":
					se.Name = "Sensory Assault"
					se.Description = "Overwhelms the target's senses"
				default:
					// Emotion-type essentia
					se.Name = comp.Name + " Aura"
					se.Description = "Radiates an aura of " + strings.ToLower(comp.Name)
				}
				exec.StatusEffects = append(exec.StatusEffects, se)
			}

		// -----------------------------------------------------------------
		// Actio → MechanicsInfo
		// -----------------------------------------------------------------
		case models.CategoryActio:
			actioNames = append(actioNames, comp.Name)
			exec.Mechanics.Actions = append(exec.Mechanics.Actions, comp.Name)
			switch comp.Name {
			case "Push", "Pull", "Lift", "Spin":
				exec.Mechanics.ForceDirection = strings.ToLower(comp.Name)
				exec.Mechanics.SavingThrow = "STR"
			case "Pierce", "Crush":
				exec.Mechanics.RequiresAttack = true
			case "Bind", "Grab":
				exec.Mechanics.SavingThrow = "STR"
			case "Create", "Destroy", "Mutate":
				// These are already appended to Actions above
			}

		// -----------------------------------------------------------------
		// Magnitudo → ModifiersInfo
		// -----------------------------------------------------------------
		case models.CategoryMagnitudo:
			magnitudes = append(magnitudes, comp.Name)
			switch comp.Name {
			case "Strong":
				exec.Modifiers.PowerScale = "strong"
			case "Weak":
				exec.Modifiers.PowerScale = "weak"
			case "Extreme":
				exec.Modifiers.CostMultiplier = 2
			case "Increase":
				exec.Modifiers.PropertyShift = "increase"
			case "Decrease":
				exec.Modifiers.PropertyShift = "decrease"
			}
		}
	}

	// Compute dice pool via Laws of Equivalency and apply to all damage layers
	if len(exec.Damage) > 0 {
		formaName := ""
		if len(formaNames) > 0 {
			formaName = formaNames[0]
		}
		scopusName := ""
		if len(scopusNames) > 0 {
			scopusName = scopusNames[0]
		}
		primaryType := models.DamageType(exec.Damage[0].DamageType)
		pool := CalculateSpellEffect(len(components), primaryType, formaName, scopusName, magnitudes)

		// Split dice evenly among damage layers
		perLayer := pool.Count / len(exec.Damage)
		remainder := pool.Count % len(exec.Damage)
		for i := range exec.Damage {
			count := perLayer
			if i < remainder {
				count++
			}
			exec.Damage[i].BaseDice = (DicePool{Count: count, Faces: pool.Faces}).String()
		}
	}

	// Build description
	exec.Description = buildDescription(formaNames, essentiaNames, actioNames)

	return exec, nil
}

// buildDescription creates a human-readable description from the component names.
func buildDescription(formaNames, essentiaNames, actioNames []string) string {
	forma := "unknown"
	if len(formaNames) > 0 {
		forma = formaNames[0]
	}

	essentia := "raw magic"
	if len(essentiaNames) > 0 {
		essentia = strings.Join(essentiaNames, ", ")
	}

	actio := "manifests"
	if len(actioNames) > 0 {
		actio = strings.ToLower(strings.Join(actioNames, ", "))
	}

	return fmt.Sprintf("A %s spell of %s that %s.", forma, essentia, actio)
}
