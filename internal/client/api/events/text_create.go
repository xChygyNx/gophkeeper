package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/vars"
)

// TextCreate - add text
func (s Event) TextCreate(name, description, password, plaintext string, token model.Token) error {
	s.logger.Info("text create")

	secretKey := encryption.AesKeySecureRandom([]byte(password))
	encryptText, err := encryption.Encrypt(plaintext, secretKey)
	if err != nil {
		myErr := fmt.Errorf("error in encrypt text: %w", err)
		s.logger.Error(myErr)
		return myErr
	}
	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	metadata := model.MetadataEntity{Name: name, Description: description, Type: vars.Text.ToString()}
	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	createdEntityID, err := s.grpc.EntityCreate(context.Background(),
		&grpc.CreateEntityRequest{Data: []byte(encryptText), Metadata: string(jsonMetadata),
			AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID, CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on create text: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("text with ID %s is created", createdEntityID)
	return nil
}
