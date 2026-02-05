package faradhaven_classes

// SpellSystemComponents returns all universal spell system components.
// These are the building blocks of the magic system, organized into categories.
func SpellSystemComponents() []ComponentSeed {
	var components []ComponentSeed
	
	components = append(components, primordialComponents()...)
	components = append(components, physicalComponents()...)
	components = append(components, thermodynamicComponents()...)
	components = append(components, spatialComponents()...)
	components = append(components, lifeComponents()...)
	components = append(components, shapeComponents()...)
	components = append(components, schoolComponents()...)
	
	return components
}