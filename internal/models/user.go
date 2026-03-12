package models

type User struct {
	ID           string `json:"id" gorm:"primaryKey"`
	Username     string `json:"username" gorm:"uniqueIndex"`
	Email        string `json:"email" gorm:"uniqueIndex"`
	DisplayName  string `json:"display_name"`
	Bio          string `json:"bio"`
	AvatarURL    string `json:"avatar_url"`
	PasswordHash string `json:"-"`
	IsAdmin      bool   `json:"is_admin"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type AccessToken struct {
	ID         string  `json:"id" gorm:"primaryKey"`
	UserID     string  `json:"user_id" gorm:"index"`
	Name       string  `json:"name"`
	TokenHash  string  `json:"-" gorm:"uniqueIndex"`
	Scopes     string  `json:"scopes"`
	ExpiresAt  *string `json:"expires_at"`
	LastUsedAt *string `json:"last_used_at"`
	CreatedAt  string  `json:"created_at"`
}
