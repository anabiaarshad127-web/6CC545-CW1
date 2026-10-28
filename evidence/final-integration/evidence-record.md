# Final Integration Evidence Record

## Student number
100859795

## Week
final-integration

## Repository URL
https://github.com/anabiaarshad127-web/6CC545-CW1.git

## Immutable tag
`final-cw1`.

## Commit hash
`a2ed3558777cfd884bb7235ce23c65ed3e0417e0` (the verified integration commit).

## Go version / OS
Go `go1.27.1`, macOS `darwin/arm64`.

## Commands actually executed
The following commands and requests were executed during final integration verification:

- `go version`
- `go mod verify`
- `go build ./...`
- `go test ./... -count=1 -timeout 30s`
- `go test -race ./... -count=1 -timeout 60s`
- `go vet ./...`
- `curl -i --max-time 5 http://127.0.0.1:8081/healthz`
- `curl -i --max-time 5 http://127.0.0.1:8082/healthz`
- `curl -i --max-time 5 http://127.0.0.1:8081/v1/files`
- `curl -i --max-time 5 http://127.0.0.1:8082/v1/files`
- GET and POST requests to `/v1/status`, including valid, invalid and malformed JSON test data.
- POST and GET requests to `/v1/peers` to register both peers and test duplicate registration.
- `go run ./cmd/ipcclient` using TCP and UDP to send messages in both directions.
- GET requests to `/v1/messages` on both nodes.
- GET requests to `/v1/remote-files?peer_id=node-b` on Node A and `/v1/remote-files?peer_id=node-a` on Node B.
- A GET request to Node A's remote-file endpoint while Node B was stopped, followed by a health check of Node A.

## Observed results
The build, unit tests, race tests, static analysis and module verification completed successfully.

Both Node A and Node B returned HTTP 200 from `/healthz`. Their local file catalogues were different: Node A returned `alpha.txt`, while Node B returned `diagram.txt` and `notes.txt`.

The status API initially returned `healthy`. A valid POST changed the status to `degraded`. Invalid status data returned HTTP 400 with `invalid_status`, and malformed JSON returned HTTP 400 with `invalid_json`. The status remained `degraded` after both invalid requests.

Registering each peer returned HTTP 201. Both peer lists showed the registered peer. Registering Node B with Node A a second time returned HTTP 409 with `duplicate_peer`.

TCP and UDP messages were accepted in both directions. The receiving node's `/v1/messages` endpoint returned HTTP 200 and showed the corresponding message content and request ID.

Both remote-file requests returned HTTP 200. Node A retrieved Node B's `diagram.txt` and `notes.txt` metadata. Node B retrieved Node A's `alpha.txt` metadata. The responses included file sizes, SHA-256 hashes and modification timestamps.

## Failure case and diagnosis
Node B was stopped while Node A remained running. A request to `http://127.0.0.1:8081/v1/remote-files?peer_id=node-b` returned HTTP 503 Service Unavailable with `{"error":"remote_unavailable"}`. Node A's `/healthz` endpoint continued to return HTTP 200 with status `ok`.

This observed result indicates that the remote-file request reports the unavailable peer, while Node A's health endpoint remains responsive. An occupied-port startup test and a complete shutdown-and-restart test have not been verified as part of this final integration run.

## Design decision and alternative
The application uses two separate nodes with HTTP endpoints for health, status, peer registration, local files, message history and remote-file requests. TCP and UDP provide message communication, while gRPC provides remote file metadata through a Protocol Buffers service. The peer registry and file catalogue allow each node to keep track of peer details and file metadata.

This design keeps HTTP handling, network communication and remote metadata retrieval in separate parts of the application. A central service or a shared database could be considered as an alternative, but it would introduce a different architecture and dependency. The current implementation instead demonstrates communication between separate nodes.

## Remaining limitations
Peer registration and catalogues are in memory, so they are not persistent across process restarts. Peers must be registered with the relevant node. The current tests use local loopback addresses and sample file metadata; they do not demonstrate deployment across separate machines. TLS and authentication have not been implemented. Occupied-port startup and a complete shutdown-and-restart cycle remain to be verified for final integration.

## Sources and AI-use disclosure
Sources included the 6CC545 CW1 assessment brief, final integration guide, project source code, official Go, Protocol Buffers and gRPC documentation, and actual test outputs.

AI helped me understand the requirements.
