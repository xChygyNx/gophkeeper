package grpchandler

import (
	"context"
	"fmt"
	grpc "github.com/xChygyNx/gophkeeper/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EntityUpdate - checks the validity of the token and update record (text, bank card or login password)
func (h *Handler) EntityUpdate(ctx context.Context, req *grpc.UpdateEntityRequest) (*grpc.UpdateEntityResponse, error) {
	h.logger.Info("entity update")

	logError, statusErr := validateAccessToken(req.AccessToken.Token, h.token)
	if statusErr != nil {
		h.logger.Error(logError)
		return &grpc.UpdateEntityResponse{}, statusErr
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
