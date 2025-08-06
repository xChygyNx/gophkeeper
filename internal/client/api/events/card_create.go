package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xChygyNx/gophkeeper/internal/client/consts"
	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/vars"
)

// CardCreate - create card
func (s Event) CardCreate(name, description, password, paymentSystem, number, holder, cvc, endDate string, token model.Token) error {
	s.logger.Info("card create ")

	validateData := &ValidateCardData{
		Cvc:         cvc,
		timeEndDate: endDate,
	}
	checker := NewCardChecker(s.logger, consts.DateFormat, validateData)
	err := checker.RunChecks()
	if err != nil {
		myErr := fmt.Errorf("error in validate data for create card: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	card := model.Card{
		Name:          name,
		Description:   description,
		PaymentSystem: paymentSystem,
		Number:        number,
		Holder:        holder,
		EndDate:       checker.result.endDate,
		CVC:           checker.result.intCvc,
	}
	jsonCard, err := json.Marshal(card)
	if err != nil {
		myErr := fmt.Errorf("error in marshal card data: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	secretKey := encryption.AesKeySecureRandom([]byte(password))
	encryptCard, err := encryption.Encrypt(string(jsonCard), secretKey)
	if err != nil {
		myErr := fmt.Errorf("error in encrypt card data: %w", err)
		s.logger.Error(myErr)
		return myErr
	}
	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	metadata := model.MetadataEntity{Name: name, Description: description, Type: vars.Card.ToString()}
	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("error in marshal card metadata: %w", err)
	}
	createdEntityID, err := s.grpc.EntityCreate(context.Background(),
		&grpc.CreateEntityRequest{Data: []byte(encryptCard), Metadata: string(jsonMetadata),
			AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID, CreatedAt: createdToken, EndDateAt: endDateToken}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on create card: %w", err)
		s.logger.Error(myErr)
		return myErr
	}

	s.logger.Debug("created card ID: %s", createdEntityID)
	return nil
}
