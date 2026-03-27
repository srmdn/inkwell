# Inkwell

A lightweight, self-hostable CMS. Single Go binary. SQLite. Ships with a default Astro theme. No Docker required.

## Why

Most open source CMS options are either too heavy (WordPress, Ghost) or too opinionated about the frontend (Payload, Strapi). Inkwell is built for developers who self-host on a cheap VPS and want:

- A real Markdown editor out of the box
- A default frontend theme they can swap
- Zero external dependencies — no PostgreSQL, no Redis, no Docker
- One systemd service and done

## Status

Early development. Not ready for production use.

## Architecture

- **Backend**: Go, Chi router, SQLite
- **Default theme**: Astro (SSR)
- **Theme contract**: any frontend that reads `content/blog/<slug>/index.md` and calls the REST API works as a theme

## License

MIT
