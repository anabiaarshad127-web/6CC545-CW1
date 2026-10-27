# API and Protocol Documentation

## Week 2 — JSON REST Status Resource

### Purpose

The node provides a simulated status resource at `/v1/status`. It supports `GET` requests to read the current status and `POST` requests to update it. The accepted status values are `healthy` and `degraded`.

The `/healthz` endpoint has a different purpose. It reports the node's identity and an `ok` value to indicate that the health handler is responding. The simulated status does not automatically stop the server or confirm that other nodes are available.

### Request and response contract

| Request | HTTP status | Meaning |
|---|---:|---|
| `GET /healthz` | 200 | Returns the node identity and health response |
| `GET /v1/status` | 200 | Returns the current simulated status |
| `POST /v1/status` with valid JSON | 200 | Updates and returns the status |
| Unsupported method | 405 | Method not allowed |
| Missing or unsupported content type | 415 | The request must use JSON |
| Malformed JSON or unknown field | 400 | The request is invalid |
| Invalid status or additional JSON value | 400 | The request does not satisfy the API rules |
| Request body exceeding 4096 bytes | 413 | The request body is too large |

The precise response body depends on the request and the error encountered. Clients should use the HTTP status code and JSON error identifier to understand a failed request.

### GET example

Request:

```http
GET /v1/status HTTP/1.1
```

Illustrative response:

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
```

```json
{"status":"healthy"}
```

The example shows one possible response. The actual value may be `healthy` or `degraded`, depending on the most recent successful update.

### POST example

Request:

```http
POST /v1/status HTTP/1.1
Content-Type: application/json
```

```json
{"status":"degraded"}
```

Illustrative response:

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
```

```json
{"status":"degraded"}
```

The handler decodes and validates the submitted status before updating the stored value. If validation fails, the request should be rejected without changing the previous status.

These request and response examples explain the API contract. They are not presented as actual command output.

### JSON body validation

The REST handler applies validation to incoming requests. The documented rules are:

- The request must contain one JSON object.
- The `status` field must contain `healthy` or `degraded`.
- Unknown fields are rejected.
- Malformed JSON, arrays, `null` and additional JSON values are rejected.
- The maximum request body size is 4096 bytes.
- A valid body of up to the permitted size can be accepted.
- `application/json; charset=utf-8` is accepted.
- A missing or unsupported content type is rejected with HTTP 415.

The JSON decoder rejects unknown fields, but Go's standard JSON decoder does not automatically reject repeated occurrences of a known field. This is a limitation to consider when extending the validation rules.

### Error response format

Errors are returned as JSON objects containing an `error` field.

For example, an invalid JSON request may produce:

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
```

```json
{"error":"invalid_json"}
```

This is an illustrative example. The actual error identifier depends on the failure. For instance, the recorded invalid-status request returned `{"error":"invalid_status"}`.

The HTTP status communicates the general outcome, while the error identifier provides more specific information about the rejected request.

### State invariant

A rejected request must not change the status stored before that request. The handler should complete decoding and validation before updating the state.

This rule prevents an invalid request from partially updating the resource. Tests should check both the error response and the stored status after a rejected request.

### Design rationale and limitations

Separating the status resource from the health endpoint gives each endpoint a clear responsibility. `GET` reads the current representation, while `POST` updates it. This design is easier to understand and test than using one endpoint for unrelated operations.

HTTP 200 is appropriate for a successful update because the response returns the updated representation of an existing resource. HTTP 201 is generally used when an operation creates a new resource.

The body-size limit restricts the amount of request data processed. Server timeouts address a different concern by limiting how long certain network operations can take.

The status is simulated rather than a complete health-monitoring mechanism. A value of `degraded` does not automatically stop the node, and a value of `healthy` does not prove that another node is reachable. The service is intended for local coursework demonstrations and does not provide authentication.

### Verification

The Week 2 evidence records that all 13 `TestWeek02_REST` cases passed. The tests cover GET and POST requests, JSON content types, malformed JSON, unknown fields, invalid status values, multiple JSON values, `null`, arrays, oversized bodies and unsupported methods.

The recorded manual checks showed an initial GET returning HTTP 200 with `{"status":"healthy"}`. A valid POST changed the status to `degraded`. A POST with `{"status":"banana"}` returned HTTP 400 with `{"error":"invalid_status"}`. A subsequent GET confirmed that the rejected update had not changed the stored status.

The evidence also records a successful full Go test run and `go vet ./...`. These results support the behaviour tested in the recorded run; they do not demonstrate authentication or secure communication between remote nodes.

---

## Week 4 — TCP and UDP Inter-Node Messaging

### Purpose

Each node provides TCP and UDP listeners for exchanging validated messages with another node. Accepted messages are stored in a bounded in-memory inbox and can be retrieved through `GET /v1/messages`.

The implementation demonstrates two transport protocols in a local coursework environment. The supplied run scripts use loopback addresses, and the implementation does not provide authentication.

### Message format

Messages are encoded as JSON and contain the following fields.

| Field | Meaning |
|---|---|
| `version` | Protocol version; the implementation uses version 1 |
| `request_id` | Identifier for an individual request |
| `node_id` | Identifier supplied by the sending node |
| `text` | Message text |

The receiver validates incoming messages before accepting them. The documented validation rules include the protocol version, identifiers, non-empty text, valid UTF-8, JSON structure and message-size limit.

The encoded message limit is documented as 4096 bytes. The submitted implementation should be checked to confirm that this limit is applied consistently to both transports.

A `node_id` identifies the sender in the application protocol, but it does not prove the sender's identity because authentication is not implemented.

### TCP framing

TCP provides an ordered byte stream rather than separate application-message boundaries. Therefore, each message is sent with a four-byte unsigned, big-endian length prefix followed by the JSON payload.

The receiver reads the header and then the declared payload. It must not assume that one network read will return a complete message, because TCP data may arrive in smaller pieces.

The receiver validates the declared length before processing the payload. Zero-length or oversized frames are rejected according to the protocol rules. Handling partial reads and incomplete frames is important because it prevents fragments of a message from being treated as complete input.

### UDP transport

UDP sends each encoded JSON message as an individual datagram. Unlike TCP, it preserves datagram boundaries and does not require the same length-prefix framing.

The receiver validates incoming datagrams and sends an acknowledgement when a valid message is accepted. However, UDP does not guarantee delivery, ordering or retransmission. Datagrams may be lost or duplicated.

A sender may encounter a timeout or network error if a response does not arrive. Such an outcome does not, by itself, prove whether the receiver processed the original datagram. Automatic retries should not be claimed unless they are supported by the implementation.

### Acknowledgement

A successful acknowledgement uses the following fields.

| Field | Meaning |
|---|---|
| `version` | Protocol version |
| `request_id` | Identifier of the request being acknowledged |
| `receiver` | Identifier of the receiving node |
| `accepted` | Indicates whether the message was accepted |

For TCP, the sender checks the acknowledgement version, request identifier, acceptance flag and receiver field, as described by the protocol implementation.

An acknowledgement indicates that the receiver accepted the message and returned a response. It does not guarantee permanent storage. Messages are held in memory and may be lost when the node process stops.

### HTTP inbox endpoint

`GET /v1/messages` returns the messages currently held in the node's inbox as a JSON array. The documented inbox retains a bounded recent history of up to 20 messages.

Request:

```http
GET /v1/messages HTTP/1.1
```

Illustrative response:

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "version": 1,
    "request_id": "tcp-1",
    "node_id": "node-a",
    "text": "hello"
  }
]
```

The response illustrates the message structure and is not a recorded test transcript. Only GET is intended for this endpoint; unsupported methods should receive HTTP 405.

### Demonstration ports

The supplied run scripts use the following loopback addresses.

| Node | HTTP | TCP | UDP | gRPC |
|---|---|---|---|---|
| Node A | `127.0.0.1:8081` | `127.0.0.1:8181` | `127.0.0.1:8281` | `127.0.0.1:9091` |
| Node B | `127.0.0.1:8082` | `127.0.0.1:8182` | `127.0.0.1:8282` | `127.0.0.1:9092` |

These ports are for the local demonstration configuration. They do not imply that the nodes are publicly accessible.

### TCP and UDP comparison

TCP provides reliable, ordered byte-stream delivery while the connection remains usable. The application adds explicit framing so that the receiver can identify individual messages. TCP is useful when ordered delivery matters, although connection handling and stream processing introduce additional complexity.

UDP preserves datagram boundaries and has lower transport overhead, but delivery is not guaranteed and packets may be lost, duplicated or reordered. Acknowledgements and request identifiers help the application observe communication outcomes, but they do not make UDP inherently reliable or provide authentication.

Using both transports demonstrates how the same application-level messaging task can be implemented over different transport protocols. Using TCP alone would simplify the design but would not demonstrate the differences between stream-based and datagram-based communication.

### Design rationale and limitations

The bounded in-memory inbox provides a simple way to inspect accepted messages through an HTTP endpoint. It avoids the additional complexity of persistent storage, but messages are not retained after the process stops.

The protocol validates incoming data before accepting messages. TCP framing handles the fact that a stream does not preserve application-message boundaries, while UDP processes individual datagrams. These approaches address different transport requirements.

The prototype is designed for local coursework testing. Node IDs and request IDs support identification and correlation within the protocol but do not authenticate the sender. The implementation should not be described as secure communication over an untrusted network.

### Verification evidence

The Week 4 evidence records that the targeted IPC tests passed. The recorded test coverage includes TCP and UDP round trips, TCP frame fragmentation and size limits, invalid messages, cancelled clients, oversized TCP frames, oversized UDP datagrams and incomplete TCP frame timeouts.

The evidence also records successful completion of the full Go test suite, the race-enabled test suite and `go vet ./...`.

Manual client commands were recorded for sending TCP and UDP messages in both directions between Node A and Node B. The commands targeted Node A's TCP and UDP ports `8181` and `8281`, and Node B's ports `8182` and `8282`. The evidence also records queries to both nodes' `/v1/messages` endpoints.

These statements summarise the recorded evidence. They do not claim that a stopped-node failure scenario was manually reproduced, because that scenario is not established by the cited Week 4 command record. The test results apply to the cases actually executed, not every possible network condition.

---

## Week 5 — gRPC Remote File Metadata

### Purpose

The node provides a gRPC service that allows a client to request file metadata from a remote node. Protocol Buffers define the request and response structures, while gRPC provides the remote procedure call between the client and server.

The service returns metadata rather than transferring file contents.

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

Each `FileMetadata` record contains the following fields.

| Field | Type | Meaning |
|---|---|---|
| `name` | string | File name |
| `size` | int64 | File size |
| `sha256` | string | SHA-256 digest |
| `modified_unix` | int64 | Modification time represented as a Unix timestamp |

The Protocol Buffers definition establishes the structure of the RPC request and response. The generated Go code is used by the client and server to exchange these messages.

### Request and response flow

The client calls `ListFiles` with a `ListFilesRequest` containing the requester's ID. The server responds with the serving node's identity and a list of file metadata.

The HTTP endpoint `GET /v1/remote-files?peer_id=node-b` provides a REST interface to this operation. The local node uses its peer registry to identify the remote peer and requests metadata through gRPC. The HTTP handler then returns the result to the HTTP client as JSON.

The HTTP response is the application's REST representation of the result. It should not be confused with the original Protocol Buffers message transmitted over gRPC.

### Error handling and limitations

A remote call can fail if the server is unavailable, the request is cancelled or its deadline expires. The client and HTTP handler must handle these outcomes instead of assuming that every request will succeed.

An empty file list can be a valid response. It should not automatically be treated as an error merely because no metadata is returned.

The `requester_id` field is an unverified claim and does not establish the requester's identity. Similarly, the returned `node_id` is not proof of the server's identity. Authentication requires additional mechanisms that are outside this prototype's demonstrated capabilities.

The HTTP remote-file endpoint translates remote-call outcomes into HTTP responses. The project evidence describes HTTP 503 for an unavailable remote service and HTTP 504 for a deadline or timeout. Other mappings, including cancellation and invalid requests, should be documented according to the submitted handler and its tests.

The prototype uses local coursework configuration and does not provide TLS or authentication. It should not be presented as secure communication for an untrusted network.

### Design rationale and limitations

Protocol Buffers provide a defined service contract, and gRPC supports remote invocation of the metadata operation. The HTTP adapter allows clients to request remote metadata through the node's REST interface.

An alternative would be to obtain metadata through a remote HTTP endpoint. The gRPC approach demonstrates the use of a separate remote procedure call interface alongside REST and TCP/UDP messaging.

The service returns metadata only; it does not transfer file contents. Peer registration is manual, and the local registry determines which peers can be addressed through the HTTP adapter. This limits automatic discovery but keeps the prototype's behaviour straightforward to demonstrate and test.

### Verification

The Week 5 evidence records successful completion of the following checks:

- `go build ./...`
- `go test ./... -count=1 -timeout 30s`
- `go test -race ./... -count=1 -timeout 60s`
- `go vet ./...`
- `git diff --check`

The recorded manual checks showed both nodes' `/healthz` endpoints returning HTTP 200. Node B's `/v1/files` endpoint returned HTTP 200 and listed `diagram.txt` and `notes.txt`.

Registering Node B in Node A returned HTTP 201. The recorded request to `/v1/remote-files?peer_id=node-b` returned HTTP 200, with `node_id` set to `node-b` and metadata for both files, including sizes, SHA-256 hashes and modification timestamps.

These results support the successful remote metadata-listing path. They do not, by themselves, prove that timeout, cancellation, empty-result and every request-validation scenario was manually demonstrated. Such cases should be described as automated-test coverage only where the actual tests confirm that coverage.

---

## Cross-cutting limitations

The prototype demonstrates several communication approaches, but it has limitations that should be considered when interpreting its results.

- The supplied configuration uses loopback addresses for local coursework testing.
- Peer registration is manual and stored in each node's local registry.
- The file catalogue stores metadata rather than file contents.
- The message inbox is in memory and does not retain messages after a process restart.
- UDP does not guarantee delivery or ordering.
- Request IDs and node IDs support correlation but do not authenticate a peer.
- The prototype does not provide TLS or authentication.
- Test results apply to the commands and cases actually run and do not prove that every possible network failure has been handled.

These limitations define the scope of the coursework prototype and identify areas that could be improved in a more complete distributed application.

## Sources and evidence

This document describes the REST, TCP, UDP and gRPC interfaces implemented in the coursework project. Verification summaries are based on the corresponding weekly evidence records and the commands and results recorded there.

The project source, Protocol Buffers definition and supplied run scripts should be used to confirm implementation-specific details, including request validation, size limits, error mappings and port configuration.

ChatGPT assisted with reviewing and improving the documentation wording. The student remains responsible for checking the document against the submitted implementation, understanding the design decisions and following the module's AI-use disclosure requirements.
