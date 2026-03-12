package models

type Webhook struct {
	ID        string `json:"id"`
	LoomID    string `json:"loom_id"`
	URL       string `json:"url"`
	Secret    string `json:"-"`
	Events    string `json:"events"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type WebhookDelivery struct {
	ID          string `json:"id"`
	WebhookID   string `json:"webhook_id"`
	Event       string `json:"event"`
	Payload     string `json:"payload"`
	StatusCode  *int   `json:"status_code"`
	Response    string `json:"response"`
	DurationMS  *int   `json:"duration_ms"`
	DeliveredAt string `json:"delivered_at"`
}
