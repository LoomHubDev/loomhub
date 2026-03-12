package store

import (
	"database/sql"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type TokenStore struct {
	db *sql.DB
}

func NewTokenStore(db *sql.DB) *TokenStore {
	return &TokenStore{db: db}
}

func (s *TokenStore) Create(t *models.AccessToken) error {
	_, err := s.db.Exec(`
		INSERT INTO access_tokens (id, user_id, name, token_hash, scopes, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.UserID, t.Name, t.TokenHash, t.Scopes, t.ExpiresAt, t.CreatedAt,
	)
	return err
}

func (s *TokenStore) ListByUser(userID string) ([]models.AccessToken, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, name, scopes, expires_at, last_used_at, created_at
		FROM access_tokens WHERE user_id = ?
		ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []models.AccessToken
	for rows.Next() {
		var t models.AccessToken
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Scopes, &t.ExpiresAt, &t.LastUsedAt, &t.CreatedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func (s *TokenStore) Delete(id, userID string) error {
	_, err := s.db.Exec("DELETE FROM access_tokens WHERE id = ? AND user_id = ?", id, userID)
	return err
}
