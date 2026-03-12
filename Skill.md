# LoomHub — Skill Context for LLMs

> This file provides context for AI assistants working with or explaining LoomHub.

## What is LoomHub?

LoomHub is a GitHub-like hosting platform for Loom projects. It provides web-based repository hosting, collaboration, and a sync API that the Loom CLI talks to.

**LoomHub is to Loom what GitHub is to Git** — but built specifically for Loom's operation-based model.

## Features

- **Loom hosting**: Create and manage Loom repositories (called "looms")
- **Sync protocol**: Receive operations and objects from `loom send`, serve them via `loom receive`
- **Collaboration**: Invite collaborators with read/write/admin permissions
- **Organizations**: Group looms and members under org namespaces
- **Work Requests**: Issue-like tracking with labels (Loom's term for pull requests / issues)
- **Checkpoint log**: Browse, filter, and paginate operation history
- **Stream comparison**: Visual diff between any two streams
- **Activity feed**: Track sends, changes, and team activity

## Vocabulary

Use Loom terminology, not Git terminology:

| Git/GitHub term | Loom/LoomHub term |
|----------------|-------------------|
| repository | loom |
| push | send |
| pull | receive |
| pull request | work request |
| branch | stream |
| commit | checkpoint / operation |
| clone | (not yet implemented) |
| remote | hub |

## Tech Stack

- **Backend**: Go 1.25 + chi router + SQLite (modernc.org/sqlite, pure Go)
- **Frontend**: Vue 3 + TypeScript + Vite + Tailwind CSS 4 + Pinia + ofetch
- **Auth**: JWT tokens with bcrypt password hashing
- **Testing**: Go integration tests (backend), Vitest + happy-dom + @vue/test-utils (frontend)

## Architecture

```
Frontend (Vue 3 SPA)
  ↓ API calls (/api/v1/...)
Backend (Go + chi)
  ↓
Hub Database (SQLite — users, looms, orgs, collaborators, activities)
  ↓
Per-Loom Databases (SQLite — operations, streams, entities, checkpoints)
  ↓
Shared Object Store (content-addressed files on disk)
```

### Two database layers:
1. **Hub DB** (`data/loomhub.db`): Users, looms, organizations, collaborators, activities, labels, work requests
2. **Per-Loom DBs** (`data/looms/{owner}/{name}/loom.db`): Operations, streams, entities — same schema as the Loom CLI's `.loom/loom.db`

## API Endpoints

### Auth
- `POST /api/v1/auth/register` — Register new user
- `POST /api/v1/auth/login` — Login, returns JWT

### User
- `GET /api/v1/user` — Current user profile
- `PATCH /api/v1/user` — Update profile
- `GET /api/v1/users/{username}` — Public profile

### Looms
- `GET /api/v1/looms` — List user's looms
- `POST /api/v1/looms` — Create loom
- `GET /api/v1/looms/{owner}/{loom}` — Get loom details
- `PATCH /api/v1/looms/{owner}/{loom}` — Update loom
- `DELETE /api/v1/looms/{owner}/{loom}` — Delete loom

### Collaboration
- `GET /api/v1/looms/{owner}/{loom}/collaborators` — List collaborators
- `POST /api/v1/looms/{owner}/{loom}/collaborators` — Add collaborator
- `DELETE /api/v1/looms/{owner}/{loom}/collaborators` — Remove collaborator

### Content (per-loom data)
- `GET /api/v1/looms/{owner}/{loom}/streams` — List streams
- `GET /api/v1/looms/{owner}/{loom}/entities` — List entities
- `GET /api/v1/looms/{owner}/{loom}/checkpoints` — List checkpoints (filterable by type, author, path)
- `GET /api/v1/looms/{owner}/{loom}/compare` — Diff two streams
- `GET /api/v1/looms/{owner}/{loom}/objects/{ref}` — Get object content

### Sync (per-loom, different route prefix)
- `POST /{owner}/{loom}/api/v1/negotiate` — Compare stream states
- `POST /{owner}/{loom}/api/v1/push` — Receive operations from CLI
- `POST /{owner}/{loom}/api/v1/pull` — Send operations to CLI

### Organizations
- `POST /api/v1/orgs` — Create org
- `GET /api/v1/orgs/{org}` — Get org
- `GET /api/v1/orgs/{org}/members` — List members
- `POST /api/v1/orgs/{org}/members` — Add member

## Permission Model

```
site admin → can do anything
owner      → full control of their looms
admin      → collaborator with admin permission
write      → can send (push) to the loom
read       → can receive (pull) and browse
public     → anyone can read public looms
```

`CheckAccess(loom, userID, isAdmin)` returns the effective permission level.

## Running

```bash
# Backend
go run ./cmd/loomhub          # starts on :8080

# Frontend
cd frontend && npm run dev    # starts on :5173, proxies to :8080
```

## Related Projects

- **Loom CLI** (github.com/LoomHubDev/loom): The versioning system itself
- **Loom Docs** (github.com/LoomHubDev/loom-docs): Documentation site
