package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func (h *Handler) CreateLoom(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Visibility  string `json:"visibility"`
		Owner       string `json:"owner"` // username or org name
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	req.Name = strings.ToLower(strings.TrimSpace(req.Name))
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "Name is required")
		return
	}
	if req.Visibility == "" {
		req.Visibility = "public"
	}
	if req.Visibility != "public" && req.Visibility != "private" {
		writeError(w, http.StatusBadRequest, "bad_request", "Visibility must be public or private")
		return
	}

	// Resolve owner
	ownerName := req.Owner
	if ownerName == "" {
		ownerName = cu.Username
	}

	owner, err := h.owners.GetByName(ownerName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Owner not found")
		return
	}

	// Check permission: user must own or be admin of the owner
	if owner.ID != cu.ID && !cu.IsAdmin {
		// TODO: check org membership
		writeError(w, http.StatusForbidden, "forbidden", "You don't have permission to create looms for this owner")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	id := ulid.Make().String()
	diskPath := filepath.Join("looms", ownerName, req.Name)

	loom := &models.Loom{
		ID:            id,
		OwnerID:       owner.ID,
		Name:          req.Name,
		Description:   req.Description,
		Visibility:    req.Visibility,
		DefaultStream: "main",
		DiskPath:      diskPath,
		CreatedAt:     now,
		UpdatedAt:     now,
		OwnerName:     ownerName,
		FullName:      ownerName + "/" + req.Name,
	}

	if err := h.looms.Create(loom); err != nil {
		writeError(w, http.StatusConflict, "conflict", "A loom with this name already exists")
		return
	}

	// Create on-disk directories
	loomDataPath := filepath.Join(h.cfg.Server.DataDir, diskPath)
	os.MkdirAll(loomDataPath, 0755)

	// Initialize per-loom database
	loomDB, err := sql.Open("sqlite", filepath.Join(loomDataPath, "loom.db")+"?_pragma=journal_mode(wal)&_pragma=foreign_keys(1)")
	if err == nil {
		initLoomDB(loomDB)
		loomDB.Close()
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":             loom.ID,
		"owner":          loom.OwnerName,
		"name":           loom.Name,
		"full_name":      loom.FullName,
		"description":    loom.Description,
		"visibility":     loom.Visibility,
		"default_stream": loom.DefaultStream,
		"replicate_url":  fmt.Sprintf("%s/%s/%s", h.cfg.Server.BaseURL, ownerName, req.Name),
		"created_at":     loom.CreatedAt,
	})
}

func (h *Handler) GetLoom(w http.ResponseWriter, r *http.Request) {
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	if loom.Visibility == "private" {
		cu := auth.GetUser(r.Context())
		if cu == nil || h.collabs.CheckAccess(loom, cu.ID, cu.IsAdmin) == "" {
			writeError(w, http.StatusNotFound, "not_found", "Loom not found")
			return
		}
	}

	writeJSON(w, http.StatusOK, loom)
}

func (h *Handler) UpdateLoom(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	if h.collabs.CheckAccess(loom, cu.ID, cu.IsAdmin) != "admin" {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	var req struct {
		Description   *string `json:"description"`
		Visibility    *string `json:"visibility"`
		DefaultStream *string `json:"default_stream"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	if req.Description != nil {
		loom.Description = *req.Description
	}
	if req.Visibility != nil {
		loom.Visibility = *req.Visibility
	}
	if req.DefaultStream != nil {
		loom.DefaultStream = *req.DefaultStream
	}
	loom.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := h.looms.Update(loom); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to update loom")
		return
	}
	writeJSON(w, http.StatusOK, loom)
}

func (h *Handler) DeleteLoom(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	if h.collabs.CheckAccess(loom, cu.ID, cu.IsAdmin) != "admin" {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	if err := h.looms.Delete(loom.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to delete loom")
		return
	}

	// Clean up disk
	loomDataPath := filepath.Join(h.cfg.Server.DataDir, loom.DiskPath)
	os.RemoveAll(loomDataPath)

	w.WriteHeader(http.StatusNoContent)
}

func initLoomDB(db *sql.DB) {
	db.Exec(`
		CREATE TABLE IF NOT EXISTS operations (
			id         TEXT PRIMARY KEY,
			seq        INTEGER NOT NULL,
			stream_id  TEXT NOT NULL,
			space_id   TEXT NOT NULL,
			entity_id  TEXT NOT NULL,
			type       TEXT NOT NULL,
			path       TEXT NOT NULL,
			delta      BLOB,
			object_ref TEXT,
			parent_seq INTEGER,
			author     TEXT NOT NULL,
			timestamp  TEXT NOT NULL,
			meta       TEXT NOT NULL DEFAULT '{}'
		);
		CREATE INDEX IF NOT EXISTS idx_ops_seq ON operations(seq);
		CREATE INDEX IF NOT EXISTS idx_ops_stream ON operations(stream_id, seq);

		CREATE TABLE IF NOT EXISTS checkpoints (
			id        TEXT PRIMARY KEY,
			stream_id TEXT NOT NULL,
			seq       INTEGER NOT NULL,
			title     TEXT NOT NULL,
			summary   TEXT NOT NULL DEFAULT '',
			author    TEXT NOT NULL,
			source    TEXT NOT NULL DEFAULT 'manual',
			tags      TEXT NOT NULL DEFAULT '{}',
			created_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_cp_stream ON checkpoints(stream_id, seq);

		CREATE TABLE IF NOT EXISTS streams (
			id        TEXT PRIMARY KEY,
			name      TEXT NOT NULL UNIQUE,
			head_seq  INTEGER NOT NULL DEFAULT 0,
			parent_id TEXT,
			status    TEXT NOT NULL DEFAULT 'active',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS entities (
			id        TEXT NOT NULL,
			space_id  TEXT NOT NULL,
			path      TEXT NOT NULL,
			kind      TEXT NOT NULL DEFAULT 'file',
			object_ref TEXT,
			size      INTEGER NOT NULL DEFAULT 0,
			mod_time  TEXT,
			PRIMARY KEY (space_id, id)
		);
		CREATE INDEX IF NOT EXISTS idx_entities_path ON entities(space_id, path);

		CREATE TABLE IF NOT EXISTS objects (
			hash       TEXT PRIMARY KEY,
			size       INTEGER NOT NULL DEFAULT 0,
			compressed BOOLEAN NOT NULL DEFAULT FALSE
		);

		CREATE TABLE IF NOT EXISTS metadata (
			key   TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);
	`)
}
