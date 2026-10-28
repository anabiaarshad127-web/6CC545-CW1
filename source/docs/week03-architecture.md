# Week 3 Architecture

## Overview

Week 3 extends the basic node with a peer registry and a file metadata catalogue. The registry records known peers, while the catalogue stores metadata about files available to the node.

## Peer Registry

The peer registry stores peer details in memory. Before a peer is added, its details are validated. Duplicate peer IDs are rejected, and the HTTP API returns `409 Conflict` for a duplicate registration.

A read-write mutex protects the registry's shared state. The registry returns a sorted snapshot of its peers rather than exposing its internal map directly.

## File Catalogue

The file catalogue stores file metadata, not file contents. It validates candidate metadata before replacing the stored catalogue and rejects duplicate filenames. Catalogue results are returned in filename order.

The node supports a `-seed` option for loading supplied synthetic file metadata. This provides test data for the catalogue; it does not demonstrate that the node scans a live filesystem or transfers the listed files.

## Limitations

The registry and catalogue are in-memory components, so their state is not persistent across process restarts. Peer registration is manual, and this checkpoint does not provide automatic peer discovery. Listing file metadata also does not prove that the corresponding file contents can be retrieved from a remote peer.
