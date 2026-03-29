// Package migrationcreature normalizes monsters.size/type and corpses.creature_type/creature_size
// to canonical CreatureSize / CreatureType strings.
//
// Run order: after migrate_spells; before or after migrationweapon (independent). Requires DATABASE_URL.
package migrationcreature

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// MonsterChange is a planned update to monsters.
type MonsterChange struct {
	ID      uuid.UUID
	Name    string
	Updates map[string]interface{}
}

// CorpseChange is a planned update to corpses.
type CorpseChange struct {
	ID      uuid.UUID
	Name    string
	Updates map[string]interface{}
}

type monsterRow struct {
	ID   uuid.UUID
	Name string
	Size string
	Type string
}

type corpseRow struct {
	ID             uuid.UUID
	Name           string
	CreatureType   string
	CreatureSize   string
}

// ComputeMonsterChanges plans monster row updates.
func ComputeMonsterChanges(db *gorm.DB) ([]MonsterChange, error) {
	var rows []monsterRow
	if err := db.Raw(`
SELECT id, name, size::text, "type"::text FROM monsters ORDER BY name
`).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("load monsters: %w", err)
	}
	var out []MonsterChange
	for _, r := range rows {
		m := map[string]interface{}{}
		if sz, ok := models.ParseCreatureSize(r.Size); ok && string(sz) != r.Size {
			m["size"] = string(sz)
		}
		if ty, ok := models.ParseCreatureType(r.Type); ok && string(ty) != r.Type {
			m["type"] = string(ty)
		}
		if len(m) == 0 {
			continue
		}
		out = append(out, MonsterChange{ID: r.ID, Name: r.Name, Updates: m})
	}
	return out, nil
}

// ComputeCorpseChanges plans corpse row updates.
func ComputeCorpseChanges(db *gorm.DB) ([]CorpseChange, error) {
	var rows []corpseRow
	if err := db.Raw(`
SELECT id, name, creature_type::text, creature_size::text FROM corpses ORDER BY name
`).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("load corpses: %w", err)
	}
	var out []CorpseChange
	for _, r := range rows {
		m := map[string]interface{}{}
		if ty, ok := models.ParseCreatureType(r.CreatureType); ok && string(ty) != r.CreatureType {
			m["creature_type"] = string(ty)
		}
		if sz, ok := models.ParseCreatureSize(r.CreatureSize); ok && string(sz) != r.CreatureSize {
			m["creature_size"] = string(sz)
		}
		if len(m) == 0 {
			continue
		}
		out = append(out, CorpseChange{ID: r.ID, Name: r.Name, Updates: m})
	}
	return out, nil
}

// ApplyMonsters writes planned monster updates.
func ApplyMonsters(db *gorm.DB, changes []MonsterChange) error {
	for _, c := range changes {
		if len(c.Updates) == 0 {
			continue
		}
		if err := db.Model(&models.Monster{}).Where("id = ?", c.ID).Updates(c.Updates).Error; err != nil {
			return fmt.Errorf("update monster %s: %w", c.ID, err)
		}
	}
	return nil
}

// ApplyCorpses writes planned corpse updates.
func ApplyCorpses(db *gorm.DB, changes []CorpseChange) error {
	for _, c := range changes {
		if len(c.Updates) == 0 {
			continue
		}
		if err := db.Model(&models.Corpse{}).Where("id = ?", c.ID).Updates(c.Updates).Error; err != nil {
			return fmt.Errorf("update corpse %s: %w", c.ID, err)
		}
	}
	return nil
}
