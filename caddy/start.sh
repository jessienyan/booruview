#!/bin/ash

set -euo pipefail

/etc/caddy/generate_caddyfile.sh

MEDIA_PROXY_HOST_SCHEMELESS=${MEDIA_PROXY_HOST#http://}
export MEDIA_PROXY_HOST_SCHEMELESS=${MEDIA_PROXY_HOST_SCHEMELESS#https://}

caddy run --config /etc/caddy/Caddyfile --adapter caddyfile
