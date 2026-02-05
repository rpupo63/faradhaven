package versioning

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// ComputeSeedHash computes a SHA256 hash of the seed data.
// The data is JSON-serialized before hashing for deterministic results.
func ComputeSeedHash(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:]), nil
}

// GetSeedMetadata retrieves the metadata for a seed by name.
// Returns nil if no record exists.
func GetSeedMetadata(db *gorm.DB, seedName string) (*models.SeedMetadata, error) {
	var meta models.SeedMetadata
	err := db.Where("seed_name = ?", seedName).First(&meta).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

// ShouldRunSeed checks if a seed needs to run based on content hash comparison.
// Returns true if the seed should run (hash differs or no previous run).
func ShouldRunSeed(db *gorm.DB, seedName, currentHash string) (bool, error) {
	meta, err := GetSeedMetadata(db, seedName)
	if err != nil {
		return true, err
	}
	if meta == nil {
		return true, nil // No previous run
	}
	if meta.Status != "success" {
		return true, nil // Previous run failed
	}
	return meta.ContentHash != currentHash, nil
}

// RecordSeedRun records the result of a seed execution.
// Updates existing record or creates new one.
func RecordSeedRun(db *gorm.DB, seedName, contentHash string, recordCount int, duration time.Duration, status string) error {
	meta, err := GetSeedMetadata(db, seedName)
	if err != nil {
		return err
	}

	durationMs := int(duration.Milliseconds())

	if meta == nil {
		// Create new record
		meta = &models.SeedMetadata{
			ID:          uuid.New(),
			SeedName:    seedName,
			ContentHash: contentHash,
			RecordCount: recordCount,
			ExecutedAt:  time.Now(),
			DurationMs:  durationMs,
			Status:      status,
		}
		return db.Create(meta).Error
	}

	// Update existing record
	return db.Model(meta).Updates(map[string]interface{}{
		"content_hash": contentHash,
		"record_count": recordCount,
		"executed_at":  time.Now(),
		"duration_ms":  durationMs,
		"status":       status,
	}).Error
}

// ClearSeedMetadata removes all seed metadata records.
// Useful for forcing a full re-seed.
func ClearSeedMetadata(db *gorm.DB) error {
	return db.Exec("DELETE FROM seed_metadata").Error
}
