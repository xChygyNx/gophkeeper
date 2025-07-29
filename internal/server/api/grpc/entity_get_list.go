package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xChygyNx/gophkeeper/internal/server/model"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/errors"
)

// EntityGetList - checks the validity of token and get list records (text, bank card or login-password)
func (h *Handler) EntityGetList(ctx context.Context, req *grpc.GetListEntityRequest) (*grpc.GetListEntityResponse, error) {
	h.logger.Info("Get list entity")

	endDateToken, err := h.token.GetEndDateToken(req.AccessToken.Token)
	if err != nil {
		finalError := fmt.Errorf("error in get end date of token from DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.GetListEntityResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}
	valid := h.token.Validate(endDateToken)
	if !valid {
		h.logger.Error(errors.ErrNotValidateToken)
		return &grpc.GetListEntityResponse{}, status.Errorf(codes.Unauthenticated, errors.ErrNotValidateToken.Error())
	}

	ListEntity, err := h.entity.GetList(req.AccessToken.UserId, req.Type)
	if err != nil {
		finalError := fmt.Errorf("error in get list of entity from DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.GetListEntityResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}
	list, err := model.GetListEntity(ListEntity)
	if err != nil {
		finalError := fmt.Errorf("error in transform list of entity in grpc format: %w", err)
		h.logger.Error(finalError)
		return &grpc.GetListEntityResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}

	h.logger.Debug("Entity list: %s", list)
	return &grpc.GetListEntityResponse{Node: list}, nil
}
