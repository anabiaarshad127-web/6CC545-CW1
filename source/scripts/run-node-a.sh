#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
exec go run ./cmd/node -id node-a -http 127.0.0.1:8081 -tcp 127.0.0.1:8181 -udp 127.0.0.1:8281 -grpc 127.0.0.1:9091 -seed testdata/node-a-files.json
