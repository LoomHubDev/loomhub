package store

import (
	"time"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PinStore struct {
	db *gorm.DB
}

func NewPinStore(db *gorm.DB) *PinStore {
	return &PinStore{db: db}
}

func (s *PinStore) Pin(userID, loomID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC().Format(time.RFC3339)
		p := models.Pin{UserID: userID, LoomID: loomID, CreatedAt: now}
		tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&p)
		return s.updatePinCount(tx, loomID)
	})
}

func (s *PinStore) Unpin(userID, loomID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Delete(&models.Pin{}, "user_id = ? AND loom_id = ?", userID, loomID)
		return s.updatePinCount(tx, loomID)
	})
}

func (s *PinStore) updatePinCount(tx *gorm.DB, loomID string) error {
	var count int64
	tx.Model(&models.Pin{}).Where("loom_id = ?", loomID).Count(&count)
	return tx.Model(&models.Loom{}).Where("id = ?", loomID).Update("pin_count", count).Error
}

func (s *PinStore) IsPinned(userID, loomID string) (bool, error) {
	var count int64
	err := s.db.Model(&models.Pin{}).Where("user_id = ? AND loom_id = ?", userID, loomID).Count(&count).Error
	return count > 0, err
}
