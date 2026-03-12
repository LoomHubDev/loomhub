package api

import (
	"net/http"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/models"
)

func (h *Handler) AdminStats(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	if !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	var userCount, loomCount, wrCount, orgCount int64
	h.db.Model(&models.User{}).Count(&userCount)
	h.db.Model(&models.Loom{}).Count(&loomCount)
	h.db.Model(&models.WeaveRequest{}).Count(&wrCount)
	h.db.Model(&models.Organization{}).Count(&orgCount)

	writeJSON(w, http.StatusOK, map[string]any{
		"users":          userCount,
		"looms":          loomCount,
		"weave_requests": wrCount,
		"organizations":  orgCount,
	})
}

func (h *Handler) AdminListUsers(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	if !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	users, err := h.search.SearchUsers("", 100)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *Handler) AdminListLooms(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	if !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	looms, _, err := h.search.ExploreLooms("recent", 100, 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	writeJSON(w, http.StatusOK, looms)
}
