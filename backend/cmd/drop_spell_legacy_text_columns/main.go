// Drops legacy spells table text columns used before dice were split into integers:
//
//	damage_dice, suggested_damage_dice, ai_recommended_damage_dice
//
// The normal path is: go run ./cmd/migrate_spells -apply
// (that migrates data into damage_dice_count / damage_die_size / etc. and drops these columns).
// Use this tool only when:
//   - You backfilled integers manually and still need to remove leftover text columns, or
//   - You want an explicit second step after verifying data (idempotent: DROP IF EXISTS).
//
// Usage (from backend/, DATABASE_URL set):
//
//	go run ./cmd/drop_spell_legacy_text_columns           # dry-run: show which columns exist
//	go run ./cmd/drop_spell_legacy_text_columns -apply    # DROP COLUMN IF EXISTS for each
//
// Does not delete data from the integer columns; only removes obsolete text columns when present.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rpupo63/faradhaven/backend/internal/bootstrap"
	"gorm.io/gorm"
)

var legacyCols = []string{
	"damage_dice",
	"suggested_damage_dice",
	"ai_recommended_damage_dice",
}

func main() {
	apply := flag.Bool("apply", false, "execute DROP COLUMN IF EXISTS for legacy text columns")
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

	existing, err := listExisting(db, "spells", legacyCols)
	if err != nil {
		log.Fatalf("inspect: %v", err)
	}

	if len(existing) == 0 {
		fmt.Println("No legacy spell dice text columns present (already dropped or never existed). Nothing to do.")
		return
	}

	fmt.Println("Legacy columns still present on spells:")
	for _, c := range existing {
		fmt.Printf("  - %s\n", c)
	}

	if !*apply {
		fmt.Println("\nDry-run only. Pass -apply to drop these columns (after confirming integer dice columns are populated).")
		return
	}

	for _, col := range existing {
		q := fmt.Sprintf(`ALTER TABLE spells DROP COLUMN IF EXISTS %s`, col)
		if err := db.Exec(q).Error; err != nil {
			log.Fatalf("drop %s: %v", col, err)
		}
		fmt.Printf("Dropped column: %s\n", col)
	}
	fmt.Println("Done.")
}

func listExisting(db *gorm.DB, table string, cols []string) ([]string, error) {
	var out []string
	for _, col := range cols {
		var n int64
		err := db.Raw(`
SELECT COUNT(*) FROM information_schema.columns
WHERE table_schema = current_schema() AND table_name = ? AND column_name = ?
`, table, col).Scan(&n).Error
		if err != nil {
			return nil, err
		}
		if n > 0 {
			out = append(out, col)
		}
	}
	return out, nil
}
