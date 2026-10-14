package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	filesv1 "example.edu/6cc545/cw1/api/files/v1"
	"example.edu/6cc545/cw1/internal/config"
	"example.edu/6cc545/cw1/internal/httpapi"
	"example.edu/6cc545/cw1/internal/ipc"
	"example.edu/6cc545/cw1/internal/model"
	"example.edu/6cc545/cw1/internal/rpcfiles"
	"example.edu/6cc545/cw1/internal/state"
	"google.golang.org/grpc"
)

const maxSeedSize = 1 << 20

func CheckStage(c config.Config) error {
	return nil
}

func Wire(ctx context.Context, c config.Config, mux *http.ServeMux) (func(), error) {
	status := state.NewStatus()
	mux.HandleFunc("/v1/status", httpapi.StatusResource(status))

	peers := state.NewRegistry()
	catalog := state.NewCatalog()

	if c.Seed != "" {
		data, err := readSeed(c.Seed)
		if err != nil {
			return nil, err
		}

		var files []model.FileMetadata
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&files); err != nil {
			return nil, fmt.Errorf("decode seed file %q: %w", c.Seed, err)
		}

		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			if err == nil {
				return nil, fmt.Errorf("seed file %q must contain one JSON value", c.Seed)
			}
			return nil, fmt.Errorf("read end of seed file %q: %w", c.Seed, err)
		}

		if err := catalog.Replace(files); err != nil {
			return nil, fmt.Errorf("validate seed file %q: %w", c.Seed, err)
		}
	}

	mux.HandleFunc("/v1/peers", httpapi.Peers(peers))
	mux.HandleFunc("/v1/files", httpapi.Files(catalog))
	mux.HandleFunc("/v1/remote-files", httpapi.RemoteFiles(peers, c.ID, c.Timeout))

	serviceCtx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	var closers []io.Closer
	var grpcServers []*grpc.Server
	var once sync.Once

	inbox := state.NewInbox()

	mux.HandleFunc("/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httpapi.ErrorJSON(w, http.StatusMethodNotAllowed, "method_not_allowed")
			return
		}
		httpapi.WriteJSON(w, http.StatusOK, inbox.Snapshot())
	})

	cleanup := func() {
		once.Do(func() {
			cancel()
			for _, server := range grpcServers {
				server.Stop()
			}
			for _, closer := range closers {
				_ = closer.Close()
			}
			wg.Wait()
		})
	}

	if c.TCP != "" {
		listener, err := net.Listen("tcp", c.TCP)
		if err != nil {
			cleanup()
			return nil, fmt.Errorf("listen on TCP %s: %w", c.TCP, err)
		}

		closers = append(closers, listener)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := ipc.ServeTCP(serviceCtx, listener, c.ID, c.Timeout, inbox.Record); err != nil && serviceCtx.Err() == nil {
				log.Printf("TCP service stopped: %v", err)
			}
		}()
	}

	if c.UDP != "" {
		address, err := net.ResolveUDPAddr("udp", c.UDP)
		if err != nil {
			cleanup()
			return nil, fmt.Errorf("resolve UDP address %s: %w", c.UDP, err)
		}

		conn, err := net.ListenUDP("udp", address)
		if err != nil {
			cleanup()
			return nil, fmt.Errorf("listen on UDP %s: %w", c.UDP, err)
		}

		closers = append(closers, conn)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := ipc.ServeUDP(serviceCtx, conn, c.ID, c.Timeout, inbox.Record); err != nil && serviceCtx.Err() == nil {
				log.Printf("UDP service stopped: %v", err)
			}
		}()
	}

	if c.GRPC != "" {
		listener, err := net.Listen("tcp", c.GRPC)
		if err != nil {
			cleanup()
			return nil, fmt.Errorf("listen on gRPC %s: %w", c.GRPC, err)
		}

		server := grpc.NewServer(
			grpc.MaxRecvMsgSize(4096),
			grpc.MaxSendMsgSize(65536),
		)
		filesv1.RegisterFileServiceServer(server, &rpcfiles.Service{
			NodeID:  c.ID,
			Catalog: catalog,
		})

		closers = append(closers, listener)
		grpcServers = append(grpcServers, server)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := server.Serve(listener); err != nil && serviceCtx.Err() == nil {
				log.Printf("gRPC service stopped: %v", err)
			}
		}()
	}

	return cleanup, nil
}

func readSeed(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open seed file %q: %w", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxSeedSize+1))
	if err != nil {
		return nil, fmt.Errorf("read seed file %q: %w", path, err)
	}
	if len(data) > maxSeedSize {
		return nil, fmt.Errorf("seed file %q exceeds the 1 MiB limit", path)
	}

	return data, nil
}
