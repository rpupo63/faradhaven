package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// SavedSpellService handles Speed Dial spell management
type SavedSpellService struct {
	db             *gorm.DB
	savedSpellRepo *database.SavedSpellRepo
	CharacterRepo  *database.CharacterRepo
	classRepo      *database.ClassRepo
}

// NewSavedSpellService creates a new SavedSpellService
func NewSavedSpellService(
	db *gorm.DB,
	savedSpellRepo *database.SavedSpellRepo,
		CharacterRepo *database.CharacterRepo,
		classRepo     *database.ClassRepo,
	) *SavedSpellService {
		return &SavedSpellService{
			db:             db,
			savedSpellRepo: savedSpellRepo,
			CharacterRepo:  CharacterRepo,
			classRepo:      classRepo,
		}
}

// SaveSpellRequest contains the data needed to save a spell to a slot
type SaveSpellRequest struct {
	Name         string      `json:"name"`
	ComponentIDs []uuid.UUID `json:"component_ids"`
	TotalCost    int         `json:"total_cost,omitempty"`
	DamageType   *string     `json:"damage_type,omitempty"`
	Description  *string     `json:"description,omitempty"`
	MaxUses      *int        `json:"max_uses,omitempty"`
}

// SaveSpell saves a spell combination to a Speed Dial slot
func (s *SavedSpellService) SaveSpell(characterID uuid.UUID, slotIndex int, req SaveSpellRequest) (*models.SavedSpell, error) {
	// Get character to check slot limits
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, errors.New("character not found")
	}

	// Get class-level resource values for Speed Dial slots and spell length
	resourceMap, err := s.classRepo.GetLevelResourceMap(character.ClassID, character.Level)
	if err != nil {
		return nil, fmt.Errorf("could not load class level resources: %w", err)
	}

	maxSlots := resourceMap["speed_dial_slots"]
	if maxSlots == 0 {
		maxSlots = 1 // Default to 1 slot if not set
	}

	if slotIndex < 0 || slotIndex >= maxSlots {
		return nil, fmt.Errorf("invalid slot index (must be 0-%d)", maxSlots-1)
	}

	// Calculate tier-based cost from component data
	computedCost, err := s.CalculateBlueprintCost(req.ComponentIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate blueprint cost: %w", err)
	}

	// Convert component UUIDs to string array
	componentStrs := make([]string, len(req.ComponentIDs))
	for i, id := range req.ComponentIDs {
		componentStrs[i] = id.String()
	}

	// Check if slot is already occupied and remove existing
	existing, _ := s.savedSpellRepo.FindByCharacterAndSlot(characterID, slotIndex)
	if existing != nil {
		_ = s.savedSpellRepo.Delete(existing.ID)
	}

	// Set uses remaining to max uses if specified
	var usesRemaining *int
	if req.MaxUses != nil {
		usesRemaining = req.MaxUses
	}

	savedSpell := &models.SavedSpell{
		CharacterID:   characterID,
		Name:          req.Name,
		ComponentIDs:  pq.StringArray(componentStrs),
		SlotIndex:     slotIndex,
		TotalCost:     computedCost,
		DamageType:    req.DamageType,
		Description:   req.Description,
		MaxUses:       req.MaxUses,
		UsesRemaining: usesRemaining,
	}

	if err := s.savedSpellRepo.Add(savedSpell); err != nil {
		return nil, err
	}

	return savedSpell, nil
}

// CalculateBlueprintCost queries component Tiers from the DB and returns the total Stability cost.
// Cost formula: sum of each component's Tier value (Tier 1 = 1 charge, Tier 2 = 2 charges).
func (s *SavedSpellService) CalculateBlueprintCost(componentIDs []uuid.UUID) (int, error) {
	if len(componentIDs) == 0 {
		return 0, nil
	}

	var components []models.Component
	err := s.db.Where("id IN ?", componentIDs).Find(&components).Error
	if err != nil {
		return 0, err
	}

	if len(components) != len(componentIDs) {
		return 0, fmt.Errorf("some components not found (expected %d, found %d)", len(componentIDs), len(components))
	}

	totalCost := 0
	for _, comp := range components {
		tier := comp.Tier
		if tier == 0 {
			tier = 1 // default to Tier 1
		}
		totalCost += tier
	}

	return totalCost, nil
}

// CastSavedSpell casts a spell from a Speed Dial slot
func (s *SavedSpellService) CastSavedSpell(characterID uuid.UUID, slotIndex int) (*models.SavedSpell, error) {
	savedSpell, err := s.savedSpellRepo.FindByCharacterAndSlot(characterID, slotIndex)
	if err != nil {
		return nil, errors.New("no spell in this slot")
	}

	// Check uses remaining
	if savedSpell.UsesRemaining != nil && *savedSpell.UsesRemaining <= 0 {
		return nil, errors.New("no uses remaining")
	}

	// Decrement uses if limited
	if savedSpell.UsesRemaining != nil {
		newUses := *savedSpell.UsesRemaining - 1
		savedSpell.UsesRemaining = &newUses
		if err := s.savedSpellRepo.Update(savedSpell); err != nil {
			return nil, err
		}
	}

	return savedSpell, nil
}

// GetSpeedDial returns all saved spells for a character
func (s *SavedSpellService) GetSpeedDial(characterID uuid.UUID) ([]*models.SavedSpell, error) {
	return s.savedSpellRepo.FindByCharacterID(characterID)
}

// ClearSlot removes a saved spell from a slot
func (s *SavedSpellService) ClearSlot(characterID uuid.UUID, slotIndex int) error {
	savedSpell, err := s.savedSpellRepo.FindByCharacterAndSlot(characterID, slotIndex)
	if err != nil {
		return errors.New("no spell in this slot")
	}

	return s.savedSpellRepo.Delete(savedSpell.ID)
}

// ResetUsesOnRest resets uses for all saved spells on rest
func (s *SavedSpellService) ResetUsesOnRest(characterID uuid.UUID) error {
	return s.savedSpellRepo.ResetUsesOnRest(characterID)
}
