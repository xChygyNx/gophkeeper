package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xChygyNx/gophkeeper/internal/server/model"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
)

// FileGetList - checks the validity of tokens and get list records file
func (h *Handler) FileGetList(ctx context.Context, req *grpc.GetListBinaryRequest) (*grpc.GetListBinaryResponse, error) {
	h.logger.Info("file get list")

	logError, statusErr := validateAccessToken(req.AccessToken.Token, h.token)
	if statusErr != nil {
		h.logger.Error(logError)
		return &grpc.GetListBinaryResponse{}, statusErr
	}

	ListFile, err := h.file.GetListFile(req.AccessToken.UserId)
	if err != nil {
		finalError := fmt.Errorf("error in get list of files: %w", err)
		h.logger.Error(finalError)
		return &grpc.GetListBinaryResponse{}, status.Errorf(
			codes.Internal, finalError.Error(),
		)
	}
	list := model.GetListFile(ListFile)

	h.logger.Debug("Get file list: %s", ListFile)
	return &grpc.GetListBinaryResponse{Node: list}, nil
}
