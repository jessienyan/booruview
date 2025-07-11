#!/bin/bash

set -ex
cd ${0%/*}
PREFIX=ghcr.io/jessienyan
COMMIT=$(git rev-parse --short master)
DATE=$(git show -s --format=%cs master)
TAG=$COMMIT-$DATE

API_IMG=$PREFIX/booruview-api
CADDY_IMG=$PREFIX/booruview-caddy
CLIENT_IMG=$PREFIX/booruview-client
VALKEY_IMG=$PREFIX/booruview-valkey

docker build -t $API_IMG --build-arg COMMIT_HASH=$COMMIT -f api/Dockerfile.prod api/
docker build -t $VALKEY_IMG valkey/
docker build -t $CLIENT_IMG client/

docker run --rm \
    -e VITE_COMMIT_SHA=$COMMIT \
    -e VITE_LAST_COMMIT_DATE=$DATE \
    -v ./client/dist:/app/dist \
    $CLIENT_IMG yarn build
docker build -t $CADDY_IMG -f caddy/Dockerfile .

docker tag $API_IMG $API_IMG:$TAG
docker tag $API_IMG $API_IMG:latest

docker tag $CADDY_IMG $CADDY_IMG:$TAG
docker tag $CADDY_IMG $CADDY_IMG:latest

docker tag $VALKEY_IMG $VALKEY_IMG:$TAG
docker tag $VALKEY_IMG $VALKEY_IMG:latest

docker push $API_IMG
docker push $CADDY_IMG
docker push $VALKEY_IMG
