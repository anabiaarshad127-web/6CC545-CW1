# Week 5 Evidence Record

**Student number:** 100859795

**Week:** Week 5

**Repository URL:** https://github.com/anabiaarshad127-web/6CC545-CW1.git

**Immutable tag:** `week05`

**Commit hash:** `6e5ae6974170ae83113d130813bf301543582b3a`

**Go version / OS:** `go1.27.1 darwin/arm64` (macOS)

**Protocol Buffers compiler:** `libprotoc 36.2`

## Commands actually executed

- `go build ./...`
- `go test ./... -count=1 -timeout 30s`
- `go test -race ./... -count=1 -timeout 60s`
- `go vet ./...`
- `git diff --check`
- `curl -i http://127.0.0.1:8081/healthz`
- `curl -i http://127.0.0.1:8082/healthz`
- `curl -i http://127.0.0.1:8082/v1/files`
- `curl -i -X POST http://127.0.0.1:8081/v1/peers -H 'Content-Type: application/json' -d '{"id":"node-b","http_address":"127.0.0.1:8082","grpc_address":"127.0.0.1:9092","status":"online"}'`
- `curl -i 'http://127.0.0.1:8081/v1/remote-files?peer_id=node-b'`

## Observed results

The build, full test suite, race-enabled test suite, and `go vet` completed successfully. The diff check reported no whitespace errors.

Both nodes returned HTTP 200 from `/healthz`. Node B's `/v1/files` endpoint returned HTTP 200 and listed `diagram.txt` and `notes.txt`.

Registering Node B in Node A returned HTTP 201. The remote-file request returned HTTP 200 with `node_id` set to `node-b` and metadata for both files, including their sizes, SHA-256 hashes, and modification timestamps. This verified the remote file-metadata request through gRPC.

## Design decision and alternative

I used Protocol Buffers to define the file-metadata service and gRPC for communication between nodes. The HTTP endpoint on Node A requests metadata from the registered peer and returns the result as JSON. This keeps the HTTP interface separate from the remote procedure call.

An alternative would be to request file metadata through the peer's HTTP endpoint. However, using gRPC demonstrates the remote invocation requirement for this checkpoint and provides a defined service contract.

## Remaining limitations

The service returns file metadata rather than file contents. The local coursework demonstration uses loopback addresses and insecure gRPC transport; it does not provide TLS or authentication. Peer registration is manual rather than automatic.

## Sources and AI-use disclosure

**Sources used:** The 6CC545 CW1 Assessment Specification and module guidance.

**AI assistance:** ChatGPT helped explain requirements, troubleshoot commands, and organise the documentation. I ran the tests and live requests and checked their results myself.
