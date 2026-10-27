# Week 4 Evidence Record

**Student number:** 100859795

**Week:** Week 4

**Repository URL:** https://github.com/anabiaarshad127-web/6CC545-CW1.git

**Immutable tag:** `week04`

**Commit hash:** `4f89950e372afeb275a8afcc31206876c1f4248f`

**Go version / OS:** `go1.27.1 darwin/arm64` (macOS)

## Commands actually executed

- `go test ./internal/ipc -count=1`
- `go test ./... -count=1 -timeout 20s`
- `go test -race ./... -count=1 -timeout 30s`
- `go vet ./...`
- `go run ./cmd/ipcclient -network tcp -target 127.0.0.1:8182 -id node-a -text "hello over TCP from A to B" -request tcp-a-b-1`
- `go run ./cmd/ipcclient -network udp -target 127.0.0.1:8282 -id node-a -text "hello over UDP from A to B" -request udp-a-b-1`
- `go run ./cmd/ipcclient -network tcp -target 127.0.0.1:8181 -id node-b -text "hello over TCP from B to A" -request tcp-b-a-1`
- `go run ./cmd/ipcclient -network udp -target 127.0.0.1:8281 -id node-b -text "hello over UDP from B to A" -request udp-b-a-1`
- `curl -sS http://127.0.0.1:8082/v1/messages`
- `curl -sS http://127.0.0.1:8081/v1/messages`

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

**Sources used:** The 6CC545 CW1 Assessment Specification and module guidance.

**AI assistance:** ChatGPT helped explain requirements, troubleshoot commands, and organise the evidence. I ran the tests and checked the results myself.
