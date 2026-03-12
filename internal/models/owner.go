package models

type Owner struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"uniqueIndex"`
	Type      string `json:"type" gorm:"column:type"`
	CreatedAt string `json:"created_at"`
}
