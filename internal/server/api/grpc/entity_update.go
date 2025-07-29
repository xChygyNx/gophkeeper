package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/errors"
)

// EntityUpdate - checks the validity of the token and update record (text, bank card or login password)
func (h *Handler) EntityUpdate(ctx context.Context, req *grpc.UpdateEntityRequest) (*grpc.UpdateEntityResponse, error) {
	h.logger.Info("entity update")

	endDateToken, err := h.token.GetEndDateToken(req.AccessToken.Token)
	if err != nil {
		finalError := fmt.Errorf("error in get end date of token from DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.UpdateEntityResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}
	valid := h.token.Validate(endDateToken)
	if !valid {
		h.logger.Error(errors.ErrNotValidateToken)
		return &grpc.UpdateEntityResponse{}, status.Errorf(codes.Unauthenticated, errors.ErrNotValidateToken.Error())
	}

	UpdatedEntityID, err := h.entity.Update(req.AccessToken.UserId, req.Name, req.Type, req.Data)
	if err != nil {
		finalError := fmt.Errorf("error in update entity record in DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.UpdateEntityResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}

	h.logger.Debug("Update entity ID: %s", UpdatedEntityID)
	return &grpc.UpdateEntityResponse{Id: UpdatedEntityID}, nil
}
