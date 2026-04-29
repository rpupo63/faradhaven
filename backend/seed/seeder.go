package seed

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/seed/versioning"
	"gorm.io/gorm"
)

// Seed represents a single seed operation.
type Seed struct {
	// Name is the identifier for this seed (e.g., "races", "classes").
	Name string
	// Run executes the seed. Return nil if successful.
	Run func(db *gorm.DB) error
	// HashData returns the data to be hashed for version checking.
	// If nil, the seed will always run.
	HashData func() (interface{}, int)
}

// Seeder manages seed operations.
type Seeder struct {
	db    *gorm.DB
	seeds []Seed
}

// NewSeeder creates a new seeder with the given database connection.
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// Register adds a seed to the seeder.
func (s *Seeder) Register(seed Seed) {
	s.seeds = append(s.seeds, seed)
}

// RegisterAll adds multiple seeds at once.
func (s *Seeder) RegisterAll(seeds []Seed) {
	s.seeds = append(s.seeds, seeds...)
}

// ClearAllData deletes seeded child data from the database.
// Parent tables (races, classes, archetypes, components, weapons) are NOT cleared
// because characters reference them. The seed functions will update these in place.
// Only child/dependent tables are cleared so they can be re-seeded fresh.
func (s *Seeder) ClearAllData() error {
	log.Println("Clearing seeded child data...")

	// Only clear child tables that don't have character references
	// Order is CRITICAL: must delete children before parents for foreign keys.
	tablesToClear := []string{
		// Effect tables
		"effects",

		// Weapon child tables
		"weapon_damages",

		// Race child tables
		"trait_options",
		"traits",
		"race_components",
		"lineages",

		// Class child tables (archetypes excluded — characters reference them)
		"class_level_resources",
		"level_features",
		"class_starting_equipment_options",
		"class_starting_equipment_choices",
		"class_weapon_requirements",
		"class_resource_definitions",
		"class_components",
		"class_levels",
	}

	for _, table := range tablesToClear {
		// Use DELETE instead of TRUNCATE to avoid CASCADE issues
		if err := s.db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			return fmt.Errorf("could not clear table %s: %w", table, err)
		}
		log.Printf("  Cleared table: %s", table)
	}

	// Also clear seed metadata to force full re-seed
	if err := versioning.ClearSeedMetadata(s.db); err != nil {
		return fmt.Errorf("could not clear seed_metadata: %w", err)
	}

	log.Println("Child data cleared successfully")
	return nil
}

// ClearAndSeed clears all seeded data and runs all seeds from scratch.
// Each seed runs in its own transaction for atomic failure handling.
func (s *Seeder) ClearAndSeed() error {
	// Clear all existing seeded data
	if err := s.ClearAllData(); err != nil {
		return fmt.Errorf("failed to clear data: %w", err)
	}

	return s.runSeeds(true) // Force all seeds to run
}

// RunAll runs all registered seeds (without clearing first).
// Uses version checking to skip unchanged seeds.
// Useful for seeding an empty database or updating after changes.
func (s *Seeder) RunAll() error {
	return s.runSeeds(false) // Use version checking
}

// runSeeds executes all seeds with optional version checking.
// Each seed runs in its own transaction.
func (s *Seeder) runSeeds(force bool) error {
	// Sort seeds by name for consistent ordering
	sort.Slice(s.seeds, func(i, j int) bool {
		return s.seeds[i].Name < s.seeds[j].Name
	})

	for _, seed := range s.seeds {
		var contentHash string
		var recordCount int

		// Check if seed needs to run (version check)
		if !force && seed.HashData != nil {
			data, count := seed.HashData()
			recordCount = count

			hash, err := versioning.ComputeSeedHash(data)
			if err != nil {
				log.Printf("Warning: could not compute hash for %s: %v", seed.Name, err)
			} else {
				contentHash = hash
				shouldRun, err := versioning.ShouldRunSeed(s.db, seed.Name, hash)
				if err != nil {
					log.Printf("Warning: could not check version for %s: %v", seed.Name, err)
				} else if !shouldRun {
					log.Printf("Skipping seed (unchanged): %s", seed.Name)
					continue
				}
			}
		}

		log.Printf("Running seed: %s", seed.Name)
		start := time.Now()

		// Run seed in a transaction
		err := s.db.Transaction(func(tx *gorm.DB) error {
			return seed.Run(tx)
		})

		duration := time.Since(start)

		if err != nil {
			// Record failed seed
			if contentHash != "" {
				_ = versioning.RecordSeedRun(s.db, seed.Name, contentHash, recordCount, duration, "failed")
			}
			return fmt.Errorf("seed %s failed: %w", seed.Name, err)
		}

		// Record successful seed
		if contentHash != "" {
			if err := versioning.RecordSeedRun(s.db, seed.Name, contentHash, recordCount, duration, "success"); err != nil {
				log.Printf("Warning: could not record seed run for %s: %v", seed.Name, err)
			}
		}

		log.Printf("Seed completed: %s (%dms)", seed.Name, duration.Milliseconds())
	}

	return nil
}

// HasData checks if any seeded data exists in the database.
func (s *Seeder) HasData() bool {
	var count int64

	// Check if races exist (quick indicator of seeded data)
	s.db.Model(&models.Race{}).Count(&count)
	if count > 0 {
		return true
	}

	// Check if classes exist
	s.db.Model(&models.Class{}).Count(&count)
	return count > 0
}
