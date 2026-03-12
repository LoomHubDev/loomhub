package sync

import "encoding/json"

// Wire types — must match Loom's Go types

type NegotiateRequest struct {
	ProjectID string            `json:"project_id"`
	Streams   []StreamSyncState `json:"streams"`
}

type StreamSyncState struct {
	StreamID string `json:"stream_id"`
	Name     string `json:"name"`
	HeadSeq  int64  `json:"head_seq"`
}

type NegotiateResponse struct {
	CommonSeqs map[string]int64 `json:"common_seqs"`
	ServerSeqs map[string]int64 `json:"server_seqs"`
	NeedsPush  bool             `json:"needs_push"`
	NeedsPull  bool             `json:"needs_pull"`
}

type PushRequest struct {
	ProjectID  string      `json:"project_id"`
	StreamID   string      `json:"stream_id"`
	FromSeq    int64       `json:"from_seq"`
	Operations []Operation `json:"operations"`
	Objects    []ObjectData `json:"objects"`
}

type Operation struct {
	ID        string          `json:"id"`
	Seq       int64           `json:"seq"`
	StreamID  string          `json:"stream_id"`
	SpaceID   string          `json:"space_id"`
	EntityID  string          `json:"entity_id"`
	Type      string          `json:"type"`
	Path      string          `json:"path"`
	Delta     json.RawMessage `json:"delta,omitempty"`
	ObjectRef string          `json:"object_ref,omitempty"`
	ParentSeq int64           `json:"parent_seq"`
	Author    string          `json:"author"`
	Timestamp string          `json:"timestamp"`
	Meta      json.RawMessage `json:"meta,omitempty"`
}

type ObjectData struct {
	Hash    string `json:"hash"`
	Content []byte `json:"content"`
}

type PushResponse struct {
	OK         bool   `json:"ok"`
	Applied    int    `json:"applied"`
	ServerHead int64  `json:"server_head"`
	Error      string `json:"error,omitempty"`
}

type PullRequest struct {
	ProjectID string `json:"project_id"`
	StreamID  string `json:"stream_id"`
	FromSeq   int64  `json:"from_seq"`
}

type PullResponse struct {
	Operations []Operation  `json:"operations"`
	Objects    []ObjectData `json:"objects"`
	ServerHead int64        `json:"server_head"`
}
