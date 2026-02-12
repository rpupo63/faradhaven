package services

// CalculateFractureDC determines if a Lorewright must make a saving throw against The Fracture
// when harvesting a creature. It returns true if a save is needed, and the DC for the check.
// Logic:
// - As the Lorewright levels up, they can safely harness more complex creature memories.
// - Safe Threshold: Creature CR <= (Level / 3). If below this, no Fracture save is required.
// - If above the threshold, a Wisdom Save is required.
// - DC = 10 + (Creature CR - Safe Threshold). The harder the creature is relative to your level, the harder the check.
// - Minimum DC is 10.
func CalculateFractureDC(level int, creatureCR float64) (bool, int) {
	// Calculate the CR threshold the character can handle safely.
	safeThreshold := float64(level) / 3.0

	if creatureCR <= safeThreshold {
		return false, 0
	}

	// Calculate DC
	excess := creatureCR - safeThreshold
	dc := 10 + int(excess+0.99) // Ceiling

	return true, dc
}
