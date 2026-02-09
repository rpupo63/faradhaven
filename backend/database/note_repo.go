package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type NoteRepository interface {
	FindAll() ([]*models.SharedNote, error)
	FindByID(id uuid.UUID) (*models.SharedNote, error)
	Add(note *models.SharedNote) error
	Update(note *models.SharedNote) error
	Delete(id uuid.UUID) error
}

type NoteRepo struct {
	db *gorm.DB
}

func NewNoteRepo(db *gorm.DB) *NoteRepo {
	return &NoteRepo{db}
}

func (r *NoteRepo) FindAll() ([]*models.SharedNote, error) {
	var notes []*models.SharedNote
	err := r.db.Order("created_at desc").Find(&notes).Error
	return notes, err
}

func (r *NoteRepo) FindByID(id uuid.UUID) (*models.SharedNote, error) {
	var note models.SharedNote
	err := r.db.First(&note, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepo) Add(note *models.SharedNote) error {
	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now
	return r.db.Create(note).Error
}

func (r *NoteRepo) Update(note *models.SharedNote) error {
	note.UpdatedAt = time.Now()
	return r.db.Save(note).Error
}

func (r *NoteRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.SharedNote{}, "id = ?", id).Error
}
