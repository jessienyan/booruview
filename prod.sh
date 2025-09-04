#!/bin/bash

set -ex
cd ${0%/*}
PREFIX=codeberg.org/jessienyan/booruview

# tag images as {unixts}-{commit}
TAG=$(git show -s --format=%ct-%h master)

API_IMG=$PREFIX/api
CADDY_IMG=$PREFIX/caddy
CLIENT_IMG=$PREFIX/client
VALKEY_IMG=$PREFIX/valkey

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

docker tag $API_IMG $API_IMG:$TAG
docker tag $API_IMG $API_IMG:latest
docker push $API_IMG:$TAG
docker push $API_IMG:latest

docker tag $CADDY_IMG $CADDY_IMG:$TAG
docker tag $CADDY_IMG $CADDY_IMG:latest
docker push $CADDY_IMG:$TAG
docker push $CADDY_IMG:latest

docker tag $VALKEY_IMG $VALKEY_IMG:$TAG
docker tag $VALKEY_IMG $VALKEY_IMG:latest
docker push $VALKEY_IMG:$TAG
docker push $VALKEY_IMG:latest
