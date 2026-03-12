package store

import (
	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

type TokenStore struct {
	db *gorm.DB
}

func NewTokenStore(db *gorm.DB) *TokenStore {
	return &TokenStore{db: db}
}

func (s *TokenStore) Create(t *models.AccessToken) error {
	return s.db.Create(t).Error
}

func (s *TokenStore) ListByUser(userID string) ([]models.AccessToken, error) {
	var tokens []models.AccessToken
	err := s.db.Select("id, user_id, name, scopes, expires_at, last_used_at, created_at").
		Where("user_id = ?", userID).Order("created_at DESC").Find(&tokens).Error
	return tokens, err
}

func (s *TokenStore) Delete(id, userID string) error {
	return s.db.Delete(&models.AccessToken{}, "id = ? AND user_id = ?", id, userID).Error
}
