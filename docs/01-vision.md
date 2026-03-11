# LoomHub — Vision & Design Principles

## What Is LoomHub?

LoomHub is the remote hosting and collaboration platform for [Loom](https://github.com/flakerimi/loom) repositories. It is to Loom what GitHub is to Git — a web-based service where teams push, pull, browse, and collaborate on Loom-versioned projects.

## Why LoomHub?

Loom tracks every operation across multiple content spaces (code, docs, design, data, config, notes) in a single, unified timeline. LoomHub makes that timeline collaborative:

- **Push/Pull** — Teams sync Loom projects through a central server
- **Browse** — View checkpoint history, diffs, and entity trees through a web UI
- **Collaborate** — Merge requests, reviews, and discussions on streams
- **Discover** — Search across projects, checkpoints, and entities
- **Automate** — Webhooks and CI/CD integration on checkpoint events

## Design Principles

### 1. Loom-Native

LoomHub speaks Loom's protocol natively. It stores operations, checkpoints, streams, and content-addressed objects exactly as Loom does locally. No translation layer, no impedance mismatch.

### 2. Multi-Space Aware

Unlike GitHub (code-only), LoomHub understands all Loom spaces: code, docs, design, data, config, notes. The web UI renders each space appropriately — syntax-highlighted code, rendered markdown, visual design components, data schema diagrams.

### 3. Operation-Level Granularity

LoomHub stores and displays the full operation log, not just snapshots. You can see every individual change, not just checkpoint-to-checkpoint diffs. This enables fine-grained audit trails and precise rollbacks.

### 4. Append-Only by Default

Like Loom itself, LoomHub's storage is append-only. Operations are never deleted (only compacted with explicit admin action). This guarantees a complete audit trail.

### 5. Simple & Fast

- Go backend, single binary deployment
- SQLite per project (no shared database cluster needed for small/medium scale)
- Minimal JavaScript — server-rendered HTML with progressive enhancement
- Fast push/pull over HTTP with delta sync

### 6. Self-Hostable First

LoomHub is designed to run on a single machine or a small cluster. No Kubernetes required. A single `loomhub serve` command starts the server.

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                    LoomHub Server                     │
│                                                       │
│  ┌──────────┐  ┌──────────┐  ┌───────────────────┐  │
│  │ Sync API │  │ REST API │  │    Web UI (SSR)    │  │
│  │ (push/   │  │ (users,  │  │  (templ + htmx)   │  │
│  │  pull,   │  │  repos,  │  │                    │  │
│  │  negot.) │  │  search) │  │                    │  │
│  └────┬─────┘  └────┬─────┘  └────────┬──────────┘  │
│       │              │                 │              │
│  ┌────┴──────────────┴─────────────────┴──────────┐  │
│  │              Application Layer                  │  │
│  │  (auth, permissions, webhooks, notifications)   │  │
│  └────────────────────┬───────────────────────────┘  │
│                       │                              │
│  ┌────────────────────┴───────────────────────────┐  │
│  │              Storage Layer                      │  │
│  │                                                 │  │
│  │  ┌─────────┐  ┌───────────┐  ┌──────────────┐ │  │
│  │  │ Hub DB  │  │ Per-Repo  │  │ Object Store │ │  │
│  │  │(SQLite) │  │  Loom DBs │  │  (shared)    │ │  │
│  │  │         │  │  (SQLite) │  │              │ │  │
│  │  └─────────┘  └───────────┘  └──────────────┘ │  │
│  └────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
```

### Key Layers

| Layer | Responsibility |
|-------|---------------|
| **Sync API** | Implements Loom's negotiate/push/pull protocol for `loom push` and `loom pull` |
| **REST API** | CRUD for users, repos, merge requests, search, webhooks |
| **Web UI** | Server-rendered pages with htmx for interactivity |
| **Application** | Business logic — auth, permissions, notifications, webhook dispatch |
| **Storage** | Hub-level SQLite (users, repos, permissions) + per-repo Loom databases + shared object store |

## Tech Stack

| Component | Technology | Why |
|-----------|-----------|-----|
| Language | Go | Matches Loom; single binary; excellent concurrency |
| HTTP Router | chi | Lightweight, composable, idiomatic Go |
| Templates | templ | Type-safe Go templates, compiles to Go code |
| Interactivity | htmx | Progressive enhancement without SPA complexity |
| CSS | Tailwind CSS | Utility-first, fast prototyping |
| Database | SQLite (modernc.org) | Same as Loom; no external DB dependency; WAL mode |
| Auth | bcrypt + JWT | Simple, proven |
| Object Storage | Filesystem (content-addressed) | Shared across repos, deduplication via SHA-256 |

## Deployment Model

### Single Server (Default)

```
loomhub serve --data /var/lib/loomhub --port 3000
```

Everything runs in a single process:
- HTTP server (sync + API + web)
- SQLite databases (hub + per-repo)
- Object store on local filesystem

### Scale-Out (Future)

For larger deployments:
- Object store → S3-compatible backend
- Hub database → PostgreSQL
- Per-repo databases → Sharded across nodes
- HTTP → Behind reverse proxy with sticky sessions

## Comparison with GitHub

| Feature | GitHub | LoomHub |
|---------|--------|---------|
| Versioning system | Git | Loom |
| Content types | Code only | Code + docs + design + data + config + notes |
| Change granularity | Commits (manual) | Operations (automatic) + Checkpoints (named) |
| Branching | Branches | Streams |
| Merge requests | Pull Requests | Merge Requests (stream → stream) |
| Storage | Packfiles | SQLite + content-addressed objects |
| Self-hosting | GitHub Enterprise ($$) | Built-in, free |
| CI/CD | GitHub Actions | Webhooks + external runners |
