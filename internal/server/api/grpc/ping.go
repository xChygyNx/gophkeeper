package grpchandler

import (
	"context"
	"fmt"
	grpc "github.com/xChygyNx/gophkeeper/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Ping - checks the database connection
func (h *Handler) Ping(ctx context.Context, req *grpc.PingRequest) (*grpc.PingResponse, error) {
	h.logger.Info("ping")
	var msg string
	err := h.database.Ping()
	if err != nil {
		finalError := fmt.Errorf("error in DB connection: %w", err)
		h.logger.Error(finalError)
		return &grpc.PingResponse{Message: finalError.Error()}, status.Errorf(
			codes.Internal, finalError.Error(),
		)
	}
	msg = "successful database connection"
	h.logger.Info(msg)
	return &grpc.PingResponse{Message: msg}, nil
}
