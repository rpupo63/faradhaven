package services

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// UpdateHP updates the character's current HP by the given delta (positive = heal, negative = damage)
func (s *LevelUpService) UpdateHP(userID uuid.UUID, characterID uuid.UUID, delta int, source string) (*models.Character, error) {
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Initialize HP if this is a legacy character
	s.ensureHPInitialized(character)

	// --- Sanguinist Moral Seesaw & Extraction Logic ---
	if character.Class.Name == "The Sanguinist" {
		// Sanguine Extraction
		if source == "Bite" || source == "Siphon" {
			unstableIchor, err := s.componentRepo.FindByName("Unstable Ichor")
			if err != nil {
				// Log error but don't fail the HP update
				fmt.Printf("Error finding Unstable Ichor component: %v\n", err)
			} else {
				var amountToGrant int
				if character.Level >= 1 && character.Level <= 4 {
					amountToGrant = 1
				} else if character.Level >= 5 && character.Level <= 10 {
					amountToGrant = 2
				} else {
					amountToGrant = 3
				}

				err := s.CharacterRepo.UpdateComponentCount(character.ID, unstableIchor.ID, amountToGrant)
				if err != nil {
					fmt.Printf("Error granting Unstable Ichor: %v\n", err)
				}
			}
		}

		// Backfire Check
		isDamageFeature := source == "Bite" || source == "Shadow Mist" || source == "Ichor Lash"
		if delta < 0 && isDamageFeature && (character.SanguineMP-character.SanguineBR) >= 3 {
			// Sanguine Backfire: damage self
			backfireDamage := rand.Intn(8) + 1
			// Call UpdateHP on self, with a specific source to prevent recursion
			_, err := s.UpdateHP(userID, characterID, -backfireDamage, "Sanguine Backfire")
			if err != nil {
				// Log the error but continue with the original damage
				fmt.Printf("Error applying sanguine backfire: %v\n", err)
			}
		}

		// Ravenous Check
		isHealingFeature := source == "Blood Graft"
		if delta > 0 && isHealingFeature && (character.SanguineBR-character.SanguineMP) >= 3 {
			// Ravenous: damage target before healing
			ravenousDamage := rand.Intn(8) + 1
			character.CurrentHP -= ravenousDamage
		}
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
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Initialize HP if this is a legacy character
	s.ensureHPInitialized(character)

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
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Initialize HP if this is a legacy character
	s.ensureHPInitialized(character)

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
	
	// Clamp to MaxHP (healing doesn't affect TempHP usually, unless special features)
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
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Initialize HP if this is a legacy character
	s.ensureHPInitialized(character)

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

// LongRest performs a long rest: restore HP to max, restore all hit dice, restore spell points
func (s *LevelUpService) LongRest(userID uuid.UUID, characterID uuid.UUID) (*models.Character, error) {
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Initialize HP if this is a legacy character
	s.ensureHPInitialized(character)

	// Restore HP to max
	character.CurrentHP = character.MaxHP
	character.TempHP = 0 // Temp HP goes away after long rest

	// Restore all hit dice
	character.HitDiceUsed = 0

	// Restore spell points to max
	classLevel, err := s.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
	if err == nil && classLevel != nil {
		character.CurrentSpellPoints = classLevel.MaxSpellPoints
	}

	// Restore class-specific resources
	s.resourceService.RestoreClassResources(character, "long_rest", classLevel)

	// Sanguinist: All extracted components decay on Long Rest
	if character.Class.Name == "The Sanguinist" {
		if err := s.CharacterRepo.ClearComponentsForCharacter(character.ID); err != nil {
			fmt.Printf("Error clearing Sanguinist components during long rest: %v\n", err)
		}
	}

	// Ironwright: Automated Logistics (Level 15) — baseline components restored on Long Rest
	if character.Class.Name == "The Ironwright" && character.Level >= 15 {
		class, classErr := s.classRepo.FindByIDWithLevels(character.ClassID)
		if classErr == nil && len(class.Components) > 0 {
			// Distribute baseline (= character level) evenly across class component pool
			perComponent := character.Level / len(class.Components)
			if perComponent < 1 {
				perComponent = 1
			}
			for _, comp := range class.Components {
				// Check current count; only top up if below baseline
				var charComp models.CharacterComponent
				err := s.db.Where("character_id = ? AND component_id = ?", character.ID, comp.ID).First(&charComp).Error
				currentCount := 0
				if err == nil {
					currentCount = charComp.Count
				}
				if currentCount < perComponent {
					delta := perComponent - currentCount
					_ = s.CharacterRepo.UpdateComponentCount(character.ID, comp.ID, delta)
				}
			}
		}
	}

	if err := s.db.Save(character).Error; err != nil {
		return nil, fmt.Errorf("failed to long rest: %w", err)
	}

	return character, nil
}

// ensureHPInitialized calculates and sets MaxHP/CurrentHP for characters that don't have them yet
func (s *LevelUpService) ensureHPInitialized(character *models.Character) {
	if character.MaxHP > 0 {
		return
	}

	conMod := models.AbilityModifier(character.Constitution)
	avgHitDie := (character.Class.HitDie + 1) / 2
	baseHP := character.Class.HitDie
	
	totalHP := baseHP + (avgHitDie * (character.Level - 1)) + (conMod * character.Level)
	if character.Level < 1 {
		totalHP = baseHP + conMod
	}
	
	if totalHP < 1 {
		totalHP = 1
	}

	character.MaxHP = totalHP
	if character.CurrentHP == 0 {
		character.CurrentHP = totalHP
	}
}
