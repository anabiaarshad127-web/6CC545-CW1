# 6CC545 CW1 — Distributed Applications and Web Services

## Project overview

This repository contains my cumulative coursework project for the 6CC545 Distributed Applications and Web Services module. The application has been developed through five weekly checkpoints.

The project includes node configuration, HTTP health and status endpoints, peer registration, file metadata catalogues, TCP and UDP messaging, and gRPC communication using Protocol Buffers to request metadata from a remote node.

This project is a coursework prototype for demonstrating distributed application concepts. It exchanges file metadata rather than transferring actual file contents.

## Repository

GitHub repository: https://github.com/anabiaarshad127-web/6CC545-CW1

To clone the repository, run:

```bash
git clone https://github.com/anabiaarshad127-web/6CC545-CW1.git
cd 6CC545-CW1
```

If the repository has already been cloned, open a terminal in the existing project directory.

## Project structure

- `source/` — Go application source code, tests, scripts and test data.
- `repository.txt` — Repository information and final checkpoint reference.
- `cw1-commentary.pdf` — Technical commentary for the coursework.
- `implementation-log.md` — Development decisions, testing and AI-use records.
- `evidence/week01/evidence-record.md` — Week 1 checkpoint evidence.
- `evidence/week02/evidence-record.md` — Week 2 checkpoint evidence.
- `evidence/week03/evidence-record.md` — Week 3 checkpoint evidence.
- `evidence/week04/evidence-record.md` — Week 4 checkpoint evidence.
- `evidence/week05/evidence-record.md` — Week 5 checkpoint evidence.
- `evidence/final-integration/evidence-record.md` — Final integration evidence.
- `docs/architecture.png` — Architecture diagram.
- `docs/api-and-protocols.md` — HTTP API, messaging and gRPC documentation.
- `docs/test-results.txt` — Recorded test results.
- `source/docs/week01-design-note.md` — Week 1 design note.
- `source/docs/week03-architecture.md` — Week 3 architecture documentation.

Check that all assessment-required files are present in the submission package before submitting.

## Prerequisites

For building and testing, you need:

- Go 1.25.0 or a compatible later version.
- A terminal with shell support.
- A supported operating system.

The project has been tested on macOS. The recorded environment and results are described in the coursework evidence.

The generated Protocol Buffers files are included in the project. The `protoc` compiler and Go protobuf generators are required only if the generated files need to be recreated.

## Build and test

From the repository root, enter the Go module directory:

```bash
cd source
```

Build the application:

```bash
go build ./...
```

Run the test suite:

```bash
go test ./... -count=1 -timeout 30s
```

Run the tests with the race detector:

```bash
go test -race ./... -count=1 -timeout 30s
```

Check the code:

```bash
go vet ./...
```

Verify the Go modules:

```bash
go mod verify
```

Run the project check script:

```bash
sh scripts/check.sh
```

The check script runs the build, tests and static checks specified in the script. It does not replace the separate race-detector and module-verification commands.

Refer to `docs/test-results.txt` for recorded results. Run the commands again against the checkout being submitted if you need to verify its current state.

## Running the nodes

Start each node in a separate terminal. Run the commands below from the repository root.

**Terminal 1 — Node A**

```bash
sh source/scripts/run-node-a.sh
```

**Terminal 2 — Node B**

```bash
sh source/scripts/run-node-b.sh
```

The scripts build and start the nodes using their configured settings. Keep both terminals open while testing.

Both nodes use loopback addresses for local testing.

| Service | Node A | Node B |
|---|---|---|
| HTTP | `127.0.0.1:8081` | `127.0.0.1:8082` |
| TCP | `127.0.0.1:8181` | `127.0.0.1:8182` |
| UDP | `127.0.0.1:8281` | `127.0.0.1:8282` |
| gRPC | `127.0.0.1:9091` | `127.0.0.1:9092` |

## Checking node health

With both nodes running, open a third terminal, navigate to your cloned repository, and enter the `source/` directory. Replace the example path below with the actual location of your repository.

```bash
cd /path/to/6CC545-CW1/source
```

Replace `/path/to/6CC545-CW1` with the actual location of your project.

Check Node A:

```bash
go run ./cmd/request -url http://127.0.0.1:8081/healthz
```

Check Node B:

```bash
go run ./cmd/request -url http://127.0.0.1:8082/healthz
```

A successful health request returns HTTP 200 and a JSON response containing the node identity and an `ok` status.

## HTTP endpoints

The application provides the following endpoints:

| Endpoint | Purpose |
|---|---|
| `/healthz` | Checks node health and reports its identity. |
| `/v1/status` | Reads or updates the simulated node status. |
| `/v1/files` | Lists local file metadata. |
| `/v1/peers` | Lists or registers peers. |
| `/v1/messages` | Reads messages accepted by the node. |
| `/v1/remote-files` | Requests file metadata from a registered peer using gRPC. |

Check Node A's health:

```bash
curl -i http://127.0.0.1:8081/healthz
```

List local file metadata:

```bash
curl -i http://127.0.0.1:8081/v1/files
```

Read the current status:

```bash
curl -i http://127.0.0.1:8081/v1/status
```

Update the status to degraded:

```bash
curl -i -X POST \
  http://127.0.0.1:8081/v1/status \
  -H 'Content-Type: application/json' \
  --data '{"status":"degraded"}'
```

The status endpoint accepts `healthy` and `degraded`. Invalid requests are rejected according to the API contract in `docs/api-and-protocols.md`.

## Peer registration and remote metadata

Run these commands from the `source/` directory while both nodes are running.

Register Node B with Node A:

```bash
curl -i -X POST \
  http://127.0.0.1:8081/v1/peers \
  -H 'Content-Type: application/json' \
  --data-binary @testdata/peer-b.json
```

A successful registration returns HTTP 201.

List registered peers:

```bash
curl -i http://127.0.0.1:8081/v1/peers
```

Request remote file metadata:

```bash
curl -i \
  'http://127.0.0.1:8081/v1/remote-files?peer_id=node-b'
```

When the peer is registered and its gRPC service is available, the remote metadata request returns HTTP 200. The recorded test data includes filenames such as `diagram.txt` and `notes.txt`.

Registering a duplicate peer returns HTTP 409 with a `duplicate_peer` error. If the remote gRPC service is unavailable, the request may return HTTP 503.

## TCP and UDP messaging

The project includes an IPC client for testing TCP and UDP communication. Run the following command from the `source/` directory while Node B is running.

Send a TCP message to Node B:

```bash
go run ./cmd/ipcclient \
  -network tcp \
  -target 127.0.0.1:8182 \
  -id node-a \
  -text "hello from node-a" \
  -request demo-tcp
```

For UDP, change `-network tcp` to `-network udp` and use a different request ID.

The client waits for an acknowledgement. Successful communication depends on the target node being available and reachable. UDP does not guarantee delivery or ordering.

Read messages accepted by Node B:

```bash
curl -i http://127.0.0.1:8082/v1/messages
```

The message inbox keeps a limited recent history in memory. Messages are not stored permanently.

## Regenerating Protocol Buffers files

Regenerate the Go bindings only if the Protocol Buffers definitions have changed and the generated files need updating.

The generation script is located at `source/scripts/generate.sh`. It requires `protoc` and the relevant Go generators to be available as specified by the project.

From the repository root, run:

```bash
cd source
sh scripts/generate.sh
```

Review any generated changes and run the build and test commands afterwards. Regeneration is not necessary for normal application use.

## Troubleshooting

- **A port is already in use:** Check whether another node process is running and stop it before restarting.
- **A health request fails:** Confirm that the relevant node is running and that the correct HTTP port is being used.
- **Remote metadata returns HTTP 503:** Check that the peer is registered and its gRPC service is available.
- **Peer registration returns HTTP 409:** The peer may already be registered. Check the output of `GET /v1/peers`.
- **The generation script cannot find `protoc`:** Install the required Protocol Buffers compiler and ensure it is available on your `PATH`.
- **A test-data file cannot be found:** Check that the command is being run from the `source/` directory when it uses a relative path.

## Cleanup

To stop the nodes, press `Ctrl+C` in Terminal 1 and Terminal 2.

If a port remains occupied, check for a node process that has not stopped before starting the application again. Do not terminate unrelated processes.

## Safety and limitations

- The supplied node scripts bind services to loopback addresses for local testing. They are not configured for public deployment.
- TLS and authentication are not implemented.
- Peer registration and file metadata catalogues are held in memory and are not persistent.
- The message inbox is held in memory and is cleared when the node stops.
- Remote file requests return metadata rather than actual file contents.
- Remote metadata requests depend on the selected peer and its gRPC service being available.
- The `requester_id` in a gRPC request does not prove the requester's identity.
- UDP delivery is not guaranteed.
- This project is a coursework prototype and should not be treated as a production-ready file-sharing service.

## Further documentation

See the following files for more information:

- `docs/api-and-protocols.md` — API and protocol documentation.
- `docs/test-results.txt` — Recorded test results.
- `docs/architecture.png` — Architecture diagram.
- `evidence/week01/` to `evidence/week05/` — Evidence records for the weekly checkpoints.
- `evidence/final-integration/` — Evidence for final integration.
- `implementation-log.md` — Implementation history and AI-use disclosure.
- `source/docs/week01-design-note.md` — Week 1 design note.
- `source/docs/week03-architecture.md` — Week 3 architecture documentation.
