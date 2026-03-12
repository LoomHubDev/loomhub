package models

type WeaveRequest struct {
	ID            string  `json:"id" gorm:"primaryKey"`
	LoomID        string  `json:"loom_id" gorm:"index:idx_wr_loom;uniqueIndex:idx_wr_loom_number"`
	Number        int     `json:"number" gorm:"uniqueIndex:idx_wr_loom_number"`
	AuthorID      string  `json:"author_id" gorm:"index"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	SourceStream  string  `json:"source_stream"`
	TargetStream  string  `json:"target_stream"`
	SourceHeadSeq int64   `json:"source_head_seq"`
	TargetHeadSeq int64   `json:"target_head_seq"`
	MergeBaseSeq  int64   `json:"merge_base_seq"`
	Status        string  `json:"status" gorm:"index:idx_wr_loom;default:open"`
	WovenBy       *string `json:"woven_by"`
	WovenAt       *string `json:"woven_at"`
	ClosedAt      *string `json:"closed_at"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`

	// Association
	Author *User `json:"-" gorm:"foreignKey:AuthorID"`

	// Computed
	AuthorName string `json:"author,omitempty" gorm:"-"`
}

type Comment struct {
	ID         string  `json:"id" gorm:"primaryKey"`
	WRID       string  `json:"wr_id" gorm:"column:wr_id;index"`
	AuthorID   string  `json:"author_id" gorm:"index"`
	Body       string  `json:"body"`
	SpaceID    *string `json:"space_id"`
	EntityPath *string `json:"entity_path"`
	LineNumber *int    `json:"line_number"`
	IsSystem   bool    `json:"is_system" gorm:"default:false"`
	EventType  *string `json:"event_type"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`

	// Association
	Author *User `json:"-" gorm:"foreignKey:AuthorID"`

	// Computed
	AuthorName string `json:"author,omitempty" gorm:"-"`
}

func (Comment) TableName() string { return "wr_comments" }

type Review struct {
	ID         string `json:"id" gorm:"primaryKey"`
	WRID       string `json:"wr_id" gorm:"column:wr_id;index"`
	ReviewerID string `json:"reviewer_id"`
	Status     string `json:"status"`
	Body       string `json:"body"`
	CreatedAt  string `json:"created_at"`

	// Association
	Reviewer *User `json:"-" gorm:"foreignKey:ReviewerID"`

	// Computed
	ReviewerName string `json:"reviewer,omitempty" gorm:"-"`
}

type Activity struct {
	ID        string  `json:"id" gorm:"primaryKey"`
	ActorID   string  `json:"actor_id" gorm:"index"`
	LoomID    *string `json:"loom_id" gorm:"index"`
	Type      string  `json:"type" gorm:"column:type;index"`
	Payload   string  `json:"payload" gorm:"default:'{}'"`
	CreatedAt string  `json:"created_at"`

	// Associations
	Actor   *User `json:"-" gorm:"foreignKey:ActorID"`
	LoomRef *Loom `json:"-" gorm:"foreignKey:LoomID"`

	// Computed
	ActorName    string  `json:"actor,omitempty" gorm:"-"`
	LoomFullName *string `json:"loom_full_name,omitempty" gorm:"-"`
}

func (a *Activity) FillJoinedFields() {
	if a.Actor != nil {
		a.ActorName = a.Actor.Username
	}
	if a.LoomRef != nil && a.LoomRef.Owner != nil {
		fn := a.LoomRef.Owner.Name + "/" + a.LoomRef.Name
		a.LoomFullName = &fn
	}
}
