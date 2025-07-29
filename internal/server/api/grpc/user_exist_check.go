package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
)

// UserExist - user existence check
func (h *Handler) UserExist(ctx context.Context, req *grpc.UserExistRequest) (*grpc.UserExistResponse, error) {
	h.logger.Info("user exist check")
	exist, err := h.user.UserExists(req.Username)
	if err != nil {
		finalError := fmt.Errorf("error in check user exists: %w", err)
		h.logger.Error(finalError)
		return &grpc.UserExistResponse{Exist: false}, status.Errorf(
			codes.Internal, finalError.Error(),
		)
	}
	h.logger.Debug("User %s is exists: %b", req.Username, exist)
	return &grpc.UserExistResponse{Exist: exist}, nil
}
