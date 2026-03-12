package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

type ContextUser struct {
	ID       string
	Username string
	IsAdmin  bool
}

func Middleware(secret string, db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := extractUser(r, secret, db)
			if user != nil {
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUser(r.Context()) == nil {
			http.Error(w, `{"error":{"code":"unauthorized","message":"Authentication required"}}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUser(ctx context.Context) *ContextUser {
	u, _ := ctx.Value(UserContextKey).(*ContextUser)
	return u
}

func extractUser(r *http.Request, secret string, db *sql.DB) *ContextUser {
	// Try JWT from cookie
	if cookie, err := r.Cookie("session"); err == nil {
		claims, err := ValidateToken(cookie.Value, secret)
		if err == nil {
			return &ContextUser{
				ID:       claims.UserID,
				Username: claims.Username,
				IsAdmin:  claims.IsAdmin,
			}
		}
	}

	// Try Bearer token (JWT or access token)
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		return nil
	}

	// Try as JWT first
	claims, err := ValidateToken(token, secret)
	if err == nil {
		return &ContextUser{
			ID:       claims.UserID,
			Username: claims.Username,
			IsAdmin:  claims.IsAdmin,
		}
	}

	// Try as access token
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	var userID, username string
	var isAdmin bool
	err = db.QueryRow(`
		SELECT u.id, u.username, u.is_admin
		FROM access_tokens t
		JOIN users u ON u.id = t.user_id
		WHERE t.token_hash = ?
		AND (t.expires_at IS NULL OR t.expires_at > datetime('now'))
	`, tokenHash).Scan(&userID, &username, &isAdmin)
	if err != nil {
		return nil
	}

	// Update last_used_at
	db.Exec("UPDATE access_tokens SET last_used_at = datetime('now') WHERE token_hash = ?", tokenHash)

	return &ContextUser{
		ID:       userID,
		Username: username,
		IsAdmin:  isAdmin,
	}
}
