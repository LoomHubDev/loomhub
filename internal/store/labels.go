package store

import (
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LabelStore struct {
	db *gorm.DB
}

func NewLabelStore(db *gorm.DB) *LabelStore {
	return &LabelStore{db: db}
}

func (s *LabelStore) Create(l *models.Label) error {
	if err := s.db.Create(l).Error; err != nil {
		return fmt.Errorf("create label: %w", err)
	}
	return nil
}

func (s *LabelStore) ListByLoom(loomID string) ([]models.Label, error) {
	var labels []models.Label
	err := s.db.Where("loom_id = ?", loomID).Order("name ASC").Find(&labels).Error
	return labels, err
}

func (s *LabelStore) Update(l *models.Label) error {
	return s.db.Model(l).Updates(map[string]any{
		"name": l.Name, "color": l.Color, "description": l.Description,
	}).Error
}

func (s *LabelStore) Delete(id string) error {
	return s.db.Delete(&models.Label{}, "id = ?", id).Error
}

func (s *LabelStore) GetByID(id string) (*models.Label, error) {
	var l models.Label
	if err := s.db.First(&l, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *LabelStore) AddToWR(wrID, labelID string) error {
	wl := models.WRLabel{WRID: wrID, LabelID: labelID}
	return s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&wl).Error
}

func (s *LabelStore) RemoveFromWR(wrID, labelID string) error {
	return s.db.Delete(&models.WRLabel{}, "wr_id = ? AND label_id = ?", wrID, labelID).Error
}

func (s *LabelStore) ListByWR(wrID string) ([]models.Label, error) {
	var labels []models.Label
	subQuery := s.db.Model(&models.WRLabel{}).Select("label_id").Where("wr_id = ?", wrID)
	err := s.db.Where("id IN (?)", subQuery).Order("name ASC").Find(&labels).Error
	return labels, err
}
