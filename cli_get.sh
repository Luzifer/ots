#!/bin/bash
set -euo pipefail

deps=(curl jq)
for cmd in "${deps[@]}"; do
	which ${cmd} >/dev/null || {
		echo "'${cmd}' util is required for this script"
		exit 1
	}
done

# Get URL from CLI argument
url="${1:-}"
[[ -n $url ]] || {
	echo "Usage: $0 'URL to get the secret'"
	exit 1
}
# normalize url and extract parts
url="${url/|/%7C}"
host="${url%%/\#*}"
idpass="${url##*\#}"
pass="${idpass##*\%7C}"
id="${idpass%%\%7C*}"
geturl="${host}/api/get/${id}"

# fetch secret and decrypt to STDOUT
curl -sSf "${geturl}" | jq -r ".secret" |
	openssl aes-256-cbc -base64 -pass "pass:${pass}" -md md5 -d 2>/dev/null
