package store

import (
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

type OwnerStore struct {
	db *gorm.DB
}

func NewOwnerStore(db *gorm.DB) *OwnerStore {
	return &OwnerStore{db: db}
}

func (s *OwnerStore) Create(owner *models.Owner) error {
	if err := s.db.Create(owner).Error; err != nil {
		return fmt.Errorf("create owner: %w", err)
	}
	return nil
}

func (s *OwnerStore) GetByName(name string) (*models.Owner, error) {
	var o models.Owner
	if err := s.db.Where("name = ?", name).First(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *OwnerStore) NameExists(name string) (bool, error) {
	var count int64
	err := s.db.Model(&models.Owner{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
