package api

import (
	"net/http"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetActivityFeed(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	activities, err := h.activities.ListByUser(cu.ID, 30)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if activities == nil {
		activities = []models.Activity{}
	}
	writeJSON(w, http.StatusOK, activities)
}

func (h *Handler) GetUserActivity(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.users.GetByUsername(username)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "User not found")
		return
	}
	activities, err := h.activities.ListByUser(user.ID, 30)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if activities == nil {
		activities = []models.Activity{}
	}
	writeJSON(w, http.StatusOK, activities)
}

func (h *Handler) GetGlobalActivity(w http.ResponseWriter, r *http.Request) {
	activities, err := h.activities.ListGlobal(50)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if activities == nil {
		activities = []models.Activity{}
	}
	writeJSON(w, http.StatusOK, activities)
}
