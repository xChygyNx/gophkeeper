package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/xChygyNx/gophkeeper/internal/client/consts"
	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/vars"
)

// CardUpdate - update card
func (s Event) CardUpdate(name, passwordSecure, paymentSystem, number, holder, cvc, endDateCard string, token model.Token) error {
	s.logger.Info("card update")

	intCvc, err := strconv.Atoi(cvc)
	if err != nil {
		myErr := fmt.Errorf("invalid cvc code of card: %w", err)
		s.logger.Error(myErr)
		return myErr
	}
	timeEndDate, err := time.Parse(consts.DateFormat, endDateCard)
	if err != nil {
		myErr := fmt.Errorf("invalid end time of card: %w", err)
		s.logger.Error(myErr)
		return myErr
	}
	card := model.Card{PaymentSystem: paymentSystem, Number: number, Holder: holder, CVC: intCvc, EndDate: timeEndDate}
	jsonCard, err := json.Marshal(card)
	if err != nil {
		myErr := fmt.Errorf("error in marshal card data: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	secretKey := encryption.AesKeySecureRandom([]byte(passwordSecure))
	encryptCard, err := encryption.Encrypt(string(jsonCard), secretKey)
	if err != nil {
		myErr := fmt.Errorf("error in encrypt card data: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	updatedCardEntityID, err := s.grpc.EntityUpdate(context.Background(),
		&grpc.UpdateEntityRequest{Name: name, Data: []byte(encryptCard), Type: vars.Card.ToString(),
			AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID, CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on update card: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("card with ID %s updated", updatedCardEntityID)
	return nil
}
