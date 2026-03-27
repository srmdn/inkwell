# Changelog

## v0.2.0 (2026-03-27)

- Webhook rebuild endpoint: `POST /api/webhook/rebuild` (protected by `WEBHOOK_SECRET`)
- Newsletter: subscriber management + SMTP send
  - `POST /api/subscribe` (public)
  - `GET /api/unsubscribe?token=xxx` (public)
  - `GET /api/admin/subscribers`
  - `DELETE /api/admin/subscribers/{id}`
  - `POST /api/admin/newsletter/send`
- Landing page live at foliocms.com
- GitHub repo About section (description, homepage, topics)

---

## v0.1.0 (2026-03-27)

First usable release.

- Backend API: post CRUD, auth (JWT + CSRF), rebuild trigger
- Migration runner with version tracking
- First-run setup wizard (`--setup`)
- `install.sh`: one-command installer
- Default Astro theme (`foliocms-theme-default`)
- Theme contract, API reference, configuration reference docs

---

Releases follow [Semantic Versioning](docs/versioning.md).
