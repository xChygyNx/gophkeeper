package events

import (
	"context"
	"fmt"
	grpc "github.com/xChygyNx/gophkeeper/internal/proto"

	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/vars"
)

// LoginPasswordDelete - delete login-password
func (s Event) LoginPasswordDelete(loginPassword []string, token model.Token) error {
	s.logger.Info("login password delete")

	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	deletedLoginPasswordEntityID, err := s.grpc.EntityDelete(context.Background(),
		&grpc.DeleteEntityRequest{Name: loginPassword[0], Type: vars.LoginPassword.ToString(),
			AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID,
				CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request to login/password delete: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("Login/password with ID %s deleted", deletedLoginPasswordEntityID)
	return nil
}
