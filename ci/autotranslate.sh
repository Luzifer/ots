#!/usr/bin/env bash
set -euo pipefail

function log() {
  echo "[$(date +%H:%M:%S)] $@" >&2
}

translation_keys=($(
  jq -r '. | keys | .[]' src/langs/en.json
))

for lang_file in src/langs/*.json; do
  lang=$(echo ${lang_file} | sed -E 's@.*/([^\/\.]*)\.json@\1@')
  log "Processing ${lang}..."

  target_lang=$(jq -r ".__lang" ${lang_file} | grep -v null || echo "")
  [[ -n $target_lang ]] || {
    log "  + Missing '__lang' key, cannot translate."
    continue
  }

  for tk in "${translation_keys[@]}"; do
    [[ $(jq -r ".[\"${tk}\"]" ${lang_file}) == null ]] || continue
    log "  + Missing '${tk}', fetching..."

    source_str=$(jq -r ".[\"${tk}\"]" src/langs/en.json)

    translation="$(
      curl -sSf -X POST "${DEEPL_API_ENDPOINT}" \
        -H "Authorization: DeepL-Auth-Key ${DEEPL_API_KEY}" \
        -F "text=${source_str}" \
        -F "target_lang=${target_lang}" |
        jq -r '.translations[0].text'
    )"

    jq -S --arg t "${translation}" ".[\"${tk}\"]=\$t" ${lang_file} >${lang_file}.tmp
    mv ${lang_file}.tmp ${lang_file}

  done

done
