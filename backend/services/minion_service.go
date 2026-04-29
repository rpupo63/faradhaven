package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

// MinionService handles minion management for constructs, echoes, etc.
type MinionService struct {
	db            *gorm.DB
	minionRepo    *database.MinionRepo
	CharacterRepo *database.CharacterRepo
	classRepo     *database.ClassRepo
	componentRepo *database.ComponentRepo
}

// NewMinionService creates a new MinionService
func NewMinionService(
	db *gorm.DB,
	minionRepo *database.MinionRepo,
	CharacterRepo *database.CharacterRepo,
	classRepo *database.ClassRepo,
	componentRepo *database.ComponentRepo,
) *MinionService {
	return &MinionService{
		db:            db,
		minionRepo:    minionRepo,
		CharacterRepo: CharacterRepo,
		classRepo:     classRepo,
		componentRepo: componentRepo,
	}
}

// CreateConstructRequest contains the data needed to create a construct
type CreateConstructRequest struct {
	TemplateKey string `json:"template_key"` // "sentry", "striker", "titan"
	CustomName  string `json:"custom_name,omitempty"`
}

// constructMetadata stores actions and recipe component IDs for reclamation
type constructMetadata struct {
	Actions            []models.ConstructAction `json:"actions"`
	RecipeComponentIDs []string                 `json:"recipe_component_ids,omitempty"`
}

// CreateConstruct creates an Ironwright construct, validating and deducting required components
func (s *MinionService) CreateConstruct(characterID uuid.UUID, req CreateConstructRequest) (*models.Minion, error) {
	// Get character to check level and concurrency limit
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, errors.New("character not found")
	}

	// Get construct template
	templates := models.GetConstructTemplates()
	template, exists := templates[req.TemplateKey]
	if !exists {
		return nil, fmt.Errorf("invalid construct template: %s", req.TemplateKey)
	}

	// Drones are not created via this path
	if req.TemplateKey == "drone" {
		return nil, errors.New("use CreateDrone for drone constructs")
	}

	// Check concurrency limit (from class level resources)
	concurrencyLimit := s.classRepo.GetLevelResourceValue(character.ClassID, character.Level, "concurrency_limit")
	if concurrencyLimit == 0 {
		concurrencyLimit = 1 // Default to 1 if not set
	}

	currentCount, err := s.minionRepo.CountByCharacterAndType(characterID, models.MinionTypeConstruct)
	if err != nil {
		return nil, err
	}

	if int(currentCount) >= concurrencyLimit {
		return nil, fmt.Errorf("construct limit reached (%d/%d)", currentCount, concurrencyLimit)
	}

	// Validate and deduct required components
	recipeComponentIDs, err := s.validateAndDeductComponents(characterID, template.RequiredComponents)
	if err != nil {
		return nil, err
	}

	// Calculate HP based on level
	hp := models.CalculateConstructHP(template, character.Level)

	// Serialize actions + recipe component IDs to metadata
	meta := constructMetadata{
		Actions:            template.Actions,
		RecipeComponentIDs: recipeComponentIDs,
	}
	metaJSON, _ := json.Marshal(meta)

	name := template.Name
	if req.CustomName != "" {
		name = req.CustomName
	}

	minion := &models.Minion{
		CharacterID:   characterID,
		Name:          name,
		MinionType:    models.MinionTypeConstruct,
		ArmorClass:    template.BaseAC,
		MaxHP:         hp,
		CurrentHP:     hp,
		Speed:         template.BaseSpeed,
		ComponentCost: len(template.RequiredComponents),
		TemplateKey:   &req.TemplateKey,
		IsActive:      true,
		Metadata:      metaJSON,
	}

	if err := s.minionRepo.Add(minion); err != nil {
		return nil, err
	}

	return minion, nil
}

// validateAndDeductComponents checks that the character has all required components
// and deducts them. Returns the component IDs that were consumed.
func (s *MinionService) validateAndDeductComponents(characterID uuid.UUID, componentNames []string) ([]string, error) {
	if len(componentNames) == 0 {
		return nil, nil
	}

	// Look up component UUIDs by name
	components, err := s.componentRepo.GetComponentsByNames(componentNames)
	if err != nil {
		return nil, fmt.Errorf("failed to look up components: %w", err)
	}

	// Build name→component map
	compMap := make(map[string]models.Component, len(components))
	for _, c := range components {
		compMap[c.Name] = c
	}

	// Verify all required components exist in the DB
	for _, name := range componentNames {
		if _, found := compMap[name]; !found {
			return nil, fmt.Errorf("unknown component: %s", name)
		}
	}

	// Check character has each component (count >= 1) and build deduction list
	// Count occurrences of each component name (a recipe might need duplicates)
	needed := make(map[string]int, len(componentNames))
	for _, name := range componentNames {
		needed[name]++
	}

	var charComp models.CharacterComponent
	for name, count := range needed {
		comp := compMap[name]
		err := s.db.Where("character_id = ? AND component_id = ?", characterID, comp.ID).First(&charComp).Error
		if err != nil || charComp.Count < count {
			return nil, fmt.Errorf("missing required component: %s (need %d)", name, count)
		}
	}

	// Deduct each component
	var recipeIDs []string
	for _, name := range componentNames {
		comp := compMap[name]
		if err := s.CharacterRepo.UpdateComponentCount(characterID, comp.ID, -1); err != nil {
			return nil, fmt.Errorf("failed to deduct component %s: %w", name, err)
		}
		recipeIDs = append(recipeIDs, comp.ID.String())
	}

	return recipeIDs, nil
}

// DestroyConstruct removes a construct and optionally returns reclaimed components to inventory
func (s *MinionService) DestroyConstruct(minionID uuid.UUID, reclaim bool) (int, error) {
	minion, err := s.minionRepo.FindByID(minionID)
	if err != nil {
		return 0, errors.New("construct not found")
	}

	if minion.MinionType != models.MinionTypeConstruct && minion.MinionType != models.MinionTypeDrone {
		return 0, errors.New("minion is not a construct")
	}

	reclaimedComponents := 0
	if reclaim && minion.Metadata != nil {
		// Try to read stored recipe component IDs from metadata
		var meta constructMetadata
		if err := json.Unmarshal(minion.Metadata, &meta); err == nil && len(meta.RecipeComponentIDs) > 0 {
			// Reclaim first 50% (rounded down) of component IDs
			reclaimCount := len(meta.RecipeComponentIDs) / 2
			for i := 0; i < reclaimCount; i++ {
				compID, err := uuid.Parse(meta.RecipeComponentIDs[i])
				if err != nil {
					continue
				}
				if err := s.CharacterRepo.UpdateComponentCount(minion.CharacterID, compID, 1); err != nil {
					continue
				}
				reclaimedComponents++
			}
		} else if minion.ComponentCost > 0 {
			// Fallback for legacy constructs without recipe metadata
			reclaimedComponents = minion.ComponentCost / 2
		}
	}

	if err := s.minionRepo.Delete(minionID); err != nil {
		return 0, err
	}

	return reclaimedComponents, nil
}

// CreateDroneRequest contains the data needed to create a Swarm Engineer drone
type CreateDroneRequest struct {
	ComponentID uuid.UUID `json:"component_id"` // Which component to spend
	CustomName  string    `json:"custom_name,omitempty"`
}

// CreateDrone creates a Swarm Engineer micro-construct
func (s *MinionService) CreateDrone(characterID uuid.UUID, req CreateDroneRequest) (*models.Minion, error) {
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, errors.New("character not found")
	}

	// Validate character has the component (count >= 1)
	var charComp models.CharacterComponent
	err = s.db.Where("character_id = ? AND component_id = ?", characterID, req.ComponentID).First(&charComp).Error
	if err != nil || charComp.Count < 1 {
		return nil, errors.New("missing required component for drone")
	}

	// Count existing active drones
	droneCount, err := s.minionRepo.CountByCharacterAndType(characterID, models.MinionTypeDrone)
	if err != nil {
		return nil, err
	}

	// Max drones = INT modifier = (Intelligence - 10) / 2
	intMod := (character.Intelligence - 10) / 2
	if intMod < 1 {
		intMod = 1
	}
	if int(droneCount) >= intMod {
		return nil, fmt.Errorf("drone limit reached (%d/%d, INT modifier)", droneCount, intMod)
	}

	// Deduct 1 of the specified component
	if err := s.CharacterRepo.UpdateComponentCount(characterID, req.ComponentID, -1); err != nil {
		return nil, fmt.Errorf("failed to deduct component: %w", err)
	}

	droneTemplate := models.GetConstructTemplates()["drone"]

	// Store recipe in metadata
	meta := constructMetadata{
		Actions:            droneTemplate.Actions,
		RecipeComponentIDs: []string{req.ComponentID.String()},
	}
	metaJSON, _ := json.Marshal(meta)

	templateKey := "drone"
	name := "Drone"
	if req.CustomName != "" {
		name = req.CustomName
	}

	minion := &models.Minion{
		CharacterID:   characterID,
		Name:          name,
		MinionType:    models.MinionTypeDrone,
		ArmorClass:    droneTemplate.BaseAC,
		MaxHP:         1,
		CurrentHP:     1,
		Speed:         droneTemplate.BaseSpeed,
		ComponentCost: 1,
		TemplateKey:   &templateKey,
		IsActive:      true,
		Metadata:      metaJSON,
	}

	if err := s.minionRepo.Add(minion); err != nil {
		return nil, err
	}

	return minion, nil
}

// StoreEchoRequest contains the data needed to store an echo ability
type StoreEchoRequest struct {
	SlotIndex          int    `json:"slot_index"`
	AbilityName        string `json:"ability_name"`
	AbilityDescription string `json:"ability_description"`
	AbilityType        string `json:"ability_type"` // "skill", "action", "trait", "language"
	SourceCreatureType string `json:"source_creature_type"`
}

// StoreEcho stores a Lorewright echo ability in a slot
func (s *MinionService) StoreEcho(characterID uuid.UUID, req StoreEchoRequest) (*models.Minion, error) {
	// Get character to check echo slots
	character, err := s.CharacterRepo.FindByID(characterID)
	if err != nil {
		return nil, errors.New("character not found")
	}

	// Get max echo slots from class level resources
	maxSlots := s.classRepo.GetLevelResourceValue(character.ClassID, character.Level, "echo_slots")
	if maxSlots == 0 {
		return nil, errors.New("character has no echo slots available")
	}

	if req.SlotIndex < 0 || req.SlotIndex >= maxSlots {
		return nil, fmt.Errorf("invalid slot index (must be 0-%d)", maxSlots-1)
	}

	// Check if slot is already occupied and remove existing echo
	existing, _ := s.minionRepo.FindEchoBySlot(characterID, req.SlotIndex)
	if existing != nil {
		_ = s.minionRepo.Delete(existing.ID)
	}

	minion := &models.Minion{
		CharacterID:        characterID,
		Name:               req.AbilityName,
		MinionType:         models.MinionTypeEcho,
		AbilityName:        &req.AbilityName,
		AbilityDescription: &req.AbilityDescription,
		AbilityType:        &req.AbilityType,
		SourceCreatureType: &req.SourceCreatureType,
		SlotIndex:          &req.SlotIndex,
		IsActive:           false, // Echoes start inactive
	}

	if err := s.minionRepo.Add(minion); err != nil {
		return nil, err
	}

	return minion, nil
}

// ActivateEcho activates an echo (sets it as the current active ability)
func (s *MinionService) ActivateEcho(minionID uuid.UUID) (*models.Minion, error) {
	minion, err := s.minionRepo.FindByID(minionID)
	if err != nil {
		return nil, errors.New("echo not found")
	}

	if minion.MinionType != models.MinionTypeEcho {
		return nil, errors.New("minion is not an echo")
	}

	// Deactivate all other echoes for this character
	echoes, err := s.minionRepo.FindByCharacterAndType(minion.CharacterID, models.MinionTypeEcho)
	if err != nil {
		return nil, err
	}

	for _, e := range echoes {
		if e.IsActive && e.ID != minionID {
			e.IsActive = false
			_ = s.minionRepo.Update(e)
		}
	}

	// Activate the selected echo
	minion.IsActive = true
	if err := s.minionRepo.Update(minion); err != nil {
		return nil, err
	}

	return minion, nil
}

// DeactivateEcho deactivates an echo
func (s *MinionService) DeactivateEcho(minionID uuid.UUID) error {
	minion, err := s.minionRepo.FindByID(minionID)
	if err != nil {
		return errors.New("echo not found")
	}

	minion.IsActive = false
	return s.minionRepo.Update(minion)
}

// GetMinions returns all minions for a character, optionally filtered by type
func (s *MinionService) GetMinions(characterID uuid.UUID, minionType *models.MinionType) ([]*models.Minion, error) {
	if minionType != nil {
		return s.minionRepo.FindByCharacterAndType(characterID, *minionType)
	}
	return s.minionRepo.FindByCharacterID(characterID)
}

// GetMinion returns a single minion by ID
func (s *MinionService) GetMinion(minionID uuid.UUID) (*models.Minion, error) {
	return s.minionRepo.FindByID(minionID)
}

// UpdateMinionHP updates a minion's HP (for constructs/summons)
func (s *MinionService) UpdateMinionHP(minionID uuid.UUID, delta int) (*models.Minion, error) {
	minion, err := s.minionRepo.FindByID(minionID)
	if err != nil {
		return nil, errors.New("minion not found")
	}

	minion.CurrentHP += delta
	if minion.CurrentHP < 0 {
		minion.CurrentHP = 0
	}
	if minion.CurrentHP > minion.MaxHP {
		minion.CurrentHP = minion.MaxHP
	}

	if err := s.minionRepo.Update(minion); err != nil {
		return nil, err
	}

	return minion, nil
}

// DeleteMinion removes a minion
func (s *MinionService) DeleteMinion(minionID uuid.UUID) error {
	return s.minionRepo.Delete(minionID)
}
