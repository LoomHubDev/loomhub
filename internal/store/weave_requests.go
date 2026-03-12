package store

import (
	"database/sql"
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type WeaveRequestStore struct {
	db *sql.DB
}

func NewWeaveRequestStore(db *sql.DB) *WeaveRequestStore {
	return &WeaveRequestStore{db: db}
}

func (s *WeaveRequestStore) Create(wr *models.WeaveRequest) error {
	// Get next number for this loom
	var maxNum int
	s.db.QueryRow("SELECT COALESCE(MAX(number), 0) FROM weave_requests WHERE loom_id = ?", wr.LoomID).Scan(&maxNum)
	wr.Number = maxNum + 1

	_, err := s.db.Exec(`
		INSERT INTO weave_requests (id, loom_id, number, author_id, title, description,
			source_stream, target_stream, source_head_seq, target_head_seq, merge_base_seq,
			status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		wr.ID, wr.LoomID, wr.Number, wr.AuthorID, wr.Title, wr.Description,
		wr.SourceStream, wr.TargetStream, wr.SourceHeadSeq, wr.TargetHeadSeq, wr.MergeBaseSeq,
		wr.Status, wr.CreatedAt, wr.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create weave request: %w", err)
	}
	return nil
}

func (s *WeaveRequestStore) GetByNumber(loomID string, number int) (*models.WeaveRequest, error) {
	var wr models.WeaveRequest
	err := s.db.QueryRow(`
		SELECT wr.id, wr.loom_id, wr.number, wr.author_id, wr.title, wr.description,
			   wr.source_stream, wr.target_stream, wr.source_head_seq, wr.target_head_seq,
			   wr.merge_base_seq, wr.status, wr.woven_by, wr.woven_at, wr.closed_at,
			   wr.created_at, wr.updated_at, u.username
		FROM weave_requests wr
		JOIN users u ON u.id = wr.author_id
		WHERE wr.loom_id = ? AND wr.number = ?`,
		loomID, number,
	).Scan(
		&wr.ID, &wr.LoomID, &wr.Number, &wr.AuthorID, &wr.Title, &wr.Description,
		&wr.SourceStream, &wr.TargetStream, &wr.SourceHeadSeq, &wr.TargetHeadSeq,
		&wr.MergeBaseSeq, &wr.Status, &wr.WovenBy, &wr.WovenAt, &wr.ClosedAt,
		&wr.CreatedAt, &wr.UpdatedAt, &wr.AuthorName,
	)
	if err != nil {
		return nil, fmt.Errorf("get weave request: %w", err)
	}
	return &wr, nil
}

func (s *WeaveRequestStore) ListByLoom(loomID, status string, limit, offset int) ([]models.WeaveRequest, int, error) {
	var total int
	query := "SELECT COUNT(*) FROM weave_requests WHERE loom_id = ?"
	args := []any{loomID}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if err := s.db.QueryRow(query, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	selectQuery := `
		SELECT wr.id, wr.loom_id, wr.number, wr.author_id, wr.title, wr.description,
			   wr.source_stream, wr.target_stream, wr.source_head_seq, wr.target_head_seq,
			   wr.merge_base_seq, wr.status, wr.woven_by, wr.woven_at, wr.closed_at,
			   wr.created_at, wr.updated_at, u.username
		FROM weave_requests wr
		JOIN users u ON u.id = wr.author_id
		WHERE wr.loom_id = ?`
	selectArgs := []any{loomID}
	if status != "" {
		selectQuery += " AND wr.status = ?"
		selectArgs = append(selectArgs, status)
	}
	selectQuery += " ORDER BY wr.updated_at DESC LIMIT ? OFFSET ?"
	selectArgs = append(selectArgs, limit, offset)

	rows, err := s.db.Query(selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var wrs []models.WeaveRequest
	for rows.Next() {
		var wr models.WeaveRequest
		if err := rows.Scan(
			&wr.ID, &wr.LoomID, &wr.Number, &wr.AuthorID, &wr.Title, &wr.Description,
			&wr.SourceStream, &wr.TargetStream, &wr.SourceHeadSeq, &wr.TargetHeadSeq,
			&wr.MergeBaseSeq, &wr.Status, &wr.WovenBy, &wr.WovenAt, &wr.ClosedAt,
			&wr.CreatedAt, &wr.UpdatedAt, &wr.AuthorName,
		); err != nil {
			return nil, 0, err
		}
		wrs = append(wrs, wr)
	}
	return wrs, total, nil
}

func (s *WeaveRequestStore) UpdateStatus(id, status string, wovenBy *string) error {
	if status == "woven" {
		_, err := s.db.Exec(
			"UPDATE weave_requests SET status = ?, woven_by = ?, woven_at = datetime('now'), updated_at = datetime('now') WHERE id = ?",
			status, wovenBy, id,
		)
		return err
	}
	if status == "closed" {
		_, err := s.db.Exec(
			"UPDATE weave_requests SET status = ?, closed_at = datetime('now'), updated_at = datetime('now') WHERE id = ?",
			status, id,
		)
		return err
	}
	_, err := s.db.Exec("UPDATE weave_requests SET status = ?, updated_at = datetime('now') WHERE id = ?", status, id)
	return err
}

func (s *WeaveRequestStore) Update(wr *models.WeaveRequest) error {
	_, err := s.db.Exec(
		"UPDATE weave_requests SET title = ?, description = ?, updated_at = datetime('now') WHERE id = ?",
		wr.Title, wr.Description, wr.ID,
	)
	return err
}
