# LoomHub — Web UI

Vue 3 SPA with Tailwind CSS. Communicates with the Go backend via REST API.

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Framework | Vue 3 (Composition API + `<script setup>`) |
| Build | Vite |
| Routing | Vue Router (history mode) |
| State | Pinia |
| CSS | Tailwind CSS |
| HTTP | ofetch |
| Syntax highlighting | Shiki (client-side, lazy-loaded per language) |
| Markdown | markdown-it |
| Icons | Heroicons (Vue components) |

---

## SPA Architecture

```
frontend/
├── src/
│   ├── main.ts                     # App entry point
│   ├── App.vue                     # Root component
│   ├── router/
│   │   └── index.ts                # Route definitions
│   ├── stores/
│   │   ├── auth.ts                 # Auth state (user, token)
│   │   ├── repo.ts                 # Active repo state
│   │   └── notification.ts         # Toast notifications
│   ├── api/
│   │   ├── client.ts               # ofetch instance with auth interceptor
│   │   ├── users.ts                # User API calls
│   │   ├── repos.ts                # Repo API calls
│   │   ├── sync.ts                 # Content API (log, entities, diff)
│   │   ├── mergeRequests.ts        # MR API calls
│   │   └── search.ts               # Search API calls
│   ├── layouts/
│   │   ├── DefaultLayout.vue       # Nav + footer wrapper
│   │   └── AuthLayout.vue          # Minimal layout for login/register
│   ├── views/
│   │   ├── Home.vue                # Landing / dashboard
│   │   ├── Explore.vue             # Browse public repos
│   │   ├── Login.vue
│   │   ├── Register.vue
│   │   ├── Profile.vue             # User/org profile
│   │   ├── repo/
│   │   │   ├── RepoLayout.vue      # Repo header + tabs wrapper
│   │   │   ├── RepoOverview.vue    # README + stats
│   │   │   ├── EntityTree.vue      # File/entity browser
│   │   │   ├── EntityView.vue      # Single entity content
│   │   │   ├── CheckpointLog.vue   # Checkpoint history
│   │   │   ├── CheckpointDetail.vue
│   │   │   ├── DiffView.vue        # Multi-space diff
│   │   │   ├── Streams.vue         # Stream list
│   │   │   ├── OperationLog.vue    # Raw operation history
│   │   │   └── Settings.vue        # Repo settings
│   │   ├── mr/
│   │   │   ├── MRList.vue          # Merge request list
│   │   │   ├── MRNew.vue           # Create MR form
│   │   │   └── MRDetail.vue        # MR discussion + diff + reviews
│   │   ├── settings/
│   │   │   ├── ProfileSettings.vue
│   │   │   ├── TokenSettings.vue
│   │   │   └── KeySettings.vue
│   │   └── admin/
│   │       ├── AdminDashboard.vue
│   │       ├── AdminUsers.vue
│   │       └── AdminRepos.vue
│   ├── components/
│   │   ├── nav/
│   │   │   ├── AppNav.vue          # Top navigation bar
│   │   │   ├── RepoNav.vue         # Repo sub-navigation tabs
│   │   │   └── UserMenu.vue        # Avatar + dropdown
│   │   ├── repo/
│   │   │   ├── EntityTreeItem.vue  # Recursive tree node
│   │   │   ├── SpaceTabs.vue       # code | docs | design | data | config | notes
│   │   │   ├── StreamSelector.vue  # Stream dropdown
│   │   │   ├── CheckpointCard.vue  # Checkpoint in log list
│   │   │   └── RepoStats.vue       # Star/fork/checkpoint counts
│   │   ├── diff/
│   │   │   ├── DiffFile.vue        # Single entity diff (collapsible)
│   │   │   ├── DiffHunk.vue        # Hunk with line numbers
│   │   │   ├── DiffLine.vue        # Add/delete/context line
│   │   │   ├── DiffStats.vue       # +/- summary bar
│   │   │   └── InlineComment.vue   # Comment anchored to diff line
│   │   ├── mr/
│   │   │   ├── MRStatus.vue        # Open/merged/closed badge
│   │   │   ├── ReviewBadge.vue     # Approved/changes requested
│   │   │   ├── CommentThread.vue   # Comment + replies
│   │   │   └── CommentEditor.vue   # Markdown editor with preview
│   │   ├── content/
│   │   │   ├── CodeViewer.vue      # Syntax highlighted code (Shiki)
│   │   │   ├── MarkdownViewer.vue  # Rendered markdown
│   │   │   ├── JsonViewer.vue      # Collapsible JSON tree
│   │   │   └── RawViewer.vue       # Plain text fallback
│   │   └── ui/
│   │       ├── Avatar.vue
│   │       ├── Badge.vue
│   │       ├── Button.vue
│   │       ├── Dropdown.vue
│   │       ├── Modal.vue
│   │       ├── Pagination.vue
│   │       ├── Spinner.vue
│   │       ├── Toast.vue
│   │       ├── Tabs.vue
│   │       └── EmptyState.vue
│   └── utils/
│       ├── format.ts               # Date, file size formatting
│       ├── space.ts                # Space-specific helpers
│       └── diff.ts                 # Diff parsing helpers
├── index.html
├── vite.config.ts
├── tailwind.config.ts
├── tsconfig.json
└── package.json
```

---

## Route Map

### Public Routes

| Route | View | Description |
|-------|------|-------------|
| `/` | Home | Landing (guest) or dashboard (logged in) |
| `/explore` | Explore | Browse public repos |
| `/login` | Login | Login form |
| `/register` | Register | Registration form |
| `/:owner` | Profile | User or org profile + repo list |
| `/:owner/:repo` | RepoOverview | Repo overview (README, stats) |

### Repository Routes

All under `/:owner/:repo/`:

| Route | View | Description |
|-------|------|-------------|
| `/` | RepoOverview | Default stream entity tree + README |
| `/tree/:stream` | EntityTree | Browse entities at stream head |
| `/tree/:stream/:space/*path` | EntityView | View single entity |
| `/log` | CheckpointLog | Checkpoint history |
| `/checkpoints/:id` | CheckpointDetail | Full checkpoint details + diff |
| `/operations` | OperationLog | Raw operation log |
| `/diff/:from...:to` | DiffView | Multi-space diff between refs |
| `/streams` | Streams | List all streams |
| `/merge-requests` | MRList | List MRs (open/closed/merged) |
| `/merge-requests/new` | MRNew | Create MR form |
| `/merge-requests/:number` | MRDetail | MR detail page |
| `/settings` | Settings | Repo settings (admin only) |

### User Settings Routes

| Route | View | Description |
|-------|------|-------------|
| `/settings/profile` | ProfileSettings | Edit profile |
| `/settings/tokens` | TokenSettings | Manage access tokens |
| `/settings/keys` | KeySettings | Manage SSH keys |

### Admin Routes

| Route | View | Description |
|-------|------|-------------|
| `/admin` | AdminDashboard | System stats |
| `/admin/users` | AdminUsers | User management |
| `/admin/repos` | AdminRepos | Repo management |

---

## Key UI Components

### Entity Tree Browser

Multi-space aware file browser with reactive updates:

```
┌─────────────────────────────────────────────┐
│  flakerimi / my-app                          │
│  Stream: [main ▾]   Checkpoint: #45 (latest) │
├─────────────────────────────────────────────┤
│  [code] [docs] [design] [config]  ← SpaceTabs│
├─────────────────────────────────────────────┤
│  📁 src/                                     │
│  📁 test/                                    │
│  📄 main.go              2.1 KB   3h ago     │
│  📄 go.mod               836 B    1d ago     │
│  📄 README.md            1.2 KB   2d ago     │
└─────────────────────────────────────────────┘
```

- `SpaceTabs` — reactive tab switching filters entities by space
- `StreamSelector` — dropdown switches stream, triggers API refetch
- `EntityTreeItem` — recursive component, click to navigate
- Breadcrumb navigation for deep paths
- Vue Router handles navigation, Pinia caches entity data

### Multi-Space Diff View

```
┌─────────────────────────────────────────────┐
│  Diff: #43 → #45                             │
│  +134 operations · 6 entities · 2 spaces     │
│  [Unified] [Side-by-Side]  ← toggle          │
├─────────────────────────────────────────────┤
│  ▼ code (5 files changed)       ← DiffFile   │
│    ┌────────────────────────────────────┐    │
│    │ src/auth/login.go (+45, -3)        │    │
│    ├────────────────────────────────────┤    │
│    │ 10  func Login(...) {    ← DiffHunk│    │
│    │ 11- return nil           ← DiffLine│    │
│    │ 11+ token, err := ...              │    │
│    │    💬 [Add comment]  ← InlineComment│    │
│    └────────────────────────────────────┘    │
│                                              │
│  ▼ docs (1 file changed)                     │
│    ┌────────────────────────────────────┐    │
│    │ api-reference.md (+12, -0)         │    │
│    │ (rendered markdown diff)           │    │
│    └────────────────────────────────────┘    │
└─────────────────────────────────────────────┘
```

Key features:
- Unified/side-by-side toggle (reactive, no re-fetch)
- Collapsible file sections
- Expand context around hunks (fetch more lines on click)
- Inline comments on diff lines (for MR review)
- Space-grouped diffs with per-space stats

### Merge Request Detail

```
┌─────────────────────────────────────────────┐
│  MR #42: Add user authentication             │
│  feature/auth → main                         │
│  [Open] Author: flakerimi   ✓ Approved       │
│  [Merge ▾] [Close]                           │
├─────────────────────────────────────────────┤
│  [Conversation] [Diff (8)] [Reviews (1)]     │
├─────────────────────────────────────────────┤
│  ┌─ Conversation tab ─────────────────────┐ │
│  │  flakerimi opened 3h ago               │ │
│  │  Implements JWT-based authentication... │ │
│  │                                        │ │
│  │  jane reviewed 1h ago — Approved ✓     │ │
│  │  "LGTM! Ship it."                      │ │
│  │                                        │ │
│  │  ┌─ Comment Editor ─────────────────┐  │ │
│  │  │ [Write] [Preview]                │  │ │
│  │  │ [                           ]    │  │ │
│  │  │               [Comment] [Submit] │  │ │
│  │  └─────────────────────────────────┘  │ │
│  └────────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
```

Features:
- Tab switching without page reload (Vue Router nested routes)
- Real-time comment updates (polling or SSE)
- Markdown preview in comment editor
- Review submission with approve/request changes
- Merge button with conflict detection

### Dashboard (Logged In)

```
┌─────────────────────────────────────────────┐
│  LoomHub                    [+ New Repo] 🔔  │
├──────────────┬──────────────────────────────┤
│  Your Repos  │  Activity Feed               │
│              │                              │
│  my-app      │  flakerimi pushed to main    │
│  ★ 3  ⑂ 1   │  my-app · 3h ago · 5 files   │
│              │                              │
│  website     │  jane opened MR #42          │
│  ★ 0  ⑂ 0   │  my-app · 5h ago              │
│              │                              │
│  [View all]  │  flakerimi created repo      │
│              │  website · 1d ago             │
└──────────────┴──────────────────────────────┘
```

---

## Space-Specific Rendering

Each space uses a different content viewer component:

| Space | Component | Rendering |
|-------|-----------|-----------|
| **code** | `CodeViewer` | Shiki syntax highlighting, line numbers, copy button |
| **docs** | `MarkdownViewer` | Rendered markdown with TOC |
| **design** | `JsonViewer` | Collapsible JSON tree |
| **data** | `CodeViewer` | SQL/YAML/JSON highlighting |
| **config** | `CodeViewer` | TOML/YAML/JSON highlighting |
| **notes** | `MarkdownViewer` | Rendered markdown / plain text |

The `EntityView` component dispatches to the correct viewer based on space type.

---

## State Management (Pinia)

### Auth Store

```ts
// stores/auth.ts
export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))

  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.is_admin ?? false)

  async function login(email: string, password: string) { ... }
  async function logout() { ... }
  async function fetchUser() { ... }

  return { user, token, isAuthenticated, isAdmin, login, logout, fetchUser }
})
```

### Repo Store

```ts
// stores/repo.ts
export const useRepoStore = defineStore('repo', () => {
  const repo = ref<Repository | null>(null)
  const activeStream = ref('main')
  const streams = ref<Stream[]>([])

  async function fetchRepo(owner: string, name: string) { ... }
  async function fetchStreams() { ... }

  return { repo, activeStream, streams, fetchRepo, fetchStreams }
})
```

---

## API Client

```ts
// api/client.ts
import { ofetch } from 'ofetch'

export const api = ofetch.create({
  baseURL: '/api/v1',
  onRequest({ options }) {
    const token = localStorage.getItem('token')
    if (token) {
      options.headers = { ...options.headers, Authorization: `Bearer ${token}` }
    }
  },
  onResponseError({ response }) {
    if (response.status === 401) {
      useAuthStore().logout()
    }
  }
})
```

---

## Development

```bash
cd frontend
npm install
npm run dev          # Vite dev server with HMR on :5173
                     # Proxies /api/* to Go backend on :3000
```

### Vite Proxy Config

```ts
// vite.config.ts
export default defineConfig({
  server: {
    proxy: {
      '/api': 'http://localhost:3000'
    }
  }
})
```

### Production Build

```bash
npm run build        # Outputs to frontend/dist/
```

The Go binary embeds `frontend/dist/` via `embed.FS` and serves it at `/`.

---

## Responsive Design

- Desktop: Full layout with sidebar (dashboard), wide diff views
- Tablet: Collapsed sidebar, stacked layout
- Mobile: Single-column, hamburger menu, simplified diffs

Tailwind breakpoints: `sm:` (640px), `md:` (768px), `lg:` (1024px), `xl:` (1280px).
