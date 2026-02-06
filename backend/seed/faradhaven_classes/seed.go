package faradhaven_classes

import (
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed/batch"
	"github.com/rpupo63/unified-personal-site-backend/seed/uuids"
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

// ClassComponentLink represents a class-component join table entry
type ClassComponentLink struct {
	ClassID     string `gorm:"type:uuid;primaryKey"`
	ComponentID string `gorm:"type:uuid;primaryKey"`
}

func (ClassComponentLink) TableName() string {
	return "class_components"
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
	b.WriteString("Starting Equipment: " + strings.Join(cs.AutomaticEquipNames, "; ") + "\n")
	return b.String()
}

// SeedComponents creates Component rows for all spell system components.
// This is run before classes to ensure components exist for linking.
func SeedComponents(tx *gorm.DB) error {
	componentSeeds := SpellSystemComponents()
	components := make([]models.Component, 0, len(componentSeeds))

	for _, comp := range componentSeeds {
		componentID := uuids.ComponentUUID(comp.Name)
		components = append(components, models.Component{
			ID:          componentID,
			Name:        comp.Name,
			Symbol:      comp.Symbol,
			Category:    comp.Category,
			Description: comp.Description,
			Element:     comp.Element,
		})
	}

	// Batch upsert components (ON CONFLICT DO UPDATE)
	if err := batch.UpsertBatchUpdateAll(tx, components, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d components", len(components))

	return nil
}

// SeedFaradhavenClasses creates Class, ClassLevel, Archetype, LevelFeature, and ClassComponent rows.
// Uses batch operations and deterministic UUIDs for efficient reseeding.
func SeedFaradhavenClasses(tx *gorm.DB) error {
	classSeeds := AllClasses()

	// Step 1: Collect all entities into slices
	classes := make([]models.Class, 0, len(classSeeds))
	var allArchetypes []models.Archetype
	var allClassLevels []models.ClassLevel
	var allLevelFeatures []models.LevelFeature
	var allEquipmentChoices []models.ClassStartingEquipmentChoice
	var allEquipmentOptions []models.ClassStartingEquipmentOption
	var allClassComponents []ClassComponentLink
	var allWeaponRequirements []models.ClassWeaponRequirement

	// Build archetype map for level feature references
	archetypeMap := make(map[string]uuid.UUID) // "ClassName:ArchetypeName" -> archetypeID

	for _, cs := range classSeeds {
		classID := uuids.ClassUUID(cs.Name)

		// Map weapon/item names to UUID strings
		weaponIDs := make([]string, len(cs.AutomaticWeaponNames))
		for i, name := range cs.AutomaticWeaponNames {
			weaponIDs[i] = uuids.WeaponUUID(name).String()
		}
		itemIDs := make([]string, len(cs.AutomaticItemNames))
		for i, name := range cs.AutomaticItemNames {
			itemIDs[i] = uuids.ItemUUID(name).String()
		}

		classes = append(classes, models.Class{
			ID:                  classID,
			Name:                cs.Name,
			Description:         cs.Description,
			HitDie:              cs.HitDie,
			PrimaryAbility:      cs.PrimaryAbility,
			PhotoURL:            cs.PhotoURL,
			ArchetypeLevel:      cs.ArchetypeLevel,
			Proficiencies:       cs.Proficiencies,
			SkillFocus:          pq.StringArray(cs.DnDSkillFocus),
			SkillChoice:         pq.StringArray(cs.SkillChoice),
			SkillChoiceCount:    2, // D&D 5e standard
			Tools:               pq.StringArray(cs.Tools),
			SavingThrows:        pq.StringArray(cs.SavingThrows),
			StartingEquip:       pq.StringArray(cs.AutomaticEquipNames),
			StartingWeaponIDs:   pq.StringArray(weaponIDs),
			StartingItemIDs:     pq.StringArray(itemIDs),
			ResourceType:        cs.ResourceType,
			ResourceName:        cs.ResourceName,
			ResourceRestoreType: cs.ResourceRestoreType,
		})

		// Collect archetypes
		for i, as := range cs.Archetypes {
			archetypeID := uuids.ArchetypeUUID(cs.Name, as.Name)
			archetypeMap[cs.Name+":"+as.Name] = archetypeID

			allArchetypes = append(allArchetypes, models.Archetype{
				ID:          archetypeID,
				ClassID:     classID,
				Name:        as.Name,
				Description: as.Description,
				SortOrder:   i,
			})
		}

		// Collect weapon requirement if defined
		if cs.WeaponRequirement != nil {
			wr := cs.WeaponRequirement
			reqID := uuids.WeaponRequirementUUID(cs.Name, wr.ModifierType, wr.SelectionLevel)
			allWeaponRequirements = append(allWeaponRequirements, models.ClassWeaponRequirement{
				ID:                reqID,
				ClassID:           classID,
				SelectionLevel:    wr.SelectionLevel,
				ModifierType:      models.ModifierType(wr.ModifierType),
				Description:       wr.Description,
				AllowedCategories: pq.StringArray(wr.AllowedCategories),
			})
		}

		// Collect equipment choices and options
		for i, choiceSeed := range cs.EquipmentChoices {
			choiceID := uuids.EquipmentChoiceUUID(cs.Name, choiceSeed.Instruction)

			allEquipmentChoices = append(allEquipmentChoices, models.ClassStartingEquipmentChoice{
				ID:          choiceID,
				ClassID:     classID,
				Instruction: choiceSeed.Instruction,
				SortOrder:   i,
			})

			for _, optionSeed := range choiceSeed.Options {
				optionID := uuids.EquipmentOptionUUID(choiceID, optionSeed.Description)

				optWeaponIDs := make([]string, len(optionSeed.WeaponNames))
				for j, name := range optionSeed.WeaponNames {
					optWeaponIDs[j] = uuids.WeaponUUID(name).String()
				}
				optItemIDs := make([]string, len(optionSeed.ItemNames))
				for j, name := range optionSeed.ItemNames {
					optItemIDs[j] = uuids.ItemUUID(name).String()
				}

				allEquipmentOptions = append(allEquipmentOptions, models.ClassStartingEquipmentOption{
					ID:          optionID,
					ChoiceID:    choiceID,
					Description: optionSeed.Description,
					Items:       pq.StringArray(optionSeed.Items),
					WeaponIDs:   pq.StringArray(optWeaponIDs),
					ItemIDs:     pq.StringArray(optItemIDs),
				})
			}
		}

		// Collect class levels (1-20) and their features
		for level := 1; level <= 20; level++ {
			classLevelID := uuids.ClassLevelUUID(cs.Name, level)
			avgHitDie := (cs.HitDie + 1) / 2

			cl := models.ClassLevel{
				ID:                      classLevelID,
				ClassID:                 classID,
				Level:                   level,
				HpGain:                  avgHitDie,
				ProficiencyBonus:        proficiencyByLevel(level),
				MaxSpellPoints:          maxSpellPointsByLevel(level),
				AbilityScoreImprovement: abilityScoreImprovementByLevel(level),
			}

			// Apply structured level progression data
			if cs.LevelProgression != nil {
				if lp, ok := cs.LevelProgression[level]; ok {
					if lp.CantripsKnown != nil {
						cl.CantripsKnown = lp.CantripsKnown
					}
					if lp.SpellsKnown != nil {
						cl.SpellsKnown = lp.SpellsKnown
					}
					if lp.MaxSpellPoints > 0 {
						cl.MaxSpellPoints = lp.MaxSpellPoints
					}
					cl.ExtraAttackCount = lp.ExtraAttackCount
					cl.SneakAttackDice = lp.SneakAttackDice
					cl.RageDamageBonus = lp.RageDamageBonus
					cl.MartialArtsDie = lp.MartialArtsDie
					cl.UnarmoredMovement = lp.UnarmoredMovement
					cl.SuperiorityDice = lp.SuperiorityDice
					cl.SuperiorityDie = lp.SuperiorityDie
					cl.BardicInspiration = lp.BardicInspiration

					// Faradhaven class resources
					cl.ConcurrencyLimit = lp.ConcurrencyLimit
					cl.YieldDie = lp.YieldDie
					cl.TimerDuration = lp.TimerDuration
					cl.SpeedDialSlots = lp.SpeedDialSlots
					cl.MadnessBaseDC = lp.MadnessBaseDC
					cl.FeralBonus = lp.FeralBonus
					cl.EchoSlots = lp.EchoSlots
					cl.MaxStability = lp.MaxStability
					cl.MaxSpellLevel = lp.MaxSpellLevel
					cl.BiteDamageDice = lp.BiteDamageDice
				}
			}

			// Set Features text
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

			allClassLevels = append(allClassLevels, cl)

			// Collect level features
			sortOrder := 0

			// Level 1: include class features
			if level == 1 {
				for _, f := range cs.ClassFeatures {
					name, desc := parseFeature(f)
					featureID := uuids.LevelFeatureUUID(classLevelID, name, sortOrder)
					allLevelFeatures = append(allLevelFeatures, models.LevelFeature{
						ID:           featureID,
						ClassLevelID: classLevelID,
						Name:         name,
						Description:  desc,
						SortOrder:    sortOrder,
					})
					sortOrder++
				}
			}

			// Level-specific shared features (nil ArchetypeID = all archetypes get this)
			if cs.LevelFeatures != nil {
				if f, ok := cs.LevelFeatures[level]; ok && f != "" {
					name, desc := parseFeature(f)
					featureID := uuids.LevelFeatureUUID(classLevelID, name, sortOrder)
					allLevelFeatures = append(allLevelFeatures, models.LevelFeature{
						ID:           featureID,
						ClassLevelID: classLevelID,
						ArchetypeID:  nil, // shared by all archetypes
						Name:         name,
						Description:  desc,
						SortOrder:    sortOrder,
					})
					sortOrder++
				}
			}

			// Archetype-specific features for this level
			for _, as := range cs.Archetypes {
				if f, ok := as.Features[level]; ok && f != "" {
					archetypeID := archetypeMap[cs.Name+":"+as.Name]
					name, desc := parseFeature(f)
					featureID := uuids.LevelFeatureUUID(classLevelID, name, sortOrder)
					allLevelFeatures = append(allLevelFeatures, models.LevelFeature{
						ID:           featureID,
						ClassLevelID: classLevelID,
						ArchetypeID:  &archetypeID,
						Name:         name,
						Description:  desc,
						SortOrder:    sortOrder,
					})
					sortOrder++
				}
			}
		}

		// Collect class-component links
		for _, m := range AllComponentClassMappings() {
			if m.ClassName == cs.Name {
				componentID := uuids.ComponentUUID(m.Component.Name)
				allClassComponents = append(allClassComponents, ClassComponentLink{
					ClassID:     classID.String(),
					ComponentID: componentID.String(),
				})
			}
		}
	}

	// Step 2: Clear child tables (order matters for foreign keys)
	tablesToClear := []string{
		"class_starting_equipment_options",
		"class_starting_equipment_choices",
		"level_features",
		"class_levels",
		"class_components",
		"class_weapon_requirements",
	}
	for _, table := range tablesToClear {
		if err := tx.Exec("DELETE FROM " + table).Error; err != nil {
			log.Printf("  Note: Could not clear %s: %v", table, err)
		}
	}

	// Step 3: Batch upsert classes (ON CONFLICT DO UPDATE)
	if err := batch.UpsertBatchUpdateAll(tx, classes, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d classes", len(classes))

	// Step 4: Batch upsert archetypes (ON CONFLICT DO UPDATE)
	if err := batch.UpsertBatchUpdateAll(tx, allArchetypes, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d archetypes", len(allArchetypes))

	// Step 5: Batch insert class levels
	if err := batch.InsertBatch(tx, allClassLevels, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d class levels", len(allClassLevels))

	// Step 6: Batch insert level features
	if err := batch.InsertBatch(tx, allLevelFeatures, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d level features", len(allLevelFeatures))

	// Step 7: Batch insert equipment choices
	if err := batch.InsertBatch(tx, allEquipmentChoices, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d equipment choices", len(allEquipmentChoices))

	// Step 8: Batch insert equipment options
	if err := batch.InsertBatch(tx, allEquipmentOptions, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d equipment options", len(allEquipmentOptions))

	// Step 9: Batch insert class-component links
	if len(allClassComponents) > 0 {
		if err := tx.CreateInBatches(allClassComponents, batch.DefaultBatchSize).Error; err != nil {
			return err
		}
		log.Printf("Inserted %d class-component links", len(allClassComponents))
	}

	// Step 10: Batch insert weapon requirements
	if len(allWeaponRequirements) > 0 {
		if err := batch.InsertBatch(tx, allWeaponRequirements, batch.DefaultBatchSize); err != nil {
			return err
		}
		log.Printf("Inserted %d weapon requirements", len(allWeaponRequirements))
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
