# LoomHub — Development Roadmap

Phased implementation plan, from MVP to full platform.

---

## Phase 1 — Foundation (MVP)

Minimal viable server: send, receive, and browse looms.

### 1.1 Project Setup
- [ ] Go module, dependencies, Makefile
- [ ] Config loading (TOML)
- [ ] HTTP server with chi router
- [ ] Hub database init + schema migrations

### 1.2 User Accounts
- [ ] Registration (username, email, password)
- [ ] Login (JWT sessions)
- [ ] Auth middleware
- [ ] Access tokens (for CLI)

### 1.3 Loom Management
- [ ] Create loom (API + web)
- [ ] List user looms
- [ ] Delete loom
- [ ] Loom settings (name, description, visibility)

### 1.4 Sync Protocol
- [ ] Negotiate endpoint
- [ ] Send endpoint (receive operations + objects from client)
- [ ] Receive endpoint (send operations + objects to client)
- [ ] Shared object store
- [ ] Per-loom Loom database management

### 1.5 Web UI — Basics
- [ ] Landing page
- [ ] Login / register pages
- [ ] User profile page
- [ ] Loom overview page
- [ ] Entity tree browser
- [ ] Entity content viewer (code highlighting, markdown)
- [ ] Checkpoint log page

**Milestone: Users can create accounts, send/receive Loom projects, and browse them on the web.**

---

## Phase 2 — Collaboration

Weave requests, reviews, and social features.

### 2.1 Weave Requests
- [ ] Create WR (stream → stream)
- [ ] List / filter WRs
- [ ] WR detail page with diff
- [ ] WR comments (general + inline)
- [ ] Weave / close WR
- [ ] WR status tracking

### 2.2 Reviews
- [ ] Submit review (approve / request changes / comment)
- [ ] Review summary on WR page
- [ ] Review requirements for weave (optional)

### 2.3 Diff Viewer
- [ ] Multi-space diff rendering
- [ ] Unified / side-by-side toggle
- [ ] Inline comments on diff lines
- [ ] Expand context around hunks

### 2.4 Social Features
- [ ] Pin looms
- [ ] Activity feed (dashboard)
- [ ] User activity on profile

### 2.5 Organizations
- [ ] Create org
- [ ] Org members (owner, admin, member)
- [ ] Org loom ownership

**Milestone: Teams can collaborate with weave requests, code review, and org-based ownership.**

---

## Phase 3 — Platform

Webhooks, search, spins, and polish.

### 3.1 Webhooks
- [ ] Webhook CRUD
- [ ] Event dispatch (send, checkpoint, WR, review)
- [ ] HMAC signing
- [ ] Delivery log with retry
- [ ] Test webhook button

### 3.2 Search
- [ ] Loom search (name, description)
- [ ] Checkpoint search (FTS on title/summary)
- [ ] User search
- [ ] Explore page (trending, recent)

### 3.3 Spinning
- [ ] Spin loom
- [ ] Cross-spin weave requests
- [ ] Spin network visualization

### 3.4 Labels & Milestones
- [ ] Labels for WRs
- [ ] Label filtering

### 3.5 Admin Panel
- [ ] Admin dashboard (system stats)
- [ ] User management
- [ ] Loom management
- [ ] Server configuration

### 3.6 Polish
- [ ] Responsive design
- [ ] Keyboard navigation
- [ ] Notification indicators
- [ ] Error pages (404, 500)
- [ ] Loading states

**Milestone: Full-featured hosting platform with webhooks, search, and admin tools.**

---

## Phase 4 — Scale & Extend (Future)

### 4.1 Performance
- [ ] Object store → S3-compatible backend option
- [ ] Hub database → PostgreSQL option
- [ ] CDN for static assets
- [ ] Caching layer (object metadata, entity trees)

### 4.2 CI/CD Integration
- [ ] Built-in runner system (or webhook-based)
- [ ] Status checks on weave requests
- [ ] Protected streams (require checks to pass)

### 4.3 Real-Time
- [ ] SSE/WebSocket for live updates (WR comments, send events)
- [ ] Live collaboration indicators

### 4.4 API Extensions
- [ ] GraphQL API
- [ ] OAuth2 provider (third-party apps)
- [ ] SSH-based send/receive

### 4.5 Import/Export
- [ ] Import from Git/GitHub
- [ ] Export to standard formats
- [ ] Loom transfer between owners

---

## Implementation Order

For the initial build, implement in this order:

```
1. cmd/loomhub/main.go          — Entry point + CLI commands
2. internal/config/              — Config loading
3. internal/database/            — Hub DB + schema
4. internal/models/              — Core types
5. internal/store/users.go       — User CRUD
6. internal/auth/                — JWT + passwords + middleware
7. internal/server/              — HTTP server setup
8. internal/api/users.go         — User registration + login endpoints
9. internal/store/looms.go       — Loom CRUD
10. internal/api/looms.go        — Loom endpoints
11. internal/sync/               — Negotiate + send + receive
12. internal/web/                — Web UI pages
13. internal/store/weave_requests.go — WR CRUD
14. internal/api/weave_requests.go   — WR endpoints
15. internal/webhook/            — Event dispatch
```

Each step builds on the previous. Tests are written alongside each component.
