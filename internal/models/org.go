package models

type Organization struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"uniqueIndex"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type OrgMember struct {
	OrgID     string `json:"org_id" gorm:"primaryKey"`
	UserID    string `json:"user_id" gorm:"primaryKey;index"`
	Role      string `json:"role" gorm:"default:member"`
	CreatedAt string `json:"created_at"`

	// Association
	User *User `json:"-" gorm:"foreignKey:UserID"`

	// Computed
	Username string `json:"username,omitempty" gorm:"-"`
}
