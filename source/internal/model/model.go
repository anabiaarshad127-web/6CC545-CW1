package model

import (
	"encoding/hex"
	"fmt"
	"strings"

	"example.edu/6cc545/cw1/internal/config"
)

type Peer struct {
	ID          string `json:"id"`
	HTTPAddress string `json:"http_address"`
	GRPCAddress string `json:"grpc_address"`
	Status      string `json:"status"`
}

type FileMetadata struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	SHA256       string `json:"sha256"`
	ModifiedUnix int64  `json:"modified_unix"`
}

func (p Peer) Validate() error {
	if !config.ValidID(p.ID) {
		return fmt.Errorf("invalid peer ID: %q", p.ID)
	}
	if !config.ValidAddress(p.HTTPAddress) {
		return fmt.Errorf("invalid peer HTTP address: %q", p.HTTPAddress)
	}
	if !config.ValidAddress(p.GRPCAddress) {
		return fmt.Errorf("invalid peer gRPC address: %q", p.GRPCAddress)
	}

	switch p.Status {
	case "unknown", "online", "offline":
		return nil
	default:
		return fmt.Errorf("invalid peer status: %q", p.Status)
	}
}

func (f FileMetadata) Validate() error {
	if f.Name == "" || len(f.Name) > 128 {
		return fmt.Errorf("file name must contain 1 to 128 bytes")
	}
	if f.Name == "." || f.Name == ".." {
		return fmt.Errorf("invalid file name: %q", f.Name)
	}
	if strings.ContainsAny(f.Name, "/\\\x00") {
		return fmt.Errorf("file name must not contain path separators or null bytes")
	}
	if f.Size < 0 {
		return fmt.Errorf("file size cannot be negative")
	}
	if f.ModifiedUnix < 0 {
		return fmt.Errorf("modified timestamp cannot be negative")
	}
	if len(f.SHA256) != 64 {
		return fmt.Errorf("SHA-256 must contain exactly 64 hexadecimal characters")
	}
	if _, err := hex.DecodeString(f.SHA256); err != nil {
		return fmt.Errorf("invalid SHA-256: %w", err)
	}
	return nil
}
