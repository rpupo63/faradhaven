package services

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// LevelDown reverts a character to their previous level using stored history
func (s *LevelUpService) LevelDown(userID uuid.UUID, characterID uuid.UUID) (*LevelUpResponse, error) {
	// 1. Load character
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	// 2. Verify ownership
	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// 3. Check min level
	if character.Level <= 1 {
		return nil, ErrMinLevelReached
	}

	// 4. Find history for current level
	history, err := s.historyRepo.FindByCharacterAndLevel(characterID, character.Level)
	if err != nil {
		return nil, ErrNoHistoryForLevel
	}

	// 5. Parse snapshot
	var snapshot models.CharacterSnapshotData
	if err := json.Unmarshal(history.CharacterSnapshot, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to parse snapshot: %w", err)
	}

	// 6. Execute in transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Restore character state from snapshot
		character.Level = snapshot.Level
		character.Strength = snapshot.Strength
		character.Dexterity = snapshot.Dexterity
		character.Constitution = snapshot.Constitution
		character.Intelligence = snapshot.Intelligence
		character.Wisdom = snapshot.Wisdom
		character.Charisma = snapshot.Charisma
		character.SpellbookIDs = snapshot.SpellbookIDs
		character.CurrentSpellPoints = snapshot.CurrentSpellPoints
		character.SanguineNotoriety = snapshot.Notoriety

		// Restore HP and hit dice from snapshot
		character.CurrentHP = snapshot.CurrentHP
		character.MaxHP = snapshot.MaxHP
		character.TempHP = snapshot.TempHP
		character.HitDiceUsed = snapshot.HitDiceUsed

		// Restore archetype from snapshot (nil if none was selected before this level)
		if snapshot.ArchetypeID != nil {
			archetypeID, err := uuid.Parse(*snapshot.ArchetypeID)
			if err == nil {
				character.ArchetypeID = &archetypeID
			}
		} else {
			character.ArchetypeID = nil
		}

		if err := tx.Save(character).Error; err != nil {
			return err
		}

		// Restore skill proficiencies from snapshot
		if err := tx.Where("character_id = ?", characterID).Delete(&models.CharacterSkill{}).Error; err != nil {
			return err
		}
		for _, skillID := range history.SkillSnapshot {
			cs := &models.CharacterSkill{
				CharacterID: characterID,
				SkillID:     skillID,
				Proficient:  true,
			}
			if err := tx.Create(cs).Error; err != nil {
				return err
			}
		}

		// Delete the history record for the level we're reverting
		return tx.Delete(history).Error
	})

	if err != nil {
		return nil, fmt.Errorf("level-down transaction failed: %w", err)
	}

	// Get the class level for the restored level
	classLevel, _ := s.findClassLevelWithFeatures(character.ClassID, character.Level)

	return &LevelUpResponse{
		Character:  character,
		NewLevel:   character.Level,
		ClassLevel: classLevel,
	}, nil
}
