package api

import (
	"encoding/json"
	"net/http"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/sync"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandleNegotiate(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(owner, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	// Private looms require at least read access
	if loom.Visibility == "private" {
		cu := auth.GetUser(r.Context())
		if cu == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
			return
		}
		if h.collabs.CheckAccess(loom, cu.ID, cu.IsAdmin) == "" {
			writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
			return
		}
	}

	var req sync.NegotiateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	resp, err := h.sync.Negotiate(loom, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) HandleSend(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	// Push requires authentication
	cu := auth.GetUser(r.Context())
	if cu == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	loom, err := h.looms.GetByOwnerAndName(owner, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	// Push requires write access
	access := h.collabs.CheckAccess(loom, cu.ID, cu.IsAdmin)
	if access != "write" && access != "admin" {
		writeError(w, http.StatusForbidden, "forbidden", "Write access required to push")
		return
	}

	var req sync.PushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	resp, err := h.sync.Send(loom, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	// Log push activity
	h.activities.Log(cu.ID, &loom.ID, "push", map[string]any{
		"ops_count": len(req.Operations),
	})

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) HandleReceive(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(owner, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	// Private looms require at least read access
	if loom.Visibility == "private" {
		cu := auth.GetUser(r.Context())
		if cu == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
			return
		}
		if h.collabs.CheckAccess(loom, cu.ID, cu.IsAdmin) == "" {
			writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
			return
		}
	}

	var req sync.PullRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	resp, err := h.sync.Receive(loom, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
