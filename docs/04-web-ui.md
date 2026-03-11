# LoomHub — Web UI

Server-rendered HTML with htmx for progressive enhancement. No SPA, no heavy JavaScript framework.

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Templates | [templ](https://templ.guide) — type-safe Go templates |
| Interactivity | [htmx](https://htmx.org) — HTML-driven AJAX |
| CSS | Tailwind CSS (compiled at build time) |
| Icons | Heroicons (SVG) |
| Syntax highlighting | server-side via chroma |
| Markdown rendering | goldmark |

---

## Page Map

### Public Pages

| Route | Page | Description |
|-------|------|-------------|
| `/` | Landing / Dashboard | Home page (logged out: landing, logged in: dashboard) |
| `/explore` | Explore | Browse public repos, trending, recent |
| `/login` | Login | Email/password login form |
| `/register` | Register | Account creation |
| `/{owner}` | Profile | User or org profile + repo list |
| `/{owner}/{repo}` | Repository | Repo overview — README, recent checkpoints, stats |

### Repository Pages

All under `/{owner}/{repo}/`:

| Route | Page | Description |
|-------|------|-------------|
| `/` | Overview | Default stream entity tree + README + stats |
| `/tree/{stream}` | Entity Tree | Browse entity tree at stream head |
| `/tree/{stream}/{space}/{path...}` | Entity View | View entity content (syntax highlight, render md, etc.) |
| `/log` | Checkpoint Log | Paginated checkpoint history |
| `/log?stream={name}` | Stream Log | Checkpoints for specific stream |
| `/checkpoints/{id}` | Checkpoint Detail | Full checkpoint info + diff + entity states |
| `/operations` | Operation Log | Raw operation history (filterable by space/entity) |
| `/diff/{from}...{to}` | Diff View | Side-by-side or unified diff between two refs |
| `/streams` | Streams | List all streams with status, head, fork point |
| `/merge-requests` | Merge Requests | List MRs (open/closed/merged tabs) |
| `/merge-requests/new` | New MR | Create merge request form |
| `/merge-requests/{number}` | MR Detail | Discussion, diff, reviews, merge button |
| `/settings` | Repo Settings | Name, description, visibility, danger zone |
| `/settings/collaborators` | Collaborators | Manage repo access |
| `/settings/webhooks` | Webhooks | Webhook management |

### User Pages

| Route | Page | Description |
|-------|------|-------------|
| `/settings/profile` | Profile Settings | Edit profile info |
| `/settings/tokens` | Access Tokens | Manage API tokens |
| `/settings/keys` | SSH Keys | Manage SSH keys |

### Organization Pages

| Route | Page | Description |
|-------|------|-------------|
| `/orgs/new` | Create Org | Organization creation form |
| `/orgs/{name}/settings` | Org Settings | Edit org info |
| `/orgs/{name}/members` | Members | Manage org members |

### Admin Pages

| Route | Page | Description |
|-------|------|-------------|
| `/admin` | Admin Dashboard | System stats, health |
| `/admin/users` | User Management | List, create, disable users |
| `/admin/repos` | Repo Management | List, delete repos |

---

## Key UI Components

### Entity Tree Browser

Similar to GitHub's file browser but multi-space aware:

```
┌─────────────────────────────────────────────┐
│  flakerimi / my-app                          │
│  Stream: main ▾   Checkpoint: #45 (latest)   │
├─────────────────────────────────────────────┤
│  Spaces: [code] [docs] [design] [config]     │
├─────────────────────────────────────────────┤
│  📁 src/                                     │
│  📁 test/                                    │
│  📄 main.go              2.1 KB   3h ago     │
│  📄 go.mod               836 B    1d ago     │
│  📄 README.md            1.2 KB   2d ago     │
└─────────────────────────────────────────────┘
```

- Space tabs filter the view to a single content space
- Clicking a directory navigates into it (htmx partial swap)
- Clicking a file shows content with appropriate rendering
- Stream selector switches the view to different timelines

### Checkpoint Log

```
┌─────────────────────────────────────────────┐
│  Checkpoint History — main                    │
│  Stream: main ▾   Search: [____________]     │
├─────────────────────────────────────────────┤
│  ● #45  Add authentication         3h ago    │
│         flakerimi · 5 files · code, docs     │
│                                              │
│  ● #44  Fix database connection     1d ago    │
│         flakerimi · 2 files · code           │
│                                              │
│  ● #43  Initial project setup       2d ago    │
│         flakerimi · 12 files · code, config  │
│                                              │
│  [Load more...]                              │
└─────────────────────────────────────────────┘
```

### Multi-Space Diff View

Diffs are organized by space, then by entity:

```
┌─────────────────────────────────────────────┐
│  Diff: #43 → #45                             │
│  +134 operations · 6 entities · 2 spaces     │
├─────────────────────────────────────────────┤
│  ▼ code (5 files changed)                    │
│    ┌────────────────────────────────────┐    │
│    │ src/auth/login.go (+45, -3)        │    │
│    ├────────────────────────────────────┤    │
│    │ 10  func Login(...) {              │    │
│    │ 11-   return nil                   │    │
│    │ 11+   token, err := generateJWT()  │    │
│    │ 12+   if err != nil {              │    │
│    │ 13+     return err                 │    │
│    │ 14+   }                            │    │
│    └────────────────────────────────────┘    │
│                                              │
│  ▼ docs (1 file changed)                     │
│    ┌────────────────────────────────────┐    │
│    │ api-reference.md (+12, -0)         │    │
│    │ (rendered markdown diff)           │    │
│    └────────────────────────────────────┘    │
└─────────────────────────────────────────────┘
```

### Merge Request Page

```
┌─────────────────────────────────────────────┐
│  MR #42: Add user authentication             │
│  feature/auth → main                         │
│  Status: Open   Author: flakerimi            │
│  [Merge] [Close]                             │
├─────────────────────────────────────────────┤
│  [Conversation] [Diff (8 files)] [Reviews]   │
├─────────────────────────────────────────────┤
│  flakerimi opened this merge request 3h ago  │
│                                              │
│  Implements JWT-based authentication with... │
│                                              │
│  ─────────────────────────────────────────── │
│                                              │
│  jane: Looks good! Consider adding rate...   │
│                                              │
│  ─────────────────────────────────────────── │
│                                              │
│  [Write a comment...]              [Submit]  │
└─────────────────────────────────────────────┘
```

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

## htmx Patterns

### Partial Page Updates

Entity tree navigation uses htmx to swap only the content area:

```html
<a hx-get="/{owner}/{repo}/tree/{stream}/{space}/{path}"
   hx-target="#content"
   hx-push-url="true">
  src/auth/
</a>
```

### Infinite Scroll

Checkpoint log uses infinite scroll:

```html
<div hx-get="/{owner}/{repo}/log?after={cursor}"
     hx-trigger="revealed"
     hx-swap="afterend">
  Loading...
</div>
```

### Live Updates

MR comments update in real-time via polling:

```html
<div hx-get="/{owner}/{repo}/merge-requests/{number}/comments"
     hx-trigger="every 10s"
     hx-swap="innerHTML">
  <!-- comments rendered here -->
</div>
```

### Form Submission

Create MR without full page reload:

```html
<form hx-post="/api/v1/repos/{owner}/{repo}/merge-requests"
      hx-target="#main"
      hx-push-url="true">
  ...
</form>
```

---

## Space-Specific Rendering

Each space renders entity content differently:

| Space | Rendering |
|-------|-----------|
| **code** | Syntax-highlighted source (chroma), line numbers, copy button |
| **docs** | Rendered markdown (goldmark), table of contents |
| **design** | JSON tree viewer, visual component preview (future) |
| **data** | Schema table viewer, SQL/YAML highlighting |
| **config** | TOML/YAML/JSON highlighting with key-value structure |
| **notes** | Rendered markdown, plain text |

---

## Responsive Design

- Desktop: Full layout with sidebar navigation
- Tablet: Collapsed sidebar, full content area
- Mobile: Single-column layout, hamburger menu

All pages work without JavaScript (htmx enhances but doesn't require JS).
