# LoomHub — Authentication & Permissions

---

## Authentication

### Registration Flow

1. User submits username, email, password
2. Server validates uniqueness (username and email)
3. Password hashed with bcrypt (cost 12)
4. User record created
5. Session JWT issued

### Login Flow

1. User submits email/username + password
2. Server verifies bcrypt hash
3. JWT session token issued (httpOnly cookie)
4. Token expires after 7 days (configurable)

### JWT Structure

```json
{
  "sub": "01HZ...",
  "username": "flakerimi",
  "is_admin": false,
  "exp": 1741996800,
  "iat": 1741392000
}
```

Signed with HMAC-SHA256 using server secret.

### Access Tokens

For CLI and API access:

```
loomhub-pat_01HZ...random-bytes...
```

- Stored as SHA-256 hash in database
- Scoped: `read`, `write`, `admin`
- Optional expiration
- Last-used tracking

### Authentication Middleware

```
Request → Extract token (cookie or header)
        → Validate token
        → Load user from DB
        → Attach user to request context
        → Next handler
```

Unauthenticated requests can access public repos (read-only).

---

## Permission Model

### Repository Access Levels

| Level | Capabilities |
|-------|-------------|
| **none** | No access (private repos) |
| **read** | View repo, clone, pull |
| **write** | Push, create streams, create MRs |
| **admin** | Settings, collaborators, webhooks, delete |

### How Access Is Determined

For a given user + repository, access is resolved in order:

```
1. Is user a site admin?          → admin access to everything
2. Is user the repo owner?        → admin access
3. Is user a repo collaborator?   → collaborator's permission level
4. Is user an org member?         → based on org role:
   - owner/admin                  → admin access
   - member                       → read (public), write (internal)
5. Is repo public?                → read access
6. Otherwise                      → no access
```

### Organization Roles

| Role | Repo Access | Org Management |
|------|------------|----------------|
| **owner** | Admin to all org repos | Full org settings, billing, delete |
| **admin** | Admin to all org repos | Manage members, create repos |
| **member** | Read to all org repos (write if internal) | View members |

### Merge Request Permissions

| Action | Required |
|--------|----------|
| View MR | Read access to repo |
| Create MR | Write access to repo |
| Comment on MR | Read access to repo |
| Submit review | Read access to repo |
| Merge MR | Write access + MR approved (or admin) |
| Close MR | MR author or write access |

### Sync Permissions

| Action | Required |
|--------|----------|
| Pull (public repo) | No auth required |
| Pull (private repo) | Read access |
| Push | Write access |

---

## Security

### Password Requirements

- Minimum 8 characters
- No maximum length (bcrypt truncates at 72 bytes)
- No complexity rules (length is more important)

### Rate Limiting

| Endpoint | Limit |
|----------|-------|
| `/login` | 10 per minute per IP |
| `/register` | 5 per hour per IP |
| `/api/v1/sync/*/push` | 60 per minute per user |
| `/api/v1/sync/*/pull` | 120 per minute per user |
| General API | 300 per minute per user |

Implemented with token bucket per IP/user, stored in memory (no Redis needed).

### Session Security

- JWT in httpOnly, Secure, SameSite=Lax cookie
- CSRF protection via SameSite + double-submit cookie for mutations
- Session revocation: token blocklist checked on each request

### Webhook Security

- Payloads signed with HMAC-SHA256
- Secret per webhook
- Delivery retries with exponential backoff (3 attempts)
- 10-second timeout per delivery

---

## First User / Admin Bootstrap

On first run with empty database:

1. First registered user automatically becomes site admin
2. Admin can promote other users via `/admin/users`
3. Admin can configure server settings (registration open/closed, etc.)

## Configuration

```toml
[auth]
jwt_secret = "..."           # Auto-generated on first run
jwt_expiry = "168h"          # 7 days
bcrypt_cost = 12
registration_open = true     # Can be closed by admin
require_email_verification = false  # Future: email verification

[rate_limit]
enabled = true
login_per_minute = 10
register_per_hour = 5
api_per_minute = 300
sync_push_per_minute = 60
sync_pull_per_minute = 120
```
