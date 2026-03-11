# LoomHub — Project Structure

Go + Vue monorepo layout and development guide.

---

## Directory Layout

```
loomhub/
├── cmd/
│   └── loomhub/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── config/
│   │   └── config.go               # Server configuration (TOML)
│   │
│   ├── database/
│   │   ├── hub.go                  # Hub database init & migrations
│   │   ├── schema.go               # Hub schema definitions
│   │   └── repo.go                 # Per-repo Loom database access
│   │
│   ├── models/
│   │   ├── owner.go                # Owner (shared namespace)
│   │   ├── user.go                 # User, SSHKey, AccessToken
│   │   ├── org.go                  # Organization, OrgMember
│   │   ├── repo.go                 # Repository, Collaborator
│   │   ├── merge_request.go        # MergeRequest, Comment, Review
│   │   ├── webhook.go              # Webhook, Delivery
│   │   ├── activity.go             # Activity feed
│   │   └── label.go                # Labels
│   │
│   ├── store/
│   │   ├── owners.go               # Owner namespace CRUD
│   │   ├── users.go                # User CRUD
│   │   ├── orgs.go                 # Org CRUD
│   │   ├── repos.go                # Repo CRUD
│   │   ├── merge_requests.go       # MR CRUD
│   │   ├── comments.go             # Comment CRUD
│   │   ├── reviews.go              # Review CRUD
│   │   ├── webhooks.go             # Webhook CRUD
│   │   ├── activities.go           # Activity logging
│   │   ├── stars.go                # Stars
│   │   └── tokens.go               # Access tokens
│   │
│   ├── auth/
│   │   ├── jwt.go                  # JWT generation & validation
│   │   ├── password.go             # bcrypt hashing & verification
│   │   ├── middleware.go           # Auth middleware
│   │   └── permissions.go          # Permission checks
│   │
│   ├── sync/
│   │   ├── handler.go              # Negotiate/push/pull HTTP handlers
│   │   ├── negotiate.go            # Negotiate logic
│   │   ├── push.go                 # Push logic (receive ops + objects)
│   │   ├── pull.go                 # Pull logic (send ops + objects)
│   │   └── objects.go              # Shared object store management
│   │
│   ├── api/
│   │   ├── router.go               # REST API router setup
│   │   ├── users.go                # User endpoints
│   │   ├── orgs.go                 # Org endpoints
│   │   ├── repos.go                # Repo endpoints
│   │   ├── content.go              # Entity/checkpoint/log endpoints
│   │   ├── merge_requests.go       # MR endpoints
│   │   ├── search.go               # Search endpoints
│   │   ├── webhooks.go             # Webhook endpoints
│   │   └── middleware.go           # API middleware (JSON, pagination, CORS)
│   │
│   ├── render/
│   │   └── diff.go                 # Server-side diff computation for API
│   │
│   ├── webhook/
│   │   ├── dispatch.go             # Webhook event dispatch
│   │   ├── delivery.go             # HTTP delivery with retries
│   │   └── sign.go                 # HMAC-SHA256 signing
│   │
│   └── server/
│       └── server.go               # HTTP server setup, SPA fallback, embed.FS
│
├── frontend/                       # Vue 3 SPA
│   ├── src/
│   │   ├── main.ts
│   │   ├── App.vue
│   │   ├── router/
│   │   ├── stores/
│   │   ├── api/
│   │   ├── layouts/
│   │   ├── views/
│   │   ├── components/
│   │   └── utils/
│   ├── index.html
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   └── package.json
│
├── test/
│   ├── integration/
│   │   ├── sync_test.go            # Push/pull integration tests
│   │   ├── api_test.go             # REST API tests
│   │   └── helpers_test.go         # Test helpers
│   └── testutil/
│       └── helpers.go              # Test helpers (test DB, test server)
│
├── docs/                           # Design documentation
│
├── go.mod
├── go.sum
├── Makefile
├── .gitignore
└── LICENSE
```

---

## Dependencies

### Go (backend)

```go
// go.mod
module github.com/flakerimi/loomhub

go 1.23

require (
    // HTTP
    github.com/go-chi/chi/v5          // Router
    github.com/go-chi/cors            // CORS middleware

    // Database
    modernc.org/sqlite                 // Pure Go SQLite

    // Auth
    golang.org/x/crypto               // bcrypt
    github.com/golang-jwt/jwt/v5      // JWT

    // Utilities
    github.com/oklog/ulid/v2          // ULID generation
    github.com/klauspost/compress     // Zstandard compression
    github.com/BurntSushi/toml        // Config parsing
)
```

### Node.js (frontend)

```json
{
  "dependencies": {
    "vue": "^3.5",
    "vue-router": "^4.5",
    "pinia": "^3",
    "ofetch": "^1",
    "@heroicons/vue": "^2"
  },
  "devDependencies": {
    "vite": "^6",
    "@vitejs/plugin-vue": "^5",
    "typescript": "^5.7",
    "tailwindcss": "^4",
    "shiki": "^3",
    "markdown-it": "^14"
  }
}
```

---

## Build & Run

### Makefile

```makefile
.PHONY: build dev test clean

# Full production build
build: frontend-build
	go build -o bin/loomhub ./cmd/loomhub

# Build Vue SPA
frontend-build:
	cd frontend && npm install && npm run build

# Development — run Go backend + Vue dev server
dev:
	@echo "Starting Go backend on :3000..."
	go run ./cmd/loomhub serve --dev &
	@echo "Starting Vue dev server on :5173..."
	cd frontend && npm run dev

# Run all tests
test: test-go test-frontend

test-go:
	go test ./...

test-frontend:
	cd frontend && npm test

# Clean
clean:
	rm -rf bin/ frontend/dist/
```

### SPA Embedding

In production, the Vue `dist/` is embedded into the Go binary:

```go
// internal/server/server.go
//go:embed all:../../frontend/dist
var frontendDist embed.FS

func (s *Server) setupRoutes() {
    // API routes
    s.router.Route("/api/v1", s.apiRoutes)

    // SPA fallback — serve index.html for all non-API routes
    s.router.Get("/*", s.spaHandler())
}
```

### CLI Commands

```bash
# Start server
loomhub serve --data /var/lib/loomhub --port 3000

# Create admin user
loomhub admin create-user --username admin --email admin@example.com --admin

# Run database migrations
loomhub migrate

# Verify system integrity
loomhub doctor

# Garbage collect orphaned objects
loomhub gc

# Show server info
loomhub info
```

---

## Configuration File

`config.toml`:

```toml
[server]
host = "0.0.0.0"
port = 3000
base_url = "https://hub.example.com"
data_dir = "/var/lib/loomhub"

[auth]
jwt_secret = ""              # Auto-generated if empty
jwt_expiry = "168h"
bcrypt_cost = 12
registration_open = true

[storage]
max_object_size = "100MB"
max_push_size = "500MB"
max_repos_per_user = 100

[rate_limit]
enabled = true
api_per_minute = 300

[logging]
level = "info"               # debug, info, warn, error
format = "text"              # text, json
```

---

## Development Workflow

1. Clone repo
2. Install Go 1.23+, Node.js 20+
3. `cd frontend && npm install`
4. Terminal 1: `go run ./cmd/loomhub serve --dev` (backend on :3000)
5. Terminal 2: `cd frontend && npm run dev` (Vite on :5173 with proxy to :3000)
6. Open `http://localhost:5173`
7. First user to register becomes admin

### Running Tests

```bash
# Go tests
go test ./...                     # All
go test ./internal/store/...      # Store only
go test ./internal/sync/...       # Sync only
go test ./test/integration/       # Integration
go test -race ./...               # Race detector

# Frontend tests
cd frontend && npm test           # Vitest
cd frontend && npm run test:e2e   # Playwright (future)
```
