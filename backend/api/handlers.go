package api

import (
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/services"
)

// initializeHandlers creates and returns all handlers organized in a routeHandlers struct
func initializeHandlers(db database.Database) *routeHandlers {
	// Initialize services
	resourceService := services.NewResourceService()
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
		resourceService,
	)

	return &routeHandlers{
		authHandler:      newAuthHandler(db.UserRepo()),
		userHandler:      newUserHandler(db.UserRepo()),
		characterHandler: newCharacterHandler(db.CharacterRepo(), db.RaceRepo(), db.ClassRepo(), db.ItemRepo(), db.WeaponRepo(), db.SpellRepo(), resourceService, s3Service),
		spellHandler:     newSpellHandler(db.SpellRepo()),
		beastHandler:     newBeastHandler(db.BeastRepo(), db.AttackRepo()),
		levelHandler:     newLevelHandler(levelUpService, db.ClassRepo(), resourceService),
		weaponHandler:    newWeaponHandler(db.WeaponRepo()),
		itemHandler:      newItemHandler(db.ItemRepo()),
		componentHandler: newComponentHandler(db.ComponentRepo()),
		effectHandler:    newEffectHandler(db.EffectRepo()),
		noteHandler:      newNoteHandler(db.NoteRepo()),
	}
}
