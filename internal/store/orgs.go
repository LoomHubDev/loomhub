package store

import (
	"database/sql"
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type OrgStore struct {
	db *sql.DB
}

func NewOrgStore(db *sql.DB) *OrgStore {
	return &OrgStore{db: db}
}

func (s *OrgStore) Create(org *models.Organization) error {
	_, err := s.db.Exec(`
		INSERT INTO organizations (id, name, display_name, description, avatar_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		org.ID, org.Name, org.DisplayName, org.Description, org.AvatarURL, org.CreatedAt, org.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create org: %w", err)
	}
	return nil
}

func (s *OrgStore) GetByName(name string) (*models.Organization, error) {
	var org models.Organization
	err := s.db.QueryRow(`
		SELECT id, name, display_name, description, avatar_url, created_at, updated_at
		FROM organizations WHERE name = ?`, name,
	).Scan(&org.ID, &org.Name, &org.DisplayName, &org.Description, &org.AvatarURL, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (s *OrgStore) Update(org *models.Organization) error {
	_, err := s.db.Exec(`
		UPDATE organizations SET display_name = ?, description = ?, avatar_url = ?, updated_at = ?
		WHERE id = ?`,
		org.DisplayName, org.Description, org.AvatarURL, org.UpdatedAt, org.ID,
	)
	return err
}

func (s *OrgStore) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM organizations WHERE id = ?", id)
	return err
}

func (s *OrgStore) AddMember(orgID, userID, role string) error {
	_, err := s.db.Exec(`
		INSERT INTO org_members (org_id, user_id, role, created_at) VALUES (?, ?, ?, datetime('now'))
		ON CONFLICT(org_id, user_id) DO UPDATE SET role = excluded.role`,
		orgID, userID, role,
	)
	return err
}

func (s *OrgStore) RemoveMember(orgID, userID string) error {
	_, err := s.db.Exec("DELETE FROM org_members WHERE org_id = ? AND user_id = ?", orgID, userID)
	return err
}

func (s *OrgStore) GetMemberRole(orgID, userID string) (string, error) {
	var role string
	err := s.db.QueryRow("SELECT role FROM org_members WHERE org_id = ? AND user_id = ?", orgID, userID).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

func (s *OrgStore) ListMembers(orgID string) ([]models.OrgMember, error) {
	rows, err := s.db.Query(`
		SELECT om.org_id, om.user_id, om.role, om.created_at, u.username
		FROM org_members om
		JOIN users u ON u.id = om.user_id
		WHERE om.org_id = ?
		ORDER BY om.created_at ASC`, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.OrgMember
	for rows.Next() {
		var m models.OrgMember
		if err := rows.Scan(&m.OrgID, &m.UserID, &m.Role, &m.CreatedAt, &m.Username); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (s *OrgStore) ListByUser(userID string) ([]models.Organization, error) {
	rows, err := s.db.Query(`
		SELECT o.id, o.name, o.display_name, o.description, o.avatar_url, o.created_at, o.updated_at
		FROM organizations o
		JOIN org_members om ON om.org_id = o.id
		WHERE om.user_id = ?
		ORDER BY o.name ASC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		var o models.Organization
		if err := rows.Scan(&o.ID, &o.Name, &o.DisplayName, &o.Description, &o.AvatarURL, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, nil
}
