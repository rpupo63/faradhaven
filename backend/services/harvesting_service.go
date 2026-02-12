package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// HarvestingService handles the logic for the Lorewright's Visceral Psychometry
type HarvestingService struct {
	charRepo   database.CharacterRepository
	beastRepo  database.BeastRepository
	consumRepo database.ConsumptionHistoryRepository
	classRepo  database.ClassRepository
}

// NewHarvestingService creates a new HarvestingService
func NewHarvestingService(
	charRepo database.CharacterRepository,
	beastRepo database.BeastRepository,
	consumRepo database.ConsumptionHistoryRepository,
	classRepo database.ClassRepository,
) *HarvestingService {
	return &HarvestingService{
		charRepo:   charRepo,
		beastRepo:  beastRepo,
		consumRepo: consumRepo,
		classRepo:  classRepo,
	}
}

// HarvestableAbility represents a single skill or attack that can be harvested.
type HarvestableAbility struct {
	ID          string `json:"id"`          // e.g., "perception", "bite"
	Name        string `json:"name"`        // e.g., "Perception", "Bite"
	Type        string `json:"type"`        // "skill", "attack", "recipe"
	Description string `json:"description"` // e.g., "+5", "Melee Weapon Attack: +4 to hit, 1d6+2 piercing damage."
}

// HarvestableAbilitiesResponse contains all abilities that can be harvested from a beast.
type HarvestableAbilitiesResponse struct {
	BeastID      uuid.UUID            `json:"beast_id"`
	Abilities    []HarvestableAbility `json:"abilities"`
	FractureDC   int                  `json:"fracture_dc"`
	RequiresSave bool                 `json:"requires_save"`
}

// GetHarvestableAbilities retrieves a list of skills, attacks, and recipes from a beast.
func (s *HarvestingService) GetHarvestableAbilities(characterID, beastID uuid.UUID) (*HarvestableAbilitiesResponse, error) {
	beast, err := s.beastRepo.FindByIDWithRelations(beastID)
	if err != nil {
		return nil, errors.New("beast not found")
	}
	character, err := s.charRepo.FindByID(characterID)
	if err != nil {
		return nil, errors.New("character not found")
	}

	var abilities []HarvestableAbility

	// Get Skills
	for _, skill := range beast.Skills {
		abilities = append(abilities, HarvestableAbility{
			ID:          skill.Name,
			Name:        skill.Name,
			Type:        "skill",
			Description: fmt.Sprintf("Value: %d", skill.Value),
		})
	}

	// Get Attacks
	for _, attack := range beast.Attacks {
		var desc string
		if attack.Description != nil {
			desc = *attack.Description
		}
		abilities = append(abilities, HarvestableAbility{
			ID:          attack.Name,
			Name:        attack.Name,
			Type:        "attack",
			Description: desc,
		})
	}

	// TODO: Add logic for deriving Recipes

	requiresSave, dc := CalculateFractureDC(character.Level, beast.ParseChallengeRating())

	return &HarvestableAbilitiesResponse{
		BeastID:      beast.ID,
		Abilities:    abilities,
		FractureDC:   dc,
		RequiresSave: requiresSave,
	}, nil
}

// ConfirmHarvestRequest contains the data needed to confirm a harvest.
type ConfirmHarvestRequest struct {
	CharacterID uuid.UUID `json:"character_id"`
	BeastID     uuid.UUID `json:"beast_id"`
	AbilityID   string    `json:"ability_id"`
	AbilityType string    `json:"ability_type"` // "skill", "attack", or "recipe"
	DidFailSave bool      `json:"did_fail_save"`
}

// ConfirmHarvest slots a chosen ability into a character's Harvest Bank.
func (s *HarvestingService) ConfirmHarvest(req ConfirmHarvestRequest) (*models.Character, error) {
	character, err := s.charRepo.FindByIDWithRelations(req.CharacterID)
	if err != nil {
		return nil, errors.New("character not found")
	}
	beast, err := s.beastRepo.FindByID(req.BeastID)
	if err != nil {
		return nil, errors.New("beast not found")
	}

	// 1. Handle Trauma from failed save
	if req.DidFailSave {
		character.Trauma++
	}

	// 2. Log consumption for Predator's Strike
	history := &models.ConsumptionHistory{
		CharacterID:  character.ID,
		CreatureType: beast.Type,
	}
	if err := s.consumRepo.Create(history); err != nil {
		// Log error but don't block harvest
		fmt.Printf("Error logging consumption history: %v\n", err)
	}

	// 3. Update HarvestedAbilities
	var harvested models.HarvestedAbilities
	if len(character.HarvestedAbilities) > 0 {
		if err := json.Unmarshal(character.HarvestedAbilities, &harvested); err != nil {
			return nil, fmt.Errorf("could not parse existing harvested abilities: %w", err)
		}
	}

	resourceMap, err := s.classRepo.GetLevelResourceMap(character.ClassID, character.Level)
	if err != nil {
		return nil, errors.New("could not determine character's class level details")
	}

	switch req.AbilityType {
	case "skill":
		if len(harvested.Skills) >= resourceMap["skill_slots"] {
			return nil, errors.New("no available skill slots")
		}
		harvested.Skills = append(harvested.Skills, req.AbilityID)
	case "attack":
		if len(harvested.Attacks) >= resourceMap["attack_slots"] {
			return nil, errors.New("no available attack slots")
		}
		harvested.Attacks = append(harvested.Attacks, req.AbilityID)
	case "recipe":
		if len(harvested.Recipes) >= resourceMap["recipe_slots"] {
			return nil, errors.New("no available recipe slots")
		}
		harvested.Recipes = append(harvested.Recipes, req.AbilityID)
	default:
		return nil, fmt.Errorf("unknown ability type: %s", req.AbilityType)
	}

	rawJson, err := json.Marshal(harvested)
	if err != nil {
		return nil, fmt.Errorf("could not serialize harvested abilities: %w", err)
	}
	character.HarvestedAbilities = rawJson

	// 4. Save updated character
	if err := s.charRepo.Update(character); err != nil {
		return nil, err
	}

	return character, nil
}
