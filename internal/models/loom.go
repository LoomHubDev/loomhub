package models

type Loom struct {
	ID              string  `json:"id"`
	OwnerID         string  `json:"owner_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Visibility      string  `json:"visibility"`
	DefaultStream   string  `json:"default_stream"`
	DiskPath        string  `json:"-"`
	SizeBytes       int64   `json:"size_bytes"`
	CheckpointCount int     `json:"checkpoint_count"`
	StreamCount     int     `json:"stream_count"`
	PinCount        int     `json:"pin_count"`
	SpinCount       int     `json:"spin_count"`
	SpunFrom        *string `json:"spun_from"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	SyncedAt        *string `json:"synced_at"`

	// Joined fields (not always populated)
	OwnerName string `json:"owner,omitempty"`
	FullName  string `json:"full_name,omitempty"`
}

type Collaborator struct {
	LoomID     string `json:"loom_id"`
	UserID     string `json:"user_id"`
	Permission string `json:"permission"`
	CreatedAt  string `json:"created_at"`
}

type CollaboratorWithUser struct {
	UserID      string `json:"user_id"`
	Permission  string `json:"permission"`
	CreatedAt   string `json:"created_at"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}
