#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

dcc() {
	docker compose -f docker-compose.test.yml $@
}

dcc build
dcc up -d --renew-anon-volumes valkey
dcc run --rm api ash -c 'rm -f $DATABASE_PATH && goose up && go test ./...'
dcc down
