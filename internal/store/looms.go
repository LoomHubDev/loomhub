package store

import (
	"database/sql"
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type LoomStore struct {
	db *sql.DB
}

func NewLoomStore(db *sql.DB) *LoomStore {
	return &LoomStore{db: db}
}

func (s *LoomStore) Create(loom *models.Loom) error {
	_, err := s.db.Exec(`
		INSERT INTO looms (id, owner_id, name, description, visibility, default_stream, disk_path, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		loom.ID, loom.OwnerID, loom.Name, loom.Description, loom.Visibility,
		loom.DefaultStream, loom.DiskPath, loom.CreatedAt, loom.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create loom: %w", err)
	}
	return nil
}

func (s *LoomStore) GetByOwnerAndName(ownerName, loomName string) (*models.Loom, error) {
	var l models.Loom
	err := s.db.QueryRow(`
		SELECT l.id, l.owner_id, l.name, l.description, l.visibility, l.default_stream,
			   l.disk_path, l.size_bytes, l.checkpoint_count, l.stream_count,
			   l.pin_count, l.spin_count, l.spun_from, l.created_at, l.updated_at, l.synced_at,
			   o.name
		FROM looms l
		JOIN owners o ON o.id = l.owner_id
		WHERE o.name = ? AND l.name = ?`,
		ownerName, loomName,
	).Scan(
		&l.ID, &l.OwnerID, &l.Name, &l.Description, &l.Visibility, &l.DefaultStream,
		&l.DiskPath, &l.SizeBytes, &l.CheckpointCount, &l.StreamCount,
		&l.PinCount, &l.SpinCount, &l.SpunFrom, &l.CreatedAt, &l.UpdatedAt, &l.SyncedAt,
		&l.OwnerName,
	)
	if err != nil {
		return nil, fmt.Errorf("get loom: %w", err)
	}
	l.FullName = l.OwnerName + "/" + l.Name
	return &l, nil
}

func (s *LoomStore) ListByOwner(ownerID string, limit, offset int) ([]models.Loom, int, error) {
	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM looms WHERE owner_id = ?", ownerID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(`
		SELECT l.id, l.owner_id, l.name, l.description, l.visibility, l.default_stream,
			   l.disk_path, l.size_bytes, l.checkpoint_count, l.stream_count,
			   l.pin_count, l.spin_count, l.spun_from, l.created_at, l.updated_at, l.synced_at,
			   o.name
		FROM looms l
		JOIN owners o ON o.id = l.owner_id
		WHERE l.owner_id = ?
		ORDER BY l.updated_at DESC
		LIMIT ? OFFSET ?`,
		ownerID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var looms []models.Loom
	for rows.Next() {
		var l models.Loom
		if err := rows.Scan(
			&l.ID, &l.OwnerID, &l.Name, &l.Description, &l.Visibility, &l.DefaultStream,
			&l.DiskPath, &l.SizeBytes, &l.CheckpointCount, &l.StreamCount,
			&l.PinCount, &l.SpinCount, &l.SpunFrom, &l.CreatedAt, &l.UpdatedAt, &l.SyncedAt,
			&l.OwnerName,
		); err != nil {
			return nil, 0, err
		}
		l.FullName = l.OwnerName + "/" + l.Name
		looms = append(looms, l)
	}
	return looms, total, nil
}

func (s *LoomStore) ListPublicByOwner(ownerName string, limit, offset int) ([]models.Loom, int, error) {
	var total int
	if err := s.db.QueryRow(`
		SELECT COUNT(*) FROM looms l JOIN owners o ON o.id = l.owner_id
		WHERE o.name = ? AND l.visibility = 'public'`,
		ownerName,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(`
		SELECT l.id, l.owner_id, l.name, l.description, l.visibility, l.default_stream,
			   l.disk_path, l.size_bytes, l.checkpoint_count, l.stream_count,
			   l.pin_count, l.spin_count, l.spun_from, l.created_at, l.updated_at, l.synced_at,
			   o.name
		FROM looms l
		JOIN owners o ON o.id = l.owner_id
		WHERE o.name = ? AND l.visibility = 'public'
		ORDER BY l.updated_at DESC
		LIMIT ? OFFSET ?`,
		ownerName, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var looms []models.Loom
	for rows.Next() {
		var l models.Loom
		if err := rows.Scan(
			&l.ID, &l.OwnerID, &l.Name, &l.Description, &l.Visibility, &l.DefaultStream,
			&l.DiskPath, &l.SizeBytes, &l.CheckpointCount, &l.StreamCount,
			&l.PinCount, &l.SpinCount, &l.SpunFrom, &l.CreatedAt, &l.UpdatedAt, &l.SyncedAt,
			&l.OwnerName,
		); err != nil {
			return nil, 0, err
		}
		l.FullName = l.OwnerName + "/" + l.Name
		looms = append(looms, l)
	}
	return looms, total, nil
}

func (s *LoomStore) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM looms WHERE id = ?", id)
	return err
}

func (s *LoomStore) Update(loom *models.Loom) error {
	_, err := s.db.Exec(`
		UPDATE looms SET description = ?, visibility = ?, default_stream = ?, updated_at = ?
		WHERE id = ?`,
		loom.Description, loom.Visibility, loom.DefaultStream, loom.UpdatedAt, loom.ID,
	)
	return err
}
