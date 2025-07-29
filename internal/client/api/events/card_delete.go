package events

import (
	"context"
	"fmt"

	"github.com/xChygyNx/gophkeeper/internal/client/model"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/vars"
)

// CardDelete -  delete card
func (s Event) CardDelete(card []string, token model.Token) error {
	s.logger.Info("card delete")

	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	deletedCardEntityID, err := s.grpc.EntityDelete(context.Background(),
		&grpc.DeleteEntityRequest{Name: card[0], Type: vars.Card.ToString(),
			AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID,
				CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request to card delete: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("Card with ID %s deleted", deletedCardEntityID)
	return nil
}
