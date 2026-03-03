#!/bin/bash

# Builds the images for running booruview in production mode.
# If TAG_AND_PUSH=1 then the images will be pushed to codeberg and
# a new release is created

set -euo pipefail
cd ${0%/*}

IMG_PREFIX=codeberg.org/jessienyan/booruview

API_IMG=$IMG_PREFIX/api
CADDY_IMG=$IMG_PREFIX/caddy
CLIENT_IMG=$IMG_PREFIX/client
VALKEY_IMG=$IMG_PREFIX/valkey

COMMIT=$(git rev-parse --short HEAD)
DATE=$(git show -s --format=%cs HEAD)

# returns the git release tag as "YYYY-MM-DD". If a tag with that name already exists
# a revision is appended, e.g. "YYYY-MM-DD--rev1"
get_release_tag() {
    CURRENT_TAG=$(git tag --points-at HEAD)
    if [[ $CURRENT_TAG != "" ]]; then
        echo $CURRENT_TAG
        return
    fi

    TAG=$DATE
    REVISION=$(git tag --list "$TAG*" | wc -l)

    if [[ $REVISION > 0 ]]; then
        TAG=$TAG--rev$REVISION
    fi

    echo $TAG
}

RELEASE_TAG=$(get_release_tag)

echo ">>> building client image"
docker build --quiet -t $CLIENT_IMG client/
echo ">>> bundling assets"
mkdir -p client/dist
docker run --rm \
    -e VITE_COMMIT_SHA=$COMMIT \
    -e VITE_LAST_COMMIT_DATE=$DATE \
    -v ./client/dist:/app/dist \
    $CLIENT_IMG yarn build
echo ">>> building API image"
docker build --quiet -t $API_IMG --build-arg COMMIT_HASH=$COMMIT -f api/Dockerfile.prod .
echo ">>> building valkey image"
docker build --quiet -t $VALKEY_IMG valkey/

echo ">>> building caddy image"
docker build --quiet -t $CADDY_IMG -f caddy/Dockerfile .

if [[ ${TAG_AND_PUSH:-0} = 1 ]]; then
    echo ">>> pushing API image"
    docker tag $API_IMG $API_IMG:$RELEASE_TAG
    docker tag $API_IMG $API_IMG:latest
    docker push --quiet $API_IMG:$RELEASE_TAG
    docker push --quiet $API_IMG:latest

    echo ">>> pushing caddy image"
    docker tag $CADDY_IMG $CADDY_IMG:$RELEASE_TAG
    docker tag $CADDY_IMG $CADDY_IMG:latest
    docker push --quiet $CADDY_IMG:$RELEASE_TAG
    docker push --quiet $CADDY_IMG:latest

    echo ">>> pushing valkey image"
    docker tag $VALKEY_IMG $VALKEY_IMG:$RELEASE_TAG
    docker tag $VALKEY_IMG $VALKEY_IMG:latest
    docker push --quiet $VALKEY_IMG:$RELEASE_TAG
    docker push --quiet $VALKEY_IMG:latest

    git tag "$RELEASE_TAG" && git push --tags
    echo ">>> release tagged as $RELEASE_TAG"
fi

echo ">>> done"
