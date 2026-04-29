// Normalizes map_tokens.token_type to pc | npc.
//
// Usage (from backend/):
//
//	go run ./cmd/migrate_map_tokens
//	go run ./cmd/migrate_map_tokens -apply
//
// Requires DATABASE_URL.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/joho/godotenv"
	"github.com/rpupo63/faradhaven/backend/internal/bootstrap"
	"github.com/rpupo63/faradhaven/backend/internal/migrationmap"
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

	changes, err := migrationmap.ComputeChanges(db)
	if err != nil {
		log.Fatalf("plan: %v", err)
	}

	fmt.Printf("map_tokens rows with planned updates: %d\n", len(changes))
	if len(changes) > 0 && !*apply {
		fmt.Println("Dry-run (no writes). Pass -apply to persist.")
	}

	for _, c := range changes {
		fmt.Printf("[%s] map_id=%s\n", c.ID, c.MapID)
		for _, k := range sortedKeys(c.Updates) {
			fmt.Printf("  %s: %v\n", k, c.Updates[k])
		}
	}

	if *apply {
		if err := migrationmap.ApplyPlanned(db, changes); err != nil {
			log.Fatalf("apply: %v", err)
		}
		fmt.Println("Applied map_tokens updates.")
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
