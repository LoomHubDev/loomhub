package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func (h *Handler) CreateWeaveRequest(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	var req struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		SourceStream string `json:"source_stream"`
		TargetStream string `json:"target_stream"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}
	if req.Title == "" || req.SourceStream == "" || req.TargetStream == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "Title, source_stream, and target_stream are required")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	wr := &models.WeaveRequest{
		ID:           ulid.Make().String(),
		LoomID:       loom.ID,
		AuthorID:     cu.ID,
		Title:        req.Title,
		Description:  req.Description,
		SourceStream: req.SourceStream,
		TargetStream: req.TargetStream,
		Status:       "open",
		CreatedAt:    now,
		UpdatedAt:    now,
		AuthorName:   cu.Username,
	}

	if err := h.weaveRequests.Create(wr); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to create weave request")
		return
	}

	h.activities.Log(cu.ID, &loom.ID, "create_wr", map[string]any{
		"wr_number": wr.Number, "title": wr.Title,
	})

	writeJSON(w, http.StatusCreated, wr)
}

func (h *Handler) ListWeaveRequests(w http.ResponseWriter, r *http.Request) {
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	status := r.URL.Query().Get("status")
	wrs, total, err := h.weaveRequests.ListByLoom(loom.ID, status, 50, 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": wrs, "total": total})
}

func (h *Handler) GetWeaveRequest(w http.ResponseWriter, r *http.Request) {
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	number, _ := strconv.Atoi(chi.URLParam(r, "number"))

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	wr, err := h.weaveRequests.GetByNumber(loom.ID, number)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Weave request not found")
		return
	}

	writeJSON(w, http.StatusOK, wr)
}

func (h *Handler) UpdateWeaveRequest(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	number, _ := strconv.Atoi(chi.URLParam(r, "number"))

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	wr, err := h.weaveRequests.GetByNumber(loom.ID, number)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Weave request not found")
		return
	}

	if wr.AuthorID != cu.ID && !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	if req.Title != nil {
		wr.Title = *req.Title
	}
	if req.Description != nil {
		wr.Description = *req.Description
	}

	if err := h.weaveRequests.Update(wr); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to update")
		return
	}
	writeJSON(w, http.StatusOK, wr)
}

func (h *Handler) WeaveWR(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	number, _ := strconv.Atoi(chi.URLParam(r, "number"))

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	wr, err := h.weaveRequests.GetByNumber(loom.ID, number)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Weave request not found")
		return
	}

	if wr.Status != "open" {
		writeError(w, http.StatusBadRequest, "bad_request", "Weave request is not open")
		return
	}

	if err := h.weaveRequests.UpdateStatus(wr.ID, "woven", &cu.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to weave")
		return
	}

	h.activities.Log(cu.ID, &loom.ID, "weave", map[string]any{"wr_number": wr.Number})

	wr.Status = "woven"
	writeJSON(w, http.StatusOK, wr)
}

func (h *Handler) CloseWR(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	number, _ := strconv.Atoi(chi.URLParam(r, "number"))

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	wr, err := h.weaveRequests.GetByNumber(loom.ID, number)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Weave request not found")
		return
	}

	if wr.AuthorID != cu.ID && !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	if err := h.weaveRequests.UpdateStatus(wr.ID, "closed", nil); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to close")
		return
	}

	wr.Status = "closed"
	writeJSON(w, http.StatusOK, wr)
}

// Comments

func (h *Handler) ListComments(w http.ResponseWriter, r *http.Request) {
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	number, _ := strconv.Atoi(chi.URLParam(r, "number"))

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	wr, err := h.weaveRequests.GetByNumber(loom.ID, number)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Weave request not found")
		return
	}

	comments, err := h.comments.ListByWR(wr.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if comments == nil {
		comments = []models.Comment{}
	}
	writeJSON(w, http.StatusOK, comments)
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	number, _ := strconv.Atoi(chi.URLParam(r, "number"))

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	wr, err := h.weaveRequests.GetByNumber(loom.ID, number)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Weave request not found")
		return
	}

	var req struct {
		Body       string  `json:"body"`
		SpaceID    *string `json:"space_id"`
		EntityPath *string `json:"entity_path"`
		LineNumber *int    `json:"line_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Body == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "Body is required")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	comment := &models.Comment{
		ID:         ulid.Make().String(),
		WRID:       wr.ID,
		AuthorID:   cu.ID,
		Body:       req.Body,
		SpaceID:    req.SpaceID,
		EntityPath: req.EntityPath,
		LineNumber: req.LineNumber,
		CreatedAt:  now,
		UpdatedAt:  now,
		AuthorName: cu.Username,
	}

	if err := h.comments.Create(comment); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to create comment")
		return
	}

	_ = loom // used above
	writeJSON(w, http.StatusCreated, comment)
}

// Reviews

func (h *Handler) ListReviews(w http.ResponseWriter, r *http.Request) {
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	number, _ := strconv.Atoi(chi.URLParam(r, "number"))

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	wr, err := h.weaveRequests.GetByNumber(loom.ID, number)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Weave request not found")
		return
	}

	reviews, err := h.reviews.ListByWR(wr.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if reviews == nil {
		reviews = []models.Review{}
	}
	writeJSON(w, http.StatusOK, reviews)
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	number, _ := strconv.Atoi(chi.URLParam(r, "number"))

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	wr, err := h.weaveRequests.GetByNumber(loom.ID, number)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Weave request not found")
		return
	}

	var req struct {
		Status string `json:"status"` // approved, changes_requested, commented
		Body   string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}
	if req.Status != "approved" && req.Status != "changes_requested" && req.Status != "commented" {
		writeError(w, http.StatusBadRequest, "bad_request", "Status must be approved, changes_requested, or commented")
		return
	}

	review := &models.Review{
		ID:           ulid.Make().String(),
		WRID:         wr.ID,
		ReviewerID:   cu.ID,
		Status:       req.Status,
		Body:         req.Body,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		ReviewerName: cu.Username,
	}

	if err := h.reviews.Create(review); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to create review")
		return
	}

	h.activities.Log(cu.ID, &loom.ID, "review", map[string]any{
		"wr_number": wr.Number, "status": req.Status,
	})

	writeJSON(w, http.StatusCreated, review)
}

// Pins

func (h *Handler) PinLoom(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	if err := h.pins.Pin(cu.ID, loom.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to pin")
		return
	}

	h.activities.Log(cu.ID, &loom.ID, "pin", nil)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UnpinLoom(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	if err := h.pins.Unpin(cu.ID, loom.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to unpin")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
