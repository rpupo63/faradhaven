package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// ApplyEffectOptions contains options for applying an effect
type ApplyEffectOptions struct {
	Stacks            int        // Number of stacks to apply (default 1)
	MaxStacks         int        // Maximum stacks allowed (default 1)
	DurationRounds    *int       // Duration in rounds (nil = indefinite)
	DurationMinutes   *int       // Duration in minutes
	Source            string     // Description of the source
	SourceCharacterID *uuid.UUID // Who applied this effect
	SourceSpellID     *uuid.UUID // What spell caused this effect
	IsConcentration   bool       // Whether this effect requires concentration
	StackIfExists     bool       // If true and effect exists, add stacks instead of replacing
}

// ExpiredEffect represents an effect that has expired
type ExpiredEffect struct {
	EffectID   uuid.UUID
	EffectName string
	Reason     string // "duration_expired", "concentration_broken", etc.
}

// TickResult represents the result of a duration tick
type TickResult struct {
	Expired       []ExpiredEffect
	RoundsTickerd int
}

// EffectService handles effect management including stacking and concentration
type EffectService struct {
	db                  *gorm.DB
	characterEffectRepo *database.CharacterEffectRepo
	effectRepo          *database.EffectRepo
	CharacterRepo       *database.CharacterRepo
}

// NewEffectService creates a new EffectService
func NewEffectService(
	db *gorm.DB,
	characterEffectRepo *database.CharacterEffectRepo,
	effectRepo *database.EffectRepo,
	CharacterRepo *database.CharacterRepo,
) *EffectService {
	return &EffectService{
		db:                  db,
		characterEffectRepo: characterEffectRepo,
		effectRepo:          effectRepo,
		CharacterRepo:       CharacterRepo,
	}
}

// ApplyEffect adds or stacks an effect on a character
func (s *EffectService) ApplyEffect(characterID, effectID uuid.UUID, opts ApplyEffectOptions) (*models.CharacterEffect, error) {
	// Validate effect exists
	effect, err := s.effectRepo.FindByID(effectID)
	if err != nil {
		return nil, errors.New("effect not found")
	}

	// Set defaults
	if opts.Stacks <= 0 {
		opts.Stacks = 1
	}
	if opts.MaxStacks <= 0 {
		opts.MaxStacks = 1
	}

	// Check if this effect already exists on the character
	existing, err := s.characterEffectRepo.FindByEffectAndCharacter(characterID, effectID)
	if err == nil && existing != nil && opts.StackIfExists {
		// Effect exists and we should stack - add stacks up to max
		existing.Stacks += opts.Stacks
		if existing.Stacks > existing.MaxStacks {
			existing.Stacks = existing.MaxStacks
		}
		// Update source if provided
		if opts.Source != "" {
			existing.Source = opts.Source
		}
		// Refresh duration if provided
		if opts.DurationRounds != nil {
			existing.DurationRounds = opts.DurationRounds
		}
		existing.UpdatedAt = time.Now()

		if err := s.characterEffectRepo.Update(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// Handle concentration - break existing concentration if applying a new one
	if opts.IsConcentration {
		if _, err := s.BreakConcentration(characterID); err != nil {
			return nil, err
		}
	}

	// Create new effect
	charEffect := &models.CharacterEffect{
		CharacterID:       characterID,
		EffectID:          effectID,
		Stacks:            opts.Stacks,
		MaxStacks:         opts.MaxStacks,
		DurationRounds:    opts.DurationRounds,
		DurationMinutes:   opts.DurationMinutes,
		Source:            opts.Source,
		SourceCharacterID: opts.SourceCharacterID,
		SourceSpellID:     opts.SourceSpellID,
		IsConcentration:   opts.IsConcentration,
	}

	if err := s.characterEffectRepo.Add(charEffect); err != nil {
		return nil, err
	}

	// Reload to get the Effect relationship populated
	charEffect.Effect = *effect
	return charEffect, nil
}

// RemoveEffect removes an effect instance, optionally reducing stacks instead
func (s *EffectService) RemoveEffect(effectInstanceID uuid.UUID, removeAllStacks bool, stacksToRemove int) error {
	effect, err := s.characterEffectRepo.FindByID(effectInstanceID)
	if err != nil {
		return errors.New("effect instance not found")
	}

	// If removing all or only one stack, just delete
	if removeAllStacks || effect.Stacks <= 1 || (stacksToRemove >= effect.Stacks) {
		return s.characterEffectRepo.Delete(effectInstanceID)
	}

	// Otherwise reduce stacks
	if stacksToRemove <= 0 {
		stacksToRemove = 1
	}
	effect.Stacks -= stacksToRemove
	if effect.Stacks < 0 {
		effect.Stacks = 0
	}

	if effect.Stacks == 0 {
		return s.characterEffectRepo.Delete(effectInstanceID)
	}

	return s.characterEffectRepo.Update(effect)
}

// ModifyStacks modifies the stack count of an effect
func (s *EffectService) ModifyStacks(effectInstanceID uuid.UUID, delta int) (*models.CharacterEffect, error) {
	effect, err := s.characterEffectRepo.FindByID(effectInstanceID)
	if err != nil {
		return nil, errors.New("effect instance not found")
	}

	effect.Stacks += delta

	// Clamp to valid range
	if effect.Stacks < 0 {
		effect.Stacks = 0
	}
	if effect.Stacks > effect.MaxStacks {
		effect.Stacks = effect.MaxStacks
	}

	// Remove if stacks hit zero
	if effect.Stacks == 0 {
		return nil, s.characterEffectRepo.Delete(effectInstanceID)
	}

	if err := s.characterEffectRepo.Update(effect); err != nil {
		return nil, err
	}

	return effect, nil
}

// TickDuration decrements duration for all round-based effects and removes expired ones
func (s *EffectService) TickDuration(characterID uuid.UUID, rounds int) (*TickResult, error) {
	if rounds <= 0 {
		rounds = 1
	}

	// Decrement durations and get expired effects
	expired, err := s.characterEffectRepo.DecrementDurationRounds(characterID, rounds)
	if err != nil {
		return nil, err
	}

	result := &TickResult{
		Expired:       make([]ExpiredEffect, 0),
		RoundsTickerd: rounds,
	}

	// Build result and delete expired effects
	for _, exp := range expired {
		result.Expired = append(result.Expired, ExpiredEffect{
			EffectID:   exp.EffectID,
			EffectName: exp.Effect.Name,
			Reason:     "duration_expired",
		})
		// Delete the expired effect
		_ = s.characterEffectRepo.Delete(exp.ID)
	}

	return result, nil
}

// BreakConcentration removes all concentration effects from a character
func (s *EffectService) BreakConcentration(characterID uuid.UUID) ([]ExpiredEffect, error) {
	concentrationEffects, err := s.characterEffectRepo.FindConcentrationEffects(characterID)
	if err != nil {
		return nil, err
	}

	expired := make([]ExpiredEffect, 0)
	for _, ce := range concentrationEffects {
		expired = append(expired, ExpiredEffect{
			EffectID:   ce.EffectID,
			EffectName: ce.Effect.Name,
			Reason:     "concentration_broken",
		})
		_ = s.characterEffectRepo.Delete(ce.ID)
	}

	// Clear concentration slots on character
	character, err := s.CharacterRepo.FindByID(characterID)
	if err == nil && character != nil {
		character.ConcentrationSlots = nil
		_ = s.CharacterRepo.Update(character)
	}

	return expired, nil
}

// GetActiveEffects returns all active effects for a character
func (s *EffectService) GetActiveEffects(characterID uuid.UUID) ([]*models.CharacterEffect, error) {
	return s.characterEffectRepo.FindByCharacterID(characterID)
}

// GetActiveEffect returns a specific effect instance
func (s *EffectService) GetActiveEffect(effectInstanceID uuid.UUID) (*models.CharacterEffect, error) {
	return s.characterEffectRepo.FindByID(effectInstanceID)
}

// ClearAllEffects removes all effects from a character
func (s *EffectService) ClearAllEffects(characterID uuid.UUID) error {
	return s.characterEffectRepo.DeleteByCharacterID(characterID)
}

// HasEffect checks if a character has a specific effect active
func (s *EffectService) HasEffect(characterID, effectID uuid.UUID) (bool, *models.CharacterEffect, error) {
	effect, err := s.characterEffectRepo.FindByEffectAndCharacter(characterID, effectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, effect, nil
}

// GetEffectsBySource returns all effects applied by a specific source character
func (s *EffectService) GetEffectsBySource(sourceCharacterID uuid.UUID) ([]*models.CharacterEffect, error) {
	return s.characterEffectRepo.FindBySourceCharacter(sourceCharacterID)
}
