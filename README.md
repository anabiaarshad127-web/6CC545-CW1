# 6CC545 CW1 — Distributed Applications and Web Services

## Project overview

This project is a cumulative Go application developed for the 6CC545 coursework. It demonstrates a distributed service using two local nodes.

The five weekly checkpoints cover node startup and health checks, REST status updates, peer and file metadata, TCP and UDP messaging, and gRPC-based remote file metadata requests.

## Project structure

- `source/` — Go application source code, tests, configuration and launch scripts.
- `repository.txt` — GitHub repository URL.
- `implementation-log.md` — development tasks, design decisions, test results and AI-use disclosures.
- `evidence/` — evidence records for the five weekly checkpoints and final integration.
- `docs/` — architecture, API and protocol documentation, and test results.
- `cw1-commentary.pdf` — coursework commentary.

Some submission files may still be under preparation. Check that the final package contains all required files before submission.

## Requirements

- Go, using a version compatible with the project.
- Git.
- macOS or Linux with a Bash-compatible shell, or Windows with a suitable shell.

## Build and test

Open Terminal and change to the Go module directory:

```bash
cd source
```

Run the following commands separately:

```bash
go build ./...
```

```bash
go test ./... -count=1 -timeout 30s
```

```bash
go test -race ./... -count=1 -timeout 30s
```

```bash
go vet ./...
```

Record the actual results of each command in the test documentation. Do not describe a test as passed unless it has been run successfully.

## Starting the nodes

On macOS or Linux, open two separate Terminal windows.

**Terminal 1 — Node A**

```bash
cd ~/Downloads/6CC545-CW1/source
sh scripts/run-node-a.sh
```

**Terminal 2 — Node B**

```bash
cd ~/Downloads/6CC545-CW1/source
sh scripts/run-node-b.sh
```

The scripts run the nodes in the foreground. Press `Ctrl+C` in each terminal to stop its node.

## Local service addresses

| Service | Node A | Node B |
|---|---|---|
| HTTP | `127.0.0.1:8081` | `127.0.0.1:8082` |
| TCP | `127.0.0.1:8181` | `127.0.0.1:8182` |
| UDP | `127.0.0.1:8281` | `127.0.0.1:8282` |
| gRPC | `127.0.0.1:9091` | `127.0.0.1:9092` |

## Available features

- `GET /healthz` — checks node health.
- `GET /v1/status` — retrieves the current node status.
- `POST /v1/status` — validates and updates the node status.
- Peer endpoints — register and inspect peer metadata.
- `GET /v1/files` — lists file metadata in the local catalogue.
- `GET /v1/messages` — inspects messages received by a node.
- TCP and UDP messaging — exchanges messages between nodes.
- `GET /v1/remote-files?peer_id=node-b` — requests remote file metadata through gRPC.

The REST endpoints use JSON. The gRPC service uses Protocol Buffers.

## Example checks

With both nodes running, open a third Terminal window and run these commands from the `source` directory:

```bash
curl -i http://127.0.0.1:8081/healthz
```

```bash
curl -i http://127.0.0.1:8082/healthz
```

Inspect the local file catalogues:

```bash
curl -i http://127.0.0.1:8081/v1/files
```

```bash
curl -i http://127.0.0.1:8082/v1/files
```

Inspect the message inboxes:

```bash
curl -i http://127.0.0.1:8081/v1/messages
```

```bash
curl -i http://127.0.0.1:8082/v1/messages
```

The `cmd/ipcclient` program can send TCP or UDP messages. For example, send a TCP message from Node A to Node B:

```bash
go run ./cmd/ipcclient -network tcp -target 127.0.0.1:8182 -id node-a -text "hello over TCP from A to B" -request tcp-a-b-1
```

Check Node B's inbox to inspect the received message.

## Data and limitations

The nodes use loopback addresses and synthetic coursework data. The file catalogue stores metadata rather than transferring file contents. Peer information is stored locally, and automatic peer discovery is not provided.

The application does not provide TLS or authentication. Keep the services on loopback addresses and do not expose them to other networks.

## Evidence and version control

The repository maintains Git checkpoint tags for Weeks 1–5. The final submission should also include a `final-cw1` tag after final integration checks have been completed.

The submission evidence should record the relevant tag, commit hash, commands actually executed, observed results, failure diagnosis, design decisions, limitations, and sources or AI assistance.

Before submission, check that the required evidence records, implementation log, architecture document, API and protocol documentation, test results, commentary PDF, and final submission archive are present.

