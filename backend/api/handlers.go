package api

import (
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/services"
)

// initializeHandlers creates and returns all handlers organized in a routeHandlers struct
func initializeHandlers(db database.Database) *routeHandlers {
	// Initialize services
	resourceService := services.NewResourceService()
	levelUpService := services.NewLevelUpService(
		db.DB(),
		db.CharacterRepo(),
		db.ClassRepo(),
		db.LevelUpHistoryRepo(),
		db.ArchetypeRepo(),
		db.WeaponRepo(),
		resourceService,
	)

	return &routeHandlers{
		authHandler:      newAuthHandler(db.UserRepo()),
		userHandler:      newUserHandler(db.UserRepo()),
		characterHandler: newCharacterHandler(db.CharacterRepo(), db.RaceRepo(), db.ClassRepo(), db.ItemRepo(), db.WeaponRepo(), resourceService),
		spellHandler:     newSpellHandler(db.SpellRepo()),
		beastHandler:     newBeastHandler(db.BeastRepo(), db.AttackRepo()),
		levelHandler:     newLevelHandler(levelUpService, db.ClassRepo(), resourceService),
		weaponHandler:    newWeaponHandler(db.WeaponRepo()),
		itemHandler:      newItemHandler(db.ItemRepo()),
		componentHandler: newComponentHandler(db.ComponentRepo()),
		effectHandler:    newEffectHandler(db.EffectRepo()),
	}
}
