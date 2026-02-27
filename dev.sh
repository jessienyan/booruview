#!/bin/bash

set -euo pipefail
cd ${0%/*}

# Updates the db schema, creating a new one if it doesn't exist
migrate_db() {
    docker run \
        --rm \
        -it \
        -v "$(pwd)/api/database:/workspace" \
        -w /workspace \
        -u "$(id -u):$(id -g)" \
        keinos/sqlite3 \
        ash -c "sqlite3 ./sqlite.db < schema.sql"
}

COMPOSE="docker compose"

if [[ -x "$(command -v docker-compose)" ]]; then
    COMPOSE="docker-compose"
fi

export VITE_COMMIT_SHA=$(git rev-parse --short)
export VITE_LAST_COMMIT_DATE=$(git show -s --format=%cs)

migrate_db

$COMPOSE up --build $@
