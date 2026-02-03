package models

import (
	"log"

	"gorm.io/gorm"
)

// GenerateModels runs auto-migration for all models
func GenerateModels(db *gorm.DB) {
	log.Println("Running migrations for Faradhaven models...")

	err := db.AutoMigrate(AllModels()...)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully!")
}
