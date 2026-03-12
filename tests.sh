#!/bin/sh

alias dcc="docker compose -f docker-compose.test.yml"

dcc up -d --renew-anon-volumes valkey
dcc run --rm api sh -c "rm -f $DATABASE_PATH && goose up && go test ./..."
dcc down
