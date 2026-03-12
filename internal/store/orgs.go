package store

import (
	"fmt"
	"time"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrgStore struct {
	db *gorm.DB
}

func NewOrgStore(db *gorm.DB) *OrgStore {
	return &OrgStore{db: db}
}

func (s *OrgStore) Create(org *models.Organization) error {
	if err := s.db.Create(org).Error; err != nil {
		return fmt.Errorf("create org: %w", err)
	}
	return nil
}

func (s *OrgStore) GetByName(name string) (*models.Organization, error) {
	var org models.Organization
	if err := s.db.Where("name = ?", name).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (s *OrgStore) Update(org *models.Organization) error {
	return s.db.Model(org).Updates(map[string]any{
		"display_name": org.DisplayName,
		"description":  org.Description,
		"avatar_url":   org.AvatarURL,
		"updated_at":   org.UpdatedAt,
	}).Error
}

func (s *OrgStore) Delete(id string) error {
	return s.db.Delete(&models.Organization{}, "id = ?", id).Error
}

func (s *OrgStore) AddMember(orgID, userID, role string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	m := models.OrgMember{OrgID: orgID, UserID: userID, Role: role, CreatedAt: now}
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "org_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role"}),
	}).Create(&m).Error
}

func (s *OrgStore) RemoveMember(orgID, userID string) error {
	return s.db.Delete(&models.OrgMember{}, "org_id = ? AND user_id = ?", orgID, userID).Error
}

func (s *OrgStore) GetMemberRole(orgID, userID string) (string, error) {
	var m models.OrgMember
	err := s.db.Select("role").Where("org_id = ? AND user_id = ?", orgID, userID).First(&m).Error
	if err != nil {
		return "", err
	}
	return m.Role, nil
}

func (s *OrgStore) ListMembers(orgID string) ([]models.OrgMember, error) {
	var members []models.OrgMember
	err := s.db.Preload("User").Where("org_id = ?", orgID).Order("created_at ASC").Find(&members).Error
	if err != nil {
		return nil, err
	}
	for i := range members {
		if members[i].User != nil {
			members[i].Username = members[i].User.Username
		}
	}
	return members, err
}

func (s *OrgStore) ListByUser(userID string) ([]models.Organization, error) {
	var orgs []models.Organization
	subQuery := s.db.Model(&models.OrgMember{}).Select("org_id").Where("user_id = ?", userID)
	err := s.db.Where("id IN (?)", subQuery).Order("name ASC").Find(&orgs).Error
	return orgs, err
}
