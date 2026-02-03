package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

var (
	ErrMaxLevelReached       = errors.New("character is already at maximum level (20)")
	ErrMinLevelReached       = errors.New("character is already at minimum level (1)")
	ErrNoHistoryForLevel     = errors.New("no level-up history found for this level")
	ErrInsufficientASIPoints = errors.New("ASI allocation exceeds available points")
	ErrUnauthorized          = errors.New("unauthorized: character belongs to different user")
)

type LevelUpService struct {
	db            *gorm.DB
	characterRepo *database.CharacterRepo
	classRepo     *database.ClassRepo
	historyRepo   *database.LevelUpHistoryRepo
}

func NewLevelUpService(
	db *gorm.DB,
	characterRepo *database.CharacterRepo,
	classRepo *database.ClassRepo,
	historyRepo *database.LevelUpHistoryRepo,
) *LevelUpService {
	return &LevelUpService{
		db:            db,
		characterRepo: characterRepo,
		classRepo:     classRepo,
		historyRepo:   historyRepo,
	}
}

// LevelUpRequest contains the choices made during level-up
type LevelUpRequest struct {
	CharacterID     uuid.UUID      `json:"character_id"`
	SkillSelections []string       `json:"skill_selections,omitempty"` // new skills to add
	ASIAllocation   map[string]int `json:"asi_allocation,omitempty"`   // {"strength": 1, "dexterity": 1}
	SpellsLearned   []string       `json:"spells_learned,omitempty"`   // new spell IDs
}

// LevelUpResponse contains the result of a level-up operation
type LevelUpResponse struct {
	Character  *models.Character   `json:"character"`
	NewLevel   int                 `json:"new_level"`
	ClassLevel *models.ClassLevel  `json:"class_level"`
	HistoryID  uuid.UUID           `json:"history_id,omitempty"`
}

// LevelUpPreview returns what will be available at the next level
type LevelUpPreview struct {
	CurrentLevel       int                `json:"current_level"`
	NextLevel          int                `json:"next_level"`
	ClassLevel         *models.ClassLevel `json:"class_level"`
	ASIPointsAvailable int                `json:"asi_points_available"`
	NewSpellsAllowed   int                `json:"new_spells_allowed"`
}

// LevelUp advances a character by one level, applying choices and saving history
func (s *LevelUpService) LevelUp(userID uuid.UUID, req LevelUpRequest) (*LevelUpResponse, error) {
	// 1. Load character with skills
	character, err := s.characterRepo.FindByIDWithSkills(req.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	// 2. Verify ownership
	if character.UserID != userID {
		return nil, ErrUnauthorized
	}

	// 3. Check max level
	if character.Level >= 20 {
		return nil, ErrMaxLevelReached
	}

	newLevel := character.Level + 1

	// 4. Get the new class level data (with LevelFeatures)
	newClassLevel, err := s.findClassLevelWithFeatures(character.ClassID, newLevel)
	if err != nil {
		return nil, fmt.Errorf("class level data not found: %w", err)
	}

	// 5. Validate ASI allocation if this level grants ASI points
	if newClassLevel.AbilityScoreImprovement > 0 && len(req.ASIAllocation) > 0 {
		totalAlloc := 0
		for _, v := range req.ASIAllocation {
			totalAlloc += v
		}
		if totalAlloc > newClassLevel.AbilityScoreImprovement {
			return nil, ErrInsufficientASIPoints
		}
	}

	// 6. Create snapshot of current state BEFORE changes
	snapshot := models.CharacterSnapshotData{
		Level:              character.Level,
		Strength:           character.Strength,
		Dexterity:          character.Dexterity,
		Constitution:       character.Constitution,
		Intelligence:       character.Intelligence,
		Wisdom:             character.Wisdom,
		Charisma:           character.Charisma,
		SpellbookIDs:       []string(character.SpellbookIDs),
		CurrentSpellPoints: character.CurrentSpellPoints,
	}
	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Copy current skill IDs for snapshot
	skillSnapshot := make([]string, len(character.SkillProficiencyIDs))
	copy(skillSnapshot, character.SkillProficiencyIDs)

	// 7. Execute in transaction
	var history *models.LevelUpHistory
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Apply level increase
		character.Level = newLevel

		// Apply ASI allocation
		for ability, points := range req.ASIAllocation {
			switch ability {
			case "strength":
				character.Strength += points
			case "dexterity":
				character.Dexterity += points
			case "constitution":
				character.Constitution += points
			case "intelligence":
				character.Intelligence += points
			case "wisdom":
				character.Wisdom += points
			case "charisma":
				character.Charisma += points
			}
		}

		// Add new spells to spellbook
		if len(req.SpellsLearned) > 0 {
			character.SpellbookIDs = append(character.SpellbookIDs, req.SpellsLearned...)
		}

		// Reset spell points to new max
		character.CurrentSpellPoints = newClassLevel.MaxSpellPoints

		// Save character
		if err := tx.Save(character).Error; err != nil {
			return err
		}

		// Add new skill proficiencies
		for _, skillID := range req.SkillSelections {
			cs := &models.CharacterSkill{
				CharacterID: character.ID,
				SkillID:     skillID,
				Proficient:  true,
			}
			if err := tx.Create(cs).Error; err != nil {
				return err
			}
		}

		// Create history record
		asiJSON, _ := json.Marshal(req.ASIAllocation)

		// Extract feature names for recording
		var featuresGained []string
		for _, f := range newClassLevel.LevelFeatures {
			featuresGained = append(featuresGained, f.Name)
		}

		history = &models.LevelUpHistory{
			CharacterID:       character.ID,
			UserID:            userID,
			Level:             newLevel,
			SkillSelections:   req.SkillSelections,
			ASIAllocation:     asiJSON,
			SpellsLearned:     req.SpellsLearned,
			FeaturesGained:    featuresGained,
			CharacterSnapshot: snapshotJSON,
			SkillSnapshot:     skillSnapshot,
		}
		return tx.Create(history).Error
	})

	if err != nil {
		return nil, fmt.Errorf("level-up transaction failed: %w", err)
	}

	return &LevelUpResponse{
		Character:  character,
		NewLevel:   newLevel,
		ClassLevel: newClassLevel,
		HistoryID:  history.ID,
	}, nil
}

// LevelDown reverts a character to their previous level using stored history
func (s *LevelUpService) LevelDown(userID uuid.UUID, characterID uuid.UUID) (*LevelUpResponse, error) {
	// 1. Load character
	character, err := s.characterRepo.FindByID(characterID)
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

// GetLevelHistory returns all level-up history for a character
func (s *LevelUpService) GetLevelHistory(userID uuid.UUID, characterID uuid.UUID) ([]*models.LevelUpHistory, error) {
	character, err := s.characterRepo.FindByID(characterID)
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
	character, err := s.characterRepo.FindByIDWithSkills(characterID)
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

	return &LevelUpPreview{
		CurrentLevel:       character.Level,
		NextLevel:          nextLevel,
		ClassLevel:         classLevel,
		ASIPointsAvailable: classLevel.AbilityScoreImprovement,
		NewSpellsAllowed:   newSpellsAllowed,
	}, nil
}

// findClassLevelWithFeatures fetches a ClassLevel with its LevelFeatures preloaded
func (s *LevelUpService) findClassLevelWithFeatures(classID uuid.UUID, level int) (*models.ClassLevel, error) {
	var cl models.ClassLevel
	if err := s.db.Preload("LevelFeatures", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&cl, "class_id = ? AND level = ?", classID, level).Error; err != nil {
		return nil, err
	}
	return &cl, nil
}
