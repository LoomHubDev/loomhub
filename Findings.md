# LoomHub Findings

Current documentation review findings for `/Users/flakerim/LoomProject/loomhub`.

## P1

### 1. Sync API is described as native Loom, but the contract differs

Affected doc:

- `docs/03-api-design.md`

Problem:

- the docs say LoomHub exposes Loom's native sync API
- the documented endpoints move sync under `/api/v1/sync/{owner}/{repo}`
- the request bodies omit `project_id`
- push/pull add `checkpoints` payloads that are not part of Loom's current documented sync protocol

Why it matters:

- a `loom` client built to the current Loom docs will not interoperate with LoomHub as documented
- implementation will drift unless LoomHub either:
  - matches Loom's current protocol exactly, or
  - explicitly defines a versioned translation layer

### 2. Storage compatibility is overstated

Affected doc:

- `docs/02-data-models.md`

Problem:

- the docs say repo-level storage is identical to local `.loom/` format
- the design also introduces:
  - a shared cross-repo object store
  - a global object reference table

Why it matters:

- Loom currently models object storage and ref counting at the repo level
- this means LoomHub is not actually using local Loom repo storage unchanged
- the docs need either:
  - a weaker compatibility claim, or
  - an explicit adaptation/storage layer description

## P2

### 3. Shared owner namespace is claimed but not enforced

Affected doc:

- `docs/02-data-models.md`

Problem:

- the docs claim usernames and org names share one namespace
- the schema keeps them in separate tables with separate unique constraints

Why it matters:

- both a user and an org could be named `acme`
- that makes `/{owner}` and `/{owner}/{repo}` routing ambiguous

Suggested direction:

- use a shared principals table, or
- add a reservation mechanism that enforces a global owner namespace

### 4. Repository ownership schema cannot support the documented cascade behavior

Affected doc:

- `docs/02-data-models.md`

Problem:

- `repositories.owner_id` is polymorphic (`user` or `org`)
- it is not a foreign key to either owner table
- the docs still claim cascading deletes from owner to repo

Why it matters:

- the database cannot enforce that invariant as written
- owner deletion cleanup must happen in application code unless the schema changes

### 5. Loom CLI integration docs drift from Loom's current command surface

Affected doc:

- `docs/09-cli-reference.md`

Problem:

- the doc references:
  - `loom auth login`
  - `loom clone`
  - `~/.loom/credentials`
- Loom's current docs describe:
  - `loom remote add`
  - `loom remote auth <name>`
- the other commands/storage format are not currently part of Loom's documented CLI

Why it matters:

- readers may assume these commands already exist
- LoomHub implementation may get ahead of Loom without an explicit compatibility plan

## Verdict

The LoomHub docs are good enough to express product direction, but not yet good enough to be treated as the final implementation source of truth.

Before building against them, fix:

1. the Loom protocol boundary
2. the storage compatibility claim
3. the owner namespace and ownership schema inconsistencies
4. the Loom CLI integration drift
