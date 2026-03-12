# LoomHub Findings

Current code review findings.

## P1

1. ~~`internal/api/sync.go`~~
   ~~Sync endpoints have no auth or permission enforcement.~~
   **FIXED**: Push requires auth + write/admin access. Private looms require auth + any access for negotiate/pull. Content browsing endpoints check private loom visibility via `checkLoomRead`. Push activity logged.

2. `internal/auth/middleware.go`
   Access-token scopes are stored but never enforced. A `read` token currently has write-level power.

3. ~~`internal/api/looms.go`~~
   ~~The org/collaborator permission model is not implemented.~~
   **FIXED**: `CollaboratorStore` implements full CRUD on `loom_collaborators` table. `CheckAccess()` resolves effective permission (admin > owner > collaborator > public-read). Loom CRUD, sync endpoints, and content browsing all use permission checks. API endpoints: `GET/POST/DELETE /looms/{owner}/{loom}/collaborators`. Frontend settings page added.

## P2

4. `internal/sync/handler.go`
   First send creates streams with `name = stream_id`, so server-side stream names become opaque IDs instead of human-readable names.

## Validation

- `go test ./...`
- `go test -race ./...`

## Test Gaps

- ~~no coverage for sync authz~~ **COVERED**: `TestSyncEndpointStubs` tests auth push + unauth rejection
- no coverage for PAT scope enforcement
- ~~no coverage for org/collaborator permissions~~ **COVERED**: `TestCollaborators` tests add/list/remove, push with/without access, outsider rejection
- no coverage for stream-name preservation on first sync
