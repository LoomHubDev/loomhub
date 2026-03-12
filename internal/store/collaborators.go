package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type CollaboratorStore struct {
	db *sql.DB
}

func NewCollaboratorStore(db *sql.DB) *CollaboratorStore {
	return &CollaboratorStore{db: db}
}

func (s *CollaboratorStore) Add(loomID, userID, permission string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(`
		INSERT INTO loom_collaborators (loom_id, user_id, permission, created_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(loom_id, user_id) DO UPDATE SET permission = excluded.permission`,
		loomID, userID, permission, now,
	)
	if err != nil {
		return fmt.Errorf("add collaborator: %w", err)
	}
	return nil
}

func (s *CollaboratorStore) Remove(loomID, userID string) error {
	_, err := s.db.Exec("DELETE FROM loom_collaborators WHERE loom_id = ? AND user_id = ?", loomID, userID)
	return err
}

func (s *CollaboratorStore) List(loomID string) ([]models.CollaboratorWithUser, error) {
	rows, err := s.db.Query(`
		SELECT c.user_id, c.permission, c.created_at,
			   u.username, u.display_name, u.avatar_url
		FROM loom_collaborators c
		JOIN users u ON u.id = c.user_id
		WHERE c.loom_id = ?
		ORDER BY c.created_at`, loomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collabs []models.CollaboratorWithUser
	for rows.Next() {
		var c models.CollaboratorWithUser
		if err := rows.Scan(&c.UserID, &c.Permission, &c.CreatedAt,
			&c.Username, &c.DisplayName, &c.AvatarURL); err != nil {
			return nil, err
		}
		collabs = append(collabs, c)
	}
	return collabs, nil
}

func (s *CollaboratorStore) GetPermission(loomID, userID string) (string, error) {
	var perm string
	err := s.db.QueryRow(
		"SELECT permission FROM loom_collaborators WHERE loom_id = ? AND user_id = ?",
		loomID, userID,
	).Scan(&perm)
	if err != nil {
		return "", err
	}
	return perm, nil
}

// CheckAccess returns the effective permission for a user on a loom.
// Returns "admin" for owners, the collaborator permission, or "" for no access.
func (s *CollaboratorStore) CheckAccess(loom *models.Loom, userID string, isAdmin bool) string {
	// Site admins have admin access to everything
	if isAdmin {
		return "admin"
	}
	// Owner has admin access
	if loom.OwnerID == userID {
		return "admin"
	}
	// Check collaborator table
	perm, err := s.GetPermission(loom.ID, userID)
	if err == nil {
		return perm
	}
	// Public looms give read access to everyone
	if loom.Visibility == "public" {
		return "read"
	}
	return ""
}
