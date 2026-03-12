package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func (h *Handler) CreateAccessToken(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())

	var req struct {
		Name   string `json:"name"`
		Scopes string `json:"scopes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "Name is required")
		return
	}
	if req.Scopes == "" {
		req.Scopes = "read,write"
	}

	// Generate random token
	raw := make([]byte, 32)
	rand.Read(raw)
	plainToken := "lh_" + hex.EncodeToString(raw)

	// Store hash
	hash := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(hash[:])

	token := &models.AccessToken{
		ID:        ulid.Make().String(),
		UserID:    cu.ID,
		Name:      req.Name,
		TokenHash: tokenHash,
		Scopes:    req.Scopes,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if err := h.tokens.Create(token); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to create token")
		return
	}

	// Return the plain token only once
	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         token.ID,
		"name":       token.Name,
		"token":      plainToken,
		"scopes":     token.Scopes,
		"created_at": token.CreatedAt,
	})
}

func (h *Handler) ListAccessTokens(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	tokens, err := h.tokens.ListByUser(cu.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if tokens == nil {
		tokens = []models.AccessToken{}
	}
	writeJSON(w, http.StatusOK, tokens)
}

func (h *Handler) DeleteAccessToken(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	tokenID := chi.URLParam(r, "tokenID")
	if err := h.tokens.Delete(tokenID, cu.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to delete token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
