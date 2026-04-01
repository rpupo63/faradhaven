package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rpupo63/unified-personal-site-backend/internal/bootstrap"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
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

	var spells []*models.Spell
	if err := db.Preload("ComponentLinks", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Preload("ComponentLinks.Component").Find(&spells).Error; err != nil {
		log.Fatalf("Failed to query spells: %v", err)
	}
	for _, sp := range spells {
		sp.HydrateComponentsFromLinks()
	}

	total := len(spells)
	fmt.Printf("Found %d spell(s)\n", total)

	updated := 0
	skipped := 0

	for _, spell := range spells {
		correctLevel := len(spell.Components)
		if correctLevel == 0 {
			correctLevel = 1 // floor at 1 for spells with no components
		}

		if spell.Level == correctLevel {
			skipped++
			continue
		}

		fmt.Printf("  %-40s  level %d → %d  (components: %d)\n", spell.Name, spell.Level, correctLevel, len(spell.Components))

		if err := db.Model(&models.Spell{}).Where("id = ?", spell.ID).Updates(map[string]interface{}{
			"level":      correctLevel,
			"updated_at": time.Now(),
		}).Error; err != nil {
			log.Printf("ERROR updating spell %s (%s): %v", spell.Name, spell.ID, err)
			continue
		}
		updated++
	}

	fmt.Printf("\nDone. %d updated, %d already correct (out of %d total)\n", updated, skipped, total)
}
