package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func (h *Handler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}
	if loom.OwnerID != cu.ID && !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	var req struct {
		URL    string `json:"url"`
		Secret string `json:"secret"`
		Events string `json:"events"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "URL is required")
		return
	}
	if req.Events == "" {
		req.Events = "send"
	}
	if req.Secret == "" {
		b := make([]byte, 20)
		rand.Read(b)
		req.Secret = hex.EncodeToString(b)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	wh := &models.Webhook{
		ID:        ulid.Make().String(),
		LoomID:    loom.ID,
		URL:       req.URL,
		Secret:    req.Secret,
		Events:    req.Events,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.webhooks.Create(wh); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to create webhook")
		return
	}
	writeJSON(w, http.StatusCreated, wh)
}

func (h *Handler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}
	if loom.OwnerID != cu.ID && !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	webhooks, err := h.webhooks.ListByLoom(loom.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if webhooks == nil {
		webhooks = []models.Webhook{}
	}
	writeJSON(w, http.StatusOK, webhooks)
}

func (h *Handler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	whID := chi.URLParam(r, "webhookID")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}
	if loom.OwnerID != cu.ID && !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	wh, err := h.webhooks.GetByID(whID)
	if err != nil || wh.LoomID != loom.ID {
		writeError(w, http.StatusNotFound, "not_found", "Webhook not found")
		return
	}

	var req struct {
		URL    *string `json:"url"`
		Events *string `json:"events"`
		Active *bool   `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}
	if req.URL != nil {
		wh.URL = *req.URL
	}
	if req.Events != nil {
		wh.Events = *req.Events
	}
	if req.Active != nil {
		wh.Active = *req.Active
	}
	wh.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := h.webhooks.Update(wh); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to update webhook")
		return
	}
	writeJSON(w, http.StatusOK, wh)
}

func (h *Handler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	whID := chi.URLParam(r, "webhookID")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}
	if loom.OwnerID != cu.ID && !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	wh, err := h.webhooks.GetByID(whID)
	if err != nil || wh.LoomID != loom.ID {
		writeError(w, http.StatusNotFound, "not_found", "Webhook not found")
		return
	}

	if err := h.webhooks.Delete(whID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to delete webhook")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListWebhookDeliveries(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	whID := chi.URLParam(r, "webhookID")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}
	if loom.OwnerID != cu.ID && !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	wh, err := h.webhooks.GetByID(whID)
	if err != nil || wh.LoomID != loom.ID {
		writeError(w, http.StatusNotFound, "not_found", "Webhook not found")
		return
	}

	deliveries, err := h.webhooks.ListDeliveries(whID, 20)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if deliveries == nil {
		deliveries = []models.WebhookDelivery{}
	}
	writeJSON(w, http.StatusOK, deliveries)
}
