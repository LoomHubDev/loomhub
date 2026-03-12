package store

import (
	"fmt"
	"time"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

type WeaveRequestStore struct {
	db *gorm.DB
}

func NewWeaveRequestStore(db *gorm.DB) *WeaveRequestStore {
	return &WeaveRequestStore{db: db}
}

func (s *WeaveRequestStore) Create(wr *models.WeaveRequest) error {
	var maxNum int
	s.db.Model(&models.WeaveRequest{}).
		Select("COALESCE(MAX(number), 0)").
		Where("loom_id = ?", wr.LoomID).
		Scan(&maxNum)
	wr.Number = maxNum + 1

	if err := s.db.Create(wr).Error; err != nil {
		return fmt.Errorf("create weave request: %w", err)
	}
	return nil
}

func (s *WeaveRequestStore) GetByNumber(loomID string, number int) (*models.WeaveRequest, error) {
	var wr models.WeaveRequest
	err := s.db.Preload("Author").
		Where("loom_id = ? AND number = ?", loomID, number).
		First(&wr).Error
	if err != nil {
		return nil, fmt.Errorf("get weave request: %w", err)
	}
	if wr.Author != nil {
		wr.AuthorName = wr.Author.Username
	}
	return &wr, nil
}

func (s *WeaveRequestStore) ListByLoom(loomID, status string, limit, offset int) ([]models.WeaveRequest, int, error) {
	q := s.db.Model(&models.WeaveRequest{}).Where("loom_id = ?", loomID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := s.db.Preload("Author").Where("loom_id = ?", loomID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var wrs []models.WeaveRequest
	err := query.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&wrs).Error
	if err != nil {
		return nil, 0, err
	}
	for i := range wrs {
		if wrs[i].Author != nil {
			wrs[i].AuthorName = wrs[i].Author.Username
		}
	}
	return wrs, int(total), nil
}

func (s *WeaveRequestStore) UpdateStatus(id, status string, wovenBy *string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	updates := map[string]any{"status": status, "updated_at": now}
	if status == "woven" {
		updates["woven_by"] = wovenBy
		updates["woven_at"] = now
	}
	if status == "closed" {
		updates["closed_at"] = now
	}
	return s.db.Model(&models.WeaveRequest{}).Where("id = ?", id).Updates(updates).Error
}

func (s *WeaveRequestStore) Update(wr *models.WeaveRequest) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return s.db.Model(&models.WeaveRequest{}).Where("id = ?", wr.ID).Updates(map[string]any{
		"title": wr.Title, "description": wr.Description, "updated_at": now,
	}).Error
}
