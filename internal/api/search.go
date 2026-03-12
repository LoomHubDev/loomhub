package api

import (
	"net/http"
	"strconv"

	"github.com/LoomHubDev/loomhub/internal/models"
)

func (h *Handler) SearchLooms(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, http.StatusOK, []models.Loom{})
		return
	}
	looms, err := h.search.SearchLooms(q, 20)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Search failed")
		return
	}
	if looms == nil {
		looms = []models.Loom{}
	}
	writeJSON(w, http.StatusOK, looms)
}

func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, http.StatusOK, []models.User{})
		return
	}
	users, err := h.search.SearchUsers(q, 20)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Search failed")
		return
	}
	if users == nil {
		users = []models.User{}
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *Handler) Explore(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "recent"
	}
	limit := 20
	offset := 0
	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && o > 0 {
		offset = o
	}

	looms, total, err := h.search.ExploreLooms(sort, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if looms == nil {
		looms = []models.Loom{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": looms, "total": total})
}
