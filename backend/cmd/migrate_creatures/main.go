// Normalizes monsters (size, type) and corpses (creature_size, creature_type).
//
// Usage (from backend/):
//
//	go run ./cmd/migrate_creatures
//	go run ./cmd/migrate_creatures -apply
//
// Requires DATABASE_URL. Optional: pg_dump monsters,corpse tables before -apply.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/joho/godotenv"
	"github.com/rpupo63/faradhaven/backend/internal/bootstrap"
	"github.com/rpupo63/faradhaven/backend/internal/migrationcreature"
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

	mChanges, err := migrationcreature.ComputeMonsterChanges(db)
	if err != nil {
		log.Fatalf("plan monsters: %v", err)
	}
	cChanges, err := migrationcreature.ComputeCorpseChanges(db)
	if err != nil {
		log.Fatalf("plan corpses: %v", err)
	}

	fmt.Printf("monsters with planned updates: %d\n", len(mChanges))
	fmt.Printf("corpses with planned updates: %d\n", len(cChanges))
	if (len(mChanges) > 0 || len(cChanges) > 0) && !*apply {
		fmt.Println("Dry-run (no writes). Pass -apply to persist.")
	}

	for _, c := range mChanges {
		fmt.Printf("[monster %s] %q\n", c.ID, c.Name)
		for _, k := range sortedKeys(c.Updates) {
			fmt.Printf("  %s: %v\n", k, c.Updates[k])
		}
	}
	for _, c := range cChanges {
		fmt.Printf("[corpse %s] %q\n", c.ID, c.Name)
		for _, k := range sortedKeys(c.Updates) {
			fmt.Printf("  %s: %v\n", k, c.Updates[k])
		}
	}

	if *apply {
		if err := migrationcreature.ApplyMonsters(db, mChanges); err != nil {
			log.Fatalf("apply monsters: %v", err)
		}
		if err := migrationcreature.ApplyCorpses(db, cChanges); err != nil {
			log.Fatalf("apply corpses: %v", err)
		}
		fmt.Println("Applied monster and corpse updates.")
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
