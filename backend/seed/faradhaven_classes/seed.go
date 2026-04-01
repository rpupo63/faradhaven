package faradhaven_classes

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed/batch"
	"github.com/rpupo63/unified-personal-site-backend/seed/uuids"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Syllogist(),
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

// maxSpellPointsByLevel scales ~50–90 for component-mixing (tuned vs D&D full-caster slot budgets).
func maxSpellPointsByLevel(level int) int {
	if level < 1 {
		level = 1
	}
	return 48 + (level * 2)
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
	var allResourceDefs []models.ClassResourceDefinition
	var allLevelResources []models.ClassLevelResource

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
			ID:                classID,
			Name:              cs.Name,
			Description:       cs.Description,
			HitDie:            cs.HitDie,
			PrimaryAbility:    cs.PrimaryAbility,
			PhotoURL:          cs.PhotoURL,
			ArchetypeLevel:    cs.ArchetypeLevel,
			Proficiencies:     cs.Proficiencies,
			SkillFocus:        pq.StringArray(cs.DnDSkillFocus),
			SkillChoice:       pq.StringArray(cs.SkillChoice),
			SkillChoiceCount:  2, // D&D 5e standard
			Tools:             pq.StringArray(cs.Tools),
			SavingThrows:      pq.StringArray(cs.SavingThrows),
			StartingEquip:     pq.StringArray(cs.AutomaticEquipNames),
			StartingWeaponIDs: pq.StringArray(weaponIDs),
			StartingItemIDs:   pq.StringArray(itemIDs),
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
					if lp.MaxSpellPoints != nil {
						cl.MaxSpellPoints = *lp.MaxSpellPoints
					}
					cl.ExtraAttackCount = lp.ExtraAttackCount
					cl.SneakAttackDice = lp.SneakAttackDice
					cl.RageDamageBonus = lp.RageDamageBonus
					cl.MartialArtsDie = lp.MartialArtsDie
					cl.UnarmoredMovement = lp.UnarmoredMovement
					cl.SuperiorityDice = lp.SuperiorityDice
					cl.SuperiorityDie = lp.SuperiorityDie
					cl.BardicInspiration = lp.BardicInspiration
					cl.MaxSpellLevel = lp.MaxSpellLevel

					// Faradhaven class resources → ClassLevelResource rows
					for key, value := range lp.Resources {
						resID := uuids.ClassLevelResourceUUID(classLevelID, key)
						allLevelResources = append(allLevelResources, models.ClassLevelResource{
							ID:           resID,
							ClassLevelID: classLevelID,
							ResourceKey:  key,
							Value:        value,
						})
					}
				}
			}

			allClassLevels = append(allClassLevels, cl)

			// Collect level features
			sortOrder := 0

			// Level-specific shared features (nil ArchetypeID = all archetypes get this)
			if cs.LevelFeatures != nil {
				if features, ok := cs.LevelFeatures[level]; ok {
					for _, f := range features {
						if f.Name == "" {
							continue
						}
						featureID := uuids.LevelFeatureUUID(classLevelID, f.Name, sortOrder)
						lf := models.LevelFeature{
							ID:             featureID,
							ClassLevelID:   classLevelID,
							ArchetypeID:    nil,
							Name:           f.Name,
							Description:    f.Description,
							SortOrder:      sortOrder,
							ActionType:     f.ActionType,
							UsesPerRest:    f.UsesPerRest,
							ResetCondition: f.ResetCondition,
						}
						if len(f.ResourceCosts) > 0 {
							type rcJSON struct {
								Key    string `json:"key"`
								Amount int    `json:"amount"`
							}
							costs := make([]rcJSON, len(f.ResourceCosts))
							for i, rc := range f.ResourceCosts {
								costs[i] = rcJSON{Key: rc.Key, Amount: rc.Amount}
							}
							if b, err := json.Marshal(costs); err == nil {
								lf.ResourceCosts = datatypes.JSON(b)
							}
						}
						if len(f.ResourceGains) > 0 {
							type rgJSON struct {
								Key    string `json:"key"`
								Amount int    `json:"amount"`
							}
							gains := make([]rgJSON, len(f.ResourceGains))
							for i, rg := range f.ResourceGains {
								gains[i] = rgJSON{Key: rg.Key, Amount: rg.Amount}
							}
							if b, err := json.Marshal(gains); err == nil {
								lf.ResourceGains = datatypes.JSON(b)
							}
						}
						allLevelFeatures = append(allLevelFeatures, lf)
						sortOrder++
					}
				}
			}

			// Archetype-specific features for this level
			for _, as := range cs.Archetypes {
				if features, ok := as.Features[level]; ok {
					for _, f := range features {
						if f.Name == "" {
							continue
						}
						archetypeID := archetypeMap[cs.Name+":"+as.Name]
						featureID := uuids.LevelFeatureUUID(classLevelID, f.Name, sortOrder)
						lf := models.LevelFeature{
							ID:             featureID,
							ClassLevelID:   classLevelID,
							ArchetypeID:    &archetypeID,
							Name:           f.Name,
							Description:    f.Description,
							SortOrder:      sortOrder,
							ActionType:     f.ActionType,
							UsesPerRest:    f.UsesPerRest,
							ResetCondition: f.ResetCondition,
						}
						if len(f.ResourceCosts) > 0 {
							type rcJSON struct {
								Key    string `json:"key"`
								Amount int    `json:"amount"`
							}
							costs := make([]rcJSON, len(f.ResourceCosts))
							for i, rc := range f.ResourceCosts {
								costs[i] = rcJSON{Key: rc.Key, Amount: rc.Amount}
							}
							if b, err := json.Marshal(costs); err == nil {
								lf.ResourceCosts = datatypes.JSON(b)
							}
						}
						if len(f.ResourceGains) > 0 {
							type rgJSON struct {
								Key    string `json:"key"`
								Amount int    `json:"amount"`
							}
							gains := make([]rgJSON, len(f.ResourceGains))
							for i, rg := range f.ResourceGains {
								gains[i] = rgJSON{Key: rg.Key, Amount: rg.Amount}
							}
							if b, err := json.Marshal(gains); err == nil {
								lf.ResourceGains = datatypes.JSON(b)
							}
						}
						allLevelFeatures = append(allLevelFeatures, lf)
						sortOrder++
					}
				}
			}
		}

		// Collect class-component links
		for _, compName := range cs.ComponentPool {
			componentID := uuids.ComponentUUID(compName)
			allClassComponents = append(allClassComponents, ClassComponentLink{
				ClassID:     classID.String(),
				ComponentID: componentID.String(),
			})
		}

		// Collect class resource definitions
		for _, rd := range cs.ResourceDefinitions {
			cat, ok := models.ParseClassResourceCategory(rd.Category)
			if !ok {
				return fmt.Errorf("class %q resource %q: invalid category %q", cs.Name, rd.Key, rd.Category)
			}
			defID := uuids.ClassResourceDefUUID(cs.Name, rd.Key)
			allResourceDefs = append(allResourceDefs, models.ClassResourceDefinition{
				ID:                 defID,
				ClassID:            classID,
				ResourceKey:        rd.Key,
				DisplayName:        rd.DisplayName,
				Category:           cat,
				Description:        rd.Description,
				DisplayOrder:       rd.DisplayOrder,
				IsTrackable:        rd.IsTrackable,
				RestoreOnShortRest: rd.RestoreOnShortRest,
				RestoreOnLongRest:  rd.RestoreOnLongRest,
			})
		}
	}

	// Step 2: Clear child tables (order matters for foreign keys)
	// (Note: This is also handled globally by Seeder.ClearAllData)
	tablesToClear := []string{
		"class_starting_equipment_options",
		"class_starting_equipment_choices",
		"class_level_resources",
		"class_resource_definitions",
		"level_features",
		"class_levels",
		"class_components",
		"class_weapon_requirements",
	}
	for _, table := range tablesToClear {
		if err := tx.Exec("DELETE FROM " + table).Error; err != nil {
			return fmt.Errorf("could not clear table %s: %w", table, err)
		}
	}

	// Step 3: Batch upsert classes (ON CONFLICT DO UPDATE)
	if err := batch.UpsertBatchUpdateAll(tx, classes, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d classes", len(classes))

	// Step 4: Batch upsert archetypes (characters reference them, so can't clear+insert)
	if err := batch.UpsertBatchUpdateAll(tx, allArchetypes, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Upserted %d archetypes", len(allArchetypes))

	// Step 5: Batch insert class levels
	if err := batch.InsertBatch(tx, allClassLevels, batch.DefaultBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d class levels", len(allClassLevels))

	// Step 5a: Batch insert class resource definitions (depends on classes)
	if len(allResourceDefs) > 0 {
		if err := batch.InsertBatch(tx, allResourceDefs, batch.DefaultBatchSize); err != nil {
			return err
		}
		log.Printf("Inserted %d class resource definitions", len(allResourceDefs))
	}

	// Step 5b: Batch insert class level resources (depends on class levels)
	if len(allLevelResources) > 0 {
		if err := batch.InsertBatch(tx, allLevelResources, batch.DefaultBatchSize); err != nil {
			return err
		}
		log.Printf("Inserted %d class level resources", len(allLevelResources))
	}

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

	// Step 9: Upsert class-component links (ON CONFLICT DO NOTHING for composite PK)
	if len(allClassComponents) > 0 {
		seen := make(map[[2]string]struct{}, len(allClassComponents))
		deduped := make([]ClassComponentLink, 0, len(allClassComponents))
		for _, link := range allClassComponents {
			key := [2]string{link.ClassID, link.ComponentID}
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				deduped = append(deduped, link)
			}
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).
			CreateInBatches(deduped, batch.DefaultBatchSize).Error; err != nil {
			return err
		}
		log.Printf("Upserted %d class-component links", len(deduped))
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
