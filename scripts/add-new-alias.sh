#!/bin/bash
set -euo pipefail

usage() {
  printf '%s\n' \
    'Usage: add-new-alias.sh <alias> <url>' \
    '' \
    'Alias: 1-64 letters, numbers, hyphens, or underscores.' \
    'URL: starts with http:// or https://.'
}

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  usage
  exit 0
fi

if [[ $# -ne 2 ]]; then
  usage >&2
  exit 2
fi

ALIAS=$1
URL=$2

if [[ ! $ALIAS =~ ^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$ ]]; then
  printf 'Invalid alias: %s\n' "$ALIAS" >&2
  exit 2
fi

if [[ $ALIAS == "list" || $ALIAS == "new" ]]; then
  printf 'Alias is reserved: %s\n' "$ALIAS" >&2
  exit 2
fi

if [[ -z $URL || ( $URL != http://* && $URL != https://* ) ]]; then
  printf 'Invalid URL: %s\n' "$URL" >&2
  exit 2
fi

curl --fail --show-error --silent \
  --request POST http://localhost:8080/new \
  --header 'Content-Type: application/json' \
  --data "{\"url\":\"$URL\",\"alias\":\"$ALIAS\"}"
