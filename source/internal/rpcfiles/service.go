package rpcfiles

import (
	"context"

	pb "example.edu/6cc545/cw1/api/files/v1"
	"example.edu/6cc545/cw1/internal/config"
	"example.edu/6cc545/cw1/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Catalog interface {
	Snapshot() []model.FileMetadata
}

type Service struct {
	pb.UnimplementedFileServiceServer
	NodeID  string
	Catalog Catalog
}

func (s *Service) ListFiles(ctx context.Context, r *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	// Stop immediately if the caller has cancelled the request.
	if err := ctx.Err(); err != nil {
		return nil, status.FromContextError(err).Err()
	}

	// A missing request or invalid requester ID is not accepted.
	if r == nil || !config.ValidID(r.GetRequesterId()) {
		return nil, status.Error(codes.InvalidArgument, "invalid requester_id")
	}

	// The service must use the node's existing catalogue.
	if s == nil || s.Catalog == nil || !config.ValidID(s.NodeID) {
		return nil, status.Error(codes.Internal, "file service is not configured")
	}

	// Snapshot returns a safe copy of the catalogue.
	metadata := s.Catalog.Snapshot()
	files := make([]*pb.FileMetadata, 0, len(metadata))

	for _, file := range metadata {
		files = append(files, &pb.FileMetadata{
			Name:         file.Name,
			Size:         file.Size,
			Sha256:       file.SHA256,
			ModifiedUnix: file.ModifiedUnix,
		})
	}

	return &pb.ListFilesResponse{
		Files:  files,
		NodeId: s.NodeID,
	}, nil
}
