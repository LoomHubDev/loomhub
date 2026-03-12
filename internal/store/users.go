package store

import (
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(user *models.User) error {
	if err := s.db.Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (s *UserStore) GetByID(id string) (*models.User, error) {
	var u models.User
	if err := s.db.First(&u, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return &u, nil
}

func (s *UserStore) GetByUsername(username string) (*models.User, error) {
	var u models.User
	if err := s.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return &u, nil
}

func (s *UserStore) GetByEmail(email string) (*models.User, error) {
	var u models.User
	if err := s.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return &u, nil
}

func (s *UserStore) Update(user *models.User) error {
	return s.db.Model(user).Updates(map[string]any{
		"display_name": user.DisplayName,
		"bio":          user.Bio,
		"avatar_url":   user.AvatarURL,
		"updated_at":   user.UpdatedAt,
	}).Error
}

func (s *UserStore) Count() (int, error) {
	var count int64
	err := s.db.Model(&models.User{}).Count(&count).Error
	return int(count), err
}
