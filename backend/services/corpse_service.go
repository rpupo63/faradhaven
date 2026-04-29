package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
)

// CorpseService handles corpse management for harvesting and psychometry
type CorpseService struct {
	corpseRepo *database.CorpseRepo
}

// NewCorpseService creates a new CorpseService
func NewCorpseService(corpseRepo *database.CorpseRepo) *CorpseService {
	return &CorpseService{corpseRepo: corpseRepo}
}

// CreateCorpseRequest contains the data needed to create a corpse
type CreateCorpseRequest struct {
	MapID               *uuid.UUID          `json:"map_id,omitempty"`
	Name                string              `json:"name"`
	CreatureType        models.CreatureType `json:"creature_type"`
	CreatureSize        models.CreatureSize `json:"creature_size,omitempty"`
	ChallengeRating     float64             `json:"challenge_rating,omitempty"`
	GridX               *int                `json:"grid_x,omitempty"`
	GridY               *int                `json:"grid_y,omitempty"`
	AvailableComponents []string            `json:"available_components,omitempty"`
	ComponentYield      int                 `json:"component_yield,omitempty"`
	SourceBeastID       *uuid.UUID          `json:"source_beast_id,omitempty"`
	ExpiresInMinutes    *int                `json:"expires_in_minutes,omitempty"`
}

// CreateCorpse creates a new corpse
func (s *CorpseService) CreateCorpse(req CreateCorpseRequest) (*models.Corpse, error) {
	size := req.CreatureSize
	if size == "" || !size.IsValid() {
		size = models.SizeMedium
	}

	yield := req.ComponentYield
	if yield == 0 {
		yield = 1
	}

	var expiresAt *time.Time
	if req.ExpiresInMinutes != nil && *req.ExpiresInMinutes > 0 {
		exp := time.Now().Add(time.Duration(*req.ExpiresInMinutes) * time.Minute)
		expiresAt = &exp
	}

	corpse := &models.Corpse{
		MapID:               req.MapID,
		Name:                req.Name,
		CreatureType:        req.CreatureType,
		CreatureSize:        size,
		ChallengeRating:     req.ChallengeRating,
		GridX:               req.GridX,
		GridY:               req.GridY,
		AvailableComponents: pq.StringArray(req.AvailableComponents),
		ComponentYield:      yield,
		SourceBeastID:       req.SourceBeastID,
		DiedAt:              time.Now(),
		ExpiresAt:           expiresAt,
	}

	if err := s.corpseRepo.Add(corpse); err != nil {
		return nil, err
	}

	return corpse, nil
}

// HarvestResult contains the result of harvesting a corpse
type HarvestResult struct {
	ComponentsYielded int      `json:"components_yielded"`
	ComponentIDs      []string `json:"component_ids,omitempty"`
	Message           string   `json:"message"`
}

// HarvestCorpse harvests components from a corpse (Ironwright Scavenge)
func (s *CorpseService) HarvestCorpse(corpseID uuid.UUID) (*HarvestResult, error) {
	corpse, err := s.corpseRepo.FindByID(corpseID)
	if err != nil {
		return nil, errors.New("corpse not found")
	}

	if !corpse.CanBeHarvested() {
		if corpse.HasBeenHarvested {
			return nil, errors.New("corpse has already been harvested")
		}
		return nil, errors.New("corpse is too old to harvest (must be within 1 hour of death)")
	}

	// Mark as harvested
	corpse.HasBeenHarvested = true
	if err := s.corpseRepo.Update(corpse); err != nil {
		return nil, err
	}

	return &HarvestResult{
		ComponentsYielded: corpse.ComponentYield,
		ComponentIDs:      []string(corpse.AvailableComponents),
		Message:           "Corpse harvested successfully",
	}, nil
}

// PsychometryResult contains the result of consuming a corpse for psychometry
type PsychometryResult struct {
	CreatureType    string  `json:"creature_type"`
	ChallengeRating float64 `json:"challenge_rating"`
	TraumaDC        int     `json:"trauma_dc"`
	Message         string  `json:"message"`
}

// ConsumeCorpse consumes a corpse for Lorewright Visceral Psychometry
func (s *CorpseService) ConsumeCorpse(corpseID uuid.UUID, characterLevel int) (*PsychometryResult, error) {
	corpse, err := s.corpseRepo.FindByID(corpseID)
	if err != nil {
		return nil, errors.New("corpse not found")
	}

	if !corpse.CanBeConsumed() {
		if corpse.HasBeenConsumed {
			return nil, errors.New("corpse has already been consumed")
		}
		return nil, errors.New("corpse is too old to consume (must be within 1 hour of death)")
	}

	// Mark as consumed
	corpse.HasBeenConsumed = true
	if err := s.corpseRepo.Update(corpse); err != nil {
		return nil, err
	}

	// Calculate Trauma DC: 10 + (CR - (Level / 3))
	traumaDC := 10 + int(corpse.ChallengeRating) - (characterLevel / 3)
	if traumaDC < 10 {
		traumaDC = 10
	}

	return &PsychometryResult{
		CreatureType:    string(corpse.CreatureType),
		ChallengeRating: corpse.ChallengeRating,
		TraumaDC:        traumaDC,
		Message:         "Psychometry ritual complete - make a Wisdom save against the Trauma DC",
	}, nil
}

// GetCorpse returns a single corpse by ID
func (s *CorpseService) GetCorpse(corpseID uuid.UUID) (*models.Corpse, error) {
	return s.corpseRepo.FindByID(corpseID)
}

// GetCorpses returns all corpses, optionally filtered by map
func (s *CorpseService) GetCorpses(mapID *uuid.UUID) ([]*models.Corpse, error) {
	if mapID != nil {
		return s.corpseRepo.FindByMapID(*mapID)
	}
	return s.corpseRepo.FindAll()
}

// UpdateCorpse persists changes to a corpse
func (s *CorpseService) UpdateCorpse(corpse *models.Corpse) error {
	return s.corpseRepo.Update(corpse)
}

// DeleteCorpse removes a corpse
func (s *CorpseService) DeleteCorpse(corpseID uuid.UUID) error {
	return s.corpseRepo.Delete(corpseID)
}

// CleanupExpiredCorpses removes all expired corpses
func (s *CorpseService) CleanupExpiredCorpses() (int64, error) {
	return s.corpseRepo.DeleteExpired()
}
