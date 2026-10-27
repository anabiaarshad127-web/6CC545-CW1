# Weekly evidence record

**Student number:** 100859795

**Week:** Week 1

**Repository URL:** https://github.com/anabiaarshad127-web/6CC545-CW1.git

**Immutable tag:** `week01`

**Commit hash:** `660cbaa3b1386be7b9793ffb8651481b3f4d497c`

**Go version / OS:** `go1.27.1 darwin/arm64` (macOS)

## Commands actually executed

- `go version`
- `go test ./internal/config -run "^TestWeek01_" -v`
- `go test ./... -count=1 -timeout 20s`
- `go vet ./...`
- `sh scripts/run-node-a.sh`
- `go run ./cmd/request -url http://127.0.0.1:8081/healthz`
- `sh scripts/run-node-b.sh`
- `go run ./cmd/request -url http://127.0.0.1:8082/healthz`
- `sh scripts/run-node-a.sh` while Node A was already running
- `go run ./cmd/request -url http://127.0.0.1:8081/healthz` while the original Node A was running
- `go run ./cmd/request -url http://127.0.0.1:8081/healthz` after Node A was stopped

## Observed results (do not paste expected examples as results)

The Week 1 configuration tests passed, including `TestWeek01_Config`, `TestWeek01_UnknownFlag` and `TestWeek01_InvalidOptionalAddress`.

The full test suite passed with:

`go test ./... -count=1 -timeout 20s`

`go vet ./...` also completed successfully with no reported problems.

Node A started successfully on `127.0.0.1:8081`. A request to `/healthz` returned HTTP 200 and the JSON response:

`{"node_id":"node-a","status":"ok"}`

Node B started successfully on `127.0.0.1:8082`. A request to `/healthz` returned HTTP 200 and:

`{"node_id":"node-b","status":"ok"}`

After Node A was stopped, a request to port 8081 returned `connection refused`, showing that the HTTP listener was no longer running.

## Failure case and diagnosis

I tested starting another Node A while the first Node A was already using port 8081. The second process failed with an address-in-use error:

`listen tcp 127.0.0.1:8081: bind: address already in use`

The original Node A continued to respond to `/healthz`. This showed that the startup error was caused by the port already being occupied.

## Design decision and alternative

I kept the configuration, application and HTTP parts in separate packages because each part has a different purpose. The configuration checks that the node ID and addresses are valid before the server starts. The `/healthz` handler returns the node ID and an OK status in JSON. I used loopback addresses so the service runs locally for this checkpoint. I also tested what happens when the HTTP port is already being used. The second process failed with an address-in-use error, while the original node continued working.

An alternative would be to put the configuration, server startup and health response in one main function. However, this could make the project harder to manage when more features are added in later weeks. The main limitation of Week 1 is that it only provides the basic node and health endpoint.

## Remaining limitations

The Week 1 implementation only provides the basic node foundation and HTTP health endpoint. The REST state operations, TCP and UDP communication, and Protocol Buffers remote file metadata are not part of this checkpoint.

## Sources and AI-use disclosure

**Sources used:** The 6CC545 CW1 Assessment Specification and the module Start Here / Week 1 guidance.

**AI assistance:** ChatGPT was used to explain the assessment requirements, troubleshoot Go and Git commands, suggest tests, and help structure the evidence. I ran and checked the commands and results myself and reviewed the submitted work.
