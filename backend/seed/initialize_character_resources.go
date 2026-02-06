package seed

import (
	"fmt"
	"log"

	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"gorm.io/gorm"
)

// InitializeCharacterResources populates the new class-specific resource fields
// for existing characters that were created before these fields existed.
func InitializeCharacterResources(db *gorm.DB) error {
	log.Println("Initializing class-specific resources for existing characters...")

	var characters []models.Character
	if err := db.Find(&characters).Error; err != nil {
		return fmt.Errorf("failed to fetch characters: %w", err)
	}

	resourceService := services.NewResourceService()

	for _, char := range characters {
		// Load the character's class and current level data
		var class models.Class
		if err := db.First(&class, "id = ?", char.ClassID).Error; err != nil {
			log.Printf("  Warning: could not find class %s for character %s", char.ClassID, char.ID)
			continue
		}

		var classLevel models.ClassLevel
		if err := db.First(&classLevel, "class_id = ? AND level = ?", char.ClassID, char.Level).Error; err != nil {
			log.Printf("  Warning: could not find class level %d for character %s", char.Level, char.ID)
			continue
		}

		updated := false

		// 1. Piston Brawler: Initialize Stability
		if class.ResourceType == "stability" && char.CurrentStability == 0 {
			char.CurrentStability = classLevel.MaxStability
			updated = true
		}

		// 2. Sanguinist: Initialize Blood Ichor
		if class.ResourceType == "blood_ichor" && char.CurrentBloodIchor == 0 {
			char.CurrentBloodIchor = resourceService.ComputeMaxBloodIchor(&char)
			updated = true
		}

		// 3. Mutagen: Reset Madness DC (Cast Count)
		if class.ResourceType == "madness" && char.MadnessCastCount != 0 {
			// This is usually 0 already for new fields, but good for completeness
			char.MadnessCastCount = 0
			updated = true
		}

		// 4. Lorewright: Reset Echo Slots
		if class.ResourceType == "echo_slots" && char.EchoSlotsUsed != 0 {
			char.EchoSlotsUsed = 0
			updated = true
		}

		// 5. Initialize Money if 0
		if char.Money == 0 {
			char.Money = 5000 // 50 gp
			updated = true
		}

		if updated {
			if err := db.Save(&char).Error; err != nil {
				log.Printf("  Error: failed to update character %s: %v", char.ID, err)
			} else {
				log.Printf("  Initialized resources for %s (%s)", char.Name, class.Name)
			}
		}
	}

	log.Println("Character resource initialization completed")
	return nil
}
