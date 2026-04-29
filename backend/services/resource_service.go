package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

// ResourceService handles class-specific resource calculations and restoration
type ResourceService struct {
	classRepo    *database.ClassRepo
	resourceRepo *database.CharacterResourceRepo
}

func NewResourceService(classRepo *database.ClassRepo, resourceRepo *database.CharacterResourceRepo) *ResourceService {
	return &ResourceService{
		classRepo:    classRepo,
		resourceRepo: resourceRepo,
	}
}

func abilityModifier(score int) int {
	return (score - 10) / 2
}

func proficiencyBonusForLevel(level int) int {
	if level < 1 {
		level = 1
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

// InitializeCharacterResources creates CharacterResource rows for all trackable resources
// defined by the character's class. Called during character creation.
func (s *ResourceService) InitializeCharacterResources(character *models.Character) error {
	if character == nil {
		return nil
	}
	characterID := character.ID
	classID := character.ClassID
	level := character.Level

	defs, err := s.classRepo.FindResourceDefinitionsByClassID(classID)
	if err != nil {
		return err
	}
	if len(defs) == 0 {
		return nil
	}

	// Get level resource values for initial max
	levelResources, err := s.classRepo.FindLevelResourcesByClassAndLevel(classID, level)
	if err != nil {
		return err
	}
	levelResourceMap := make(map[string]int)
	for _, lr := range levelResources {
		levelResourceMap[lr.ResourceKey] = lr.Value
	}

	for _, def := range defs {
		if !def.IsTrackable {
			continue // Only create CharacterResource for trackable resources (pools, counters)
		}

		maxVal := levelResourceMap[def.ResourceKey]
		maxValPtr := &maxVal

		res := &models.CharacterResource{
			CharacterID:        characterID,
			ResourceKey:        def.ResourceKey,
			ResourceName:       def.DisplayName,
			CurrentValue:       maxVal, // Start at max
			MaxValue:           maxValPtr,
			RestoreOnShortRest: def.RestoreOnShortRest,
			RestoreOnLongRest:  def.RestoreOnLongRest,
		}
		if err := s.resourceRepo.Add(res); err != nil {
			return err
		}
	}
	return nil
}

// EnsureTrackableClassResources creates missing CharacterResource rows for trackable pools/counters.
// Characters created before a class gained resource definitions, or whose init partially failed,
// may have trait rows but no pool rows — without this, spend/gain APIs no-op or error.
// Idempotent: existing rows are not modified.
func (s *ResourceService) EnsureTrackableClassResources(character *models.Character) error {
	if character == nil || s.resourceRepo == nil {
		return nil
	}
	characterID := character.ID
	classID := character.ClassID
	level := character.Level

	defs, err := s.classRepo.FindResourceDefinitionsByClassID(classID)
	if err != nil {
		return err
	}
	if len(defs) == 0 {
		return nil
	}

	levelResources, err := s.classRepo.FindLevelResourcesByClassAndLevel(classID, level)
	if err != nil {
		return err
	}
	levelResourceMap := make(map[string]int)
	for _, lr := range levelResources {
		levelResourceMap[lr.ResourceKey] = lr.Value
	}

	for _, def := range defs {
		if !def.IsTrackable {
			continue
		}
		_, err := s.resourceRepo.FindByCharacterAndKey(characterID, def.ResourceKey)
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		maxVal := levelResourceMap[def.ResourceKey]
		maxValPtr := &maxVal

		res := &models.CharacterResource{
			CharacterID:        characterID,
			ResourceKey:        def.ResourceKey,
			ResourceName:       def.DisplayName,
			CurrentValue:       maxVal,
			MaxValue:           maxValPtr,
			RestoreOnShortRest: def.RestoreOnShortRest,
			RestoreOnLongRest:  def.RestoreOnLongRest,
		}
		if err := s.resourceRepo.Add(res); err != nil {
			// Concurrent ensure or rare duplicate — treat as satisfied
			if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				continue
			}
			return err
		}
	}
	return nil
}

// UpdateCharacterResourcesForLevel updates CharacterResource MaxValue from the new level's
// ClassLevelResource values. Called during level-up.
func (s *ResourceService) UpdateCharacterResourcesForLevel(character *models.Character) error {
	if character == nil {
		return nil
	}
	characterID := character.ID
	classID := character.ClassID
	level := character.Level

	levelResources, err := s.classRepo.FindLevelResourcesByClassAndLevel(classID, level)
	if err != nil {
		return err
	}
	levelResourceMap := make(map[string]int)
	for _, lr := range levelResources {
		levelResourceMap[lr.ResourceKey] = lr.Value
	}

	existingResources, err := s.resourceRepo.FindByCharacterID(characterID)
	if err != nil {
		return err
	}

	for _, res := range existingResources {
		newMax, ok := levelResourceMap[res.ResourceKey]
		if !ok {
			continue
		}
		maxVal := newMax
		res.MaxValue = &maxVal
		// If current exceeds new max, cap it
		if res.CurrentValue > newMax {
			res.CurrentValue = newMax
		}
		if err := s.resourceRepo.Update(res); err != nil {
			return err
		}
	}
	return nil
}

// GetResourceValue returns the current value of a CharacterResource by key, or 0 if not found
func (s *ResourceService) GetResourceValue(characterID uuid.UUID, key string) int {
	res, err := s.resourceRepo.FindByCharacterAndKey(characterID, key)
	if err != nil {
		return 0
	}
	return res.CurrentValue
}

// GetResourceMaxValue returns the max value of a CharacterResource by key, or 0 if not found
func (s *ResourceService) GetResourceMaxValue(characterID uuid.UUID, key string) int {
	res, err := s.resourceRepo.FindByCharacterAndKey(characterID, key)
	if err != nil || res.MaxValue == nil {
		return 0
	}
	return *res.MaxValue
}

// GetAllCharacterResources returns all CharacterResource rows for a character
func (s *ResourceService) GetAllCharacterResources(characterID uuid.UUID) ([]*models.CharacterResource, error) {
	if s.resourceRepo == nil {
		return nil, nil
	}
	return s.resourceRepo.FindByCharacterID(characterID)
}

// DeductResource deducts the specified amount from a character's resource.
func (s *ResourceService) DeductResource(characterID uuid.UUID, resourceKey string, amount int) error {
	return s.resourceRepo.UpdateCurrentValue(characterID, resourceKey, -amount)
}

// RestoreResourceToMax sets a trackable pool to its current max (used to keep spell_points in sync with rests).
func (s *ResourceService) RestoreResourceToMax(characterID uuid.UUID, resourceKey string) error {
	return s.resourceRepo.RestoreResourceToMax(characterID, resourceKey)
}

// RestoreClassResources handles resource restoration during rests.
// Generic pool restoration (stability, blood ichor, etc.) is handled by CharacterResource.ProcessRestoration.
// Special-case Character column logic is kept here for mechanics that require custom behavior.
func (s *ResourceService) RestoreClassResources(character *models.Character, restType string, classLevel *models.ClassLevel) {
	// Mutagen: Madness DC and cast count reset on short or long rest
	if restType == "short_rest" || restType == "long_rest" {
		character.MadnessCastCount = 0
		character.CurrentMadnessDC = 10 // Reset to base DC
	}

	// Lorewright: Trauma recovery - 1 point per long rest (full recovery at level 13+)
	if restType == "long_rest" && character.Trauma > 0 {
		if character.Level >= 13 {
			character.Trauma = 0 // Full recovery at level 13+
		} else {
			character.Trauma--
			if character.Trauma < 0 {
				character.Trauma = 0
			}
		}
	}

	// Powder Mage: Clear last component used on rest
	if restType == "long_rest" {
		character.LastComponentUsed = nil
	}

	// --- Generic CharacterResource restoration (pools like stability, blood ichor, etc.) ---
	if s.resourceRepo != nil {
		isLongRest := restType == "long_rest"
		_ = s.resourceRepo.ProcessRestoration(character.ID, isLongRest)
	}
}
