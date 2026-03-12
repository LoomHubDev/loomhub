package sync

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type Stream struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	HeadSeq   int64  `json:"head_seq"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Entity struct {
	ID        string `json:"id"`
	SpaceID   string `json:"space_id"`
	Path      string `json:"path"`
	ObjectRef string `json:"object_ref"`
	Size      int64  `json:"size"`
	ModTime   string `json:"mod_time"`
}

type Checkpoint struct {
	ID        string          `json:"id"`
	Seq       int64           `json:"seq"`
	StreamID  string          `json:"stream_id"`
	SpaceID   string          `json:"space_id"`
	EntityID  string          `json:"entity_id"`
	Type      string          `json:"type"`
	Path      string          `json:"path"`
	ObjectRef string          `json:"object_ref"`
	ParentSeq int64           `json:"parent_seq"`
	Author    string          `json:"author"`
	Timestamp string          `json:"timestamp"`
	Meta      json.RawMessage `json:"meta,omitempty"`
}

// resolveStreamID resolves a stream name or ID to a stream ID.
func resolveStreamID(db *sql.DB, nameOrID string) string {
	if nameOrID == "" {
		return ""
	}
	var id string
	// Try as ID first
	err := db.QueryRow("SELECT id FROM streams WHERE id = ?", nameOrID).Scan(&id)
	if err == nil {
		return id
	}
	// Try as name
	err = db.QueryRow("SELECT id FROM streams WHERE name = ?", nameOrID).Scan(&id)
	if err == nil {
		return id
	}
	return nameOrID
}

func (h *Handler) ListStreams(loom *models.Loom) ([]Stream, error) {
	db, err := h.loomDBs.Open(loom.DiskPath)
	if err != nil {
		return []Stream{}, nil // loom DB not yet created = no streams
	}

	rows, err := db.Query("SELECT id, name, head_seq, status, created_at, updated_at FROM streams ORDER BY name ASC")
	if err != nil {
		return []Stream{}, nil
	}
	defer rows.Close()

	var streams []Stream
	for rows.Next() {
		var s Stream
		if err := rows.Scan(&s.ID, &s.Name, &s.HeadSeq, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		streams = append(streams, s)
	}
	if streams == nil {
		streams = []Stream{}
	}
	return streams, nil
}

func (h *Handler) ListEntities(loom *models.Loom, streamNameOrID string) ([]Entity, error) {
	db, err := h.loomDBs.Open(loom.DiskPath)
	if err != nil {
		return []Entity{}, nil
	}

	streamID := resolveStreamID(db, streamNameOrID)

	// Get entities that belong to operations on this stream
	// If streamID is empty, return all entities
	var rows *sql.Rows
	if streamID != "" {
		rows, err = db.Query(`
			SELECT DISTINCT e.id, e.space_id, e.path, e.object_ref, e.size, e.mod_time
			FROM entities e
			WHERE e.space_id IN (
				SELECT DISTINCT space_id FROM operations WHERE stream_id = ?
			)
			ORDER BY e.path ASC`, streamID)
	} else {
		rows, err = db.Query(`
			SELECT id, space_id, path, object_ref, size, mod_time
			FROM entities ORDER BY path ASC`)
	}
	if err != nil {
		return []Entity{}, nil
	}
	defer rows.Close()

	var entities []Entity
	for rows.Next() {
		var e Entity
		if err := rows.Scan(&e.ID, &e.SpaceID, &e.Path, &e.ObjectRef, &e.Size, &e.ModTime); err != nil {
			continue
		}
		entities = append(entities, e)
	}
	if entities == nil {
		entities = []Entity{}
	}
	return entities, nil
}

type CheckpointFilter struct {
	Stream string
	Type   string
	Author string
	Path   string
	Limit  int
	Offset int
}

func (h *Handler) ListCheckpoints(loom *models.Loom, f CheckpointFilter) ([]Checkpoint, int, error) {
	db, err := h.loomDBs.Open(loom.DiskPath)
	if err != nil {
		return []Checkpoint{}, 0, nil
	}

	streamID := resolveStreamID(db, f.Stream)

	where := ""
	args := []any{}
	addFilter := func(clause string, val any) {
		if where == "" {
			where = " WHERE " + clause
		} else {
			where += " AND " + clause
		}
		args = append(args, val)
	}
	if streamID != "" {
		addFilter("stream_id = ?", streamID)
	}
	if f.Type != "" {
		addFilter("type = ?", f.Type)
	}
	if f.Author != "" {
		addFilter("author LIKE ?", "%"+f.Author+"%")
	}
	if f.Path != "" {
		addFilter("path LIKE ?", "%"+f.Path+"%")
	}

	var total int
	db.QueryRow("SELECT COUNT(*) FROM operations"+where, args...).Scan(&total)

	query := `SELECT id, seq, stream_id, space_id, entity_id, type, path, object_ref, parent_seq, author, timestamp, meta
		FROM operations` + where + ` ORDER BY seq DESC LIMIT ? OFFSET ?`
	args = append(args, f.Limit, f.Offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return []Checkpoint{}, 0, nil
	}
	defer rows.Close()

	var checkpoints []Checkpoint
	for rows.Next() {
		var c Checkpoint
		var meta []byte
		if err := rows.Scan(&c.ID, &c.Seq, &c.StreamID, &c.SpaceID, &c.EntityID, &c.Type, &c.Path, &c.ObjectRef, &c.ParentSeq, &c.Author, &c.Timestamp, &meta); err != nil {
			continue
		}
		if len(meta) > 0 {
			c.Meta = json.RawMessage(meta)
		}
		checkpoints = append(checkpoints, c)
	}
	if checkpoints == nil {
		checkpoints = []Checkpoint{}
	}
	return checkpoints, total, nil
}

func (h *Handler) GetEntityContent(loom *models.Loom, objectRef string) ([]byte, error) {
	if objectRef == "" {
		return nil, fmt.Errorf("no object ref")
	}
	return h.objects.Read(objectRef)
}
