package models

type Loom struct {
	ID              string  `json:"id" gorm:"primaryKey"`
	OwnerID         string  `json:"owner_id" gorm:"index;uniqueIndex:idx_looms_owner_name"`
	Name            string  `json:"name" gorm:"index;uniqueIndex:idx_looms_owner_name"`
	Description     string  `json:"description"`
	Visibility      string  `json:"visibility" gorm:"index;default:public"`
	DefaultStream   string  `json:"default_stream" gorm:"default:main"`
	DiskPath        string  `json:"-"`
	SizeBytes       int64   `json:"size_bytes" gorm:"default:0"`
	CheckpointCount int     `json:"checkpoint_count" gorm:"default:0"`
	StreamCount     int     `json:"stream_count" gorm:"default:0"`
	PinCount        int     `json:"pin_count" gorm:"index;default:0"`
	SpinCount       int     `json:"spin_count" gorm:"default:0"`
	SpunFrom        *string `json:"spun_from"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	SyncedAt        *string `json:"synced_at" gorm:"index"`

	// Association
	Owner *Owner `json:"-" gorm:"foreignKey:OwnerID"`

	// Computed (not stored)
	OwnerName string `json:"owner,omitempty" gorm:"-"`
	FullName  string `json:"full_name,omitempty" gorm:"-"`
}

func (l *Loom) FillOwnerFields() {
	if l.Owner != nil {
		l.OwnerName = l.Owner.Name
		l.FullName = l.Owner.Name + "/" + l.Name
	}
}

type Collaborator struct {
	LoomID     string `json:"loom_id" gorm:"primaryKey"`
	UserID     string `json:"user_id" gorm:"primaryKey"`
	Permission string `json:"permission" gorm:"default:read"`
	CreatedAt  string `json:"created_at"`

	// Association
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

func (Collaborator) TableName() string { return "loom_collaborators" }

type CollaboratorWithUser struct {
	UserID      string `json:"user_id"`
	Permission  string `json:"permission"`
	CreatedAt   string `json:"created_at"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type Pin struct {
	UserID    string `json:"user_id" gorm:"primaryKey"`
	LoomID    string `json:"loom_id" gorm:"primaryKey;index"`
	CreatedAt string `json:"created_at"`
}

type WRLabel struct {
	WRID    string `json:"wr_id" gorm:"primaryKey;column:wr_id"`
	LabelID string `json:"label_id" gorm:"primaryKey"`
}

func (WRLabel) TableName() string { return "wr_labels" }

type ObjectRef struct {
	Hash   string `json:"hash" gorm:"primaryKey"`
	LoomID string `json:"loom_id" gorm:"primaryKey;index"`
	Size   int64  `json:"size" gorm:"default:0"`
}
