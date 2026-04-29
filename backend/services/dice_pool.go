package services

import (
	"fmt"

	"github.com/rpupo63/faradhaven/backend/models"
)

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

// CalculateSpellEffect computes the dice pool for a spell based on level, damage type, shape (Forma), and magnitudo modifiers.
// level: the total number of components in the spell.
// damageType: the primary damage type of the spell.
// formaName: the name of the Forma component (shape).
// scopusName: the name of the Scopus component (targeting).
// magnitudes: slice of Magnitudo component names present.
func CalculateSpellEffect(level int, damageType models.DamageType, formaName string, scopusName string, magnitudes []string) DicePool {
	// Level (number of components) → Dice Count
	var count int
	switch {
	case level <= 3:
		count = 2
	case level == 4:
		count = 3
	case level == 5:
		count = 4
	case level == 6:
		count = 6
	case level == 7:
		count = 8
	default: // level 8+
		count = 10
	}

	// Base faces based on Shape (Forma) and Targeting (Scopus)
	// touch (Self Scopus) is most powerful.
	// Projectile/Beam is standard.
	// AoE is less powerful.
	faces := 8

	// Scopus modifier: targeting constraints can trade flexibility for potency.
	switch scopusName {
	case "Self":
		faces += 2
	case "Enemy", "LOS-Only", "Object":
		faces += 1
	case "Marked":
		faces += 2
	case "Chain", "Area-First", "Through-Walls":
		faces -= 2
	}

	// Forma modifier:
	switch formaName {
	case "Projectile", "Beam":
		// Standard (d8)
	case "Touch", "Lance":
		// Highly constrained precision/melee delivery (d10)
		faces += 2
	case "Arc":
		// Curved precision with minor flexibility cost offset (d9)
		faces += 1
	case "Nova", "Cone":
		// Restricted AoE (d6)
		faces -= 2
	case "Ring", "Pillar":
		// Mid-sized positional AoE (d6)
		faces -= 2
	case "Wall", "Zone", "Aura":
		// Large/Persistent AoE (d4, -1 die)
		faces -= 4
		count -= 1
	case "Orbit":
		// Mobile persistent field has similar pressure to other sustained AoE.
		faces -= 4
		count -= 1
	}

	// Damage Type scaling: more resistant types get larger dice to compensate.
	switch damageType {
	case models.DamageForce, models.DamageRadiant, models.DamagePsychic, models.DamageThunder:
		// Rare resistance -> Base
	case models.DamageAcid, models.DamageCold, models.DamageLightning, models.DamageNecrotic:
		// Moderate resistance -> Step up (+2 faces)
		faces += 2
	case models.DamageFire, models.DamagePoison, models.DamageBludgeoning, models.DamagePiercing, models.DamageSlashing:
		// Common resistance -> High step (+4 faces)
		faces += 4
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

	// Ensure at least 1 die
	if count < 1 {
		count = 1
	}

	return DicePool{Count: count, Faces: faces}
}

// isAoEForma returns true if the Forma component creates an area of effect.
func isAoEForma(formaName string) bool {
	switch formaName {
	case "Nova", "Cone", "Zone", "Wall", "Aura", "Ring", "Pillar", "Orbit":
		return true
	default:
		return false
	}
}
