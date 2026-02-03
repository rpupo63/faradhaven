package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rpupo63/unified-personal-site-backend/internal/bootstrap"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_classes"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_races"
)

func main() {
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
	if err := db.AutoMigrate(models.AllModels()...); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Seeding Faradhaven races...")
	if err := faradhaven_races.SeedFaradhavenRaces(db); err != nil {
		log.Fatalf("Race seed failed: %v", err)
	}

	fmt.Println("Seeding Faradhaven classes...")
	if err := faradhaven_classes.SeedFaradhavenClasses(db); err != nil {
		log.Fatalf("Class seed failed: %v", err)
	}
	fmt.Println("Seed completed successfully.")
}
