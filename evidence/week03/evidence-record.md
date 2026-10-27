# Week 3 Evidence Record

**Student number:** 100859795

**Week:** Week 3

**Repository URL:** https://github.com/anabiaarshad127-web/6CC545-CW1.git

**Immutable tag:** `week03`

**Commit hash:** `95bc6c07a4ccc24341608def3d24e5ab4181b596`

**Go version / OS:** `go1.27.1 darwin/arm64` (macOS)

## Commands actually executed

- `go test ./internal/model -count=1`
- `go test ./internal/model ./internal/state -count=1`
- `go test ./internal/state -run "^TestWeek03_" -v -count=1`
- `go test ./internal/app ./internal/httpapi ./internal/state ./internal/model -count=1`
- `go test ./... -count=1 -timeout 30s`
- `go vet ./...`
- `go test -race ./... -count=1 -timeout 30s`
- `go test ./internal/state ./internal/httpapi ./internal/model -list 'TestWeek03_.*'`
- `sh -n scripts/run-node-a.sh`
- `sh -n scripts/run-node-b.sh`
- `curl -i http://127.0.0.1:8082/v1/files`

## Observed results

The model, peer registry, file catalogue, and peer API tests passed. The full Go test suite passed, `go vet ./...` completed without reported problems, and the race-enabled test suite passed without race warnings.

Node B's file catalogue returned HTTP 200 and contained `diagram.txt` and `notes.txt`, sorted by filename. Node A's catalogue returned `alpha.txt`. Posting a valid peer to Node A returned HTTP 201. Repeating the same peer request returned HTTP 409, and a GET request showed one peer in the registry.

## Design decision and alternative

I separated peer and file metadata models from the registry, catalogue, and HTTP handlers. This makes each part easier to test. A mutex protects shared state, and snapshots return copies instead of exposing the internal collections. The file catalogue validates a complete replacement before updating the stored data, so invalid input does not replace the existing catalogue.

An alternative would be to keep peers and file metadata directly inside the HTTP handlers. However, this would mix storage, validation, and request handling, making the application harder to maintain as more distributed features are added.

## Remaining limitations

The catalogue stores metadata only and does not transfer file contents. The peer registry is local to each node; this checkpoint does not yet provide automatic peer discovery. TCP, UDP, and gRPC communication remain for later checkpoints.

## Sources and AI-use disclosure

**Sources used:** The 6CC545 CW1 Assessment Specification and module guidance.

**AI assistance:** ChatGPT helped explain requirements, troubleshoot commands, and organise the documentation. I ran the tests and requests and checked the results myself.
