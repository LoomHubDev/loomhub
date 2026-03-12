package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/LoomHubDev/loomhub/internal/store"
	"github.com/oklog/ulid/v2"
)

type Dispatcher struct {
	webhooks *store.WebhookStore
	client   *http.Client
}

func NewDispatcher(webhooks *store.WebhookStore) *Dispatcher {
	return &Dispatcher{
		webhooks: webhooks,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (d *Dispatcher) Dispatch(loomID, event string, payload any) {
	hooks, err := d.webhooks.ListByLoomAndEvent(loomID, event)
	if err != nil || len(hooks) == 0 {
		return
	}

	data, _ := json.Marshal(payload)
	for _, wh := range hooks {
		go d.deliver(wh, event, data)
	}
}

func (d *Dispatcher) deliver(wh models.Webhook, event string, payload []byte) {
	start := time.Now()

	req, err := http.NewRequest("POST", wh.URL, bytes.NewReader(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-LoomHub-Event", event)
	req.Header.Set("X-LoomHub-Delivery", ulid.Make().String())

	if wh.Secret != "" {
		mac := hmac.New(sha256.New, []byte(wh.Secret))
		mac.Write(payload)
		sig := hex.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-LoomHub-Signature", fmt.Sprintf("sha256=%s", sig))
	}

	var statusCode *int
	var respBody string

	resp, err := d.client.Do(req)
	if err != nil {
		respBody = err.Error()
	} else {
		code := resp.StatusCode
		statusCode = &code
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		resp.Body.Close()
		respBody = string(buf[:n])
	}

	durationMS := int(time.Since(start).Milliseconds())
	delivery := &models.WebhookDelivery{
		ID:          ulid.Make().String(),
		WebhookID:   wh.ID,
		Event:       event,
		Payload:     string(payload),
		StatusCode:  statusCode,
		Response:    respBody,
		DurationMS:  &durationMS,
		DeliveredAt: time.Now().UTC().Format(time.RFC3339),
	}
	d.webhooks.LogDelivery(delivery)
}
