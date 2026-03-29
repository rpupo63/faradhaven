// Deprecated: use go run ./cmd/migrate_spells -mechanics-only (same behavior).
//
// Migrates legacy spell rows: SpellType, integer range (feet), duration strings per
// models.ValidateSpellDuration, and AI recommended type/range/duration.
//
// Usage (from backend/):
//
//	go run ./cmd/migrate_spell_mechanics           # dry-run
//	go run ./cmd/migrate_spell_mechanics -apply    # write updates
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rpupo63/unified-personal-site-backend/internal/bootstrap"
	"github.com/rpupo63/unified-personal-site-backend/internal/migrationspell"
)

func main() {
	apply := flag.Bool("apply", false, "apply updates to the database (default: dry-run only)")
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

	opts := migrationspell.Options{Apply: *apply, MechanicsOnly: true}

	changes, err := migrationspell.ComputeChanges(db, opts)
	if err != nil {
		log.Fatalf("plan: %v", err)
	}

	total, _ := migrationspell.CountSpells(db)
	fmt.Printf("Found %d spell row(s); %d with planned updates\n", total, len(changes))
	if len(changes) > 0 && !*apply {
		fmt.Println("Dry-run (no writes). Pass -apply to update the database.")
	}

	for _, c := range changes {
		fmt.Printf("[%s] %q\n", c.ID, c.Name)
		for k, v := range c.Updates {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}
	if *apply && len(changes) > 0 {
		if err := migrationspell.ApplyPlanned(db, opts, changes); err != nil {
			log.Fatalf("apply: %v", err)
		}
		fmt.Println("Updates applied.")
	}

	fmt.Printf("\nDone. %d spell(s) with planned changes", len(changes))
	if !*apply && len(changes) > 0 {
		fmt.Print(". Re-run with -apply to persist.")
	}
	fmt.Println()
}
