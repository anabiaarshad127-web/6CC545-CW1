#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
exec go run ./cmd/node -id node-b -http 127.0.0.1:8082 -tcp 127.0.0.1:8182 -udp 127.0.0.1:8282 -grpc 127.0.0.1:9092 -seed testdata/node-b-files.json
