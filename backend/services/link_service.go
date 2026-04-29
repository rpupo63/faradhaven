package services

import (
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
)

type LinkService interface {
	CreateLink(sourceCharacterID, targetCharacterID uuid.UUID, linkType models.LinkType, notes *string) (*models.CharacterLink, error)
	RemoveLink(linkID uuid.UUID) error
	GetLinks(characterID uuid.UUID) ([]*models.CharacterLink, error)
}

type linkService struct {
	repo database.CharacterLinkRepository
}

func NewLinkService(repo database.CharacterLinkRepository) LinkService {
	return &linkService{repo}
}

func (s *linkService) CreateLink(sourceCharacterID, targetCharacterID uuid.UUID, linkType models.LinkType, notes *string) (*models.CharacterLink, error) {
	link := &models.CharacterLink{
		SourceCharacterID: sourceCharacterID,
		TargetCharacterID: targetCharacterID,
		LinkType:          linkType,
		IsActive:          true,
		Notes:             notes,
	}

	err := s.repo.Add(link)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (s *linkService) RemoveLink(linkID uuid.UUID) error {
	return s.repo.Delete(linkID)
}

func (s *linkService) GetLinks(characterID uuid.UUID) ([]*models.CharacterLink, error) {
	return s.repo.FindByCharacterID(characterID)
}
