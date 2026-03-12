package sync

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type Handler struct {
	hubDB   *sql.DB
	objects *ObjectStore
	loomDBs *LoomDBManager
}

func NewHandler(hubDB *sql.DB, dataDir string) *Handler {
	return &Handler{
		hubDB:   hubDB,
		objects: NewObjectStore(dataDir),
		loomDBs: NewLoomDBManager(dataDir),
	}
}

func (h *Handler) Close() {
	h.loomDBs.Close()
}

func (h *Handler) Negotiate(loom *models.Loom, req *NegotiateRequest) (*NegotiateResponse, error) {
	db, err := h.loomDBs.Open(loom.DiskPath)
	if err != nil {
		return nil, fmt.Errorf("open loom db: %w", err)
	}

	resp := &NegotiateResponse{
		CommonSeqs: make(map[string]int64),
		ServerSeqs: make(map[string]int64),
	}

	// Get server's stream states
	rows, err := db.Query("SELECT id, name, head_seq FROM streams")
	if err != nil {
		// Empty loom — no streams yet
		for _, cs := range req.Streams {
			resp.CommonSeqs[cs.StreamID] = 0
			resp.ServerSeqs[cs.StreamID] = 0
			if cs.HeadSeq > 0 {
				resp.NeedsPush = true
			}
		}
		return resp, nil
	}
	defer rows.Close()

	serverStreams := make(map[string]int64) // stream_id → head_seq
	serverByName := make(map[string]string) // name → stream_id
	for rows.Next() {
		var id, name string
		var headSeq int64
		rows.Scan(&id, &name, &headSeq)
		serverStreams[id] = headSeq
		serverByName[name] = id
	}

	for _, cs := range req.Streams {
		serverSeq, exists := serverStreams[cs.StreamID]
		if !exists {
			// Try matching by name
			if sid, ok := serverByName[cs.Name]; ok {
				serverSeq = serverStreams[sid]
				exists = true
			}
		}

		if !exists {
			resp.CommonSeqs[cs.StreamID] = 0
			resp.ServerSeqs[cs.StreamID] = 0
			if cs.HeadSeq > 0 {
				resp.NeedsPush = true
			}
		} else {
			common := min(cs.HeadSeq, serverSeq)
			resp.CommonSeqs[cs.StreamID] = common
			resp.ServerSeqs[cs.StreamID] = serverSeq
			if cs.HeadSeq > serverSeq {
				resp.NeedsPush = true
			}
			if serverSeq > cs.HeadSeq {
				resp.NeedsPull = true
			}
		}
	}

	// Check for server streams the client doesn't know about
	clientStreams := make(map[string]bool)
	for _, cs := range req.Streams {
		clientStreams[cs.StreamID] = true
		clientStreams[cs.Name] = true
	}
	for id, seq := range serverStreams {
		if !clientStreams[id] && seq > 0 {
			resp.ServerSeqs[id] = seq
			resp.CommonSeqs[id] = 0
			resp.NeedsPull = true
		}
	}

	return resp, nil
}

func (h *Handler) Send(loom *models.Loom, req *PushRequest) (*PushResponse, error) {
	// Lock this loom for writing
	h.loomDBs.Lock(loom.ID)
	defer h.loomDBs.Unlock(loom.ID)

	db, err := h.loomDBs.Open(loom.DiskPath)
	if err != nil {
		return nil, fmt.Errorf("open loom db: %w", err)
	}

	// Store objects in shared store
	for _, obj := range req.Objects {
		if err := h.objects.Write(obj.Hash, obj.Content); err != nil {
			return &PushResponse{OK: false, Error: fmt.Sprintf("store object %s: %v", obj.Hash, err)}, nil
		}
		// Track reference in hub db
		h.hubDB.Exec(
			"INSERT OR IGNORE INTO object_refs (hash, loom_id, size) VALUES (?, ?, ?)",
			obj.Hash, loom.ID, len(obj.Content),
		)
		// Track in loom's object index
		db.Exec(
			"INSERT OR IGNORE INTO objects (hash, size) VALUES (?, ?)",
			obj.Hash, len(obj.Content),
		)
	}

	// Apply operations in a transaction
	tx, err := db.Begin()
	if err != nil {
		return &PushResponse{OK: false, Error: "begin transaction: " + err.Error()}, nil
	}

	applied := 0
	var maxSeq int64

	for _, op := range req.Operations {
		deltaBytes, _ := json.Marshal(op.Delta)
		metaBytes, _ := json.Marshal(op.Meta)

		_, err := tx.Exec(`
			INSERT OR IGNORE INTO operations (id, seq, stream_id, space_id, entity_id, type, path, delta, object_ref, parent_seq, author, timestamp, meta)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			op.ID, op.Seq, op.StreamID, op.SpaceID, op.EntityID, op.Type, op.Path,
			deltaBytes, op.ObjectRef, op.ParentSeq, op.Author, op.Timestamp, metaBytes,
		)
		if err != nil {
			tx.Rollback()
			return &PushResponse{OK: false, Error: fmt.Sprintf("write op %s: %v", op.ID, err)}, nil
		}
		applied++
		if op.Seq > maxSeq {
			maxSeq = op.Seq
		}

		// Update entity state
		if op.Type == "delete" {
			tx.Exec("DELETE FROM entities WHERE space_id = ? AND id = ?", op.SpaceID, op.EntityID)
		} else {
			tx.Exec(`
				INSERT INTO entities (id, space_id, path, object_ref, size, mod_time)
				VALUES (?, ?, ?, ?, 0, ?)
				ON CONFLICT(space_id, id) DO UPDATE SET
					path = excluded.path, object_ref = excluded.object_ref, mod_time = excluded.mod_time`,
				op.EntityID, op.SpaceID, op.Path, op.ObjectRef, op.Timestamp,
			)
		}
	}

	// Ensure stream exists and update head
	if req.StreamID != "" && maxSeq > 0 {
		tx.Exec(`
			INSERT INTO streams (id, name, head_seq, status, created_at, updated_at)
			VALUES (?, ?, ?, 'active', datetime('now'), datetime('now'))
			ON CONFLICT(id) DO UPDATE SET head_seq = MAX(head_seq, excluded.head_seq), updated_at = datetime('now')`,
			req.StreamID, req.StreamID, maxSeq,
		)
	}

	if err := tx.Commit(); err != nil {
		return &PushResponse{OK: false, Error: "commit: " + err.Error()}, nil
	}

	// Update hub-level stats
	h.hubDB.Exec("UPDATE looms SET synced_at = datetime('now'), updated_at = datetime('now') WHERE id = ?", loom.ID)

	return &PushResponse{
		OK:         true,
		Applied:    applied,
		ServerHead: maxSeq,
	}, nil
}

func (h *Handler) Receive(loom *models.Loom, req *PullRequest) (*PullResponse, error) {
	db, err := h.loomDBs.Open(loom.DiskPath)
	if err != nil {
		return nil, fmt.Errorf("open loom db: %w", err)
	}

	// Get operations after fromSeq for the given stream
	rows, err := db.Query(`
		SELECT id, seq, stream_id, space_id, entity_id, type, path, delta, object_ref, parent_seq, author, timestamp, meta
		FROM operations
		WHERE stream_id = ? AND seq > ?
		ORDER BY seq ASC`,
		req.StreamID, req.FromSeq,
	)
	if err != nil {
		return &PullResponse{Operations: []Operation{}, Objects: []ObjectData{}, ServerHead: req.FromSeq}, nil
	}
	defer rows.Close()

	var ops []Operation
	objectRefs := make(map[string]bool)

	for rows.Next() {
		var op Operation
		var delta, meta []byte
		err := rows.Scan(
			&op.ID, &op.Seq, &op.StreamID, &op.SpaceID, &op.EntityID,
			&op.Type, &op.Path, &delta, &op.ObjectRef, &op.ParentSeq,
			&op.Author, &op.Timestamp, &meta,
		)
		if err != nil {
			continue
		}
		if len(delta) > 0 {
			op.Delta = json.RawMessage(delta)
		}
		if len(meta) > 0 {
			op.Meta = json.RawMessage(meta)
		}
		ops = append(ops, op)
		if op.ObjectRef != "" {
			objectRefs[op.ObjectRef] = true
		}
	}

	// Read referenced objects from shared store
	var objects []ObjectData
	for hash := range objectRefs {
		content, err := h.objects.Read(hash)
		if err != nil {
			continue
		}
		objects = append(objects, ObjectData{Hash: hash, Content: content})
	}

	// Get server head for this stream
	var serverHead int64
	db.QueryRow("SELECT COALESCE(MAX(head_seq), 0) FROM streams WHERE id = ?", req.StreamID).Scan(&serverHead)

	if ops == nil {
		ops = []Operation{}
	}
	if objects == nil {
		objects = []ObjectData{}
	}

	return &PullResponse{
		Operations: ops,
		Objects:    objects,
		ServerHead: serverHead,
	}, nil
}
