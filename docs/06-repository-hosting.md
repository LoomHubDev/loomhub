# LoomHub — Repository Hosting

How LoomHub stores, serves, and manages Loom repositories on the server side.

---

## Repository Lifecycle

### Creation

When a user creates a repository via API or web UI:

1. Validate name (lowercase, alphanumeric + hyphens, 1-100 chars)
2. Check uniqueness within owner namespace
3. Insert record in hub database
4. Create disk directory: `{data_dir}/repos/{owner}/{repo}/`
5. Initialize empty Loom database: `{disk_path}/loom.db`
6. Create default stream ("main") in the Loom database
7. Log activity event

### First Push

The first `loom push` to an empty repo:

1. Client calls negotiate → server reports empty state
2. Client pushes all operations + objects
3. Server writes operations to repo's `loom.db`
4. Server writes objects to shared object store
5. Server updates stream head
6. Server updates repo stats (checkpoint_count, etc.)
7. Server dispatches `push` webhook

### Subsequent Pushes

1. Negotiate → find common ancestor
2. Client sends only new operations + new objects
3. Server appends operations (append-only)
4. Server stores new objects (deduplicates against shared store)
5. Update stream head, stats, `pushed_at`

### Forking

When a user forks a repository:

1. Create new repo record with `forked_from` reference
2. Copy the Loom database file (atomic copy)
3. Objects are already shared — no copying needed
4. Increment fork_count on source repo

This is fast because:
- Loom database is a single file (SQLite)
- Objects are content-addressed and shared

### Deletion

1. Mark repo as deleted (soft delete with grace period)
2. After grace period (24h), hard delete:
   - Remove repo record from hub database
   - Delete Loom database file
   - Run object GC (remove objects with zero references)

---

## Storage Architecture

### Per-Repo Loom Database

Each repository has its own SQLite database at `repos/{owner}/{repo}/loom.db`.

This database contains the exact same schema as a local Loom project:

```sql
-- Operations (append-only log)
-- Checkpoints (named points)
-- Streams (timeline management)
-- Entities (current state per space)
-- Objects (index — hash, size, compressed, ref_count)
-- Metadata (key-value config)
```

The server opens this database when handling sync requests and web UI queries. WAL mode enables concurrent reads during a push.

### Shared Object Store

All repos share a single content-addressed object store:

```
{data_dir}/objects/{hash[0:2]}/{hash}
```

Example:
```
objects/ab/abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789
objects/cd/cdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789ab
```

Benefits:
- **Deduplication** — Same file across repos stored once
- **Fork efficiency** — Forked repos share all existing objects
- **Simple backup** — `rsync` the objects directory

### Object Reference Counting

The hub maintains a global reference count per object:

```sql
-- In hub.db
CREATE TABLE object_refs (
    hash        TEXT PRIMARY KEY,
    ref_count   INTEGER NOT NULL DEFAULT 1,
    size        INTEGER NOT NULL,
    compressed  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TEXT NOT NULL
);
```

When a push adds objects:
1. For each object, check if it exists in the shared store
2. If exists: increment `ref_count`
3. If not: write to disk, insert with `ref_count = 1`

When a repo is deleted:
1. For each object referenced by the repo, decrement `ref_count`
2. Objects with `ref_count = 0` are garbage collected

### Object Retrieval for Pull

When a client pulls and needs objects:

1. Server reads operation `object_ref` fields
2. Checks which objects the client doesn't have (based on negotiate)
3. Reads objects from shared store
4. Sends in pull response

---

## Concurrency & Locking

### Read Operations (Browse, Pull)

- SQLite WAL mode allows concurrent reads
- No locking needed for reads
- Multiple users can browse/pull simultaneously

### Write Operations (Push)

- One push at a time per repository (mutex per repo)
- Push acquires repo lock → writes operations → updates stream → releases lock
- Lock is in-memory (Go sync.Mutex per repo ID)
- Lock timeout: 30 seconds

### Database Connection Pooling

```go
// Per-repo connection pool
db.SetMaxOpenConns(4)       // Max concurrent connections
db.SetMaxIdleConns(2)       // Keep 2 idle connections
db.SetConnMaxLifetime(1h)   // Recycle after 1 hour
```

### Repo Database Caching

Frequently accessed repos keep their database connections open in an LRU cache:

```go
type RepoDBCache struct {
    cache    *lru.Cache  // repo_id → *sql.DB
    capacity int         // e.g., 256 repos
}
```

When a repo is evicted from the cache, its database connection is closed. This prevents opening too many databases simultaneously.

---

## Size Limits & Quotas

| Resource | Default Limit |
|----------|--------------|
| Object size | 100 MB per object |
| Push batch | 500 MB per push |
| Repo count per user | 100 |
| Repo count per org | 500 |
| Total storage per user | 10 GB |

Configurable via server config:

```toml
[limits]
max_object_size = "100MB"
max_push_size = "500MB"
max_repos_per_user = 100
max_repos_per_org = 500
max_storage_per_user = "10GB"
```

---

## Backup & Recovery

### Full Backup

```bash
# Stop writes (or accept point-in-time inconsistency)
rsync -a /var/lib/loomhub/ /backup/loomhub/
```

### Incremental Backup

```bash
# Hub database
sqlite3 /var/lib/loomhub/hub.db ".backup /backup/hub.db"

# Per-repo databases (online backup)
for db in /var/lib/loomhub/repos/*/*/loom.db; do
    sqlite3 "$db" ".backup /backup/$(dirname $db | sed 's|.*/repos/||')/loom.db"
done

# Objects (rsync handles incremental naturally)
rsync -a /var/lib/loomhub/objects/ /backup/objects/
```

### Recovery

1. Stop server
2. Restore files to data directory
3. Start server
4. Run `loomhub doctor` to verify integrity

---

## Garbage Collection

### Object GC

Runs periodically (default: daily) or on demand:

1. Scan `object_refs` table for entries with `ref_count = 0`
2. Verify no repo references the object (safety check)
3. Delete object from disk
4. Remove entry from `object_refs`

### Database Compaction

Per-repo Loom databases can be compacted:

1. `VACUUM` to reclaim space from deleted rows
2. Optimize FTS5 indexes
3. Rebuild statistics for query planner

Triggered when:
- Repo database exceeds size threshold
- Admin requests compaction
- After large deletions

---

## Health Checks

### `GET /health`

Returns server health status:

```json
{
  "status": "ok",
  "version": "0.1.0",
  "uptime": "72h15m",
  "storage": {
    "hub_db_size": "2.1 MB",
    "object_store_size": "1.3 GB",
    "repo_count": 42,
    "total_repo_db_size": "156 MB"
  }
}
```

### `loomhub doctor`

CLI command to verify integrity:

1. Check hub database schema version
2. Verify all repo directories exist
3. Verify object store consistency (files match hashes)
4. Check for orphaned objects
5. Verify each repo database integrity (`PRAGMA integrity_check`)
