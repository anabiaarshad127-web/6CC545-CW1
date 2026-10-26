# API and Protocol Documentation

## Week 2 — JSON REST Status Resource

### Purpose

The node provides a simulated status resource at `/v1/status`. It supports GET requests to read the current status and POST requests to update it. The accepted status values are `healthy` and `degraded`.

The `/healthz` endpoint serves a different purpose. It reports the node's identity and the value `ok` to indicate that the health handler is responding. The simulated status does not automatically stop the server or confirm that other nodes are available.

### Request and response contract

| Request | Success or error status | Meaning |
|---|---|---|
| `GET /healthz` | 200 | Returns the node identity and liveness status |
| `GET /v1/status` | 200 | Returns the current simulated status |
| `POST /v1/status` with valid JSON | 200 | Updates and returns the status |
| Unsupported method | 405 | Method not allowed; the handler sets the `Allow` header |
| Missing or unsupported content type | 415 | The request must use JSON |
| Malformed JSON or unknown field | 400 | The request is invalid |
| Invalid status or extra JSON value | 400 | The request does not satisfy the API rules |
| Request body exceeding 4096 bytes | 413 | The request body is too large |

### GET example

Request:

`GET /v1/status`

Example response:

```http
HTTP 200
Content-Type: application/json; charset=utf-8
```

```json
{"status":"healthy"}
```

The actual value can be `healthy` or `degraded`, depending on the most recent successful POST request.

### POST example

Request:

```http
POST /v1/status
Content-Type: application/json
```

```json
{"status":"degraded"}
```

Example response:

```http
HTTP 200
Content-Type: application/json; charset=utf-8
```

```json
{"status":"degraded"}
```

The handler first decodes the request into a temporary value and validates it. The stored status is updated only after the request passes validation.

### JSON body rules

- The request must contain one JSON object.
- The `status` field must contain `healthy` or `degraded`.
- Unknown fields are rejected.
- Malformed JSON, arrays, `null` and additional JSON values are rejected.
- The maximum request body size is 4096 bytes.
- A body of exactly 4096 bytes can be accepted if it is valid.
- A body exceeding the limit receives HTTP 413.
- `application/json; charset=utf-8` is accepted.
- A missing or unsupported content type receives HTTP 415.

### Error response format

Errors are returned as JSON objects containing an `error` field.

For example:

```http
HTTP 400
Content-Type: application/json; charset=utf-8
```

```json
{"error":"invalid_json"}
```

The HTTP status indicates whether the request succeeded. The error label gives a general explanation of why the request was rejected.

### State invariant

A rejected request must not change the status stored before that request. The handler performs decoding and validation before updating the state. The tests should verify that invalid requests preserve the previous status.

### Design rationale and limitations

Separating the status resource from the health endpoint gives each endpoint a clear responsibility. GET reads the current representation, while POST updates it. This is easier to understand and test than using one endpoint for unrelated operations.

HTTP 200 is appropriate for a successful update because the response contains the updated representation of an existing resource. HTTP 201 is more appropriate when an operation creates a new resource.

The body-size limit restricts the amount of request data read. Server timeouts control how long network operations may take, so these measures address different problems.

The JSON decoder rejects unknown fields and additional JSON values. However, Go's standard JSON decoder does not automatically reject repeated occurrences of a known field, which remains a limitation. The service is intended for local coursework demonstrations and does not provide authentication.

### Verification

The Week 2 evidence should record the REST tests, body-size boundary test, full Go test suite and `go vet` results that were actually observed. The existing evidence record describes successful GET and POST requests, a malformed POST returning HTTP 400, and a subsequent GET confirming that the status remained `degraded`. These claims should be retained only if the corresponding recorded test evidence supports them.

---

## Week 4 — TCP and UDP Inter-Node Messaging

### Purpose

Each node provides TCP and UDP listeners for exchanging validated messages with another node. Accepted messages are stored in a bounded in-memory inbox and can be retrieved through `GET /v1/messages`.

The service is designed for local coursework demonstrations. The supplied run scripts use loopback addresses, and the implementation does not provide authentication.

### Message format

Messages are encoded as JSON and contain the following fields:

| Field | Meaning |
|---|---|
| `version` | Protocol version; the current implementation accepts version 1 |
| `request_id` | Identifier for an individual request |
| `node_id` | Identifier of the sending node |
| `text` | Non-empty message text containing valid UTF-8 |

The implementation rejects messages with an unsupported protocol version, invalid identifiers, empty text or invalid UTF-8. Encoded messages are limited to 4096 bytes. The decoder also rejects unknown JSON fields and additional JSON values.

### TCP framing

TCP provides a byte stream rather than separate message boundaries. Therefore, each message is sent with a four-byte unsigned big-endian length prefix followed by the JSON payload.

The receiver reads the complete header and payload, rejecting zero-length or oversized frames. This allows the receiver to process fragmented reads without assuming that one network read contains an entire message.

### UDP transport

UDP sends each encoded JSON message as an individual datagram. The receiver validates incoming datagrams and sends an acknowledgement when it processes a valid message.

UDP does not guarantee delivery, ordering or retransmission. A sender may experience a timeout or network error when a receiver is unavailable. The implementation should not be described as providing automatic retries unless additional code confirms that behaviour.

### Acknowledgement

A successful acknowledgement uses the following JSON fields:

| Field | Meaning |
|---|---|
| `version` | Protocol version |
| `request_id` | Identifier of the request being acknowledged |
| `receiver` | Identifier of the receiving node |
| `accepted` | Indicates whether the message was accepted |

For TCP, the sender checks the acknowledgement version, request identifier, acceptance flag and receiver field.

An acknowledgement means that the receiver processed and accepted the message. It does not guarantee permanent storage. Messages are kept in memory, and the inbox is lost when the node process stops.

### HTTP inbox endpoint

`GET /v1/messages` returns the messages currently held in the node's inbox as a JSON array. The inbox retains a bounded recent history of up to 20 messages.

Example request:

`GET /v1/messages`

Illustrative response:

```http
HTTP 200
Content-Type: application/json; charset=utf-8
```

```json
[{"version":1,"request_id":"tcp-1","node_id":"node-a","text":"hello"}]
```

The response illustrates the message structure; it is not a recorded test transcript. Only GET is accepted. Other methods receive HTTP 405.

### Demonstration ports

| Node | HTTP | TCP | UDP | gRPC |
|---|---|---|---|---|
| Node A | `127.0.0.1:8081` | `127.0.0.1:8181` | `127.0.0.1:8281` | `127.0.0.1:9091` |
| Node B | `127.0.0.1:8082` | `127.0.0.1:8182` | `127.0.0.1:8282` | `127.0.0.1:9092` |

These addresses match the supplied shell run scripts.

### TCP and UDP comparison

TCP provides a reliable, ordered byte stream while the connection remains usable. The application adds explicit framing so the receiver can identify each message. TCP is suitable when ordered delivery matters, although connection management and stream handling introduce additional complexity.

UDP preserves datagram boundaries and has lower transport overhead, but packets can be lost or duplicated and delivery is not guaranteed. This implementation uses request identifiers, acknowledgements and deadlines to make communication outcomes observable. Automatic retries should not be claimed without supporting implementation evidence.

### Verification evidence

The existing Week 4 evidence record describes TCP and UDP messages being accepted in both directions between Node A and Node B. It also records successful responses from `GET /v1/messages` and connection errors when Node B was stopped.

Keep these statements only where the actual recorded commands and results support them. The submission evidence should include the targeted IPC tests, full Go test suite, race detector and `go vet` results that were actually observed.

---

## Week 5 — gRPC Remote File Metadata

### Purpose

The node provides a gRPC service that allows a client to request file metadata from a remote node. Protocol Buffers define the request and response structures, while gRPC provides the remote procedure call between the client and server.

### Service contract

The service is defined in `api/files/v1/files.proto`.

| Item | Definition |
|---|---|
| Service | `FileService` |
| RPC method | `ListFiles` |
| Request | `ListFilesRequest` |
| Response | `ListFilesResponse` |
| Request field | `requester_id` — string |
| Response fields | `files` — repeated `FileMetadata`; `node_id` — string |

Each `FileMetadata` record contains the following fields:

| Field | Type | Meaning |
|---|---|---|
| `name` | string | File name |
| `size` | int64 | File size |
| `sha256` | string | SHA-256 digest |
| `modified_unix` | int64 | Modification time represented as a Unix timestamp |

### Request and response

The client calls `ListFiles` with a `ListFilesRequest` containing the requester's ID. The server responds with file metadata and the identity of the responding node.

The service returns metadata rather than file contents. Protocol Buffers define the structure of the exchanged messages, while gRPC provides the communication mechanism.

### Error handling and limitations

A remote call may fail if the server is unavailable, the request is cancelled or the configured deadline expires. The client must handle these outcomes instead of assuming that every request will succeed. An empty file list is also a valid response and should not automatically be treated as an error.

The `requester_id` field is an unverified claim and does not establish the requester's identity. The coursework service does not provide authentication through this field. It is intended for a local teaching demonstration rather than unprotected public deployment.

The HTTP remote-file endpoint translates gRPC errors into HTTP responses. For example, an unavailable remote service produces HTTP 503, while a deadline or timeout produces HTTP 504.

### Verification

The Week 5 evidence should cover successful remote metadata listing, unavailable services, cancelled or timed-out requests, empty results and request validation where these cases are supported by the actual tests.

The full Go test suite, race detector and `go vet ./...` results should be included only when their corresponding commands were run and their observed outputs were recorded.

For submission evidence, include the actual server and client commands and their observed results. Example responses must be labelled as examples rather than presented as real test transcripts.