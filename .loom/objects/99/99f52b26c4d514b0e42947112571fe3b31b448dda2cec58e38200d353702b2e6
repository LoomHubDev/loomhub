package store

import (
	"database/sql"
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type OwnerStore struct {
	db *sql.DB
}

func NewOwnerStore(db *sql.DB) *OwnerStore {
	return &OwnerStore{db: db}
}

func (s *OwnerStore) Create(owner *models.Owner) error {
	_, err := s.db.Exec(
		"INSERT INTO owners (id, name, type, created_at) VALUES (?, ?, ?, ?)",
		owner.ID, owner.Name, owner.Type, owner.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create owner: %w", err)
	}
	return nil
}

func (s *OwnerStore) GetByName(name string) (*models.Owner, error) {
	var o models.Owner
	err := s.db.QueryRow(
		"SELECT id, name, type, created_at FROM owners WHERE name = ?", name,
	).Scan(&o.ID, &o.Name, &o.Type, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *OwnerStore) NameExists(name string) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM owners WHERE name = ?", name).Scan(&count)
	return count > 0, err
}
