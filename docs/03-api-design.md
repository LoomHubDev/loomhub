# LoomHub — API Design

LoomHub exposes two API surfaces:

1. **Sync API** — Implements Loom's native negotiate/send/receive protocol, routed by `{owner}/{loom}` (used by `loom send`/`loom receive`)
2. **REST API** — CRUD for hub-level resources (users, looms, weave requests, etc.)

The Vue SPA is served as static files and consumes the REST API.

All APIs share the same HTTP server and authentication layer.

---

## Base URL & Versioning

```
https://loomhub.dev/api/v1/...                     # REST API (hub-level)
https://loomhub.dev/{owner}/{loom}/api/v1/...      # Sync API (per-loom)
https://loomhub.dev/                               # Vue SPA (static files)
```

---

## Authentication

### Methods

| Method | Use Case | Header |
|--------|----------|--------|
| JWT Token | Web UI sessions | `Cookie: session=<jwt>` |
| Access Token | CLI / API clients | `Authorization: Bearer <token>` |
| Basic Auth | Simple scripts | `Authorization: Basic <base64>` |

### Token Scopes

| Scope | Access |
|-------|--------|
| `read` | Read public + private looms the user has access to |
| `write` | Send to looms, create WRs, post comments |
| `admin` | Manage looms, users, webhooks, permissions |

---

## Sync API

These endpoints implement Loom's native sync protocol as defined in [Loom's sync spec](../../loom/docs/06-systems/sync.md).

> **Design decision:** LoomHub does NOT extend or modify Loom's sync protocol. The request/response types are identical to what `loom-server` uses.

### How Loom CLI Reaches LoomHub

Loom's sync endpoints are at fixed paths (`/api/v1/negotiate`, `/api/v1/push`, `/api/v1/pull`). The Loom client constructs these from the hub URL base. When a user adds a LoomHub hub:

```bash
loom hub add origin https://loomhub.dev/flakerimi/my-app
```

The Loom client uses `https://loomhub.dev/flakerimi/my-app` as the base URL and appends the standard sync paths:

```
POST https://loomhub.dev/flakerimi/my-app/api/v1/negotiate
POST https://loomhub.dev/flakerimi/my-app/api/v1/push
POST https://loomhub.dev/flakerimi/my-app/api/v1/pull
```

This matches how Loom's sync client already works — it uses the hub URL as the base, not a hardcoded server root. No Loom client changes are needed.

### LoomHub Server Routing

LoomHub routes these requests by extracting `{owner}/{loom}` from the path prefix:

| Incoming Request | Routing |
|-----------------|---------|
| `POST /{owner}/{loom}/api/v1/negotiate` | Resolve loom from `{owner}/{loom}`, handle negotiate |
| `POST /{owner}/{loom}/api/v1/push` | Resolve loom, handle send |
| `POST /{owner}/{loom}/api/v1/pull` | Resolve loom, handle receive |

```go
// Router setup
r.Route("/{owner}/{loom}/api/v1", func(r chi.Router) {
    r.Use(resolveLoomMiddleware) // extracts owner/loom, loads loom DB
    r.Post("/negotiate", s.handleNegotiate)
    r.Post("/push", s.handleSend)
    r.Post("/pull", s.handleReceive)
})
```

The `resolveLoomMiddleware` resolves the `{owner}/{loom}` path into a loom database connection and injects it into the request context. The sync handlers then operate on the loom DB using Loom's standard protocol types — no translation needed.

### POST `/{owner}/{loom}/api/v1/negotiate`

Find common ancestor and determine what needs to sync.

**Request** (matches `NegotiateRequest` from Loom):
```json
{
  "project_id": "my-app-id",
  "streams": [
    {
      "stream_id": "01HZ...",
      "name": "main",
      "head_seq": 1234
    }
  ]
}
```

> Note: `project_id` is included per Loom's protocol. LoomHub validates it matches the resolved loom but primarily uses `{owner}/{loom}` for routing.

**Response** (matches `NegotiateResponse` from Loom):
```json
{
  "common_seqs": {
    "01HZ...": 1100
  },
  "server_seqs": {
    "01HZ...": 1180
  },
  "needs_push": true,
  "needs_pull": true
}
```

### POST `/{owner}/{loom}/api/v1/push`

Send operations and objects from client to server.

**Request** (matches `PushRequest` from Loom):
```json
{
  "project_id": "my-app-id",
  "stream_id": "01HZ...",
  "from_seq": 1100,
  "operations": [
    {
      "id": "01HZ...",
      "seq": 1101,
      "stream_id": "01HZ...",
      "space_id": "code",
      "entity_id": "src/main.go",
      "type": "modify",
      "path": "src/main.go",
      "delta": "...",
      "object_ref": "ab3def...",
      "parent_seq": 1100,
      "author": "flakerimi",
      "timestamp": "2026-03-12T10:30:00Z",
      "meta": {}
    }
  ],
  "objects": [
    {
      "hash": "ab3def...",
      "content": "<base64-encoded>"
    }
  ]
}
```

**Response** (matches `PushResponse` from Loom):
```json
{
  "ok": true,
  "applied": 134,
  "server_head": 1234
}
```

> **Server-side handling:** When objects arrive, LoomHub writes them to the shared object store (not per-loom). This is transparent to the client — the protocol is unchanged.

### POST `/{owner}/{loom}/api/v1/pull`

Receive new operations and objects from server to client.

**Request** (matches `PullRequest` from Loom):
```json
{
  "project_id": "my-app-id",
  "stream_id": "01HZ...",
  "from_seq": 1180
}
```

**Response** (matches `PullResponse` from Loom):
```json
{
  "operations": [...],
  "objects": [...],
  "server_head": 1234
}
```

---

## REST API

### Users

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/users` | Register new user |
| GET | `/api/v1/users/{username}` | Get user profile |
| PATCH | `/api/v1/user` | Update authenticated user's profile |
| GET | `/api/v1/user` | Get authenticated user |
| GET | `/api/v1/user/looms` | List authenticated user's looms |
| GET | `/api/v1/users/{username}/looms` | List user's public looms |

#### POST `/api/v1/users` — Register

```json
// Request
{
  "username": "flakerimi",
  "email": "flakerimi@example.com",
  "password": "...",
  "display_name": "Flakerimi"
}

// Response (201)
{
  "id": "01HZ...",
  "username": "flakerimi",
  "email": "flakerimi@example.com",
  "display_name": "Flakerimi",
  "created_at": "2026-03-12T10:00:00Z"
}
```

### Organizations

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/orgs` | Create organization |
| GET | `/api/v1/orgs/{name}` | Get organization |
| PATCH | `/api/v1/orgs/{name}` | Update organization |
| DELETE | `/api/v1/orgs/{name}` | Delete organization |
| GET | `/api/v1/orgs/{name}/members` | List members |
| PUT | `/api/v1/orgs/{name}/members/{username}` | Add member |
| DELETE | `/api/v1/orgs/{name}/members/{username}` | Remove member |
| GET | `/api/v1/orgs/{name}/looms` | List org looms |

### Looms

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/looms` | Create loom |
| GET | `/api/v1/looms/{owner}/{loom}` | Get loom |
| PATCH | `/api/v1/looms/{owner}/{loom}` | Update loom |
| DELETE | `/api/v1/looms/{owner}/{loom}` | Delete loom |
| POST | `/api/v1/looms/{owner}/{loom}/spin` | Spin a loom |

#### POST `/api/v1/looms` — Create Loom

```json
// Request
{
  "name": "my-app",
  "description": "My application",
  "visibility": "public",
  "owner": "flakerimi"
}

// Response (201)
{
  "id": "01HZ...",
  "owner": "flakerimi",
  "name": "my-app",
  "full_name": "flakerimi/my-app",
  "description": "My application",
  "visibility": "public",
  "default_stream": "main",
  "replicate_url": "https://loomhub.dev/flakerimi/my-app",
  "created_at": "2026-03-12T10:00:00Z"
}
```

### Loom Content

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/looms/{owner}/{loom}/streams` | List streams |
| GET | `/api/v1/looms/{owner}/{loom}/streams/{name}` | Get stream |
| GET | `/api/v1/looms/{owner}/{loom}/log` | Checkpoint log |
| GET | `/api/v1/looms/{owner}/{loom}/checkpoints/{id}` | Get checkpoint |
| GET | `/api/v1/looms/{owner}/{loom}/diff` | Diff between refs |
| GET | `/api/v1/looms/{owner}/{loom}/entities` | List entities |
| GET | `/api/v1/looms/{owner}/{loom}/entities/{path}` | Get entity content |
| GET | `/api/v1/looms/{owner}/{loom}/operations` | List operations |
| GET | `/api/v1/looms/{owner}/{loom}/search` | Search checkpoints |

#### GET `/api/v1/looms/{owner}/{loom}/log` — Checkpoint Log

```
GET /api/v1/looms/flakerimi/my-app/log?stream=main&limit=20&offset=0
```

```json
{
  "checkpoints": [
    {
      "id": "01HZ...",
      "title": "Add authentication",
      "summary": "JWT-based auth system",
      "author": "flakerimi",
      "timestamp": "2026-03-12T10:35:00Z",
      "source": "manual",
      "seq": 1234,
      "stream": "main",
      "space_summary": {
        "code": {"entities_changed": 5},
        "docs": {"entities_changed": 1}
      }
    }
  ],
  "total": 45,
  "has_more": true
}
```

#### GET `/api/v1/looms/{owner}/{loom}/diff`

```
GET /api/v1/looms/flakerimi/my-app/diff?from=01HX...&to=01HZ...&space=code
```

```json
{
  "from": "01HX...",
  "to": "01HZ...",
  "spaces": [
    {
      "space_id": "code",
      "entities": [
        {
          "path": "src/auth/login.go",
          "change": "modified",
          "hunks": [
            {
              "old_start": 10,
              "old_count": 5,
              "new_start": 10,
              "new_count": 8,
              "lines": [
                {"type": "context", "content": "func Login(w http.ResponseWriter, r *http.Request) {"},
                {"type": "delete", "content": "\treturn nil"},
                {"type": "add", "content": "\ttoken, err := generateJWT(user)"},
                {"type": "add", "content": "\tif err != nil {"},
                {"type": "add", "content": "\t\treturn err"},
                {"type": "add", "content": "\t}"}
              ]
            }
          ]
        }
      ]
    }
  ],
  "stats": {
    "spaces_changed": 2,
    "entities_changed": 6,
    "operations": 134
  }
}
```

#### GET `/api/v1/looms/{owner}/{loom}/entities`

Browse entities at a checkpoint or stream head.

```
GET /api/v1/looms/flakerimi/my-app/entities?stream=main&space=code&path=src/
```

```json
{
  "entities": [
    {
      "path": "src/main.go",
      "kind": "file",
      "space": "code",
      "size": 2048,
      "object_ref": "ab3def...",
      "mod_time": "2026-03-12T10:30:00Z"
    },
    {
      "path": "src/auth/",
      "kind": "directory",
      "children": 4
    }
  ]
}
```

### Weave Requests

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/looms/{owner}/{loom}/weave-requests` | Create WR |
| GET | `/api/v1/looms/{owner}/{loom}/weave-requests` | List WRs |
| GET | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}` | Get WR |
| PATCH | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}` | Update WR |
| POST | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}/weave` | Weave |
| POST | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}/close` | Close |
| GET | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}/diff` | WR diff |

#### POST — Create Weave Request

```json
// Request
{
  "title": "Add user authentication",
  "description": "Implements JWT-based auth with...",
  "source_stream": "feature/auth",
  "target_stream": "main"
}

// Response (201)
{
  "number": 42,
  "title": "Add user authentication",
  "author": "flakerimi",
  "source_stream": "feature/auth",
  "target_stream": "main",
  "status": "open",
  "created_at": "2026-03-12T11:00:00Z",
  "diff_stats": {
    "spaces_changed": 2,
    "entities_changed": 8,
    "operations": 56
  }
}
```

### Comments

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}/comments` | List comments |
| POST | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}/comments` | Add comment |
| PATCH | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}/comments/{id}` | Edit comment |
| DELETE | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}/comments/{id}` | Delete comment |

### Reviews

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}/reviews` | Submit review |
| GET | `/api/v1/looms/{owner}/{loom}/weave-requests/{number}/reviews` | List reviews |

### Pins

| Method | Path | Description |
|--------|------|-------------|
| PUT | `/api/v1/looms/{owner}/{loom}/pin` | Pin loom |
| DELETE | `/api/v1/looms/{owner}/{loom}/pin` | Unpin loom |
| GET | `/api/v1/user/pinned` | List pinned looms |

### Webhooks

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/looms/{owner}/{loom}/webhooks` | Create webhook |
| GET | `/api/v1/looms/{owner}/{loom}/webhooks` | List webhooks |
| PATCH | `/api/v1/looms/{owner}/{loom}/webhooks/{id}` | Update webhook |
| DELETE | `/api/v1/looms/{owner}/{loom}/webhooks/{id}` | Delete webhook |
| GET | `/api/v1/looms/{owner}/{loom}/webhooks/{id}/deliveries` | Delivery log |
| POST | `/api/v1/looms/{owner}/{loom}/webhooks/{id}/test` | Test webhook |

### Search

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/search/looms` | Search looms |
| GET | `/api/v1/search/users` | Search users |
| GET | `/api/v1/search/checkpoints` | Search checkpoints |

---

## Webhook Events

When configured events occur, LoomHub sends POST requests to webhook URLs.

### Event Types

| Event | Trigger | Payload |
|-------|---------|---------|
| `send` | Operations sent to loom | Stream info, operation count, checkpoint list |
| `checkpoint` | New checkpoint created | Checkpoint details |
| `weave_request` | WR opened/closed/woven | WR details, action type |
| `review` | Review submitted | Review details |
| `comment` | Comment posted | Comment details |
| `pin` | Loom pinned/unpinned | User, action |
| `spin` | Loom spun | Spin details |

### Webhook Payload Format

```json
{
  "event": "send",
  "timestamp": "2026-03-12T10:35:00Z",
  "sender": {
    "id": "01HZ...",
    "username": "flakerimi"
  },
  "loom": {
    "id": "01HZ...",
    "full_name": "flakerimi/my-app"
  },
  "payload": {
    "stream": "main",
    "operations_count": 134,
    "from_seq": 1100,
    "to_seq": 1234,
    "checkpoints": [
      {
        "id": "01HZ...",
        "title": "Add authentication"
      }
    ]
  }
}
```

Signed with HMAC-SHA256:
```
X-LoomHub-Signature: sha256=<hex-digest>
X-LoomHub-Event: send
X-LoomHub-Delivery: <delivery-id>
```

---

## Pagination

All list endpoints support cursor-based pagination:

```
GET /api/v1/looms/flakerimi/my-app/log?limit=20&after=01HZ...
```

Response includes:
```json
{
  "items": [...],
  "total": 156,
  "has_more": true,
  "next_cursor": "01HY..."
}
```

## Error Format

```json
{
  "error": {
    "code": "not_found",
    "message": "Loom not found",
    "details": {}
  }
}
```

Standard error codes: `bad_request`, `unauthorized`, `forbidden`, `not_found`, `conflict`, `rate_limited`, `internal_error`.
