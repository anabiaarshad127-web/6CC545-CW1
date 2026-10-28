# Week 4 Evidence Record

**Student number:** 100859795

**Week:** Week 4

**Repository URL:** https://github.com/anabiaarshad127-web/6CC545-CW1.git

**Immutable tag:** `week04`

**Commit hash:** `4f89950e372afeb275a8afcc31206876c1f4248f`

**Go version / OS:** `go1.27.1 darwin/arm64` (macOS)

## Commands actually executed

- `go test ./internal/ipc -count=1`
- `go test ./... -count=1 -timeout 30s`
- `go test -race ./... -count=1 -timeout 30s`
- `go vet ./...`

Both nodes were started using the following commands:

- `go run ./cmd/node -id node-b -http 127.0.0.1:8082 -tcp 127.0.0.1:8182 -udp 127.0.0.1:8282`
- `go run ./cmd/node -id node-a -http 127.0.0.1:8081 -tcp 127.0.0.1:8181 -udp 127.0.0.1:8281`

The following client commands tested TCP and UDP in both directions:

- `go run ./cmd/ipcclient -network tcp -target 127.0.0.1:8182 -id node-a -text "hello over TCP from A to B" -request tcp-a-b-1`
- `go run ./cmd/ipcclient -network udp -target 127.0.0.1:8282 -id node-a -text "hello over UDP from A to B" -request udp-a-b-1`
- `go run ./cmd/ipcclient -network tcp -target 127.0.0.1:8181 -id node-b -text "hello over TCP from B to A" -request tcp-b-a-1`
- `go run ./cmd/ipcclient -network udp -target 127.0.0.1:8281 -id node-b -text "hello over UDP from B to A" -request udp-b-a-1`

The inbox endpoints were checked using:

- `curl -i --max-time 5 http://127.0.0.1:8081/v1/messages`
- `curl -i --max-time 5 http://127.0.0.1:8082/v1/messages`

## Observed results

The IPC package tests passed. The full Go test suite passed, the race-enabled test suite passed without reported race warnings, and `go vet ./...` completed without diagnostics.

Both nodes started successfully and reported that they were ready on their configured loopback addresses.

All four client requests returned successful JSON acknowledgements with `"accepted":true`. The TCP and UDP requests from Node A were acknowledged by Node B, and the requests from Node B were acknowledged by Node A.

Both inbox endpoints returned HTTP `200 OK`. Node A's inbox contained the TCP and UDP messages sent by Node B. Node B's inbox contained the TCP and UDP messages sent by Node A.

These live checks were performed on 27 October 2026 against a temporary copy of the immutable `week04` snapshot. They do not establish that the same checks were performed on the original checkpoint date.

## Failure case and diagnosis

The Week 4 failure tests cover oversized TCP frames, malformed UDP messages, oversized UDP datagrams, and incomplete TCP frames that time out. These tests check that invalid or incomplete messages are rejected or dropped without changing the inbox. They also check that the network servers shut down within the expected timeout. This helps prevent invalid messages from being accepted and ensures that incomplete TCP communication does not produce a successful acknowledgement.

The tests use local network connections and check the results programmatically. Their results should be confirmed by running the Week 4 IPC tests against the immutable `week04` snapshot.

## Design decision and alternative

I separated the TCP and UDP communication functions into the `internal/ipc` package. TCP uses message framing so the receiver can identify complete messages, while UDP processes individual datagrams. The protocol validates messages and applies size limits. Each node stores received messages in an inbox, which can be inspected through the HTTP `/v1/messages` endpoint.

An alternative would be to place the communication code directly in the main application package. Keeping it separate makes the protocol and network behaviour easier to test and maintain.

## Remaining limitations

The demonstration uses loopback addresses on one machine. The service log states that authentication is not enabled, so the prototype should not be treated as secure communication. UDP does not provide the same delivery guarantees as TCP. This checkpoint demonstrates message exchange and acknowledgement, rather than production-ready networking. The gRPC file-metadata functionality belongs to Week 5.

## Sources and AI-use disclosure

I have consulted RFC 9293 for TCP and RFC 768 for UDP. I have received partial help from artificial intelligence. 
