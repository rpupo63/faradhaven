package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	api "github.com/rpupo63/unified-personal-site-backend/api"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/internal/bootstrap"
)

// @title           Faradhaven API
// @version         1.0
// @description     API for the Faradhaven game - managing characters, spells, and beasts

// @host      localhost:8080
// @BasePath  /

// @schemes   http https

func main() {
	fmt.Println("Initializing Faradhaven backend...")

	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Info: No .env file found (using system environment variables): %v\n", err)
	}

	// 2. Get Connection String
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Error: DATABASE_URL not set. Please set it in your .env file.")
	}

	fmt.Println("Connecting to database...")

	// 3, 4, 5. Connect, Extensions, Pool (via Bootstrap)
	db, err := bootstrap.InitDB(dsn)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	fmt.Println("Connected to database successfully")

	// 6. Wrap for Application
	currentDB := database.New(db)

	// 7. Run migrations based on GENERATE_MODELS env var
	// GENERATE_MODELS=only  -> run migrations and exit (for CI/CD or standalone migration)
	// GENERATE_MODELS=true  -> run migrations and continue to start server
	// GENERATE_MODELS=false or unset -> skip migrations
	migrateMode := strings.ToLower(strings.TrimSpace(os.Getenv("GENERATE_MODELS")))
	if migrateMode == "only" || migrateMode == "true" {
		fmt.Println("Running database migrations...")
		if err := currentDB.AutoMigrate(); err != nil {
			log.Fatalf("Error running migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")

		if migrateMode == "only" {
			fmt.Println("GENERATE_MODELS=only set, exiting after migrations.")
			return
		}
	}

	// 8. Start Server
	startServer(currentDB)
}

func startServer(db database.Database) {
	errChannel := make(chan error)
	defer close(errChannel)

	server, err := api.NewServer(db)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	go server.Start(errChannel)
	go listenToInterrupt(errChannel)

	fatalErr := <-errChannel
	fmt.Printf("Closing server: %v\n", fatalErr)

	server.ShutdownGracefully(30 * time.Second)
}

func listenToInterrupt(errChannel chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errChannel <- fmt.Errorf("%s", <-c)
}
