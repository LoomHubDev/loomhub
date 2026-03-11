# LoomHub — Project Structure

Go project layout and development guide.

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
│   │   ├── user.go                 # User, SSHKey, AccessToken
│   │   ├── org.go                  # Organization, OrgMember
│   │   ├── repo.go                 # Repository, Collaborator
│   │   ├── merge_request.go        # MergeRequest, Comment, Review
│   │   ├── webhook.go              # Webhook, Delivery
│   │   ├── activity.go             # Activity feed
│   │   └── label.go                # Labels
│   │
│   ├── store/
│   │   ├── users.go                # User CRUD
│   │   ├── orgs.go                 # Org CRUD
│   │   ├── repos.go                # Repo CRUD
│   │   ├── merge_requests.go       # MR CRUD
│   │   ├── comments.go             # Comment CRUD
│   │   ├── reviews.go              # Review CRUD
│   │   ├── webhooks.go             # Webhook CRUD
│   │   ├── activities.go           # Activity logging
│   │   ├── stars.go                # Stars
│   │   └── tokens.go              # Access tokens
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
│   │   └── objects.go              # Object store management
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
│   │   └── middleware.go           # API-specific middleware (JSON, pagination)
│   │
│   ├── web/
│   │   ├── router.go               # Web UI router setup
│   │   ├── pages/                  # Page handlers
│   │   │   ├── home.go
│   │   │   ├── auth.go             # Login, register
│   │   │   ├── profile.go          # User/org profiles
│   │   │   ├── repo.go             # Repo pages
│   │   │   ├── tree.go             # Entity tree browser
│   │   │   ├── log.go              # Checkpoint log
│   │   │   ├── diff.go             # Diff viewer
│   │   │   ├── merge_request.go    # MR pages
│   │   │   ├── settings.go         # Settings pages
│   │   │   └── admin.go            # Admin pages
│   │   └── templates/              # templ templates
│   │       ├── layouts/
│   │       │   ├── base.templ      # Base HTML layout
│   │       │   ├── app.templ       # Authenticated layout
│   │       │   └── public.templ    # Public layout
│   │       ├── components/
│   │       │   ├── nav.templ       # Navigation bar
│   │       │   ├── footer.templ
│   │       │   ├── pagination.templ
│   │       │   ├── avatar.templ
│   │       │   ├── diff.templ      # Diff rendering component
│   │       │   ├── entity_tree.templ
│   │       │   ├── checkpoint.templ
│   │       │   └── flash.templ     # Flash messages
│   │       └── pages/
│   │           ├── home.templ
│   │           ├── login.templ
│   │           ├── register.templ
│   │           ├── explore.templ
│   │           ├── profile.templ
│   │           ├── repo.templ
│   │           ├── tree.templ
│   │           ├── entity.templ
│   │           ├── log.templ
│   │           ├── checkpoint.templ
│   │           ├── diff.templ
│   │           ├── streams.templ
│   │           ├── mr_list.templ
│   │           ├── mr_new.templ
│   │           ├── mr_detail.templ
│   │           ├── settings.templ
│   │           └── admin.templ
│   │
│   ├── render/
│   │   ├── code.go                 # Syntax highlighting (chroma)
│   │   ├── markdown.go             # Markdown rendering (goldmark)
│   │   ├── diff.go                 # Diff formatting
│   │   └── json.go                 # JSON tree rendering
│   │
│   ├── webhook/
│   │   ├── dispatch.go             # Webhook event dispatch
│   │   ├── delivery.go             # HTTP delivery with retries
│   │   └── sign.go                 # HMAC-SHA256 signing
│   │
│   └── server/
│       └── server.go               # HTTP server setup, middleware chain
│
├── static/
│   ├── css/
│   │   └── style.css               # Tailwind output
│   ├── js/
│   │   └── htmx.min.js            # htmx library
│   └── img/
│       └── logo.svg
│
├── test/
│   ├── integration/
│   │   ├── sync_test.go            # Push/pull integration tests
│   │   ├── api_test.go             # REST API tests
│   │   └── web_test.go             # Web UI smoke tests
│   └── testutil/
│       └── helpers.go              # Test helpers
│
├── docs/                           # Design documentation (this folder)
│
├── go.mod
├── go.sum
├── Makefile
├── tailwind.config.js
├── .gitignore
└── LICENSE
```

---

## Dependencies

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

    // Templates
    github.com/a-h/templ              // Type-safe templates

    // Rendering
    github.com/alecthomas/chroma/v2   // Syntax highlighting
    github.com/yuin/goldmark           // Markdown

    // Utilities
    github.com/oklog/ulid/v2          // ULID generation
    github.com/klauspost/compress     // Zstandard compression
    github.com/BurntSushi/toml        // Config parsing
)
```

---

## Build & Run

### Makefile

```makefile
.PHONY: build run dev test clean

# Build
build: templ-generate tailwind
	go build -o bin/loomhub ./cmd/loomhub

# Development (with hot reload)
dev:
	templ generate --watch &
	npx tailwindcss -i static/css/input.css -o static/css/style.css --watch &
	go run ./cmd/loomhub serve --dev

# Generate templ templates
templ-generate:
	templ generate

# Build Tailwind CSS
tailwind:
	npx tailwindcss -i static/css/input.css -o static/css/style.css --minify

# Run tests
test:
	go test ./...

# Clean
clean:
	rm -rf bin/
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
2. Install Go 1.23+, Node.js (for Tailwind), templ CLI
3. `make dev` — starts everything with hot reload
4. Server runs at `http://localhost:3000`
5. First user to register becomes admin

### Running Tests

```bash
make test                    # All tests
go test ./internal/store/... # Store tests only
go test ./internal/sync/...  # Sync tests only
go test ./test/integration/  # Integration tests
```
