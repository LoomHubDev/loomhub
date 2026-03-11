# LoomHub — Data Models

LoomHub has two layers of data:

1. **Hub-level** — Users, organizations, repositories, permissions, merge requests, webhooks (stored in a central Hub database)
2. **Repo-level** — Operations, checkpoints, streams, entities, objects (stored in per-repo Loom databases, based on the local `.loom/` schema with a server-side adaptation layer)

This document covers the Hub-level models. Repo-level models are defined by [Loom's data model](../../loom/docs/04-data-models.md).

> **Note on storage compatibility:** Per-repo databases use the same *schema* as a local `.loom/` project (operations, checkpoints, streams, entities, metadata tables). However, object storage differs: locally, Loom stores objects on disk inside `.loom/objects/` with per-repo reference counts. On LoomHub, objects are stored in a **shared cross-repo object store** with a **global reference count table** in the hub database. The sync protocol handles this transparently — the `objects` table in each repo database stores the hash index, but actual blob storage is redirected to the shared store. See [Repository Hosting](06-repository-hosting.md) for details on this adaptation layer.

---

## Hub Database Schema

### Owners (Shared Namespace)

Users and organizations share a single namespace for URL routing (`/{owner}/{repo}`). The `owners` table enforces this globally — no user and org can have the same name.

```sql
CREATE TABLE owners (
    id          TEXT PRIMARY KEY,       -- ULID (same as user.id or org.id)
    name        TEXT NOT NULL UNIQUE,    -- lowercase, alphanumeric + hyphens
    type        TEXT NOT NULL,           -- 'user' or 'org'
    created_at  TEXT NOT NULL
);

CREATE INDEX idx_owners_name ON owners(name);
```

When a user registers or an org is created, a row is inserted into `owners` first. If the name is taken (by either a user or an org), the insert fails.

### Users

```sql
CREATE TABLE users (
    id          TEXT PRIMARY KEY REFERENCES owners(id) ON DELETE CASCADE,
    username    TEXT NOT NULL UNIQUE,    -- lowercase, must match owners.name
    email       TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL DEFAULT '',
    bio         TEXT NOT NULL DEFAULT '',
    avatar_url  TEXT NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL,         -- bcrypt
    is_admin    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TEXT NOT NULL,           -- RFC 3339
    updated_at  TEXT NOT NULL
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

### SSH Keys

```sql
CREATE TABLE ssh_keys (
    id          TEXT PRIMARY KEY,       -- ULID
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,           -- "MacBook Pro", "CI Server"
    fingerprint TEXT NOT NULL UNIQUE,    -- SHA-256 fingerprint
    public_key  TEXT NOT NULL,
    last_used_at TEXT,
    created_at  TEXT NOT NULL
);

CREATE INDEX idx_ssh_keys_user ON ssh_keys(user_id);
CREATE INDEX idx_ssh_keys_fingerprint ON ssh_keys(fingerprint);
```

### Access Tokens

```sql
CREATE TABLE access_tokens (
    id          TEXT PRIMARY KEY,       -- ULID
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,           -- "CI token", "Deploy key"
    token_hash  TEXT NOT NULL UNIQUE,    -- SHA-256 of token
    scopes      TEXT NOT NULL DEFAULT 'read', -- comma-separated: read,write,admin
    expires_at  TEXT,                    -- NULL = never expires
    last_used_at TEXT,
    created_at  TEXT NOT NULL
);

CREATE INDEX idx_tokens_user ON access_tokens(user_id);
CREATE INDEX idx_tokens_hash ON access_tokens(token_hash);
```

### Organizations

```sql
CREATE TABLE organizations (
    id          TEXT PRIMARY KEY REFERENCES owners(id) ON DELETE CASCADE,
    name        TEXT NOT NULL UNIQUE,    -- lowercase slug, must match owners.name
    display_name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    avatar_url  TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);

CREATE INDEX idx_orgs_name ON organizations(name);
```

### Organization Members

```sql
CREATE TABLE org_members (
    org_id      TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role        TEXT NOT NULL DEFAULT 'member', -- owner, admin, member
    created_at  TEXT NOT NULL,
    PRIMARY KEY (org_id, user_id)
);

CREATE INDEX idx_org_members_user ON org_members(user_id);
```

### Repositories

```sql
CREATE TABLE repositories (
    id          TEXT PRIMARY KEY,       -- ULID
    owner_id    TEXT NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,           -- lowercase slug
    description TEXT NOT NULL DEFAULT '',
    visibility  TEXT NOT NULL DEFAULT 'public', -- public, private
    default_stream TEXT NOT NULL DEFAULT 'main',
    disk_path   TEXT NOT NULL,           -- relative path to repo data on disk
    size_bytes  INTEGER NOT NULL DEFAULT 0,

    -- Stats (denormalized, updated periodically)
    checkpoint_count INTEGER NOT NULL DEFAULT 0,
    stream_count     INTEGER NOT NULL DEFAULT 0,
    star_count       INTEGER NOT NULL DEFAULT 0,
    fork_count       INTEGER NOT NULL DEFAULT 0,

    forked_from TEXT REFERENCES repositories(id) ON DELETE SET NULL,

    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL,
    pushed_at   TEXT,                    -- last push timestamp

    UNIQUE(owner_id, name)
);

CREATE INDEX idx_repos_owner ON repositories(owner_id);
CREATE INDEX idx_repos_name ON repositories(name);
CREATE INDEX idx_repos_visibility ON repositories(visibility);
CREATE INDEX idx_repos_pushed ON repositories(pushed_at);
CREATE INDEX idx_repos_stars ON repositories(star_count);
```

Owner type (user or org) can be resolved via `JOIN owners ON owners.id = repositories.owner_id`. Cascading deletes now work through the database: deleting an owner cascades to the user/org row and to all their repositories.

### Repository Collaborators

```sql
CREATE TABLE repo_collaborators (
    repo_id     TEXT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission  TEXT NOT NULL DEFAULT 'read', -- read, write, admin
    created_at  TEXT NOT NULL,
    PRIMARY KEY (repo_id, user_id)
);

CREATE INDEX idx_collabs_user ON repo_collaborators(user_id);
```

### Stars

```sql
CREATE TABLE stars (
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    repo_id     TEXT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    created_at  TEXT NOT NULL,
    PRIMARY KEY (user_id, repo_id)
);

CREATE INDEX idx_stars_repo ON stars(repo_id);
```

### Merge Requests

```sql
CREATE TABLE merge_requests (
    id          TEXT PRIMARY KEY,       -- ULID
    repo_id     TEXT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    number      INTEGER NOT NULL,       -- sequential per-repo number
    author_id   TEXT NOT NULL REFERENCES users(id),

    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',

    source_stream TEXT NOT NULL,         -- stream name being merged
    target_stream TEXT NOT NULL,         -- stream name to merge into

    -- Sequence tracking for diff computation
    source_head_seq INTEGER NOT NULL,    -- head seq of source at creation/update
    target_head_seq INTEGER NOT NULL,    -- head seq of target at creation/update
    merge_base_seq  INTEGER NOT NULL,    -- common ancestor seq

    status      TEXT NOT NULL DEFAULT 'open', -- open, merged, closed
    merged_by   TEXT REFERENCES users(id),
    merged_at   TEXT,
    closed_at   TEXT,

    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL,

    UNIQUE(repo_id, number)
);

CREATE INDEX idx_mr_repo ON merge_requests(repo_id, status);
CREATE INDEX idx_mr_author ON merge_requests(author_id);
CREATE INDEX idx_mr_number ON merge_requests(repo_id, number);
```

### Merge Request Comments

```sql
CREATE TABLE mr_comments (
    id          TEXT PRIMARY KEY,       -- ULID
    mr_id       TEXT NOT NULL REFERENCES merge_requests(id) ON DELETE CASCADE,
    author_id   TEXT NOT NULL REFERENCES users(id),
    body        TEXT NOT NULL,

    -- Inline comment fields (NULL for general comments)
    space_id    TEXT,                    -- which space
    entity_path TEXT,                    -- which file/entity
    line_number INTEGER,                -- line in diff

    -- System events
    is_system   BOOLEAN NOT NULL DEFAULT FALSE,
    event_type  TEXT,                    -- status_change, review, merge, etc.

    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);

CREATE INDEX idx_mr_comments_mr ON mr_comments(mr_id);
CREATE INDEX idx_mr_comments_author ON mr_comments(author_id);
```

### Reviews

```sql
CREATE TABLE reviews (
    id          TEXT PRIMARY KEY,       -- ULID
    mr_id       TEXT NOT NULL REFERENCES merge_requests(id) ON DELETE CASCADE,
    reviewer_id TEXT NOT NULL REFERENCES users(id),
    status      TEXT NOT NULL,           -- approved, changes_requested, commented
    body        TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL
);

CREATE INDEX idx_reviews_mr ON reviews(mr_id);
```

### Labels

```sql
CREATE TABLE labels (
    id          TEXT PRIMARY KEY,       -- ULID
    repo_id     TEXT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    color       TEXT NOT NULL DEFAULT '#888888', -- hex color
    description TEXT NOT NULL DEFAULT '',
    UNIQUE(repo_id, name)
);

CREATE TABLE mr_labels (
    mr_id       TEXT NOT NULL REFERENCES merge_requests(id) ON DELETE CASCADE,
    label_id    TEXT NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (mr_id, label_id)
);
```

### Webhooks

```sql
CREATE TABLE webhooks (
    id          TEXT PRIMARY KEY,       -- ULID
    repo_id     TEXT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    url         TEXT NOT NULL,
    secret      TEXT NOT NULL DEFAULT '', -- HMAC signing secret
    events      TEXT NOT NULL DEFAULT 'push', -- comma-separated: push,mr,checkpoint,review
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);

CREATE INDEX idx_webhooks_repo ON webhooks(repo_id);
```

### Webhook Deliveries

```sql
CREATE TABLE webhook_deliveries (
    id          TEXT PRIMARY KEY,       -- ULID
    webhook_id  TEXT NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event       TEXT NOT NULL,
    payload     TEXT NOT NULL,           -- JSON
    status_code INTEGER,
    response    TEXT,
    duration_ms INTEGER,
    delivered_at TEXT NOT NULL
);

CREATE INDEX idx_deliveries_webhook ON webhook_deliveries(webhook_id);
```

### Activity Feed

```sql
CREATE TABLE activities (
    id          TEXT PRIMARY KEY,       -- ULID
    actor_id    TEXT NOT NULL REFERENCES users(id),
    repo_id     TEXT REFERENCES repositories(id) ON DELETE CASCADE,
    type        TEXT NOT NULL,           -- push, create_repo, create_mr, merge, star, fork, comment
    payload     TEXT NOT NULL DEFAULT '{}', -- JSON with event-specific data
    created_at  TEXT NOT NULL
);

CREATE INDEX idx_activities_actor ON activities(actor_id, created_at);
CREATE INDEX idx_activities_repo ON activities(repo_id, created_at);
CREATE INDEX idx_activities_type ON activities(type, created_at);
```

---

## Disk Layout

```
/var/lib/loomhub/
├── hub.db                          # Hub-level SQLite database
├── objects/                        # Shared content-addressed object store
│   ├── ab/
│   │   └── abcdef0123456789...
│   └── cd/
│       └── cdef0123456789...
├── repos/                          # Per-repo Loom databases
│   ├── flakerimi/
│   │   ├── my-app/
│   │   │   └── loom.db            # Repo's Loom database
│   │   └── another-project/
│   │       └── loom.db
│   └── acme-org/
│       └── platform/
│           └── loom.db
└── tmp/                            # Temporary files (uploads, etc.)
```

### Object Store Sharing

All repositories share a single content-addressed object store. Since objects are identified by SHA-256 hash, identical content across repos is stored only once. This provides:

- **Deduplication** — Forked repos share most objects
- **Efficiency** — Common libraries/dependencies stored once
- **Simplicity** — Single directory to back up

### Per-Repo Database

Each repository has its own SQLite database containing tables based on the Loom schema:
- `operations` — append-only operation log
- `checkpoints` — named points in time
- `streams` — timeline management
- `entities` — current entity states
- `objects` — object hash index (hash, size, compressed flag — but **not** ref_count or blob storage)
- `metadata` — repo-level key-value config

The schema is based on a local `.loom/` database with one key difference: **object blob storage is delegated to the shared object store** rather than stored per-repo. The `objects` table in each repo database serves as an index of which objects the repo references, but the actual bytes live in the shared store. The sync handler translates between Loom's native protocol (which assumes per-repo object storage) and LoomHub's shared store transparently.

---

## Entity Relationships

```
Owner (shared namespace) ─── type=user ──── User
        │                └── type=org  ──── Organization
        │
        └── owns ──── Repository
                          │
User ─────┬── member of ── Organization
          │
          ├── stars ────── Repository
          │
          ├── authors ──── MergeRequest ──── has ── Comments
          │                     │                     Reviews
          │                     │                     Labels
          └── creates ──── AccessToken
                            SSHKey
                            Activity

Repository ──── has ──── Webhooks ──── has ──── Deliveries
```

## Key Invariants

1. **Owner namespace uniqueness** — Enforced by the `owners` table: usernames and org names share a single `UNIQUE(name)` constraint. No user and org can have the same name. This makes `/{owner}/{repo}` routing unambiguous.
2. **Repo uniqueness** — Repository names are unique within an owner (`UNIQUE(owner_id, name)`)
3. **MR numbering** — Merge request numbers are sequential per repository
4. **Cascading deletes** — Deleting an owner cascades through `owners → users/organizations → repositories → (MRs, collaborators, stars, webhooks, etc.)`. All enforced by foreign keys with `ON DELETE CASCADE`.
5. **Object immutability** — Objects in the shared store are never modified, only added or garbage-collected
