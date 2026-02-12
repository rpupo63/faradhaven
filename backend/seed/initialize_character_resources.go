package seed

import (
	"fmt"
	"log"

	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// InitializeCharacterResources ensures existing characters have money and
// generic CharacterResource rows. Legacy class-specific column initialization
// has been replaced by the CharacterResource system (see BackfillCharacterResources).
func InitializeCharacterResources(db *gorm.DB) error {
	log.Println("Initializing resources for existing characters...")

	var characters []models.Character
	if err := db.Find(&characters).Error; err != nil {
		return fmt.Errorf("failed to fetch characters: %w", err)
	}

	for _, char := range characters {
		updated := false

		// Initialize Money if 0
		if char.Money == 0 {
			char.Money = 5000 // 50 gp
			updated = true
		}

		if updated {
			if err := db.Save(&char).Error; err != nil {
				log.Printf("  Error: failed to update character %s: %v", char.ID, err)
			} else {
				log.Printf("  Initialized resources for %s", char.Name)
			}
		}
	}

	log.Println("Character resource initialization completed")
	return nil
}

// BackfillCharacterResources creates CharacterResource rows from existing Character column values
// for characters that were created before the generic resource system was introduced.
// This is idempotent: it skips characters that already have CharacterResource rows.
func BackfillCharacterResources(db *gorm.DB) error {
	log.Println("Backfilling CharacterResource rows from existing character columns...")

	classRepo := database.NewClassRepo(db)
	resourceRepo := database.NewCharacterResourceRepo(db)

	var characters []models.Character
	if err := db.Find(&characters).Error; err != nil {
		return fmt.Errorf("failed to fetch characters: %w", err)
	}

	backfilled := 0
	for _, char := range characters {
		// Check if this character already has generic CharacterResource rows
		existing, err := resourceRepo.FindByCharacterID(char.ID)
		if err == nil && len(existing) > 0 {
			continue // Already has resources, skip
		}

		// Get resource definitions for the class
		defs, err := classRepo.FindResourceDefinitionsByClassID(char.ClassID)
		if err != nil || len(defs) == 0 {
			continue
		}

		// Get level resource values
		levelResources, err := classRepo.FindLevelResourcesByClassAndLevel(char.ClassID, char.Level)
		if err != nil {
			continue
		}
		levelResourceMap := make(map[string]int)
		for _, lr := range levelResources {
			levelResourceMap[lr.ResourceKey] = lr.Value
		}

		// Load legacy column values from character
		legacyValues := map[string]int{
			"max_stability":   char.CurrentStability,  // Current value from Character
			"max_blood_ichor": char.CurrentBloodIchor, // Current value from Character
			"madness_base_dc": char.CurrentMadnessDC,  // Current DC value
		}

		for _, def := range defs {
			if !def.IsTrackable {
				continue
			}

			maxVal := levelResourceMap[def.ResourceKey]
			maxValPtr := &maxVal

			// Use legacy column value if available, otherwise default to max
			currentVal := maxVal
			if legacyVal, ok := legacyValues[def.ResourceKey]; ok && legacyVal > 0 {
				currentVal = legacyVal
			}

			res := &models.CharacterResource{
				CharacterID:        char.ID,
				ResourceKey:        def.ResourceKey,
				ResourceName:       def.DisplayName,
				CurrentValue:       currentVal,
				MaxValue:           maxValPtr,
				RestoreOnShortRest: def.RestoreOnShortRest,
				RestoreOnLongRest:  def.RestoreOnLongRest,
			}
			if err := resourceRepo.Add(res); err != nil {
				log.Printf("  Error: failed to create resource %s for character %s: %v", def.ResourceKey, char.ID, err)
				continue
			}
		}
		backfilled++
	}

	log.Printf("Backfilled CharacterResource rows for %d characters", backfilled)
	return nil
}
