package store

import (
	"fmt"
	"time"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CollaboratorStore struct {
	db *gorm.DB
}

func NewCollaboratorStore(db *gorm.DB) *CollaboratorStore {
	return &CollaboratorStore{db: db}
}

func (s *CollaboratorStore) Add(loomID, userID, permission string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	c := models.Collaborator{
		LoomID: loomID, UserID: userID, Permission: permission, CreatedAt: now,
	}
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "loom_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"permission"}),
	}).Create(&c).Error
	if err != nil {
		return fmt.Errorf("add collaborator: %w", err)
	}
	return nil
}

func (s *CollaboratorStore) Remove(loomID, userID string) error {
	return s.db.Delete(&models.Collaborator{}, "loom_id = ? AND user_id = ?", loomID, userID).Error
}

func (s *CollaboratorStore) List(loomID string) ([]models.CollaboratorWithUser, error) {
	var collabs []models.Collaborator
	err := s.db.Preload("User").Where("loom_id = ?", loomID).Order("created_at ASC").Find(&collabs).Error
	if err != nil {
		return nil, err
	}

	var result []models.CollaboratorWithUser
	for _, c := range collabs {
		cwu := models.CollaboratorWithUser{
			UserID:     c.UserID,
			Permission: c.Permission,
			CreatedAt:  c.CreatedAt,
		}
		if c.User != nil {
			cwu.Username = c.User.Username
			cwu.DisplayName = c.User.DisplayName
			cwu.AvatarURL = c.User.AvatarURL
		}
		result = append(result, cwu)
	}
	return result, nil
}

func (s *CollaboratorStore) GetPermission(loomID, userID string) (string, error) {
	var c models.Collaborator
	err := s.db.Select("permission").Where("loom_id = ? AND user_id = ?", loomID, userID).First(&c).Error
	if err != nil {
		return "", err
	}
	return c.Permission, nil
}

func (s *CollaboratorStore) CheckAccess(loom *models.Loom, userID string, isAdmin bool) string {
	if isAdmin {
		return "admin"
	}
	if loom.OwnerID == userID {
		return "admin"
	}
	perm, err := s.GetPermission(loom.ID, userID)
	if err == nil {
		return perm
	}
	if loom.Visibility == "public" {
		return "read"
	}
	return ""
}
