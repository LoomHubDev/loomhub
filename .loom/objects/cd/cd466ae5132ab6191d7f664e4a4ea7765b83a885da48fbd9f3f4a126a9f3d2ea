package store

import (
	"database/sql"
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type CommentStore struct {
	db *sql.DB
}

func NewCommentStore(db *sql.DB) *CommentStore {
	return &CommentStore{db: db}
}

func (s *CommentStore) Create(c *models.Comment) error {
	_, err := s.db.Exec(`
		INSERT INTO wr_comments (id, wr_id, author_id, body, space_id, entity_path, line_number, is_system, event_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.WRID, c.AuthorID, c.Body, c.SpaceID, c.EntityPath, c.LineNumber,
		c.IsSystem, c.EventType, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

func (s *CommentStore) ListByWR(wrID string) ([]models.Comment, error) {
	rows, err := s.db.Query(`
		SELECT c.id, c.wr_id, c.author_id, c.body, c.space_id, c.entity_path, c.line_number,
			   c.is_system, c.event_type, c.created_at, c.updated_at, u.username
		FROM wr_comments c
		JOIN users u ON u.id = c.author_id
		WHERE c.wr_id = ?
		ORDER BY c.created_at ASC`, wrID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(
			&c.ID, &c.WRID, &c.AuthorID, &c.Body, &c.SpaceID, &c.EntityPath, &c.LineNumber,
			&c.IsSystem, &c.EventType, &c.CreatedAt, &c.UpdatedAt, &c.AuthorName,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (s *CommentStore) Update(id, body, updatedAt string) error {
	_, err := s.db.Exec("UPDATE wr_comments SET body = ?, updated_at = ? WHERE id = ?", body, updatedAt, id)
	return err
}

func (s *CommentStore) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM wr_comments WHERE id = ?", id)
	return err
}

func (s *CommentStore) GetByID(id string) (*models.Comment, error) {
	var c models.Comment
	err := s.db.QueryRow(`
		SELECT id, wr_id, author_id, body, space_id, entity_path, line_number,
			   is_system, event_type, created_at, updated_at
		FROM wr_comments WHERE id = ?`, id,
	).Scan(
		&c.ID, &c.WRID, &c.AuthorID, &c.Body, &c.SpaceID, &c.EntityPath, &c.LineNumber,
		&c.IsSystem, &c.EventType, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
