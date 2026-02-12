package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type MapRepo struct {
	db *gorm.DB
}

func NewMapRepo(db *gorm.DB) *MapRepo {
	return &MapRepo{db: db}
}

// CreateMap creates a new game map
func (r *MapRepo) CreateMap(gameMap *models.GameMap) error {
	return r.db.Create(gameMap).Error
}

// GetMapByID retrieves a map by its ID, including tokens
func (r *MapRepo) GetMapByID(id uuid.UUID) (*models.GameMap, error) {
	var gameMap models.GameMap
	if err := r.db.Preload("Tokens").First(&gameMap, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &gameMap, nil
}

// GetMapByRoomCode retrieves a map by its room code, including tokens
func (r *MapRepo) GetMapByRoomCode(roomCode string) (*models.GameMap, error) {
	var gameMap models.GameMap
	if err := r.db.Preload("Tokens").First(&gameMap, "room_code = ?", roomCode).Error; err != nil {
		return nil, err
	}
	return &gameMap, nil
}

// GetMapsByOwner retrieves all maps owned by a specific user
func (r *MapRepo) GetMapsByOwner(ownerID uuid.UUID) ([]models.GameMap, error) {
	var maps []models.GameMap
	if err := r.db.Where("owner_id = ?", ownerID).Find(&maps).Error; err != nil {
		return nil, err
	}
	return maps, nil
}

// UpdateMap updates an existing map
func (r *MapRepo) UpdateMap(gameMap *models.GameMap) error {
	return r.db.Save(gameMap).Error
}

// DeleteMap deletes a map by its ID
func (r *MapRepo) DeleteMap(id uuid.UUID) error {
	return r.db.Delete(&models.GameMap{}, "id = ?", id).Error
}

// AddToken adds a new token to a map
func (r *MapRepo) AddToken(token *models.MapToken) error {
	return r.db.Create(token).Error
}

// UpdateToken updates an existing token
func (r *MapRepo) UpdateToken(token *models.MapToken) error {
	return r.db.Save(token).Error
}

// GetTokenByID retrieves a token by its ID
func (r *MapRepo) GetTokenByID(id uuid.UUID) (*models.MapToken, error) {
	var token models.MapToken
	if err := r.db.First(&token, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

// DeleteToken deletes a token by its ID
func (r *MapRepo) DeleteToken(id uuid.UUID) error {
	return r.db.Delete(&models.MapToken{}, "id = ?", id).Error
}

// GetTokensByInitiativeOrder returns tokens with initiative set, ordered ASC
func (r *MapRepo) GetTokensByInitiativeOrder(mapID uuid.UUID) ([]models.MapToken, error) {
	var tokens []models.MapToken
	err := r.db.Where("map_id = ? AND initiative_order IS NOT NULL", mapID).
		Order("initiative_order ASC").Find(&tokens).Error
	return tokens, err
}

// SetTokenInitiativeOrder updates a single token's initiative order
func (r *MapRepo) SetTokenInitiativeOrder(tokenID uuid.UUID, order *int) error {
	return r.db.Model(&models.MapToken{}).Where("id = ?", tokenID).
		Update("initiative_order", order).Error
}
