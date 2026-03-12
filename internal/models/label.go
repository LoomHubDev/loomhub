package models

type Label struct {
	ID          string `json:"id"`
	LoomID      string `json:"loom_id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}
