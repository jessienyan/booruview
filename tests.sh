#!/bin/bash

set -euo pipefail
IFS=$'\n\t'
cd ${0%/*}

dcc() {
	docker compose -f docker-compose.test.yml $@
}

suppress() {
	set +e;

	output=$($@ 2>&1)
	exitcode=$?
	if [[ $exitcode != 0 ]]; then
		echo "ERROR: $output"
		exit $exitcode
	fi

	set -e;
}

mkdir -p .gotestcache

echo ">>> building..."
suppress dcc build
suppress dcc up -d --renew-anon-volumes valkey
dcc run --rm --quiet api
suppress dcc down
