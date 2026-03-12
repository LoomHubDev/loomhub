package api

import (
	"net/http"

	"github.com/LoomHubDev/loomhub/internal/auth"
)

func (h *Handler) AdminStats(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	if !cu.IsAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	var userCount, loomCount, wrCount, orgCount int
	h.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	h.db.QueryRow("SELECT COUNT(*) FROM looms").Scan(&loomCount)
	h.db.QueryRow("SELECT COUNT(*) FROM weave_requests").Scan(&wrCount)
	h.db.QueryRow("SELECT COUNT(*) FROM organizations").Scan(&orgCount)

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
