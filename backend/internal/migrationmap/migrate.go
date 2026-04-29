// Package migrationmap normalizes map_tokens.token_type to "pc" or "npc".
//
// Run order: independent of other migration CLIs. Requires DATABASE_URL.
package migrationmap

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

// Change is one map_tokens row update.
type Change struct {
	ID      uuid.UUID
	MapID   uuid.UUID
	Updates map[string]interface{}
}

type row struct {
	ID        uuid.UUID
	MapID     uuid.UUID
	TokenType string
}

// ComputeChanges plans token_type normalization.
func ComputeChanges(db *gorm.DB) ([]Change, error) {
	var rows []row
	if err := db.Raw(`
SELECT id, map_id, token_type::text FROM map_tokens ORDER BY id
`).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("load map_tokens: %w", err)
	}
	var out []Change
	for _, r := range rows {
		raw := strings.TrimSpace(strings.ToLower(r.TokenType))
		tt, ok := models.ParseMapTokenType(raw)
		if !ok {
			continue
		}
		if string(tt) == r.TokenType {
			continue
		}
		out = append(out, Change{
			ID:    r.ID,
			MapID: r.MapID,
			Updates: map[string]interface{}{
				"token_type": string(tt),
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
		if err := db.Model(&models.MapToken{}).Where("id = ?", c.ID).Updates(c.Updates).Error; err != nil {
			return fmt.Errorf("update map_token %s: %w", c.ID, err)
		}
	}
	return nil
}
