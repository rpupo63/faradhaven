package models

// AllModels returns all models for database migration.
// Order matters: parent tables must be created before child tables with foreign keys.
func AllModels() []interface{} {
	return []interface{}{
		// Core entities (no foreign keys to other app tables)
		&User{},
		&Component{},

		// Race hierarchy
		&Race{},
		&Lineage{},
		&Trait{},
		&TraitOption{},

		// Class hierarchy
		&Class{},
		&ClassLevel{},
		&LevelFeature{},
		&ClassComponent{},

		// Character hierarchy
		&Character{},
		&CharacterSkill{},
		&LevelUpHistory{},

		// Spell hierarchy
		&Spell{},
		&SpellComponent{},

		// Beast hierarchy
		&Beast{},
		&Attack{},
	}
}
