package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/errors"
)

// EntityDelete - check the validity of the token and delete record (text, bank card or login password)
func (h *Handler) EntityDelete(ctx context.Context, req *grpc.DeleteEntityRequest) (*grpc.DeleteEntityResponse, error) {
	h.logger.Info("delete entity")

	endDateToken, err := h.token.GetEndDateToken(req.AccessToken.Token)
	if err != nil {
		finalError := fmt.Errorf("error in get end date of token from DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.DeleteEntityResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}
	valid := h.token.Validate(endDateToken)
	if !valid {
		h.logger.Error(errors.ErrNotValidateToken)
		return &grpc.DeleteEntityResponse{}, status.Errorf(codes.Unauthenticated, errors.ErrNotValidateToken.Error())
	}

	DeletedEntityID, err := h.entity.Delete(req.AccessToken.UserId, req.Name, req.Type)
	if err != nil {
		finalError := fmt.Errorf("error in delete entity from DB: %w", err)
		h.logger.Error(finalError)
		return &grpc.DeleteEntityResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}

	h.logger.Debug("Deleted Entity ID: %d", DeletedEntityID)
	return &grpc.DeleteEntityResponse{Id: DeletedEntityID}, nil
}
