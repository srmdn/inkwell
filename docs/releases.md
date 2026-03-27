# Releases

Every git tag gets a GitHub release. No tags without releases.

## What a Release Contains

| Item | Required | Notes |
|------|----------|-------|
| Annotated git tag | Yes | `git tag -a v0.1.0 -m "..."` |
| GitHub release page | Yes | Title = tag, body = CHANGELOG section |
| Pre-built binaries | Yes | See build targets below |
| Source code archive | Auto | GitHub generates `.zip` and `.tar.gz` automatically |

## Binary Build Targets

Inkwell is a single Go binary. Attach pre-built binaries for:

| Target | GOOS | GOARCH | Filename |
|--------|------|--------|----------|
| Linux 64-bit | `linux` | `amd64` | `inkwell-linux-amd64` |
| Linux ARM | `linux` | `arm64` | `inkwell-linux-arm64` |
| macOS Intel | `darwin` | `amd64` | `inkwell-darwin-amd64` |
| macOS Apple Silicon | `darwin` | `arm64` | `inkwell-darwin-arm64` |

Build all targets before publishing:

```bash
VERSION=v0.1.0

GOOS=linux  GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o inkwell-linux-amd64  ./cmd/server/
GOOS=linux  GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o inkwell-linux-arm64  ./cmd/server/
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o inkwell-darwin-amd64 ./cmd/server/
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o inkwell-darwin-arm64 ./cmd/server/
```

## Release Process

```bash
# 1. Ensure main is clean and all changes are committed
git status

# 2. Update CHANGELOG.md with this release's section
# 3. Commit it
git add CHANGELOG.md && git commit -m "chore: release v0.1.0"

# 4. Tag
git tag -a v0.1.0 -m "v0.1.0: first usable release"
git push origin main
git push origin v0.1.0

# 5. Build binaries (see targets above)

# 6. Create GitHub release via gh CLI
gh release create v0.1.0 \
  --title "v0.1.0" \
  --notes "$(cat <<'EOF'
## What's Changed

- Initial release

## Install

See [README](https://github.com/srmdn/inkwell#install) for full instructions.

Download the binary for your platform below, or build from source:

```
go install github.com/srmdn/inkwell/cmd/server@v0.1.0
```

## Full Changelog

https://github.com/srmdn/inkwell/commits/v0.1.0
EOF
)" \
  inkwell-linux-amd64 \
  inkwell-linux-arm64 \
  inkwell-darwin-amd64 \
  inkwell-darwin-arm64
```

## Release Notes Format

```markdown
## What's Changed

- Added X feature
- Fixed Y bug

## Breaking Changes (if any)

- Renamed env var FOO to BAR — update your .env

## How to Upgrade

Replace the binary and run with --migrate if prompted.

## Full Changelog

https://github.com/srmdn/inkwell/compare/v0.1.0...v0.2.0
```

## Patch Releases

For critical bug fixes: branch from the tag, fix, tag `v0.1.1`, release.
No need to wait for the next minor version.

## AI-Assisted Release Notes

Draft from commit history:

```bash
git log v0.1.0..v0.2.0 --oneline
# paste to AI: "draft release notes from these commits"
```

You review and edit before publishing. You publish — never automate.
