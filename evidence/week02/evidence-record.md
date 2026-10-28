# Week 2 evidence record

**Student number:** 100859795

**Week:** Week 2

**Repository URL:** https://github.com/anabiaarshad127-web/6CC545-CW1.git

**Immutable tag:** `week02`

**Commit hash:** `4392a6fe9e34ef2c9edccee7ae494bddf2c6400a`

**Go version / OS:** `go1.27.1 darwin/arm64` (macOS)

## Commands actually executed

The following checks were run against a temporary copy of the immutable `week02` Git snapshot.

- `go version`
- `go test ./internal/httpapi -run '^TestWeek02_REST$' -v -count=1`
- `go test ./... -count=1 -timeout 20s`
- `go vet ./...`
- `go build -o /tmp/cw1-week02-node ./cmd/node`
- `curl -i --max-time 5 http://127.0.0.1:8081/v1/status`
- `curl -i --max-time 5 -X POST http://127.0.0.1:8081/v1/status -H "Content-Type: application/json" -d '{"status":"degraded"}'`
- `curl -i --max-time 5 http://127.0.0.1:8081/v1/status` after the valid POST
- `curl -i --max-time 5 -X POST http://127.0.0.1:8081/v1/status -H "Content-Type: application/json" -d '{"status":"banana"}'`
- `curl -i --max-time 5 http://127.0.0.1:8081/v1/status` after the invalid POST

## Observed results

**Fresh verification against the `week02` snapshot:**

The Go version was `go1.27.1 darwin/arm64`. All 13 cases in `TestWeek02_REST` passed. These covered GET, valid POST, a JSON content type with a charset, malformed JSON, unknown fields, invalid values, multiple JSON objects, null, arrays, missing or incorrect content types, oversized bodies and unsupported methods.

The full Go test suite passed with `go test ./... -count=1 -timeout 20s`. The `go vet ./...` command completed without diagnostic output, and the node executable built successfully.

The initial GET request to `/v1/status` returned HTTP 200 with `{"status":"healthy"}`. A valid POST containing `{"status":"degraded"}` returned HTTP 200 with `{"status":"degraded"}`. A subsequent GET confirmed that the updated status was retained.

Posting `{"status":"banana"}` returned HTTP 400 with `{"error":"invalid_status"}`. A further GET returned HTTP 200 with `{"status":"degraded"}`, confirming that the invalid update did not change the existing status.

The temporary node process was terminated after the HTTP checks. These fresh checks establish the behaviour of the code at the `week02` tag in the current environment; they do not, by themselves, establish the exact results observed on the original checkpoint date.

## Failure case and diagnosis

I submitted an invalid status value, `banana`, to `/v1/status`. The server returned HTTP 400 with `{"error":"invalid_status"}`. A subsequent GET confirmed that the status remained `degraded`. This demonstrated that the server rejected the invalid value without changing the existing state.

## Design decision and alternative

I added a REST resource at `/v1/status` so a client can read and update the node's status using JSON. GET returns the current status, while POST validates the submitted status before changing the state. I separated the HTTP request handling from the status state so each part has a clear responsibility. The request decoder checks the content type, limits the body to 4096 bytes, rejects unknown fields and accepts only one JSON object. Invalid requests return JSON errors and do not update the state.

An alternative would be to handle all request validation directly inside the resource handler. However, keeping decoding in a separate function makes the validation easier to maintain and test. This checkpoint still uses a local status resource; remote communication features are developed in later weeks.

## Remaining limitations

This checkpoint implements the REST status resource. It does not yet demonstrate the later TCP, UDP and gRPC file-metadata features.

## Sources and AI-use disclosure

 I have consulted RFC 9110 for HTTP semantics and the OWASP REST Security Cheat Sheet for API security guidance. AI also provided me with partial support when reviewing this week's work.
