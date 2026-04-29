package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type CharacterLinkRepository interface {
	FindByCharacterID(characterID uuid.UUID) ([]*models.CharacterLink, error)
	FindBySourceCharacter(sourceID uuid.UUID) ([]*models.CharacterLink, error)
	FindByTargetCharacter(targetID uuid.UUID) ([]*models.CharacterLink, error)
	FindByID(id uuid.UUID) (*models.CharacterLink, error)
	Add(link *models.CharacterLink) error
	Update(link *models.CharacterLink) error
	Delete(id uuid.UUID) error
}

type CharacterLinkRepo struct {
	db *gorm.DB
}

func NewCharacterLinkRepo(db *gorm.DB) *CharacterLinkRepo {
	return &CharacterLinkRepo{db}
}

// FindByCharacterID returns all links involving a character (as source or target)
func (r *CharacterLinkRepo) FindByCharacterID(characterID uuid.UUID) ([]*models.CharacterLink, error) {
	var links []*models.CharacterLink
	err := r.db.Where("source_character_id = ? OR target_character_id = ?", characterID, characterID).
		Order("created_at desc").Find(&links).Error
	return links, err
}

// FindBySourceCharacter returns all links where the character is the source
func (r *CharacterLinkRepo) FindBySourceCharacter(sourceID uuid.UUID) ([]*models.CharacterLink, error) {
	var links []*models.CharacterLink
	err := r.db.Where("source_character_id = ?", sourceID).Order("created_at desc").Find(&links).Error
	return links, err
}

// FindByTargetCharacter returns all links where the character is the target
func (r *CharacterLinkRepo) FindByTargetCharacter(targetID uuid.UUID) ([]*models.CharacterLink, error) {
	var links []*models.CharacterLink
	err := r.db.Where("target_character_id = ?", targetID).Order("created_at desc").Find(&links).Error
	return links, err
}

// FindActiveByType returns active links of a specific type for a character
func (r *CharacterLinkRepo) FindActiveByType(characterID uuid.UUID, linkType models.LinkType) ([]*models.CharacterLink, error) {
	var links []*models.CharacterLink
	err := r.db.Where("(source_character_id = ? OR target_character_id = ?) AND link_type = ? AND is_active = true",
		characterID, characterID, linkType).Find(&links).Error
	return links, err
}

// FindByID returns a single link by ID
func (r *CharacterLinkRepo) FindByID(id uuid.UUID) (*models.CharacterLink, error) {
	var link models.CharacterLink
	err := r.db.First(&link, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// FindExisting checks if a link already exists between two characters
func (r *CharacterLinkRepo) FindExisting(sourceID, targetID uuid.UUID, linkType models.LinkType) (*models.CharacterLink, error) {
	var link models.CharacterLink
	err := r.db.Where("source_character_id = ? AND target_character_id = ? AND link_type = ?",
		sourceID, targetID, linkType).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// Add creates a new character link
func (r *CharacterLinkRepo) Add(link *models.CharacterLink) error {
	now := time.Now()
	link.CreatedAt = now
	link.UpdatedAt = now
	return r.db.Create(link).Error
}

// Update updates an existing character link
func (r *CharacterLinkRepo) Update(link *models.CharacterLink) error {
	link.UpdatedAt = time.Now()
	return r.db.Save(link).Error
}

// Delete removes a character link by ID
func (r *CharacterLinkRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.CharacterLink{}, "id = ?", id).Error
}

// DeleteExpired removes all expired links
func (r *CharacterLinkRepo) DeleteExpired() (int64, error) {
	result := r.db.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Delete(&models.CharacterLink{})
	return result.RowsAffected, result.Error
}

// Deactivate sets a link as inactive
func (r *CharacterLinkRepo) Deactivate(id uuid.UUID) error {
	return r.db.Model(&models.CharacterLink{}).Where("id = ?", id).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()}).Error
}
