package seed

import (
	"fmt"
	"log"
	"sort"

	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// Seed represents a single seed operation.
type Seed struct {
	// Name is the identifier for this seed (e.g., "races", "classes").
	Name string
	// Run executes the seed. Return nil if successful.
	Run func(db *gorm.DB) error
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
	// Parent tables (races, classes, components, weapons, archetypes) are updated in-place by seeds
	tablesToClear := []string{
		// Weapon child tables
		"weapon_damages",

		// Class child tables (not the classes themselves - characters reference them)
		"class_components",
		"level_features",
		"class_levels",

		// Race child tables (not the races themselves - characters reference them)
		"race_components",
		"trait_options",
		"traits",
		"lineages",

		// Drop the old seed_versions table if it exists
		"seed_versions",
	}

	for _, table := range tablesToClear {
		// Use DELETE instead of TRUNCATE to avoid CASCADE issues
		if err := s.db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			log.Printf("  Note: Could not clear table %s (may not exist): %v", table, err)
		} else {
			log.Printf("  Cleared table: %s", table)
		}
	}

	log.Println("Child data cleared")
	return nil
}

// ClearAndSeed clears all seeded data and runs all seeds from scratch.
// This is a simple approach: delete everything, then reseed.
func (s *Seeder) ClearAndSeed() error {
	// Clear all existing seeded data
	if err := s.ClearAllData(); err != nil {
		return fmt.Errorf("failed to clear data: %w", err)
	}

	// Sort seeds by name for consistent ordering
	sort.Slice(s.seeds, func(i, j int) bool {
		return s.seeds[i].Name < s.seeds[j].Name
	})

	// Run all seeds
	for _, seed := range s.seeds {
		log.Printf("Running seed: %s", seed.Name)
		if err := seed.Run(s.db); err != nil {
			return fmt.Errorf("seed %s failed: %w", seed.Name, err)
		}
		log.Printf("Seed completed: %s", seed.Name)
	}

	return nil
}

// RunAll runs all registered seeds (without clearing first).
// Useful for seeding an empty database.
func (s *Seeder) RunAll() error {
	// Sort seeds by name for consistent ordering
	sort.Slice(s.seeds, func(i, j int) bool {
		return s.seeds[i].Name < s.seeds[j].Name
	})

	for _, seed := range s.seeds {
		log.Printf("Running seed: %s", seed.Name)
		if err := seed.Run(s.db); err != nil {
			return fmt.Errorf("seed %s failed: %w", seed.Name, err)
		}
		log.Printf("Seed completed: %s", seed.Name)
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
