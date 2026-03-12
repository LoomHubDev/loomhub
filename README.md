# LoomHub

A GitHub-like hosting platform for [Loom](https://github.com/LoomHubDev/loom) projects.

LoomHub lets you host, browse, and collaborate on Loom repositories. Push and pull streams, manage collaborators, browse checkpoint logs, compare streams, and manage your projects through a web interface.

## Features

- **Repository hosting** — Create and manage Loom repositories with public/private visibility
- **Sync protocol** — Push/pull operations and objects between the Loom CLI and LoomHub
- **Collaboration** — Invite collaborators with read/write/admin permissions
- **Checkpoint log** — Browse, filter, and paginate operation history
- **Stream comparison** — Visual diff between any two streams
- **Organizations** — Group repos and members under org namespaces
- **Work requests** — Issue-like work tracking with labels
- **Activity feed** — Track pushes, repo changes, and team activity
- **Admin dashboard** — Site-wide stats, user and loom management

## Tech Stack

- **Backend**: Go + [chi](https://github.com/go-chi/chi) router + SQLite ([modernc](https://pkg.go.dev/modernc.org/sqlite))
- **Frontend**: Vue 3 + TypeScript + Vite + Tailwind CSS 4 + Pinia
- **Auth**: JWT tokens with bcrypt password hashing

## Quick Start

### Backend

```bash
cd backend
go run ./cmd/loomhub
```

The server starts on `http://localhost:8080` by default.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

The dev server starts on `http://localhost:5173` and proxies API calls to the backend.

## Connecting the Loom CLI

```bash
# In your Loom project
loom hub add origin http://localhost:8080/username/reponame
loom hub auth origin
loom send
```

## API

All API endpoints are under `/api/v1/`:

- `POST /api/v1/auth/register` — Register
- `POST /api/v1/auth/login` — Login (returns JWT)
- `GET /api/v1/looms` — List looms
- `POST /api/v1/looms` — Create loom
- `GET/PATCH/DELETE /api/v1/looms/{owner}/{loom}` — Manage a loom
- `GET/POST/DELETE /api/v1/looms/{owner}/{loom}/collaborators` — Manage collaborators
- `GET /api/v1/looms/{owner}/{loom}/streams` — List streams
- `GET /api/v1/looms/{owner}/{loom}/checkpoints` — List checkpoints (filterable)
- `GET /api/v1/looms/{owner}/{loom}/entities` — List entities

Sync endpoints are per-loom at `/{owner}/{loom}/api/v1/`:

- `POST /negotiate` — Compare client/server stream state
- `POST /push` — Send operations and objects
- `POST /pull` — Receive operations and objects

## License

MIT
