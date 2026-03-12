package models

type Webhook struct {
	ID        string `json:"id" gorm:"primaryKey"`
	LoomID    string `json:"loom_id" gorm:"index"`
	URL       string `json:"url"`
	Secret    string `json:"-"`
	Events    string `json:"events" gorm:"default:send"`
	Active    bool   `json:"active" gorm:"default:true"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type WebhookDelivery struct {
	ID          string `json:"id" gorm:"primaryKey"`
	WebhookID   string `json:"webhook_id" gorm:"index"`
	Event       string `json:"event"`
	Payload     string `json:"payload"`
	StatusCode  *int   `json:"status_code"`
	Response    string `json:"response"`
	DurationMS  *int   `json:"duration_ms"`
	DeliveredAt string `json:"delivered_at"`
}
