#!/usr/bin/env bash
# Build the admin dashboard and copy the output into the Go embed directory.
# Run this before `go build` whenever admin-ui source changes.
#
# Usage: bash scripts/build-admin.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADMIN_SRC="$REPO_ROOT/admin-ui"
EMBED_DEST="$REPO_ROOT/internal/adminui/dist"

if [ ! -d "$ADMIN_SRC" ]; then
  echo "Error: admin-ui directory not found at $ADMIN_SRC" >&2
  exit 1
fi

echo "Installing admin-ui dependencies..."
(cd "$ADMIN_SRC" && npm install)

echo "Building admin-ui..."
(cd "$ADMIN_SRC" && npm run build)

echo "Copying dist to $EMBED_DEST..."
rm -rf "$EMBED_DEST"
cp -r "$ADMIN_SRC/dist" "$EMBED_DEST"

echo "Done. Run 'go build ./cmd/server/' to compile the binary."
