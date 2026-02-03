package bootstrap

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB connects to the database and sets up extensions
// It mimics the startup logic in main.go
func InitDB(dsn string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disable implicit prepared statement usage (required for connection pooling)
	}), &gorm.Config{
		Logger:      newLogger,
		PrepareStmt: false, // Disable prepared statements to avoid "already exists" errors with connection pooling
	})
	if err != nil {
		return nil, err
	}

	setupExtensions(db)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func setupExtensions(db *gorm.DB) {
	ctx := db.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err := ctx.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Error enabling uuid-ossp extension: %v", err)
	}

	if err := ctx.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		log.Printf("Error enabling vector extension: %v", err)
	}
}

// LoadEnv attempts to load .env from various locations
// Useful for tests running in subdirectories
func LoadEnv() {
	// Try current directory
	_ = godotenv.Load()
	// Try parent directory
	_ = godotenv.Load("../.env")
	// Try grand-parent directory
	_ = godotenv.Load("../../.env")
}
