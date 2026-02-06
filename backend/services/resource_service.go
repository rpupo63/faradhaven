package services

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// ResourceService handles class-specific resource calculations and restoration
type ResourceService struct{}

func NewResourceService() *ResourceService {
	return &ResourceService{}
}

// ComputeMaxBloodIchor calculates max ichor: Level + CHA (+ CON at 11+)
func (s *ResourceService) ComputeMaxBloodIchor(character *models.Character) int {
	chaMod := (character.Charisma - 10) / 2
	if character.Charisma < 10 && (character.Charisma-10)%2 != 0 {
		chaMod--
	}
	
	maxIchor := character.Level + chaMod
	
	// Add CON at level 11+
	if character.Level >= 11 {
		conMod := (character.Constitution - 10) / 2
		if character.Constitution < 10 && (character.Constitution-10)%2 != 0 {
			conMod--
		}
		maxIchor += conMod
	}
	
	return maxIchor
}

// RestoreClassResources handles resource restoration during rests
// classLevel is optional - if provided, allows restoring Stability to max
func (s *ResourceService) RestoreClassResources(character *models.Character, restType string, classLevel *models.ClassLevel) {
	// Mutagen: Madness DC resets on short or long rest
	if restType == "short_rest" || restType == "long_rest" {
		character.MadnessCastCount = 0
	}

	// Sanguinist: Ichor resets on long rest
	if restType == "long_rest" {
		maxIchor := s.ComputeMaxBloodIchor(character)
		character.CurrentBloodIchor = maxIchor
	}

	// Piston Brawler: Stability resets on long rest (normally via weapon action, but long rest fully restores)
	if restType == "long_rest" && classLevel != nil && classLevel.MaxStability > 0 {
		character.CurrentStability = classLevel.MaxStability
	}

	// Lorewright: Echo slots used resets on long rest (slots are still filled, but usage counter resets)
	// Note: Echo slots themselves aren't consumed, they're replaced - so no restoration needed
}
