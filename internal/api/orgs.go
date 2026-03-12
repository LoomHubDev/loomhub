package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func (h *Handler) CreateOrg(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())

	var req struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}
	if !validUsername.MatchString(req.Name) {
		writeError(w, http.StatusBadRequest, "bad_request", "Name must be 3-40 lowercase alphanumeric characters or hyphens")
		return
	}

	exists, err := h.owners.NameExists(req.Name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if exists {
		writeError(w, http.StatusConflict, "conflict", "Name is already taken")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	id := ulid.Make().String()

	owner := &models.Owner{
		ID:        id,
		Name:      req.Name,
		Type:      "org",
		CreatedAt: now,
	}
	if err := h.owners.Create(owner); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to create owner")
		return
	}

	org := &models.Organization{
		ID:          id,
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := h.orgs.Create(org); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to create organization")
		return
	}

	if err := h.orgs.AddMember(id, cu.ID, "owner"); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to add owner member")
		return
	}

	h.activities.Log(cu.ID, nil, "create_org", map[string]any{"org": req.Name})

	writeJSON(w, http.StatusCreated, org)
}

func (h *Handler) GetOrg(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "org")
	org, err := h.orgs.GetByName(name)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Organization not found")
		return
	}
	writeJSON(w, http.StatusOK, org)
}

func (h *Handler) UpdateOrg(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	name := chi.URLParam(r, "org")

	org, err := h.orgs.GetByName(name)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Organization not found")
		return
	}

	role, err := h.orgs.GetMemberRole(org.ID, cu.ID)
	if err != nil || (role != "owner" && role != "admin" && !cu.IsAdmin) {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	var req struct {
		DisplayName *string `json:"display_name"`
		Description *string `json:"description"`
		AvatarURL   *string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	if req.DisplayName != nil {
		org.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		org.Description = *req.Description
	}
	if req.AvatarURL != nil {
		org.AvatarURL = *req.AvatarURL
	}
	org.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := h.orgs.Update(org); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to update organization")
		return
	}
	writeJSON(w, http.StatusOK, org)
}

func (h *Handler) ListOrgMembers(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "org")
	org, err := h.orgs.GetByName(name)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Organization not found")
		return
	}

	members, err := h.orgs.ListMembers(org.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if members == nil {
		members = []models.OrgMember{}
	}
	writeJSON(w, http.StatusOK, members)
}

func (h *Handler) AddOrgMember(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	orgName := chi.URLParam(r, "org")

	org, err := h.orgs.GetByName(orgName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Organization not found")
		return
	}

	role, err := h.orgs.GetMemberRole(org.ID, cu.ID)
	if err != nil || (role != "owner" && role != "admin" && !cu.IsAdmin) {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	var req struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}
	if req.Role != "owner" && req.Role != "admin" && req.Role != "member" {
		writeError(w, http.StatusBadRequest, "bad_request", "Role must be owner, admin, or member")
		return
	}

	user, err := h.users.GetByUsername(req.Username)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "User not found")
		return
	}

	if err := h.orgs.AddMember(org.ID, user.ID, req.Role); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to add member")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemoveOrgMember(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	orgName := chi.URLParam(r, "org")
	username := chi.URLParam(r, "username")

	org, err := h.orgs.GetByName(orgName)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Organization not found")
		return
	}

	role, err := h.orgs.GetMemberRole(org.ID, cu.ID)
	if err != nil || (role != "owner" && !cu.IsAdmin) {
		writeError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}

	user, err := h.users.GetByUsername(username)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "User not found")
		return
	}

	if err := h.orgs.RemoveMember(org.ID, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to remove member")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListUserOrgs(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	orgs, err := h.orgs.ListByUser(cu.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if orgs == nil {
		orgs = []models.Organization{}
	}
	writeJSON(w, http.StatusOK, orgs)
}

func (h *Handler) ListOrgLooms(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "org")
	org, err := h.orgs.GetByName(name)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Organization not found")
		return
	}

	looms, total, err := h.looms.ListPublicByOwner(org.Name, 50, 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": looms, "total": total})
}
