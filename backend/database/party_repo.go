package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PartyRepository defines the interface for interacting with Party data.
type PartyRepository interface {
	Add(party *models.Party) error
	FindByID(id uuid.UUID) (*models.Party, error)
	FindByIDWithMembers(id uuid.UUID) (*models.Party, error)
	FindByIDWithIdentifiedBeasts(id uuid.UUID) (*models.Party, error)
	FindByIDWithMembersAndBeasts(id uuid.UUID) (*models.Party, error)
	FindByOwnerID(ownerID uuid.UUID) ([]*models.Party, error)
	Update(party *models.Party) error
	Delete(id uuid.UUID) error
	AddMember(partyID, characterID uuid.UUID) error
	RemoveMember(partyID, characterID uuid.UUID) error
	AddIdentifiedBeast(partyID, beastID uuid.UUID) error
	RemoveIdentifiedBeast(partyID, beastID uuid.UUID) error
	GetMembers(partyID uuid.UUID) ([]*models.Character, error)
	GetIdentifiedBeasts(partyID uuid.UUID) ([]*models.Beast, error)
}

// PartyRepo implements PartyRepository using GORM.
type PartyRepo struct {
	db *gorm.DB
}

// NewPartyRepo creates a new PartyRepo.
func NewPartyRepo(db *gorm.DB) *PartyRepo {
	return &PartyRepo{db: db}
}

// Add creates a new Party record.
func (r *PartyRepo) Add(party *models.Party) error {
	party.CreatedAt = time.Now()
	party.UpdatedAt = time.Now()
	return r.db.Create(party).Error
}

// FindByID retrieves a Party by its ID.
func (r *PartyRepo) FindByID(id uuid.UUID) (*models.Party, error) {
	var party models.Party
	err := r.db.Preload("Owner").First(&party, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &party, nil
}

// FindByIDWithMembers retrieves a Party by its ID, preloading its members.
func (r *PartyRepo) FindByIDWithMembers(id uuid.UUID) (*models.Party, error) {
	var party models.Party
	err := r.db.Preload("Owner").Preload("Members").First(&party, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &party, nil
}

// FindByIDWithIdentifiedBeasts retrieves a Party by its ID, preloading its identified beasts.
func (r *PartyRepo) FindByIDWithIdentifiedBeasts(id uuid.UUID) (*models.Party, error) {
	var party models.Party
	err := r.db.Preload("Owner").Preload("IdentifiedBeasts").First(&party, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &party, nil
}

// FindByIDWithMembersAndBeasts retrieves a Party by its ID, preloading members and identified beasts.
func (r *PartyRepo) FindByIDWithMembersAndBeasts(id uuid.UUID) (*models.Party, error) {
	var party models.Party
	err := r.db.Preload("Owner").
		Preload("Members").
		Preload("IdentifiedBeasts").
		First(&party, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &party, nil
}

// FindByOwnerID retrieves all Parties owned by a specific user.
func (r *PartyRepo) FindByOwnerID(ownerID uuid.UUID) ([]*models.Party, error) {
	var parties []*models.Party
	err := r.db.Where("owner_id = ?", ownerID).Find(&parties).Error
	if err != nil {
		return nil, err
	}
	return parties, nil
}

// Update updates an existing Party record.
func (r *PartyRepo) Update(party *models.Party) error {
	party.UpdatedAt = time.Now()
	return r.db.Save(party).Error
}

// Delete removes a Party record by its ID.
func (r *PartyRepo) Delete(id uuid.UUID) error {
	return r.db.Select(clause.Associations).Delete(&models.Party{}, "id = ?", id).Error
}

// AddMember associates a character with a party.
func (r *PartyRepo) AddMember(partyID, characterID uuid.UUID) error {
	return r.db.Model(&models.Character{}).Where("id = ?", characterID).Update("party_id", partyID).Error
}

// RemoveMember disassociates a character from a party.
func (r *PartyRepo) RemoveMember(partyID, characterID uuid.UUID) error {
	// Set the character's party_id to NULL
	return r.db.Model(&models.Character{}).Where("id = ? AND party_id = ?", characterID, partyID).Update("party_id", nil).Error
}

// AddIdentifiedBeast adds a beast to a party's identified beasts.
func (r *PartyRepo) AddIdentifiedBeast(partyID, beastID uuid.UUID) error {
	return r.db.Exec("INSERT INTO party_beasts (party_id, beast_id) VALUES (?, ?) ON CONFLICT DO NOTHING", partyID, beastID).Error
}

// RemoveIdentifiedBeast removes a beast from a party's identified beasts.
func (r *PartyRepo) RemoveIdentifiedBeast(partyID, beastID uuid.UUID) error {
	return r.db.Exec("DELETE FROM party_beasts WHERE party_id = ? AND beast_id = ?", partyID, beastID).Error
}

// GetMembers retrieves all characters belonging to a party.
func (r *PartyRepo) GetMembers(partyID uuid.UUID) ([]*models.Character, error) {
	var members []*models.Character
	err := r.db.Where("party_id = ?", partyID).Find(&members).Error
	return members, err
}

// GetIdentifiedBeasts retrieves all beasts identified by a party.
func (r *PartyRepo) GetIdentifiedBeasts(partyID uuid.UUID) ([]*models.Beast, error) {
	var party models.Party
	if err := r.db.Preload("IdentifiedBeasts").First(&party, "id = ?", partyID).Error; err != nil {
		return nil, err
	}
	return party.IdentifiedBeasts, nil
}
