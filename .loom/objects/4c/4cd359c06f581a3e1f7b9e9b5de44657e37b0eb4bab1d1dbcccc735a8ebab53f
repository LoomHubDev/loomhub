package store

import (
	"database/sql"
)

type PinStore struct {
	db *sql.DB
}

func NewPinStore(db *sql.DB) *PinStore {
	return &PinStore{db: db}
}

func (s *PinStore) Pin(userID, loomID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"INSERT OR IGNORE INTO pins (user_id, loom_id, created_at) VALUES (?, ?, datetime('now'))",
		userID, loomID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Exec("UPDATE looms SET pin_count = (SELECT COUNT(*) FROM pins WHERE loom_id = ?) WHERE id = ?", loomID, loomID)
	return tx.Commit()
}

func (s *PinStore) Unpin(userID, loomID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	tx.Exec("DELETE FROM pins WHERE user_id = ? AND loom_id = ?", userID, loomID)
	tx.Exec("UPDATE looms SET pin_count = (SELECT COUNT(*) FROM pins WHERE loom_id = ?) WHERE id = ?", loomID, loomID)
	return tx.Commit()
}

func (s *PinStore) IsPinned(userID, loomID string) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM pins WHERE user_id = ? AND loom_id = ?", userID, loomID).Scan(&count)
	return count > 0, err
}
