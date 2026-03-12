# LoomHub — CLI Reference

The `loomhub` binary serves as both the server and admin tool.

---

## Commands

### `loomhub serve`

Start the LoomHub server.

```bash
loomhub serve [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.toml` | Path to config file |
| `--data` | `./data` | Data directory |
| `--port` | `3000` | HTTP port |
| `--host` | `0.0.0.0` | Bind address |
| `--dev` | `false` | Development mode (verbose logging, no caching) |

### `loomhub migrate`

Run database migrations.

```bash
loomhub migrate [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--data` | `./data` | Data directory |

### `loomhub admin create-user`

Create a user from the command line.

```bash
loomhub admin create-user --username NAME --email EMAIL [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--username` | (required) | Username |
| `--email` | (required) | Email address |
| `--password` | (prompted) | Password (prompted if not provided) |
| `--admin` | `false` | Make site admin |
| `--data` | `./data` | Data directory |

### `loomhub doctor`

Verify system integrity.

```bash
loomhub doctor [flags]
```

Checks:
- Hub database schema version and integrity
- Per-loom database integrity
- Object store consistency (files vs. index)
- Orphaned objects
- Disk space

### `loomhub gc`

Garbage collect unreferenced objects.

```bash
loomhub gc [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--dry-run` | `false` | Show what would be deleted |
| `--data` | `./data` | Data directory |

### `loomhub info`

Show server information.

```bash
loomhub info [flags]
```

Output:
```
LoomHub v0.1.0
Data directory: /var/lib/loomhub
Hub database: 2.1 MB
Object store: 1.3 GB (4,523 objects)
Looms: 42
Users: 15
```

---

## Loom CLI Integration

Users interact with LoomHub through Loom's existing CLI commands (see [Loom CLI reference](../../loom/docs/11-cli-reference.md) and [sync protocol](../../loom/docs/06-systems/sync.md)):

```bash
# Add LoomHub as a hub
loom hub add origin https://loomhub.dev/flakerimi/my-app

# Authenticate — uses Loom's existing hub auth command
# Prompts for token or opens browser for OAuth
loom hub auth origin

# Send to LoomHub
loom send origin

# Receive from LoomHub
loom receive origin

# Show sync status
loom hub status
```

> **Note:** LoomHub does not require any new `loom` CLI commands. All interaction uses `loom hub add`, `loom hub auth`, `loom send`, and `loom receive` which are part of Loom's current command surface. The Loom client uses the hub URL as the base path and appends `/api/v1/negotiate`, `/api/v1/push`, `/api/v1/pull` — this works naturally with LoomHub's `/{owner}/{loom}/api/v1/...` routing. Future commands like `loom replicate` would need to be added to Loom itself first.

### Access Token Setup

Users create an access token in LoomHub's web UI (Settings > Access Tokens), then provide it when running `loom hub auth origin`. Loom stores credentials as documented in its own config format.
