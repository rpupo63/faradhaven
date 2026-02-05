package faradhaven_classes

import (
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/models"
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
// archetypeMap maps archetype name to its database record (for linking archetype-specific features)
func createLevelFeatures(db *gorm.DB, cl models.ClassLevel, cs FaradhavenClassSeed, level int, archetypeMap map[string]*models.Archetype) error {
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

	// Add level-specific shared features (nil ArchetypeID = all archetypes get this)
	if cs.LevelFeatures != nil {
		if f, ok := cs.LevelFeatures[level]; ok && f != "" {
			name, desc := parseFeature(f)
			lf := models.LevelFeature{
				ClassLevelID: cl.ID,
				ArchetypeID:  nil, // shared by all archetypes
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

	// Add archetype-specific features for this level
	for _, as := range cs.Archetypes {
		if f, ok := as.Features[level]; ok && f != "" {
			archetype := archetypeMap[as.Name]
			if archetype == nil {
				continue
			}
			name, desc := parseFeature(f)
			lf := models.LevelFeature{
				ClassLevelID: cl.ID,
				ArchetypeID:  &archetype.ID,
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
	b.WriteString("Starting Equipment: " + strings.Join(cs.AutomaticEquipNames, "; ") + "\n")
	return b.String()
}

// seedEquipmentChoices persists the starting equipment choices for a class
func seedEquipmentChoices(db *gorm.DB, classID uuid.UUID, className string, choices []EquipmentChoiceSeed) error {
	for i, choiceSeed := range choices {
		choiceID := uuids.EquipmentChoiceUUID(className, choiceSeed.Instruction)

		var choice models.ClassStartingEquipmentChoice
		err := db.Where("id = ?", choiceID).First(&choice).Error
		if err == gorm.ErrRecordNotFound {
			choice = models.ClassStartingEquipmentChoice{
				ID:          choiceID,
				ClassID:     classID,
				Instruction: choiceSeed.Instruction,
				SortOrder:   i,
			}
			if err := db.Create(&choice).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// Seed options for this choice
		for _, optionSeed := range choiceSeed.Options {
			optionID := uuids.EquipmentOptionUUID(choiceID, optionSeed.Description)

			// Map WeaponNames and ItemNames to deterministic UUIDs
			weaponIDs := make([]string, len(optionSeed.WeaponNames))
			for i, name := range optionSeed.WeaponNames {
				weaponIDs[i] = uuids.WeaponUUID(name).String()
			}

			itemIDs := make([]string, len(optionSeed.ItemNames))
			for i, name := range optionSeed.ItemNames {
				itemIDs[i] = uuids.ItemUUID(name).String()
			}

			var option models.ClassStartingEquipmentOption
			err := db.Where("id = ?", optionID).First(&option).Error
			if err == gorm.ErrRecordNotFound {
				option = models.ClassStartingEquipmentOption{
					ID:          optionID,
					ChoiceID:    choiceID,
					Description: optionSeed.Description,
					Items:       pq.StringArray(optionSeed.Items),
					WeaponIDs:   pq.StringArray(weaponIDs),
					ItemIDs:     pq.StringArray(itemIDs),
				}
				if err := db.Create(&option).Error; err != nil {
					return err
				}
			} else if err == nil {
				// Update existing to include IDs
				option.WeaponIDs = pq.StringArray(weaponIDs)
				option.ItemIDs = pq.StringArray(itemIDs)
				option.Items = pq.StringArray(optionSeed.Items)
				if err := db.Save(&option).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
		}
	}
	return nil
}

// SeedFaradhavenClasses creates Class, ClassLevel (1-20), Archetype, and ClassComponent rows for all Faradhaven classes.
// Uses deterministic UUIDs so reseeding doesn't break character references.
func SeedFaradhavenClasses(db *gorm.DB) error {
	for _, cs := range AllClasses() {
		// Generate deterministic UUID for this class
		classID := uuids.ClassUUID(cs.Name)

		var c models.Class
		err := db.Where("id = ?", classID).First(&c).Error
		if err == nil {
			// Map Automatic names to UUID strings
			weaponIDs := make([]string, len(cs.AutomaticWeaponNames))
			for i, name := range cs.AutomaticWeaponNames {
				weaponIDs[i] = uuids.WeaponUUID(name).String()
			}
			itemIDs := make([]string, len(cs.AutomaticItemNames))
			for i, name := range cs.AutomaticItemNames {
				itemIDs[i] = uuids.ItemUUID(name).String()
			}

			// Update existing class with new fields
			c.Name = cs.Name
			c.Description = cs.Description
			c.PhotoURL = cs.PhotoURL
			c.ArchetypeLevel = cs.ArchetypeLevel
			c.StartingEquip = pq.StringArray(cs.AutomaticEquipNames)
			c.StartingWeaponIDs = pq.StringArray(weaponIDs)
			c.StartingItemIDs = pq.StringArray(itemIDs)

			if err := db.Save(&c).Error; err != nil {
				return err
			}
			log.Printf("Updated class: %s", cs.Name)
		} else if err == gorm.ErrRecordNotFound {
			// Map Automatic names to UUID strings
			weaponIDs := make([]string, len(cs.AutomaticWeaponNames))
			for i, name := range cs.AutomaticWeaponNames {
				weaponIDs[i] = uuids.WeaponUUID(name).String()
			}
			itemIDs := make([]string, len(cs.AutomaticItemNames))
			for i, name := range cs.AutomaticItemNames {
				itemIDs[i] = uuids.ItemUUID(name).String()
			}

			c = models.Class{
				ID:               classID,
				Name:             cs.Name,
				Description:      cs.Description,
				HitDie:           cs.HitDie,
				PrimaryAbility:   cs.PrimaryAbility,
				PhotoURL:         cs.PhotoURL,
				ArchetypeLevel:   cs.ArchetypeLevel,
				Proficiencies:    cs.Proficiencies,
				SkillFocus:       pq.StringArray(cs.DnDSkillFocus),
				SkillChoice:      pq.StringArray(cs.SkillChoice),
				SkillChoiceCount: 2, // D&D 5e standard: choose 2 from SkillChoice
				Tools:            pq.StringArray(cs.Tools),
				SavingThrows:     pq.StringArray(cs.SavingThrows),
				StartingEquip:    pq.StringArray(cs.AutomaticEquipNames),
				StartingWeaponIDs: pq.StringArray(weaponIDs),
				StartingItemIDs:   pq.StringArray(itemIDs),
			}
			if err := db.Create(&c).Error; err != nil {
				return err
			}
			log.Printf("Created class: %s (ID: %s)", c.Name, c.ID)
		} else {
			return err
		}

		// Seed Equipment Choices
		if err := seedEquipmentChoices(db, c.ID, c.Name, cs.EquipmentChoices); err != nil {
			return err
		}

		// Create archetypes for this class with deterministic UUIDs
		archetypeMap := make(map[string]*models.Archetype) // name -> archetype for later lookup
		for i, as := range cs.Archetypes {
			archetypeID := uuids.ArchetypeUUID(cs.Name, as.Name)

			var archetype models.Archetype
			err := db.Where("id = ?", archetypeID).First(&archetype).Error
			if err == gorm.ErrRecordNotFound {
				archetype = models.Archetype{
					ID:          archetypeID,
					ClassID:     c.ID,
					Name:        as.Name,
					Description: as.Description,
					SortOrder:   i,
				}
				if err := db.Create(&archetype).Error; err != nil {
					return err
				}
				log.Printf("  Created archetype: %s (ID: %s)", as.Name, archetypeID)
			} else if err != nil {
				return err
			}
			archetypeMap[as.Name] = &archetype
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

			// Apply structured level progression data from seed
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
				}
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

			// Create structured LevelFeature records (including archetype-specific ones)
			if err := createLevelFeatures(db, cl, cs, level, archetypeMap); err != nil {
				return err
			}
		}

		// Components and ClassComponent links (from faradhaven_components.go)
		for _, m := range AllComponentClassMappings() {
			if m.ClassName != cs.Name {
				continue
			}
			comp := m.Component
			// Generate deterministic UUID for this component
			componentID := uuids.ComponentUUID(comp.Name)

			var component models.Component
			err := db.Where("id = ?", componentID).First(&component).Error
			if err == gorm.ErrRecordNotFound {
				component = models.Component{
					ID:          componentID,
					Name:        comp.Name,
					Symbol:      comp.Symbol,
					Category:    comp.Category,
					Description: comp.Description,
					Element:     comp.Element,
				}
				if err := db.Create(&component).Error; err != nil {
					return err
				}
				log.Printf("  Created component: %s [%s] (%s) (ID: %s)", comp.Name, comp.Symbol, comp.Category, componentID)
			} else if err != nil {
				return err
			}

			// Check if ClassComponent link already exists
			var classComponent models.ClassComponent
			err = db.Where("class_id = ? AND component_id = ?", c.ID, componentID).First(&classComponent).Error
			if err == gorm.ErrRecordNotFound {
				classComponent = models.ClassComponent{
					ClassID:     c.ID,
					ComponentID: componentID,
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
