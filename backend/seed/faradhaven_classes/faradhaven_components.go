package faradhaven_classes

// ComponentClassMapping pairs a component with the class that has access to it.
// Supports many-to-many: a component could appear in multiple ClassNames in the future.
type ComponentClassMapping struct {
	Component ComponentSeed
	ClassName string
}

// ClassComponentNames maps class names to the spell system component names they have access to.
// These reference components defined in SpellSystemComponents().
func ClassComponentNames() map[string][]string {
	return map[string][]string{
		// The Rift Weaver - dimensional/psychic magic specialist
		"The Rift Weaver": {"Arcanum", "Psi", "Teleport", "Sight", "Echo"},

		// The Mutagen - biological transformation specialist
		"The Mutagen": {"Vita", "Acidum", "Mend", "Transmute", "Haste"},

		// The Ironwright - lightning/metal tech specialist
		"The Ironwright": {"Fulgur", "Ferrum", "Conduct", "Summon", "Extreme"},

		// The Powder Mage - ranged elemental combat specialist
		"The Powder Mage": {"Ignis", "Cool", "Pierce", "Homing", "Projectile", "Scatter"},

		// The Piston Brawler - physical force melee specialist
		"The Piston Brawler": {"Push", "Crush", "Sonus", "Nova", "Focus"},

		// The Sanguinist - life/death magic specialist
		"The Sanguinist": {"Drain", "Mend", "Wither", "Chain", "Sanctus"},

		// The Vapor Blade - stealth/poison specialist
		"The Vapor Blade": {"Umbra", "Acidum", "Silence", "Invisible", "Pierce"},

		// The Lorewright - information/divination specialist
		"The Lorewright": {"Sight", "Identify", "Predict", "Link", "Lux"},
	}
}

// AllComponentClassMappings returns all spell components and their class associations.
// Components reference the universal spell system components by name.
// ClassComponent links are created from these mappings.
func AllComponentClassMappings() []ComponentClassMapping {
	// Build mappings by looking up components from the spell system
	spellComponents := SpellSystemComponents()
	componentMap := make(map[string]ComponentSeed)
	for _, c := range spellComponents {
		componentMap[c.Name] = c
	}

	var mappings []ComponentClassMapping
	for className, componentNames := range ClassComponentNames() {
		for _, compName := range componentNames {
			if comp, ok := componentMap[compName]; ok {
				mappings = append(mappings, ComponentClassMapping{
					Component: comp,
					ClassName: className,
				})
			}
		}
	}
	return mappings
}
