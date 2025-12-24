#!/bin/bash

set -euo pipefail
cd ${0%/*}

create_db() {
    if [[ -e database/sqlite.db ]]; then
        return
    fi

    docker run \
        --rm \
        -it \
        -v "$(pwd)/database:/workspace" \
        -w /workspace \
        -u "$(id -u):$(id -g)" \
        keinos/sqlite3 \
        ash -c "sqlite3 sqlite.db < schema.sql"

    echo "created db"
}

COMPOSE="docker compose"

if [[ -x "$(command -v docker-compose)" ]]; then
    COMPOSE="docker-compose"
fi

export VITE_COMMIT_SHA=$(git rev-parse --short master)
export VITE_LAST_COMMIT_DATE=$(git show -s --format=%cs master)

create_db

$COMPOSE up --build $@
