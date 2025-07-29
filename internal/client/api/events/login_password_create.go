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

// LoginPasswordCreate - add login-password
func (s Event) LoginPasswordCreate(name, description, passwordSecure, login, password string, token model.Token) error {
	s.logger.Info("login password create")

	loginPassword := model.LoginPassword{Login: login, Password: password}
	jsonLoginPassword, err := json.Marshal(loginPassword)
	if err != nil {
		myErr := fmt.Errorf("error in marshal login/password data: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	secretKey := encryption.AesKeySecureRandom([]byte(passwordSecure))
	encryptLoginPassword, err := encryption.Encrypt(string(jsonLoginPassword), secretKey)
	if err != nil {
		myErr := fmt.Errorf("error in decrypt login/password data: %w", err)
		s.logger.Error(myErr)
		return myErr
	}
	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	metadata := model.MetadataEntity{Name: name, Description: description, Type: vars.LoginPassword.ToString()}
	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("error in marshal login/password metadata: %w", err)
	}
	createdEntityID, err := s.grpc.EntityCreate(context.Background(),
		&grpc.CreateEntityRequest{Data: []byte(encryptLoginPassword), Metadata: string(jsonMetadata),
			AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID, CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on create login/password for user: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("created login/password with ID: %s", createdEntityID)
	return nil
}
