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
	ErrArchetypeRequired     = errors.New("archetype selection required at this level")
	ErrInvalidArchetype      = errors.New("invalid archetype for this class")
)

type LevelUpService struct {
	db            *gorm.DB
	characterRepo *database.CharacterRepo
	classRepo     *database.ClassRepo
	historyRepo   *database.LevelUpHistoryRepo
	archetypeRepo *database.ArchetypeRepo
}

func NewLevelUpService(
	db *gorm.DB,
	characterRepo *database.CharacterRepo,
	classRepo *database.ClassRepo,
	historyRepo *database.LevelUpHistoryRepo,
	archetypeRepo *database.ArchetypeRepo,
) *LevelUpService {
	return &LevelUpService{
		db:            db,
		characterRepo: characterRepo,
		classRepo:     classRepo,
		historyRepo:   historyRepo,
		archetypeRepo: archetypeRepo,
	}
}

// LevelUpRequest contains the choices made during level-up
type LevelUpRequest struct {
	CharacterID     uuid.UUID      `json:"character_id"`
	SkillSelections []string       `json:"skill_selections,omitempty"` // new skills to add
	ASIAllocation   map[string]int `json:"asi_allocation,omitempty"`   // {"strength": 1, "dexterity": 1}
	SpellsLearned   []string       `json:"spells_learned,omitempty"`   // new spell IDs
	HPRollResult    *int           `json:"hp_roll_result,omitempty"`   // nil = use average, otherwise the rolled value
	ArchetypeID     *uuid.UUID     `json:"archetype_id,omitempty"`     // required when reaching archetype level
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
	CurrentLevel            int                 `json:"current_level"`
	NextLevel               int                 `json:"next_level"`
	ClassLevel              *models.ClassLevel  `json:"class_level"`
	ASIPointsAvailable      int                 `json:"asi_points_available"`
	NewSpellsAllowed        int                 `json:"new_spells_allowed"`
	HitDie                  int                 `json:"hit_die"`                            // Class hit die (e.g., 10 for d10)
	ConMod                  int                 `json:"con_mod"`                            // Current CON modifier
	HPGainAverage           int                 `json:"hp_gain_average"`                    // Average HP gain (hitDie/2 + 1 + conMod)
	CurrentMaxHP            int                 `json:"current_max_hp"`                     // Current max HP before level up
	RequiresArchetypeChoice bool                `json:"requires_archetype_choice"`          // true if this level requires archetype selection
	AvailableArchetypes     []*models.Archetype `json:"available_archetypes,omitempty"`     // archetypes to choose from (if required)
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

	// 5b. Check if archetype selection is required
	var archetypeSelected *uuid.UUID
	if character.Class.ArchetypeLevel != nil && newLevel == *character.Class.ArchetypeLevel && character.ArchetypeID == nil {
		// This is the archetype selection level and character doesn't have one yet
		if req.ArchetypeID == nil {
			return nil, ErrArchetypeRequired
		}
		// Validate the archetype belongs to this class
		archetype, err := s.archetypeRepo.FindByID(*req.ArchetypeID)
		if err != nil || archetype.ClassID != character.ClassID {
			return nil, ErrInvalidArchetype
		}
		archetypeSelected = req.ArchetypeID
	}

	// 6. Create snapshot of current state BEFORE changes
	var archetypeIDStr *string
	if character.ArchetypeID != nil {
		str := character.ArchetypeID.String()
		archetypeIDStr = &str
	}
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
		CurrentHP:          character.CurrentHP,
		MaxHP:              character.MaxHP,
		TempHP:             character.TempHP,
		HitDiceUsed:        character.HitDiceUsed,
		ArchetypeID:        archetypeIDStr,
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

		// Calculate HP gain for this level
		conMod := (character.Constitution - 10) / 2
		var hpGain int
		if req.HPRollResult != nil {
			// Use the player's roll
			hpGain = *req.HPRollResult + conMod
		} else {
			// Use average (rounded up): (hitDie / 2) + 1
			avgHitDie := (character.Class.HitDie / 2) + 1
			hpGain = avgHitDie + conMod
		}
		// Minimum 1 HP per level
		if hpGain < 1 {
			hpGain = 1
		}
		character.MaxHP += hpGain
		// On level up, restore to new max HP (design decision)
		character.CurrentHP = character.MaxHP

		// Add new spells to spellbook
		if len(req.SpellsLearned) > 0 {
			character.SpellbookIDs = append(character.SpellbookIDs, req.SpellsLearned...)
		}

		// Reset spell points to new max
		character.CurrentSpellPoints = newClassLevel.MaxSpellPoints

		// Apply archetype selection if made at this level
		if archetypeSelected != nil {
			character.ArchetypeID = archetypeSelected
		}

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
			ArchetypeSelected: archetypeSelected,
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

// filterFeaturesByArchetype returns features that are either shared (nil archetype) or match the given archetype
func filterFeaturesByArchetype(features []models.LevelFeature, archetypeID *uuid.UUID) []models.LevelFeature {
	var filtered []models.LevelFeature
	for _, f := range features {
		// Include shared features (nil ArchetypeID) or features matching the character's archetype
		if f.ArchetypeID == nil || (archetypeID != nil && *f.ArchetypeID == *archetypeID) {
			filtered = append(filtered, f)
		}
	}
	return filtered
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