package sync

import (
	"database/sql"
	"os"
	"testing"

	"github.com/LoomHubDev/loomhub/internal/database"
	"github.com/LoomHubDev/loomhub/internal/models"
	_ "modernc.org/sqlite"
)

func setupTestSync(t *testing.T) (*Handler, *models.Loom) {
	t.Helper()
	dir := t.TempDir()

	// Hub database
	hubDB, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(hubDB); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { hubDB.Close() })

	// Create owner + user + loom in hub
	hubDB.Exec("INSERT INTO owners (id, name, type, created_at) VALUES ('u1', 'testuser', 'user', datetime('now'))")
	hubDB.Exec(`INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
		VALUES ('u1', 'testuser', 'test@test.com', 'hash', datetime('now'), datetime('now'))`)
	hubDB.Exec(`INSERT INTO looms (id, owner_id, name, visibility, default_stream, disk_path, created_at, updated_at)
		VALUES ('l1', 'u1', 'test-loom', 'public', 'main', 'looms/testuser/test-loom', datetime('now'), datetime('now'))`)

	loom := &models.Loom{
		ID:       "l1",
		OwnerID:  "u1",
		Name:     "test-loom",
		DiskPath: "looms/testuser/test-loom",
	}

	// Initialize per-loom database
	handler := NewHandler(hubDB, dir)
	t.Cleanup(func() { handler.Close() })

	loomDBPath := dir + "/looms/testuser/test-loom"
	initTestLoomDB(t, loomDBPath)

	return handler, loom
}

func initTestLoomDB(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", path+"/loom.db?_pragma=journal_mode(wal)&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS operations (
			id TEXT PRIMARY KEY, seq INTEGER NOT NULL, stream_id TEXT NOT NULL,
			space_id TEXT NOT NULL, entity_id TEXT NOT NULL, type TEXT NOT NULL,
			path TEXT NOT NULL, delta BLOB, object_ref TEXT, parent_seq INTEGER,
			author TEXT NOT NULL, timestamp TEXT NOT NULL, meta TEXT NOT NULL DEFAULT '{}'
		);
		CREATE INDEX IF NOT EXISTS idx_ops_seq ON operations(seq);
		CREATE INDEX IF NOT EXISTS idx_ops_stream ON operations(stream_id, seq);
		CREATE TABLE IF NOT EXISTS streams (
			id TEXT PRIMARY KEY, name TEXT NOT NULL UNIQUE, head_seq INTEGER NOT NULL DEFAULT 0,
			parent_id TEXT, status TEXT NOT NULL DEFAULT 'active',
			created_at TEXT NOT NULL, updated_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS entities (
			id TEXT NOT NULL, space_id TEXT NOT NULL, path TEXT NOT NULL,
			kind TEXT NOT NULL DEFAULT 'file', object_ref TEXT,
			size INTEGER NOT NULL DEFAULT 0, mod_time TEXT,
			PRIMARY KEY (space_id, id)
		);
		CREATE INDEX IF NOT EXISTS idx_entities_path ON entities(space_id, path);
		CREATE TABLE IF NOT EXISTS objects (
			hash TEXT PRIMARY KEY, size INTEGER NOT NULL DEFAULT 0,
			compressed BOOLEAN NOT NULL DEFAULT FALSE
		);
		CREATE TABLE IF NOT EXISTS checkpoints (
			id TEXT PRIMARY KEY, stream_id TEXT NOT NULL, seq INTEGER NOT NULL,
			title TEXT NOT NULL, summary TEXT NOT NULL DEFAULT '', author TEXT NOT NULL,
			source TEXT NOT NULL DEFAULT 'manual', tags TEXT NOT NULL DEFAULT '{}',
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS metadata (key TEXT PRIMARY KEY, value TEXT NOT NULL);
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNegotiateEmptyLoom(t *testing.T) {
	h, loom := setupTestSync(t)

	resp, err := h.Negotiate(loom, &NegotiateRequest{
		ProjectID: "test",
		Streams: []StreamSyncState{
			{StreamID: "s1", Name: "main", HeadSeq: 100},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if !resp.NeedsPush {
		t.Error("NeedsPush should be true for empty server")
	}
	if resp.NeedsPull {
		t.Error("NeedsPull should be false for empty server")
	}
}

func TestSendAndReceive(t *testing.T) {
	h, loom := setupTestSync(t)

	content := []byte("package main\n\nfunc main() {}\n")
	objHash := HashContent(content)

	// Send operations
	pushResp, err := h.Send(loom, &PushRequest{
		ProjectID: "test",
		StreamID:  "s1",
		FromSeq:   0,
		Operations: []Operation{
			{
				ID: "op1", Seq: 1, StreamID: "s1", SpaceID: "code",
				EntityID: "main.go", Type: "create", Path: "main.go",
				ObjectRef: objHash, Author: "testuser", Timestamp: "2026-03-12T10:00:00Z",
			},
			{
				ID: "op2", Seq: 2, StreamID: "s1", SpaceID: "code",
				EntityID: "main.go", Type: "modify", Path: "main.go",
				ObjectRef: objHash, Author: "testuser", Timestamp: "2026-03-12T10:01:00Z",
			},
		},
		Objects: []ObjectData{
			{Hash: objHash, Content: content},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !pushResp.OK {
		t.Fatalf("Send failed: %s", pushResp.Error)
	}
	if pushResp.Applied != 2 {
		t.Errorf("Applied = %d, want 2", pushResp.Applied)
	}
	if pushResp.ServerHead != 2 {
		t.Errorf("ServerHead = %d, want 2", pushResp.ServerHead)
	}

	// Verify object exists in shared store
	if !h.objects.Has(objHash) {
		t.Error("object should be in shared store")
	}

	// Negotiate after send
	negResp, _ := h.Negotiate(loom, &NegotiateRequest{
		ProjectID: "test",
		Streams: []StreamSyncState{
			{StreamID: "s1", Name: "s1", HeadSeq: 2},
		},
	})
	if negResp.NeedsPush {
		t.Error("NeedsPush should be false after send with matching head")
	}

	// Receive from seq 0 (another client)
	pullResp, err := h.Receive(loom, &PullRequest{
		ProjectID: "test",
		StreamID:  "s1",
		FromSeq:   0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(pullResp.Operations) != 2 {
		t.Errorf("Operations = %d, want 2", len(pullResp.Operations))
	}
	if len(pullResp.Objects) != 1 {
		t.Errorf("Objects = %d, want 1", len(pullResp.Objects))
	}
	if pullResp.Objects[0].Hash != objHash {
		t.Errorf("Object hash = %s, want %s", pullResp.Objects[0].Hash, objHash)
	}

	// Receive from seq 1 (partial)
	pullResp2, _ := h.Receive(loom, &PullRequest{
		ProjectID: "test",
		StreamID:  "s1",
		FromSeq:   1,
	})
	if len(pullResp2.Operations) != 1 {
		t.Errorf("Partial receive operations = %d, want 1", len(pullResp2.Operations))
	}
}

func TestSendIdempotent(t *testing.T) {
	h, loom := setupTestSync(t)

	content := []byte("test content")
	objHash := HashContent(content)

	req := &PushRequest{
		ProjectID: "test",
		StreamID:  "s1",
		FromSeq:   0,
		Operations: []Operation{
			{
				ID: "op1", Seq: 1, StreamID: "s1", SpaceID: "code",
				EntityID: "file.txt", Type: "create", Path: "file.txt",
				ObjectRef: objHash, Author: "testuser", Timestamp: "2026-03-12T10:00:00Z",
			},
		},
		Objects: []ObjectData{{Hash: objHash, Content: content}},
	}

	// Send twice — should be idempotent (INSERT OR IGNORE)
	h.Send(loom, req)
	resp, _ := h.Send(loom, req)
	if !resp.OK {
		t.Fatalf("second send failed: %s", resp.Error)
	}
}

