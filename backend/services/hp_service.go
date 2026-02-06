package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// UpdateHP updates the character's current HP by the given delta (positive = heal, negative = damage)
func (s *LevelUpService) UpdateHP(userID uuid.UUID, characterID uuid.UUID, delta int) (*models.Character, error) {
	character, err := s.characterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Apply delta
	character.CurrentHP += delta

	// Clamp to valid range [0, MaxHP + TempHP]
	maxWithTemp := character.MaxHP + character.TempHP
	if character.CurrentHP > maxWithTemp {
		character.CurrentHP = maxWithTemp
	}
	if character.CurrentHP < 0 {
		character.CurrentHP = 0
	}

	if err := s.db.Save(character).Error; err != nil {
		return nil, fmt.Errorf("failed to update HP: %w", err)
	}

	return character, nil
}

// SetTempHP sets the character's temporary HP (replaces, doesn't stack per 5e rules)
func (s *LevelUpService) SetTempHP(userID uuid.UUID, characterID uuid.UUID, tempHP int) (*models.Character, error) {
	character, err := s.characterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Temp HP doesn't stack, only use new value if higher (player choice in 5e)
	if tempHP < 0 {
		tempHP = 0
	}
	character.TempHP = tempHP

	if err := s.db.Save(character).Error; err != nil {
		return nil, fmt.Errorf("failed to set temp HP: %w", err)
	}

	return character, nil
}

// UseHitDiceResult contains the result of using hit dice
type UseHitDiceResult struct {
	Character   *models.Character `json:"character"`
	HPHealed    int               `json:"hp_healed"`
	DiceUsed    int               `json:"dice_used"`
	DiceResults []int             `json:"dice_results"`
}

// UseHitDice spends hit dice during a short rest and heals HP
// rolls is the array of roll results from the frontend (d{HitDie} + ConMod each)
func (s *LevelUpService) UseHitDice(userID uuid.UUID, characterID uuid.UUID, rolls []int) (*UseHitDiceResult, error) {
	character, err := s.characterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Calculate available hit dice
	available := character.Level - character.HitDiceUsed
	if len(rolls) > available {
		return nil, fmt.Errorf("not enough hit dice available: have %d, requested %d", available, len(rolls))
	}

	// Sum the healing
	totalHealing := 0
	for _, roll := range rolls {
		if roll < 1 {
			roll = 1 // Minimum 1 HP per die
		}
		totalHealing += roll
	}

	// Apply healing
	character.HitDiceUsed += len(rolls)
	character.CurrentHP += totalHealing
	if character.CurrentHP > character.MaxHP {
		character.CurrentHP = character.MaxHP
	}

	if err := s.db.Save(character).Error; err != nil {
		return nil, fmt.Errorf("failed to use hit dice: %w", err)
	}

	return &UseHitDiceResult{
		Character:   character,
		HPHealed:    totalHealing,
		DiceUsed:    len(rolls),
		DiceResults: rolls,
	}, nil
}

// ShortRest performs a short rest: allows using hit dice (call UseHitDice), restores spell points
func (s *LevelUpService) ShortRest(userID uuid.UUID, characterID uuid.UUID) (*models.Character, error) {
	character, err := s.characterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Restore spell points to max
	classLevel, err := s.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
	if err == nil && classLevel != nil {
		character.CurrentSpellPoints = classLevel.MaxSpellPoints
	}

	// Restore class-specific resources
	s.resourceService.RestoreClassResources(character, "short_rest", classLevel)

	if err := s.db.Save(character).Error; err != nil {
		return nil, fmt.Errorf("failed to short rest: %w", err)
	}

	return character, nil
}

// LongRest performs a long rest: restore HP to max, restore half hit dice (minimum 1), restore spell points
func (s *LevelUpService) LongRest(userID uuid.UUID, characterID uuid.UUID) (*models.Character, error) {
	character, err := s.characterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Restore HP to max
	character.CurrentHP = character.MaxHP
	character.TempHP = 0 // Temp HP goes away after long rest

	// Restore half hit dice (rounded down, minimum 1)
	diceToRestore := character.Level / 2
	if diceToRestore < 1 {
		diceToRestore = 1
	}
	character.HitDiceUsed -= diceToRestore
	if character.HitDiceUsed < 0 {
		character.HitDiceUsed = 0
	}

	// Restore spell points to max
	classLevel, err := s.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
	if err == nil && classLevel != nil {
		character.CurrentSpellPoints = classLevel.MaxSpellPoints
	}

	// Restore class-specific resources
	s.resourceService.RestoreClassResources(character, "long_rest", classLevel)

	if err := s.db.Save(character).Error; err != nil {
		return nil, fmt.Errorf("failed to long rest: %w", err)
	}

	return character, nil
}
