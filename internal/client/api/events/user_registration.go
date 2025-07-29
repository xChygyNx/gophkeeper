package events

import (
	"fmt"
	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
)

// Registration - user registration
func (s Event) Registration(username, password string) (model.Token, error) {
	s.logger.Info("registration")

	token := model.Token{}
	password, err := encryption.HashPassword(password)
	if err != nil {
		myErr := fmt.Errorf("error in hash password: %w", err)
		s.logger.Error(myErr)
		return token, myErr
	}
	registeredUser, err := s.grpc.Registration(s.context, &grpc.RegistrationRequest{Username: username, Password: password})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on user registration: %w", err)
		s.logger.Error(myErr)
		return token, myErr
	}
	createdToken, err := service.ConvertTimestampToTime(registeredUser.AccessToken.CreatedAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		s.logger.Error(myErr)
		return token, myErr
	}
	endDateToken, err := service.ConvertTimestampToTime(registeredUser.AccessToken.EndDateAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		s.logger.Error(myErr)
		return token, myErr
	}
	token = model.Token{AccessToken: registeredUser.AccessToken.Token, UserID: registeredUser.AccessToken.UserId,
		CreatedAt: createdToken, EndDateAt: endDateToken}

	err = service.CreateStorageUser(s.config.FileFolder, token.UserID)
	if err != nil {
		myErr := fmt.Errorf("error in create user storage: %w", err)
		s.logger.Error(myErr)
		return token, myErr
	}

	return token, nil
}
