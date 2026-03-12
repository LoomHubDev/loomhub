package models

type Owner struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"` // "user" or "org"
	CreatedAt string `json:"created_at"`
}
