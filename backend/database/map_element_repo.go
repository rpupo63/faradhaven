package database

import (
	"log"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type MapElementRepo struct {
	db *gorm.DB
}

func NewMapElementRepo(db *gorm.DB) *MapElementRepo {
	return &MapElementRepo{db: db}
}

func (r *MapElementRepo) Create(element *models.MapElement) error {
	if err := r.db.Create(element).Error; err != nil {
		log.Printf("Failed to create map element in DB: %v", err)
		return err
	}
	return nil
}

// CreateWithProperties creates a MapElement and its associated property record
// in a single transaction.
func (r *MapElementRepo) CreateWithProperties(element *models.MapElement) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		trapProps := element.TrapProperties
		dtProps := element.DifficultTerrainProperties
		elevProps := element.ElevationProperties
		wallProps := element.WallProperties

		element.TrapProperties = nil
		element.DifficultTerrainProperties = nil
		element.ElevationProperties = nil
		element.WallProperties = nil

		if err := tx.Create(element).Error; err != nil {
			return err
		}

		if trapProps != nil {
			trapProps.MapElementID = element.ID
			if err := tx.Create(trapProps).Error; err != nil {
				return err
			}
			element.TrapProperties = trapProps
		}
		if dtProps != nil {
			dtProps.MapElementID = element.ID
			if err := tx.Create(dtProps).Error; err != nil {
				return err
			}
			element.DifficultTerrainProperties = dtProps
		}
		if elevProps != nil {
			elevProps.MapElementID = element.ID
			if err := tx.Create(elevProps).Error; err != nil {
				return err
			}
			element.ElevationProperties = elevProps
		}
		if wallProps != nil {
			wallProps.MapElementID = element.ID
			if err := tx.Create(wallProps).Error; err != nil {
				return err
			}
			element.WallProperties = wallProps
		}

		return nil
	})
}

// CreateMultipleWithProperties creates several MapElements and their property records
// in a single transaction, replacing the per-element loop in the bulk-create handler.
func (r *MapElementRepo) CreateMultipleWithProperties(elements []*models.MapElement) error {
	if len(elements) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, element := range elements {
			trapProps := element.TrapProperties
			dtProps := element.DifficultTerrainProperties
			elevProps := element.ElevationProperties
			wallProps := element.WallProperties

			element.TrapProperties = nil
			element.DifficultTerrainProperties = nil
			element.ElevationProperties = nil
			element.WallProperties = nil

			if err := tx.Create(element).Error; err != nil {
				return err
			}

			if trapProps != nil {
				trapProps.MapElementID = element.ID
				if err := tx.Create(trapProps).Error; err != nil {
					return err
				}
				element.TrapProperties = trapProps
			}
			if dtProps != nil {
				dtProps.MapElementID = element.ID
				if err := tx.Create(dtProps).Error; err != nil {
					return err
				}
				element.DifficultTerrainProperties = dtProps
			}
			if elevProps != nil {
				elevProps.MapElementID = element.ID
				if err := tx.Create(elevProps).Error; err != nil {
					return err
				}
				element.ElevationProperties = elevProps
			}
			if wallProps != nil {
				wallProps.MapElementID = element.ID
				if err := tx.Create(wallProps).Error; err != nil {
					return err
				}
				element.WallProperties = wallProps
			}
		}
		return nil
	})
}

func (r *MapElementRepo) GetByID(id uuid.UUID) (*models.MapElement, error) {
	var element models.MapElement
	if err := r.db.
		Preload("TrapProperties").
		Preload("DifficultTerrainProperties").
		Preload("ElevationProperties").
		Preload("WallProperties").
		First(&element, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &element, nil
}

func (r *MapElementRepo) Update(element *models.MapElement) error {
	return r.db.Save(element).Error
}

func (r *MapElementRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.MapElement{}, "id = ?", id).Error
}
