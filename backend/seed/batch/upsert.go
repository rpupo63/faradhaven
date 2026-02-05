package batch

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DefaultBatchSize = 100

// UpsertBatchUpdateAll performs batch upsert with ON CONFLICT DO UPDATE SET all columns.
// This is useful when you want to fully replace existing records.
func UpsertBatchUpdateAll[T any](tx *gorm.DB, records []T, batchSize int) error {
	if len(records) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = DefaultBatchSize
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).CreateInBatches(records, batchSize).Error
}

// UpsertBatchUpdateColumns performs batch upsert with ON CONFLICT DO UPDATE for specific columns.
// This is useful when you want to update only certain fields on conflict.
func UpsertBatchUpdateColumns[T any](tx *gorm.DB, records []T, batchSize int, columns []string) error {
	if len(records) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = DefaultBatchSize
	}

	assignmentColumns := make([]clause.Column, len(columns))
	for i, col := range columns {
		assignmentColumns[i] = clause.Column{Name: col}
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(columns),
	}).CreateInBatches(records, batchSize).Error
}

// InsertBatch performs batch insert without any conflict handling.
// Use this for child tables that are cleared before seeding.
func InsertBatch[T any](tx *gorm.DB, records []T, batchSize int) error {
	if len(records) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = DefaultBatchSize
	}
	return tx.CreateInBatches(records, batchSize).Error
}
