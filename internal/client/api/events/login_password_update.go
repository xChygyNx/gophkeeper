package events

import (
	"context"
	"encoding/json"
	"fmt"
	grpc "github.com/xChygyNx/gophkeeper/internal/proto"

	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/vars"
)

// LoginPasswordUpdate - update login-password
func (s Event) LoginPasswordUpdate(name, passwordSecure, login, password string, token model.Token) error {
	s.logger.Info("login password update")

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
		myErr := fmt.Errorf("error in encrypt login/password data: %w", err)
		s.logger.Error(myErr)
		return myErr
	}
	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	updatedLoginPasswordEntityID, err := s.grpc.EntityUpdate(context.Background(),
		&grpc.UpdateEntityRequest{Name: name, Data: []byte(encryptLoginPassword), Type: vars.LoginPassword.ToString(),
			AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID, CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on update login/password: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("login/password with ID %s updated", updatedLoginPasswordEntityID)
	return nil
}
