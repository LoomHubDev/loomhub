package models

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	DisplayName  string `json:"display_name"`
	Bio          string `json:"bio"`
	AvatarURL    string `json:"avatar_url"`
	PasswordHash string `json:"-"`
	IsAdmin      bool   `json:"is_admin"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type AccessToken struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	Name       string  `json:"name"`
	TokenHash  string  `json:"-"`
	Scopes     string  `json:"scopes"`
	ExpiresAt  *string `json:"expires_at"`
	LastUsedAt *string `json:"last_used_at"`
	CreatedAt  string  `json:"created_at"`
}
