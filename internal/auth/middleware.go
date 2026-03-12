package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/LoomHubDev/loomhub/internal/models"
	"gorm.io/gorm"
)

type contextKey string

const UserContextKey contextKey = "user"

type ContextUser struct {
	ID       string
	Username string
	IsAdmin  bool
}

func Middleware(secret string, db *gorm.DB) func(http.Handler) http.Handler {
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

func extractUser(r *http.Request, secret string, db *gorm.DB) *ContextUser {
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
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
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
	now := time.Now().UTC().Format(time.RFC3339)

	var at models.AccessToken
	err = db.Where("token_hash = ? AND (expires_at IS NULL OR expires_at > ?)", tokenHash, now).
		First(&at).Error
	if err != nil {
		return nil
	}

	var user models.User
	if err = db.Where("id = ?", at.UserID).First(&user).Error; err != nil {
		return nil
	}

	// Update last_used_at
	db.Model(&at).Update("last_used_at", now)

	return &ContextUser{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}
}
