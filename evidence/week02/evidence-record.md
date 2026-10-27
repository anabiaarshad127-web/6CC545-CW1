# Week 2 evidence record

**Student number:** 100859795

**Week:** Week 2

**Repository URL:** https://github.com/anabiaarshad127-web/6CC545-CW1.git

**Immutable tag:** `week02`

**Commit hash:** `4392a6f`

## Commands actually executed

- `go test ./internal/httpapi -run "^TestWeek02_REST$" -v -count=1`
- `curl -i http://127.0.0.1:8081/v1/status`
- `curl -i -X POST http://127.0.0.1:8081/v1/status -H "Content-Type: application/json" -d '{"status":"degraded"}'`
- `curl -i -X POST http://127.0.0.1:8081/v1/status -H "Content-Type: application/json" -d '{"status":"banana"}'`
- `go test ./... -count=1 -timeout 20s`
- `go vet ./...`
- `git diff --cached --check`

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

**Sources used:** The 6CC545 CW1 Assessment Specification and module guidance.

**AI assistance:** ChatGPT helped explain requirements, troubleshoot commands, and organise the evidence. I ran the tests and requests and checked the results myself.
