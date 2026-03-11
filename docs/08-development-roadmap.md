# LoomHub — Development Roadmap

Phased implementation plan, from MVP to full platform.

---

## Phase 1 — Foundation (MVP)

Minimal viable server: push, pull, and browse repos.

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

### 1.3 Repository Management
- [ ] Create repo (API + web)
- [ ] List user repos
- [ ] Delete repo
- [ ] Repo settings (name, description, visibility)

### 1.4 Sync Protocol
- [ ] Negotiate endpoint
- [ ] Push endpoint (receive operations + objects)
- [ ] Pull endpoint (send operations + objects)
- [ ] Shared object store
- [ ] Per-repo Loom database management

### 1.5 Web UI — Basics
- [ ] Landing page
- [ ] Login / register pages
- [ ] User profile page
- [ ] Repo overview page
- [ ] Entity tree browser
- [ ] Entity content viewer (code highlighting, markdown)
- [ ] Checkpoint log page

**Milestone: Users can create accounts, push/pull Loom repos, and browse them on the web.**

---

## Phase 2 — Collaboration

Merge requests, reviews, and social features.

### 2.1 Merge Requests
- [ ] Create MR (stream → stream)
- [ ] List / filter MRs
- [ ] MR detail page with diff
- [ ] MR comments (general + inline)
- [ ] Merge / close MR
- [ ] MR status tracking

### 2.2 Reviews
- [ ] Submit review (approve / request changes / comment)
- [ ] Review summary on MR page
- [ ] Review requirements for merge (optional)

### 2.3 Diff Viewer
- [ ] Multi-space diff rendering
- [ ] Unified / side-by-side toggle
- [ ] Inline comments on diff lines
- [ ] Expand context around hunks

### 2.4 Social Features
- [ ] Star repos
- [ ] Activity feed (dashboard)
- [ ] User activity on profile

### 2.5 Organizations
- [ ] Create org
- [ ] Org members (owner, admin, member)
- [ ] Org repo ownership

**Milestone: Teams can collaborate with merge requests, code review, and org-based ownership.**

---

## Phase 3 — Platform

Webhooks, search, forks, and polish.

### 3.1 Webhooks
- [ ] Webhook CRUD
- [ ] Event dispatch (push, checkpoint, MR, review)
- [ ] HMAC signing
- [ ] Delivery log with retry
- [ ] Test webhook button

### 3.2 Search
- [ ] Repository search (name, description)
- [ ] Checkpoint search (FTS on title/summary)
- [ ] User search
- [ ] Explore page (trending, recent)

### 3.3 Forking
- [ ] Fork repository
- [ ] Cross-fork merge requests
- [ ] Fork network visualization

### 3.4 Labels & Milestones
- [ ] Labels for MRs
- [ ] Label filtering

### 3.5 Admin Panel
- [ ] Admin dashboard (system stats)
- [ ] User management
- [ ] Repo management
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
- [ ] Status checks on merge requests
- [ ] Protected streams (require checks to pass)

### 4.3 Real-Time
- [ ] SSE/WebSocket for live updates (MR comments, push events)
- [ ] Live collaboration indicators

### 4.4 API Extensions
- [ ] GraphQL API
- [ ] OAuth2 provider (third-party apps)
- [ ] SSH-based push/pull

### 4.5 Import/Export
- [ ] Import from Git/GitHub
- [ ] Export to standard formats
- [ ] Repository transfer between owners

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
9. internal/store/repos.go       — Repo CRUD
10. internal/api/repos.go        — Repo endpoints
11. internal/sync/               — Negotiate + push + pull
12. internal/web/                — Web UI pages
13. internal/store/merge_requests.go — MR CRUD
14. internal/api/merge_requests.go   — MR endpoints
15. internal/webhook/            — Event dispatch
```

Each step builds on the previous. Tests are written alongside each component.
