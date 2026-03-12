package api

import (
	"encoding/json"
	"net/http"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func (h *Handler) CreateLabel(w http.ResponseWriter, r *http.Request) {
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
		Name        string `json:"name"`
		Color       string `json:"color"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "Name is required")
		return
	}
	if req.Color == "" {
		req.Color = "#888888"
	}

	label := &models.Label{
		ID:          ulid.Make().String(),
		LoomID:      loom.ID,
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
	}

	if err := h.labels.Create(label); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to create label")
		return
	}
	writeJSON(w, http.StatusCreated, label)
}

func (h *Handler) ListLabels(w http.ResponseWriter, r *http.Request) {
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	labels, err := h.labels.ListByLoom(loom.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if labels == nil {
		labels = []models.Label{}
	}
	writeJSON(w, http.StatusOK, labels)
}

func (h *Handler) UpdateLabel(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	labelID := chi.URLParam(r, "labelID")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}
	if loom.OwnerID != cu.ID && !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	label, err := h.labels.GetByID(labelID)
	if err != nil || label.LoomID != loom.ID {
		writeError(w, http.StatusNotFound, "not_found", "Label not found")
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Color       *string `json:"color"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}
	if req.Name != nil {
		label.Name = *req.Name
	}
	if req.Color != nil {
		label.Color = *req.Color
	}
	if req.Description != nil {
		label.Description = *req.Description
	}

	if err := h.labels.Update(label); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to update label")
		return
	}
	writeJSON(w, http.StatusOK, label)
}

func (h *Handler) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	labelID := chi.URLParam(r, "labelID")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}
	if loom.OwnerID != cu.ID && !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	label, err := h.labels.GetByID(labelID)
	if err != nil || label.LoomID != loom.ID {
		writeError(w, http.StatusNotFound, "not_found", "Label not found")
		return
	}

	if err := h.labels.Delete(labelID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to delete label")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AddLabelToWR(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	labelID := chi.URLParam(r, "labelID")

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
		WRID string `json:"wr_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.WRID == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "wr_id is required")
		return
	}

	if err := h.labels.AddToWR(req.WRID, labelID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to add label")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
