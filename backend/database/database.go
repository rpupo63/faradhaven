package database

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type Database struct {
	db                 *gorm.DB
	userRepo           *UserRepo
	characterRepo      *CharacterRepo
	raceRepo           *RaceRepo
	classRepo          *ClassRepo
	archetypeRepo      *ArchetypeRepo
	spellRepo          *SpellRepo
	beastRepo          *BeastRepo
	attackRepo         *AttackRepo
	levelUpHistoryRepo *LevelUpHistoryRepo
	weaponRepo         *WeaponRepo
	itemRepo           *ItemRepo
	componentRepo      *ComponentRepo
	effectRepo         *EffectRepo
}

// New initializes a new Database struct with each repository using a shared GORM database instance
func New(db *gorm.DB) Database {
	return Database{
		db:                 db,
		userRepo:           NewUserRepo(db),
		characterRepo:      NewCharacterRepo(db),
		raceRepo:           NewRaceRepo(db),
		classRepo:          NewClassRepo(db),
		archetypeRepo:      NewArchetypeRepo(db),
		spellRepo:          NewSpellRepo(db),
		beastRepo:          NewBeastRepo(db),
		attackRepo:         NewAttackRepo(db),
		levelUpHistoryRepo: NewLevelUpHistoryRepo(db),
		weaponRepo:         NewWeaponRepo(db),
		itemRepo:           NewItemRepo(db),
		componentRepo:      NewComponentRepo(db),
		effectRepo:         NewEffectRepo(db),
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

// AutoMigrate runs GORM auto-migration for all models
func (d Database) AutoMigrate() error {
	return d.db.AutoMigrate(models.AllModels()...)
}

// DB returns the underlying GORM database instance
func (d Database) DB() *gorm.DB {
	return d.db
}

// Transaction executes a function within a database transaction
func (d Database) Transaction(fn func(tx *gorm.DB) error) error {
	return d.db.Transaction(fn)
}
