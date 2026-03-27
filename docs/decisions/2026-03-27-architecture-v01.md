# Decision Record: v0.1 Architecture

- Date: 2026-03-27
- Status: Accepted
- Owners: @srmdn

## Context

First architectural decisions for the Inkwell CMS project. Needed to settle
structure before scaffolding begins.

## Decisions

### 1. Theme in separate repo

The default Astro theme lives in its own repository, not inside the core Go
repo. An `install.sh` script handles cloning and wiring at install time.

**Rationale:** Keeps the core repo focused on the Go binary. Allows themes
to be versioned and released independently. Community theme authors don't
need to touch core. Mirrors how Ghost and WordPress handle themes.

### 2. Migration strategy: version table

SQL migrations are numbered files with a version tracking table in SQLite.

**Rationale:** Required for safe upgrades when other people are running their
own instances. Additive-only works for a single deployment but breaks down
across multiple upgrade paths.

### 3. Configurable content path, default `content/blog/`

Content directory is configurable via `.env` but defaults to `content/blog/`
at the project root.

**Rationale:** Sensible default that matches the theme contract. Configurable
for users with non-standard layouts.

### 4. Admin setup via CLI flag `--setup`

First-run wizard triggered by `./inkwell --setup`. Creates the admin account
interactively.

**Rationale:** Better UX for open source installs than requiring manual env
var setup. Standard pattern for self-hosted tools.

### 5. Rebuild mechanism: subprocess for v0.1

Backend triggers theme build by calling `npm run build` as a subprocess and
restarting the theme service.

**Rationale:** Simplest path for v0.1. A webhook-based approach is better
long-term and can replace this in a future version without breaking the API.

## Consequences

- Theme repo must be created separately before a full end-to-end install is possible
- Migration files must be committed and never deleted
- `--setup` flag must be documented clearly in the README

## Follow-up Actions

- [x] Create `inkwell-theme-default` repo (separate, Astro SSR)
- [x] Write migration runner in `internal/db/`
- [ ] Document `--setup` flow in README
