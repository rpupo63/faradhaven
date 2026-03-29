// Normalizes weapon_damages.damage_type and damage_category (canonical enum strings).
//
// Run order (recommended): migrate_spells first if migrating spells; then migrate_weapon_damages,
// migrate_creatures, migrate_map_tokens, migrate_class_resources (independent).
//
// Usage (from backend/):
//
//	go run ./cmd/migrate_weapon_damages
//	go run ./cmd/migrate_weapon_damages -apply
//
// Set DATABASE_URL. Optional safety: pg_dump or CSV export of weapon_damages before -apply.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/joho/godotenv"
	"github.com/rpupo63/unified-personal-site-backend/internal/bootstrap"
	"github.com/rpupo63/unified-personal-site-backend/internal/migrationweapon"
)

func main() {
	apply := flag.Bool("apply", false, "write normalized values")
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

	changes, err := migrationweapon.ComputeChanges(db)
	if err != nil {
		log.Fatalf("plan: %v", err)
	}

	fmt.Printf("weapon_damages rows with planned updates: %d\n", len(changes))
	if len(changes) > 0 && !*apply {
		fmt.Println("Dry-run (no writes). Pass -apply to persist.")
	}

	for _, c := range changes {
		fmt.Printf("[%s] weapon_id=%s\n", c.ID, c.WeaponID)
		for _, k := range sortedKeys(c.Updates) {
			fmt.Printf("  %s: %v\n", k, c.Updates[k])
		}
	}

	if *apply {
		if err := migrationweapon.ApplyPlanned(db, changes); err != nil {
			log.Fatalf("apply: %v", err)
		}
		fmt.Println("Applied weapon_damages updates.")
	}
	fmt.Println("Done.")
}

func sortedKeys(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
