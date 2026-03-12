package store

import (
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

type LoomStore struct {
	db *gorm.DB
}

func NewLoomStore(db *gorm.DB) *LoomStore {
	return &LoomStore{db: db}
}

func (s *LoomStore) Create(loom *models.Loom) error {
	if err := s.db.Create(loom).Error; err != nil {
		return fmt.Errorf("create loom: %w", err)
	}
	return nil
}

func (s *LoomStore) GetByOwnerAndName(ownerName, loomName string) (*models.Loom, error) {
	var owner models.Owner
	if err := s.db.Where("name = ?", ownerName).First(&owner).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	var l models.Loom
	if err := s.db.Where("owner_id = ? AND name = ?", owner.ID, loomName).First(&l).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	l.OwnerName = ownerName
	l.FullName = ownerName + "/" + l.Name
	return &l, nil
}

func (s *LoomStore) ListByOwner(ownerID string, limit, offset int) ([]models.Loom, int, error) {
	var total int64
	if err := s.db.Model(&models.Loom{}).Where("owner_id = ?", ownerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var looms []models.Loom
	err := s.db.Preload("Owner").
		Where("owner_id = ?", ownerID).
		Order("updated_at DESC").
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

func (s *LoomStore) ListPublicByOwner(ownerName string, limit, offset int) ([]models.Loom, int, error) {
	var owner models.Owner
	if err := s.db.Where("name = ?", ownerName).First(&owner).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	s.db.Model(&models.Loom{}).Where("owner_id = ? AND visibility = ?", owner.ID, "public").Count(&total)

	var looms []models.Loom
	err := s.db.Where("owner_id = ? AND visibility = ?", owner.ID, "public").
		Order("updated_at DESC").
		Limit(limit).Offset(offset).
		Find(&looms).Error
	if err != nil {
		return nil, 0, err
	}

	for i := range looms {
		looms[i].OwnerName = ownerName
		looms[i].FullName = ownerName + "/" + looms[i].Name
	}
	return looms, int(total), nil
}

func (s *LoomStore) Delete(id string) error {
	return s.db.Delete(&models.Loom{}, "id = ?", id).Error
}

func (s *LoomStore) Update(loom *models.Loom) error {
	return s.db.Model(loom).Updates(map[string]any{
		"description":    loom.Description,
		"visibility":     loom.Visibility,
		"default_stream": loom.DefaultStream,
		"updated_at":     loom.UpdatedAt,
	}).Error
}
