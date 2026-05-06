#!/usr/bin/env bash
set -euo pipefail

MSG_FILE="${1:-}"
if [[ -z "$MSG_FILE" || ! -f "$MSG_FILE" ]]; then
  echo "commit-msg: message file not found"
  exit 1
fi

FIRST_LINE="$(head -n1 "$MSG_FILE" | tr -d '\r')"

if [[ "$FIRST_LINE" =~ ^Merge\  ]] || [[ "$FIRST_LINE" =~ ^Revert\ \" ]]; then
  exit 0
fi

PATTERN='^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([a-zA-Z0-9._-]+\))?(!)?: .+'

if [[ ! "$FIRST_LINE" =~ $PATTERN ]]; then
  echo
  echo "❌ Bad commit message:"
  echo "   $FIRST_LINE"
  echo
  echo "✅ Use Conventional Commits, examples:"
  echo "   feat(api): add pagination"
  echo "   fix(auth)!: drop legacy token"
  echo "   chore: bump deps"
  echo
  exit 1
fi
