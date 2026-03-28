#!/bin/bash

set -euo pipefail
IFS=$'\n\t'
cd ${0%/*}

dcc() {
	docker compose -f docker-compose.test.yml $@
}

mkdir -p .gotestcache

echo ">>> building..."
dcc build --quiet
dcc up -d --renew-anon-volumes valkey
dcc run --rm --quiet api ash -c 'rm -f $DATABASE_PATH && goose up && go test ./...'
dcc down --quiet
