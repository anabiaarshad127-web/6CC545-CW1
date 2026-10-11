# 6CC545 CW1 — Distributed Applications and Web Services

## Project overview

This project is a cumulative Go application developed for the 6CC545 coursework. It demonstrates a basic distributed service with two local nodes.

The implementation currently covers the first four weekly checkpoints: node startup, REST status updates, peer and file metadata, and TCP/UDP message communication.

## Requirements

- Go toolchain compatible with the project.
- macOS or Linux with a Bash-compatible shell, or Windows PowerShell.
- Git for version control.

## Build and test

Run these commands from the `source` directory:

```bash
go build ./...
go test ./... -count=1 -timeout 30s
go test -race ./... -count=1 -timeout 30s
go vet ./...
```

## Starting the nodes

On macOS or Linux, open two separate Terminal windows.

In the first window, run:

```bash
sh scripts/run-node-a.sh
```

In the second window, run:

```bash
sh scripts/run-node-b.sh
```

The launch scripts start the nodes in the foreground. Press `Ctrl+C` in each window to stop the corresponding node.

## Local addresses

| Service | Node A | Node B |
|---|---|---|
| HTTP | `127.0.0.1:8081` | `127.0.0.1:8082` |
| TCP | `127.0.0.1:8181` | `127.0.0.1:8182` |
| UDP | `127.0.0.1:8281` | `127.0.0.1:8282` |

## Available features

- `GET /healthz` — checks whether a node is healthy.
- `GET /v1/status` — reads the current node status.
- `POST /v1/status` — validates and updates the node status.
- Peer endpoints — manage and inspect peer metadata.
- `GET /v1/files` — lists file metadata stored in the local catalogue.
- `GET /v1/messages` — inspects messages received by a node.
- TCP and UDP messaging — sends messages between the two local nodes.

The REST resources use JSON. Input validation and tests cover invalid requests. The TCP and UDP implementation includes message validation and size limits.

## Example checks

After starting both nodes, check their health:

```bash
curl -i http://127.0.0.1:8081/healthz
curl -i http://127.0.0.1:8082/healthz
```

Inspect their file catalogues and message inboxes:

```bash
curl -i http://127.0.0.1:8081/v1/files
curl -i http://127.0.0.1:8082/v1/files
curl -i http://127.0.0.1:8081/v1/messages
curl -i http://127.0.0.1:8082/v1/messages
```

The `cmd/ipcclient` program can send TCP or UDP messages. For example, with both nodes running:

```bash
go run ./cmd/ipcclient -network tcp -target 127.0.0.1:8182 -id node-a -text "hello over TCP from A to B" -request tcp-a-b-1
```

Check Node B's inbox to inspect the received message.

## Data and limitations

The nodes use loopback addresses and synthetic coursework data. The file catalogue stores metadata; it does not transfer file contents. Peer information is stored locally, and the application does not yet provide automatic peer discovery.

This implementation does not provide TLS or authentication. Keep the services on loopback addresses and do not expose them to other networks.

The gRPC file-metadata request remains to be completed and verified for Week 5.

## Evidence and version control

Weekly implementation details, test results, design decisions, and AI-use disclosures are recorded in:

- `docs/evidence-record.md`
- `docs/implementation-log.md`

Weekly checkpoint tags are maintained in Git. Each checkpoint should be tested and reviewed before it is committed and tagged.
