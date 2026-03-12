# LoomHub — Agent Development Guide

> This file helps AI coding agents work on the LoomHub codebase effectively.

## Project Structure

```
cmd/loomhub/                  Server entry point
internal/
  api/                        HTTP handlers (chi)
    handler.go                Handler struct, dependencies
    auth.go                   Register, login, JWT middleware
    looms.go                  CRUD for looms
    collaborators.go          Collaborator management
    content.go                Stream/entity/checkpoint read endpoints
    sync.go                   Negotiate/push/pull handlers
    orgs.go                   Organization endpoints
  auth/                       JWT utilities, context helpers
  models/                     Shared types (Loom, User, etc.)
  store/                      Database access layer
    users.go, looms.go, collaborators.go, activities.go, etc.
  sync/
    handler.go                Sync protocol logic (negotiate, send, receive)
    protocol.go               Wire types (request/response structs)
    content.go                Per-loom data queries (streams, entities, checkpoints)
    objects.go                Object store (shared across looms)
    loomdb.go                 Per-loom SQLite manager
    diff.go                   Stream comparison
  server/
    server.go                 Router setup, middleware, route registration
frontend/
  src/
    api/                      API client functions (ofetch)
    stores/                   Pinia stores (auth, loom)
    views/                    Vue components/pages
    router/                   Vue Router config
    utils/                    Shared utilities
```

## Backend Patterns

### Adding an API endpoint

1. Add handler method to `internal/api/handler.go`'s `Handler` struct
2. Implement in the appropriate file (e.g., `looms.go`, `content.go`)
3. Register route in `internal/server/server.go`
4. Use `writeJSON(w, status, data)` and `writeError(w, status, code, message)` for responses

```go
func (h *Handler) MyEndpoint(w http.ResponseWriter, r *http.Request) {
    // Auth check if needed
    cu := auth.GetUser(r.Context())
    if cu == nil {
        writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
        return
    }

    // Get URL params
    owner := chi.URLParam(r, "owner")

    // Do work...
    result, err := h.someStore.DoThing(owner)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "internal_error", "Failed")
        return
    }

    writeJSON(w, http.StatusOK, result)
}
```

### Auth patterns

- `auth.GetUser(r.Context())` returns `*auth.ContextUser` or nil
- `auth.RequireAuth` middleware rejects unauthenticated requests
- `h.collabs.CheckAccess(loom, userID, isAdmin)` returns permission: "admin", "write", "read", or ""

### Store pattern

Each store wraps a `*sql.DB` and provides typed methods:

```go
type MyStore struct { db *sql.DB }
func NewMyStore(db *sql.DB) *MyStore { return &MyStore{db: db} }
func (s *MyStore) Get(id string) (*MyModel, error) { ... }
```

## Frontend Patterns

### API calls

Use functions from `src/api/` that wrap `ofetch`:

```typescript
// src/api/client.ts provides the base `api()` function
import { api } from './client'

export function getThings(): Promise<Thing[]> {
  return api('/things')
}
```

### Stores

Pinia stores in `src/stores/`. Auth store manages JWT token and user state.

### Views

Vue 3 `<script setup lang="ts">` with Composition API. Tailwind CSS for styling (dark theme).

## Sync Protocol

The sync protocol is operation-based:

1. **Negotiate**: Client sends stream states → server returns common seqs, tells if send/receive needed
2. **Send (push)**: Client sends operations + objects → server stores them
3. **Receive (pull)**: Client requests ops after a seq → server returns ops + objects

Wire types are in `internal/sync/protocol.go`. Operations use `json.RawMessage` for delta and meta fields.

Important: The sync endpoints have a **different route prefix** (`/{owner}/{loom}/api/v1/`) than the main API (`/api/v1/`).

## Testing

```bash
# Backend
go test ./... -v

# Frontend
cd frontend && npx vitest run
```

Backend tests are integration tests in `test/integration/api_test.go` that spin up a real server.
Frontend tests use Vitest + happy-dom + @vue/test-utils.

## Common Mistakes

- Don't confuse the two route prefixes: `/api/v1/looms/{owner}/{loom}/...` (content API) vs `/{owner}/{loom}/api/v1/...` (sync API)
- Don't forget auth checks on new endpoints — check `collabs.CheckAccess()` for permission
- Don't use Git vocabulary in user-facing strings (use send/receive, stream, checkpoint, loom)
- Per-loom databases have a different schema than the hub database — they mirror the Loom CLI's schema
- Objects are stored in a shared directory, not per-loom
