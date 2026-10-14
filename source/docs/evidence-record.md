# Weekly evidence record

Student number: 100859795

Week: Week 1

Repository URL: https://github.com/anabiaarshad127-web/6CC545-CW1.git

Immutable tag: week01

Commit hash: 660cbaa3b1386be7b9793ffb8651481b3f4d497c

Go version / OS: go1.27.1 darwin/arm64 (macOS)

## Commands actually executed

* `go version`
* `go test ./internal/config -run "^TestWeek01_" -v`
* `go test ./... -count=1 -timeout 20s`
* `go vet ./...`
* `sh scripts/run-node-a.sh`
* `go run ./cmd/request -url http://127.0.0.1:8081/healthz`
* `sh scripts/run-node-b.sh`
* `go run ./cmd/request -url http://127.0.0.1:8082/healthz`
* `sh scripts/run-node-a.sh` while Node A was already running
* `go run ./cmd/request -url http://127.0.0.1:8081/healthz` while the original Node A was running
* `go run ./cmd/request -url http://127.0.0.1:8081/healthz` after Node A was stopped

## Observed results (do not paste expected examples as results)

The Week 1 configuration tests passed, including `TestWeek01_Config`, `TestWeek01_UnknownFlag`, and `TestWeek01_InvalidOptionalAddress`.

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

I kept the configuration, application, and HTTP parts in separate packages because each part has a different purpose. The configuration checks that the node ID and addresses are valid before the server starts. The `/healthz` handler returns the node ID and an OK status in JSON. I used loopback addresses so the service runs locally for this checkpoint. I also tested what happens when the HTTP port is already being used. The second process failed with an address-in-use error, while the original node continued working.

An alternative would be to put the configuration, server startup, and health response in one main function. However, this could make the project harder to manage when more features are added in later weeks. The main limitation of Week 1 is that it only provides the basic node and health endpoint.

## Remaining limitations

The Week 1 implementation only provides the basic node foundation and HTTP health endpoint. The REST state operations, TCP and UDP communication, and Protocol Buffers remote file metadata are not part of this checkpoint.

## Sources and AI-use disclosure

Sources used: the 6CC545 CW1 Assessment Specification and the module Start Here / Week 1 guidance.

AI assistance: ChatGPT was used to explain the assessment requirements, troubleshoot Go and Git commands, suggest tests, and help structure the evidence. I ran and checked the commands and results myself and reviewed the submitted work.

---

# Week 2 evidence record

Week: Week 2

Repository URL: https://github.com/anabiaarshad127-web/6CC545-CW1.git

Immutable tag: week02

Commit hash: 4392a6f

## Commands actually executed

* `go test ./internal/httpapi -run "^TestWeek02_REST$" -v -count=1`
* `curl -i http://127.0.0.1:8081/v1/status`
* `curl -i -X POST http://127.0.0.1:8081/v1/status -H "Content-Type: application/json" -d '{"status":"degraded"}'`
* `curl -i -X POST http://127.0.0.1:8081/v1/status -H "Content-Type: application/json" -d '{"status":"banana"}'`
* `go test ./... -count=1 -timeout 20s`
* `go vet ./...`
* `git diff --cached --check`

## Observed results

All 13 cases in `TestWeek02_REST` passed. These covered GET, valid POST, a JSON content type with a charset, malformed JSON, unknown fields, invalid values, multiple JSON objects, null, arrays, missing or incorrect content types, oversized bodies, and unsupported methods.

The initial GET returned HTTP 200 with `{"status":"healthy"}`. A valid POST changed the status to `degraded` and returned HTTP 200. A subsequent GET confirmed that the changed status was retained.

Posting `{"status":"banana"}` returned HTTP 400 with `{"error":"invalid_status"}`. A follow-up GET confirmed that the status remained `degraded`, so the invalid update did not change the state.

The full Go test suite passed, and `go vet ./...` completed without reported problems. The staged-diff check produced no output.

## Design decision and alternative

I added a REST resource at `/v1/status` to let a client read and update the node's status using JSON. GET returns the current status, while POST validates the submitted status before changing the state. I separated the HTTP request handling from the status state so each part has a clear responsibility. The request decoder checks the content type, limits the body to 4096 bytes, rejects unknown fields, and accepts only one JSON object. Invalid requests return JSON errors and do not update the state.

An alternative would be to handle all request validation directly inside the resource handler. However, keeping decoding in a separate function makes the validation easier to maintain and test. This checkpoint still uses a local status resource; remote communication features will be developed in later weeks.

## Remaining limitations

This checkpoint implements the REST status resource only. The later TCP, UDP, and gRPC file-metadata features are not yet demonstrated by this Week 2 evidence.

## Sources and AI-use disclosure

Sources used: the 6CC545 CW1 Assessment Specification and module guidance.

AI assistance: ChatGPT helped explain requirements, troubleshoot commands, and organise the evidence. I ran the tests and requests and checked the results myself.

---

# Week 3 evidence record

Week: Week 3

Repository URL: https://github.com/anabiaarshad127-web/6CC545-CW1.git

Immutable tag: week03

Commit hash: 95bc6c07a4ccc24341608def3d24e5ab4181b596

Go version / OS: go1.27.1 darwin/arm64 (macOS)

## Commands actually executed

* `go test ./internal/model -count=1`
* `go test ./internal/model ./internal/state -count=1`
* `go test ./internal/state -run "^TestWeek03_" -v -count=1`
* `go test ./internal/app ./internal/httpapi ./internal/state ./internal/model -count=1`
* `go test ./... -count=1 -timeout 30s`
* `go vet ./...`
* `go test -race ./... -count=1 -timeout 30s`
* `go test ./internal/state ./internal/httpapi ./internal/model -list 'TestWeek03_.*'`
* `sh -n scripts/run-node-a.sh`
* `sh -n scripts/run-node-b.sh`
* `curl -i http://127.0.0.1:8082/v1/files`

## Observed results

The model, peer registry, file catalogue, and peer API tests passed. The full Go test suite passed, `go vet ./...` completed without reported problems, and the race-enabled test suite passed without race warnings.

Node B's file catalogue returned HTTP 200 and contained `diagram.txt` and `notes.txt`, sorted by filename. Node A's catalogue returned `alpha.txt`. Posting a valid peer to Node A returned HTTP 201. Repeating the same peer request returned HTTP 409, and a GET request showed one peer in the registry.

## Design decision and alternative

I separated peer and file metadata models from the registry, catalogue, and HTTP handlers. This makes each part easier to test. A mutex protects shared state, and snapshots return copies instead of exposing the internal collections. The file catalogue validates a complete replacement before updating the stored data, so invalid input does not replace the existing catalogue.

An alternative would be to keep peers and file metadata directly inside the HTTP handlers. However, this would mix storage, validation, and request handling, making the application harder to maintain as more distributed features are added.

## Remaining limitations

The catalogue stores metadata only and does not transfer file contents. The peer registry is local to each node; this checkpoint does not yet provide automatic peer discovery. TCP, UDP, and gRPC communication remain for later checkpoints.

## Sources and AI-use disclosure

Sources used: the 6CC545 CW1 Assessment Specification and module guidance.

AI assistance: ChatGPT helped explain requirements, troubleshoot commands, and organise the documentation. I ran the tests and requests and checked the results myself.

---

# Week 4 evidence record

Week: Week 4

Repository URL: https://github.com/anabiaarshad127-web/6CC545-CW1.git

Immutable tag: week04

Commit hash: 4f89950e372afeb275a8afcc31206876c1f4248f

## Commands actually executed

* `go test ./internal/ipc -count=1`
* `go test ./... -count=1 -timeout 20s`
* `go test -race ./... -count=1 -timeout 30s`
* `go vet ./...`
* `go run ./cmd/ipcclient -network tcp -target 127.0.0.1:8182 -id node-a -text "hello over TCP from A to B" -request tcp-a-b-1`
* `go run ./cmd/ipcclient -network udp -target 127.0.0.1:8282 -id node-a -text "hello over UDP from A to B" -request udp-a-b-1`
* `go run ./cmd/ipcclient -network tcp -target 127.0.0.1:8181 -id node-b -text "hello over TCP from B to A" -request tcp-b-a-1`
* `go run ./cmd/ipcclient -network udp -target 127.0.0.1:8281 -id node-b -text "hello over UDP from B to A" -request udp-b-a-1`
* `curl -sS http://127.0.0.1:8082/v1/messages`
* `curl -sS http://127.0.0.1:8081/v1/messages`

## Observed results

The Week 4 IPC tests passed, including TCP and UDP round trips, TCP frame fragmentation and size limits, invalid messages, cancelled clients, oversized TCP frames, oversized UDP datagrams, and incomplete TCP frame timeouts.

The full Go test suite passed. The race-enabled test suite also passed without reported race warnings. `go vet ./...` completed without reported problems.

Both nodes started successfully on their configured loopback addresses. TCP and UDP client commands were used to send messages in both directions, and the message inbox endpoints were queried to inspect received messages.

## Design decision and alternative

I separated the TCP and UDP communication code into the IPC package. The message protocol validates incoming messages and applies size limits. TCP uses framing so the receiver can identify complete messages, while UDP handles individual datagrams. The inbox stores received messages so they can be inspected through the HTTP API.

An alternative would be to put the network communication code directly inside the main application package. However, separating it makes the protocol and network behaviour easier to test and maintain.

## Remaining limitations

The messages are exchanged between local nodes using loopback addresses. This checkpoint demonstrates message communication, not secure communication. The gRPC file-metadata request remains for Week 5.

## Sources and AI-use disclosure

Sources used: the 6CC545 CW1 Assessment Specification and module guidance.

AI assistance: ChatGPT helped explain requirements, troubleshoot commands, and organise the evidence. I ran the tests and checked the results myself.

---

# Week 5 evidence record

Week: Week 5

Repository URL: https://github.com/anabiaarshad127-web/6CC545-CW1.git

Immutable tag: week05

Commit hash: 6e5ae6974170ae83113d130813bf301543582b3a

Go version / OS: go1.27.1 darwin/arm64 (macOS)

Protocol Buffers compiler: `libprotoc 36.2`

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

Sources used: the 6CC545 CW1 Assessment Specification and module guidance.

AI assistance: ChatGPT helped explain requirements, troubleshoot commands, and organise the documentation. I ran the tests and live requests and checked their results myself.