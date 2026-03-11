# LoomHub — API Design

LoomHub exposes three API surfaces:

1. **Sync API** — Loom's native push/pull protocol (used by `loom push`/`loom pull`)
2. **REST API** — CRUD for hub-level resources (users, repos, merge requests, etc.)
3. **Web UI** — Server-rendered HTML pages

All APIs share the same HTTP server and authentication layer.

---

## Base URL & Versioning

```
https://hub.example.com/api/v1/...     # REST API
https://hub.example.com/api/v1/sync/...  # Sync API
https://hub.example.com/{owner}/{repo}   # Web UI
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
| `read` | Read public + private repos the user has access to |
| `write` | Push to repos, create MRs, post comments |
| `admin` | Manage repos, users, webhooks, permissions |

---

## Sync API

These endpoints implement Loom's native sync protocol. The `loom` CLI calls these directly.

### POST `/api/v1/sync/{owner}/{repo}/negotiate`

Find common ancestor and determine what needs to sync.

**Request:**
```json
{
  "streams": [
    {
      "stream_id": "01HZ...",
      "name": "main",
      "head_seq": 1234
    }
  ]
}
```

**Response:**
```json
{
  "common_seqs": {
    "main": 1100
  },
  "server_seqs": {
    "main": 1180
  },
  "needs_push": true,
  "needs_pull": true
}
```

### POST `/api/v1/sync/{owner}/{repo}/push`

Push operations and objects from client to server.

**Request:**
```json
{
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
      "content": "<base64-encoded>",
      "size": 4096,
      "compressed": true
    }
  ],
  "checkpoints": [
    {
      "id": "01HZ...",
      "stream_id": "01HZ...",
      "seq": 1234,
      "title": "Add authentication",
      "summary": "JWT-based auth system",
      "author": "flakerimi",
      "timestamp": "2026-03-12T10:35:00Z",
      "source": "manual",
      "spaces": [],
      "parent_id": "01HY..."
    }
  ]
}
```

**Response:**
```json
{
  "ok": true,
  "applied": 134,
  "server_head": 1234
}
```

### POST `/api/v1/sync/{owner}/{repo}/pull`

Pull new operations and objects from server to client.

**Request:**
```json
{
  "stream_id": "01HZ...",
  "from_seq": 1180
}
```

**Response:**
```json
{
  "operations": [...],
  "objects": [...],
  "checkpoints": [...],
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
| GET | `/api/v1/user/repos` | List authenticated user's repos |
| GET | `/api/v1/users/{username}/repos` | List user's public repos |

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
| GET | `/api/v1/orgs/{name}/repos` | List org repos |

### Repositories

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/repos` | Create repository |
| GET | `/api/v1/repos/{owner}/{repo}` | Get repository |
| PATCH | `/api/v1/repos/{owner}/{repo}` | Update repository |
| DELETE | `/api/v1/repos/{owner}/{repo}` | Delete repository |
| POST | `/api/v1/repos/{owner}/{repo}/fork` | Fork repository |

#### POST `/api/v1/repos` — Create Repository

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
  "clone_url": "https://hub.example.com/flakerimi/my-app",
  "created_at": "2026-03-12T10:00:00Z"
}
```

### Repository Content

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/repos/{owner}/{repo}/streams` | List streams |
| GET | `/api/v1/repos/{owner}/{repo}/streams/{name}` | Get stream |
| GET | `/api/v1/repos/{owner}/{repo}/log` | Checkpoint log |
| GET | `/api/v1/repos/{owner}/{repo}/checkpoints/{id}` | Get checkpoint |
| GET | `/api/v1/repos/{owner}/{repo}/diff` | Diff between refs |
| GET | `/api/v1/repos/{owner}/{repo}/entities` | List entities |
| GET | `/api/v1/repos/{owner}/{repo}/entities/{path}` | Get entity content |
| GET | `/api/v1/repos/{owner}/{repo}/operations` | List operations |
| GET | `/api/v1/repos/{owner}/{repo}/search` | Search checkpoints |

#### GET `/api/v1/repos/{owner}/{repo}/log` — Checkpoint Log

```
GET /api/v1/repos/flakerimi/my-app/log?stream=main&limit=20&offset=0
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

#### GET `/api/v1/repos/{owner}/{repo}/diff`

```
GET /api/v1/repos/flakerimi/my-app/diff?from=01HX...&to=01HZ...&space=code
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

#### GET `/api/v1/repos/{owner}/{repo}/entities`

Browse entities at a checkpoint or stream head.

```
GET /api/v1/repos/flakerimi/my-app/entities?stream=main&space=code&path=src/
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

### Merge Requests

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/repos/{owner}/{repo}/merge-requests` | Create MR |
| GET | `/api/v1/repos/{owner}/{repo}/merge-requests` | List MRs |
| GET | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}` | Get MR |
| PATCH | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}` | Update MR |
| POST | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}/merge` | Merge |
| POST | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}/close` | Close |
| GET | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}/diff` | MR diff |

#### POST — Create Merge Request

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
| GET | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}/comments` | List comments |
| POST | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}/comments` | Add comment |
| PATCH | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}/comments/{id}` | Edit comment |
| DELETE | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}/comments/{id}` | Delete comment |

### Reviews

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}/reviews` | Submit review |
| GET | `/api/v1/repos/{owner}/{repo}/merge-requests/{number}/reviews` | List reviews |

### Stars

| Method | Path | Description |
|--------|------|-------------|
| PUT | `/api/v1/repos/{owner}/{repo}/star` | Star repo |
| DELETE | `/api/v1/repos/{owner}/{repo}/star` | Unstar repo |
| GET | `/api/v1/user/starred` | List starred repos |

### Webhooks

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/repos/{owner}/{repo}/webhooks` | Create webhook |
| GET | `/api/v1/repos/{owner}/{repo}/webhooks` | List webhooks |
| PATCH | `/api/v1/repos/{owner}/{repo}/webhooks/{id}` | Update webhook |
| DELETE | `/api/v1/repos/{owner}/{repo}/webhooks/{id}` | Delete webhook |
| GET | `/api/v1/repos/{owner}/{repo}/webhooks/{id}/deliveries` | Delivery log |
| POST | `/api/v1/repos/{owner}/{repo}/webhooks/{id}/test` | Test webhook |

### Search

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/search/repos` | Search repositories |
| GET | `/api/v1/search/users` | Search users |
| GET | `/api/v1/search/checkpoints` | Search checkpoints |

---

## Webhook Events

When configured events occur, LoomHub sends POST requests to webhook URLs.

### Event Types

| Event | Trigger | Payload |
|-------|---------|---------|
| `push` | Operations pushed to repo | Stream info, operation count, checkpoint list |
| `checkpoint` | New checkpoint created | Checkpoint details |
| `merge_request` | MR opened/closed/merged | MR details, action type |
| `review` | Review submitted | Review details |
| `comment` | Comment posted | Comment details |
| `star` | Repo starred/unstarred | User, action |
| `fork` | Repo forked | Fork details |

### Webhook Payload Format

```json
{
  "event": "push",
  "timestamp": "2026-03-12T10:35:00Z",
  "sender": {
    "id": "01HZ...",
    "username": "flakerimi"
  },
  "repository": {
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
X-LoomHub-Event: push
X-LoomHub-Delivery: <delivery-id>
```

---

## Pagination

All list endpoints support cursor-based pagination:

```
GET /api/v1/repos/flakerimi/my-app/log?limit=20&after=01HZ...
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
    "message": "Repository not found",
    "details": {}
  }
}
```

Standard error codes: `bad_request`, `unauthorized`, `forbidden`, `not_found`, `conflict`, `rate_limited`, `internal_error`.
