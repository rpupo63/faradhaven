package models

import "time"

// SeedVersion tracks which seeds have been applied to the database.
// Similar to schema migrations but for data seeding.
type SeedVersion struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;not null"` // e.g., "20240101_initial_races"
	AppliedAt time.Time `gorm:"autoCreateTime"`
}
