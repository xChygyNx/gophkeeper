package events

import (
	"context"
	"fmt"
	grpc "github.com/xChygyNx/gophkeeper/internal/proto"

	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/vars"
)

// TextDelete - delete text
func (s Event) TextDelete(text []string, token model.Token) error {
	s.logger.Info("text delete")

	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)

	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	deletedTextEntityID, err := s.grpc.EntityDelete(context.Background(),
		&grpc.DeleteEntityRequest{Name: text[0], Type: vars.Text.ToString(),
			AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID,
				CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request to text delete: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("Text with ID %s is deleted", deletedTextEntityID)
	return nil
}
