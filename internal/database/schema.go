package database

import "database/sql"

func Migrate(db *sql.DB) error {
	_, err := db.Exec(hubSchema)
	return err
}

const hubSchema = `
CREATE TABLE IF NOT EXISTS owners (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    type        TEXT NOT NULL CHECK(type IN ('user', 'org')),
    created_at  TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_owners_name ON owners(name);

CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY REFERENCES owners(id) ON DELETE CASCADE,
    username      TEXT NOT NULL UNIQUE,
    email         TEXT NOT NULL UNIQUE,
    display_name  TEXT NOT NULL DEFAULT '',
    bio           TEXT NOT NULL DEFAULT '',
    avatar_url    TEXT NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL,
    is_admin      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TEXT NOT NULL,
    updated_at    TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS access_tokens (
    id           TEXT PRIMARY KEY,
    user_id      TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    token_hash   TEXT NOT NULL UNIQUE,
    scopes       TEXT NOT NULL DEFAULT 'read',
    expires_at   TEXT,
    last_used_at TEXT,
    created_at   TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_tokens_user ON access_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_tokens_hash ON access_tokens(token_hash);

CREATE TABLE IF NOT EXISTS organizations (
    id           TEXT PRIMARY KEY REFERENCES owners(id) ON DELETE CASCADE,
    name         TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL DEFAULT '',
    description  TEXT NOT NULL DEFAULT '',
    avatar_url   TEXT NOT NULL DEFAULT '',
    created_at   TEXT NOT NULL,
    updated_at   TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_orgs_name ON organizations(name);

CREATE TABLE IF NOT EXISTS org_members (
    org_id     TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role       TEXT NOT NULL DEFAULT 'member' CHECK(role IN ('owner', 'admin', 'member')),
    created_at TEXT NOT NULL,
    PRIMARY KEY (org_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_org_members_user ON org_members(user_id);

CREATE TABLE IF NOT EXISTS looms (
    id               TEXT PRIMARY KEY,
    owner_id         TEXT NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    name             TEXT NOT NULL,
    description      TEXT NOT NULL DEFAULT '',
    visibility       TEXT NOT NULL DEFAULT 'public' CHECK(visibility IN ('public', 'private')),
    default_stream   TEXT NOT NULL DEFAULT 'main',
    disk_path        TEXT NOT NULL,
    size_bytes       INTEGER NOT NULL DEFAULT 0,
    checkpoint_count INTEGER NOT NULL DEFAULT 0,
    stream_count     INTEGER NOT NULL DEFAULT 0,
    pin_count        INTEGER NOT NULL DEFAULT 0,
    spin_count       INTEGER NOT NULL DEFAULT 0,
    spun_from        TEXT REFERENCES looms(id) ON DELETE SET NULL,
    created_at       TEXT NOT NULL,
    updated_at       TEXT NOT NULL,
    synced_at        TEXT,
    UNIQUE(owner_id, name)
);
CREATE INDEX IF NOT EXISTS idx_looms_owner ON looms(owner_id);
CREATE INDEX IF NOT EXISTS idx_looms_name ON looms(name);
CREATE INDEX IF NOT EXISTS idx_looms_visibility ON looms(visibility);
CREATE INDEX IF NOT EXISTS idx_looms_synced ON looms(synced_at);
CREATE INDEX IF NOT EXISTS idx_looms_pins ON looms(pin_count);

CREATE TABLE IF NOT EXISTS loom_collaborators (
    loom_id    TEXT NOT NULL REFERENCES looms(id) ON DELETE CASCADE,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission TEXT NOT NULL DEFAULT 'read' CHECK(permission IN ('read', 'write', 'admin')),
    created_at TEXT NOT NULL,
    PRIMARY KEY (loom_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_collabs_user ON loom_collaborators(user_id);

CREATE TABLE IF NOT EXISTS pins (
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    loom_id    TEXT NOT NULL REFERENCES looms(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL,
    PRIMARY KEY (user_id, loom_id)
);
CREATE INDEX IF NOT EXISTS idx_pins_loom ON pins(loom_id);

CREATE TABLE IF NOT EXISTS weave_requests (
    id              TEXT PRIMARY KEY,
    loom_id         TEXT NOT NULL REFERENCES looms(id) ON DELETE CASCADE,
    number          INTEGER NOT NULL,
    author_id       TEXT NOT NULL REFERENCES users(id),
    title           TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    source_stream   TEXT NOT NULL,
    target_stream   TEXT NOT NULL,
    source_head_seq INTEGER NOT NULL,
    target_head_seq INTEGER NOT NULL,
    merge_base_seq  INTEGER NOT NULL,
    status          TEXT NOT NULL DEFAULT 'open' CHECK(status IN ('open', 'woven', 'closed')),
    woven_by        TEXT REFERENCES users(id),
    woven_at        TEXT,
    closed_at       TEXT,
    created_at      TEXT NOT NULL,
    updated_at      TEXT NOT NULL,
    UNIQUE(loom_id, number)
);
CREATE INDEX IF NOT EXISTS idx_wr_loom ON weave_requests(loom_id, status);
CREATE INDEX IF NOT EXISTS idx_wr_author ON weave_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_wr_number ON weave_requests(loom_id, number);

CREATE TABLE IF NOT EXISTS wr_comments (
    id          TEXT PRIMARY KEY,
    wr_id       TEXT NOT NULL REFERENCES weave_requests(id) ON DELETE CASCADE,
    author_id   TEXT NOT NULL REFERENCES users(id),
    body        TEXT NOT NULL,
    space_id    TEXT,
    entity_path TEXT,
    line_number INTEGER,
    is_system   BOOLEAN NOT NULL DEFAULT FALSE,
    event_type  TEXT,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_wr_comments_wr ON wr_comments(wr_id);
CREATE INDEX IF NOT EXISTS idx_wr_comments_author ON wr_comments(author_id);

CREATE TABLE IF NOT EXISTS reviews (
    id          TEXT PRIMARY KEY,
    wr_id       TEXT NOT NULL REFERENCES weave_requests(id) ON DELETE CASCADE,
    reviewer_id TEXT NOT NULL REFERENCES users(id),
    status      TEXT NOT NULL CHECK(status IN ('approved', 'changes_requested', 'commented')),
    body        TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_reviews_wr ON reviews(wr_id);

CREATE TABLE IF NOT EXISTS labels (
    id          TEXT PRIMARY KEY,
    loom_id     TEXT NOT NULL REFERENCES looms(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    color       TEXT NOT NULL DEFAULT '#888888',
    description TEXT NOT NULL DEFAULT '',
    UNIQUE(loom_id, name)
);

CREATE TABLE IF NOT EXISTS wr_labels (
    wr_id    TEXT NOT NULL REFERENCES weave_requests(id) ON DELETE CASCADE,
    label_id TEXT NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (wr_id, label_id)
);

CREATE TABLE IF NOT EXISTS webhooks (
    id         TEXT PRIMARY KEY,
    loom_id    TEXT NOT NULL REFERENCES looms(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    secret     TEXT NOT NULL DEFAULT '',
    events     TEXT NOT NULL DEFAULT 'send',
    active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_webhooks_loom ON webhooks(loom_id);

CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id           TEXT PRIMARY KEY,
    webhook_id   TEXT NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event        TEXT NOT NULL,
    payload      TEXT NOT NULL,
    status_code  INTEGER,
    response     TEXT,
    duration_ms  INTEGER,
    delivered_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_deliveries_webhook ON webhook_deliveries(webhook_id);

CREATE TABLE IF NOT EXISTS activities (
    id         TEXT PRIMARY KEY,
    actor_id   TEXT NOT NULL REFERENCES users(id),
    loom_id    TEXT REFERENCES looms(id) ON DELETE CASCADE,
    type       TEXT NOT NULL,
    payload    TEXT NOT NULL DEFAULT '{}',
    created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_activities_actor ON activities(actor_id, created_at);
CREATE INDEX IF NOT EXISTS idx_activities_loom ON activities(loom_id, created_at);
CREATE INDEX IF NOT EXISTS idx_activities_type ON activities(type, created_at);

CREATE TABLE IF NOT EXISTS object_refs (
    hash      TEXT NOT NULL,
    loom_id   TEXT NOT NULL REFERENCES looms(id) ON DELETE CASCADE,
    size      INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (hash, loom_id)
);
CREATE INDEX IF NOT EXISTS idx_object_refs_loom ON object_refs(loom_id);
`
