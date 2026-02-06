package api

import (
	"github.com/go-chi/chi/v5"
)

// setupFrontendRoutes sets up all routes with authentication
func setupFrontendRoutes(r chi.Router, handlers *routeHandlers, authMiddleware authMiddleware) {
	// Public API routes (no auth)
	r.Group(func(r chi.Router) {
		r.Use(ColoredHTTPLoggingMiddleware)
		r.Post("/api/auth/login", handlers.authHandler.login())
		r.Post("/api/auth/register", handlers.authHandler.register())
		// Class compendium (reference data, no auth required)
		r.Get("/api/classes", handlers.characterHandler.getAllClasses())
		r.Get("/api/classes/{classID}", handlers.characterHandler.getClassByID())
		// Race compendium (reference data, no auth required)
		r.Get("/api/races", handlers.characterHandler.getAllRaces())
		r.Get("/api/races/{raceID}", handlers.characterHandler.getRaceByID())
		// Character creation options (combines races + classes, reference data)
		r.Get("/api/characters/options", handlers.characterHandler.getCreationOptions())
		// Weapon compendium (reference data, no auth required)
		r.Get("/api/weapons", handlers.weaponHandler.getAllWeapons())
		// Item compendium (reference data, no auth required)
		r.Get("/api/items", handlers.itemHandler.getAllItems())
		r.Get("/api/items/{itemID}", handlers.itemHandler.getItemByID())
		// Component compendium (reference data for periodic table, no auth required)
		r.Get("/api/components", handlers.componentHandler.getAllComponents())
		r.Get("/api/components/{componentID}", handlers.componentHandler.getComponentByID())
		r.Get("/api/components/category/{category}", handlers.componentHandler.getComponentsByCategory())
	})

	// Protected API routes (token required)
	r.Group(func(r chi.Router) {
		r.Use(ColoredHTTPLoggingMiddleware)
		r.Use(authMiddleware.authenticate)

		// User endpoints
		r.Get("/api/users", handlers.userHandler.getAllUsers())
		r.Get("/api/user/{userID}", handlers.userHandler.getUser())
		r.Get("/api/user/{userID}/full", handlers.userHandler.getUserFull())
		r.Post("/api/user", handlers.userHandler.createUser())
		r.Put("/api/user/{userID}", handlers.userHandler.updateUser())
		r.Put("/api/user/{userID}/active-character", handlers.userHandler.setActiveCharacter())
		r.Delete("/api/user/{userID}", handlers.userHandler.deleteUser())

		// Note: /api/races, /api/races/{raceID}, /api/classes, /api/classes/{classID}, /api/characters/options are also available publicly

		// Character endpoints
		r.Get("/api/characters", handlers.characterHandler.getAllCharacters())
		r.Get("/api/character/{characterID}", handlers.characterHandler.getCharacter())
		r.Get("/api/character/{characterID}/sheet", handlers.characterHandler.getCharacterSheet())
		r.Post("/api/character/{characterID}/rest", handlers.characterHandler.restSpellPoints())
		r.Get("/api/user/{userID}/characters", handlers.characterHandler.getCharactersByUser())
		r.Post("/api/character", handlers.characterHandler.createCharacter())
		r.Put("/api/character/{characterID}", handlers.characterHandler.updateCharacter())
		r.Delete("/api/character/{characterID}", handlers.characterHandler.deleteCharacter())
		r.Post("/api/character/{characterID}/purchase", handlers.characterHandler.purchaseItem())

		// Character backstory endpoints
		r.Get("/api/character/{characterID}/backstory", handlers.characterHandler.getBackstory())
		r.Put("/api/character/{characterID}/backstory", handlers.characterHandler.updateBackstory())

		// Level-up/Level-down endpoints
		r.Get("/api/character/{characterID}/level-up/preview", handlers.levelHandler.getLevelUpPreview())
		r.Post("/api/character/{characterID}/level-up", handlers.levelHandler.levelUp())
		r.Post("/api/character/{characterID}/level-down", handlers.levelHandler.levelDown())
		r.Get("/api/character/{characterID}/level-history", handlers.levelHandler.getLevelHistory())

		// HP and hit dice management endpoints
		r.Patch("/api/character/{characterID}/hp", handlers.levelHandler.updateHP())
		r.Put("/api/character/{characterID}/temp-hp", handlers.levelHandler.setTempHP())
		r.Post("/api/character/{characterID}/hit-dice", handlers.levelHandler.useHitDice())
		r.Post("/api/character/{characterID}/rest/short", handlers.levelHandler.shortRest())
		r.Post("/api/character/{characterID}/rest/long", handlers.levelHandler.longRest())

		// Spell endpoints
		r.Get("/api/spells", handlers.spellHandler.getAllSpells())
		r.Get("/api/spell/{spellID}", handlers.spellHandler.getSpell())
		r.Get("/api/user/{userID}/spells", handlers.spellHandler.getSpellsByUser())
		r.Get("/api/character/{characterID}/spells", handlers.spellHandler.getSpellsByCharacter())
		r.Post("/api/spell", handlers.spellHandler.createSpell())
		r.Put("/api/spell/{spellID}", handlers.spellHandler.updateSpell())
		r.Delete("/api/spell/{spellID}", handlers.spellHandler.deleteSpell())

		// Beast endpoints
		r.Get("/api/beasts", handlers.beastHandler.getAllBeasts())
		r.Get("/api/beast/{beastID}", handlers.beastHandler.getBeast())
		r.Get("/api/user/{userID}/beasts", handlers.beastHandler.getBeastsByUser())
		r.Post("/api/beast", handlers.beastHandler.createBeast())
		r.Put("/api/beast/{beastID}", handlers.beastHandler.updateBeast())
		r.Delete("/api/beast/{beastID}", handlers.beastHandler.deleteBeast())
	})
}
