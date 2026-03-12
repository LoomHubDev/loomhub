package store

import (
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

type ReviewStore struct {
	db *gorm.DB
}

func NewReviewStore(db *gorm.DB) *ReviewStore {
	return &ReviewStore{db: db}
}

func (s *ReviewStore) Create(r *models.Review) error {
	if err := s.db.Create(r).Error; err != nil {
		return fmt.Errorf("create review: %w", err)
	}
	return nil
}

func (s *ReviewStore) ListByWR(wrID string) ([]models.Review, error) {
	var reviews []models.Review
	err := s.db.Preload("Reviewer").
		Where("wr_id = ?", wrID).
		Order("created_at ASC").
		Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	for i := range reviews {
		if reviews[i].Reviewer != nil {
			reviews[i].ReviewerName = reviews[i].Reviewer.Username
		}
	}
	return reviews, nil
}
