#!/usr/bin/env bash
set -euo pipefail

osarch=(
  darwin/amd64
  darwin/arm64
  linux/amd64
  linux/arm
  linux/arm64
  windows/amd64
)

function go_package() {
  cd "${4}"

  local outname="${3}"
  [[ $1 == windows ]] && outname="${3}.exe"

  log "=> Building ${3} for ${1}/${2}..."
  CGO_ENABLED=0 GOARCH=$2 GOOS=$1 go build \
    -ldflags "-s -w -X main.version=${version}" \
    -mod=readonly \
    -trimpath \
    -o "${outname}"

  if [[ $1 == linux ]]; then
    log "=> Packging ${3} as ${3}_${1}_${2}.tgz..."
    tar -czf "${builddir}/${3}_${1}_${2}.tgz" "${outname}"
  else
    log "=> Packging ${3} as ${3}_${1}_${2}.zip..."
    zip "${builddir}/${3}_${1}_${2}.zip" "${outname}"
  fi

  rm "${outname}"
}

function go_package_all() {
  for oa in "${osarch[@]}"; do
    local os=$(cut -d / -f 1 <<<"${oa}")
    local arch=$(cut -d / -f 2 <<<"${oa}")
    (go_package "${os}" "${arch}" "${1}" "${2}")
  done
}

function log() {
  echo "[$(date +%H:%M:%S)] $@" >&2
}

root=$(pwd)
builddir="${root}/.build"
version="$(git describe --tags --always || echo dev)"

log "Building version ${version}..."

log "Resetting output directory..."
rm -rf "${builddir}"
mkdir -p "${builddir}"

log "Building API-Server..."
go_package_all "ots" "."

log "Building OTS-CLI..."
go_package_all "ots-cli" "./cmd/ots-cli"

log "Generating SHA256SUMS file..."
(cd "${builddir}" && sha256sum * | tee SHA256SUMS)
