package api

import (
	"os" // Added os import for ConsoleWriter

	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/services"
	"github.com/rs/zerolog"     // Added zerolog import
	"github.com/rs/zerolog/log" // Added zerolog/log for global logger setup
)

// initializeHandlers creates and returns all handlers organized in a routeHandlers struct
func initializeHandlers(db database.Database, hub *Hub) *routeHandlers {
	// Initialize a logger for handlers that require it
	// In a real application, this would likely be passed down from main or a global config.
	// For now, create a basic console logger.
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	log.Logger = logger // Set global logger (optional, but good for consistency)

	// Initialize services
	resourceService := services.NewResourceService(db.ClassRepo(), db.CharacterResourceRepo())
	notorietyService := services.NewNotorietyService(db.CharacterRepo())
	s3Service, err := services.NewS3Service()
	if err != nil {
		log.Warn().Err(err).Msg("S3 service not initialized; uploads and store owner portrait URLs disabled")
		s3Service = nil
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

	lootService := services.NewLootService(db.DB(), db.ItemRepo(), db.WeaponRepo(), db.CharacterRepo())

	corpseService := services.NewCorpseService(db.CorpseRepo())
	linkService := services.NewLinkService(db.CharacterLinkRepo())

	// Use real Gemini client when GOOGLE_API_KEY is set, otherwise fall back to mock.
	var llmClient services.LLMClient
	if apiKey := os.Getenv("GOOGLE_API_KEY"); apiKey != "" {
		llmClient = services.NewGeminiLLMClient(apiKey)
	} else {
		llmClient = services.NewMockLLMClient()
	}

	// NEW: Initialize MonsterGenerationService
	monsterGenerationService := services.NewMonsterGenerationService(
		db.MonsterRepo(),
		db.MonsterGenerationEventRepo(),
		s3Service,
		llmClient,
	)

	// UPDATED: newLevelHandler now takes characterResourceRepo, beastRepo, and consumptionHistoryRepo
	levelHandlerInstance := newLevelHandler(levelUpService, db.ClassRepo(), db.CharacterResourceRepo(), db.BeastRepo(), db.ConsumptionHistoryRepo())

	// UPDATED: newCorpseHandler now takes all required services for proper function
	corpseHandlerInstance := newCorpseHandler(corpseService, db.CharacterRepo(), db.ComponentRepo())

	// NEW: Initialize MadnessHandler
	madnessHandlerInstance := newMadnessHandler(madnessService)

	// Initialize ComponentInterpreterService
	effectRepo := database.NewEffectRepo(db.DB())
	componentInterpreterService := services.NewComponentInterpreterService(db.ComponentRepo(), effectRepo)

	// Initialize SpellSynthesisService
	synthesisService := services.NewSpellSynthesisService(db.ComponentRepo())

	// NEW: Initialize SpellAIService
	spellAIService := services.NewSpellAIService(llmClient)

	// NEW: Initialize PartyHandler
	partyHandlerInstance := newPartyHandler(db.PartyRepo(), db.CharacterRepo(), db.BeastRepo())

	savedSpellService := services.NewSavedSpellService(
		db.DB(),
		db.SavedSpellRepo(),
		db.CharacterRepo(),
		db.ClassRepo(),
	)

	return &routeHandlers{
		authHandler:            newAuthHandler(db.UserRepo()),
		userHandler:            newUserHandler(db.UserRepo()),
		characterHandler:       newCharacterHandler(db.CharacterRepo(), db.RaceRepo(), db.ClassRepo(), db.CharacterResourceRepo(), db.ItemRepo(), db.WeaponRepo(), db.SpellRepo(), resourceService, notorietyService, s3Service, componentInterpreterService, db.PartyRepo(), db.ComponentRepo(), db.StoreOwnerRepo(), savedSpellService),
		spellHandler:           newSpellHandler(db.SpellRepo(), db.CharacterRepo(), db.ClassRepo(), db.RaceRepo(), db.UserRepo(), synthesisService, componentInterpreterService, spellAIService),
		beastHandler:           newBeastHandler(db.BeastRepo(), db.AttackRepo()),
		levelHandler:           levelHandlerInstance, // Use the instance
		weaponHandler:          newWeaponHandler(db.WeaponRepo()),
		itemHandler:            newItemHandler(db.ItemRepo()),
		componentHandler:       newComponentHandler(db.ComponentRepo()),
		effectHandler:          newEffectHandler(db.EffectRepo()),
		characterEffectHandler: newCharacterEffectHandler(effectService),
		resourceHandler:        newResourceHandler(db.CharacterResourceRepo(), db.CharacterRepo(), resourceService),
		minionHandler:          newMinionHandler(minionService),
		noteHandler:            newNoteHandler(db.NoteRepo(), db.CharacterRepo(), s3Service),
		gameMapHandler:         newGameMapHandler(db.GameMapRepo(), hub, s3Service),
		mapTokenHandler:        newMapTokenHandler(db.MapTokenRepo(), db.GameMapRepo(), hub),
		mapElementHandler:      newMapElementHandler(db.MapElementRepo(), db.GameMapRepo()),
		mechanicsHandler:       NewMechanicsHandler(db.DB()),
		corpseHandler:          corpseHandlerInstance, // Use the instance
		linkHandler:            newLinkHandler(linkService),
		harvestHandler:         newHarvestHandler(harvestingService, db.CharacterRepo()),                                       // NEW: Harvest Handler
		madnessHandler:         madnessHandlerInstance,                                                                         // NEW: Madness Handler
		lootHandler:            newLootHandler(lootService, logger),                                                            // Pass logger
		monsterHandler:         newMonsterHandler(db.MonsterRepo(), db.MonsterGenerationEventRepo(), monsterGenerationService), // NEW: Monster Handler
		partyHandler:           partyHandlerInstance,                                                                           // NEW: Party Handler
		abilityHandler:         newAbilityHandler(db.CharacterRepo(), db.CharacterResourceRepo(), resourceService),
		storeOwnerHandler:      newStoreOwnerHandler(db.StoreOwnerRepo(), s3Service),
	}
}
