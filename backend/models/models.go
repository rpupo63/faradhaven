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
		&Archetype{},
		&ClassLevel{},
		&ClassResourceDefinition{},
		&ClassLevelResource{},
		&LevelFeature{},
		&ClassComponent{},
		&ClassStartingEquipmentChoice{},
		&ClassStartingEquipmentOption{},

		// Character hierarchy
		&Character{},
		&CharacterSkill{},
		&CharacterComponent{},
		&CharacterResource{},
		&Minion{},
		&SavedSpell{},
		&CharacterLink{},
		&LevelUpHistory{},

		// World objects
		&Corpse{},

		// Spell hierarchy
		&Spell{},
		&SpellComponent{},

		// Beast hierarchy
		&Beast{},
		&Attack{},

		// Item hierarchy
		&Weapon{},
		&WeaponDamage{},
		&Item{},
		&Effect{},
		&CharacterEffect{},

		// Character weapon system (explicit join table with modifiers)
		&CharacterWeapon{},
		&WeaponModifier{},
		&ClassWeaponRequirement{},

		// Seed tracking
		&SeedMetadata{},

		// Community
		&SharedNote{},

		// Game Map
		&GameMap{},
		&MapToken{},
	}
}
