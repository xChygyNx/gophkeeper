package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xChygyNx/gophkeeper/internal/server/model"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/errors"
)

// FileRemove - checks the validity of the token, delete record, remove file on server
func (h *Handler) FileRemove(ctx context.Context, req *grpc.DeleteBinaryRequest) (*grpc.DeleteBinaryResponse, error) {
	h.logger.Info("file remove")

	endDateToken, err := h.token.GetEndDateToken(req.AccessToken.Token)
	if err != nil {
		finalError := fmt.Errorf("error in get end date of token from DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.DeleteBinaryResponse{}, status.Errorf(
			codes.Internal, finalError.Error())
	}
	valid := h.token.Validate(endDateToken)
	if !valid {
		h.logger.Error(errors.ErrNotValidateToken)
		return &grpc.DeleteBinaryResponse{}, status.Errorf(
			codes.Unauthenticated, errors.ErrNotValidateToken.Error(),
		)
	}

	FileData := &model.FileRequest{}
	FileData.UserID = req.AccessToken.UserId
	FileData.Name = req.Name

	BinaryId, err := h.file.DeleteFile(FileData)
	if err != nil {
		finalError := fmt.Errorf("error in delete file from DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.DeleteBinaryResponse{}, status.Errorf(
			codes.Internal, finalError.Error(),
		)
	}

	err = service.RemoveFile(h.config.FileFolder, req.AccessToken.UserId, req.Name)
	if err != nil {
		finalError := fmt.Errorf("error in delete file from keeper: %w", err)
		h.logger.Error(finalError)
		return &grpc.DeleteBinaryResponse{}, status.Errorf(
			codes.Internal, finalError.Error(),
		)
	}

	return &grpc.DeleteBinaryResponse{Id: BinaryId}, nil
}
