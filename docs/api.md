# API Reference

Base URL: wherever Folio is running (e.g. `http://localhost:8090`)

All responses are JSON. All request bodies must be `Content-Type: application/json`.

---

## Authentication

Protected endpoints require a valid JWT. The token is issued on login and
stored in an `HttpOnly` cookie named `token`. API clients may also pass it
as a `Authorization: Bearer <token>` header.

Protected mutation endpoints (POST/PUT/DELETE) additionally require a
`X-CSRF-Token` header. The CSRF token is returned in the login response and
can be refreshed via `GET /api/csrf-token`.

---

## Public Endpoints

### `GET /health`

Returns `200 OK` with body `ok`. Use for uptime checks.

---

### `POST /api/login`

Authenticate and receive a session token.

**Request body**
```json
{
  "email": "admin@example.com",
  "password": "yourpassword"
}
```

**Response `200`**
```json
{
  "csrf_token": "abc123..."
}
```

Sets an `HttpOnly` cookie named `token` (24-hour TTL).
Store the `csrf_token` and include it as `X-CSRF-Token` on all subsequent
mutation requests.

**Errors**: `400` invalid body, `401` invalid credentials

---

### `POST /api/logout`

Clear the session cookie.

**Response `204`** No content.

---

### `GET /api/posts`

List all published posts. Does not include post body.

**Response `200`**
```json
[
  {
    "id": 1,
    "slug": "my-first-post",
    "title": "My First Post",
    "description": "A short description.",
    "tags": "go,cms",
    "draft": false,
    "publish_date": "2026-03-27T00:00:00Z",
    "created_at": "2026-03-27T10:00:00Z",
    "updated_at": "2026-03-27T10:00:00Z"
  }
]
```

Returns `[]` (empty array) if no posts are published.

---

### `GET /api/posts/{slug}`

Get a single published post including its Markdown body.

**Response `200`**
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
  "body": "Post body in **Markdown**."
}
```

**Errors**: `404` post not found or is a draft

---

## Authenticated Endpoints

Requires: valid `token` cookie or `Authorization: Bearer` header.

---

### `GET /api/csrf-token`

Get the CSRF token for the current session. Use this after a page refresh
when the token stored in memory is lost.

**Response `200`**
```json
{
  "csrf_token": "abc123..."
}
```

---

## Authenticated + CSRF Endpoints

Requires: valid `token` cookie/header **and** `X-CSRF-Token` header.

---

### `GET /api/admin/posts`

List all posts including drafts.

**Response `200`** — same shape as `GET /api/posts` but includes drafts.

---

### `GET /api/admin/posts/{slug}`

Get a single post regardless of draft status, including body.

**Response `200`** — same shape as `GET /api/posts/{slug}`.

**Errors**: `404` post not found

---

### `POST /api/admin/posts/{slug}`

Create a new post. The slug is set in the URL path.

**Slug rules**: lowercase letters, numbers, and hyphens only. Example: `my-first-post`

**Request body**
```json
{
  "title": "My First Post",
  "description": "A short description.",
  "tags": ["go", "cms"],
  "draft": true,
  "publish_date": "2026-03-27",
  "body": "Post body in **Markdown**."
}
```

| Field | Required | Notes |
|-------|----------|-------|
| `title` | Yes | |
| `description` | No | |
| `tags` | No | Array of strings |
| `draft` | No | Defaults to `false` |
| `publish_date` | No | Format `YYYY-MM-DD`. Defaults to today. |
| `body` | No | Markdown string |

**Response `201`** No content.

**Errors**: `400` invalid slug or missing title, `409` post already exists

---

### `PUT /api/admin/posts/{slug}`

Update an existing post. Replaces all fields.

**Request body** — same shape as `POST /api/admin/posts/{slug}`

**Response `204`** No content.

**Errors**: `400` missing title, `404` post not found

---

### `DELETE /api/admin/posts/{slug}`

Delete a post and its content files from disk.

**Response `204`** No content.

**Errors**: `404` post not found

---

### `POST /api/admin/rebuild`

Trigger an async theme rebuild. Returns immediately — poll status to track
progress.

**Response `202`** Build started.

**Response `409`** A rebuild is already in progress.

---

### `GET /api/admin/rebuild/status`

Get the current rebuild status.

**Response `200`**
```json
{
  "status": "success",
  "output": "...",
  "started_at": "2026-03-27T10:00:00Z",
  "finished_at": "2026-03-27T10:01:05Z",
  "error": ""
}
```

| `status` value | Meaning |
|----------------|---------|
| `idle` | No build has run yet |
| `running` | Build in progress |
| `success` | Last build succeeded |
| `failed` | Last build failed — see `error` field |
