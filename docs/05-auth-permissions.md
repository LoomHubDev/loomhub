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

Unauthenticated requests can access public looms (read-only).

---

## Permission Model

### Loom Access Levels

| Level | Capabilities |
|-------|-------------|
| **none** | No access (private looms) |
| **read** | View loom, replicate, receive |
| **write** | Send, create streams, create WRs |
| **admin** | Settings, collaborators, webhooks, delete |

### How Access Is Determined

For a given user + loom, access is resolved in order:

```
1. Is user a site admin?          → admin access to everything
2. Is user the loom owner?        → admin access
3. Is user a loom collaborator?   → collaborator's permission level
4. Is user an org member?         → based on org role:
   - owner/admin                  → admin access
   - member                       → read access to all org looms
5. Is loom public?                → read access
6. Otherwise                      → no access
```

> Note: Loom visibility is `public` or `private`. Org members get read access to all org looms (including private ones). For write access, org members must be added as loom collaborators.

### Organization Roles

| Role | Loom Access | Org Management |
|------|------------|----------------|
| **owner** | Admin to all org looms | Full org settings, billing, delete |
| **admin** | Admin to all org looms | Manage members, create looms |
| **member** | Read to all org looms | View members |

### Weave Request Permissions

| Action | Required |
|--------|----------|
| View WR | Read access to loom |
| Create WR | Write access to loom |
| Comment on WR | Read access to loom |
| Submit review | Read access to loom |
| Weave WR | Write access + WR approved (or admin) |
| Close WR | WR author or write access |

### Sync Permissions

| Action | Required |
|--------|----------|
| Receive (public loom) | No auth required |
| Receive (private loom) | Read access |
| Send | Write access |

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
| `/{owner}/{loom}/api/v1/push` | 60 per minute per user |
| `/{owner}/{loom}/api/v1/pull` | 120 per minute per user |
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
sync_send_per_minute = 60
sync_receive_per_minute = 120
```
