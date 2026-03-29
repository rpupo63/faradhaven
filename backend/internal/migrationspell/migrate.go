// Package migrationspell migrates legacy spells table rows to match models.Spell:
// - mechanics: SpellType, integer range (feet), ValidateSpellDuration, AI recommended type/range/duration
// - damage: damage_dice_count + damage_die_size (from legacy text columns), save/damage-type normalization
// - drops legacy text columns: damage_dice, suggested_damage_dice, ai_recommended_damage_dice (when present)
// - ScanLegacyDamageParseIssues (report.go) lists rows that fail legacy parsers (use before -apply).
//
// Run via cmd/migrate_spells (full pipeline) with MechanicsOnly / DamageOnly flags.
package migrationspell

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// Options controls which phases run and whether writes are applied.
type Options struct {
	Apply         bool
	MechanicsOnly bool // only type/range/duration + AI mechanics columns
	DamageOnly    bool // only dice + save/damage normalization + legacy column drop
}

// SpellChange is one spell row worth of updates (may combine mechanics + damage).
type SpellChange struct {
	ID      uuid.UUID
	Name    string
	Updates map[string]interface{}
}

// ComputeChanges builds planned updates without writing. Use for dry-run output.
func ComputeChanges(db *gorm.DB, opts Options) ([]SpellChange, error) {
	if opts.MechanicsOnly && opts.DamageOnly {
		return nil, fmt.Errorf("cannot set both MechanicsOnly and DamageOnly")
	}

	hasLegacyDamageDice, err := columnExists(db, "spells", "damage_dice")
	if err != nil {
		return nil, err
	}

	mechRows, err := loadMechanicsRows(db)
	if err != nil {
		return nil, fmt.Errorf("load mechanics rows: %w", err)
	}

	var dmgByID map[uuid.UUID]legacyDamageRow
	if !opts.MechanicsOnly && hasLegacyDamageDice {
		dmgRows, err := loadLegacyDamageRows(db)
		if err != nil {
			return nil, fmt.Errorf("load damage rows: %w", err)
		}
		dmgByID = make(map[uuid.UUID]legacyDamageRow, len(dmgRows))
		for _, r := range dmgRows {
			dmgByID[r.ID] = r
		}
	}

	var out []SpellChange
	for _, mr := range mechRows {
		var mech plannedMechanics
		if !opts.DamageOnly {
			mech = planMechanics(mr)
		}
		var dmg rowDamageUpdates
		if !opts.MechanicsOnly && hasLegacyDamageDice {
			if dr, ok := dmgByID[mr.ID]; ok {
				dmg = planDamage(dr)
			}
		}
		m := mergeUpdates(mech, dmg)
		if len(m) == 0 {
			continue
		}
		out = append(out, SpellChange{ID: mr.ID, Name: mr.Name, Updates: m})
	}
	return out, nil
}

// ApplyPlanned writes planned updates and drops legacy text columns. Call after ComputeChanges when opts.Apply is true.
func ApplyPlanned(db *gorm.DB, opts Options, changes []SpellChange) error {
	hasLegacyDamageDice, err := columnExists(db, "spells", "damage_dice")
	if err != nil {
		return err
	}

	if !opts.MechanicsOnly {
		if err := EnsureDamageDiceColumns(db); err != nil {
			return fmt.Errorf("ensure schema: %w", err)
		}
	}

	for _, c := range changes {
		m := c.Updates
		m["updated_at"] = time.Now()
		if err := db.Model(&models.Spell{}).Where("id = ?", c.ID).Updates(m).Error; err != nil {
			return fmt.Errorf("update spell %s: %w", c.ID, err)
		}
	}

	if !opts.MechanicsOnly && hasLegacyDamageDice {
		if err := dropLegacyDamageDiceColumns(db); err != nil {
			return fmt.Errorf("drop legacy columns: %w", err)
		}
	}
	return nil
}

// Run computes changes and optionally applies (convenience for wrapper CLIs).
func Run(db *gorm.DB, opts Options) (int, error) {
	changes, err := ComputeChanges(db, opts)
	if err != nil {
		return 0, err
	}
	if !opts.Apply {
		return len(changes), nil
	}
	if err := ApplyPlanned(db, opts, changes); err != nil {
		return 0, err
	}
	return len(changes), nil
}

// CountSpells returns the number of rows in spells.
func CountSpells(db *gorm.DB) (int64, error) {
	var n int64
	err := db.Raw(`SELECT COUNT(*) FROM spells`).Scan(&n).Error
	return n, err
}

func columnExists(db *gorm.DB, table, col string) (bool, error) {
	var n int64
	err := db.Raw(`
SELECT COUNT(*) FROM information_schema.columns
WHERE table_schema = current_schema() AND table_name = ? AND column_name = ?
`, table, col).Scan(&n).Error
	return n > 0, err
}

// EnsureDamageDiceColumns adds integer dice columns if missing (idempotent).
func EnsureDamageDiceColumns(db *gorm.DB) error {
	stmts := []string{
		`ALTER TABLE spells ADD COLUMN IF NOT EXISTS damage_dice_count INTEGER`,
		`ALTER TABLE spells ADD COLUMN IF NOT EXISTS damage_die_size INTEGER`,
		`ALTER TABLE spells ADD COLUMN IF NOT EXISTS suggested_damage_dice_count INTEGER`,
		`ALTER TABLE spells ADD COLUMN IF NOT EXISTS suggested_damage_die_size INTEGER`,
		`ALTER TABLE spells ADD COLUMN IF NOT EXISTS ai_recommended_damage_dice_count INTEGER`,
		`ALTER TABLE spells ADD COLUMN IF NOT EXISTS ai_recommended_damage_die_size INTEGER`,
	}
	for _, s := range stmts {
		if err := db.Exec(s).Error; err != nil {
			return err
		}
	}
	return nil
}

func dropLegacyDamageDiceColumns(db *gorm.DB) error {
	for _, col := range []string{"damage_dice", "suggested_damage_dice", "ai_recommended_damage_dice"} {
		if err := db.Exec(fmt.Sprintf(`ALTER TABLE spells DROP COLUMN IF EXISTS %s`, col)).Error; err != nil {
			return err
		}
	}
	return nil
}

// --- Mechanics (from migrate_spell_mechanics) ---

type mechanicsRow struct {
	ID                    uuid.UUID
	Name                  string
	Type                  string
	RangeText             sql.NullString
	Duration              sql.NullString
	AIRecommendedType     sql.NullString
	AIRecommendedRange    sql.NullString
	AIRecommendedDuration sql.NullString
}

func loadMechanicsRows(db *gorm.DB) ([]mechanicsRow, error) {
	const q = `
SELECT s.id, s.name,
       COALESCE(s.type::text, '') AS type,
       CASE WHEN s."range" IS NULL THEN NULL ELSE s."range"::text END AS range_text,
       s.duration,
       s.ai_recommended_type,
       CASE WHEN s.ai_recommended_range IS NULL THEN NULL ELSE s.ai_recommended_range::text END AS ai_recommended_range,
       s.ai_recommended_duration
FROM spells s
ORDER BY s.name
`
	raw, err := db.Raw(q).Rows()
	if err != nil {
		return nil, err
	}
	defer raw.Close()

	var out []mechanicsRow
	for raw.Next() {
		var r mechanicsRow
		if err := raw.Scan(&r.ID, &r.Name, &r.Type, &r.RangeText, &r.Duration, &r.AIRecommendedType, &r.AIRecommendedRange, &r.AIRecommendedDuration); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, raw.Err()
}

type plannedMechanics struct {
	Type                       *string
	Range                      *int
	RangeClear                 bool
	Duration                   *string
	DurationClear              bool
	AIRecommendedType          *string
	AIRecommendedTypeClear     bool
	AIRecommendedRange         *int
	AIRecommendedRangeClear    bool
	AIRecommendedDuration      *string
	AIRecommendedDurationClear bool
}

func (p plannedMechanics) empty() bool {
	return p.Type == nil && !p.RangeClear && p.Range == nil && !p.DurationClear && p.Duration == nil &&
		p.AIRecommendedType == nil && !p.AIRecommendedTypeClear && !p.AIRecommendedRangeClear && p.AIRecommendedRange == nil &&
		!p.AIRecommendedDurationClear && p.AIRecommendedDuration == nil
}

func planMechanics(r mechanicsRow) plannedMechanics {
	var p plannedMechanics

	wantType := models.SpellTypeUtility
	if st, ok := models.ParseSpellTypeRecommendation(r.Type); ok {
		wantType = st
	}
	if string(wantType) != strings.TrimSpace(r.Type) {
		s := string(wantType)
		p.Type = &s
	}

	if r.RangeText.Valid {
		raw := strings.TrimSpace(r.RangeText.String)
		if raw == "" {
			p.RangeClear = true
		} else if n, ok := models.LegacyRangeTextToFeet(raw); ok {
			if strings.TrimSpace(r.RangeText.String) != strconv.Itoa(n) {
				p.Range = &n
			}
		} else {
			p.RangeClear = true
		}
	}

	if r.Duration.Valid {
		d := strings.TrimSpace(r.Duration.String)
		if d == "" {
			p.DurationClear = true
		} else {
			fixed := d
			if models.ValidateSpellDuration(fixed) != nil {
				fixed = models.NormalizeSpellDurationCandidate(d)
			}
			if models.ValidateSpellDuration(fixed) != nil {
				p.DurationClear = true
			} else if fixed != r.Duration.String {
				p.Duration = &fixed
			}
		}
	}

	if r.AIRecommendedType.Valid {
		t := strings.TrimSpace(r.AIRecommendedType.String)
		if t == "" {
			p.AIRecommendedTypeClear = true
		} else if st, ok := models.ParseSpellTypeRecommendation(t); ok {
			s := string(st)
			if s != t {
				p.AIRecommendedType = &s
			}
		} else {
			p.AIRecommendedTypeClear = true
		}
	}

	if r.AIRecommendedRange.Valid {
		raw := strings.TrimSpace(r.AIRecommendedRange.String)
		if raw == "" {
			p.AIRecommendedRangeClear = true
		} else {
			var n int
			ok := false
			if v, err := strconv.Atoi(raw); err == nil && strconv.Itoa(v) == raw {
				n, ok = v, true
			} else if v, ok2 := models.LegacyRangeTextToFeet(raw); ok2 {
				n, ok = v, true
			}
			if !ok {
				p.AIRecommendedRangeClear = true
			} else if strings.TrimSpace(r.AIRecommendedRange.String) != strconv.Itoa(n) {
				p.AIRecommendedRange = &n
			}
		}
	}

	if r.AIRecommendedDuration.Valid {
		d := strings.TrimSpace(r.AIRecommendedDuration.String)
		if d == "" {
			p.AIRecommendedDurationClear = true
		} else {
			fixed := d
			if models.ValidateSpellDuration(fixed) != nil {
				fixed = models.NormalizeSpellDurationCandidate(d)
			}
			if models.ValidateSpellDuration(fixed) != nil {
				p.AIRecommendedDurationClear = true
			} else if fixed != r.AIRecommendedDuration.String {
				p.AIRecommendedDuration = &fixed
			}
		}
	}

	return p
}

// --- Damage (from migrate_spell_damage_fields) ---

type legacyDamageRow struct {
	ID                      uuid.UUID
	Name                    string
	DamageDice              sql.NullString
	SuggestedDamageDice     sql.NullString
	AIRecommendedDamageDice sql.NullString
	SaveAttr                sql.NullString
	AIRecommendedSaveAttr   sql.NullString
	AIRecommendedDamageType sql.NullString
}

func loadLegacyDamageRows(db *gorm.DB) ([]legacyDamageRow, error) {
	const q = `
SELECT id, name,
       damage_dice, suggested_damage_dice, ai_recommended_damage_dice,
       save_attr, ai_recommended_save_attr, ai_recommended_damage_type
FROM spells
ORDER BY name
`
	raw, err := db.Raw(q).Rows()
	if err != nil {
		return nil, err
	}
	defer raw.Close()

	var out []legacyDamageRow
	for raw.Next() {
		var r legacyDamageRow
		if err := raw.Scan(&r.ID, &r.Name,
			&r.DamageDice, &r.SuggestedDamageDice, &r.AIRecommendedDamageDice,
			&r.SaveAttr, &r.AIRecommendedSaveAttr, &r.AIRecommendedDamageType); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, raw.Err()
}

type rowDamageUpdates struct {
	DamageDiceCount, DamageDieSize                             *int
	SuggestedDamageDiceCount, SuggestedDamageDieSize           *int
	AIRecommendedDamageDiceCount, AIRecommendedDamageDieSize   *int
	SaveAttr, AIRecommendedSaveAttr                            *string
	AIRecommendedDamageType                                    *string
}

func (u rowDamageUpdates) empty() bool {
	return u.DamageDiceCount == nil && u.DamageDieSize == nil &&
		u.SuggestedDamageDiceCount == nil && u.SuggestedDamageDieSize == nil &&
		u.AIRecommendedDamageDiceCount == nil && u.AIRecommendedDamageDieSize == nil &&
		u.SaveAttr == nil && u.AIRecommendedSaveAttr == nil && u.AIRecommendedDamageType == nil
}

func planDamage(r legacyDamageRow) rowDamageUpdates {
	var u rowDamageUpdates

	if r.DamageDice.Valid {
		if c, sz, ok := models.ParseLegacyDamageDiceString(r.DamageDice.String); ok {
			u.DamageDiceCount = &c
			u.DamageDieSize = &sz
		}
	}
	if r.SuggestedDamageDice.Valid {
		if c, sz, ok := models.ParseLegacyDamageDiceString(r.SuggestedDamageDice.String); ok {
			u.SuggestedDamageDiceCount = &c
			u.SuggestedDamageDieSize = &sz
		}
	}
	if r.AIRecommendedDamageDice.Valid {
		if c, sz, ok := models.ParseLegacyDamageDiceString(r.AIRecommendedDamageDice.String); ok {
			u.AIRecommendedDamageDiceCount = &c
			u.AIRecommendedDamageDieSize = &sz
		}
	}

	if r.SaveAttr.Valid {
		if a, ok := models.ParseSaveAttribute(r.SaveAttr.String); ok {
			s := string(a)
			u.SaveAttr = &s
		}
	}
	if r.AIRecommendedSaveAttr.Valid {
		if a, ok := models.ParseSaveAttribute(r.AIRecommendedSaveAttr.String); ok {
			s := string(a)
			u.AIRecommendedSaveAttr = &s
		}
	}
	if r.AIRecommendedDamageType.Valid {
		if d, ok := models.ParseDamageType(r.AIRecommendedDamageType.String); ok {
			s := string(d)
			u.AIRecommendedDamageType = &s
		}
	}
	return u
}

func mergeUpdates(mech plannedMechanics, dmg rowDamageUpdates) map[string]interface{} {
	m := map[string]interface{}{}
	if !mech.empty() {
		if mech.Type != nil {
			m["type"] = *mech.Type
		}
		if mech.RangeClear {
			m["range"] = nil
		} else if mech.Range != nil {
			m["range"] = *mech.Range
		}
		if mech.DurationClear {
			m["duration"] = nil
		} else if mech.Duration != nil {
			m["duration"] = *mech.Duration
		}
		if mech.AIRecommendedTypeClear {
			m["ai_recommended_type"] = nil
		} else if mech.AIRecommendedType != nil {
			m["ai_recommended_type"] = *mech.AIRecommendedType
		}
		if mech.AIRecommendedRangeClear {
			m["ai_recommended_range"] = nil
		} else if mech.AIRecommendedRange != nil {
			m["ai_recommended_range"] = *mech.AIRecommendedRange
		}
		if mech.AIRecommendedDurationClear {
			m["ai_recommended_duration"] = nil
		} else if mech.AIRecommendedDuration != nil {
			m["ai_recommended_duration"] = *mech.AIRecommendedDuration
		}
	}
	if !dmg.empty() {
		if dmg.DamageDiceCount != nil {
			m["damage_dice_count"] = *dmg.DamageDiceCount
			m["damage_die_size"] = *dmg.DamageDieSize
		}
		if dmg.SuggestedDamageDiceCount != nil {
			m["suggested_damage_dice_count"] = *dmg.SuggestedDamageDiceCount
			m["suggested_damage_die_size"] = *dmg.SuggestedDamageDieSize
		}
		if dmg.AIRecommendedDamageDiceCount != nil {
			m["ai_recommended_damage_dice_count"] = *dmg.AIRecommendedDamageDiceCount
			m["ai_recommended_damage_die_size"] = *dmg.AIRecommendedDamageDieSize
		}
		if dmg.SaveAttr != nil {
			m["save_attr"] = *dmg.SaveAttr
		}
		if dmg.AIRecommendedSaveAttr != nil {
			m["ai_recommended_save_attr"] = *dmg.AIRecommendedSaveAttr
		}
		if dmg.AIRecommendedDamageType != nil {
			m["ai_recommended_damage_type"] = *dmg.AIRecommendedDamageType
		}
	}
	return m
}
