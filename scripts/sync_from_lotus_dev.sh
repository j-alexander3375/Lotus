#!/usr/bin/env bash
set -euo pipefail

# Sync this repository from the local sibling lotus_dev repository.
REMOTE_NAME="${REMOTE_NAME:-lotus_dev}"
REMOTE_PATH="${REMOTE_PATH:-/home/josh_35p/Documents/Code/lotus_dev}"
SOURCE_BRANCH="${SOURCE_BRANCH:-priv}"
TARGET_BRANCH="${1:-master}"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "Error: run this script from inside a git repository." >&2
  exit 1
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Error: tracked changes detected. Commit or stash changes before syncing." >&2
  exit 1
fi

if ! git remote get-url "$REMOTE_NAME" >/dev/null 2>&1; then
  git remote add "$REMOTE_NAME" "$REMOTE_PATH"
fi

git fetch "$REMOTE_NAME" --prune

current_branch="$(git branch --show-current)"
if [ "$current_branch" != "$TARGET_BRANCH" ]; then
  git checkout "$TARGET_BRANCH"
fi

git merge --ff-only "$REMOTE_NAME/$SOURCE_BRANCH"

echo "Synced $TARGET_BRANCH from $REMOTE_NAME/$SOURCE_BRANCH"