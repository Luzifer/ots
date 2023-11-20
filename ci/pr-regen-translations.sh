#!/usr/bin/env bash
set -euo pipefail

files=(
  i18n.yaml
  src/langs/langs.js
)
translation_branch=upd-translate

function log() {
  echo "$@" >&2
}

PR_REMOTE_URL=${1:-}
[[ -n $PR_REMOTE_URL ]] || {
  log "Usage: $(basename $0) <github branch URL i.e. https://github.com/Luzifer/ots/tree/translate-de>"
  exit 1
}

remote="$(cut -d '/' -f 1-5 <<<"${PR_REMOTE_URL}").git"
branch=$(cut -d '/' -f 7 <<<"${PR_REMOTE_URL}")

git diff --exit-code >/dev/null || {
  log "FATAL: Local changes detected, stopping now."
  exit 1
}

switch_back_branch=$(git branch --show-current)
trap "git switch ${switch_back_branch}" EXIT

log "Updating branch '${branch}' of remote '${remote}'..."

log "+ Fetching remote..."
git fetch "${remote}" "${branch}"

log "+ Creating work-branch..."
if git branch | grep -q ${translation_branch}; then
  git branch -D ${translation_branch}
fi
git branch ${translation_branch} FETCH_HEAD

log "+ Switching to work-branch..."
git switch ${translation_branch}

log "+ Updating translations..."
make translate

if git diff --exit-code "${files[@]}" >/dev/null; then
  log "No changed introduced, stopping now."
fi

log "+ Committing changes..."
git add "${files[@]}"
git commit -m 'CI: Update embedded translations'

log "+ Please review these changes:"
git show

log "[Enter] to continue, [Ctrl+C] to cancel..."
read

log "+ Updating remote branch..."
git push ${remote} ${translation_branch}:${branch}

log "Updated remote PR, switching back to previous branch..."
