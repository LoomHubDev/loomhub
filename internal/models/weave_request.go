package models

type WeaveRequest struct {
	ID            string  `json:"id"`
	LoomID        string  `json:"loom_id"`
	Number        int     `json:"number"`
	AuthorID      string  `json:"author_id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	SourceStream  string  `json:"source_stream"`
	TargetStream  string  `json:"target_stream"`
	SourceHeadSeq int64   `json:"source_head_seq"`
	TargetHeadSeq int64   `json:"target_head_seq"`
	MergeBaseSeq  int64   `json:"merge_base_seq"`
	Status        string  `json:"status"` // open, woven, closed
	WovenBy       *string `json:"woven_by"`
	WovenAt       *string `json:"woven_at"`
	ClosedAt      *string `json:"closed_at"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`

	// Joined
	AuthorName string `json:"author,omitempty"`
}

type Comment struct {
	ID         string  `json:"id"`
	WRID       string  `json:"wr_id"`
	AuthorID   string  `json:"author_id"`
	Body       string  `json:"body"`
	SpaceID    *string `json:"space_id"`
	EntityPath *string `json:"entity_path"`
	LineNumber *int    `json:"line_number"`
	IsSystem   bool    `json:"is_system"`
	EventType  *string `json:"event_type"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`

	// Joined
	AuthorName string `json:"author,omitempty"`
}

type Review struct {
	ID         string `json:"id"`
	WRID       string `json:"wr_id"`
	ReviewerID string `json:"reviewer_id"`
	Status     string `json:"status"` // approved, changes_requested, commented
	Body       string `json:"body"`
	CreatedAt  string `json:"created_at"`

	// Joined
	ReviewerName string `json:"reviewer,omitempty"`
}

type Activity struct {
	ID        string  `json:"id"`
	ActorID   string  `json:"actor_id"`
	LoomID    *string `json:"loom_id"`
	Type      string  `json:"type"`
	Payload   string  `json:"payload"`
	CreatedAt string  `json:"created_at"`

	// Joined
	ActorName    string  `json:"actor,omitempty"`
	LoomFullName *string `json:"loom_full_name,omitempty"`
}
