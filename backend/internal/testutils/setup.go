// Package testutils provides a shared, seeded Postgres database for API and
// seed tests. It connects to the server pointed at by TEST_DATABASE_URL
// (e.g. the backend/docker-compose.yml Postgres), drops and recreates a
// dedicated faradhaven_test database, migrates, and runs the full seed
// pipeline once per test binary.
package testutils

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/internal/bootstrap"
	"github.com/rpupo63/faradhaven/backend/seed"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const testDBName = "faradhaven_test"

var (
	once      sync.Once
	sharedDB  *gorm.DB
	sharedApp database.Database
	setupErr  error
)

// SetupSeededTestDB returns a migrated AND seeded database shared by every
// test in the binary. Tests must therefore create their own users/characters
// and never mutate seeded reference data (classes, races, components, ...).
//
// Skips when TEST_DATABASE_URL is unset:
//
//	docker compose -f backend/docker-compose.yml up -d
//	export TEST_DATABASE_URL=postgresql://postgres:postgres@localhost:5432/faradhaven?sslmode=disable
func SetupSeededTestDB(t *testing.T) (*gorm.DB, database.Database) {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set; start backend/docker-compose.yml postgres and set it to run DB-backed tests")
	}

	once.Do(func() {
		setupErr = initSharedDB(dsn)
	})
	if setupErr != nil {
		t.Fatalf("test database setup failed: %v", setupErr)
	}
	return sharedDB, sharedApp
}

func initSharedDB(adminDSN string) error {
	admin, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("connecting to admin database: %w", err)
	}

	if err := admin.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", testDBName)).Error; err != nil {
		return fmt.Errorf("dropping old test database: %w", err)
	}
	if err := admin.Exec("CREATE DATABASE " + testDBName).Error; err != nil {
		return fmt.Errorf("creating test database: %w", err)
	}
	if sqlDB, err := admin.DB(); err == nil {
		_ = sqlDB.Close()
	}

	u, err := url.Parse(adminDSN)
	if err != nil {
		return fmt.Errorf("parsing TEST_DATABASE_URL: %w", err)
	}
	u.Path = "/" + testDBName

	db, err := bootstrap.InitDB(u.String())
	if err != nil {
		return fmt.Errorf("initializing test database: %w", err)
	}

	app := database.New(db)
	if err := app.AutoMigrate(); err != nil {
		return fmt.Errorf("migrating test database: %w", err)
	}

	seeder := seed.NewSeeder(db)
	seeder.RegisterAll(seed.AllSeeds())
	if err := seeder.ClearAndSeed(); err != nil {
		return fmt.Errorf("seeding test database: %w", err)
	}

	sharedDB = db
	sharedApp = app
	return nil
}
