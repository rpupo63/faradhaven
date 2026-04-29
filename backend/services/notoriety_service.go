package services

import (
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
)

type NotorietyService interface {
	UpdateNotoriety(characterID uuid.UUID, mpChange int, brChange int) error
}

type notorietyService struct {
	repo database.CharacterRepository
}

func NewNotorietyService(repo database.CharacterRepository) NotorietyService {
	return &notorietyService{repo}
}

func (s *notorietyService) UpdateNotoriety(characterID uuid.UUID, mpChange int, brChange int) error {
	character, err := s.repo.FindByID(characterID)
	if err != nil {
		return err
	}

	character.SanguineMP += mpChange
	character.SanguineBR += brChange
	character.SanguineNotoriety = character.SanguineMP - character.SanguineBR

	return s.repo.Update(character)
}
