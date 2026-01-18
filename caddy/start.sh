#!/bin/ash

set -euo pipefail

/etc/caddy/generate_caddyfile.sh
caddy run --config /etc/caddy/Caddyfile --adapter caddyfile
