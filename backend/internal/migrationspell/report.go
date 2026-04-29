package migrationspell

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

// LegacyDamageParseIssue records a spell row and field whose legacy text could not be normalized by parsers.
type LegacyDamageParseIssue struct {
	SpellID uuid.UUID
	Name    string
	Field   string
	Raw     string
	Detail  string
}

// ScanLegacyDamageParseIssues inspects legacy spell columns (when present) and lists parse failures.
// Use before migrate_spells -apply to fix or export bad rows. Safe when damage_dice column is already dropped (returns nil, nil).
func ScanLegacyDamageParseIssues(db *gorm.DB) ([]LegacyDamageParseIssue, error) {
	hasLegacy, err := columnExists(db, "spells", "damage_dice")
	if err != nil {
		return nil, err
	}
	if !hasLegacy {
		return nil, nil
	}

	rows, err := loadLegacyDamageRows(db)
	if err != nil {
		return nil, err
	}

	var issues []LegacyDamageParseIssue
	for _, r := range rows {
		if r.DamageDice.Valid && strings.TrimSpace(r.DamageDice.String) != "" {
			if _, _, ok := models.ParseLegacyDamageDiceString(r.DamageDice.String); !ok {
				issues = append(issues, LegacyDamageParseIssue{
					SpellID: r.ID, Name: r.Name, Field: "damage_dice",
					Raw: r.DamageDice.String, Detail: "ParseLegacyDamageDiceString failed",
				})
			}
		}
		if r.SuggestedDamageDice.Valid && strings.TrimSpace(r.SuggestedDamageDice.String) != "" {
			if _, _, ok := models.ParseLegacyDamageDiceString(r.SuggestedDamageDice.String); !ok {
				issues = append(issues, LegacyDamageParseIssue{
					SpellID: r.ID, Name: r.Name, Field: "suggested_damage_dice",
					Raw: r.SuggestedDamageDice.String, Detail: "ParseLegacyDamageDiceString failed",
				})
			}
		}
		if r.AIRecommendedDamageDice.Valid && strings.TrimSpace(r.AIRecommendedDamageDice.String) != "" {
			if _, _, ok := models.ParseLegacyDamageDiceString(r.AIRecommendedDamageDice.String); !ok {
				issues = append(issues, LegacyDamageParseIssue{
					SpellID: r.ID, Name: r.Name, Field: "ai_recommended_damage_dice",
					Raw: r.AIRecommendedDamageDice.String, Detail: "ParseLegacyDamageDiceString failed",
				})
			}
		}
		if r.SaveAttr.Valid && strings.TrimSpace(r.SaveAttr.String) != "" {
			if _, ok := models.ParseSaveAttribute(r.SaveAttr.String); !ok {
				issues = append(issues, LegacyDamageParseIssue{
					SpellID: r.ID, Name: r.Name, Field: "save_attr",
					Raw: r.SaveAttr.String, Detail: "ParseSaveAttribute failed",
				})
			}
		}
		if r.AIRecommendedSaveAttr.Valid && strings.TrimSpace(r.AIRecommendedSaveAttr.String) != "" {
			if _, ok := models.ParseSaveAttribute(r.AIRecommendedSaveAttr.String); !ok {
				issues = append(issues, LegacyDamageParseIssue{
					SpellID: r.ID, Name: r.Name, Field: "ai_recommended_save_attr",
					Raw: r.AIRecommendedSaveAttr.String, Detail: "ParseSaveAttribute failed",
				})
			}
		}
		if r.AIRecommendedDamageType.Valid && strings.TrimSpace(r.AIRecommendedDamageType.String) != "" {
			if _, ok := models.ParseDamageType(r.AIRecommendedDamageType.String); !ok {
				issues = append(issues, LegacyDamageParseIssue{
					SpellID: r.ID, Name: r.Name, Field: "ai_recommended_damage_type",
					Raw: r.AIRecommendedDamageType.String, Detail: "ParseDamageType failed",
				})
			}
		}
	}

	if ok, err := columnExists(db, "spells", "damage_type"); err != nil {
		return nil, err
	} else if ok {
		type dtRow struct {
			ID         uuid.UUID
			Name       string
			DamageType sql.NullString
		}
		var dts []dtRow
		if err := db.Raw(`SELECT id, name, damage_type FROM spells WHERE damage_type IS NOT NULL AND trim(damage_type::text) <> ''`).Scan(&dts).Error; err != nil {
			return nil, err
		}
		for _, dr := range dts {
			if !dr.DamageType.Valid {
				continue
			}
			if _, ok := models.ParseDamageType(dr.DamageType.String); !ok {
				issues = append(issues, LegacyDamageParseIssue{
					SpellID: dr.ID, Name: dr.Name, Field: "damage_type",
					Raw: dr.DamageType.String, Detail: "ParseDamageType failed",
				})
			}
		}
	}

	return issues, nil
}
