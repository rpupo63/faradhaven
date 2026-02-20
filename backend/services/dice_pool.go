package services

import "fmt"

// DicePool represents a set of dice to roll.
type DicePool struct {
	Count int
	Faces int
}

// String returns dice notation like "2d8".
func (d DicePool) String() string {
	return fmt.Sprintf("%dd%d", d.Count, d.Faces)
}

// Average returns the expected average of the dice pool.
func (d DicePool) Average() float64 {
	return float64(d.Count) * (float64(d.Faces) + 1) / 2
}

// CalculateSpellEffect computes the dice pool for a spell based on tier, AoE, and magnitudo modifiers.
// baseTier: the highest component tier in the spell.
// isAoE: true if the Forma is an area shape (Nova, Cone, Zone, Wall, Aura).
// magnitudes: slice of Magnitudo component names present (e.g. ["Strong"], ["Weak", "Extreme"]).
func CalculateSpellEffect(baseTier int, isAoE bool, magnitudes []string) DicePool {
	// Tier → Dice Count
	var count int
	switch {
	case baseTier <= 1:
		count = 2
	case baseTier == 2:
		count = 3
	case baseTier == 3:
		count = 5
	case baseTier == 4:
		count = 7
	default: // tier 5+
		count = 10
	}

	// AoE Tax: AoE shapes use d6, single-target uses d8
	faces := 8
	if isAoE {
		faces = 6
	}

	// Magnitudo dice stepping
	hasStrong := false
	hasWeak := false
	hasExtreme := false
	for _, m := range magnitudes {
		switch m {
		case "Strong":
			hasStrong = true
		case "Weak":
			hasWeak = true
		case "Extreme":
			hasExtreme = true
		}
	}

	if hasStrong {
		faces += 2
	}
	if hasWeak {
		faces -= 2
	}

	// Clamp to d4–d12
	if faces < 4 {
		faces = 4
	}
	if faces > 12 {
		faces = 12
	}

	// Extreme: maximize to d12, add +2 dice
	if hasExtreme {
		faces = 12
		count += 2
	}

	return DicePool{Count: count, Faces: faces}
}

// isAoEForma returns true if the Forma component creates an area of effect.
func isAoEForma(formaName string) bool {
	switch formaName {
	case "Nova", "Cone", "Zone", "Wall", "Aura":
		return true
	default:
		return false
	}
}
