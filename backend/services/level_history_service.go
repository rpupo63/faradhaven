package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// GetLevelHistory returns all level-up history for a character
func (s *LevelUpService) GetLevelHistory(userID uuid.UUID, characterID uuid.UUID) ([]*models.LevelUpHistory, error) {
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	return s.historyRepo.FindByCharacterID(characterID)
}

// GetLevelUpPreview returns what will be available at the next level
func (s *LevelUpService) GetLevelUpPreview(userID uuid.UUID, characterID uuid.UUID) (*LevelUpPreview, error) {
	character, err := s.CharacterRepo.FindByIDWithSkills(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	if character.Level >= 20 {
		return nil, ErrMaxLevelReached
	}

	nextLevel := character.Level + 1
	classLevel, err := s.findClassLevelWithFeatures(character.ClassID, nextLevel)
	if err != nil {
		return nil, fmt.Errorf("class level data not found: %w", err)
	}

	// Calculate new spells allowed
	currentClassLevel, _ := s.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
	newSpellsAllowed := 0
	if classLevel.SpellsKnown != nil && currentClassLevel != nil && currentClassLevel.SpellsKnown != nil {
		newSpellsAllowed = *classLevel.SpellsKnown - *currentClassLevel.SpellsKnown
	}

	// Calculate HP preview
	conMod := (character.Constitution - 10) / 2
	avgHitDie := (character.Class.HitDie / 2) + 1
	hpGainAvg := avgHitDie + conMod
	if hpGainAvg < 1 {
		hpGainAvg = 1
	}

	// Check if archetype selection is required at the next level
	requiresArchetype := false
	var availableArchetypes []*models.Archetype
	if character.Class.ArchetypeLevel != nil && nextLevel == *character.Class.ArchetypeLevel && character.ArchetypeID == nil {
		requiresArchetype = true
		availableArchetypes, _ = s.archetypeRepo.FindByClassID(character.ClassID)
	}

	// Check if weapon selection is required at the next level
	requiresWeaponSelection := false
	var weaponSelectionInfo *WeaponSelectionInfo
	weaponReq, err := s.classRepo.FindWeaponRequirementByClassAndLevel(character.ClassID, nextLevel)
	if err == nil && weaponReq != nil {
		// Check if character already has a primary weapon with this modifier
		// For now, we assume this is a new requirement at this level
		requiresWeaponSelection = true

		// Get eligible weapons based on categories
		eligibleWeapons, _ := s.weaponRepo.FindByCategories([]string(weaponReq.AllowedCategories))

		weaponSelectionInfo = &WeaponSelectionInfo{
			Description:     weaponReq.Description,
			ModifierType:    string(weaponReq.ModifierType),
			EligibleWeapons: eligibleWeapons,
		}
	}

	// Filter level features based on character's archetype (if they have one)
	if character.ArchetypeID != nil {
		classLevel.LevelFeatures = filterFeaturesByArchetype(classLevel.LevelFeatures, character.ArchetypeID)
	}

	return &LevelUpPreview{
		CurrentLevel:            character.Level,
		NextLevel:               nextLevel,
		ClassLevel:              classLevel,
		ASIPointsAvailable:      classLevel.AbilityScoreImprovement,
		NewSpellsAllowed:        newSpellsAllowed,
		HitDie:                  character.Class.HitDie,
		ConMod:                  conMod,
		HPGainAverage:           hpGainAvg,
		CurrentMaxHP:            character.MaxHP,
		RequiresArchetypeChoice: requiresArchetype,
		AvailableArchetypes:     availableArchetypes,
		RequiresWeaponSelection: requiresWeaponSelection,
		WeaponSelectionInfo:     weaponSelectionInfo,
	}, nil
}
