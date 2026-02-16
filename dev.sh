#!/bin/bash

set -e

COMPOSE="docker compose"

if [[ -x "$(command -v docker-compose)" ]]; then
    COMPOSE="docker-compose"
fi

export VITE_COMMIT_SHA=$(git rev-parse --short)
export VITE_LAST_COMMIT_DATE=$(git show -s --format=%cs)

$COMPOSE up --build $@
