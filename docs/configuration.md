# Configuration

Folio is configured entirely through environment variables. Copy
`.env.example` to `.env` and fill in the required values before starting.

```bash
cp .env.example .env
```

---

## Required

| Variable | Description |
|----------|-------------|
| `JWT_SECRET` | Secret key for signing JWT tokens. Generate with: `openssl rand -hex 32` |

---

## Server

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8090` | Port the backend listens on |

---

## Database

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `data/folio.db` | Path to the SQLite database file. Relative to the binary or absolute. |

The database directory is created automatically if it does not exist.

---

## Content

| Variable | Default | Description |
|----------|---------|-------------|
| `CONTENT_DIR` | `content/blog` | Directory where post Markdown files are stored. Relative to the binary or absolute. |

---

## Admin Setup

These variables are only used during `--setup`. They are optional; if not
set, the setup wizard will prompt interactively.

| Variable | Default | Description |
|----------|---------|-------------|
| `ADMIN_EMAIL` | _(prompted)_ | Email address for the admin account |
| `ADMIN_PASSWORD` | _(prompted)_ | Password for the admin account (min 8 characters) |

Do not leave real credentials in `.env` after setup if the file could be
exposed. Consider unsetting them or using interactive setup instead.

---

## Theme

| Variable | Default | Description |
|----------|---------|-------------|
| `THEME_DIR` | `theme` | Directory where the theme lives. The build command is executed here. |
| `THEME_BUILD_CMD` | `npm run build` | Command to build the theme. Runs in `THEME_DIR`. |
| `THEME_SERVICE` | _(empty)_ | Optional. Systemd service name to restart after a successful build. Leave empty to skip. |

---

## Webhook

| Variable | Default | Description |
|----------|---------|-------------|
| `WEBHOOK_SECRET` | _(empty)_ | Optional. If set, enables `POST /api/webhook/rebuild`. Pass the secret via `X-Webhook-Secret` header or `Authorization: Bearer`. |

---

## Media Library

| Variable | Default | Description |
|----------|---------|-------------|
| `SITE_URL` | `http://localhost:8090` | Base URL of the site. Used to build absolute URLs for locally stored media files. |
| `MEDIA_STORAGE` | `local` | Storage backend. `local` stores files on disk; `s3` uses any S3-compatible provider. |

When `MEDIA_STORAGE=local`, uploaded files are stored in a `media/` directory next to `CONTENT_DIR`
and served at `GET /media/{key}`.

---

## S3-Compatible Storage

Required when `MEDIA_STORAGE=s3`. Supports AWS S3, Cloudflare R2, NevaObjects, MinIO, and any
S3-compatible provider.

| Variable | Default | Description |
|----------|---------|-------------|
| `S3_ENDPOINT` | _(required)_ | API endpoint of the S3 provider. |
| `S3_BUCKET` | _(required)_ | Bucket name. |
| `S3_REGION` | `auto` | Region. Use `auto` for providers that don't require one (R2, NevaObjects). |
| `S3_ACCESS_KEY` | _(required)_ | Access key ID. |
| `S3_SECRET_KEY` | _(required)_ | Secret access key. |
| `S3_PUBLIC_URL` | _(required)_ | Base URL prepended to file keys for public access. |

Provider-specific endpoint and public URL examples:

| Provider | `S3_ENDPOINT` | `S3_PUBLIC_URL` |
|----------|--------------|-----------------|
| AWS S3 | `https://s3.<region>.amazonaws.com` | `https://<bucket>.s3.<region>.amazonaws.com` |
| Cloudflare R2 | `https://<account_id>.r2.cloudflarestorage.com` | your custom domain or R2 public URL |
| NevaObjects | `https://s3.nevaobjects.id` | `https://s3.nevaobjects.id/<bucket>` |
| MinIO | `https://minio.example.com` | `https://minio.example.com/<bucket>` |

---

## SMTP (Newsletter)

Required to send newsletters via `POST /api/admin/newsletter/send`.

| Variable | Default | Description |
|----------|---------|-------------|
| `SMTP_HOST` | _(empty)_ | SMTP server hostname, e.g. `smtp.mailgun.org` |
| `SMTP_PORT` | `587` | SMTP port |
| `SMTP_USERNAME` | _(empty)_ | SMTP username |
| `SMTP_PASSWORD` | _(empty)_ | SMTP password |
| `SMTP_FROM` | _(empty)_ | Sender address, e.g. `newsletter@example.com` |

If `SMTP_HOST` is not set, the newsletter send endpoint returns `503`.

---

## Example `.env`

```env
PORT=8090
DATABASE_URL=data/folio.db
CONTENT_DIR=content/blog

JWT_SECRET=your-generated-secret-here

THEME_DIR=theme
THEME_BUILD_CMD=npm run build
THEME_SERVICE=

# Media (local storage, default)
SITE_URL=https://example.com
MEDIA_STORAGE=local

# Media (S3-compatible storage, uncomment to use)
# MEDIA_STORAGE=s3
# S3_ENDPOINT=https://s3.nevaobjects.id
# S3_BUCKET=my-bucket
# S3_REGION=auto
# S3_ACCESS_KEY=
# S3_SECRET_KEY=
# S3_PUBLIC_URL=https://s3.nevaobjects.id/my-bucket

# Newsletter (optional)
# SMTP_HOST=smtp.mailgun.org
# SMTP_PORT=587
# SMTP_USERNAME=
# SMTP_PASSWORD=
# SMTP_FROM=newsletter@example.com
```

---

## Notes

- Never commit your `.env` file. It is gitignored by default.
- `.env.example` is committed and lists all variable names with empty or
  default values. Add new variables there when introducing them; do not
  remove existing ones, as that would break existing installs.
- The binary loads `.env` by default. Use `--env /path/to/.env` to specify
  a different file.
