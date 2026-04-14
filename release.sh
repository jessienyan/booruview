#!/bin/bash

# Pushes the docker images to codeberg and adds a git tag.
# This automatically calls build_prod.sh

set -euo pipefail
IFS=$'\n\t'

cd ${0%/*}

source ./build_prod.sh

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

echo ">>> starting release"
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
echo ">>> done"
