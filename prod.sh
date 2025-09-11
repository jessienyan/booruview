#!/bin/bash

set -e
cd ${0%/*}

# tag images as {unixts}-{commit}
IMG_TAG=$(git show -s --format=%ct-%h master)
IMG_PREFIX=codeberg.org/jessienyan/booruview

API_IMG=$IMG_PREFIX/api
CADDY_IMG=$IMG_PREFIX/caddy
CLIENT_IMG=$IMG_PREFIX/client
VALKEY_IMG=$IMG_PREFIX/valkey

COMMIT=$(git rev-parse --short master)
DATE=$(git show -s --format=%cs master)

# returns the git release tag as "YYYY-MM-DD". If a tag with that name already exists
# a revision is appended, e.g. "YYYY-MM-DD.1"
get_release_tag() {
    TAG=$(date "+%F")
    # looking at previous commits prevents re-tagging the same commit but with a higher revision
    PREV_TAGS=$(git tag -l --points-at HEAD~1)
    REVISION=$($PREV_TAGS | grep $TAG | wc -l)

    if [[ $REVISION > 0 ]]; then
        TAG=$TAG.$REVISION
    fi

    echo $TAG
}

docker build -t $API_IMG --build-arg COMMIT_HASH=$COMMIT -f api/Dockerfile.prod api/
docker build -t $VALKEY_IMG valkey/
docker build -t $CLIENT_IMG client/

mkdir -p client/dist
docker run --rm \
    -e VITE_COMMIT_SHA=$COMMIT \
    -e VITE_LAST_COMMIT_DATE=$DATE \
    -v ./client/dist:/app/dist \
    $CLIENT_IMG yarn build
docker build -t $CADDY_IMG -f caddy/Dockerfile .

docker tag $API_IMG $API_IMG:$IMG_TAG
docker tag $API_IMG $API_IMG:latest
docker push $API_IMG:$IMG_TAG
docker push $API_IMG:latest

docker tag $CADDY_IMG $CADDY_IMG:$IMG_TAG
docker tag $CADDY_IMG $CADDY_IMG:latest
docker push $CADDY_IMG:$IMG_TAG
docker push $CADDY_IMG:latest

docker tag $VALKEY_IMG $VALKEY_IMG:$IMG_TAG
docker tag $VALKEY_IMG $VALKEY_IMG:latest
docker push $VALKEY_IMG:$IMG_TAG
docker push $VALKEY_IMG:latest

git tag "$(get_release_tag)" && git push --tags
