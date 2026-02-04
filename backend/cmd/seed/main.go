package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rpupo63/unified-personal-site-backend/internal/bootstrap"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed"
)

func main() {
	// CLI flags
	listPending := flag.Bool("pending", false, "List pending seeds")
	listApplied := flag.Bool("applied", false, "List applied seeds")
	forceRun := flag.String("force", "", "Force run a specific seed by name")
	resetSeed := flag.String("reset", "", "Reset a seed to allow re-running")
	resetAll := flag.Bool("reset-all", false, "Reset all seeds (dangerous!)")
	migrateOnly := flag.Bool("migrate-only", false, "Run migrations without seeding")
	flag.Parse()

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

	fmt.Println("Running migrations...")
	if err := db.AutoMigrate(models.AllModels()...); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if *migrateOnly {
		fmt.Println("Migrations completed (seeding skipped).")
		return
	}

	// Initialize seeder with all registered seeds
	seeder := seed.NewSeeder(db)
	seeder.RegisterAll(seed.AllSeeds())

	// Handle CLI commands
	switch {
	case *listPending:
		pending, err := seeder.Pending()
		if err != nil {
			log.Fatalf("Failed to list pending seeds: %v", err)
		}
		if len(pending) == 0 {
			fmt.Println("No pending seeds.")
		} else {
			fmt.Println("Pending seeds:")
			for _, name := range pending {
				fmt.Printf("  - %s\n", name)
			}
		}
		return

	case *listApplied:
		applied, err := seeder.Applied()
		if err != nil {
			log.Fatalf("Failed to list applied seeds: %v", err)
		}
		if len(applied) == 0 {
			fmt.Println("No seeds applied yet.")
		} else {
			fmt.Println("Applied seeds:")
			for _, sv := range applied {
				fmt.Printf("  - %s (applied: %s)\n", sv.Name, sv.AppliedAt.Format("2006-01-02 15:04:05"))
			}
		}
		return

	case *forceRun != "":
		if err := seeder.ForceRun(*forceRun); err != nil {
			log.Fatalf("Force run failed: %v", err)
		}
		fmt.Printf("Force ran seed: %s\n", *forceRun)
		return

	case *resetSeed != "":
		if err := seeder.Reset(*resetSeed); err != nil {
			log.Fatalf("Reset failed: %v", err)
		}
		fmt.Printf("Reset seed: %s (will run on next seed)\n", *resetSeed)
		return

	case *resetAll:
		fmt.Print("Are you sure you want to reset ALL seeds? This will cause them to re-run. (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Aborted.")
			return
		}
		if err := seeder.ResetAll(); err != nil {
			log.Fatalf("Reset all failed: %v", err)
		}
		fmt.Println("All seeds reset.")
		return
	}

	// Default: run pending seeds
	fmt.Println("Running pending seeds...")
	applied, err := seeder.Run()
	if err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	if applied == 0 {
		fmt.Println("No new seeds to apply.")
	} else {
		fmt.Printf("Applied %d seed(s) successfully.\n", applied)
	}
}
