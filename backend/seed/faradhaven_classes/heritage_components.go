package faradhaven_classes

// HeritageSpeciesComponents lists spell-component names formerly granted by playable races.
// They are merged into every class ComponentPool at seed time so ancestry-linked magic is
// expressed through the class component list, not race_components.
func HeritageSpeciesComponents() []string {
	return []string{
		"Anger", "Aer", "Aqua", "Arcanum", "Beam", "Bind", "Create", "Crush",
		"Decrease", "Extreme", "Fear", "Fulgur", "Ignis", "Increase", "Lux",
		"Mortis", "Nova", "Projectile", "Pull", "Self", "Sonus", "Spatium",
		"Strong", "Terra", "Umbra", "Vita", "Zone",
	}
}
