package database

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type Database struct {
	db                     *gorm.DB
	userRepo               *UserRepo
	characterRepo          *CharacterRepo
	raceRepo               *RaceRepo
	classRepo              *ClassRepo
	archetypeRepo          *ArchetypeRepo
	spellRepo              *SpellRepo
	beastRepo              *BeastRepo
	attackRepo             *AttackRepo
	levelUpHistoryRepo     *LevelUpHistoryRepo
	weaponRepo             *WeaponRepo
	itemRepo               *ItemRepo
	componentRepo          *ComponentRepo
	effectRepo             *EffectRepo
	characterEffectRepo    *CharacterEffectRepo
	characterResourceRepo  *CharacterResourceRepo
	minionRepo             *MinionRepo
	consumptionHistoryRepo *ConsumptionHistoryRepo
	savedSpellRepo         *SavedSpellRepo
	corpseRepo             *CorpseRepo
	characterLinkRepo      *CharacterLinkRepo
	noteRepo               *NoteRepo
	gameMapRepo            *GameMapRepo
	mapTokenRepo           *MapTokenRepo
	mapElementRepo         *MapElementRepo
	monsterRepo            *MonsterRepo
	partyRepo              *PartyRepo // NEW: Party Repository

	// New Map Element Properties Repos
	trapPropertiesRepo             *TrapPropertiesRepo
	difficultTerrainPropertiesRepo *DifficultTerrainPropertiesRepo
	elevationPropertiesRepo        *ElevationPropertiesRepo
	wallPropertiesRepo             *WallPropertiesRepo
}

// New initializes a new Database struct with each repository using a shared GORM database instance
func New(db *gorm.DB) Database {
	return Database{
		db:                     db,
		userRepo:               NewUserRepo(db),
		characterRepo:          NewCharacterRepo(db),
		raceRepo:               NewRaceRepo(db),
		classRepo:              NewClassRepo(db),
		archetypeRepo:          NewArchetypeRepo(db),
		spellRepo:              NewSpellRepo(db),
		beastRepo:              NewBeastRepo(db),
		attackRepo:             NewAttackRepo(db),
		levelUpHistoryRepo:     NewLevelUpHistoryRepo(db),
		weaponRepo:             NewWeaponRepo(db),
		itemRepo:               NewItemRepo(db),
		componentRepo:          NewComponentRepo(db),
		effectRepo:             NewEffectRepo(db),
		characterEffectRepo:    NewCharacterEffectRepo(db),
		characterResourceRepo:  NewCharacterResourceRepo(db),
		consumptionHistoryRepo: NewConsumptionHistoryRepo(db),
		minionRepo:             NewMinionRepo(db),
		savedSpellRepo:         NewSavedSpellRepo(db),
		corpseRepo:             NewCorpseRepo(db),
		characterLinkRepo:      NewCharacterLinkRepo(db),
		noteRepo:               NewNoteRepo(db),
		gameMapRepo:            NewGameMapRepo(db),
		mapTokenRepo:           NewMapTokenRepo(db),
		mapElementRepo:         NewMapElementRepo(db),
		monsterRepo:            NewMonsterRepo(db),
		partyRepo:              NewPartyRepo(db), // NEW: Party Repository Initialization

		// New Map Element Properties Repos
		trapPropertiesRepo:             NewTrapPropertiesRepo(db),
		difficultTerrainPropertiesRepo: NewDifficultTerrainPropertiesRepo(db),
		elevationPropertiesRepo:        NewElevationPropertiesRepo(db),
		wallPropertiesRepo:             NewWallPropertiesRepo(db),
	}
}

// Accessor methods for each repository

func (d Database) UserRepo() *UserRepo {
	return d.userRepo
}

func (d Database) CharacterRepo() *CharacterRepo {
	return d.characterRepo
}

func (d Database) RaceRepo() *RaceRepo {
	return d.raceRepo
}

func (d Database) ClassRepo() *ClassRepo {
	return d.classRepo
}

func (d Database) ArchetypeRepo() *ArchetypeRepo {
	return d.archetypeRepo
}

func (d Database) SpellRepo() *SpellRepo {
	return d.spellRepo
}

func (d Database) BeastRepo() *BeastRepo {
	return d.beastRepo
}

func (d Database) AttackRepo() *AttackRepo {
	return d.attackRepo
}

func (d Database) LevelUpHistoryRepo() *LevelUpHistoryRepo {
	return d.levelUpHistoryRepo
}

func (d Database) WeaponRepo() *WeaponRepo {
	return d.weaponRepo
}

func (d Database) ItemRepo() *ItemRepo {
	return d.itemRepo
}

func (d Database) ComponentRepo() *ComponentRepo {
	return d.componentRepo
}

func (d Database) EffectRepo() *EffectRepo {
	return d.effectRepo
}

func (d Database) CharacterEffectRepo() *CharacterEffectRepo {
	return d.characterEffectRepo
}

func (d Database) CharacterResourceRepo() *CharacterResourceRepo {
	return d.characterResourceRepo
}

func (d Database) MinionRepo() *MinionRepo {
	return d.minionRepo
}

func (d Database) SavedSpellRepo() *SavedSpellRepo {
	return d.savedSpellRepo
}

func (d Database) CorpseRepo() *CorpseRepo {
	return d.corpseRepo
}

func (d Database) CharacterLinkRepo() *CharacterLinkRepo {
	return d.characterLinkRepo
}

func (d Database) NoteRepo() *NoteRepo {
	return d.noteRepo
}

func (d Database) ConsumptionHistoryRepo() *ConsumptionHistoryRepo {
	return d.consumptionHistoryRepo
}

func (d Database) GameMapRepo() *GameMapRepo {
	return d.gameMapRepo
}

func (d Database) MapTokenRepo() *MapTokenRepo {
	return d.mapTokenRepo
}

func (d Database) MapElementRepo() *MapElementRepo {
	return d.mapElementRepo
}

func (d Database) MonsterRepo() *MonsterRepo {
	return d.monsterRepo
}

func (d Database) PartyRepo() *PartyRepo {
	return d.partyRepo
}

// Accessor methods for new Map Element Properties Repos
func (d Database) TrapPropertiesRepo() *TrapPropertiesRepo {
	return d.trapPropertiesRepo
}

func (d Database) DifficultTerrainPropertiesRepo() *DifficultTerrainPropertiesRepo {
	return d.difficultTerrainPropertiesRepo
}

func (d Database) ElevationPropertiesRepo() *ElevationPropertiesRepo {
	return d.elevationPropertiesRepo
}

func (d Database) WallPropertiesRepo() *WallPropertiesRepo {
	return d.wallPropertiesRepo
}

// AutoMigrate runs GORM auto-migration for all models
func (d Database) AutoMigrate() error {
	if err := d.db.AutoMigrate(models.AllModels()...); err != nil {
		return err
	}

	// Drop legacy columns that GORM AutoMigrate won't remove
	if d.db.Migrator().HasColumn(&models.ClassLevel{}, "features") {
		if err := d.db.Migrator().DropColumn(&models.ClassLevel{}, "features"); err != nil {
			return err
		}
	}

	return nil
}

// DB returns the underlying GORM database instance
func (d Database) DB() *gorm.DB {
	return d.db
}

// Transaction executes a function within a database transaction
func (d Database) Transaction(fn func(tx *gorm.DB) error) error {
	return d.db.Transaction(fn)
}
