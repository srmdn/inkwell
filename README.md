# Inkwell

A lightweight, self-hostable CMS. Single Go binary. SQLite. Ships with a default Astro theme. No Docker required.

## Why

Most open source CMS options are either too heavy (WordPress, Ghost) or too opinionated about the frontend (Payload, Strapi). Inkwell is built for developers who self-host on a cheap VPS and want:

- A real Markdown editor out of the box
- A default frontend theme they can swap
- Zero external dependencies — no PostgreSQL, no Redis, no Docker
- One systemd service and done

## Requirements

- Go 1.21+
- Node.js 18+ and npm
- Git
- Linux or macOS

> **Note:** The default theme uses Astro 5 with `@astrojs/node@9`. Astro 6 support (requiring Node 22+) is planned for v0.2.0.

## Install

Download and review the installer, then run it:

```bash
curl -O https://raw.githubusercontent.com/srmdn/inkwell/main/install.sh
bash install.sh
```

This will:

1. Clone Inkwell and build the binary from source
2. Clone and build the default Astro theme
3. Create a `.env` file with a generated `JWT_SECRET`
4. Run the first-time setup wizard to create your admin account

By default, everything is installed into `./inkwell/`. To use a different directory:

```bash
bash install.sh --dir /opt/inkwell
```

## Running

After install:

```bash
cd inkwell

# Start the backend API (default port 8090)
./inkwell

# In a separate terminal — start the theme (default port 4321)
cd theme && node dist/server/entry.mjs
```

For production, run both as systemd services. See [docs/configuration.md](docs/configuration.md) for all environment variables.

## First-time setup

The installer runs `--setup` for you. If you need to re-run it later:

```bash
./inkwell --setup
```

This creates the admin account interactively. You can also set `ADMIN_EMAIL` and `ADMIN_PASSWORD` in `.env` instead.

## Configuration

All configuration is via `.env`. See [docs/configuration.md](docs/configuration.md) for the full reference.

Key variables:

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8090` | Backend API port |
| `JWT_SECRET` | — | Required. Generate with `openssl rand -hex 32` |
| `CONTENT_DIR` | `content/blog` | Path to Markdown content |
| `THEME_DIR` | `theme` | Path to the theme directory |
| `THEME_BUILD_CMD` | `npm run build` | Command to rebuild the theme |
| `THEME_SERVICE` | — | systemd service name to restart after rebuild |

## API

Full API reference: [docs/api.md](docs/api.md)

## Theme

Inkwell ships with [inkwell-theme-default](https://github.com/srmdn/inkwell-theme-default) — an Astro SSR theme. The installer sets it up automatically.

To build a custom theme, see [docs/theme-contract.md](docs/theme-contract.md).

## Architecture

- **Backend**: Go, Chi router, SQLite (`modernc.org/sqlite` — pure Go, no CGo)
- **Auth**: JWT (cookie + Bearer header) + stateless CSRF (HMAC-SHA256)
- **Default theme**: Astro SSR (separate repo)
- **Content**: Markdown files with YAML frontmatter in `content/blog/<slug>/index.md`

## License

MIT
