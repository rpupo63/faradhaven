package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	FindByIDWithAllRelations(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByRefreshToken(token string) (*models.User, error)
	Add(user *models.User) error
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

// FindAll returns all users
func (r *UserRepo) FindAll() ([]*models.User, error) {
	var users []*models.User
	err := r.db.Find(&users).Error
	return users, err
}

// FindByID returns a user by ID
func (r *UserRepo) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByIDWithAllRelations returns a user with all nested relations loaded
func (r *UserRepo) FindByIDWithAllRelations(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.
		Preload("Characters").
		Preload("Spells").
		Preload("Beasts").
		Preload("Beasts.Attacks").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail returns a user by email
func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByRefreshToken returns a user by their refresh token
func (r *UserRepo) FindByRefreshToken(token string) (*models.User, error) {
	if token == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var user models.User
	err := r.db.First(&user, "refresh_token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Add inserts a new user
func (r *UserRepo) Add(user *models.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	return r.db.Create(user).Error
}

// Update updates an existing user
func (r *UserRepo) Update(user *models.User) error {
	user.UpdatedAt = time.Now()
	return r.db.Save(user).Error
}

// Delete removes a user by ID
func (r *UserRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}
