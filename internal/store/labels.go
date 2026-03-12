package store

import (
	"database/sql"
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type LabelStore struct {
	db *sql.DB
}

func NewLabelStore(db *sql.DB) *LabelStore {
	return &LabelStore{db: db}
}

func (s *LabelStore) Create(l *models.Label) error {
	_, err := s.db.Exec(`
		INSERT INTO labels (id, loom_id, name, color, description)
		VALUES (?, ?, ?, ?, ?)`,
		l.ID, l.LoomID, l.Name, l.Color, l.Description,
	)
	if err != nil {
		return fmt.Errorf("create label: %w", err)
	}
	return nil
}

func (s *LabelStore) ListByLoom(loomID string) ([]models.Label, error) {
	rows, err := s.db.Query(`
		SELECT id, loom_id, name, color, description
		FROM labels WHERE loom_id = ? ORDER BY name ASC`, loomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []models.Label
	for rows.Next() {
		var l models.Label
		if err := rows.Scan(&l.ID, &l.LoomID, &l.Name, &l.Color, &l.Description); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, nil
}

func (s *LabelStore) Update(l *models.Label) error {
	_, err := s.db.Exec(`
		UPDATE labels SET name = ?, color = ?, description = ? WHERE id = ?`,
		l.Name, l.Color, l.Description, l.ID,
	)
	return err
}

func (s *LabelStore) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM labels WHERE id = ?", id)
	return err
}

func (s *LabelStore) GetByID(id string) (*models.Label, error) {
	var l models.Label
	err := s.db.QueryRow(`
		SELECT id, loom_id, name, color, description FROM labels WHERE id = ?`, id,
	).Scan(&l.ID, &l.LoomID, &l.Name, &l.Color, &l.Description)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *LabelStore) AddToWR(wrID, labelID string) error {
	_, err := s.db.Exec(
		"INSERT OR IGNORE INTO wr_labels (wr_id, label_id) VALUES (?, ?)",
		wrID, labelID,
	)
	return err
}

func (s *LabelStore) RemoveFromWR(wrID, labelID string) error {
	_, err := s.db.Exec("DELETE FROM wr_labels WHERE wr_id = ? AND label_id = ?", wrID, labelID)
	return err
}

func (s *LabelStore) ListByWR(wrID string) ([]models.Label, error) {
	rows, err := s.db.Query(`
		SELECT l.id, l.loom_id, l.name, l.color, l.description
		FROM labels l
		JOIN wr_labels wl ON wl.label_id = l.id
		WHERE wl.wr_id = ?
		ORDER BY l.name ASC`, wrID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []models.Label
	for rows.Next() {
		var l models.Label
		if err := rows.Scan(&l.ID, &l.LoomID, &l.Name, &l.Color, &l.Description); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, nil
}
