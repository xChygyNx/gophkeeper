package events

import (
	"context"
	"fmt"

	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/vars"
)

// TextUpdate - update text
func (s Event) TextUpdate(name, passwordSecure, text string, token model.Token) error {
	s.logger.Info("text update")

	secretKey := encryption.AesKeySecureRandom([]byte(passwordSecure))
	encryptText, err := encryption.Encrypt(text, secretKey)
	if err != nil {
		if err != nil {
			myErr := fmt.Errorf("error in encrypt text: %w", err)
			s.logger.Error(myErr)
			return myErr
		}
	}

	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	updatedTextEntityID, err := s.grpc.EntityUpdate(context.Background(),
		&grpc.UpdateEntityRequest{Name: name, Data: []byte(encryptText), Type: vars.Text.ToString(),
			AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID, CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on update text: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("text with ID %s updated", updatedTextEntityID)
	return nil
}
