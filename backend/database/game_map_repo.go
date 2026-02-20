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

func (r *GameMapRepo) Update(gameMap *models.GameMap) error {
	return r.db.Save(gameMap).Error
}

func (r *GameMapRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.GameMap{}, "id = ?", id).Error
}
