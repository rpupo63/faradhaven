// Deprecated: use go run ./cmd/migrate_spells -damage-only (same behavior for the damage phase).
//
// Migrates spells from legacy damage_dice text columns to damage_dice_count + damage_die_size,
// normalizes save_attr / ai fields, drops legacy text columns.
//
// Usage (from backend/):
//
//	go run ./cmd/migrate_spell_damage_fields           # dry-run
//	go run ./cmd/migrate_spell_damage_fields -apply    # migrate and drop legacy columns
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
	apply := flag.Bool("apply", false, "apply migration and drop legacy columns")
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

	opts := migrationspell.Options{Apply: *apply, DamageOnly: true}

	var legacyColCount int64
	if err := db.Raw(`SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = current_schema() AND table_name = 'spells' AND column_name = 'damage_dice'`).Scan(&legacyColCount).Error; err != nil {
		log.Fatalf("schema check: %v", err)
	}
	if legacyColCount == 0 {
		fmt.Println("No legacy column damage_dice on spells — nothing to migrate from text columns.")
		fmt.Println("(Tip: go run ./cmd/migrate_spells still normalizes mechanics + ensures integer dice columns.)")
		if *apply {
			if err := migrationspell.EnsureDamageDiceColumns(db); err != nil {
				log.Fatalf("ensure columns: %v", err)
			}
			fmt.Println("Ensured integer dice columns exist.")
		}
		return
	}

	changes, err := migrationspell.ComputeChanges(db, opts)
	if err != nil {
		log.Fatalf("plan: %v", err)
	}

	total, _ := migrationspell.CountSpells(db)
	fmt.Printf("Found %d spell row(s); %d with planned updates\n", total, len(changes))
	if len(changes) > 0 && !*apply {
		fmt.Println("Dry-run (no writes). Pass -apply to migrate and drop legacy columns.")
	}

	for _, c := range changes {
		fmt.Printf("[%s] %q\n", c.ID, c.Name)
		for k, v := range c.Updates {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	if *apply {
		if err := migrationspell.ApplyPlanned(db, opts, changes); err != nil {
			log.Fatalf("apply: %v", err)
		}
		fmt.Println("Dropped legacy text columns (if present): damage_dice, suggested_damage_dice, ai_recommended_damage_dice")
	}

	fmt.Printf("\nDone. %d spell(s) with planned changes.\n", len(changes))
}
