package events

import (
	"context"
	"fmt"
	grpc "github.com/xChygyNx/gophkeeper/internal/proto"

	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
)

// FileUpload - upload file
func (s Event) FileUpload(name string, password string, file []byte, token model.Token) (string, error) {
	s.logger.Info("file upload")

	secretKey := encryption.AesKeySecureRandom([]byte(password))
	encryptFile, err := encryption.Encrypt(string(file), secretKey)
	if err != nil {
		if err != nil {
			myErr := fmt.Errorf("error in encrypt file data: %w", err)
			s.logger.Error(myErr)
			return "", myErr
		}
	}
	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)

	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	uploadFile, err := s.grpc.FileUpload(context.Background(),
		&grpc.UploadBinaryRequest{Name: name, Data: []byte(encryptFile),
			AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID,
				CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on file upload: %w", err)
		s.logger.Error(myErr)
		return "", myErr
	}
	s.logger.Debug("file %s is uploaded", uploadFile.Name)
	return uploadFile.Name, nil
}
