# Versioning

Inkwell follows Semantic Versioning: `vMAJOR.MINOR.PATCH`

## The Three Numbers

| Segment | Bump when... | Example |
|---------|-------------|---------|
| `MAJOR` | Breaking change — config keys renamed, API routes changed, install steps change | `v1.0.0` → `v2.0.0` |
| `MINOR` | New feature, nothing existing breaks | `v0.1.0` → `v0.2.0` |
| `PATCH` | Bug fix only, nothing existing breaks | `v0.1.0` → `v0.1.1` |

If in doubt: did you add anything new? MINOR. Did you only fix something broken? PATCH.

## Milestone Definitions

### `v0.1.0` — First usable release

Tag when a stranger can clone the repo, run `--setup`, connect a theme,
and have a working CMS. Requires:

- [ ] Backend API complete
- [ ] Default Astro theme ships and works out of the box
- [ ] `install.sh` covers the full setup flow
- [ ] README has real install instructions
- [ ] Pre-built binary attached to the GitHub release

### `v0.x.x` — Unstable

API and config can still change freely between minor versions.
Document breaking changes clearly in CHANGELOG.

### `v1.0.0` — Stable

Cut when:
- Running in real-world use
- Public API (routes, request/response shapes) is settled
- Config keys and `.env` variables are stable
- Ready to commit to backwards compatibility on breaking changes

## Tagging a Release

```bash
# 1. Finish and commit all changes for the release
# 2. Update CHANGELOG.md
# 3. Commit the changelog
git add CHANGELOG.md && git commit -m "chore: release v0.1.0"

# 4. Create an annotated tag
git tag -a v0.1.0 -m "v0.1.0: first usable release"

# 5. Push commits and tag
git push origin main
git push origin v0.1.0
```

Always tag **after** committing CHANGELOG. The tag points to the final release commit.

## CHANGELOG

Keep `CHANGELOG.md` in the repo root, newest version at top. One line per change:

```markdown
## v0.2.0 — 2026-04-15

- Added media upload support
- Fixed slug validation for unicode characters

## v0.1.0 — 2026-03-27

- Initial release
```

## Working with AI on Releases

AI tools can draft CHANGELOG and release notes from commit history:

```bash
git log v0.1.0..v0.2.0 --oneline
# paste to AI: "draft release notes from these commits"
```

Rules:
- You decide the version number — AI does not bump versions
- You review the draft before publishing
- You run `git tag` and publish the release — never automate this
