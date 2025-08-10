package grpchandler

import (
	"context"
	"fmt"
	grpc "github.com/xChygyNx/gophkeeper/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xChygyNx/gophkeeper/internal/server/model"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
)

// Authentication - user authentication, create access token
func (h *Handler) Authentication(ctx context.Context, req *grpc.AuthenticationRequest) (*grpc.AuthenticationResponse, error) {
	h.logger.Info("authentication")
	UserData := &model.UserRequest{
		Username: req.Username,
		Password: req.Password,
	}

	authenticatedUser, err := h.user.Authentication(UserData)
	if err != nil {
		finalError := fmt.Errorf("error in authinticate user: %w", err)
		h.logger.Error(finalError)
		return &grpc.AuthenticationResponse{}, status.Errorf(
			codes.Unauthenticated, finalError.Error(),
		)
	}
	user := model.GetUserData(authenticatedUser)

	token, err := h.token.Create(user.UserId, h.config.AccessTokenLifetime)
	if err != nil {
		finalError := fmt.Errorf("error in create user access token: %w", err)
		h.logger.Error(finalError)
		return &grpc.AuthenticationResponse{}, status.Errorf(
			codes.Internal, finalError.Error(),
		)
	}

	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)

	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	h.logger.Debug("User %v authenticate", authenticatedUser)
	return &grpc.AuthenticationResponse{AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID, CreatedAt: createdToken, EndDateAt: endDateToken}}, nil
}
