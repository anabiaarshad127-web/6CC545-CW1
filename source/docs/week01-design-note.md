# Week 1 Design Note

## Purpose

Week 1 establishes the basic node for the distributed application. Each node has its own identity and configuration. The `/healthz` endpoint allows a client to check whether the node is responding.

## Design decisions

The node ID is separate from its HTTP address because identity describes the participant, while the address tells a client where to connect. Configuration is validated before the server starts so invalid settings can be rejected early. The application manages the HTTP listener and uses a bounded shutdown timeout to help the node stop cleanly.

The health endpoint returns the node ID and an `ok` status for GET requests. However, a successful response only shows that the node is responding. It does not prove that another node is reachable, remote files can be accessed, or the complete distributed application is working.