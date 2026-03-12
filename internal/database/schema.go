package database

import (
	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Owner{},
		&models.User{},
		&models.AccessToken{},
		&models.Organization{},
		&models.OrgMember{},
		&models.Loom{},
		&models.Collaborator{},
		&models.Pin{},
		&models.WeaveRequest{},
		&models.Comment{},
		&models.Review{},
		&models.Label{},
		&models.WRLabel{},
		&models.Webhook{},
		&models.WebhookDelivery{},
		&models.Activity{},
		&models.ObjectRef{},
	)
}
