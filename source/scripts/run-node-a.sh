#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
exec go run ./cmd/node -id node-a -http 127.0.0.1:8081 -seed testdata/node-a-files.json
