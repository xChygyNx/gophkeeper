package events

import (
	"context"
	"fmt"

	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
)

// FileDownload - download file
func (s Event) FileDownload(name string, password string, token model.Token) error {
	s.logger.Info("file download")

	secretKey := encryption.AesKeySecureRandom([]byte(password))
	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	downloadFile, err := s.grpc.FileDownload(context.Background(),
		&grpc.DownloadBinaryRequest{Name: name, AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID,
			CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on file download: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	file, err := encryption.Decrypt(string(downloadFile.Data), secretKey)
	if err != nil {
		myErr := fmt.Errorf("error in decrypt file data: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	err = service.UploadFile(s.config.FileFolder, token.UserID, name, []byte(file))
	if err != nil {
		myErr := fmt.Errorf("error in upload file: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("file %s is downloaded", name)
	return nil
}
