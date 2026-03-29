// Package migrationweapon normalizes weapon_damages.damage_type and .damage_category
// to canonical enum strings (models.DamageType, models.WeaponDamageCategory).
//
// Run order: after migrate_spells (unrelated but listed together in cmd docs). Requires DATABASE_URL.
// Optional: pg_dump weapon_damages before -apply.
package migrationweapon

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// Change is one weapon_damages row worth of column updates.
type Change struct {
	ID       uuid.UUID
	WeaponID uuid.UUID
	Updates  map[string]interface{}
}

type row struct {
	ID             uuid.UUID
	WeaponID       uuid.UUID
	DamageDice     string
	DamageType     string
	DamageCategory string
}

// ComputeChanges plans UPDATEs without writing.
func ComputeChanges(db *gorm.DB) ([]Change, error) {
	var rows []row
	if err := db.Raw(`
SELECT id, weapon_id, damage_dice, damage_type::text, damage_category::text
FROM weapon_damages
ORDER BY id
`).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("load weapon_damages: %w", err)
	}

	var out []Change
	for _, r := range rows {
		m := map[string]interface{}{}
		if dt, ok := models.ParseDamageType(r.DamageType); ok && string(dt) != r.DamageType {
			m["damage_type"] = string(dt)
		}
		if dc, ok := models.ParseWeaponDamageCategory(r.DamageCategory); ok && string(dc) != r.DamageCategory {
			m["damage_category"] = string(dc)
		}
		if len(m) == 0 {
			continue
		}
		out = append(out, Change{ID: r.ID, WeaponID: r.WeaponID, Updates: m})
	}
	return out, nil
}

// ApplyPlanned runs UPDATEs for each planned change.
func ApplyPlanned(db *gorm.DB, changes []Change) error {
	for _, c := range changes {
		if len(c.Updates) == 0 {
			continue
		}
		if err := db.Model(&models.WeaponDamage{}).Where("id = ?", c.ID).Updates(c.Updates).Error; err != nil {
			return fmt.Errorf("update weapon_damage %s: %w", c.ID, err)
		}
	}
	return nil
}
