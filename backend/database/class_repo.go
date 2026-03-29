package database

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type ClassRepository interface {
	FindAll() ([]*models.Class, error)
	FindByID(id uuid.UUID) (*models.Class, error)
	FindByIDWithLevels(id uuid.UUID) (*models.Class, error)
	FindByName(name string) (*models.Class, error)
	FindLevelByClassAndLevel(classID uuid.UUID, level int) (*models.ClassLevel, error)
	FindEquipmentOptionsByIDs(ids []uuid.UUID) ([]models.ClassStartingEquipmentOption, error)
	FindWeaponRequirementByClassAndLevel(classID uuid.UUID, level int) (*models.ClassWeaponRequirement, error)
	FindResourceDefinitionsByClassID(classID uuid.UUID) ([]models.ClassResourceDefinition, error)
	FindLevelResourcesByClassLevel(classLevelID uuid.UUID) ([]models.ClassLevelResource, error)
	FindLevelResourcesByClassAndLevel(classID uuid.UUID, level int) ([]models.ClassLevelResource, error)
	GetLevelResourceMap(classID uuid.UUID, level int) (map[string]int, error)
	GetLevelResourceValue(classID uuid.UUID, level int, key string) int
	FindResourceInfo(classID uuid.UUID, level int) ([]models.ClassResourceDefinition, map[string]int, error)
	Add(class *models.Class) error
	AddLevel(level *models.ClassLevel) error
}

type ClassRepo struct {
	db *gorm.DB
}

func NewClassRepo(db *gorm.DB) *ClassRepo {
	return &ClassRepo{db}
}

func (r *ClassRepo) FindAll() ([]*models.Class, error) {
	var classes []*models.Class
	err := r.db.Preload("Components").Preload("Levels", func(db *gorm.DB) *gorm.DB {
		return db.Where("level = ?", 1)
	}).Preload("Levels.LevelFeatures", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Preload("EquipmentChoices", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Preload("EquipmentChoices.Options").Find(&classes).Error
	return classes, err
}

func (r *ClassRepo) FindByID(id uuid.UUID) (*models.Class, error) {
	var c models.Class
	if err := r.db.First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClassRepo) FindByIDWithLevels(id uuid.UUID) (*models.Class, error) {
	var c models.Class
	if err := r.db.Preload("Levels", func(db *gorm.DB) *gorm.DB {
		return db.Order("level ASC")
	}).Preload("Levels.LevelFeatures", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Preload("Components").First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClassRepo) FindByName(name string) (*models.Class, error) {
	var c models.Class
	if err := r.db.First(&c, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClassRepo) FindLevelByClassAndLevel(classID uuid.UUID, level int) (*models.ClassLevel, error) {
	var cl models.ClassLevel
	if err := r.db.Preload("LevelFeatures", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&cl, "class_id = ? AND level = ?", classID, level).Error; err != nil {
		return nil, err
	}
	return &cl, nil
}

func (r *ClassRepo) FindEquipmentOptionsByIDs(ids []uuid.UUID) ([]models.ClassStartingEquipmentOption, error) {
	var options []models.ClassStartingEquipmentOption
	if len(ids) == 0 {
		return options, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&options).Error
	return options, err
}

// FindResourceInfo returns resource definitions and a map of resource values for a class at a given level.
// This consolidates multiple lookups into a single call.
func (r *ClassRepo) FindResourceInfo(classID uuid.UUID, level int) ([]models.ClassResourceDefinition, map[string]int, error) {
	var class models.Class
	err := r.db.Preload("ResourceDefinitions", func(db *gorm.DB) *gorm.DB {
		return db.Order("display_order ASC")
	}).Preload("Levels", "level = ?", level).
		Preload("Levels.ResourceValues").
		First(&class, "id = ?", classID).Error

	if err != nil {
		return nil, nil, err
	}

	resourceMap := make(map[string]int)
	if len(class.Levels) > 0 {
		for _, rv := range class.Levels[0].ResourceValues {
			resourceMap[rv.ResourceKey] = rv.Value
		}
	}

	return class.ResourceDefinitions, resourceMap, nil
}

func (r *ClassRepo) Add(class *models.Class) error {
	return r.db.Create(class).Error
}

func (r *ClassRepo) AddLevel(level *models.ClassLevel) error {
	return r.db.Create(level).Error
}

// FindWeaponRequirementByClassAndLevel returns a weapon requirement if one exists for the given class and level
func (r *ClassRepo) FindWeaponRequirementByClassAndLevel(classID uuid.UUID, level int) (*models.ClassWeaponRequirement, error) {
	var req models.ClassWeaponRequirement
	if err := r.db.First(&req, "class_id = ? AND selection_level = ?", classID, level).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

// FindResourceDefinitionsByClassID returns all resource definitions for a class, ordered by display_order
func (r *ClassRepo) FindResourceDefinitionsByClassID(classID uuid.UUID) ([]models.ClassResourceDefinition, error) {
	var defs []models.ClassResourceDefinition
	err := r.db.Where("class_id = ?", classID).Order("display_order ASC").Find(&defs).Error
	return defs, err
}

// FindLevelResourcesByClassLevel returns all resource values for a given class level ID
func (r *ClassRepo) FindLevelResourcesByClassLevel(classLevelID uuid.UUID) ([]models.ClassLevelResource, error) {
	var resources []models.ClassLevelResource
	err := r.db.Where("class_level_id = ?", classLevelID).Find(&resources).Error
	return resources, err
}

// FindLevelResourcesByClassAndLevel returns all resource values for a class at a given level
func (r *ClassRepo) FindLevelResourcesByClassAndLevel(classID uuid.UUID, level int) ([]models.ClassLevelResource, error) {
	var resources []models.ClassLevelResource
	err := r.db.Raw(`
		SELECT clr.* FROM class_level_resources clr
		JOIN class_levels cl ON cl.id = clr.class_level_id
		WHERE cl.class_id = ? AND cl.level = ?
	`, classID, level).Scan(&resources).Error
	return resources, err
}

// GetLevelResourceMap returns a map of resource_key -> value for a class at a given level.
// This is a convenience method for services that need to look up specific resource values.
func (r *ClassRepo) GetLevelResourceMap(classID uuid.UUID, level int) (map[string]int, error) {
	resources, err := r.FindLevelResourcesByClassAndLevel(classID, level)
	if err != nil {
		return nil, err
	}
	m := make(map[string]int, len(resources))
	for _, r := range resources {
		m[r.ResourceKey] = r.Value
	}
	return m, nil
}

// GetLevelResourceValue returns a single resource value for a class at a given level.
// Returns 0 if the resource is not found.
func (r *ClassRepo) GetLevelResourceValue(classID uuid.UUID, level int, key string) int {
	m, err := r.GetLevelResourceMap(classID, level)
	if err != nil {
		return 0
	}
	return m[key]
}
