# Theme Contract

This document defines the interface between the Folio backend and any
frontend theme. A theme that follows this contract is compatible with Folio
out of the box.

You do not need to read or modify the Go backend to build a theme.

---

## What a Theme Is

A theme is any frontend that:

1. Reads post content from the filesystem (content contract)
2. Optionally calls the Folio REST API for dynamic data
3. Exposes a build command the backend can call
4. Optionally runs as a service the backend can restart after a build

Folio ships with one default theme (`foliocms-theme-default`, Astro SSR).
You can replace it with any framework (Astro, Next.js, SvelteKit, plain HTML,
or anything else) as long as it follows this contract.

---

## 1. Content Contract

Posts are stored as Markdown files on disk. The backend writes them; the
theme reads them.

### Directory Structure

```
<CONTENT_DIR>/
  <slug>/
    index.md        ← post content + frontmatter
    hero.webp       ← optional hero/OG image
```

`CONTENT_DIR` defaults to `content/blog` relative to the Folio binary.
It is configurable via the `CONTENT_DIR` environment variable.

### Frontmatter Spec

Every `index.md` starts with a YAML frontmatter block:

```yaml
---
title: "Post Title"
description: "A short description of the post."
publishDate: "2026-03-27"
draft: false
tags:
  - go
  - cms
---

Post body in Markdown follows here.
```

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `title` | string | Yes | Quote if the value contains a colon |
| `description` | string | No | Used for meta description and previews |
| `publishDate` | string | Yes | Format: `YYYY-MM-DD` |
| `draft` | boolean | Yes | `true` = not published, `false` = published |
| `tags` | string list | No | Used for categorization |

**Important**: any frontmatter string value that contains a colon (`:`) must
be quoted, or YAML parsing will fail.

### Theme Responsibility

- Read `index.md` files at build time (static) or request time (SSR)
- Respect `draft: true`: do not render draft posts on public pages
- Respect `publishDate`: do not render posts with a future publish date
- The hero image (`hero.webp`) is optional; handle its absence gracefully

---

## 2. API Contract

The Folio backend exposes a REST API. Themes may use the public endpoints
for dynamic data (view counts, etc.) or skip them entirely if building static.

### Base URL

Configured by the theme. Typically the same origin or a configured API URL.

### Public Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/posts` | List all published posts (metadata only, no body) |
| `GET` | `/api/posts/{slug}` | Single published post with body |
| `POST` | `/api/posts/{slug}/view` | Increment view count (optional, fire-and-forget) |

### Response: Post Object

```json
{
  "id": 1,
  "slug": "my-first-post",
  "title": "My First Post",
  "description": "A short description.",
  "tags": "go,cms",
  "draft": false,
  "publish_date": "2026-03-27T00:00:00Z",
  "created_at": "2026-03-27T10:00:00Z",
  "updated_at": "2026-03-27T10:00:00Z",
  "body": "Post body in Markdown..."
}
```

`body` is only included in the single-post response (`GET /api/posts/{slug}`).
`tags` is a comma-separated string.

### Static Themes

Static themes do not need to call the API at all. They can read everything
from the filesystem at build time. The API is available for SSR themes that
need runtime data.

---

## 3. Build Contract

The Folio backend triggers a theme rebuild via a shell command. This happens
when the admin clicks "Rebuild Site" in the dashboard.

### Configuration

```env
THEME_DIR=theme
THEME_BUILD_CMD=npm run build
```

- `THEME_DIR`: directory where the build command is executed (relative to
  the Folio binary, or absolute)
- `THEME_BUILD_CMD`: the command to run. Executed as a subprocess with
  `THEME_DIR` as the working directory.

### What the Backend Does

1. Runs `THEME_BUILD_CMD` in `THEME_DIR`
2. Captures stdout and stderr
3. On success: optionally restarts `THEME_SERVICE` (see section 4)
4. On failure: stores the error output, reports `failed` status

### What the Theme Must Do

- Accept a build command (e.g. `npm run build`, `bun run build`)
- Exit with code `0` on success, non-zero on failure
- Write build output (HTML, JS, CSS) wherever the web server expects it

### Build Status Polling

The dashboard polls `GET /api/admin/rebuild/status` after triggering a build.
Response:

```json
{
  "status": "running",
  "output": "...",
  "started_at": "2026-03-27T10:00:00Z",
  "finished_at": null,
  "error": ""
}
```

Status values: `idle` | `running` | `success` | `failed`

---

## 4. Service Contract (Optional)

For SSR themes running as a persistent process (e.g. Astro SSR, Next.js),
the backend can restart the service after a successful build.

```env
THEME_SERVICE=my-theme-service
```

If `THEME_SERVICE` is set, the backend runs `systemctl restart <THEME_SERVICE>`
after a successful build. Leave it empty for static themes or non-systemd setups.

---

## 5. Summary: What Each Side Owns

| Responsibility | Backend | Theme |
|----------------|---------|-------|
| Writing post files | Yes | No |
| Reading post files | No | Yes |
| Enforcing draft/publish rules | API only | Yes (at render time) |
| Running the build | Triggers it | Implements it |
| Serving the frontend | No | Yes |
| Auth and admin dashboard | Yes | No |

---

## 6. Building a Custom Theme

Minimum requirements:

1. Read Markdown files from `CONTENT_DIR` (default: `content/blog/`)
2. Parse YAML frontmatter using the field spec above
3. Skip posts where `draft: true` or `publishDate` is in the future
4. Expose a build command (any tool: npm, bun, make, etc.)
5. Accept `THEME_DIR` and `THEME_BUILD_CMD` configuration in Folio's `.env`

Optional:
- Call `POST /api/posts/{slug}/view` on each page visit to track views
- Set `THEME_SERVICE` if running as an SSR service that needs restarting

See `foliocms-theme-default` for a reference implementation using Astro SSR.
