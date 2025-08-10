package events

import (
	"encoding/json"
	"fmt"
	grpc "github.com/xChygyNx/gophkeeper/internal/proto"

	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	"github.com/xChygyNx/gophkeeper/internal/client/service/table"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/labels"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/vars"
)

// Synchronization - run synchronization
func (s Event) Synchronization(password string, token model.Token) ([][]string, [][]string, [][]string, [][]string, error) {
	s.logger.Info("synchronization")

	dataTblText := make([][]string, 0)
	dataTblCard := make([][]string, 0)
	dataTblLoginPassword := make([][]string, 0)
	dataTblBinary := make([][]string, 0)

	created := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDate := service.ConvertTimeToTimestamp(token.EndDateAt)
	//-----------------------------------------------
	var plaintext string
	secretKey := encryption.AesKeySecureRandom([]byte(password))

	titleText := []string{labels.NameItem, labels.DescriptionItem, labels.DataItem, labels.CreatedAtItem, labels.UpdatedAtItem}
	titleCard := []string{labels.NameItem, labels.DescriptionItem, labels.PaymentSystemItem, labels.NumberItem, labels.HolderItem,
		labels.CVCItem, labels.EndDateItem, labels.CreatedAtItem, labels.UpdatedAtItem}
	titleLoginPassword := []string{labels.NameItem, labels.DescriptionItem, labels.LoginItem, labels.PasswordItem,
		labels.CreatedAtItem, labels.UpdatedAtItem}
	titleBinary := []string{labels.NameItem, labels.CreatedAtItem}

	dataTblText = append(dataTblText, titleText)
	dataTblCard = append(dataTblCard, titleCard)
	dataTblLoginPassword = append(dataTblLoginPassword, titleLoginPassword)
	dataTblBinary = append(dataTblBinary, titleBinary)

	dataTblTextPointer := &dataTblText
	dataTblCardPointer := &dataTblCard
	dataTblLoginPasswordPointer := &dataTblLoginPassword
	dataTblBinaryPointer := &dataTblBinary

	//-----------------------------------------------
	nodesTextEntity, err := s.grpc.EntityGetList(s.context,
		&grpc.GetListEntityRequest{Type: vars.Text.ToString(),
			AccessToken: &grpc.Token{Token: token.AccessToken,
				UserId: token.UserID, CreatedAt: created, EndDateAt: endDate}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on get list of textes: %w", err)
		s.logger.Error(myErr)
		return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
	}

	for _, node := range nodesTextEntity.Node {
		plaintext, err = encryption.Decrypt(string(node.Data), secretKey)
		if err != nil {
			myErr := fmt.Errorf("error in decrypt text: %w", err)
			s.logger.Error(myErr)
			return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
		}
		err = table.AppendTextEntity(node, dataTblTextPointer, plaintext)
		if err != nil {
			myErr := fmt.Errorf("error in add text in table: %w", err)
			s.logger.Error(myErr)
			return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
		}
	}

	//-----------------------------------------------
	nodesCardEntity, err := s.grpc.EntityGetList(s.context,
		&grpc.GetListEntityRequest{Type: vars.Card.ToString(),
			AccessToken: &grpc.Token{Token: token.AccessToken,
				UserId: token.UserID, CreatedAt: created, EndDateAt: endDate}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on get list of cards: %w", err)
		s.logger.Error(myErr)
		return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
	}

	for _, node := range nodesCardEntity.Node {
		plaintext, err = encryption.Decrypt(string(node.Data), secretKey)
		if err != nil {
			myErr := fmt.Errorf("error in decrypt card data: %w", err)
			s.logger.Error(myErr)
			return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
		}

		var card model.Card
		err = json.Unmarshal([]byte(plaintext), &card)
		if err != nil {
			myErr := fmt.Errorf("error in unmarshal raw card data: %w", err)
			s.logger.Error(myErr)
			return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
		}
		err = table.AppendCardEntity(node, dataTblCardPointer, card)
		if err != nil {
			myErr := fmt.Errorf("error in add card data in table: %w", err)
			s.logger.Error(myErr)
			return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
		}
	}
	//-----------------------------------------------
	nodesLoginPasswordEntity, err := s.grpc.EntityGetList(s.context,
		&grpc.GetListEntityRequest{Type: vars.LoginPassword.ToString(),
			AccessToken: &grpc.Token{Token: token.AccessToken,
				UserId: token.UserID, CreatedAt: created, EndDateAt: endDate}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on get list of login/password: %w", err)
		s.logger.Error(myErr)
		return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
	}

	for _, node := range nodesLoginPasswordEntity.Node {
		plaintext, err = encryption.Decrypt(string(node.Data), secretKey)
		if err != nil {
			myErr := fmt.Errorf("error in decrypt login/password data: %w", err)
			s.logger.Error(myErr)
			return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
		}

		var loginPassword model.LoginPassword
		err = json.Unmarshal([]byte(plaintext), &loginPassword)
		if err != nil {
			myErr := fmt.Errorf("error in unmarshal raw login/password data: %w", err)
			s.logger.Error(myErr)
			return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
		}
		err = table.AppendLoginPasswordEntity(node, dataTblLoginPasswordPointer, loginPassword)
		if err != nil {
			myErr := fmt.Errorf("error in add card login/password in table: %w", err)
			s.logger.Error(myErr)
			return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
		}
	}
	//-----------------------------------------------
	nodesBinary, err := s.grpc.FileGetList(s.context,
		&grpc.GetListBinaryRequest{AccessToken: &grpc.Token{Token: token.AccessToken,
			UserId: token.UserID, CreatedAt: created, EndDateAt: endDate}})
	if err != nil {
		myErr := fmt.Errorf("error in grpc request on get list of bynaries: %w", err)
		s.logger.Error(myErr)
		return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
	}

	for _, node := range nodesBinary.Node {
		err = table.AppendBinary(node, dataTblBinaryPointer)
		if err != nil {
			myErr := fmt.Errorf("error in add binary in table: %w", err)
			s.logger.Error(myErr)
			return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, myErr
		}
	}
	//-----------------------------------------------
	return dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, nil
}
