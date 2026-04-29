// Normalizes class_resource_definitions.category strings.
//
// Usage (from backend/):
//
//	go run ./cmd/migrate_class_resources
//	go run ./cmd/migrate_class_resources -apply
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
	"github.com/rpupo63/faradhaven/backend/internal/migrationclass"
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

	changes, err := migrationclass.ComputeChanges(db)
	if err != nil {
		log.Fatalf("plan: %v", err)
	}

	fmt.Printf("class_resource_definitions rows with planned updates: %d\n", len(changes))
	if len(changes) > 0 && !*apply {
		fmt.Println("Dry-run (no writes). Pass -apply to persist.")
	}

	for _, c := range changes {
		fmt.Printf("[%s] class=%s key=%s\n", c.ID, c.ClassID, c.ResourceKey)
		for _, k := range sortedKeys(c.Updates) {
			fmt.Printf("  %s: %v\n", k, c.Updates[k])
		}
	}

	if *apply {
		if err := migrationclass.ApplyPlanned(db, changes); err != nil {
			log.Fatalf("apply: %v", err)
		}
		fmt.Println("Applied class_resource_definitions updates.")
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
