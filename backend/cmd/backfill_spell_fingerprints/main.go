// Backfills the component_fingerprint column on spells that currently have NULL.
//
// Spells written before the fingerprint column was added are the target.
// The fingerprint encodes the canonical component sequence so two spells with the
// same components (accounting for Logica phase order) share an identical hash.
//
// Usage (from backend/):
//
//	go run ./cmd/backfill_spell_fingerprints              # dry-run: list planned updates
//	go run ./cmd/backfill_spell_fingerprints -apply       # write component_fingerprint for all NULL rows
//
// Requires DATABASE_URL in the environment or a .env file.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rpupo63/faradhaven/backend/internal/bootstrap"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

func main() {
	apply := flag.Bool("apply", false, "write computed fingerprints to the database")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("info: no .env in cwd: %v", err)
	}
	bootstrap.LoadEnv()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := bootstrap.InitDB(dsn)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	// Count total spells and those already fingerprinted.
	var total, already int64
	db.Model(&models.Spell{}).Count(&total)
	db.Model(&models.Spell{}).Where("component_fingerprint IS NOT NULL").Count(&already)
	fmt.Printf("Total spells: %d  |  Already fingerprinted: %d  |  Needs backfill: %d\n",
		total, already, total-already)

	// Load all spells missing a fingerprint, with components preloaded.
	var spells []*models.Spell
	err = db.
		Preload("ComponentLinks", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("ComponentLinks.Component").
		Where("component_fingerprint IS NULL").
		Find(&spells).Error
	if err != nil {
		log.Fatalf("query: %v", err)
	}

	if len(spells) == 0 {
		fmt.Println("Nothing to backfill.")
		return
	}

	// Compute fingerprints.
	type row struct {
		id          uuid.UUID
		name        string
		fingerprint string
		nComponents int
	}
	rows := make([]row, len(spells))
	for i, sp := range spells {
		fp := models.ComponentFingerprint(sp.ComponentLinks)
		rows[i] = row{
			id:          sp.ID,
			name:        sp.Name,
			fingerprint: fp,
			nComponents: len(sp.ComponentLinks),
		}
	}

	// Print plan.
	fmt.Printf("\nSpells to fingerprint: %d\n", len(rows))
	for _, r := range rows {
		fmt.Printf("  [%s] %q  components=%d  fingerprint=%s\n",
			r.id, r.name, r.nComponents, r.fingerprint[:12]+"…")
	}

	if !*apply {
		fmt.Println("\nDry-run (no writes). Pass -apply to persist.")
		return
	}

	// Apply updates in batches of 100.
	const batchSize = 100
	updated := 0
	for i := 0; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}
		batch := rows[i:end]
		err := db.Transaction(func(tx *gorm.DB) error {
			for _, r := range batch {
				if err := tx.Model(&models.Spell{}).
					Where("id = ?", r.id).
					Update("component_fingerprint", r.fingerprint).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("batch %d–%d: %v", i, end, err)
		}
		updated += len(batch)
		fmt.Printf("  updated %d/%d\n", updated, len(rows))
	}

	fmt.Printf("\nDone. %d spell(s) fingerprinted.\n", updated)
}
