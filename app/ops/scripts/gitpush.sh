#!/bin/bash
set -euo pipefail

MESSAGE="${1:-chore: update work-in-progress}"

echo "[*] Add changes"
git add .

echo "[*] Commit with message: $MESSAGE"
git commit -m "$MESSAGE" || echo "[!] Nothing to commit"

BRANCH=$(git symbolic-ref --short HEAD)

echo "[*] Push to origin/$BRANCH"
git push origin "$BRANCH"

echo "[OK] Push done"
