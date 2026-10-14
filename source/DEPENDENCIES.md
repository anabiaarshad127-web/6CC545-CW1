# Week 5 Dependencies and Code Generation

The Week 5 implementation uses `google.golang.org/grpc` and `google.golang.org/protobuf`. The selected dependency versions are recorded in `go.mod` and `go.sum`.

The Protocol Buffers service definition is stored in `api/files/v1/files.proto`. Its generated Go bindings are stored alongside the service definition. The generation scripts are `scripts/generate.sh` and `scripts/generate.ps1`.

The development environment used for this checkpoint was Go `go1.27.1 darwin/arm64` on macOS, with `protoc` version `36.2`.

The following verification commands passed:

```bash
go build ./...
go test ./... -count=1 -timeout 30s
go test -race ./... -count=1 -timeout 60s
go vet ./...
```

The application uses loopback addresses and insecure gRPC transport for this local coursework demonstration. It does not provide TLS or authentication and should not be exposed to an untrusted network.

The dependency versions are compatibility choices for this project, not a claim that they are the newest releases or approved for production use. Keep generated bindings, `go.mod`, and `go.sum` under version control. Do not commit tool binaries or dependency caches.