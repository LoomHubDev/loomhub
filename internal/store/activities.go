package store

import (
	"encoding/json"
	"time"

	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type ActivityStore struct {
	db *gorm.DB
}

func NewActivityStore(db *gorm.DB) *ActivityStore {
	return &ActivityStore{db: db}
}

func (s *ActivityStore) Log(actorID string, loomID *string, actType string, payload any) error {
	data, _ := json.Marshal(payload)
	now := time.Now().UTC().Format(time.RFC3339)
	a := models.Activity{
		ID: ulid.Make().String(), ActorID: actorID, LoomID: loomID,
		Type: actType, Payload: string(data), CreatedAt: now,
	}
	return s.db.Create(&a).Error
}

func (s *ActivityStore) ListByUser(userID string, limit int) ([]models.Activity, error) {
	var activities []models.Activity
	err := s.db.Preload("Actor").Preload("LoomRef.Owner").
		Where("actor_id = ?", userID).
		Order("created_at DESC").Limit(limit).
		Find(&activities).Error
	if err != nil {
		return nil, err
	}
	for i := range activities {
		activities[i].FillJoinedFields()
	}
	return activities, nil
}

func (s *ActivityStore) ListGlobal(limit int) ([]models.Activity, error) {
	var activities []models.Activity
	err := s.db.Preload("Actor").Preload("LoomRef.Owner").
		Order("created_at DESC").Limit(limit).
		Find(&activities).Error
	if err != nil {
		return nil, err
	}
	for i := range activities {
		activities[i].FillJoinedFields()
	}
	return activities, nil
}
