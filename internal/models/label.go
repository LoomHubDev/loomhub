package models

type Label struct {
	ID          string `json:"id" gorm:"primaryKey"`
	LoomID      string `json:"loom_id" gorm:"uniqueIndex:idx_labels_loom_name"`
	Name        string `json:"name" gorm:"uniqueIndex:idx_labels_loom_name"`
	Color       string `json:"color" gorm:"default:#888888"`
	Description string `json:"description"`
}
