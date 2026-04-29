// Package migrationclass normalizes class_resource_definitions.category to canonical values.
package migrationclass

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

// Change is one class_resource_definitions row update.
type Change struct {
	ID          uuid.UUID
	ClassID     uuid.UUID
	ResourceKey string
	Updates     map[string]interface{}
}

type row struct {
	ID          uuid.UUID
	ClassID     uuid.UUID
	ResourceKey string
	Category    string
}

// ComputeChanges plans category normalization.
func ComputeChanges(db *gorm.DB) ([]Change, error) {
	var rows []row
	if err := db.Raw(`
SELECT id, class_id, resource_key, category::text FROM class_resource_definitions ORDER BY class_id, resource_key
`).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("load class_resource_definitions: %w", err)
	}
	var out []Change
	for _, r := range rows {
		cat, ok := models.ParseClassResourceCategory(r.Category)
		if !ok {
			continue
		}
		if string(cat) == r.Category {
			continue
		}
		out = append(out, Change{
			ID:          r.ID,
			ClassID:     r.ClassID,
			ResourceKey: r.ResourceKey,
			Updates: map[string]interface{}{
				"category": string(cat),
			},
		})
	}
	return out, nil
}

// ApplyPlanned writes updates.
func ApplyPlanned(db *gorm.DB, changes []Change) error {
	for _, c := range changes {
		if len(c.Updates) == 0 {
			continue
		}
		if err := db.Model(&models.ClassResourceDefinition{}).Where("id = ?", c.ID).Updates(c.Updates).Error; err != nil {
			return fmt.Errorf("update class_resource_definition %s: %w", c.ID, err)
		}
	}
	return nil
}
