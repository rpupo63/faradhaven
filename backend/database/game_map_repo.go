package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type GameMapRepo struct {
	db *gorm.DB
}

func NewGameMapRepo(db *gorm.DB) *GameMapRepo {
	return &GameMapRepo{db: db}
}

func (r *GameMapRepo) Create(gameMap *models.GameMap) error {
	return r.db.Create(gameMap).Error
}

func (r *GameMapRepo) GetByID(id uuid.UUID) (*models.GameMap, error) {
	var gameMap models.GameMap
	if err := r.db.
		Preload("Tokens").
		Preload("Elements").
		Preload("Elements.TrapProperties").
		Preload("Elements.DifficultTerrainProperties").
		Preload("Elements.ElevationProperties").
		Preload("Elements.WallProperties").
		First(&gameMap, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &gameMap, nil
}

func (r *GameMapRepo) GetOwnerID(id uuid.UUID) (uuid.UUID, error) {
	var result struct {
		OwnerID uuid.UUID
	}
	if err := r.db.Model(&models.GameMap{}).Select("owner_id").First(&result, "id = ?", id).Error; err != nil {
		return uuid.Nil, err
	}
	return result.OwnerID, nil
}

func (r *GameMapRepo) GetByRoomCode(roomCode string) (*models.GameMap, error) {
	var gameMap models.GameMap
	if err := r.db.
		Preload("Tokens").
		Preload("Elements").
		Preload("Elements.TrapProperties").
		Preload("Elements.DifficultTerrainProperties").
		Preload("Elements.ElevationProperties").
		Preload("Elements.WallProperties").
		First(&gameMap, "room_code = ?", roomCode).Error; err != nil {
		return nil, err
	}
	return &gameMap, nil
}

func (r *GameMapRepo) GetByOwner(ownerID uuid.UUID) ([]models.GameMap, error) {
	var maps []models.GameMap
	if err := r.db.Where("owner_id = ?", ownerID).Find(&maps).Error; err != nil {
		return nil, err
	}
	return maps, nil
}

// GetByAssignedUser returns maps where the user has at least one assigned token but is not the owner.
func (r *GameMapRepo) GetByAssignedUser(userID uuid.UUID) ([]models.GameMap, error) {
	var maps []models.GameMap
	if err := r.db.
		Joins("JOIN map_tokens ON map_tokens.map_id = game_maps.id").
		Where("map_tokens.assigned_user_id = ? AND game_maps.owner_id != ?", userID, userID).
		Distinct("game_maps.*").
		Find(&maps).Error; err != nil {
		return nil, err
	}
	return maps, nil
}

func (r *GameMapRepo) Update(gameMap *models.GameMap) error {
	return r.db.Save(gameMap).Error
}

func (r *GameMapRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.GameMap{}, "id = ?", id).Error
}
