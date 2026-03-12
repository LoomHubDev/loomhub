package store

import (
	"fmt"
	"strings"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

type WebhookStore struct {
	db *gorm.DB
}

func NewWebhookStore(db *gorm.DB) *WebhookStore {
	return &WebhookStore{db: db}
}

func (s *WebhookStore) Create(wh *models.Webhook) error {
	if err := s.db.Create(wh).Error; err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}
	return nil
}

func (s *WebhookStore) GetByID(id string) (*models.Webhook, error) {
	var wh models.Webhook
	if err := s.db.First(&wh, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &wh, nil
}

func (s *WebhookStore) ListByLoom(loomID string) ([]models.Webhook, error) {
	var webhooks []models.Webhook
	err := s.db.Where("loom_id = ?", loomID).Order("created_at ASC").Find(&webhooks).Error
	return webhooks, err
}

func (s *WebhookStore) Update(wh *models.Webhook) error {
	return s.db.Model(wh).Updates(map[string]any{
		"url":        wh.URL,
		"secret":     wh.Secret,
		"events":     wh.Events,
		"active":     wh.Active,
		"updated_at": wh.UpdatedAt,
	}).Error
}

func (s *WebhookStore) Delete(id string) error {
	return s.db.Delete(&models.Webhook{}, "id = ?", id).Error
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
	return s.db.Create(d).Error
}

func (s *WebhookStore) ListDeliveries(webhookID string, limit int) ([]models.WebhookDelivery, error) {
	var deliveries []models.WebhookDelivery
	err := s.db.Where("webhook_id = ?", webhookID).Order("delivered_at DESC").Limit(limit).Find(&deliveries).Error
	return deliveries, err
}
