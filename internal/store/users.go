package store

import (
	"database/sql"
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(user *models.User) error {
	_, err := s.db.Exec(`
		INSERT INTO users (id, username, email, display_name, bio, avatar_url, password_hash, is_admin, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Username, user.Email, user.DisplayName, user.Bio, user.AvatarURL,
		user.PasswordHash, user.IsAdmin, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (s *UserStore) GetByID(id string) (*models.User, error) {
	return s.scanUser(s.db.QueryRow(`
		SELECT id, username, email, display_name, bio, avatar_url, password_hash, is_admin, created_at, updated_at
		FROM users WHERE id = ?`, id,
	))
}

func (s *UserStore) GetByUsername(username string) (*models.User, error) {
	return s.scanUser(s.db.QueryRow(`
		SELECT id, username, email, display_name, bio, avatar_url, password_hash, is_admin, created_at, updated_at
		FROM users WHERE username = ?`, username,
	))
}

func (s *UserStore) GetByEmail(email string) (*models.User, error) {
	return s.scanUser(s.db.QueryRow(`
		SELECT id, username, email, display_name, bio, avatar_url, password_hash, is_admin, created_at, updated_at
		FROM users WHERE email = ?`, email,
	))
}

func (s *UserStore) Update(user *models.User) error {
	_, err := s.db.Exec(`
		UPDATE users SET display_name = ?, bio = ?, avatar_url = ?, updated_at = ?
		WHERE id = ?`,
		user.DisplayName, user.Bio, user.AvatarURL, user.UpdatedAt, user.ID,
	)
	return err
}

func (s *UserStore) Count() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

func (s *UserStore) scanUser(row *sql.Row) (*models.User, error) {
	var u models.User
	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.DisplayName, &u.Bio, &u.AvatarURL,
		&u.PasswordHash, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return &u, nil
}
