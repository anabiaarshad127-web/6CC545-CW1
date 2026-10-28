# Week 1 Design Note

## Purpose

The aim of Week 1 is to create the basic node for the distributed application. The node provides a health endpoint so that its availability can be checked.

## Design

The node uses Go and provides an HTTP `/healthz` endpoint. The endpoint allows a client to check whether the node is responding. The node address is configurable, which allows different node instances to use different addresses and ports during local testing.

The application also supports controlled shutdown so that its lifecycle can be managed when it is stopped. A shutdown timeout helps prevent the application from waiting indefinitely during shutdown.

## Limitations

At this stage, the health endpoint only checks the node's basic availability. It does not prove that peer communication, remote file access or the complete distributed application is working. These features require additional implementation and testing in later weeks.
