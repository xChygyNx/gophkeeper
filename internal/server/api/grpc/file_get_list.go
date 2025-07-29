package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xChygyNx/gophkeeper/internal/server/model"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/errors"
)

// FileGetList - checks the validity of tokens and get list records file
func (h *Handler) FileGetList(ctx context.Context, req *grpc.GetListBinaryRequest) (*grpc.GetListBinaryResponse, error) {
	h.logger.Info("file get list")

	endDateToken, err := h.token.GetEndDateToken(req.AccessToken.Token)
	if err != nil {
		finalError := fmt.Errorf("error in get end date of token from DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.GetListBinaryResponse{}, status.Errorf(
			codes.Internal, finalError.Error())
	}
	valid := h.token.Validate(endDateToken)
	if !valid {
		h.logger.Error(errors.ErrNotValidateToken)
		return &grpc.GetListBinaryResponse{}, status.Errorf(
			codes.Unauthenticated, errors.ErrNotValidateToken.Error(),
		)
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
