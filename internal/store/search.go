package store

import (
	"database/sql"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type SearchStore struct {
	db *sql.DB
}

func NewSearchStore(db *sql.DB) *SearchStore {
	return &SearchStore{db: db}
}

func (s *SearchStore) SearchLooms(query string, limit int) ([]models.Loom, error) {
	pattern := "%" + query + "%"
	rows, err := s.db.Query(`
		SELECT l.id, l.owner_id, l.name, l.description, l.visibility,
		       l.default_stream, l.checkpoint_count, l.stream_count, l.pin_count,
		       l.spin_count, l.created_at, l.updated_at, o.name
		FROM looms l
		JOIN owners o ON o.id = l.owner_id
		WHERE l.visibility = 'public'
		  AND (l.name LIKE ? OR l.description LIKE ?)
		ORDER BY l.pin_count DESC, l.updated_at DESC
		LIMIT ?`, pattern, pattern, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var looms []models.Loom
	for rows.Next() {
		var l models.Loom
		if err := rows.Scan(
			&l.ID, &l.OwnerID, &l.Name, &l.Description, &l.Visibility,
			&l.DefaultStream, &l.CheckpointCount, &l.StreamCount, &l.PinCount,
			&l.SpinCount, &l.CreatedAt, &l.UpdatedAt, &l.OwnerName,
		); err != nil {
			return nil, err
		}
		l.FullName = l.OwnerName + "/" + l.Name
		looms = append(looms, l)
	}
	return looms, nil
}

func (s *SearchStore) SearchUsers(query string, limit int) ([]models.User, error) {
	pattern := "%" + query + "%"
	rows, err := s.db.Query(`
		SELECT id, username, email, display_name, bio, avatar_url, is_admin, created_at, updated_at
		FROM users
		WHERE username LIKE ? OR display_name LIKE ?
		ORDER BY username ASC
		LIMIT ?`, pattern, pattern, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.DisplayName, &u.Bio, &u.AvatarURL, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *SearchStore) ExploreLooms(sort string, limit, offset int) ([]models.Loom, int, error) {
	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM looms WHERE visibility = 'public'").Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := "l.updated_at DESC"
	if sort == "pins" {
		orderBy = "l.pin_count DESC, l.updated_at DESC"
	} else if sort == "created" {
		orderBy = "l.created_at DESC"
	}

	rows, err := s.db.Query(`
		SELECT l.id, l.owner_id, l.name, l.description, l.visibility,
		       l.default_stream, l.checkpoint_count, l.stream_count, l.pin_count,
		       l.spin_count, l.created_at, l.updated_at, o.name
		FROM looms l
		JOIN owners o ON o.id = l.owner_id
		WHERE l.visibility = 'public'
		ORDER BY `+orderBy+`
		LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var looms []models.Loom
	for rows.Next() {
		var l models.Loom
		if err := rows.Scan(
			&l.ID, &l.OwnerID, &l.Name, &l.Description, &l.Visibility,
			&l.DefaultStream, &l.CheckpointCount, &l.StreamCount, &l.PinCount,
			&l.SpinCount, &l.CreatedAt, &l.UpdatedAt, &l.OwnerName,
		); err != nil {
			return nil, 0, err
		}
		l.FullName = l.OwnerName + "/" + l.Name
		looms = append(looms, l)
	}
	return looms, total, nil
}
