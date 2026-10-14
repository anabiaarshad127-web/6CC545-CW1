package rpcfiles

import (
	"context"
	"time"

	pb "example.edu/6cc545/cw1/api/files/v1"
	"example.edu/6cc545/cw1/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func Fetch(ctx context.Context, address, requester string, timeout time.Duration) (*pb.ListFilesResponse, error) {
	if ctx == nil {
		return nil, status.Error(codes.InvalidArgument, "missing context")
	}
	if !config.ValidAddress(address) {
		return nil, status.Error(codes.InvalidArgument, "invalid remote address")
	}
	if !config.ValidID(requester) {
		return nil, status.Error(codes.InvalidArgument, "invalid requester_id")
	}
	if timeout <= 0 {
		return nil, status.Error(codes.InvalidArgument, "timeout must be positive")
	}
	if err := ctx.Err(); err != nil {
		return nil, status.FromContextError(err).Err()
	}

	opCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(4096),
			grpc.MaxCallRecvMsgSize(65536),
		),
	)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "create gRPC client: %v", err)
	}
	defer conn.Close()

	response, err := pb.NewFileServiceClient(conn).ListFiles(
		opCtx,
		&pb.ListFilesRequest{RequesterId: requester},
	)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, status.Error(codes.Internal, "remote service returned an empty response")
	}

	return response, nil
}
