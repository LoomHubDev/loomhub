package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

var validUsername = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,38}[a-z0-9]$`)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if !h.cfg.Auth.RegistrationOpen {
		writeError(w, http.StatusForbidden, "forbidden", "Registration is closed")
		return
	}

	var req struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if !validUsername.MatchString(req.Username) {
		writeError(w, http.StatusBadRequest, "bad_request", "Username must be 3-40 lowercase alphanumeric characters or hyphens")
		return
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid email")
		return
	}
	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "bad_request", "Password must be at least 8 characters")
		return
	}

	exists, err := h.owners.NameExists(req.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	if exists {
		writeError(w, http.StatusConflict, "conflict", "Username is already taken")
		return
	}

	hash, err := auth.HashPassword(req.Password, h.cfg.Auth.BcryptCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to hash password")
		return
	}

	// First user becomes admin
	userCount, _ := h.users.Count()
	isAdmin := userCount == 0

	now := time.Now().UTC().Format(time.RFC3339)
	id := ulid.Make().String()

	owner := &models.Owner{
		ID:        id,
		Name:      req.Username,
		Type:      "user",
		CreatedAt: now,
	}
	if err := h.owners.Create(owner); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to create owner")
		return
	}

	user := &models.User{
		ID:           id,
		Username:     req.Username,
		Email:        req.Email,
		DisplayName:  req.DisplayName,
		PasswordHash: hash,
		IsAdmin:      isAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := h.users.Create(user); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to create user")
		return
	}

	expiry, _ := time.ParseDuration(h.cfg.Auth.JWTExpiry)
	token, _ := auth.GenerateToken(user.ID, user.Username, user.IsAdmin, h.cfg.Auth.JWTSecret, expiry)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(expiry.Seconds()),
	})

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"display_name": user.DisplayName,
		"is_admin":     user.IsAdmin,
		"created_at":   user.CreatedAt,
		"token":        token,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"` // username or email
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	req.Login = strings.ToLower(strings.TrimSpace(req.Login))

	var user *models.User
	var err error

	if strings.Contains(req.Login, "@") {
		user, err = h.users.GetByEmail(req.Login)
	} else {
		user, err = h.users.GetByUsername(req.Login)
	}

	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Invalid credentials")
		return
	}

	if err := auth.CheckPassword(user.PasswordHash, req.Password); err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Invalid credentials")
		return
	}

	expiry, _ := time.ParseDuration(h.cfg.Auth.JWTExpiry)
	token, _ := auth.GenerateToken(user.ID, user.Username, user.IsAdmin, h.cfg.Auth.JWTSecret, expiry)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(expiry.Seconds()),
	})

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"display_name": user.DisplayName,
		"is_admin":     user.IsAdmin,
		"token":        token,
	})
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	user, err := h.users.GetByID(cu.ID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "User not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	user, err := h.users.GetByID(cu.ID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "User not found")
		return
	}

	var req struct {
		DisplayName *string `json:"display_name"`
		Bio         *string `json:"bio"`
		AvatarURL   *string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
		return
	}

	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}
	user.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := h.users.Update(user); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to update user")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.users.GetByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) ListCurrentUserLooms(w http.ResponseWriter, r *http.Request) {
	cu := auth.GetUser(r.Context())
	looms, total, err := h.looms.ListByOwner(cu.ID, 50, 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": looms,
		"total": total,
	})
}

func (h *Handler) ListUserLooms(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	looms, total, err := h.looms.ListPublicByOwner(username, 50, 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Database error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": looms,
		"total": total,
	})
}
