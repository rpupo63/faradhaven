package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type SpellRepository interface {
	FindAll(page, limit, levelFilter int) ([]*models.Spell, int64, error)
	FindAllWithComponents() ([]*models.Spell, error)
	FindByID(id uuid.UUID) (*models.Spell, error)
	FindByUserID(userID uuid.UUID) ([]*models.Spell, error)
	FindByUserIDPaginated(userID uuid.UUID, page, limit, levelFilter int) ([]*models.Spell, int64, error)
	FindByCharacterID(characterID uuid.UUID, page, limit, levelFilter int) ([]*models.Spell, int64, error)
	FindUnchecked() ([]*models.Spell, error)
	FindByFingerprint(fingerprint string) ([]*models.Spell, error)
	Add(spell *models.Spell, componentIDs []uuid.UUID) error
	Update(spell *models.Spell) error
	UpdateFields(id uuid.UUID, fields map[string]interface{}) error
	ReplaceComponents(spellID uuid.UUID, componentIDs []uuid.UUID) error
	Delete(id uuid.UUID) error
}

type SpellRepo struct {
	db *gorm.DB
}

func NewSpellRepo(db *gorm.DB) *SpellRepo {
	return &SpellRepo{db}
}

func preloadSpellComponentLinks(db *gorm.DB) *gorm.DB {
	return db.Order("sort_order ASC")
}

func hydrateSpellComponents(spell *models.Spell) {
	if spell != nil {
		spell.HydrateComponentsFromLinks()
	}
}

func hydrateSpellsComponents(spells []*models.Spell) {
	models.HydrateComponentsFromLinksSlice(spells)
}

// FindAll returns all spells with components preloaded, with pagination and filtering
func (r *SpellRepo) FindAll(page, limit, levelFilter int) ([]*models.Spell, int64, error) {
	var spells []*models.Spell
	var totalCount int64

	query := r.db.Model(&models.Spell{}).Preload("ComponentLinks", preloadSpellComponentLinks).Preload("ComponentLinks.Component").Preload("Character")

	if levelFilter > 0 {
		query = query.Where("level = ?", levelFilter)
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Order("level ASC, name ASC").Offset(offset).Limit(limit).Find(&spells).Error
	if err != nil {
		return nil, totalCount, err
	}
	hydrateSpellsComponents(spells)
	return spells, totalCount, err
}

// FindAllWithComponents returns all spells with components preloaded, with no pagination limit.
// Use this when the full catalog is needed in memory (e.g. spellbook eligibility checks).
func (r *SpellRepo) FindAllWithComponents() ([]*models.Spell, error) {
	var spells []*models.Spell
	err := r.db.Preload("ComponentLinks", preloadSpellComponentLinks).Preload("ComponentLinks.Component").Order("level ASC, name ASC").Find(&spells).Error
	if err != nil {
		return nil, err
	}
	hydrateSpellsComponents(spells)
	return spells, err
}

// FindByID returns a spell by ID with components preloaded
func (r *SpellRepo) FindByID(id uuid.UUID) (*models.Spell, error) {
	var spell models.Spell
	err := r.db.Preload("ComponentLinks", preloadSpellComponentLinks).Preload("ComponentLinks.Component").First(&spell, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	hydrateSpellComponents(&spell)
	return &spell, nil
}

// FindByUserID returns all spells for a user with components preloaded
func (r *SpellRepo) FindByUserID(userID uuid.UUID) ([]*models.Spell, error) {
	var spells []*models.Spell
	err := r.db.Preload("ComponentLinks", preloadSpellComponentLinks).Preload("ComponentLinks.Component").Where("user_id = ?", userID).Find(&spells).Error
	if err != nil {
		return nil, err
	}
	hydrateSpellsComponents(spells)
	return spells, err
}

// FindByUserIDPaginated returns all spells for a user with components preloaded and pagination
func (r *SpellRepo) FindByUserIDPaginated(userID uuid.UUID, page, limit, levelFilter int) ([]*models.Spell, int64, error) {
	var spells []*models.Spell
	var totalCount int64

	query := r.db.Model(&models.Spell{}).Preload("ComponentLinks", preloadSpellComponentLinks).Preload("ComponentLinks.Component").Where("user_id = ?", userID)

	if levelFilter > 0 {
		query = query.Where("level = ?", levelFilter)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("level ASC, name ASC").Offset(offset).Limit(limit).Find(&spells).Error
	if err != nil {
		return nil, totalCount, err
	}
	hydrateSpellsComponents(spells)
	return spells, totalCount, err
}

// FindByCharacterID returns all spells prepared by a character with components preloaded, with pagination and filtering
func (r *SpellRepo) FindByCharacterID(characterID uuid.UUID, page, limit, levelFilter int) ([]*models.Spell, int64, error) {
	var spells []*models.Spell
	var totalCount int64

	query := r.db.Model(&models.Spell{}).Preload("ComponentLinks", preloadSpellComponentLinks).Preload("ComponentLinks.Component").Where("character_id = ?", characterID)

	if levelFilter > 0 {
		query = query.Where("level = ?", levelFilter)
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Order("level ASC, name ASC").Offset(offset).Limit(limit).Find(&spells).Error
	if err != nil {
		return nil, totalCount, err
	}
	hydrateSpellsComponents(spells)
	return spells, totalCount, err
}

// FindUnchecked returns all spells that have not been reviewed by the GM
func (r *SpellRepo) FindUnchecked() ([]*models.Spell, error) {
	var spells []*models.Spell
	err := r.db.Preload("ComponentLinks", preloadSpellComponentLinks).Preload("ComponentLinks.Component").Preload("Character").Where("checked = false").Order("created_at ASC").Find(&spells).Error
	if err != nil {
		return nil, err
	}
	hydrateSpellsComponents(spells)
	return spells, err
}

// FindByFingerprint returns all spells matching the given component fingerprint.
func (r *SpellRepo) FindByFingerprint(fingerprint string) ([]*models.Spell, error) {
	var spells []*models.Spell
	err := r.db.
		Preload("ComponentLinks", preloadSpellComponentLinks).
		Preload("ComponentLinks.Component").
		Where("component_fingerprint = ?", fingerprint).
		Find(&spells).Error
	if err != nil {
		return nil, err
	}
	hydrateSpellsComponents(spells)
	return spells, nil
}

// computeAndStoreFingerprint fetches component categories, computes the canonical
// fingerprint for the given ordered componentIDs, and writes it to the spell row.
// db may be r.db or a transaction handle (following the replaceComponentsTx pattern).
func (r *SpellRepo) computeAndStoreFingerprint(db *gorm.DB, spellID uuid.UUID, componentIDs []uuid.UUID) error {
	var comps []models.Component
	if err := db.Select("id, category").Where("id IN ?", componentIDs).Find(&comps).Error; err != nil {
		return err
	}
	compMap := make(map[uuid.UUID]models.ComponentCategory, len(comps))
	for _, c := range comps {
		compMap[c.ID] = c.Category
	}
	links := make([]models.SpellComponent, len(componentIDs))
	for i, id := range componentIDs {
		links[i] = models.SpellComponent{
			SortOrder:   i,
			ComponentID: id,
			Component:   models.Component{ID: id, Category: compMap[id]},
		}
	}
	fp := models.ComponentFingerprint(links)
	return db.Model(&models.Spell{}).Where("id = ?", spellID).Update("component_fingerprint", fp).Error
}

// Add inserts a new spell and links it to the given components
func (r *SpellRepo) Add(spell *models.Spell, componentIDs []uuid.UUID) error {
	now := time.Now()
	spell.CreatedAt = now
	spell.UpdatedAt = now
	if err := r.db.Create(spell).Error; err != nil {
		return err
	}
	if len(componentIDs) > 0 {
		if err := replaceComponentsTx(r.db, spell.ID, componentIDs); err != nil {
			return err
		}
	}
	return r.computeAndStoreFingerprint(r.db, spell.ID, componentIDs)
}

// ReplaceComponents removes existing spell-component links and creates new ones
func (r *SpellRepo) ReplaceComponents(spellID uuid.UUID, componentIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("spell_id = ?", spellID).Delete(&models.SpellComponent{}).Error; err != nil {
			return err
		}
		if len(componentIDs) > 0 {
			if err := replaceComponentsTx(tx, spellID, componentIDs); err != nil {
				return err
			}
		}
		return r.computeAndStoreFingerprint(tx, spellID, componentIDs)
	})
}

func replaceComponentsTx(tx *gorm.DB, spellID uuid.UUID, componentIDs []uuid.UUID) error {
	scs := make([]models.SpellComponent, len(componentIDs))
	now := time.Now()
	for i, compID := range componentIDs {
		scs[i] = models.SpellComponent{
			SpellID:     spellID,
			SortOrder:   i,
			ComponentID: compID,
			CreatedAt:   now,
		}
	}
	if len(scs) > 0 {
		return tx.Create(&scs).Error
	}
	return nil
}

// Update updates an existing spell
func (r *SpellRepo) Update(spell *models.Spell) error {
	spell.UpdatedAt = time.Now()
	return r.db.Save(spell).Error
}

// UpdateFields updates only the specified columns on a spell by ID.
// Use this instead of Update when you want to touch a subset of fields without
// triggering a full Save (which would cascade into associations).
func (r *SpellRepo) UpdateFields(id uuid.UUID, fields map[string]interface{}) error {
	fields["updated_at"] = time.Now()
	return r.db.Model(&models.Spell{}).Where("id = ?", id).Updates(fields).Error
}

// Delete removes a spell by ID
func (r *SpellRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Spell{}, "id = ?", id).Error
}
