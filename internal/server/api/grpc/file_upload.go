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

// FileUpload - checks the validity of the token, upload file on client
func (h *Handler) FileUpload(ctx context.Context, req *grpc.UploadBinaryRequest) (*grpc.UploadBinaryResponse, error) {
	h.logger.Info("file upload")

	endDateToken, err := h.token.GetEndDateToken(req.AccessToken.Token)
	if err != nil {
		finalError := fmt.Errorf("error in get end date of token from DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.UploadBinaryResponse{}, status.Errorf(
			codes.Internal, finalError.Error())
	}
	valid := h.token.Validate(endDateToken)
	if !valid {
		h.logger.Error(errors.ErrNotValidateToken)
		return &grpc.UploadBinaryResponse{}, status.Errorf(
			codes.Unauthenticated, errors.ErrNotValidateToken.Error(),
		)
	}

	FileData := &model.FileRequest{}
	FileData.UserID = req.AccessToken.UserId
	FileData.Name = req.Name

	exists, err := h.file.FileExists(FileData)
	if err != nil {
		finalError := fmt.Errorf("error in check file exists: %w", err)
		h.logger.Error(finalError)
		return &grpc.UploadBinaryResponse{}, status.Errorf(
			codes.Internal, finalError.Error(),
		)
	}
	if exists {
		err = errors.ErrNameAlreadyExists
		h.logger.Error(err)
		return &grpc.UploadBinaryResponse{}, status.Errorf(
			codes.AlreadyExists, err.Error(),
		)
	}

	UploadFile, err := h.file.UploadFile(FileData)
	if err != nil {
		finalError := fmt.Errorf("error in upload file in DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.UploadBinaryResponse{}, status.Errorf(
			codes.Internal, finalError.Error(),
		)
	}

	err = service.UploadFile(h.config.FileFolder, req.AccessToken.UserId, req.Name, req.Data)
	if err != nil {
		finalError := fmt.Errorf("error in upload file in keeper: %w", err)
		h.logger.Error(finalError)
		return &grpc.UploadBinaryResponse{}, status.Errorf(
			codes.Internal, finalError.Error(),
		)
	}

	h.logger.Debug(UploadFile.Name)
	return &grpc.UploadBinaryResponse{Name: UploadFile.Name}, nil
}
