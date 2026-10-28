# Week 5 Evidence Record

**Student number:** 100859795

**Week:** Week 5

**Repository URL:** https://github.com/anabiaarshad127-web/6CC545-CW1.git

**Immutable tag:** `week05`

**Commit hash:** `6e5ae6974170ae83113d130813bf301543582b3a`

**Go version / OS:** `go1.27.1 darwin/arm64` (macOS)

**Protocol Buffers compiler:** `libprotoc 36.2`

## Commands actually executed

The following build and verification commands were run from the archived Week 5 source snapshot:

- `go build ./...`
- `go test ./... -count=1 -timeout 30s`
- `go test -race ./... -count=1 -timeout 60s`
- `go vet ./...`

The following commands were used to inspect the node options and locate the supplied seed data:

- `go run ./cmd/node -h`
- `grep -R -n -E 'diagram\.txt|notes\.txt|data-dir|data_dir|file.*dir|root.*dir' . --exclude='*.pb.go' --exclude='*.sum' --exclude-dir=.git`

The live checks used the temporary audit copy at `tmp/anabia-week05-audit/source`. Node B was started with `-seed testdata/node-b-files.json`.

From the archived source directory, the nodes were started with these commands in separate terminals:
- Node A: `go run ./cmd/node -id node-a -http 127.0.0.1:8081 -grpc 127.0.0.1:9091`
- Node B: `go run ./cmd/node -id node-b -http 127.0.0.1:8082 -grpc 127.0.0.1:9092 -seed testdata/node-b-files.json`
- `curl -i --max-time 5 http://127.0.0.1:8081/healthz`
- `curl -i --max-time 5 http://127.0.0.1:8082/healthz`
- `curl -i --max-time 5 http://127.0.0.1:8082/v1/files`
- `curl -i --max-time 5 -X POST http://127.0.0.1:8081/v1/peers -H 'Content-Type: application/json' -d '{"id":"node-b","http_address":"127.0.0.1:8082","grpc_address":"127.0.0.1:9092","status":"online"}'`
- `curl -i --max-time 10 'http://127.0.0.1:8081/v1/remote-files?peer_id=node-b'`

## Observed results

The build, full test suite, race-enabled test suite, and `go vet` completed successfully against the archived Week 5 source snapshot.

During live verification on 27 October 2026, both nodes returned HTTP 200 from `/healthz`. With the supplied seed file loaded, Node B's `/v1/files` endpoint returned HTTP 200 and listed `diagram.txt` and `notes.txt`.

Registering Node B with Node A returned HTTP 201 Created. The remote-file request through Node A returned HTTP 200 with `node_id` set to `node-b` and metadata for both files, including their sizes, SHA-256 hashes, and modification timestamps. This verified the remote file-metadata request through gRPC.

These live results were obtained during the later verification session using a temporary archive of the `week05` checkpoint. They should not be interpreted as evidence that the live requests were run on the original checkpoint date.

## Failure case and diagnosis

The live verification recorded successful peer registration and a successful remote-file metadata request. However, the evidence collected here does not record a failed gRPC request or an unavailable peer. Therefore, a specific remote-service failure cannot be confirmed from these results. An additional failure test should be run against the Week 5 snapshot to verify how the application responds when remote metadata cannot be retrieved.

## Design decision and alternative

I used Protocol Buffers to define the file-metadata service and gRPC for communication between nodes. The HTTP endpoint on Node A requests metadata from the registered peer and returns the result as JSON. This keeps the HTTP interface separate from the remote procedure call.

An alternative would be to request file metadata through the peer's HTTP endpoint. However, using gRPC demonstrates the remote invocation requirement for this checkpoint and provides a defined service contract.

## Remaining limitations

The service returns file metadata rather than file contents. The local coursework demonstration uses loopback addresses and insecure gRPC transport; it does not provide TLS or authentication. Peer registration is manual rather than automatic. The file catalogue in this demonstration is loaded from supplied synthetic metadata rather than being proof of a live filesystem scan.

## Sources and AI-use disclosure

I have consulted the official Protocol Buffers documentation and gRPC guidance on deadlines. I have also got partial support from artificial intelligence.
