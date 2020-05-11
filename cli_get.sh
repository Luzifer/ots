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
URL=${@:-}
[[ -n $URL ]] || {
	echo "Usage: $0 'URL to get the secret'"
	exit 1
}
# normalize url and extract parts
URL=${URL/|/%7C}
HOST=${URL%%/\#*}
IDPASS=${URL##*\#}
PASS=${IDPASS##*\%7C}
ID=${IDPASS%%\%7C*}
GETURL=${HOST}/api/get/${ID}

# Get the secret
resp=$(curl -s $GETURL) 
# decrypt
deciphertext=$(echo $resp| jq -r ".secret" | openssl aes-256-cbc -base64 -pass "pass:$PASS" -md md5 -d 2> /dev/null)
echo $deciphertext
