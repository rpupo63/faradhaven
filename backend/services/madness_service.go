package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_classes"
)

// MadnessService handles the Lorewright's Madness Die rolls
type MadnessService struct {
	charRepo  database.CharacterRepository
	classRepo database.ClassRepository
	// Other repos/services needed for applying complex effects would go here
	// For now, this service only identifies the effect.
}

// NewMadnessService creates a new MadnessService
func NewMadnessService(
	charRepo database.CharacterRepository,
	classRepo database.ClassRepository,
) *MadnessService {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	return &MadnessService{
		charRepo:  charRepo,
		classRepo: classRepo,
	}
}

// MadnessRollResult contains the outcome of a Madness Die roll
type MadnessRollResult struct {
	CharacterID uuid.UUID `json:"character_id"`
	Roll        int       `json:"roll"`
	MadnessDie  int       `json:"madness_die"`
	Effect      string    `json:"effect,omitempty"`
}

// RollMadness rolls the Lorewright's Madness Die and determines an effect if necessary.
// It returns the roll result. The application of the effect is handled externally.
func (s *MadnessService) RollMadness(characterID uuid.UUID) (*MadnessRollResult, error) {
	character, err := s.charRepo.FindByIDWithRelations(characterID) // Need relations for CurrentClassLevel
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	if character.Class.Name != "The Lorewright" {
		return nil, errors.New("character is not a Lorewright")
	}

	classLevel, err := s.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
	if err != nil {
		return nil, errors.New("could not determine character's current class level")
	}

	madnessDie := s.resolveLorewrightMadnessDie(character.ClassID, character.Level, classLevel)
	if madnessDie <= 0 {
		return nil, errors.New("could not determine madness die size for this class level")
	}
	roll := rand.Intn(madnessDie) + 1 // Roll 1 to madnessDie

	result := &MadnessRollResult{
		CharacterID: characterID,
		Roll:        roll,
		MadnessDie:  madnessDie,
	}

	if roll == 1 {
		madnessTable := faradhaven_classes.LorewrightMadnessTable()
		if n := len(madnessTable); n > 0 {
			effectIndex := rand.Intn(n) + 1
			result.Effect = madnessTable[effectIndex]
		}
	}

	return result, nil
}

// resolveLorewrightMadnessDie prefers DB class_level_resources, then BardicInspiration, then seed data.
func (s *MadnessService) resolveLorewrightMadnessDie(classID uuid.UUID, level int, classLevel *models.ClassLevel) int {
	die := s.classRepo.GetLevelResourceValue(classID, level, "madness_die")
	if die > 0 {
		return die
	}
	if rows, err := s.classRepo.FindLevelResourcesByClassLevel(classLevel.ID); err == nil {
		for _, row := range rows {
			if row.ResourceKey == "madness_die" && row.Value > 0 {
				return row.Value
			}
		}
	}
	if classLevel.BardicInspiration > 0 {
		return classLevel.BardicInspiration
	}
	return faradhaven_classes.LorewrightMadnessDieForLevel(level)
}
