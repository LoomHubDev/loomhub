package store

import (
	"database/sql"
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type ReviewStore struct {
	db *sql.DB
}

func NewReviewStore(db *sql.DB) *ReviewStore {
	return &ReviewStore{db: db}
}

func (s *ReviewStore) Create(r *models.Review) error {
	_, err := s.db.Exec(`
		INSERT INTO reviews (id, wr_id, reviewer_id, status, body, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		r.ID, r.WRID, r.ReviewerID, r.Status, r.Body, r.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create review: %w", err)
	}
	return nil
}

func (s *ReviewStore) ListByWR(wrID string) ([]models.Review, error) {
	rows, err := s.db.Query(`
		SELECT r.id, r.wr_id, r.reviewer_id, r.status, r.body, r.created_at, u.username
		FROM reviews r
		JOIN users u ON u.id = r.reviewer_id
		WHERE r.wr_id = ?
		ORDER BY r.created_at ASC`, wrID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var r models.Review
		if err := rows.Scan(&r.ID, &r.WRID, &r.ReviewerID, &r.Status, &r.Body, &r.CreatedAt, &r.ReviewerName); err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}
