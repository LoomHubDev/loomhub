package store

import (
	"database/sql"
	"encoding/json"

	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/oklog/ulid/v2"
)

type ActivityStore struct {
	db *sql.DB
}

func NewActivityStore(db *sql.DB) *ActivityStore {
	return &ActivityStore{db: db}
}

func (s *ActivityStore) Log(actorID string, loomID *string, actType string, payload any) error {
	data, _ := json.Marshal(payload)
	_, err := s.db.Exec(
		"INSERT INTO activities (id, actor_id, loom_id, type, payload, created_at) VALUES (?, ?, ?, ?, ?, datetime('now'))",
		ulid.Make().String(), actorID, loomID, actType, string(data),
	)
	return err
}

func (s *ActivityStore) ListByUser(userID string, limit int) ([]models.Activity, error) {
	rows, err := s.db.Query(`
		SELECT a.id, a.actor_id, a.loom_id, a.type, a.payload, a.created_at,
			   u.username, l.name, o.name
		FROM activities a
		JOIN users u ON u.id = a.actor_id
		LEFT JOIN looms l ON l.id = a.loom_id
		LEFT JOIN owners o ON o.id = l.owner_id
		WHERE a.actor_id = ?
		ORDER BY a.created_at DESC LIMIT ?`, userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var a models.Activity
		var loomID, loomName, ownerName sql.NullString
		if err := rows.Scan(&a.ID, &a.ActorID, &loomID, &a.Type, &a.Payload, &a.CreatedAt,
			&a.ActorName, &loomName, &ownerName); err != nil {
			return nil, err
		}
		if loomID.Valid {
			a.LoomID = &loomID.String
		}
		if loomName.Valid && ownerName.Valid {
			fn := ownerName.String + "/" + loomName.String
			a.LoomFullName = &fn
		}
		activities = append(activities, a)
	}
	return activities, nil
}

func (s *ActivityStore) ListGlobal(limit int) ([]models.Activity, error) {
	rows, err := s.db.Query(`
		SELECT a.id, a.actor_id, a.loom_id, a.type, a.payload, a.created_at,
			   u.username, l.name, o.name
		FROM activities a
		JOIN users u ON u.id = a.actor_id
		LEFT JOIN looms l ON l.id = a.loom_id
		LEFT JOIN owners o ON o.id = l.owner_id
		ORDER BY a.created_at DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var a models.Activity
		var loomID, loomName, ownerName sql.NullString
		if err := rows.Scan(&a.ID, &a.ActorID, &loomID, &a.Type, &a.Payload, &a.CreatedAt,
			&a.ActorName, &loomName, &ownerName); err != nil {
			return nil, err
		}
		if loomID.Valid {
			a.LoomID = &loomID.String
		}
		if loomName.Valid && ownerName.Valid {
			fn := ownerName.String + "/" + loomName.String
			a.LoomFullName = &fn
		}
		activities = append(activities, a)
	}
	return activities, nil
}
