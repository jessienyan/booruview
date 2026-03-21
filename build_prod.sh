#!/bin/bash

# Builds the images for running booruview in production mode

set -euo pipefail
IFS=$'\n\t'

cd ${0%/*}

IMG_PREFIX=codeberg.org/jessienyan/booruview

# Exported for re-use in release.sh
export API_IMG=$IMG_PREFIX/api
export CADDY_IMG=$IMG_PREFIX/caddy
export CLIENT_IMG=$IMG_PREFIX/client
export VALKEY_IMG=$IMG_PREFIX/valkey

COMMIT=$(git rev-parse --short HEAD)
DATE=$(git show -s --format=%cs HEAD)

echo ">>> building API image"
docker build --quiet -t $API_IMG --build-arg COMMIT_HASH=$COMMIT -f api/Dockerfile.prod api/
echo ">>> building valkey image"
docker build --quiet -t $VALKEY_IMG valkey/
echo ">>> building client image"
docker build --quiet -t $CLIENT_IMG client/

echo ">>> bundling assets"
mkdir -p client/dist
docker run --rm \
    -e VITE_COMMIT_SHA=$COMMIT \
    -e VITE_LAST_COMMIT_DATE=$DATE \
    -v ./client/dist:/app/dist \
    $CLIENT_IMG yarn build

echo ">>> building caddy image"
docker build --quiet -t $CADDY_IMG -f caddy/Dockerfile .
echo ">>> done"
