package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
)

// EntityDelete - check the validity of the token and delete record (text, bank card or login password)
func (h *Handler) EntityDelete(ctx context.Context, req *grpc.DeleteEntityRequest) (*grpc.DeleteEntityResponse, error) {
	h.logger.Info("delete entity")

	logError, statusErr := validateAccessToken(req.AccessToken.Token, h.token)
	if statusErr != nil {
		h.logger.Error(logError)
		return &grpc.DeleteEntityResponse{}, statusErr
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
