package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

// InitiativeEntry pairs a token ID with its desired initiative order.
type InitiativeEntry struct {
	TokenID uuid.UUID
	Order   int
}

type MapTokenRepo struct {
	db *gorm.DB
}

func NewMapTokenRepo(db *gorm.DB) *MapTokenRepo {
	return &MapTokenRepo{db: db}
}

func (r *MapTokenRepo) Create(token *models.MapToken) error {
	return r.db.Create(token).Error
}

func (r *MapTokenRepo) GetByID(id uuid.UUID) (*models.MapToken, error) {
	var token models.MapToken
	if err := r.db.First(&token, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *MapTokenRepo) Update(token *models.MapToken) error {
	return r.db.Save(token).Error
}

func (r *MapTokenRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.MapToken{}, "id = ?", id).Error
}

func (r *MapTokenRepo) GetByInitiativeOrder(mapID uuid.UUID) ([]models.MapToken, error) {
	var tokens []models.MapToken
	err := r.db.Where("map_id = ? AND initiative_order IS NOT NULL", mapID).
		Order("initiative_order ASC").Find(&tokens).Error
	return tokens, err
}

func (r *MapTokenRepo) SetInitiativeOrder(tokenID uuid.UUID, order *int) error {
	return r.db.Model(&models.MapToken{}).Where("id = ?", tokenID).
		Update("initiative_order", order).Error
}

// BulkSetInitiativeOrder updates initiative_order for multiple tokens in a single query.
func (r *MapTokenRepo) BulkSetInitiativeOrder(entries []InitiativeEntry) error {
	if len(entries) == 0 {
		return nil
	}

	idsStr := "{"
	ordersStr := "{"
	for i, e := range entries {
		if i > 0 {
			idsStr += ","
			ordersStr += ","
		}
		idsStr += e.TokenID.String()
		ordersStr += fmt.Sprintf("%d", e.Order)
	}
	idsStr += "}"
	ordersStr += "}"

	return r.db.Exec(`
		UPDATE map_tokens SET initiative_order = v.ord
		FROM (SELECT unnest(?::uuid[]) AS id, unnest(?::int[]) AS ord) v
		WHERE map_tokens.id = v.id
	`, idsStr, ordersStr).Error
}

// ClearInitiativeByMapID sets initiative_order to NULL for all tokens on a map in one query.
func (r *MapTokenRepo) ClearInitiativeByMapID(mapID uuid.UUID) error {
	return r.db.Model(&models.MapToken{}).Where("map_id = ?", mapID).
		Update("initiative_order", nil).Error
}
