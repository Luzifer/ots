#!/usr/bin/env bash
set -euo pipefail

function log() {
  echo "[$(date +%H:%M:%S)] $@" >&2
}

[[ -n ${GITHUB_REF_NAME:-} ]] || {
  log "ERR: This script is intended to run on a Github Action only."
  exit 1
}

repo="ghcr.io/${GITHUB_REPOSITORY,,}"
tags=()

case "${GITHUB_REF_TYPE}" in
branch)
  # Generic build to develop: Workflow has to limit branches to master
  tags+=("${repo}:develop")
  ;;
tag)
  # Build to latest & tag: Older tags are not intended to rebuild
  tags+=("${repo}:latest" "${repo}:${GITHUB_REF_NAME}")
  ;;
*)
  log "ERR: The ref type ${GITHUB_REF_TYPE} is not handled."
  exit 1
  ;;
esac

export IFS=,
echo "docker_build_tags=${tags[*]}" >>${GITHUB_OUTPUT}
