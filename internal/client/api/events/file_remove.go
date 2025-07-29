package events

import (
	"context"
	"fmt"

	"github.com/xChygyNx/gophkeeper/internal/client/model"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
)

// FileRemove - delete file
func (s Event) FileRemove(binary []string, token model.Token) error {
	s.logger.Info("file remove")

	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)
	deletedCard, err := s.grpc.FileRemove(context.Background(),
		&grpc.DeleteBinaryRequest{Name: binary[0], AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID,
			CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on delete file: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("file %s is deleted", deletedCard)
	return nil
}
