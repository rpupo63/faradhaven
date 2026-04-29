package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	api "github.com/rpupo63/faradhaven/backend/api"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/internal/bootstrap"
	"github.com/rpupo63/faradhaven/backend/seed"
)

// @title           Faradhaven API
// @version         1.0
// @description     API for the Faradhaven game - managing characters, spells, and beasts

// @host      localhost:8080
// @BasePath  /

// @schemes   http https

func main() {
	fmt.Println("Initializing Faradhaven backend...")

	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Info: No .env file found (using system environment variables): %v\n", err)
	}

	// 2. Get Connection String
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Error: DATABASE_URL not set. Please set it in your .env file.")
	}

	fmt.Println("Connecting to database...")

	// 3, 4, 5. Connect, Extensions, Pool (via Bootstrap)
	db, err := bootstrap.InitDB(dsn)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	fmt.Println("Connected to database successfully")

	// 6. Wrap for Application
	currentDB := database.New(db)

	// 7. Run migrations and seeding based on GENERATE_MODELS env var
	// GENERATE_MODELS=only  -> run migrations/seeding and exit (for CI/CD or standalone)
	// GENERATE_MODELS=true  -> run migrations/seeding and continue to start server
	// GENERATE_MODELS=false or unset -> skip migrations and seeding
	migrateMode := strings.ToLower(strings.TrimSpace(os.Getenv("GENERATE_MODELS")))
	if migrateMode == "only" || migrateMode == "true" {
		fmt.Println("Running database migrations...")
		if err := currentDB.AutoMigrate(); err != nil {
			log.Fatalf("Error running migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")

		// Migrate existing entities to stable/deterministic UUIDs (one-time migration)
		// This preserves character references when reseeding
		if seed.NeedsUUIDMigration(db) {
			fmt.Println("Migrating to stable UUIDs...")
			if err := seed.MigrateToStableUUIDs(db); err != nil {
				log.Fatalf("Error migrating UUIDs: %v", err)
			}
			fmt.Println("UUID migration completed")
		}

		// Migrate character weapons to V2 (explicit join table)
		fmt.Println("Migrating character weapons to v2...")
		if err := seed.MigrateCharacterWeaponsV2(db); err != nil {
			log.Fatalf("Error migrating character weapons: %v", err)
		}
		fmt.Println("Character weapons migration completed")

		// Initialize class-specific resources for existing characters
		fmt.Println("Initializing character resources...")
		if err := seed.InitializeCharacterResources(db); err != nil {
			log.Fatalf("Error initializing character resources: %v", err)
		}
		fmt.Println("Character resources initialization completed")

		// Clear all seeded data and reseed from scratch
		fmt.Println("Clearing and reseeding database...")
		seeder := seed.NewSeeder(db)
		seeder.RegisterAll(seed.AllSeeds())
		if err := seeder.ClearAndSeed(); err != nil {
			log.Fatalf("Error seeding database: %v", err)
		}
		fmt.Println("Seeding completed successfully")

		if migrateMode == "only" {
			fmt.Println("GENERATE_MODELS=only set, exiting after migrations and seeding.")
			return
		}
	}

	// 8. Start Server
	startServer(currentDB)
}

func startServer(db database.Database) {
	errChannel := make(chan error)
	defer close(errChannel)

	server, err := api.NewServer(db)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	go server.Start(errChannel)
	go listenToInterrupt(errChannel)

	fatalErr := <-errChannel
	fmt.Printf("Closing server: %v\n", fatalErr)

	server.ShutdownGracefully(30 * time.Second)
}

func listenToInterrupt(errChannel chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errChannel <- fmt.Errorf("%s", <-c)
}
