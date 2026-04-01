# Changelog

## v0.3.0 (2026-04-01)

Built-in admin dashboard. No separate service or install step required.

- Admin dashboard embedded in the Go binary at `/admin`
- Post list: all posts with draft/published/scheduled filter tabs, delete with confirm
- Post editor: Milkdown WYSIWYG Markdown editor with full toolbar (headings, lists, code blocks, tables, links), hero image upload, slug auto-generation, description counter, publish date, draft toggle
- Slug rename on update: content directory moves with all assets
- Subscribers page: list all subscribers, remove with inline confirm
- Settings page: rebuild trigger with live status polling, build output and error display
- Collapsible sidebar with mobile drawer, responsive layout (375px/768px/1280px verified)
- Cmd/Ctrl+S keyboard shortcut to save in post editor
- Security: backend now binds to 127.0.0.1 instead of 0.0.0.0
- Unit tests for storage, auth, and post handlers
- `docs/api.md` updated with all v0.2.0 and v0.3.0 endpoints

---

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
