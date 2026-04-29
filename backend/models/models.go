package models

// AllModels returns all models for database migration.
// Order matters: parent tables must be created before child tables with foreign keys.
func AllModels() []interface{} {
	return []interface{}{
		// Core entities (no foreign keys to other app tables)
		&User{},
		&Party{}, // NEW: Party model
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
		&ConsumptionHistory{},
		&Minion{},
		&SavedSpell{},
		&CharacterLink{},
		&LevelUpHistory{},

		// World objects
		&Corpse{},

		// Spell hierarchy
		&Spell{},
		&SpellComponent{},

		// Monster hierarchy
		&Monster{},
		&MonsterAttack{},
		&MonsterAction{},
		&MonsterGenerationEvent{},

		// Beast hierarchy
		&Beast{},
		&Attack{},
		&BeastSkill{},
		&PartyBeast{}, // NEW: PartyBeast join table

		// Item hierarchy
		&Weapon{},
		&WeaponDamage{},
		&Item{},
		&ItemLootTheme{},
		&ItemLootLocation{},
		&ItemLootSource{},
		&ItemLootTier{},
		&ItemLootRewardAmount{},
		&ItemLootLevelBand{},
		&WeaponLootTheme{},
		&WeaponLootLocation{},
		&WeaponLootSource{},
		&WeaponLootTier{},
		&WeaponLootRewardAmount{},
		&WeaponLootLevelBand{},
		&LootLevelBudgetProfile{},
		&Effect{},
		&CharacterEffect{},

		// Character weapon system (explicit join table with modifiers)
		&CharacterWeapon{},
		&WeaponModifier{},
		&ClassWeaponRequirement{},

		// Vendors
		&StoreOwner{},
		&StoreOwnerCatalogRule{},

		// Seed tracking
		&SeedMetadata{},

		// Community
		&SharedNote{},

		// Game Map
		&GameMap{},
		&MapToken{},
		&MapElement{},
		&TrapProperties{},
		&DifficultTerrainProperties{},
		&ElevationProperties{},
		&WallProperties{},
	}
}
