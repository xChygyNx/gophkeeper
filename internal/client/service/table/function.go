package table

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/xChygyNx/gophkeeper/internal/client/consts"
	"github.com/xChygyNx/gophkeeper/internal/client/model"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
)

func SearchByColumn(slice [][]string, targetColumn int, targetValue string) bool {
	for i := 1; i < len(slice) && len(slice) > 1; i++ {
		if slice[i][targetColumn] == targetValue {
			return true
		}
	}
	return false
}

func RemoveRow(slice [][]string, indexRow int) [][]string {
	return append(slice[:indexRow], slice[indexRow+1:]...)
}

// ------------------------------------------------- append

func AppendTextEntity(node *grpc.Entity, dataTblText *[][]string, plaintext string) error {
	created, err := service.ConvertTimestampToTime(node.CreatedAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		return myErr
	}
	updated, err := service.ConvertTimestampToTime(node.UpdatedAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		return myErr
	}
	var metadata model.MetadataEntity
	err = json.Unmarshal([]byte(node.Metadata), &metadata)
	if err != nil {
		return err
	}
	row := []string{metadata.Name, metadata.Description, plaintext, created.Format(
		consts.DateAndTimeFormat), updated.Format(consts.DateAndTimeFormat)}
	*dataTblText = append(*dataTblText, row)
	return nil
}

func AppendLoginPasswordEntity(node *grpc.Entity, dataTblLoginPassword *[][]string, jsonLoginPassword model.LoginPassword) error {
	created, err := service.ConvertTimestampToTime(node.CreatedAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		return myErr
	}
	updated, err := service.ConvertTimestampToTime(node.UpdatedAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		return myErr
	}
	var metadata model.MetadataEntity
	err = json.Unmarshal([]byte(node.Metadata), &metadata)
	if err != nil {
		return err
	}
	row := []string{metadata.Name, metadata.Description, jsonLoginPassword.Login, jsonLoginPassword.Password,
		created.Format(consts.DateAndTimeFormat), updated.Format(consts.DateAndTimeFormat)}
	*dataTblLoginPassword = append(*dataTblLoginPassword, row)
	return nil
}

func AppendCardEntity(node *grpc.Entity, dataTblCard *[][]string, jsonCard model.Card) error {
	created, err := service.ConvertTimestampToTime(node.CreatedAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		return myErr
	}
	updated, err := service.ConvertTimestampToTime(node.UpdatedAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		return myErr
	}
	var metadata model.MetadataEntity
	err = json.Unmarshal([]byte(node.Metadata), &metadata)
	if err != nil {
		return err
	}
	row := []string{metadata.Name, metadata.Description, jsonCard.PaymentSystem, jsonCard.Number,
		jsonCard.Holder, strconv.Itoa(jsonCard.CVC), jsonCard.EndDate.Format(consts.DateFormat),
		created.Format(consts.DateAndTimeFormat), updated.Format(consts.DateAndTimeFormat)}
	*dataTblCard = append(*dataTblCard, row)
	return nil
}

func AppendBinary(node *grpc.Binary, dataTblBinary *[][]string) error {
	created, err := service.ConvertTimestampToTime(node.CreatedAt)
	if err != nil {
		myErr := fmt.Errorf("error in convert grpc time format: %w", err)
		return myErr
	}
	row := []string{node.Name, created.Format(consts.DateAndTimeFormat)}
	*dataTblBinary = append(*dataTblBinary, row)
	return nil
}

// ------------------------------------------------- update

func UpdateRowLoginPassword(login, password string, slice [][]string, indexRow int) [][]string {
	const indexColLogin = 2
	const indexColPassword = 3
	const indexColUpdateAt = 5
	slice[indexRow][indexColLogin] = login
	slice[indexRow][indexColPassword] = password
	slice[indexRow][indexColUpdateAt] = time.Now().Format(consts.DateAndTimeFormat)
	return slice
}

func UpdateRowText(text string, slice [][]string, indexRow int) [][]string {
	const indexColText = 2
	const indexColUpdateAt = 4
	slice[indexRow][indexColText] = text
	slice[indexRow][indexColUpdateAt] = time.Now().Format(consts.DateAndTimeFormat)
	return slice
}

func UpdateRowCard(paymentSystem, number, holder, cvc, endDate string, slice [][]string, indexRow int) [][]string {
	const indexColPaymentSystem = 2
	const indexColNumber = 3
	const indexColHolder = 4
	const indexColCvc = 5
	const indexColEndDate = 6
	const indexColUpdateAt = 8
	slice[indexRow][indexColPaymentSystem] = paymentSystem
	slice[indexRow][indexColNumber] = number
	slice[indexRow][indexColHolder] = holder
	slice[indexRow][indexColEndDate] = endDate
	slice[indexRow][indexColCvc] = cvc
	slice[indexRow][indexColUpdateAt] = time.Now().Format(consts.DateAndTimeFormat)
	return slice
}
