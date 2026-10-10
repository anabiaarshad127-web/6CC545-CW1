#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
exec go run ./cmd/node -id node-b -http 127.0.0.1:8082 -seed testdata/node-b-files.json
