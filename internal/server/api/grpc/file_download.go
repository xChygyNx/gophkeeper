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

// FileDownload - checks the validity of the token, save record, upload file on client
func (h *Handler) FileDownload(ctx context.Context, req *grpc.DownloadBinaryRequest) (*grpc.DownloadBinaryResponse, error) {
	h.logger.Info("file download")

	logError, statusErr := validateAccessToken(req.AccessToken.Token, h.token)
	if statusErr != nil {
		h.logger.Error(logError)
		return &grpc.DownloadBinaryResponse{}, statusErr
	}

	FileData := &model.FileRequest{}
	FileData.UserID = req.AccessToken.UserId
	FileData.Name = req.Name

	exists, err := h.file.FileExists(FileData)
	if err != nil {
		finalError := fmt.Errorf("error in check file exists: %w", err)
		h.logger.Error(finalError)
		return &grpc.DownloadBinaryResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}
	if !exists {
		err = errors.ErrFileNotExists
		h.logger.Error(err)
		return &grpc.DownloadBinaryResponse{}, status.Errorf(codes.NotFound, err.Error())
	}

	data, err := service.DownloadFile(h.config.FileFolder, req.AccessToken.UserId, req.Name)
	if err != nil {
		finalError := fmt.Errorf("error in download file from keeper: %w", err)
		h.logger.Error(finalError)
		return &grpc.DownloadBinaryResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}

	return &grpc.DownloadBinaryResponse{Data: data}, nil
}
