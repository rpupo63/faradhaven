package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/internal/bootstrap"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed"
)

func main() {
	// CLI flags
	migrateOnly := flag.Bool("migrate-only", false, "Run migrations without seeding")
	clearOnly := flag.Bool("clear-only", false, "Clear all seeded data without reseeding")
	migrateUUIDs := flag.Bool("migrate-uuids", false, "Migrate existing entities to stable UUIDs (one-time)")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("Info: No .env file found: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set. Set it in .env or environment.")
	}

	db, err := bootstrap.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Running migrations...")
	if err := database.MigrateSpellComponentsOrdered(db); err != nil {
		log.Fatalf("spell_components migration failed: %v", err)
	}
	if err := db.AutoMigrate(models.AllModels()...); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations completed.")

	if *migrateOnly {
		return
	}

	// One-time UUID migration (preserves character references)
	if *migrateUUIDs || seed.NeedsUUIDMigration(db) {
		fmt.Println("Migrating to stable UUIDs...")
		if err := seed.MigrateToStableUUIDs(db); err != nil {
			log.Fatalf("UUID migration failed: %v", err)
		}
		fmt.Println("UUID migration completed.")
		if *migrateUUIDs {
			return
		}
	}

	// Initialize seeder with all registered seeds
	seeder := seed.NewSeeder(db)
	seeder.RegisterAll(seed.AllSeeds())

	if *clearOnly {
		fmt.Println("Clearing seeded child data...")
		if err := seeder.ClearAllData(); err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
		fmt.Println("Child data cleared.")
		return
	}

	// Default: clear child data and reseed
	fmt.Println("Reseeding database...")
	if err := seeder.ClearAndSeed(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}
	fmt.Println("Seeding completed successfully.")
}
