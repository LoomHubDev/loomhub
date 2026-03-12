package models

type Organization struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type OrgMember struct {
	OrgID     string `json:"org_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"` // owner, admin, member
	CreatedAt string `json:"created_at"`

	// Joined
	Username string `json:"username,omitempty"`
}
