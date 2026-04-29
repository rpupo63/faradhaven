package api

import (
	"os"

	"github.com/go-chi/chi/v5"
)

// setupFrontendRoutes sets up all routes with authentication
func setupFrontendRoutes(r chi.Router, handlers *routeHandlers, authMiddleware authMiddleware, hub *Hub) {
	// Public API routes (no auth)
	r.Group(func(r chi.Router) {
		r.Use(ColoredHTTPLoggingMiddleware)
		r.Post("/api/auth/login", handlers.authHandler.login())
		r.Post("/api/auth/register", handlers.authHandler.register())
		r.Post("/api/auth/refresh", handlers.authHandler.refresh())
		r.Post("/api/auth/logout", handlers.authHandler.logout())
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
		r.Get("/api/loot/options", handlers.lootHandler.GetLootOptions)
		r.Get("/api/store-owners", handlers.storeOwnerHandler.getStoreOwners())
		r.Get("/api/store-owners/{storeOwnerID}/portrait", handlers.storeOwnerHandler.getStoreOwnerPortrait())
		// Component compendium (reference data for periodic table, no auth required)
		r.Get("/api/components", handlers.componentHandler.getAllComponents())
		r.Get("/api/components/{componentID}", handlers.componentHandler.getComponentByID())
		r.Get("/api/components/category/{category}", handlers.componentHandler.getComponentsByCategory())
		// Effect compendium (reference data, no auth required)
		r.Get("/api/effects", handlers.effectHandler.getAllEffects())
		r.Get("/api/effects/{effectID}", handlers.effectHandler.getEffectByID())
		// Spell compendium (reference data, no auth required)
		r.Get("/api/spells", handlers.spellHandler.getAllSpells())
		r.Get("/api/spell/{spellID}", handlers.spellHandler.getSpell())
		r.Get("/api/spell/{spellID}/execution", handlers.spellHandler.getSpellExecution())
		// Beast compendium (reference data, no auth required)
		r.Get("/api/beasts", handlers.beastHandler.getAllBeasts())
		r.Get("/api/beast/{beastID}", handlers.beastHandler.getBeast())
		// Maps (Public read access)
		r.Get("/api/map/{mapID}", handlers.gameMapHandler.getMap())
		r.Get("/api/map/room/{roomCode}", handlers.gameMapHandler.getMapByRoom())

		// WebSocket endpoint
		r.Get("/api/map/{mapID}/ws", ServeWs(hub, authMiddleware))
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
		r.Post("/api/character/{characterID}/cast", handlers.characterHandler.castSpell())
		r.Get("/api/character/{characterID}/speed-dial", handlers.characterHandler.getSpeedDial())
		r.Put("/api/character/{characterID}/speed-dial/{slotIndex}", handlers.characterHandler.saveSpeedDialSlot())
		r.Delete("/api/character/{characterID}/speed-dial/{slotIndex}", handlers.characterHandler.clearSpeedDialSlot())
		r.Post("/api/character/{characterID}/component/{componentID}/consume", handlers.characterHandler.consumeComponent())
		r.Post("/api/character/{characterID}/component/{componentID}/gain", handlers.characterHandler.gainComponent())
		r.Get("/api/user/{userID}/characters", handlers.characterHandler.getCharactersByUser())
		r.Post("/api/character", handlers.characterHandler.createCharacter())
		r.Put("/api/character/{characterID}", handlers.characterHandler.updateCharacter())
		r.Post("/api/character/{characterID}/notoriety", handlers.characterHandler.updateNotoriety())
		r.Post("/api/character/{characterID}/sanguine-notoriety", handlers.characterHandler.updateSanguineNotorietyPoints())
		r.Post("/api/character/{characterID}/equipment", handlers.characterHandler.updateEquipment())
		r.Delete("/api/character/{characterID}", handlers.characterHandler.deleteCharacter())
		r.Post("/api/character/{characterID}/purchase", handlers.characterHandler.purchaseItem())
		r.Post("/api/character/{characterID}/sell", handlers.characterHandler.sellItem())
		r.Post("/api/character/{characterID}/toss", handlers.characterHandler.tossItem())
		r.Post("/api/character/{characterID}/extract", handlers.characterHandler.extractComponents())
		r.Post("/api/character/{characterID}/forage-components", handlers.characterHandler.forageComponents())
		r.Post("/api/character/{characterID}/loot", handlers.lootHandler.GenerateLoot)
		r.Post("/api/character/{characterID}/loot/confirm", handlers.lootHandler.ConfirmLootPickup)
		r.Post("/api/character/{characterID}/image", handlers.characterHandler.uploadProfilePicture())
		r.Post("/api/characters/{characterID}/madness/roll", handlers.madnessHandler.rollMadness())

		// Character Link Endpoints
		r.Route("/api/character/{characterID}/links", func(r chi.Router) {
			r.Get("/", handlers.linkHandler.getLinks())
			r.Post("/", handlers.linkHandler.createLink())
			r.Delete("/{linkID}", handlers.linkHandler.removeLink())
		})

		// Character backstory endpoints
		r.Get("/api/character/{characterID}/backstory", handlers.characterHandler.getBackstory())
		r.Put("/api/character/{characterID}/backstory", handlers.characterHandler.updateBackstory())

		// Mechanics Routes
		r.Route("/api/characters/{id}/mechanics", func(r chi.Router) {
			r.Post("/roll-table", handlers.mechanicsHandler.RollTable)
			r.Post("/mutagen-cast", handlers.mechanicsHandler.MutagenCast)
			r.Get("/effects", handlers.mechanicsHandler.GetActiveEffects)
			r.Delete("/effects/{effectID}", handlers.mechanicsHandler.RemoveEffect)
		})

		// Character Effect Routes (enhanced effect management)
		r.Route("/api/characters/{id}/effects", func(r chi.Router) {
			r.Get("/", handlers.characterEffectHandler.GetActiveEffects())
			r.Post("/", handlers.characterEffectHandler.ApplyEffect())
			r.Put("/{effectInstanceID}/stack", handlers.characterEffectHandler.ModifyStacks())
			r.Delete("/{effectInstanceID}", handlers.characterEffectHandler.RemoveEffect())
			r.Post("/tick", handlers.characterEffectHandler.TickDuration())
			r.Post("/break-concentration", handlers.characterEffectHandler.BreakConcentration())
		})

		// Character Resource Routes
		r.Route("/api/characters/{id}/resources", func(r chi.Router) {
			r.Get("/", handlers.resourceHandler.GetResources())
			r.Post("/", handlers.resourceHandler.CreateResource())
			r.Get("/{key}", handlers.resourceHandler.GetResource())
			r.Post("/{key}/spend", handlers.resourceHandler.SpendResource())
			r.Post("/{key}/gain", handlers.resourceHandler.GainResource())
			r.Delete("/{key}", handlers.resourceHandler.DeleteResource())
		})

		// Minion Routes (constructs, echoes, etc.)
		r.Route("/api/characters/{id}/minions", func(r chi.Router) {
			r.Get("/", handlers.minionHandler.GetMinions())
			r.Post("/", handlers.minionHandler.CreateMinion())
			r.Get("/templates", handlers.minionHandler.GetConstructTemplates())
			r.Get("/{minionID}", handlers.minionHandler.GetMinion())
			r.Patch("/{minionID}/hp", handlers.minionHandler.UpdateMinionHP())
			r.Post("/{minionID}/activate", handlers.minionHandler.ActivateMinion())
			r.Post("/{minionID}/deactivate", handlers.minionHandler.DeactivateMinion())
			r.Delete("/{minionID}", handlers.minionHandler.DeleteMinion())
		})

		// Level Routes
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

		// Spell endpoints (protected)
		r.Get("/api/user/{userID}/spells", handlers.spellHandler.getSpellsByUser())
		r.Get("/api/character/{characterID}/spells", handlers.spellHandler.getSpellsByCharacter())
		r.Get("/api/character/{characterID}/spellbook", handlers.spellHandler.getCharacterSpellbook())
		r.Post("/api/spell", handlers.spellHandler.createSpell())
		r.Post("/api/spell/preview-ai-opinion", handlers.spellHandler.previewSpellAIOpinion())
		r.Post("/api/spell/synthesize", handlers.spellHandler.synthesizeSpell())
		r.Get("/api/spell/{spellID}/opinion", handlers.spellHandler.getSpellOpinion())
		r.Post("/api/spell/{spellID}/ai/retry", handlers.spellHandler.retryAIField())
		r.Put("/api/spell/{spellID}", handlers.spellHandler.updateSpell())
		r.Delete("/api/spell/{spellID}", handlers.spellHandler.deleteSpell())

		// GM-only endpoints
		r.Get("/api/gm/spells/unchecked", handlers.spellHandler.getUncheckedSpells())
		r.Get("/api/gm/spells/checked", handlers.spellHandler.getCheckedSpells())

		// Beast endpoints (protected)
		r.Get("/api/user/{userID}/beasts", handlers.beastHandler.getBeastsByUser())
		r.Post("/api/beast", handlers.beastHandler.createBeast())
		r.Put("/api/beast/{beastID}", handlers.beastHandler.updateBeast())
		r.Delete("/api/beast/{beastID}", handlers.beastHandler.deleteBeast())

		// Shared Notes (protected)
		r.Get("/api/notes", handlers.noteHandler.getAllNotes())
		r.Post("/api/notes", handlers.noteHandler.createNote())

		// Map endpoints (protected)
		r.Get("/api/user/{userID}/maps", handlers.gameMapHandler.getUserMaps())
		r.Post("/api/map", handlers.gameMapHandler.createMap())
		r.Post("/api/map/{mapID}/background", handlers.gameMapHandler.uploadBackgroundImage())
		r.Put("/api/map/{mapID}", handlers.gameMapHandler.updateMap())
		r.Delete("/api/map/{mapID}", handlers.gameMapHandler.deleteMap())

		// Map Token endpoints (protected)
		r.Post("/api/map/{mapID}/token", handlers.mapTokenHandler.addToken())
		r.Put("/api/map/{mapID}/token/{tokenID}", handlers.mapTokenHandler.updateToken())
		r.Delete("/api/map/{mapID}/token/{tokenID}", handlers.mapTokenHandler.deleteToken())

		// Map Element endpoints (DM only)
		r.Route("/api/map/{mapID}/elements", func(r chi.Router) {
			r.Post("/", handlers.mapElementHandler.createMapElement())
			r.Post("/batch", handlers.mapElementHandler.createMultipleMapElements())
			r.Put("/{elementID}", handlers.mapElementHandler.updateMapElement())
			r.Delete("/{elementID}", handlers.mapElementHandler.deleteMapElement())
		})

		// Initiative endpoints
		r.Get("/api/map/{mapID}/initiative", handlers.mapTokenHandler.getInitiative())
		r.Put("/api/map/{mapID}/initiative", handlers.mapTokenHandler.setInitiative())
		r.Delete("/api/map/{mapID}/initiative", handlers.mapTokenHandler.clearInitiative())

		// Corpse endpoints
		r.Route("/api/corpses", func(r chi.Router) {
			r.Get("/", handlers.corpseHandler.GetCorpses())
			r.Post("/", handlers.corpseHandler.CreateCorpse())
			r.Get("/{corpseID}", handlers.corpseHandler.GetCorpse())
			r.Delete("/{corpseID}", handlers.corpseHandler.DeleteCorpse())
			r.Post("/{corpseID}/harvest", handlers.corpseHandler.HarvestCorpse())
			r.Post("/{corpseID}/consume", handlers.corpseHandler.ConsumeCorpse())
			r.Post("/{corpseID}/scavenge-components", handlers.corpseHandler.ScavengeComponents())
		})

		// Harvesting endpoints (for Lorewright)
		r.Get("/api/beasts/{beastID}/harvestable-abilities", handlers.harvestHandler.getHarvestableAbilities())
		r.Post("/api/characters/{characterID}/harvest", handlers.harvestHandler.confirmHarvest())

		enableMonsterAIV2 := os.Getenv("ENABLE_MONSTER_AI_V2") == "true"

		// Monster endpoints (protected)
		r.Route("/api/monsters", func(r chi.Router) {
			r.Post("/", handlers.monsterHandler.createMonster())
			if enableMonsterAIV2 {
				r.Post("/preview", handlers.monsterHandler.previewMonster())
			}
			r.Get("/{monsterID}", handlers.monsterHandler.getMonster())
			r.Put("/{monsterID}", handlers.monsterHandler.updateMonster())
			r.Delete("/{monsterID}", handlers.monsterHandler.deleteMonster())
			if enableMonsterAIV2 {
				r.Post("/{monsterID}/regenerate-section", handlers.monsterHandler.regenerateSection())
				r.Post("/{monsterID}/variant", handlers.monsterHandler.createVariant())
				r.Post("/{monsterID}/duplicate", handlers.monsterHandler.duplicateMonster())
			}
		})
		r.Get("/api/user/{userID}/monsters", handlers.monsterHandler.getMonstersByUser())
		if enableMonsterAIV2 {
			r.Get("/api/user/{userID}/monsters/generation-summary", handlers.monsterHandler.getGenerationSummary())
		}

		// Party endpoints (protected)
		r.Route("/api/parties", func(r chi.Router) {
			r.Post("/", handlers.partyHandler.createParty())
			r.Get("/{partyID}", handlers.partyHandler.getParty())
			r.Put("/{partyID}", handlers.partyHandler.updateParty())
			r.Delete("/{partyID}", handlers.partyHandler.deleteParty())
			r.Post("/{partyID}/members", handlers.partyHandler.addCharacterToParty())
			r.Delete("/{partyID}/members/{characterID}", handlers.partyHandler.removeCharacterFromParty())
			r.Post("/{partyID}/identified-beasts", handlers.partyHandler.addIdentifiedBeast())
			r.Delete("/{partyID}/identified-beasts/{beastID}", handlers.partyHandler.removeIdentifiedBeast())
		})
		r.Get("/api/user/{userID}/parties", handlers.partyHandler.getPartiesByOwner())
		// Route to set/update a character's party affiliation
		r.Put("/api/characters/{characterID}/party", handlers.partyHandler.setCharacterParty())

		// Active ability use endpoints
		r.Post("/api/characters/{characterID}/traits/{traitID}/use", handlers.abilityHandler.useTraitAbility())
		r.Post("/api/characters/{characterID}/features/{featureID}/use", handlers.abilityHandler.useFeatureAbility())
	})
}
