package grpchandler

import (
	"fmt"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/errors"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateAccessToken(token string, tokenRepository *token.Token) (error, error) {
	endDateToken, err := tokenRepository.GetEndDateToken(token)
	if err != nil {
		finalError := fmt.Errorf("error in get end date of token from DB: %w", err)
		return finalError, status.Errorf(codes.Internal, finalError.Error())
	}
	valid := tokenRepository.Validate(endDateToken)
	if !valid {
		return errors.ErrNotValidateToken, status.Errorf(codes.Unauthenticated, errors.ErrNotValidateToken.Error())
	}
	return nil, nil
}
