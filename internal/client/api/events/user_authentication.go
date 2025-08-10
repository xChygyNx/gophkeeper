package events

import (
	"fmt"
	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	grpc "github.com/xChygyNx/gophkeeper/internal/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
)

// Authentication - auth user and return token
func (s Event) Authentication(username, password string) (model.Token, error) {
	s.logger.Info("authentication")

	token := model.Token{}
	password, err := encryption.HashPassword(password)
	if err != nil {
		myErr := fmt.Errorf("error in hash password: %w", err)
		s.logger.Error(myErr)
		return token, myErr
	}
	authenticatedUser, err := s.grpc.Authentication(s.context, &grpc.AuthenticationRequest{Username: username, Password: password})
	if err != nil {
		myErr := fmt.Errorf("error in user authenticate: %w", err)
		s.logger.Error(myErr)
		return token, myErr
	}

	createdToken, err := service.ConvertTimestampToTime(authenticatedUser.AccessToken.CreatedAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		s.logger.Error(myErr)
		return token, myErr
	}
	endDateToken, err := service.ConvertTimestampToTime(authenticatedUser.AccessToken.EndDateAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		s.logger.Error(myErr)
		return token, myErr
	}
	token = model.Token{AccessToken: authenticatedUser.AccessToken.Token, UserID: authenticatedUser.AccessToken.UserId,
		CreatedAt: createdToken, EndDateAt: endDateToken}

	err = service.CreateStorageNotExistsUser(s.config.FileFolder, token.UserID)
	if err != nil {
		myErr := fmt.Errorf("error in create storage for new user: %w", err)
		s.logger.Error(myErr)
		return token, myErr
	}

	return token, nil
}
