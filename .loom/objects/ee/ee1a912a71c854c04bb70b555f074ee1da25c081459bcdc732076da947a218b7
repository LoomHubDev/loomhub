package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type WebhookStore struct {
	db *sql.DB
}

func NewWebhookStore(db *sql.DB) *WebhookStore {
	return &WebhookStore{db: db}
}

func (s *WebhookStore) Create(wh *models.Webhook) error {
	_, err := s.db.Exec(`
		INSERT INTO webhooks (id, loom_id, url, secret, events, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		wh.ID, wh.LoomID, wh.URL, wh.Secret, wh.Events, wh.Active, wh.CreatedAt, wh.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}
	return nil
}

func (s *WebhookStore) GetByID(id string) (*models.Webhook, error) {
	var wh models.Webhook
	err := s.db.QueryRow(`
		SELECT id, loom_id, url, secret, events, active, created_at, updated_at
		FROM webhooks WHERE id = ?`, id,
	).Scan(&wh.ID, &wh.LoomID, &wh.URL, &wh.Secret, &wh.Events, &wh.Active, &wh.CreatedAt, &wh.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

func (s *WebhookStore) ListByLoom(loomID string) ([]models.Webhook, error) {
	rows, err := s.db.Query(`
		SELECT id, loom_id, url, secret, events, active, created_at, updated_at
		FROM webhooks WHERE loom_id = ? ORDER BY created_at ASC`, loomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []models.Webhook
	for rows.Next() {
		var wh models.Webhook
		if err := rows.Scan(&wh.ID, &wh.LoomID, &wh.URL, &wh.Secret, &wh.Events, &wh.Active, &wh.CreatedAt, &wh.UpdatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, wh)
	}
	return webhooks, nil
}

func (s *WebhookStore) Update(wh *models.Webhook) error {
	_, err := s.db.Exec(`
		UPDATE webhooks SET url = ?, secret = ?, events = ?, active = ?, updated_at = ?
		WHERE id = ?`,
		wh.URL, wh.Secret, wh.Events, wh.Active, wh.UpdatedAt, wh.ID,
	)
	return err
}

func (s *WebhookStore) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM webhooks WHERE id = ?", id)
	return err
}

func (s *WebhookStore) ListByLoomAndEvent(loomID, event string) ([]models.Webhook, error) {
	webhooks, err := s.ListByLoom(loomID)
	if err != nil {
		return nil, err
	}
	var matched []models.Webhook
	for _, wh := range webhooks {
		if !wh.Active {
			continue
		}
		events := strings.Split(wh.Events, ",")
		for _, e := range events {
			if strings.TrimSpace(e) == event || strings.TrimSpace(e) == "*" {
				matched = append(matched, wh)
				break
			}
		}
	}
	return matched, nil
}

func (s *WebhookStore) LogDelivery(d *models.WebhookDelivery) error {
	_, err := s.db.Exec(`
		INSERT INTO webhook_deliveries (id, webhook_id, event, payload, status_code, response, duration_ms, delivered_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.WebhookID, d.Event, d.Payload, d.StatusCode, d.Response, d.DurationMS, d.DeliveredAt,
	)
	return err
}

func (s *WebhookStore) ListDeliveries(webhookID string, limit int) ([]models.WebhookDelivery, error) {
	rows, err := s.db.Query(`
		SELECT id, webhook_id, event, payload, status_code, response, duration_ms, delivered_at
		FROM webhook_deliveries WHERE webhook_id = ?
		ORDER BY delivered_at DESC LIMIT ?`, webhookID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []models.WebhookDelivery
	for rows.Next() {
		var d models.WebhookDelivery
		if err := rows.Scan(&d.ID, &d.WebhookID, &d.Event, &d.Payload, &d.StatusCode, &d.Response, &d.DurationMS, &d.DeliveredAt); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}
	return deliveries, nil
}
