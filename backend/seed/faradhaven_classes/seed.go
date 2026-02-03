package faradhaven_classes

import (
	"log"
	"strings"

	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// AllClasses returns all Faradhaven class seeds for seeding
func AllClasses() []FaradhavenClassSeed {
	return []FaradhavenClassSeed{
		Mutagen(),
		Ironwright(),
		Sanguinist(),
		PistonBrawler(),
		PowderMage(),
		RiftWeaver(),
		VaporBlade(),
		Lorewright(),
	}
}

// proficiencyByLevel returns proficiency bonus for level 1-20 (D&D-style)
func proficiencyByLevel(level int) int {
	if level < 1 {
		return 2
	}
	switch {
	case level <= 4:
		return 2
	case level <= 8:
		return 3
	case level <= 12:
		return 4
	case level <= 16:
		return 5
	default:
		return 6
	}
}

// maxSpellPointsByLevel scales 50-100 for component-mixing
func maxSpellPointsByLevel(level int) int {
	if level < 1 {
		level = 1
	}
	return 50 + (level * 2)
}

// abilityScoreImprovementByLevel returns ASI points at D&D standard levels (4, 8, 12, 16, 19)
func abilityScoreImprovementByLevel(level int) int {
	switch level {
	case 4, 8, 12, 16, 19:
		return 2 // +2 to one ability or +1 to two
	default:
		return 0
	}
}

// parseFeature splits a feature string like "Name — Description" into name and description
func parseFeature(feature string) (name, description string) {
	// Try em-dash first (—)
	if idx := strings.Index(feature, " — "); idx != -1 {
		return strings.TrimSpace(feature[:idx]), strings.TrimSpace(feature[idx+len(" — "):])
	}
	// Try colon with space
	if idx := strings.Index(feature, ": "); idx != -1 {
		return strings.TrimSpace(feature[:idx]), strings.TrimSpace(feature[idx+2:])
	}
	// No separator found, use whole string as name
	return strings.TrimSpace(feature), ""
}

// createLevelFeatures creates structured LevelFeature records from the seed data
func createLevelFeatures(db *gorm.DB, cl models.ClassLevel, cs FaradhavenClassSeed, level int) error {
	sortOrder := 0

	// For level 1, include class features
	if level == 1 {
		for _, f := range cs.ClassFeatures {
			name, desc := parseFeature(f)
			lf := models.LevelFeature{
				ClassLevelID: cl.ID,
				Name:         name,
				Description:  desc,
				SortOrder:    sortOrder,
			}
			if err := db.Create(&lf).Error; err != nil {
				return err
			}
			sortOrder++
		}
	}

	// Add level-specific features
	if cs.LevelFeatures != nil {
		if f, ok := cs.LevelFeatures[level]; ok && f != "" {
			name, desc := parseFeature(f)
			lf := models.LevelFeature{
				ClassLevelID: cl.ID,
				Name:         name,
				Description:  desc,
				SortOrder:    sortOrder,
			}
			if err := db.Create(&lf).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// buildLevel1Features formats class features, proficiencies, and other metadata for ClassLevel.Features
func buildLevel1Features(cs FaradhavenClassSeed) string {
	var b strings.Builder
	b.WriteString("Archetype: " + cs.Archetype + "\n")
	b.WriteString("Concept: " + cs.Concept + "\n\n")
	b.WriteString("Class Features:\n")
	for _, f := range cs.ClassFeatures {
		b.WriteString("• " + f + "\n")
	}
	b.WriteString("\nD&D Skill Focus: " + strings.Join(cs.DnDSkillFocus, ", ") + "\n")
	b.WriteString("Proficiencies: " + cs.Proficiencies + "\n")
	b.WriteString("Skill Choice: " + strings.Join(cs.SkillChoice, ", ") + "\n")
	b.WriteString("Tools: " + strings.Join(cs.Tools, ", ") + "\n")
	b.WriteString("Saving Throws: " + strings.Join(cs.SavingThrows, ", ") + "\n")
	b.WriteString("Starting Equipment: " + strings.Join(cs.StartingEquip, "; ") + "\n")
	return b.String()
}

// SeedFaradhavenClasses creates Class, ClassLevel (1-20), and ClassComponent rows for all Faradhaven classes.
// Skips classes that already exist by name.
func SeedFaradhavenClasses(db *gorm.DB) error {
	for _, cs := range AllClasses() {
		var c models.Class
		err := db.Where("name = ?", cs.Name).First(&c).Error
		if err == nil {
			log.Printf("Class already exists: %s", cs.Name)
		} else if err == gorm.ErrRecordNotFound {
			c = models.Class{
				Name:             cs.Name,
				HitDie:           cs.HitDie,
				PrimaryAbility:   cs.PrimaryAbility,
				PhotoURL:         cs.PhotoURL,
				Proficiencies:    cs.Proficiencies,
				SkillFocus:       pq.StringArray(cs.DnDSkillFocus),
				SkillChoice:      pq.StringArray(cs.SkillChoice),
				SkillChoiceCount: 2, // D&D 5e standard: choose 2 from SkillChoice
				Tools:            pq.StringArray(cs.Tools),
				SavingThrows:     pq.StringArray(cs.SavingThrows),
				StartingEquip:    pq.StringArray(cs.StartingEquip),
			}
			if err := db.Create(&c).Error; err != nil {
				return err
			}
			log.Printf("Created class: %s", c.Name)
		} else {
			return err
		}

		// ClassLevel 1-20
		for level := 1; level <= 20; level++ {
			var cl models.ClassLevel
			err := db.Where("class_id = ? AND level = ?", c.ID, level).First(&cl).Error
			if err == nil {
				continue
			}
			if err != gorm.ErrRecordNotFound {
				return err
			}

			avgHitDie := (cs.HitDie + 1) / 2
			cl = models.ClassLevel{
				ClassID:                 c.ID,
				Level:                   level,
				HpGain:                  avgHitDie,
				ProficiencyBonus:        proficiencyByLevel(level),
				MaxSpellPoints:          maxSpellPointsByLevel(level),
				AbilityScoreImprovement: abilityScoreImprovementByLevel(level),
			}
			if level == 1 {
				cl.Features = buildLevel1Features(cs)
				if cs.LevelFeatures != nil {
					if f, ok := cs.LevelFeatures[1]; ok && f != "" {
						cl.Features = cl.Features + "\n\nLevel 1 Features:\n• " + f
					}
				}
			} else if cs.LevelFeatures != nil {
				if f, ok := cs.LevelFeatures[level]; ok {
					cl.Features = f
				}
			}
			if err := db.Create(&cl).Error; err != nil {
				return err
			}

			// Create structured LevelFeature records
			if err := createLevelFeatures(db, cl, cs, level); err != nil {
				return err
			}
		}

		// Components and ClassComponent links (from faradhaven_components.go)
		for _, m := range AllComponentClassMappings() {
			if m.ClassName != cs.Name {
				continue
			}
			comp := m.Component
			var component models.Component
			err := db.Where("name = ?", comp.Name).First(&component).Error
			if err == gorm.ErrRecordNotFound {
				component = models.Component{
					Name:        comp.Name,
					Type:        comp.Type,
					Description: comp.Description,
					Element:     comp.Element,
				}
				if err := db.Create(&component).Error; err != nil {
					return err
				}
				log.Printf("  Created component: %s (%s)", comp.Name, comp.Type)
			} else if err != nil {
				return err
			}

			// Check if ClassComponent link already exists
			var classComponent models.ClassComponent
			err = db.Where("class_id = ? AND component_id = ?", c.ID, component.ID).First(&classComponent).Error
			if err == gorm.ErrRecordNotFound {
				classComponent = models.ClassComponent{
					ClassID:     c.ID,
					ComponentID: component.ID,
				}
				if err := db.Create(&classComponent).Error; err != nil {
					return err
				}
				log.Printf("  Linked component %s to class %s", comp.Name, cs.Name)
			} else if err != nil {
				return err
			}
		}
	}
	return nil
}

// SeedFaradhavenClassesIfEmpty runs SeedFaradhavenClasses only when no Class exists (fresh DB).
func SeedFaradhavenClassesIfEmpty(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Class{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return SeedFaradhavenClasses(db)
	}
	return nil
}

// BackfillAbilityScoreImprovements updates existing ClassLevel rows to set ASI for levels 4, 8, 12, 16, 19.
// Run after migrations to ensure existing databases get the new field populated.
func BackfillAbilityScoreImprovements(db *gorm.DB) error {
	result := db.Model(&models.ClassLevel{}).
		Where("level IN ?", []int{4, 8, 12, 16, 19}).
		Update("ability_score_improvement", 2)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		log.Printf("Backfilled ability_score_improvement for %d ClassLevel rows", result.RowsAffected)
	}
	return nil
}
