#!/bin/bash

set -euo pipefail
cd ${0%/*}

export VITE_COMMIT_SHA=$(git rev-parse --short)
export VITE_LAST_COMMIT_DATE=$(git show -s --format=%cs)

docker compose up --build $@
