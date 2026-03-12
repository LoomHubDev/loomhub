# LoomHub — Vision & Design Principles

## What Is LoomHub?

LoomHub is the remote hosting and collaboration platform for [Loom](https://github.com/flakerimi/loom) looms. It is to Loom what GitHub is to Git — a web-based service where teams send, receive, browse, and collaborate on Loom-versioned projects.

## Why LoomHub?

Loom tracks every operation across multiple content spaces (code, docs, design, data, config, notes) in a single, unified timeline. LoomHub makes that timeline collaborative:

- **Send/Receive** — Teams sync Loom projects through a central server
- **Browse** — View checkpoint history, diffs, and entity trees through a web UI
- **Collaborate** — Weave requests, reviews, and discussions on streams
- **Discover** — Search across projects, checkpoints, and entities
- **Automate** — Webhooks and CI/CD integration on checkpoint events

## Design Principles

### 1. Loom-Native

LoomHub implements Loom's sync protocol (negotiate/send/receive) exactly as defined in Loom's spec. The `loom` CLI interoperates with LoomHub without modifications. Server-side storage uses the same Loom database schema per loom, with object blob storage adapted to a shared cross-loom store (see [Loom Hosting](06-loom-hosting.md)).

### 2. Multi-Space Aware

Unlike GitHub (code-only), LoomHub understands all Loom spaces: code, docs, design, data, config, notes. The web UI renders each space appropriately — syntax-highlighted code, rendered markdown, visual design components, data schema diagrams.

### 3. Operation-Level Granularity

LoomHub stores and displays the full operation log, not just snapshots. You can see every individual change, not just checkpoint-to-checkpoint diffs. This enables fine-grained audit trails and precise rollbacks.

### 4. Append-Only by Default

Like Loom itself, LoomHub's storage is append-only. Operations are never deleted (only compacted with explicit admin action). This guarantees a complete audit trail.

### 5. Simple & Fast

- Go backend, single binary deployment (Vue SPA embedded via `embed.FS`)
- SQLite per project (no shared database cluster needed for small/medium scale)
- Vue 3 SPA for a reactive UI with real-time updates
- Fast send/receive over HTTP with delta sync

### 6. Self-Hostable First

LoomHub is designed to run on a single machine or a small cluster. No Kubernetes required. A single `loomhub serve` command starts the server.

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                    LoomHub Server                     │
│                                                       │
│  ┌──────────┐  ┌──────────┐  ┌───────────────────┐  │
│  │ Sync API │  │ REST API │  │   Static Files    │  │
│  │ (send/   │  │ (users,  │  │  (Vue SPA dist)   │  │
│  │  recv,   │  │  looms,  │  │  via embed.FS     │  │
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
│  │  │ Hub DB  │  │ Per-Loom  │  │ Object Store │ │  │
│  │  │(SQLite) │  │  Loom DBs │  │  (shared)    │ │  │
│  │  │         │  │  (SQLite) │  │              │ │  │
│  │  └─────────┘  └───────────┘  └──────────────┘ │  │
│  └────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│                    Vue 3 SPA (Browser)                │
│                                                       │
│  Vue Router ── Pinia Store ── API Client              │
│       │              │              │                 │
│  Pages/Views    State Mgmt    fetch(/api/v1/...)      │
│  Components     Reactivity                            │
│  Tailwind CSS                                         │
└─────────────────────────────────────────────────────┘
```

### Key Layers

| Layer | Responsibility |
|-------|---------------|
| **Sync API** | Implements Loom's negotiate/send/receive protocol at `/{owner}/{loom}/api/v1/...` |
| **REST API** | CRUD for users, looms, weave requests, search, webhooks — consumed by the Vue SPA |
| **Vue SPA** | Client-side rendered UI with reactive components, compiled and embedded in the Go binary |
| **Application** | Business logic — auth, permissions, notifications, webhook dispatch |
| **Storage** | Hub-level SQLite (users, looms, permissions) + per-loom Loom databases + shared object store |

## Tech Stack

### Backend (Go)

| Component | Technology | Why |
|-----------|-----------|-----|
| Language | Go | Matches Loom; single binary; excellent concurrency |
| HTTP Router | chi | Lightweight, composable, idiomatic Go |
| Database | SQLite (modernc.org) | Same as Loom; no external DB dependency; WAL mode |
| Auth | bcrypt + JWT | Simple, proven |
| Object Storage | Filesystem (content-addressed) | Shared across looms, deduplication via SHA-256 |
| SPA Embedding | `embed.FS` | Vue dist bundled into Go binary for single-binary deploy |

### Frontend (Vue)

| Component | Technology | Why |
|-----------|-----------|-----|
| Framework | Vue 3 (Composition API) | Reactive UI; you already know Vue/Nuxt; large ecosystem |
| Build Tool | Vite | Fast HMR, optimized builds |
| Routing | Vue Router | Client-side navigation |
| State | Pinia | Lightweight, type-safe state management |
| CSS | Tailwind CSS | Utility-first, fast prototyping |
| HTTP Client | ofetch | Lightweight, interceptors, auto-retry |
| Syntax Highlighting | Shiki | Accurate, theme-able, same engine as VS Code |
| Markdown | markdown-it | Fast, extensible markdown rendering |
| Diff Viewer | Custom component | Multi-space aware, inline comments |

## Deployment Model

### Single Server (Default)

```
loomhub serve --data /var/lib/loomhub --port 3000
```

Everything runs in a single process:
- HTTP server (sync + API + web)
- SQLite databases (hub + per-loom)
- Object store on local filesystem

### Scale-Out (Future)

For larger deployments:
- Object store → S3-compatible backend
- Hub database → PostgreSQL
- Per-loom databases → Sharded across nodes
- HTTP → Behind reverse proxy with sticky sessions

## Comparison with GitHub

| Feature | GitHub | LoomHub |
|---------|--------|---------|
| Versioning system | Git | Loom |
| Content types | Code only | Code + docs + design + data + config + notes |
| Change granularity | Commits (manual) | Operations (automatic) + Checkpoints (named) |
| Branching | Branches | Streams |
| Collaboration | Pull Requests | Weave Requests (stream → stream) |
| Storage | Packfiles | SQLite + content-addressed objects |
| Self-hosting | GitHub Enterprise ($$) | Built-in, free |
| CI/CD | GitHub Actions | Webhooks + external runners |
