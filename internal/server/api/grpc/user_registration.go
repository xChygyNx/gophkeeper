package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xChygyNx/gophkeeper/internal/server/model"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/errors"
)

// Registration - registration new user, create access token
func (h *Handler) Registration(ctx context.Context, req *grpc.RegistrationRequest) (*grpc.RegistrationResponse, error) {
	h.logger.Info("registration")

	UserData := &model.UserRequest{}
	UserData.Username = req.Username
	UserData.Password = req.Password

	exists, err := h.user.UserExists(UserData.Username)
	if err != nil {
		finalError := fmt.Errorf("error in check user exists: %w", err)
		h.logger.Error(finalError)
		return &grpc.RegistrationResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}
	if exists {
		err = errors.ErrUsernameAlreadyExists
		h.logger.Error(err)
		return &grpc.RegistrationResponse{}, status.Errorf(
			codes.AlreadyExists, err.Error(),
		)
	}
	registeredUser, err := h.user.Registration(UserData)
	if err != nil {
		h.logger.Error(fmt.Errorf("error in insert record about user in DB: %w", err))
		return &grpc.RegistrationResponse{}, status.Errorf(codes.Internal, err.Error())
	}
	user := model.GetUserData(registeredUser)

	token, err := h.token.Create(user.UserId, h.config.AccessTokenLifetime)
	if err != nil {
		finalError := fmt.Errorf("error in create user access token: %w", err)
		h.logger.Error(finalError)
		return &grpc.RegistrationResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}

	createdToken := service.ConvertTimeToTimestamp(token.CreatedAt)
	endDateToken := service.ConvertTimeToTimestamp(token.EndDateAt)

	err = service.CreateStorageUser(h.config.FileFolder, token.UserID)
	if err != nil {
		finalError := fmt.Errorf("error in create user directory: %w", err)
		h.logger.Error(finalError)
		return &grpc.RegistrationResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}

	h.logger.Debug(registeredUser)
	return &grpc.RegistrationResponse{AccessToken: &grpc.Token{Token: token.AccessToken, UserId: token.UserID,
		CreatedAt: createdToken, EndDateAt: endDateToken}}, nil
}
