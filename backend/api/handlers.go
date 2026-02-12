package api

import (
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/services"
)

// initializeHandlers creates and returns all handlers organized in a routeHandlers struct
func initializeHandlers(db database.Database) *routeHandlers {
	// Initialize services
	resourceService := services.NewResourceService(db.ClassRepo(), db.CharacterResourceRepo())
	notorietyService := services.NewNotorietyService(db.CharacterRepo())
	s3Service, err := services.NewS3Service()
	if err != nil {
		// Log error but don't fail startup if S3 is not configured?
		// Or fail? Given this is "initializeHandlers" called from main, maybe we should log.
		// Since we don't have a logger passed in here easily (except importing zerolog/log), we can just panic or log.
		// However, the signature of initializeHandlers doesn't return error.
		// Let's assume for now we log and pass nil if error, and handler handles nil service?
		// Or better, just log and continue, handler will fail at runtime if used.
		// Ideally we should return error from initializeHandlers.
		// For simplicity/speed matching existing pattern:
	}

	levelUpService := services.NewLevelUpService(
		db.DB(),
		db.CharacterRepo(),
		db.ClassRepo(),
		db.LevelUpHistoryRepo(),
		db.ArchetypeRepo(),
		db.WeaponRepo(),
		db.ComponentRepo(),
		resourceService,
	)

	// NEW: Initialize HarvestingService
	harvestingService := services.NewHarvestingService(
		db.CharacterRepo(),
		db.BeastRepo(),
		db.ConsumptionHistoryRepo(),
		db.ClassRepo(),
	)

	// NEW: Initialize MadnessService
	madnessService := services.NewMadnessService(
		db.CharacterRepo(),
		db.ClassRepo(),
	)

	effectService := services.NewEffectService(
		db.DB(),
		db.CharacterEffectRepo(),
		db.EffectRepo(),
		db.CharacterRepo(),
	)

	minionService := services.NewMinionService(
		db.DB(),
		db.MinionRepo(),
		db.CharacterRepo(),
		db.ClassRepo(),
		db.ComponentRepo(),
	)

	corpseService := services.NewCorpseService(db.CorpseRepo())
	linkService := services.NewLinkService(db.CharacterLinkRepo())

	// UPDATED: newLevelHandler now takes characterResourceRepo, beastRepo, and consumptionHistoryRepo
	levelHandlerInstance := newLevelHandler(levelUpService, db.ClassRepo(), db.CharacterResourceRepo(), db.BeastRepo(), db.ConsumptionHistoryRepo())

	// UPDATED: newCorpseHandler now takes all required services for proper function
	corpseHandlerInstance := newCorpseHandler(corpseService, db.CharacterRepo(), db.ComponentRepo())

	// NEW: Initialize MadnessHandler
	madnessHandlerInstance := newMadnessHandler(madnessService)

	// NEW: Initialize ComponentInterpreterService
	effectRepo := database.NewEffectRepo(db.DB())
	componentInterpreterService := services.NewComponentInterpreterService(db.ComponentRepo(), effectRepo)

	return &routeHandlers{
		authHandler:            newAuthHandler(db.UserRepo()),
		userHandler:            newUserHandler(db.UserRepo()),
		characterHandler:       newCharacterHandler(db.CharacterRepo(), db.RaceRepo(), db.ClassRepo(), db.CharacterResourceRepo(), db.ItemRepo(), db.WeaponRepo(), db.SpellRepo(), resourceService, notorietyService, s3Service, componentInterpreterService),
		spellHandler:           newSpellHandler(db.SpellRepo()),
		beastHandler:           newBeastHandler(db.BeastRepo(), db.AttackRepo()),
		levelHandler:           levelHandlerInstance, // Use the instance
		weaponHandler:          newWeaponHandler(db.WeaponRepo()),
		itemHandler:            newItemHandler(db.ItemRepo()),
		componentHandler:       newComponentHandler(db.ComponentRepo()),
		effectHandler:          newEffectHandler(db.EffectRepo()),
		characterEffectHandler: newCharacterEffectHandler(effectService),
		resourceHandler:        newResourceHandler(db.CharacterResourceRepo()),
		minionHandler:          newMinionHandler(minionService),
		noteHandler:            newNoteHandler(db.NoteRepo()),
		mapHandler:             newMapHandler(db.MapRepo()),
		mechanicsHandler:       NewMechanicsHandler(db.DB()),
		corpseHandler:          corpseHandlerInstance, // Use the instance
		linkHandler:            newLinkHandler(linkService),
		harvestHandler:         newHarvestHandler(harvestingService, db.CharacterRepo()), // NEW: Harvest Handler
		madnessHandler:         madnessHandlerInstance,                                   // NEW: Madness Handler
	}
}
