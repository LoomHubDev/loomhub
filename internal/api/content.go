package api

import (
	"net/http"
	"strconv"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/LoomHubDev/loomhub/internal/sync"
	"github.com/go-chi/chi/v5"
)

// checkLoomRead resolves a loom and verifies the caller has read access.
// Returns nil loom on error (response already written).
func (h *Handler) checkLoomRead(w http.ResponseWriter, r *http.Request) *models.Loom {
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return nil
	}
	if loom.Visibility == "private" {
		cu := auth.GetUser(r.Context())
		if cu == nil || h.collabs.CheckAccess(loom, cu.ID, cu.IsAdmin) == "" {
			writeError(w, http.StatusNotFound, "not_found", "Loom not found")
			return nil
		}
	}
	return loom
}

func (h *Handler) ListStreams(w http.ResponseWriter, r *http.Request) {
	loom := h.checkLoomRead(w, r)
	if loom == nil {
		return
	}

	streams, err := h.sync.ListStreams(loom)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to list streams")
		return
	}
	writeJSON(w, http.StatusOK, streams)
}

func (h *Handler) ListEntities(w http.ResponseWriter, r *http.Request) {
	loom := h.checkLoomRead(w, r)
	if loom == nil {
		return
	}
	stream := r.URL.Query().Get("stream")

	entities, err := h.sync.ListEntities(loom, stream)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to list entities")
		return
	}
	writeJSON(w, http.StatusOK, entities)
}

func (h *Handler) ListCheckpoints(w http.ResponseWriter, r *http.Request) {
	loom := h.checkLoomRead(w, r)
	if loom == nil {
		return
	}

	limit := 50
	offset := 0
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && o > 0 {
		offset = o
	}

	f := sync.CheckpointFilter{
		Stream: r.URL.Query().Get("stream"),
		Type:   r.URL.Query().Get("type"),
		Author: r.URL.Query().Get("author"),
		Path:   r.URL.Query().Get("path"),
		Limit:  limit,
		Offset: offset,
	}

	checkpoints, total, err := h.sync.ListCheckpoints(loom, f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to list checkpoints")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": checkpoints, "total": total})
}

func (h *Handler) CompareStreams(w http.ResponseWriter, r *http.Request) {
	loom := h.checkLoomRead(w, r)
	if loom == nil {
		return
	}
	source := r.URL.Query().Get("source")
	target := r.URL.Query().Get("target")

	if source == "" || target == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "source and target query params required")
		return
	}

	diffs, err := h.sync.CompareStreams(loom, source, target)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to compute diff")
		return
	}
	writeJSON(w, http.StatusOK, diffs)
}

func (h *Handler) GetEntityContent(w http.ResponseWriter, r *http.Request) {
	loom := h.checkLoomRead(w, r)
	if loom == nil {
		return
	}
	ref := chi.URLParam(r, "ref")

	content, err := h.sync.GetEntityContent(loom, ref)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Object not found")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
