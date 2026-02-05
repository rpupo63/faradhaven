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

	if err := s.db.Save(character).Error; err != nil {
		return nil, fmt.Errorf("failed to long rest: %w", err)
	}

	return character, nil
}
