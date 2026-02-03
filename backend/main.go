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
	"github.com/rpupo63/unified-personal-site-backend/models"
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

	// 7. Run migrations if GENERATE_MODELS=true (then exit)
	if strings.ToLower(strings.TrimSpace(os.Getenv("GENERATE_MODELS"))) == "true" {
		fmt.Println("Generating models and running migrations...")
		models.GenerateModels(db)
		return
	}

	// 8. Auto-migrate if AUTO_MIGRATE=true (for local development)
	if strings.ToLower(strings.TrimSpace(os.Getenv("AUTO_MIGRATE"))) == "true" {
		fmt.Println("Running auto-migration...")
		if err := currentDB.AutoMigrate(); err != nil {
			log.Fatalf("Error running auto-migration: %v", err)
		}
		fmt.Println("Auto-migration completed successfully")
	}

	// 9. Start Server
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
