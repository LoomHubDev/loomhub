package api

import (
	"encoding/json"
	"net/http"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListCollaborators(w http.ResponseWriter, r *http.Request) {
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	// Only owner/admin/collaborators with admin can see collaborator list
	cu := auth.GetUser(r.Context())
	if cu == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}
	access := h.collabs.CheckAccess(loom, cu.ID, cu.IsAdmin)
	if access != "admin" && access != "write" {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	collabs, err := h.collabs.List(loom.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to list collaborators")
		return
	}

	if collabs == nil {
		writeJSON(w, http.StatusOK, []struct{}{})
		return
	}
	writeJSON(w, http.StatusOK, collabs)
}

func (h *Handler) AddCollaborator(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	// Only loom admin can add collaborators
	access := h.collabs.CheckAccess(loom, cu.ID, cu.IsAdmin)
	if access != "admin" {
		writeError(w, http.StatusForbidden, "forbidden", "Only admins can manage collaborators")
		return
	}

	var req struct {
		Username   string `json:"username"`
		Permission string `json:"permission"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	if req.Username == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "Username is required")
		return
	}
	if req.Permission == "" {
		req.Permission = "write"
	}
	if req.Permission != "read" && req.Permission != "write" && req.Permission != "admin" {
		writeError(w, http.StatusBadRequest, "bad_request", "Permission must be read, write, or admin")
		return
	}

	// Resolve username to user ID
	user, err := h.users.GetByUsername(req.Username)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "User not found")
		return
	}

	// Can't add the owner as a collaborator
	if user.ID == loom.OwnerID {
		writeError(w, http.StatusBadRequest, "bad_request", "Owner is already an admin")
		return
	}

	if err := h.collabs.Add(loom.ID, user.ID, req.Permission); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to add collaborator")
		return
	}

	h.activities.Log(cu.ID, &loom.ID, "collaborator_added", map[string]string{
		"username":   req.Username,
		"permission": req.Permission,
	})

	writeJSON(w, http.StatusCreated, map[string]string{
		"username":   req.Username,
		"permission": req.Permission,
	})
}

func (h *Handler) RemoveCollaborator(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	ownerName := chi.URLParam(r, "owner")
	loomName := chi.URLParam(r, "loom")
	username := chi.URLParam(r, "username")

	loom, err := h.looms.GetByOwnerAndName(ownerName, loomName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Loom not found")
		return
	}

	access := h.collabs.CheckAccess(loom, cu.ID, cu.IsAdmin)
	if access != "admin" {
		writeError(w, http.StatusForbidden, "forbidden", "Only admins can manage collaborators")
		return
	}

	user, err := h.users.GetByUsername(username)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "User not found")
		return
	}

	if err := h.collabs.Remove(loom.ID, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to remove collaborator")
		return
	}

	h.activities.Log(cu.ID, &loom.ID, "collaborator_removed", map[string]string{
		"username": username,
	})

	w.WriteHeader(http.StatusNoContent)
}
