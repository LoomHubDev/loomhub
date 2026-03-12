package store

import (
	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

type SearchStore struct {
	db *gorm.DB
}

func NewSearchStore(db *gorm.DB) *SearchStore {
	return &SearchStore{db: db}
}

func (s *SearchStore) SearchLooms(query string, limit int) ([]models.Loom, error) {
	pattern := "%" + query + "%"
	var looms []models.Loom
	err := s.db.Preload("Owner").
		Where("visibility = ? AND (name LIKE ? OR description LIKE ?)", "public", pattern, pattern).
		Order("pin_count DESC, updated_at DESC").
		Limit(limit).
		Find(&looms).Error
	if err != nil {
		return nil, err
	}
	for i := range looms {
		looms[i].FillOwnerFields()
	}
	return looms, nil
}

func (s *SearchStore) SearchUsers(query string, limit int) ([]models.User, error) {
	pattern := "%" + query + "%"
	var users []models.User
	err := s.db.Where("username LIKE ? OR display_name LIKE ?", pattern, pattern).
		Order("username ASC").
		Limit(limit).
		Find(&users).Error
	return users, err
}

func (s *SearchStore) ExploreLooms(sort string, limit, offset int) ([]models.Loom, int, error) {
	var total int64
	if err := s.db.Model(&models.Loom{}).Where("visibility = ?", "public").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := "updated_at DESC"
	if sort == "pins" {
		orderClause = "pin_count DESC, updated_at DESC"
	} else if sort == "created" {
		orderClause = "created_at DESC"
	}

	var looms []models.Loom
	err := s.db.Preload("Owner").
		Where("visibility = ?", "public").
		Order(orderClause).
		Limit(limit).Offset(offset).
		Find(&looms).Error
	if err != nil {
		return nil, 0, err
	}
	for i := range looms {
		looms[i].FillOwnerFields()
	}
	return looms, int(total), nil
}
