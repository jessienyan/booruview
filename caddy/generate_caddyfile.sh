#!/bin/ash

# Generates a Caddyfile by including optional components, if they are enabled.
# Components are enabled via ENV vars. Each component has its own Caddyfile
# that, if included, is appended to the end of the final Caddyfile.
#
# This script lives in /etc/caddy. Component files are in /etc/caddy/components

set -euo pipefail

cd ${0%/*}
COMPONENT_FILES=components/

cp $COMPONENT_FILES/Caddyfile .
[[ $USE_GRAFANA = 1 ]] && cat < $COMPONENT_FILES/Caddyfile.grafana >> Caddyfile
[[ $USE_MEDIA_PROXY = 1 ]] && cat < $COMPONENT_FILES/Caddyfile.proxy >> Caddyfile

exit 0
