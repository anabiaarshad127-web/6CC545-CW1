package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"example.edu/6cc545/cw1/internal/config"
	"example.edu/6cc545/cw1/internal/httpapi"
	"example.edu/6cc545/cw1/internal/model"
	"example.edu/6cc545/cw1/internal/state"
)

const maxSeedSize = 1 << 20 // 1 MiB

func CheckStage(c config.Config) error {
	if c.TCP != "" || c.UDP != "" || c.GRPC != "" {
		return errors.New("Week 3: TCP, UDP and gRPC are not enabled yet")
	}
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

	return func() {}, nil
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
