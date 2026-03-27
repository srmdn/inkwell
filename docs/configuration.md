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

## Example `.env`

```env
PORT=8090
DATABASE_URL=data/folio.db
CONTENT_DIR=content/blog

JWT_SECRET=your-generated-secret-here

THEME_DIR=theme
THEME_BUILD_CMD=npm run build
THEME_SERVICE=
```

---

## Notes

- Never commit your `.env` file. It is gitignored by default.
- `.env.example` is committed and lists all variable names with empty or
  default values. Add new variables there when introducing them; do not
  remove existing ones, as that would break existing installs.
- The binary loads `.env` by default. Use `--env /path/to/.env` to specify
  a different file.
