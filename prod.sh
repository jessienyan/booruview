#!/bin/bash

set -euo pipefail

cd ${0%/*}

mkdir -p database

# Update the DB to the latest schema, or create a new .db file if none exists
docker run \
	--rm \
	-it \
	-v "$(pwd)/database:/workspace" \
	-w /workspace \
	-u "$(id -u):$(id -g)" \
	keinos/sqlite3 \
	ash -c "sqlite3 ./sqlite.db < schema.sql"

docker compose -f docker-compose.prod.yml up -d
