package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// FaradhavenClassSeed defines the full class data for seeding
type FaradhavenClassSeed struct {
	Name           string
	Description    string // Short 2-sentence description of how this class uniquely plays
	HitDie         int
	PrimaryAbility string
	PhotoURL       string // URL to class artwork/photo
	Archetype      string
	Concept        string
	ClassFeatures  []string
	DnDSkillFocus  []string
	Proficiencies  string
	SkillChoice    []string
	Tools          []string
	SavingThrows   []string
	StartingEquip  []string

	// LevelFeatures maps level (1-20) to features gained at that level.
	// Used for D&D-style level progression. Level 1 base info is always built from the seed;
	// LevelFeatures[1] appends to it; LevelFeatures[2-20] set Features for those levels.
	LevelFeatures map[int]string
}

// ComponentSeed defines component data for seeding (used in faradhaven_components.go)
type ComponentSeed struct {
	Name        string
	Type        models.ComponentType
	Category    models.ComponentCategory
	Description string
	Element     string
}
