# Week 3 Evidence Record

**Student number:** 100859795

**Week:** Week 3

**Repository URL:** https://github.com/anabiaarshad127-web/6CC545-CW1.git

**Immutable tag:** `week03`

**Commit hash:** `95bc6c07a4ccc24341608def3d24e5ab4181b596`

**Go version / OS:** `go1.27.1 darwin/arm64` (macOS)

## Commands actually executed

The following commands were run against a temporary copy of the source stored at the immutable `week03` tag on 27 October 2026:

- `go test ./internal/model ./internal/state -count=1`
- `go test ./internal/state -run '^TestWeek03_' -v -count=1`
- `go test ./... -count=1 -timeout 30s`
- `go vet ./...`
- `go test -race ./... -count=1 -timeout 30s`
- `sh -n scripts/run-node-a.sh`
- `sh -n scripts/run-node-b.sh`
- `curl -i --max-time 5 http://127.0.0.1:8081/v1/files`
- `curl -i --max-time 5 http://127.0.0.1:8082/v1/files`
- `curl -i --max-time 5 -X POST http://127.0.0.1:8081/v1/peers -H "Content-Type: application/json" -d '{"id":"node-b","http_address":"127.0.0.1:8082","grpc_address":"127.0.0.1:9092","status":"unknown"}'` (run twice to check duplicate handling)

## Observed results

The model and state tests passed. All four Week 3 state tests passed: `TestWeek03_Catalog`, `TestWeek03_Registry`, `TestWeek03_ConcurrentRegistry`, and `TestWeek03_InvalidPeer`. The complete Go test suite passed, `go vet ./...` reported no diagnostics, the race-enabled test suite passed without race warnings, and both node scripts passed shell syntax checks.

The file catalogue endpoint returned HTTP 200 on both nodes. Node A returned `alpha.txt`. Node B returned `diagram.txt` and `notes.txt` in filename order.

In a clean peer-registration test, Node A started successfully and logged that it was ready. The first POST to `/v1/peers` returned HTTP 201 Created and the submitted peer details. Repeating the same POST returned HTTP 409 Conflict with `{"error":"duplicate_peer"}`.

These checks were performed on 27 October 2026 against a temporary copy of the immutable `week03` snapshot. They verify the checkpoint code as tested on that date; they do not establish that the same commands were run on the original checkpoint date.

## Failure case and diagnosis

I submitted the same peer-registration request twice to Node A. The first request returned HTTP 201 Created, while the repeated request returned HTTP 409 Conflict with `{"error":"duplicate_peer"}`. This showed that the registry rejected the duplicate registration instead of accepting the same peer twice. The observed response identified the cause of the rejected request as a duplicate peer ID.

## Design decision and alternative

I separated peer and file metadata models from the registry, catalogue, and HTTP handlers. This makes the components easier to test independently. A mutex protects shared registry state, and snapshots return copies rather than exposing the internal collections. The file catalogue validates replacement data before updating its stored contents.

An alternative would be to keep peer and file metadata directly inside the HTTP handlers. However, this would mix storage, validation, and request handling, making the application harder to maintain as more distributed features are added.

## Remaining limitations

The catalogue stores file metadata rather than transferring file contents. The peer registry is local to each node, and this checkpoint does not provide automatic peer discovery. Later distributed communication features are outside the scope of this Week 3 checkpoint.

## Sources and AI-use disclosure

**Sources used:** The 6CC545 CW1 Assessment Specification and module guidance.

**AI assistance:** ChatGPT helped explain requirements, troubleshoot commands, and organise the documentation. I ran the listed tests and HTTP requests and checked their observed results myself.
