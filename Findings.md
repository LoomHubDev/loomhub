# LoomHub Findings

Current documentation review findings for `/Users/flakerim/LoomProject/loomhub`.

## P1

### 1. LoomHub remote URLs still do not line up with Loom's current sync contract

Affected docs:

- `docs/03-api-design.md`
- `docs/09-cli-reference.md`

Problem:

- Loom's current docs define sync endpoints at:
  - `POST /api/v1/negotiate`
  - `POST /api/v1/push`
  - `POST /api/v1/pull`
- Loom's current remote examples use URLs like:
  - `https://loom.example.com/project/my-app`
- LoomHub documents remotes like:
  - `https://hub.example.com/flakerimi/my-app`
- LoomHub's sync endpoints are documented as:
  - `/api/v1/sync/{owner}/{repo}/negotiate`
  - `/api/v1/sync/{owner}/{repo}/push`
  - `/api/v1/sync/{owner}/{repo}/pull`

Why it matters:

- the docs still claim LoomHub works with the existing `loom` CLI without changes
- but the current Loom client/docs contract does not explain how a remote URL of `/{owner}/{repo}` reaches sync endpoints under `/api/v1/sync/{owner}/{repo}/...`
- this needs one of:
  - Loom client changes
  - a documented discovery/translation layer
  - or LoomHub remote URLs that match Loom's current `/project/{id}` style

## P2

### 2. Repository hosting docs contradict the shared-store adaptation model

Affected docs:

- `docs/02-data-models.md`
- `docs/06-repository-hosting.md`

Problem:

- `docs/02-data-models.md` correctly says LoomHub uses:
  - a shared cross-repo object store
  - a global object reference table
  - a per-repo `objects` table that is only an index
- `docs/06-repository-hosting.md` still says each repo DB uses the exact same Loom schema and lists:
  - `Objects (index — hash, size, compressed, ref_count)`

Why it matters:

- these two docs describe different repo-level storage contracts
- `ref_count` cannot be both repo-local and globally owned by the hub
- implementation will drift unless `06-repository-hosting.md` matches the newer adaptation model

### 3. Permission rules reference an `internal` repo visibility that does not exist elsewhere

Affected docs:

- `docs/05-auth-permissions.md`
- `docs/02-data-models.md`

Problem:

- org member permissions say members get `write (internal)`
- repository visibility is only documented as `public` or `private`

Why it matters:

- access control rules are underspecified for org repos
- either `internal` needs to become a real third visibility level, or the permission examples need to use the existing two-level model

### 4. The `go:embed` example is invalid as written

Affected doc:

- `docs/07-project-structure.md`

Problem:

- the example uses:
  - `//go:embed all:../../frontend/dist`
- this pattern is not valid for Go embedding

Why it matters:

- the example cannot compile as shown
- I verified this with a temporary Go module: `go test` fails with `pattern all:../../frontend/dist: invalid pattern syntax`
- the docs should show an embeddable path layout or a build step that copies frontend assets into a valid location before embedding

### 5. Testing strategy still assumes a `loom clone` command that the CLI docs mark as future work

Affected docs:

- `docs/10-testing-strategy.md`
- `docs/09-cli-reference.md`

Problem:

- the end-to-end test example uses:
  - `testutil.RunLoom(t, dir2, "clone", srv.URL+"/testuser/my-app")`
- but the CLI reference explicitly says `loom clone` is a future command and not part of the current Loom surface

Why it matters:

- the test plan currently depends on a command the rest of the docs say does not exist yet
- for the current MVP, the example should use `init + remote add + pull` or clearly mark `clone` as future-only

## Verdict

The docs are much closer now. The schema and owner-model problems from the earlier review are fixed.

They are still not fully implementation-ready yet because the LoomHub ↔ Loom protocol boundary is unresolved, and a few smaller docs now disagree on storage, permissions, embedding, and test assumptions.
