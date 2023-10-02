#!/bin/bash
set -euo pipefail

: ${INSTANCE:=https://ots.fyi} # Where to reach the API of the instance (omit trailing slash)

deps=(curl jq)
for cmd in "${deps[@]}"; do
  which ${cmd} >/dev/null || {
    echo "'${cmd}' util is required for this script"
    exit 1
  }
done

# Get secret from CLI argument
SECRET=${1:-}
[[ -n $SECRET ]] || {
  echo "Usage: $0 'secret to share'"
  exit 1
}

# Generate a random 20 character password
pass=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 20 || true)

# Encrypt the secret
ciphertext=$(echo "${SECRET}" | openssl aes-256-cbc -base64 -A -pass "pass:${pass}" -iter 300000 -md sha512 2>/dev/null)

# Create a secret and extract the secret ID
id=$(
  curl -sSf \
    -X POST \
    -H 'content-type: application/json' \
    -d "$(jq --arg secret "${ciphertext}" -cn '{"secret": $secret}')" \
    "${INSTANCE}/api/create" |
    jq -r '.secret_id'
)

# Display URL to user
echo -e "Secret is now available at:\n${INSTANCE}/#${id}%7C${pass}"
