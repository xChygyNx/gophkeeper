package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xChygyNx/gophkeeper/internal/server/model"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
)

// EntityGetList - checks the validity of token and get list records (text, bank card or login-password)
func (h *Handler) EntityGetList(ctx context.Context, req *grpc.GetListEntityRequest) (*grpc.GetListEntityResponse, error) {
	h.logger.Info("Get list entity")

	logError, statusErr := validateAccessToken(req.AccessToken.Token, h.token)
	if statusErr != nil {
		h.logger.Error(logError)
		return &grpc.GetListEntityResponse{}, statusErr
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
