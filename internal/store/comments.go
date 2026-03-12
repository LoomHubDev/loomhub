package store

import (
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

type CommentStore struct {
	db *gorm.DB
}

func NewCommentStore(db *gorm.DB) *CommentStore {
	return &CommentStore{db: db}
}

func (s *CommentStore) Create(c *models.Comment) error {
	if err := s.db.Create(c).Error; err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

func (s *CommentStore) ListByWR(wrID string) ([]models.Comment, error) {
	var comments []models.Comment
	err := s.db.Preload("Author").
		Where("wr_id = ?", wrID).
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	for i := range comments {
		if comments[i].Author != nil {
			comments[i].AuthorName = comments[i].Author.Username
		}
	}
	return comments, nil
}

func (s *CommentStore) Update(id, body, updatedAt string) error {
	return s.db.Model(&models.Comment{}).Where("id = ?", id).Updates(map[string]any{
		"body": body, "updated_at": updatedAt,
	}).Error
}

func (s *CommentStore) Delete(id string) error {
	return s.db.Delete(&models.Comment{}, "id = ?", id).Error
}

func (s *CommentStore) GetByID(id string) (*models.Comment, error) {
	var c models.Comment
	if err := s.db.First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
