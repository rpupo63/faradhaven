package seed

import (
	"fmt"
	"log"
	"sort"

	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// Seed represents a single versioned seed operation.
type Seed struct {
	// Name is the unique identifier for this seed (e.g., "20240101_initial_races").
	// Use a date prefix for ordering: YYYYMMDD_description
	Name string
	// Run executes the seed. Return nil if successful.
	Run func(db *gorm.DB) error
}

// Seeder manages versioned seed operations.
type Seeder struct {
	db    *gorm.DB
	seeds []Seed
}

// NewSeeder creates a new seeder with the given database connection.
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// Register adds a seed to the seeder. Seeds are run in alphabetical order by name,
// so use date prefixes (YYYYMMDD_) for proper ordering.
func (s *Seeder) Register(seed Seed) {
	s.seeds = append(s.seeds, seed)
}

// RegisterAll adds multiple seeds at once.
func (s *Seeder) RegisterAll(seeds []Seed) {
	s.seeds = append(s.seeds, seeds...)
}

// Run executes all pending seeds that haven't been applied yet.
// Returns the number of seeds applied and any error encountered.
func (s *Seeder) Run() (int, error) {
	// Ensure seed_versions table exists
	if err := s.db.AutoMigrate(&models.SeedVersion{}); err != nil {
		return 0, fmt.Errorf("failed to create seed_versions table: %w", err)
	}

	// Sort seeds by name for consistent ordering
	sort.Slice(s.seeds, func(i, j int) bool {
		return s.seeds[i].Name < s.seeds[j].Name
	})

	applied := 0
	for _, seed := range s.seeds {
		// Check if already applied
		var existing models.SeedVersion
		err := s.db.Where("name = ?", seed.Name).First(&existing).Error
		if err == nil {
			log.Printf("Seed already applied: %s", seed.Name)
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return applied, fmt.Errorf("failed to check seed status: %w", err)
		}

		// Run the seed
		log.Printf("Running seed: %s", seed.Name)
		if err := seed.Run(s.db); err != nil {
			return applied, fmt.Errorf("seed %s failed: %w", seed.Name, err)
		}

		// Record as applied
		if err := s.db.Create(&models.SeedVersion{Name: seed.Name}).Error; err != nil {
			return applied, fmt.Errorf("failed to record seed %s: %w", seed.Name, err)
		}

		applied++
		log.Printf("Seed completed: %s", seed.Name)
	}

	return applied, nil
}

// Pending returns the names of seeds that haven't been applied yet.
func (s *Seeder) Pending() ([]string, error) {
	if err := s.db.AutoMigrate(&models.SeedVersion{}); err != nil {
		return nil, err
	}

	var pending []string
	for _, seed := range s.seeds {
		var existing models.SeedVersion
		err := s.db.Where("name = ?", seed.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			pending = append(pending, seed.Name)
		} else if err != nil {
			return nil, err
		}
	}
	return pending, nil
}

// Applied returns the names of seeds that have been applied.
func (s *Seeder) Applied() ([]models.SeedVersion, error) {
	if err := s.db.AutoMigrate(&models.SeedVersion{}); err != nil {
		return nil, err
	}

	var applied []models.SeedVersion
	if err := s.db.Order("applied_at").Find(&applied).Error; err != nil {
		return nil, err
	}
	return applied, nil
}

// ForceRun runs a specific seed by name regardless of whether it was applied.
// Useful for updating seed data. Does NOT re-record the seed version.
func (s *Seeder) ForceRun(name string) error {
	for _, seed := range s.seeds {
		if seed.Name == name {
			log.Printf("Force running seed: %s", name)
			return seed.Run(s.db)
		}
	}
	return fmt.Errorf("seed not found: %s", name)
}

// Reset removes a seed version record, allowing it to run again on next Run().
func (s *Seeder) Reset(name string) error {
	result := s.db.Where("name = ?", name).Delete(&models.SeedVersion{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("seed not found: %s", name)
	}
	log.Printf("Reset seed: %s", name)
	return nil
}

// ResetAll removes all seed version records.
func (s *Seeder) ResetAll() error {
	if err := s.db.Exec("DELETE FROM seed_versions").Error; err != nil {
		return err
	}
	log.Printf("Reset all seeds")
	return nil
}
