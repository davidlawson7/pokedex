#!/usr/bin/env bash
set -euo pipefail

# Clone the PokeAPI static JSON mirror into _data/api-data.
# This is a one-time operation; the output is git-ignored.
# Re-run to refresh data.

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DEST="$REPO_ROOT/_data/api-data"

if [ -d "$DEST" ]; then
    echo "_data/api-data already exists. Remove it first to re-fetch."
    exit 0
fi

mkdir -p "$REPO_ROOT/_data"
git clone --depth=1 https://github.com/PokeAPI/api-data.git "$DEST"
echo "Done. Data available at: $DEST"
