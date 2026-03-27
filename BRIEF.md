# Project Brief: Folio

## What it does
A lightweight, self-hostable CMS. Single Go binary, SQLite, ships with a
default Astro theme out of the box. No Docker, no external dependencies.

## Who uses it
Solo developers and indie developers who self-host on a cheap VPS and want
a real CMS (editor, frontend, dashboard) without the
overhead of WordPress or Ghost.

## Stack
- Backend: Go, Chi router, SQLite
- Default theme: Astro SSR
- Binary: single compiled Go binary
- Deploy: systemd service, no Docker

## Deploy target
Any Linux VPS. One systemd service for the backend. Separate optional
service for the Astro theme if used.

## Hard requirements
- Single Go binary: no Docker, no external runtime dependencies
- SQLite only: no PostgreSQL or MySQL
- Ships with a default Astro theme (swappable via theme contract)
- Real Markdown editor (Milkdown) in the admin dashboard
- Theme contract: any frontend that reads `content/blog/<slug>/index.md`
  and calls the REST API works as a theme

## I don't care about (AI can decide)
- Folder structure: follow Go and Astro conventions
- Dashboard visual design: clean and minimal is enough
- API versioning scheme

## Out of scope (not in this version)
- Multi-site / multi-tenant support
- Plugin system
- Media library, S3, or CDN integration
- Comments
- Multi-user / roles
- Newsletter
- Anything beyond a blog CMS (no portfolio sections, no link directories, etc.)

## Post-v0.1 (revisit after first release)
- Domain: check `getfolio.com`, `foliocms.com`, `folio.dev`, `folio.app`
- Landing page: static Astro site, hosted independently
- Default Astro theme: ship as separate repo (`foliocms-theme-default`)
- Rebuild mechanism: replace subprocess approach with webhook
- Newsletter support
- Consider GitHub org (`folio-cms`) if project grows
